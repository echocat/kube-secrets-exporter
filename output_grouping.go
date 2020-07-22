package kube_secrets_exporter

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"text/template"
)

type OutputGrouping struct {
	projection      *template.Template
	plainProjection string
	matcher         *regexp.Regexp
	plainMatcher    string
	file            *template.Template
	plainFile       string
}

func (instance *OutputGrouping) Set(plain string) (err error) {
	plain = strings.TrimSpace(plain)
	if plain == "" {
		*instance = OutputGrouping{}
		return nil
	}
	parts := strings.Split(plain, "=")

	var result OutputGrouping
	if len(parts) == 3 {
		result.plainProjection = parts[0]
		result.plainMatcher = parts[1]
		result.plainFile = parts[2]
	} else if len(parts) == 1 {
		result.plainProjection = ""
		result.plainMatcher = ".*"
		result.plainFile = parts[0]
	} else {
		return fmt.Errorf("illegal output grouping, only 1 or 3 segmets are allowed")
	}

	result.projection, err = parseOutputGroupingTemplate("projection", result.plainProjection)
	if err != nil {
		return fmt.Errorf("illegal projection of output grouping: %w", err)
	}

	result.matcher, err = regexp.Compile(result.plainMatcher)
	if err != nil {
		return fmt.Errorf("illegal matcher of output grouping: %w", err)
	}

	result.file, err = parseOutputGroupingTemplate("file", result.plainFile)
	if err != nil {
		return fmt.Errorf("illegal file of output grouping: %w", err)
	}

	*instance = result
	return nil
}

func (instance OutputGrouping) String() string {
	if instance.plainProjection == "" && instance.plainMatcher == ".*" {
		return instance.plainFile
	}
	return fmt.Sprintf("%s=%s=%s", instance.plainProjection, instance.plainMatcher, instance.plainFile)
}

func (instance OutputGrouping) Project(context interface{}) (*string, error) {
	t := instance.projection
	if t == nil {
		return nil, nil
	}
	buf := new(bytes.Buffer)
	if err := t.Execute(buf, context); err != nil {
		return nil, fmt.Errorf("cannot project %v using %v: %w", context, t, err)
	}
	projected := buf.String()
	return &projected, nil
}

func (instance OutputGrouping) fileNamePatternFor(context interface{}) (string, error) {
	t := instance.file
	if t == nil {
		return Stdout.String(), nil
	}
	buf := new(bytes.Buffer)
	if err := t.Execute(buf, context); err != nil {
		return "", fmt.Errorf("cannot evaluate file name for %v using %v: %w", context, t, err)
	}
	return buf.String(), nil
}

func (instance OutputGrouping) Apply(context interface{}) (File, bool, error) {
	projected, err := instance.Project(context)
	if err != nil {
		return "", false, err
	}

	fileName, err := instance.fileNamePatternFor(context)
	if err != nil {
		return "", false, err
	}

	matcher := instance.matcher
	if matcher != nil && projected != nil {
		if !matcher.MatchString(*projected) {
			return "", false, nil
		}
		fileName = instance.matcher.ReplaceAllString(*projected, fileName)
	}

	var f File
	if err := f.Set(fileName); err != nil {
		return "", false, err
	}

	return f, true, nil
}

type OutputGroupings []OutputGrouping

func (instance *OutputGroupings) Set(plain string) error {
	var result OutputGroupings
	for _, candidate := range strings.Split(plain, ",") {
		var v OutputGrouping
		if err := v.Set(candidate); err != nil {
			return err
		}
		result = append(result, v)
	}
	*instance = result
	return nil
}

func (instance OutputGroupings) String() string {
	return strings.Join(instance.Strings(), ",")
}

func (instance OutputGroupings) Strings() []string {
	strs := make([]string, len(instance))
	for i, v := range instance {
		strs[i] = v.String()
	}
	return strs
}

func (instance OutputGroupings) Apply(context interface{}) (File, error) {
	for _, candidate := range instance {
		if match, ok, err := candidate.Apply(context); err != nil {
			return "", err
		} else if ok {
			return match, nil
		}
	}
	return Stdout, nil
}

func parseOutputGroupingTemplate(name, pattern string) (*template.Template, error) {
	return template.New(name).Parse(pattern)
}
