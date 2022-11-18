package project_generator

import (
	"os"
	"path"
)

func genConfigs(name string, appPath string) error {
	configsPath := path.Join(appPath, "configs")
	if err := os.MkdirAll(configsPath, 0o700); err != nil {
		return err
	}

	if err := os.WriteFile(path.Join(configsPath, "config.yaml"), []byte(configContent), 0o644); err != nil {
		return err
	}
	return nil
}
