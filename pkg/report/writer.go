package report

import (
	"io"
	"io/ioutil"
	"text/template"

	"github.com/Masterminds/sprig"
	"github.com/edersonbrilhante/vilicus/pkg/types"
	"golang.org/x/xerrors"
)

// Writer represents format writer interface
type Writer interface {
	Write(types.Analyses) error
}

// WriteAnalyses writes the analyses to output in custom template
func WriteAnalyses(output io.Writer, analyses types.Analyses, templateFile string) error {
	var writer Writer
	var err error

	if writer, err = newTemplate(output, templateFile); err != nil {
		return xerrors.Errorf("failed to initialize template writer: %w", err)
	}

	if err := writer.Write(analyses); err != nil {
		return xerrors.Errorf("failed to write analyses: %w", err)
	}
	return nil
}

// Template writes analyses in custom format defined by template
type Template struct {
	Output   io.Writer
	Template *template.Template
}

// Write writes analyses to output
func (tw Template) Write(analyses types.Analyses) error {
	err := tw.Template.Execute(tw.Output, analyses)
	if err != nil {
		return xerrors.Errorf("failed to write with template: %w", err)
	}
	return nil
}

func newTemplate(output io.Writer, templateFile string) (*Template, error) {
	buf, err := ioutil.ReadFile(templateFile)
	if err != nil {
		return nil, xerrors.Errorf("error retrieving template from path: %w", err)
	}
	templateContent := string(buf)

	var templateFuncMap template.FuncMap
	templateFuncMap = sprig.GenericFuncMap()
	templateFuncMap["vulnList"] = func(input types.Results) []types.Vuln {
		vulnList := input.VulnList()
		return vulnList
	}

	tmpl, err := template.New("output template").Funcs(templateFuncMap).Parse(templateContent)
	if err != nil {
		return nil, xerrors.Errorf("error parsing template: %w", err)
	}
	return &Template{Output: output, Template: tmpl}, nil
}
