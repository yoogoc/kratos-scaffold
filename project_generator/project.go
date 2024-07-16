package project_generator

import (
	_ "embed"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/pkg/errors"
	"github.com/yoogoc/kratos-scaffold/pkg/util"
)

type ProjectType int

const (
	ProjectTypeMono = iota
	ProjectTypeSubMono
	ProjectTypeSingle
)

func IsProjectTypeSingle() bool {
	wd, _ := os.Getwd()
	if _, err := os.Stat(path.Join(wd, "app")); err == nil {
		return false
	} else {
		return true
	}
}

type Project struct {
	Type  ProjectType
	Name  string
	IsBff bool
}

func NewProject() *Project {
	return &Project{}
}

func (p *Project) SetProjectType(isMono bool) {
	if isMono {
		p.Type = ProjectTypeMono
	} else {
		isSubMono := true
		wd, _ := os.Getwd()
		if _, err := os.Stat(path.Join(wd, "go.mod")); err != nil {
			if os.IsNotExist(err) {
				isSubMono = false
			}
		}
		if isSubMono {
			// name 只能是 xx
			p.Type = ProjectTypeSubMono
		} else {
			// name 可能是 xx, aa.com/xx/yy
			p.Type = ProjectTypeSingle
		}
	}
}

func (p *Project) Gen() error {
	switch p.Type {
	case ProjectTypeMono:
		return GenMono(p.Name)
	case ProjectTypeSubMono:
		return GenSubMono(p.Name, p.IsBff)
	case ProjectTypeSingle:
		return GenSingle(p.Name)
	default:
		return errors.New("unknown type")
	}
}

func GenMono(name string) error {
	// 1. find project name
	ss := strings.Split(name, "/")
	projectName := ss[len(ss)-1]
	// 2. mkdir <project name>
	wd, _ := os.Getwd()
	projectPath := path.Join(wd, projectName)
	if err := os.MkdirAll(projectPath, 0o700); err != nil {
		return err
	}
	// 3. Chdir <project path>
	err := os.Chdir(projectPath)
	if err != nil {
		return err
	}
	// 4. init go.mod, readme
	if err = genGoMod(name, projectPath); err != nil {
		return err
	}

	err = util.Go(
		"get",
		"github.com/go-kratos/kratos/v2",
		"google.golang.org/grpc",
		"google.golang.org/protobuf",
		"github.com/google/wire",
		"github.com/pkg/errors",
		"github.com/gorilla/handlers",
	)
	if err != nil {
		return err
	}
	if err = os.WriteFile(path.Join(projectPath, "README.md"), []byte(name), 0o644); err != nil {
		return err
	}
	// 5. mkdir api, app, pkg
	fmt.Println("generating api/ ...")
	if err := util.GenNullPath(path.Join(projectPath, "api")); err != nil {
		return err
	}
	fmt.Println("generating app/ ...")
	if err := util.GenNullPath(path.Join(projectPath, "app")); err != nil {
		return err
	}
	fmt.Println("generating pkg/ ...")
	servicePath := path.Join(projectPath, "pkg/contrib")
	if err := os.MkdirAll(servicePath, 0o700); err != nil {
		return err
	}
	// 6. cp common proto
	if err := cpProto(projectPath); err != nil {
		return err
	}
	// 7. cp create-migration.sh
	if err = os.WriteFile(path.Join(projectPath, "create-migration.sh"), []byte(createMigrationSh), 0o644); err != nil {
		return err
	}
	// 7. cp log
	if err = os.WriteFile(path.Join(projectPath, "pkg/contrib/zap.go"), []byte(zapGoExample), 0o644); err != nil {
		return err
	}
	return nil
}

func GenSubMono(name string, isBff bool) error {
	appDirName := "app"
	// 1. gen sub app path
	wd, _ := os.Getwd()
	subAppPath := path.Join(wd, appDirName, name)
	if err := os.MkdirAll(subAppPath, 0o700); err != nil {
		return err
	}
	// 2. gen internal: biz,data,service,server,conf
	if isBff {
		if err := genBffInternal(name, path.Join(subAppPath, "internal"), true); err != nil {
			return err
		}
	} else {
		if err := genInternal(name, path.Join(subAppPath, "internal"), true); err != nil {
			return err
		}
	}

	// 3. gen cmd
	if err := genCmd(name, subAppPath, true, isBff); err != nil {
		return err
	}
	// 4. init configs/conf.yaml
	if err := genConfigs(name, subAppPath); err != nil {
		return err
	}
	return nil
}

func GenSingle(name string) error {
	// 1. mkdir app path
	wd, _ := os.Getwd()
	appPath := path.Join(wd, name)
	if err := os.MkdirAll(appPath, 0o700); err != nil {
		return err
	}
	if err := os.Chdir(appPath); err != nil {
		return err
	}
	// 2. cp common proto
	if err := cpProto(appPath); err != nil {
		return err
	}
	// 3. gen go.mod
	if err := genGoMod(name, appPath); err != nil {
		return err
	}
	// 4. gen internal: biz,data,service,server,conf
	if err := genInternal(name, path.Join(appPath, "internal"), false); err != nil {
		return err
	}
	// 5. gen cmd
	if err := genCmd(name, appPath, false, false); err != nil {
		return err
	}
	// 6. init configs/conf.yaml
	if err := genConfigs(name, appPath); err != nil {
		return err
	}
	return nil
}
