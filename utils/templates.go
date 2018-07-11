package utils

import (
	"bytes"
	"html/template"

	"github.com/pkg/errors"
)

func RenderTemplateToBytes(name, str string, data interface{}) ([]byte, error) {
	var bytesBuffer bytes.Buffer
	tmpl, err := template.New(name).Parse(str)
	if err != nil {
		return nil, errors.Wrap(err, "RenderGoTemplateToString: template.Parse failed")
	}
	err = tmpl.Execute(&bytesBuffer, data)
	if err != nil {
		return nil, errors.Wrap(err, "RenderGoTemplateToString: template.Execute failed")
	}
	return bytesBuffer.Bytes(), nil
}
