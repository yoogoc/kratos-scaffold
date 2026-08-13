package generator

import (
	"bytes"
	_ "embed"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"

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

func (p *Proto) PatchFields() []*field.Field {
	return p.Fields.PatchFields(p.PrimaryField())
}

func (p *Proto) EntityFields() []*field.Field {
	return p.Fields.EntityFields()
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

func (p *Proto) OrderByParamName() string {
	return field.StyleFieldMap[p.FieldStyle]("orderBy")
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
	fmt.Println("exec buf generate...")
	if err := p.genClient(p.OutPath(), p.GenHttp); err != nil {
		return errors.Wrap(err, "gen proto client error")
	}
	return nil
}

// genClient runs buf against the single proto just written. It mirrors the
// `proto` recipe of the generated justfile, which cannot be reused here because
// the proto may not sit under api/<namespace>.
func (p *Proto) genClient(source string, genHttp bool) error {
	if err := util.Exec("buf", "generate", "--path", source); err != nil {
		return err
	}
	if genHttp {
		out := path.Join("openapi", p.Namespace)
		if err := os.MkdirAll(out, 0o700); err != nil {
			return err
		}
		if err := util.Exec("buf", "generate",
			"--template", "buf.gen.openapi.yaml",
			"--path", source,
			"-o", out,
		); err != nil {
			return err
		}
	}
	return nil
}

func (p *Proto) OutPath() string {
	return path.Join(p.Path(), strings.ToLower(p.Name)+".proto")
}

func (p *Proto) Path() string {
	return path.Join(p.ApiDirName, p.Namespace, p.ApiVersion)
}

func (p *Proto) GoPackage() string {
	s := strings.Split(p.Path(), "/")
	return util.ModName() + "/" + p.Path() + ";" + s[len(s)-1]
}

func (p *Proto) JavaPackage() string {
	return p.Package()
}

func (p *Proto) Package() string {
	return strings.ReplaceAll(strings.ReplaceAll(p.Path(), "/", "."), "-", "_")
}
