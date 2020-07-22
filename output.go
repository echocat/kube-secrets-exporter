package kube_secrets_exporter

import (
	"fmt"
	"github.com/blaubaer/kingpin"
)

type Output struct {
	File   OutputGroupings
	Format OutputFormat
}

func (instance *Output) RegisterFlags(fg kingpin.FlagGroup) {
	g := fg.FlagGroup("output.").
		EnvarNamePrefix("OUTPUT_")

	g.Flag("file", "Where to write the output to. It can be a regular file, 'stdout' or 'stderr'."+
		"This can result in grouping of files, too by using golang template evaluation"+
		"; example '{{.Namespace}}.yaml' will create an extra file for each element per namespace.").
		Default(Stdout.String()).
		Envar("FILE").
		SetValue(&instance.File)
	g.Flag("format", fmt.Sprintf("Which format should be used for output. Can be: %v", AllOutputFormats.String())).
		Default(OutputFormatYaml.String()).
		Envar("FORMAT").
		SetValue(&instance.Format)
}
