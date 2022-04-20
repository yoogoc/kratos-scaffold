package project_generator

import (
	_ "embed"
	"os"
	"path"
)

//go:embed resources/config.yaml.example
var configContent string

//go:embed resources/cmd.main.go.tmpl
var cmdMainTmpl string

//go:embed resources/cmd.wire.go.tmpl
var cmdWireTmpl string

//go:embed resources/biz.tx.go.example
var bizTxGo string

//go:embed resources/data.ent.schema.needremove.go.example
var entSchemaNeedRemove string

//go:embed resources/data.data.go.tmpl
var dataTmpl string

//go:embed resources/data.tx.go.tmpl
var txTmpl string

//go:embed resources/server.grpc.go.tmpl
var grpcServerTmpl string

//go:embed resources/server.http.go.tmpl
var httpServerTmpl string

//go:embed resources/conf.proto.example
var confProto string

//go:embed resources/annotations.proto
var annotationsProto string

//go:embed resources/descriptor.proto
var descriptorProto string

//go:embed resources/duration.proto
var durationProto string

//go:embed resources/errors.proto
var errorsProto string

//go:embed resources/http.proto
var httpProto string

//go:embed resources/httpbody.proto
var httpbodyProto string

//go:embed resources/struct.proto
var structProto string

//go:embed resources/timestamp.proto
var timestampProto string

//go:embed resources/wrappers.proto
var wrappersProto string

//go:embed resources/validate.proto
var validateProto string

func cpProto(projectPath string) error {
	if err := os.MkdirAll(path.Join(projectPath, "third_party"), 0o700); err != nil {
		return err
	}
	if err := os.MkdirAll(path.Join(projectPath, "third_party/errors"), 0o700); err != nil {
		return err
	}
	if err := os.WriteFile(path.Join(projectPath, "third_party/errors/errors.proto"), []byte(errorsProto), 0o644); err != nil {
		return err
	}
	if err := os.MkdirAll(path.Join(projectPath, "third_party/google"), 0o700); err != nil {
		return err
	}
	if err := os.MkdirAll(path.Join(projectPath, "third_party/google/api"), 0o700); err != nil {
		return err
	}
	if err := os.WriteFile(path.Join(projectPath, "third_party/google/api/annotations.proto"), []byte(annotationsProto), 0o644); err != nil {
		return err
	}
	if err := os.WriteFile(path.Join(projectPath, "third_party/google/api/http.proto"), []byte(httpProto), 0o644); err != nil {
		return err
	}
	if err := os.WriteFile(path.Join(projectPath, "third_party/google/api/httpbody.proto"), []byte(httpbodyProto), 0o644); err != nil {
		return err
	}
	if err := os.MkdirAll(path.Join(projectPath, "third_party/google/protobuf"), 0o700); err != nil {
		return err
	}
	if err := os.WriteFile(path.Join(projectPath, "third_party/google/protobuf/descriptor.proto"), []byte(descriptorProto), 0o644); err != nil {
		return err
	}
	if err := os.WriteFile(path.Join(projectPath, "third_party/google/protobuf/duration.proto"), []byte(durationProto), 0o644); err != nil {
		return err
	}
	if err := os.WriteFile(path.Join(projectPath, "third_party/google/protobuf/struct.proto"), []byte(structProto), 0o644); err != nil {
		return err
	}
	if err := os.WriteFile(path.Join(projectPath, "third_party/google/protobuf/timestamp.proto"), []byte(timestampProto), 0o644); err != nil {
		return err
	}
	if err := os.WriteFile(path.Join(projectPath, "third_party/google/protobuf/wrappers.proto"), []byte(wrappersProto), 0o644); err != nil {
		return err
	}
	if err := os.MkdirAll(path.Join(projectPath, "third_party/validate"), 0o700); err != nil {
		return err
	}
	if err := os.WriteFile(path.Join(projectPath, "third_party/validate/validate.proto"), []byte(validateProto), 0o644); err != nil {
		return err
	}
	return nil
}
