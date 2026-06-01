package project_generator

import (
	"fmt"
	"os"
	"path"
)

const kratosV3PseudoVersion = "v3.0.0-20260526000039-30da04b769dc"

func genGoMod(name, projectPath string) error {
	goModContent := fmt.Sprintf(`module %s

go 1.25

replace github.com/go-kratos/kratos/v3 v3.0.0 => github.com/go-kratos/kratos/v3 %s
`, name, kratosV3PseudoVersion)
	if err := os.WriteFile(path.Join(projectPath, "go.mod"), []byte(goModContent), 0o644); err != nil {
		return err
	}
	return nil
}
