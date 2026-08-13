package project_generator

import (
	"fmt"
	"os"
	"path"

	"github.com/yoogoc/kratos-scaffold/pkg/util"
)

const (
	goVersion = "1.26"
	// kratos/v3 is released, the contrib modules are not tagged yet.
	kratosV3Version              = "v3.0.0"
	kratosV3ContribPseudoVersion = "v3.0.0-20260526000039-30da04b769dc"
)

// genTools pins the protoc plugins and wire in go.mod's `tool` block. buf.gen.*
// resolves them with `go tool`, so nothing has to be installed by hand.
func genTools() error {
	return util.Exec(
		"go",
		"get",
		"-tool",
		"google.golang.org/protobuf/cmd/protoc-gen-go",
		"google.golang.org/grpc/cmd/protoc-gen-go-grpc",
		"github.com/go-kratos/kratos/cmd/protoc-gen-go-http/v3",
		"github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2",
		"github.com/envoyproxy/protoc-gen-validate",
		"github.com/google/wire/cmd/wire",
	)
}

func genGoMod(name, projectPath string) error {
	goModContent := fmt.Sprintf(`module %s

go %s
`, name, goVersion)
	if err := os.WriteFile(path.Join(projectPath, "go.mod"), []byte(goModContent), 0o644); err != nil {
		return err
	}
	return nil
}
