package project_generator

import (
	"os"
	"path"

	"github.com/yoogoc/kratos-scaffold/pkg/util"
)

func genConf(name, appPath, text string) error {
	confPath := path.Join(appPath, "conf")
	if err := os.MkdirAll(confPath, 0o700); err != nil {
		return err
	}
	if err := os.WriteFile(path.Join(confPath, "conf.proto"), []byte(text), 0o644); err != nil {
		return err
	}
	// conf.proto is its own buf module, see resources/buf.conf.yaml.example
	if err := os.WriteFile(path.Join(confPath, "buf.yaml"), []byte(bufConfYaml), 0o644); err != nil {
		return err
	}

	if err := util.Exec("just", "config", name); err != nil {
		return err
	}
	return nil
}

func genConfSingle(appPath, text string) error {
	confPath := path.Join(appPath, "conf")
	if err := os.MkdirAll(confPath, 0o700); err != nil {
		return err
	}
	if err := os.WriteFile(path.Join(confPath, "conf.proto"), []byte(text), 0o644); err != nil {
		return err
	}

	if err := util.Exec("just", "config"); err != nil {
		return err
	}
	return nil
}
