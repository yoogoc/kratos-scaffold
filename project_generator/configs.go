package project_generator

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

const configEmbedGo = "package configs\n\nimport \"embed\"\n\n//go:embed default\nvar DefaultFS embed.FS\n"
const localIgnoreRules = "configs/local/**\n!configs/local/\n!configs/local/**/\n!configs/local/**/.keep\n"

func genConfigRoot(projectPath string) error {
	defaultPath := filepath.Join(projectPath, "configs", "default")
	localPath := filepath.Join(projectPath, "configs", "local")
	if err := os.MkdirAll(defaultPath, 0o700); err != nil {
		return err
	}
	if err := os.MkdirAll(localPath, 0o700); err != nil {
		return err
	}
	if err := writeIfMissing(filepath.Join(defaultPath, "config.yaml"), []byte(configContent)); err != nil {
		return err
	}
	if err := writeIfMissing(filepath.Join(localPath, ".keep"), nil); err != nil {
		return err
	}
	if err := writeIfMissing(filepath.Join(projectPath, "configs", "embed.go"), []byte(configEmbedGo)); err != nil {
		return err
	}
	return appendUnique(filepath.Join(projectPath, ".gitignore"), localIgnoreRules)
}

func genServiceConfig(projectPath, name string) error {
	defaultPath := filepath.Join(projectPath, "configs", "default", filepath.FromSlash(name))
	localPath := filepath.Join(projectPath, "configs", "local", filepath.FromSlash(name))
	if err := os.MkdirAll(defaultPath, 0o700); err != nil {
		return err
	}
	if err := os.MkdirAll(localPath, 0o700); err != nil {
		return err
	}
	if err := writeIfMissing(filepath.Join(defaultPath, "config.yaml"), []byte(configContent)); err != nil {
		return err
	}
	return writeIfMissing(filepath.Join(localPath, ".keep"), nil)
}

func writeIfMissing(name string, content []byte) error {
	if _, err := os.Stat(name); err == nil {
		return nil
	}
	return os.WriteFile(name, content, 0o644)
}

func appendUnique(name, block string) error {
	content, err := os.ReadFile(name)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	if strings.Contains(string(content), strings.TrimSpace(block)) {
		return nil
	}
	prefix := ""
	if len(content) > 0 && content[len(content)-1] != '\n' {
		prefix = "\n"
	}
	if len(content) > 0 {
		prefix += "\n"
	}
	return os.WriteFile(name, append(content, []byte(prefix+block)...), 0o644)
}
