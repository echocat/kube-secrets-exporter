package kube_secrets_exporter

import (
	"github.com/blaubaer/kingpin"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
)

type Filter struct {
	RemoveCreationTimestamp bool
	RemoveResourceVersion   bool
	RemoveSelfLink          bool
	RemoveUid               bool
	RemoveKubectlHints      bool
}

func (instance *Filter) RegisterFlags(fg kingpin.FlagGroup) {
	g := fg.FlagGroup("filter.").
		EnvarNamePrefix("FILTER_")

	g.Flag("remove-creation-timestamp", "").
		Default("true").
		BoolVar(&instance.RemoveCreationTimestamp)
	g.Flag("remove-resource-version", "").
		Default("true").
		BoolVar(&instance.RemoveResourceVersion)
	g.Flag("remove-self-link", "").
		Default("true").
		BoolVar(&instance.RemoveSelfLink)
	g.Flag("remove-uid", "").
		Default("true").
		BoolVar(&instance.RemoveUid)
	g.Flag("remove-kubectl-hints", "").
		Default("true").
		BoolVar(&instance.RemoveKubectlHints)
}

func (instance Filter) Apply(secret *v1.Secret) error {
	if instance.RemoveCreationTimestamp {
		secret.ObjectMeta.CreationTimestamp = metav1.Now()
	}
	if instance.RemoveResourceVersion {
		secret.ObjectMeta.ResourceVersion = ""
	}
	if instance.RemoveSelfLink {
		secret.ObjectMeta.SelfLink = ""
	}
	if instance.RemoveUid {
		secret.ObjectMeta.UID = ""
	}
	if instance.RemoveKubectlHints {
		for n := range secret.Annotations {
			if strings.HasPrefix(n, "kubectl.kubernetes.io/") {
				delete(secret.Annotations, n)
			}
		}
	}
	return nil
}
