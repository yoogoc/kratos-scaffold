package project_generator

import (
	"embed"
	"io/fs"
	"os"
	"path"
	"strings"
)

//go:embed resources/config.yaml.example
var configContent string

//go:embed resources/cmd.main.go.tmpl
var cmdMainTmpl string

//go:embed resources/cmd.wire.go.tmpl
var cmdWireTmpl string

//go:embed resources/cmd.do.go.tmpl
var cmdDoTmpl string

//go:embed resources/cmd.server.go.tmpl
var cmdServerTmpl string

//go:embed resources/cmd.migration.go.tmpl
var cmdMigrationTmpl string

//go:embed resources/biz.tx.go.example
var bizTxGo string

//go:embed resources/biz.pagination.go.example
var bizPaginationGo string

//go:embed resources/data.ent.schema.needremove.go.example
var entSchemaNeedRemove string

//go:embed resources/data.ent.data.go.tmpl
var dataEntTmpl string

//go:embed resources/data.grpc.data.go.tmpl
var dataGrpcTmpl string

//go:embed resources/server.grpc.go.tmpl
var grpcServerTmpl string

//go:embed resources/server.http.go.tmpl
var httpServerTmpl string

//go:embed resources/demo.greeter.proto.tmpl
var demoGreeterProtoTmpl string

//go:embed resources/demo.greeter.biz.go.example
var demoGreeterBizGo string

//go:embed resources/demo.greeter.service.go.tmpl
var demoGreeterServiceTmpl string

//go:embed resources/demo.greeter.data.go.tmpl
var demoGreeterDataTmpl string

//go:embed resources/demo.greeter.bff.data.go.tmpl
var demoGreeterBffDataTmpl string

//go:embed resources/conf.proto.example
var confProto string

//go:embed resources/justfile.mono.example
var justfileMonoContent string

//go:embed resources/justfile.single.example
var justfileSingleContent string

//go:embed resources/buf.yaml.mono.example
var bufYamlMono string

//go:embed resources/buf.yaml.single.example
var bufYamlSingle string

//go:embed resources/buf.gen.yaml.example
var bufGenYaml string

//go:embed resources/buf.gen.openapi.yaml.example
var bufGenOpenapiYaml string

//go:embed resources/buf.gen.validate.yaml.example
var bufGenValidateYaml string

//go:embed resources/buf.conf.yaml.example
var bufConfYaml string

//go:embed resources/conf.bff.proto.example
var confBffProto string

//go:embed resources/third_party
var thirdPartyFS embed.FS

// cpBuildFiles writes the justfile and the buf configs into the project root.
// The buf workspace layout differs between mono and single, everything else is
// shared.
func cpBuildFiles(projectPath string, isMono bool) error {
	justfile := justfileSingleContent
	bufYaml := bufYamlSingle
	if isMono {
		justfile = justfileMonoContent
		bufYaml = bufYamlMono
	}
	if isDo() {
		justfile = strings.ReplaceAll(justfile, "\n# Wire dependency injection\nwire:\n    go tool wire ./cmd\n", "")
		justfile = strings.ReplaceAll(justfile, "\n# Wire dependency injection\nwire service:\n    go tool wire ./app/{{service}}/cmd\n", "")
	}
	files := map[string]string{
		"justfile":              justfile,
		"buf.yaml":              bufYaml,
		"buf.gen.yaml":          bufGenYaml,
		"buf.gen.openapi.yaml":  bufGenOpenapiYaml,
		"buf.gen.validate.yaml": bufGenValidateYaml,
	}
	for name, content := range files {
		if err := os.WriteFile(path.Join(projectPath, name), []byte(content), 0o644); err != nil {
			return err
		}
	}
	return nil
}

func cpProto(projectPath string) error {
	return fs.WalkDir(thirdPartyFS, "resources/third_party", func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, _ := path.Split(p)
		_ = rel
		// strip "resources/third_party" prefix to get relative path
		relPath := p[len("resources/third_party"):]
		if relPath == "" {
			return nil
		}
		destPath := path.Join(projectPath, "third_party", relPath)
		if d.IsDir() {
			return os.MkdirAll(destPath, 0o700)
		}
		content, err := thirdPartyFS.ReadFile(p)
		if err != nil {
			return err
		}
		return os.WriteFile(destPath, content, 0o644)
	})
}
