package project_generator

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"text/template"
)

type OtelTmpl struct {
	AppPkgPath  string
	ServiceName string
	OtelPath    string
}

func NewOtelTmpl(appPkgPath, serviceName, otelPath string) *OtelTmpl {
	return &OtelTmpl{
		AppPkgPath:  appPkgPath,
		ServiceName: serviceName,
		OtelPath:    otelPath,
	}
}

func (ot OtelTmpl) Generate() error {
	fmt.Println("generate otel/otel.go ...")
	buf := new(bytes.Buffer)
	tmpl, err := template.New("otelTmpl").Parse(otelTmpl)
	if err != nil {
		return err
	}
	if err := tmpl.Execute(buf, ot); err != nil {
		return err
	}
	if err := os.WriteFile(path.Join(ot.OtelPath, "otel.go"), buf.Bytes(), 0o644); err != nil {
		return err
	}
	return nil
}