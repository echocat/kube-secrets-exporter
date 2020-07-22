package kube_secrets_exporter

import (
	"github.com/blaubaer/kingpin"
	v1 "k8s.io/api/core/v1"
)

type Selector struct {
	Type SecretTypeMatcher
	Name IdentifierMatcher
}

func (instance *Selector) RegisterFlags(fg kingpin.FlagGroup) {
	g := fg.FlagGroup("selector.").
		EnvarNamePrefix("SELECTOR_")

	g.Flag("type", "Which types should be exported. Empty means all.").
		Envar("TYPE").
		HintAction(AllSecretTypes.Strings).
		SetValue(&instance.Type)
	g.Flag("name", "Which names should be exported. This has to match '<namespace>/<name>' of the secret. Empty means all.").
		Envar("NAME").
		SetValue(&instance.Name)
}

func (instance Selector) Matches(secret v1.Secret) bool {
	return instance.Type.Matches(SecretType(secret.Type)) &&
		instance.Name.Matches(Identifier(secret.Namespace+"/"+secret.Name))
}
