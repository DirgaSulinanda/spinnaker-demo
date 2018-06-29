package template

import (
	"bytes"
	"fmt"
	htmltemplate "html/template"
	"io"
	"net/http"
	"path/filepath"
)

// Template wraps template
type Template struct {
	baseTemplate string
	templatePath string

	funcs    htmltemplate.FuncMap
	getAsset func(string) ([]byte, error)
}

// New returns a new Template.
func New(
	baseTemplate string,
	templatePath string,
	getAsset func(string) ([]byte, error),
) *Template {
	return &Template{
		baseTemplate: baseTemplate,
		templatePath: templatePath,

		funcs:    htmltemplate.FuncMap{},
		getAsset: getAsset,
	}
}

// Funcs adds the functions in fm to the Template's function map.
// Existing functions will be overwritten in case of conflict.
func (t *Template) Funcs(fm htmltemplate.FuncMap) {
	for k, v := range fm {
		t.funcs[k] = v
	}
}

// getTemplate gets template based on base template.
func (t *Template) getTemplate(name string) (string, error) {
	baseTmpl, err := t.getAsset(t.baseTemplate)
	if err != nil {
		return "", fmt.Errorf("error reading base template: %s", err)
	}

	pageTmpl, err := t.getAsset(filepath.Join(t.templatePath, name))
	if err != nil {
		return "", fmt.Errorf("error reading page template %s: %s", name, err)
	}

	return string(baseTmpl) + string(pageTmpl), nil
}

// ExecuteTemplate executes template.
func (t *Template) ExecuteTemplate(
	w http.ResponseWriter,
	name string,
	data interface{},
) {
	templateBody, err := t.getTemplate(name)
	if err != nil {
		http.Error(w, "error getting template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	htmlTemplate := htmltemplate.New(name).Funcs(htmltemplate.FuncMap(t.funcs))

	htmlTemplate, err = htmlTemplate.Parse(templateBody)
	if err != nil {
		http.Error(w, "error parsing template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var buffer bytes.Buffer

	err = htmlTemplate.Execute(&buffer, data)
	if err != nil {
		http.Error(w, "error executing template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	io.WriteString(w, buffer.String())
}
