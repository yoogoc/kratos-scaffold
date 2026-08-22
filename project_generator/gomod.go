package project_generator

import (
	"fmt"
	"os"
	"path"

	"github.com/yoogoc/kratos-scaffold/pkg/util"
)

const (
	goVersion = "1.26.5"
	// kratos/v3 is released, the contrib modules are not tagged yet.
	kratosV3Version              = "v3.0.0"
	kratosV3ContribPseudoVersion = "v3.0.0-20260526000039-30da04b769dc"
)

// genTools pins the protoc plugins and wire in go.mod's `tool` block. buf.gen.*
// resolves them with `go tool`, so nothing has to be installed by hand.
func genTools(di string) error {
	args := []string{
		"get",
		"-tool",
		"google.golang.org/protobuf/cmd/protoc-gen-go",
		"google.golang.org/grpc/cmd/protoc-gen-go-grpc",
		"github.com/go-kratos/kratos/cmd/protoc-gen-go-http/v3",
		"github.com/google/gnostic/cmd/protoc-gen-openapi",
		"github.com/envoyproxy/protoc-gen-validate",
	}
	if di == "wire" {
		args = append(args, "github.com/google/wire/cmd/wire")
	}
	return util.Exec("go", args...)
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
