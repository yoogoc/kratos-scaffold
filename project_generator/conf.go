package project_generator

import (
	"os"
	"os/exec"
	"path"
)

func genConf(name, appPath, text string) error {
	confPath := path.Join(appPath, "conf")
	if err := os.MkdirAll(confPath, 0o700); err != nil {
		return err
	}
	if err := os.WriteFile(path.Join(confPath, "conf.proto"), []byte(text), 0o644); err != nil {
		return err
	}

	confProtoPath := path.Join("app", name, "internal/conf", "conf.proto")

	fd := exec.Command("protoc",
		"--proto_path=.",
		"--proto_path=./third_party",
		"--go_out=paths=source_relative:.",
		path.Join(".", confProtoPath),
	)
	fd.Stdout = os.Stdout
	fd.Stderr = os.Stderr
	fd.Dir = "."
	if err := fd.Run(); err != nil {
		return err
	}
	return nil
}
