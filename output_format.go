package kube_secrets_exporter

import (
	"fmt"
	"strings"
)

type OutputFormat uint8

const (
	OutputFormatYaml = OutputFormat(0)
	OutputFormatJson = OutputFormat(1)
)

func (instance *OutputFormat) Set(plain string) error {
	if v, ok := nameToOutputFormat[strings.ToLower(plain)]; ok {
		*instance = v
		return nil
	}
	return fmt.Errorf("illegal output format: %s", plain)
}

func (instance OutputFormat) String() string {
	if v, ok := outputFormatToName[instance]; ok {
		return v
	}
	return fmt.Sprintf("illegal output format: %d", instance)
}

type OutputFormats []OutputFormat

func (instance OutputFormats) String() string {
	return strings.Join(instance.Strings(), ",")
}

func (instance OutputFormats) Strings() []string {
	strs := make([]string, len(instance))
	for i, v := range instance {
		strs[i] = v.String()
	}
	return strs
}

var (
	outputFormatToName = map[OutputFormat]string{
		OutputFormatYaml: "yaml",
		OutputFormatJson: "json",
	}

	nameToOutputFormat = func(in map[OutputFormat]string) map[string]OutputFormat {
		result := make(map[string]OutputFormat)
		for f, n := range in {
			result[n] = f
		}
		return result
	}(outputFormatToName)

	AllOutputFormats = func(in map[OutputFormat]string) OutputFormats {
		result := make(OutputFormats, len(in))
		var i int
		for f := range in {
			result[i] = f
			i++
		}
		return result
	}(outputFormatToName)
)
