package generator

import (
	"bytes"
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/iancoleman/strcase"
	"github.com/pkg/errors"
	"github.com/yoogoc/kratos-scaffold/pkg/cli"
	"github.com/yoogoc/kratos-scaffold/pkg/field"
	"github.com/yoogoc/kratos-scaffold/pkg/util"
)

type Proto struct {
	Base
	GenHttp    bool
	FieldStyle string
}

func NewProto(setting *cli.EnvSettings) *Proto {
	return &Proto{
		Base:       NewBase(setting, false),
		FieldStyle: setting.FieldStyle,
	}
}

func (p *Proto) CreateFields() []*field.Field {
	return p.Fields.CreateFields(p.PrimaryKey)
}

func (p *Proto) UpdateFields() []*field.Field {
	return p.Fields.UpdateFields(p.PrimaryField())
}

func (p *Proto) PrimaryField() *field.Field {
	return p.Fields.PrimaryField(p.PrimaryKey)
}

func (p *Proto) PrimaryFieldName() string {
	return p.PrimaryField().Name
}

func (p *Proto) PrimaryFieldURLName() string {
	return fmt.Sprintf("{%s}", p.PrimaryField().Name)
}

func (p *Proto) PageParamName() string {
	return field.StyleFieldMap[p.FieldStyle]("page")
}

func (p *Proto) PageSizeParamName() string {
	return field.StyleFieldMap[p.FieldStyle]("pageSize")
}

//go:embed tmpl/proto.tmpl
var protoTmpl string

func (p *Proto) Generate() error {
	fmt.Println("generating proto...")
	buf := new(bytes.Buffer)

	funcMap := template.FuncMap{
		"ToLower":    strings.ToLower,
		"ToPlural":   util.Plural,
		"add":        func(a int, b int) int { return a + b },
		"fieldStyle": func(s string) string { return field.StyleFieldMap[p.FieldStyle](s) },
	}
	tmpl, err := template.New("protoTmpl").Funcs(funcMap).Parse(protoTmpl)
	if err != nil {
		return err
	}
	err = tmpl.Execute(buf, p)
	if err != nil {
		return err
	}

	out := p.OutPath()
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	if !filepath.IsAbs(out) {
		out = path.Join(wd, out)
	}
	n := strings.LastIndex(out, "/")
	op := out[:n]
	if _, err := os.Stat(op); os.IsNotExist(err) {
		if err := os.MkdirAll(op, 0o700); err != nil {
			return err
		}
	}
	if err := os.WriteFile(out, buf.Bytes(), 0o644); err != nil {
		return err
	}
	fmt.Println("exec protoc...")
	if err := p.genClient(p.OutPath(), p.GenHttp); err != nil {
		return errors.Wrap(err, "gen proto client error")
	}
	return nil
}

// genClient OPTIMIZE
func (p *Proto) genClient(source string, genHttp bool) error {
	args := []string{
		"--proto_path=.",
		"--proto_path=./third_party",
		"--go_out=paths=source_relative:.",
		"--go-grpc_out=paths=source_relative:.",
		// "--validate_out=paths=source_relative,lang=go:.",
	}
	if genHttp {
		args = append(args, "--go-http_out=paths=source_relative:.",
			"--openapiv2_out=.",
			"--openapiv2_opt=logtostderr=true",
			"--openapiv2_opt=json_names_for_fields=false",
		)
	}
	cmd := exec.Command("protoc", append(args, source)...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = "."
	return cmd.Run()
}

func (p *Proto) OutPath() string {
	return path.Join(p.Path(), strings.ToLower(p.Name)+".proto")
}

func (p *Proto) Path() string {
	return path.Join(p.ApiDirName, p.Namespace, util.DefaultApiVersion)
}

func (p *Proto) GoPackage() string {
	s := strings.Split(p.Path(), "/")
	return strcase.ToSnake(util.ModName()) + "/" + p.Path() + ";" + s[len(s)-1]
}

func (p *Proto) JavaPackage() string {
	return p.Package()
}

func (p *Proto) Package() string {
	return strings.ReplaceAll(p.Path(), "/", ".")
}
