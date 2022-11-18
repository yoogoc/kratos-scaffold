package project_generator

import (
	"fmt"
	"os"
	"path"

	"github.com/yoogoc/kratos-scaffold/generator"
	"github.com/yoogoc/kratos-scaffold/pkg/util"
)

func initEnt(appDataPath string) error {
	entPath := path.Join(appDataPath, "ent")
	entSchemaPath := path.Join(entPath, "schema")
	if err := os.MkdirAll(entSchemaPath, 0o700); err != nil {
		return err
	}
	if err := os.WriteFile(path.Join(entSchemaPath, "needremove.go"), []byte(entSchemaNeedRemove), 0o644); err != nil {
		return err
	}
	if err := generator.GenEntBase(entPath); err != nil {
		return err
	}
	if err := util.Go("generate", entPath); err != nil {
		return err
	}
	fmt.Println("ent至少有一个schema才会生成客户端代码,项目生成后,请自行删除ent/下needremove开头的文件和文件夹")
	return nil
}
