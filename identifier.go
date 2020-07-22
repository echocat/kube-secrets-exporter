package kube_secrets_exporter

import (
	"fmt"
	"regexp"
	"strings"
)

type Identifier string

func (instance Identifier) String() string {
	return string(instance)
}

type IdentifierPattern struct {
	plain  string
	regexp *regexp.Regexp
}

func (instance *IdentifierPattern) Set(plain string) error {
	if plain == "" {
		instance.regexp = nil
		instance.plain = ""
		return nil
	}
	if r, err := regexp.Compile("^" + plain + "$"); err != nil {
		return fmt.Errorf("illegal identifier pattern: %w", err)
	} else {
		instance.plain = plain
		instance.regexp = r
		return nil
	}
}

func (instance IdentifierPattern) String() string {
	return instance.plain
}

func (instance IdentifierPattern) Matches(candidate Identifier) bool {
	if r := instance.regexp; r != nil {
		return r.MatchString(candidate.String())
	}
	return candidate == ""
}

type IdentifierPatterns []IdentifierPattern

func (instance IdentifierPatterns) String() string {
	return strings.Join(instance.Strings(), ",")
}

func (instance IdentifierPatterns) Strings() []string {
	strs := make([]string, len(instance))
	for i, candidate := range instance {
		strs[i] = candidate.String()
	}
	return strs
}

func (instance IdentifierPatterns) Matches(candidate Identifier) bool {
	for _, v := range instance {
		if v.Matches(candidate) {
			return true
		}
	}
	return false
}

type IdentifierMatcher struct {
	Includes IdentifierPatterns
	Excludes IdentifierPatterns
}

func (instance IdentifierMatcher) Matches(candidate Identifier) bool {
	if v := instance.Includes; len(v) > 0 && !v.Matches(candidate) {
		return false
	}
	if v := instance.Excludes; len(v) > 0 && v.Matches(candidate) {
		return false
	}
	return true
}

func (instance *IdentifierMatcher) Set(plain string) error {
	var result IdentifierMatcher
	for _, part := range strings.Split(plain, ",") {
		var v IdentifierPattern
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

func (instance IdentifierMatcher) String() string {
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
