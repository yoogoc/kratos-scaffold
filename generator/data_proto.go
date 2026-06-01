package generator

import (
	"bytes"
	_ "embed"
	"fmt"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/iancoleman/strcase"
	confmod "github.com/yoogoc/kratos-scaffold/generator/conf"
	"github.com/yoogoc/kratos-scaffold/generator/wire"
	"github.com/yoogoc/kratos-scaffold/pkg/field"
	"github.com/yoogoc/kratos-scaffold/pkg/util"
	"golang.org/x/tools/imports"
)

//go:embed tmpl/data_proto_data.tmpl
var dataProtoDataTmpl string

//go:embed tmpl/data_proto_transfer.tmpl
var dataProtoTransferTmpl string

func (d *Data) GenerateProto() error {
	// 1. ensure conf.proto has Service message for target model
	if err := confmod.EnsureServiceInConf(d.ConfPath(), d.TargetModel, d.Namespace, d.IsSingle); err != nil {
		fmt.Printf("failed to modify conf.proto, please add Service config manually: %v\n", err)
	}

	// 2. gen grpc client constructor into data dir
	if err := d.genClientProto(); err != nil {
		return err
	}
	// 3. gen data transfer
	if err := d.genTransferProto(); err != nil {
		return err
	}
	// 4. gen data
	if err := d.genDataProto(); err != nil {
		return err
	}

	providerSetPath := path.Join(d.OutPath(), "data.go")
	if err := wire.AddToProviderSet(providerSetPath, "New"+d.Name+"ServiceClient", "New"+d.Name+"Repo"); err != nil {
		fmt.Printf("failed to add to ProviderSet, please add manually: %v\n", err)
	}
	return nil
}

func (d *Data) genTransferProto() error {
	fmt.Println("generating data transfer...")
	buf := new(bytes.Buffer)
	funcMap := template.FuncMap{
		"ToLower":      strings.ToLower,
		"ToPlural":     util.Plural,
		"ToCamel":      strcase.ToCamel,
		"ToLowerCamel": strcase.ToLowerCamel,
		"ToEntName":    field.EntName,
	}
	tmpl, err := template.New("dataProtoTransferTmpl").Funcs(funcMap).Parse(dataProtoTransferTmpl)
	if err != nil {
		return err
	}
	err = tmpl.Execute(buf, d)
	if err != nil {
		return err
	}
	p := path.Join(d.OutPath(), strcase.ToSnake(d.Name)+"_transfer.go")
	content, err := imports.Process(p, buf.Bytes(), nil)
	if err != nil {
		return err
	}
	return os.WriteFile(p, content, 0o644)
}

func (d *Data) genClientProto() error {
	fmt.Println("generating grpc client constructor...")
	name := strcase.ToCamel(d.Name)
	lowerCamelName := strcase.ToLowerCamel(d.Name)
	targetModelCamel := strcase.ToCamel(d.TargetModel)
	protoPkg := d.ProtoPkgPath()
	confPkg := d.ConfPkgPath()
	alias := lowerCamelName + "V1"

	funcCode := fmt.Sprintf(`
func New%sServiceClient(confData *conf.Data, logger *slog.Logger) (%s.%sServiceClient, error) {
	conn, err := grpc.NewClient(
		context.Background(),
		grpc.WithTimeout(confData.%s.Timeout.AsDuration()),
		grpc.WithEndpoint(confData.%s.Address),
		grpc.WithMiddleware(
			mmd.Client(),
			tracing.Client(),
			mklog.Client(logger),
		),
	)
	if err != nil {
		return nil, err
	}
	return %s.New%sServiceClient(conn), nil
}
`, name, alias, name, targetModelCamel, targetModelCamel, alias, name)

	dataGoPath := path.Join(d.OutPath(), "data.go")
	existing, err := os.ReadFile(dataGoPath)
	if err != nil {
		return err
	}

	src := string(existing)
	checkFuncName := "func New" + name + "ServiceClient("
	if strings.Contains(src, checkFuncName) {
		return nil
	}

	src += funcCode

	if !strings.Contains(src, `"`+protoPkg+`"`) {
		src = addImport(src, alias, protoPkg)
	}
	if !strings.Contains(src, `"`+confPkg+`"`) {
		src = addImport(src, "", confPkg)
	}
	for _, imp := range [][2]string{
		{"", "context"},
		{"", "log/slog"},
		{"", "github.com/go-kratos/kratos/v3/transport/grpc"},
		{"", "github.com/go-kratos/kratos/contrib/otel/v3/tracing"},
		{"mklog", "github.com/go-kratos/kratos/v3/middleware/logging"},
		{"mmd", "github.com/go-kratos/kratos/v3/middleware/metadata"},
	} {
		if !strings.Contains(src, `"`+imp[1]+`"`) {
			src = addImport(src, imp[0], imp[1])
		}
	}

	content, err := imports.Process(dataGoPath, []byte(src), nil)
	if err != nil {
		return err
	}
	return os.WriteFile(dataGoPath, content, 0o644)
}

func addImport(src, alias, pkg string) string {
	importLine := ""
	if alias != "" {
		importLine = fmt.Sprintf("\t%s \"%s\"", alias, pkg)
	} else {
		importLine = fmt.Sprintf("\t\"%s\"", pkg)
	}
	idx := strings.Index(src, "import (")
	if idx < 0 {
		return src
	}
	closeIdx := strings.Index(src[idx:], ")")
	if closeIdx < 0 {
		return src
	}
	insertAt := idx + closeIdx
	return src[:insertAt] + importLine + "\n" + src[insertAt:]
}

func (d *Data) genDataProto() error {
	fmt.Println("generating data...")
	buf := new(bytes.Buffer)
	funcMap := template.FuncMap{
		"ToLower":       strings.ToLower,
		"ToPlural":      util.Plural,
		"ToCamel":       strcase.ToCamel,
		"ToCamelPlural": func(s string) string { return strcase.ToCamel(util.Plural(s)) },
		"ToLowerCamel":  strcase.ToLowerCamel,
		"last": func(x int, a []*field.Field) bool {
			return x == len(a)-1
		},
	}
	tmpl, err := template.New("dataProtoDataTmpl").Funcs(funcMap).Parse(dataProtoDataTmpl)
	if err != nil {
		return err
	}
	err = tmpl.Execute(buf, d)
	if err != nil {
		return err
	}
	p := path.Join(d.OutPath(), strcase.ToSnake(d.Name)+".go")
	// content := buf.Bytes()
	content, err := imports.Process(p, buf.Bytes(), nil)
	if err != nil {
		return err
	}
	return os.WriteFile(p, content, 0o644)
}
