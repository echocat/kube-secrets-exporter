package kube_secrets_exporter

import (
	"fmt"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/client-go/kubernetes/scheme"
	"sync"
)

type OutputConsumer struct {
	Output

	groups map[File][]runtime.Object
	mutex  sync.Mutex
}

func (instance *OutputConsumer) Consume(context runtime.Object) error {
	instance.mutex.Lock()
	defer instance.mutex.Unlock()

	f, err := instance.File.Apply(context)
	if err != nil {
		return err
	}

	if instance.groups == nil {
		instance.groups = make(map[File][]runtime.Object)
	}

	group := instance.groups[f]
	group = append(group, context)
	instance.groups[f] = group

	return nil
}

func (instance *OutputConsumer) Finalize() error {
	instance.mutex.Lock()
	defer instance.mutex.Unlock()

	if groups := instance.groups; groups != nil {
		for f, values := range groups {
			if err := instance.writeGroup(f, values); err != nil {
				return err
			}
		}
	}
	return nil
}

func (instance *OutputConsumer) writeGroup(f File, values []runtime.Object) error {
	switch instance.Format {
	case OutputFormatYaml:
		return instance.writeGroupAsYaml(f, values)
	default:
		return fmt.Errorf("cannot handle output format: %v", instance.Format)
	}
}

func (instance *OutputConsumer) writeGroupAsYaml(f File, values []runtime.Object) error {
	w, err := f.Open()
	if err != nil {
		return err
	}
	defer func() { _ = w.Close() }()

	enc := json.NewSerializerWithOptions(json.DefaultMetaFactory, scheme.Scheme, scheme.Scheme, json.SerializerOptions{Yaml: true, Pretty: true})
	for i, value := range values {
		if i > 0 {
			if _, err := w.Write([]byte("---\n")); err != nil {
				return fmt.Errorf("cannot write output file %v: %w", f, err)
			}
		}
		if err := enc.Encode(value, w); err != nil {
			return fmt.Errorf("cannot write output file %v: %w", f, err)
		}
	}

	return nil
}
