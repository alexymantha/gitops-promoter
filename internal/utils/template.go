package utils

import (
	"bytes"
	"text/template"

	sprig "github.com/go-task/slim-sprig/v3"
)

var (
	SanitizedSprigFuncMap = sprig.GenericFuncMap()
)

func init() {
	delete(SanitizedSprigFuncMap, "env")
	delete(SanitizedSprigFuncMap, "expandenv")
	delete(SanitizedSprigFuncMap, "getHostByName")
}

func RenderStringTemplate(templateStr string, data any) (string, error) {
	tmpl, err := template.New("").Funcs(SanitizedSprigFuncMap).Parse(templateStr)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}
