package utils

import (
	"bytes"
	"text/template"

	"github.com/pkg/errors"
)

// RenderTemplateToBytes parses text as a template body for a new text/template
// with the given name and executes it with the given data object, returning the
// output as a byte slice.
func RenderTemplateToBytes(name, text string, data interface{}) ([]byte, error) {
	var bytesBuffer bytes.Buffer
	tmpl, err := template.New(name).Parse(text)
	if err != nil {
		return nil, errors.Wrap(err, "RenderGoTemplateToString: template.Parse failed")
	}
	err = tmpl.Execute(&bytesBuffer, data)
	if err != nil {
		return nil, errors.Wrap(err, "RenderGoTemplateToString: template.Execute failed")
	}
	return bytesBuffer.Bytes(), nil
}
