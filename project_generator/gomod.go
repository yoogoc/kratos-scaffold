package project_generator

import (
	"fmt"
	"os"
	"path"
)

func genGoMod(name, projectPath string) error {
	goModContent := fmt.Sprintf("module %s\ngo 1.18", name)
	if err := os.WriteFile(path.Join(projectPath, "go.mod"), []byte(goModContent), 0o644); err != nil {
		return err
	}
	return nil
}
