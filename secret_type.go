package kube_secrets_exporter

import (
	"fmt"
	v1 "k8s.io/api/core/v1"
	"strings"
)

type SecretType v1.SecretType

const (
	SecretTypeOpaque              = SecretType(v1.SecretTypeOpaque)
	SecretTypeServiceAccountToken = SecretType(v1.SecretTypeServiceAccountToken)
	SecretTypeDockercfg           = SecretType(v1.SecretTypeDockercfg)
	SecretTypeDockerConfigJson    = SecretType(v1.SecretTypeDockerConfigJson)
	SecretTypeBasicAuth           = SecretType(v1.SecretTypeBasicAuth)
	SecretTypeSSHAuth             = SecretType(v1.SecretTypeSSHAuth)
	SecretTypeTLS                 = SecretType(v1.SecretTypeTLS)
	SecretTypeBootstrapToken      = SecretType(v1.SecretTypeBootstrapToken)
)

var AllSecretTypes = SecretTypes{
	SecretTypeOpaque,
	SecretTypeServiceAccountToken,
	SecretTypeDockercfg,
	SecretTypeDockerConfigJson,
	SecretTypeBasicAuth,
	SecretTypeSSHAuth,
	SecretTypeTLS,
	SecretTypeBootstrapToken,
}
var validSecretTypes = func(in []SecretType) map[SecretType]bool {
	result := make(map[SecretType]bool)
	for _, v := range in {
		result[v] = true
	}
	return result
}(AllSecretTypes)

func (instance *SecretType) Set(plain string) error {
	candidate := SecretType(plain)
	if ok := validSecretTypes[candidate]; ok {
		*instance = candidate
		return nil
	}
	return fmt.Errorf("illegal secret type: %s", plain)
}

func (instance SecretType) String() string {
	return string(instance)
}

type SecretTypes []SecretType

func (instance SecretTypes) String() string {
	return strings.Join(instance.Strings(), ",")
}

func (instance SecretTypes) Strings() []string {
	strs := make([]string, len(instance))
	for i, candidate := range instance {
		strs[i] = candidate.String()
	}
	return strs
}

func (instance SecretTypes) Contains(secretType SecretType) bool {
	for _, candidate := range instance {
		if candidate == secretType {
			return true
		}
	}
	return len(instance) == 0
}

func (instance SecretTypes) Matches(candidate SecretType) bool {
	return instance.Contains(candidate)
}

type SecretTypeMatcher struct {
	Includes SecretTypes
	Excludes SecretTypes
}

func (instance SecretTypeMatcher) Matches(candidate SecretType) bool {
	if v := instance.Includes; len(v) > 0 && !v.Matches(candidate) {
		return false
	}
	if v := instance.Excludes; len(v) > 0 && v.Matches(candidate) {
		return false
	}
	return true
}

func (instance *SecretTypeMatcher) Set(plain string) error {
	var result SecretTypeMatcher
	for _, part := range strings.Split(plain, ",") {
		var v SecretType
		part = strings.TrimSpace(part)
		if len(part) > 0 {
			include := true
			if part[0] == '!' {
				include = false
				part = part[1:]
			}
			if err := v.Set(strings.TrimSpace(part)); err != nil {
				return err
			}
			if include {
				result.Includes = append(result.Includes, v)
			} else {
				result.Excludes = append(result.Excludes, v)
			}
		}
	}
	*instance = result
	return nil
}

func (instance SecretTypeMatcher) String() string {
	li, le := len(instance.Excludes), len(instance.Includes)
	strs := make([]string, li+le)
	for i, t := range instance.Includes {
		strs[i] = t.String()
	}
	for i, t := range instance.Excludes {
		strs[li+i] = "!" + t.String()
	}
	return strings.Join(strs, ",")
}
