package kube_secrets_exporter

import (
	"fmt"
	"strings"
)

type OutputBundling uint8

const (
	OutputBundlingList       = OutputBundling(0)
	OutputBundlingSeparation = OutputBundling(1)
)

func (instance *OutputBundling) Set(plain string) error {
	if v, ok := nameToOutputBundling[strings.ToLower(plain)]; ok {
		*instance = v
		return nil
	}
	return fmt.Errorf("illegal output bundling: %s", plain)
}

func (instance OutputBundling) String() string {
	if v, ok := outputBundlingToName[instance]; ok {
		return v
	}
	return fmt.Sprintf("illegal output bundling: %d", instance)
}

type OutputBundlings []OutputBundling

func (instance OutputBundlings) String() string {
	return strings.Join(instance.Strings(), ",")
}

func (instance OutputBundlings) Strings() []string {
	strs := make([]string, len(instance))
	for i, v := range instance {
		strs[i] = v.String()
	}
	return strs
}

var (
	outputBundlingToName = map[OutputBundling]string{
		OutputBundlingList:       "list",
		OutputBundlingSeparation: "separation",
	}

	nameToOutputBundling = func(in map[OutputBundling]string) map[string]OutputBundling {
		result := make(map[string]OutputBundling)
		for f, n := range in {
			result[n] = f
		}
		return result
	}(outputBundlingToName)

	AllOutputBundlings = func(in map[OutputBundling]string) OutputBundlings {
		result := make(OutputBundlings, len(in))
		var i int
		for f := range in {
			result[i] = f
			i++
		}
		return result
	}(outputBundlingToName)
)
