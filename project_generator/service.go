package project_generator

import (
	"os"
	"path"
)

func genService(appPath string) error {
	servicePath := path.Join(appPath, "service")
	if err := os.MkdirAll(servicePath, 0o700); err != nil {
		return err
	}
	serviceContent := `package service

import "github.com/google/wire"

// ProviderSet is service providers.
var ProviderSet = wire.NewSet()
`
	if err := os.WriteFile(path.Join(servicePath, "service.go"), []byte(serviceContent), 0o644); err != nil {
		return err
	}

	return nil
}
