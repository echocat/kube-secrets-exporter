package main

import (
	"github.com/blaubaer/kingpin"
	kse "github.com/echocat/kube-secrets-exporter"
)

func main() {
	kingpin.CommandLine.EnvarNamePrefix("KSE_")
	kingpin.CmdsOf(&kse.KubeSecretsExporter{})
	kingpin.Parse()
}
