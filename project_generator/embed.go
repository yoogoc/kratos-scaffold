package project_generator

import (
	"embed"
	"io/fs"
	"os"
	"path"
)

//go:embed resources/config.yaml.example
var configContent string

//go:embed resources/zap.go.example
var zapGoExample string

//go:embed resources/create-migration.sh.example
var createMigrationSh string

//go:embed resources/cmd.main.go.tmpl
var cmdMainTmpl string

//go:embed resources/cmd.wire.go.tmpl
var cmdWireTmpl string

//go:embed resources/cmd.server.go.tmpl
var cmdServerTmpl string

//go:embed resources/cmd.migration.go.tmpl
var cmdMigrationTmpl string

//go:embed resources/biz.tx.go.example
var bizTxGo string

//go:embed resources/biz.pagination.go.example
var bizPaginationGo string

//go:embed resources/log.log.go.tmpl
var logTmpl string

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

//go:embed resources/otel.otel.go.tmpl
var otelTmpl string

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

//go:embed resources/Makefile.example
var makefileContent string

//go:embed resources/conf.bff.proto.example
var confBffProto string

//go:embed resources/third_party
var thirdPartyFS embed.FS

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
