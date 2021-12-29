package notifybot

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
)

const (
	slash = "/"
	com   = "."
)

type Templater interface {
	Do(date interface{}) (string, error)
}

type templates struct {
	templates map[string]*template.Template
	content   map[string]string
	asset     assetFunc
	assetName assetNameFunc
}

type assetFunc func(name string) ([]byte, error)
type assetNameFunc func() []string

func buildTemplates(asset assetFunc, assetName assetNameFunc) (Templater, error) {

	t := &templates{
		templates: make(map[string]*template.Template),
		content:   make(map[string]string),
		asset:     asset,
		assetName: assetName,
	}

	names := t.assetName()
	for _, name := range names {
		b, err := t.asset(name)
		if err != nil {
			return nil, err
		}

		i := strings.Index(name, slash)
		j := strings.Index(name, com)
		service := name[i+1 : j]

		temple, err := template.New(service).Parse(string(b))
		if err != nil {
			return nil, err
		}

		t.templates[service] = temple
		t.content[name] = string(b)
	}

	return t, nil
}

func (t *templates) Do(data interface{}) (string, error) {
	if stringer, ok := data.(fmt.Stringer); ok {
		name := stringer.String()
		buff := &bytes.Buffer{}
		temple, ok := t.templates[name]
		if !ok {
			return "", fmt.Errorf("template %s not found", name)
		}

		err := temple.Execute(buff, data)
		if err != nil {
			return "", err
		}

		return buff.String(), nil
	}

	return "", fmt.Errorf("data %v not implement fmt.Stringer", data)
}
