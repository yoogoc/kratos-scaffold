package project_generator

import (
	"os"
	"os/exec"
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

	if err := util.Exec("make", "config-"+name); err != nil {
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

	fd := exec.Command("protoc",
		"--proto_path=./internal/conf",
		"--proto_path=./third_party",
		"--go_out=paths=source_relative:./internal/conf",
		"conf.proto",
	)
	fd.Stdout = os.Stdout
	fd.Stderr = os.Stderr
	if err := fd.Run(); err != nil {
		return err
	}
	return nil
}
