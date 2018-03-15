package dux

import (
	"bytes"
	"html/template"
	"io"
	"io/ioutil"
	"path/filepath"
)

// TemplateEngine defines the interface to access a text-based
// templating system with templates stored in a file system.
type TemplateEngine interface {
	// RenderTemplates writes the data produced by rendering the
	// template identified by templateFileName to destination.
	// The optional data parameter is used as the context object
	// when rendering the template.
	RenderTemplate(destination io.Writer, templateFileName string, data ...interface{}) error

	// RenderString renders a single string as a template with the
	// given context.
	RenderString(tmpl string, data interface{}) (string, error)
}

// HTMLTemplateEngine implements TemplateEngine using html/template.
// It parses all templates in the root directory, but only renders the
// one specified by the template file name.
type HTMLTemplateEngine struct {
	dir string
	fs  FileSystem
}

// NewHTMLTemplateEngine returns a new HTMLTemplateEngine reading
// templates from the provided directory in the given file system.
func NewHTMLTemplateEngine(dir string, fs FileSystem) *HTMLTemplateEngine {
	return &HTMLTemplateEngine{
		dir: dir,
		fs:  fs,
	}
}

// RenderTemplate implements TemplateEngine
func (t *HTMLTemplateEngine) RenderTemplate(out io.Writer, templateName string, data ...interface{}) error {
	context := (interface{})(nil)
	if len(data) > 0 {
		context = data[0]
	}

	templateFiles, err := t.fs.List(t.dir)
	if err != nil {
		return err
	}

	tmpl := template.New(templateName).Funcs(t.TemplateFuncs())
	for _, filename := range templateFiles {
		templateFile, err := t.fs.Open(filepath.Join(t.dir, filename))
		if err != nil {
			return err
		}
		contents, err := ioutil.ReadAll(templateFile)
		if err != nil {
			templateFile.Close()
			return err
		}
		tmpl, err = tmpl.Parse(string(contents))
		if err != nil {
			return err
		}
	}

	return tmpl.ExecuteTemplate(out, templateName, context)
}

// RenderString implements TemplateEngine
func (t *HTMLTemplateEngine) RenderString(tmpl string, data interface{}) (string, error) {
	parsedTemplate, err := template.New("main").Funcs(t.TemplateFuncs()).Parse(tmpl)
	if err != nil {
		return tmpl, err
	}

	out := bytes.NewBufferString("")
	if err := parsedTemplate.Execute(out, data); err != nil {
		return tmpl, err
	}

	return out.String(), nil
}

// TemplateFuncs returns a template.FuncMap containing the functions that should be made available to all templates.
//
// Modifying the map returned by this functions makes it possible to add more functions to a template.
func (t *HTMLTemplateEngine) TemplateFuncs() template.FuncMap {
	return template.FuncMap{
		"identifier": ParseIdentifier,
	}
}
