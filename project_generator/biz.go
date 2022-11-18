package project_generator

import (
	"os"
	"path"
)

func genBffBiz(appPath string) error {
	bizPath := path.Join(appPath, "biz")
	if err := os.MkdirAll(bizPath, 0o700); err != nil {
		return err
	}
	bizContent := `package biz

import "github.com/google/wire"

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet()
`
	if err := os.WriteFile(path.Join(bizPath, "biz.go"), []byte(bizContent), 0o644); err != nil {
		return err
	}
	return nil
}
