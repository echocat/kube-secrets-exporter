package kube_secrets_exporter

import (
	"fmt"
	"github.com/blaubaer/kingpin"
	"github.com/echocat/kube-secrets-exporter/kubernetes"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type KubeSecretsExporter struct {
	Environment kubernetes.Environment
	Selector    Selector
	Filter      Filter
	Output      Output

	PageSize uint32
}

func (instance *KubeSecretsExporter) RegisterCommands(fe kingpin.Cmd) {
	fe.Flag("page-size", "How many secrets should be queried at once.").
		Default("100").
		Envar("PAGE_SIZE").
		Uint32Var(&instance.PageSize)
	fe.AddAction(func(*kingpin.ParseContext) error {
		return instance.Export()
	})
	fe.RegisterFlagsOf(
		&instance.Environment,
		&instance.Selector,
		&instance.Filter,
		&instance.Output,
	)
}

func (instance *KubeSecretsExporter) Export() error {
	client, err := instance.Environment.NewClient()
	if err != nil {
		return fmt.Errorf("cannot create kubernetes client: %w", err)
	}
	consumer := OutputConsumer{
		Output: instance.Output,
	}
	opts := metav1.ListOptions{
		Limit: int64(instance.PageSize),
	}
	secrets := client.CoreV1().Secrets(instance.Environment.Namespace)
	for {
		resp, err := secrets.List(opts)
		if err != nil {
			return fmt.Errorf("cannot retrieve secrets from kubernetes: %w", err)
		}
		for _, elem := range resp.Items {
			if err := instance.onElement(elem, &consumer); err != nil {
				return fmt.Errorf("cannot handle secret %s/%s: %w", elem.Namespace, elem.Name, err)
			}
		}
		if v := resp.Continue; v != "" {
			opts.Continue = v
		} else {
			break
		}
	}
	return consumer.Finalize()
}

func (instance *KubeSecretsExporter) onElement(secret v1.Secret, consumer *OutputConsumer) error {
	if instance.Selector.Matches(secret) {
		if err := instance.Filter.Apply(&secret); err != nil {
			return err
		}
		if err := consumer.Consume(&secret); err != nil {
			return err
		}
	}
	return nil
}
