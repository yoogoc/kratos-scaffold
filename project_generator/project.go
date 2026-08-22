package project_generator

import (
	_ "embed"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/pkg/errors"
	"github.com/yoogoc/kratos-scaffold/pkg/scaffoldconfig"
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
	DI    string
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
	if p.Type == ProjectTypeSubMono && p.DI == "" {
		cfg, _, err := scaffoldconfig.Find(".")
		if err != nil {
			return err
		}
		p.DI = cfg.DI
	}
	di, err := scaffoldconfig.NormalizeDI(p.DI)
	if err != nil {
		return err
	}
	p.DI = di
	activeDI = di
	switch p.Type {
	case ProjectTypeMono:
		return GenMono(p.Name, p.DI)
	case ProjectTypeSubMono:
		return GenSubMono(p.Name, p.IsBff, p.DI)
	case ProjectTypeSingle:
		return GenSingle(p.Name, p.DI)
	default:
		return errors.New("unknown type")
	}
}

func GenMono(name, di string) error {
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
	if err = scaffoldconfig.Write(projectPath, di); err != nil {
		return err
	}

	err = util.Go(
		"get",
		"github.com/go-kratos/kratos/v3@"+kratosV3Version,
		"github.com/go-kratos/kratos/contrib/encoding/json/v3@"+kratosV3ContribPseudoVersion,
		"google.golang.org/grpc",
		"google.golang.org/protobuf",
		"github.com/yoogoc/kratos-kit",
		"github.com/pkg/errors",
		"github.com/gorilla/handlers",
		"github.com/spf13/cobra",
		"entgo.io/ent",
		"github.com/go-sql-driver/mysql",
		"github.com/jackc/pgx/v5",
		"github.com/samber/lo",
		"github.com/pressly/goose/v3",
	)
	if err != nil {
		return err
	}
	var diErr error
	if di == "wire" {
		diErr = util.Go("get", "github.com/google/wire")
	} else {
		diErr = util.Go("get", "github.com/samber/do/v2")
	}
	if diErr != nil {
		return diErr
	}
	if err = genTools(di); err != nil {
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
	if err := util.GenNullPath(path.Join(projectPath, "pkg")); err != nil {
		return err
	}
	// 6. cp common proto
	if err := cpProto(projectPath); err != nil {
		return err
	}
	// 7. cp justfile + buf configs
	if err = cpBuildFiles(projectPath, true); err != nil {
		return err
	}
	if err = genConfigRoot(projectPath); err != nil {
		return err
	}
	return nil
}

func GenSubMono(name string, isBff bool, di string) error {
	appDirName := "app"
	// 1. gen sub app path
	wd, _ := os.Getwd()
	subAppPath := path.Join(wd, appDirName, name)
	if err := os.MkdirAll(subAppPath, 0o700); err != nil {
		return err
	}
	parentConfig, _, err := scaffoldconfig.Find(wd)
	if err != nil {
		return err
	}
	if parentConfig.DI != di {
		return fmt.Errorf("DI mode conflict: parent uses %s but service requested %s", parentConfig.DI, di)
	}
	if err := scaffoldconfig.Write(subAppPath, di); err != nil {
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
	if err := genServiceConfig(wd, name); err != nil {
		return err
	}
	return nil
}

func GenSingle(name, di string) error {
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
	if err := scaffoldconfig.Write(appPath, di); err != nil {
		return err
	}
	// 3.1 go get dependencies
	if err := util.Go(
		"get",
		"github.com/go-kratos/kratos/v3@"+kratosV3Version,
		"github.com/go-kratos/kratos/contrib/encoding/json/v3@"+kratosV3ContribPseudoVersion,
		"google.golang.org/grpc",
		"google.golang.org/protobuf",
		"github.com/yoogoc/kratos-kit",
		"github.com/pkg/errors",
		"github.com/gorilla/handlers",
		"github.com/spf13/cobra",
		"entgo.io/ent",
		"github.com/go-sql-driver/mysql",
		"github.com/jackc/pgx/v5",
		"github.com/samber/lo",
		"github.com/pressly/goose/v3",
	); err != nil {
		return err
	}
	var diErr error
	if di == "wire" {
		diErr = util.Go("get", "github.com/google/wire")
	} else {
		diErr = util.Go("get", "github.com/samber/do/v2")
	}
	if diErr != nil {
		return diErr
	}
	if err := genTools(di); err != nil {
		return err
	}
	// 3.2 justfile + buf configs, needed by the conf/proto generation below
	if err := cpBuildFiles(appPath, false); err != nil {
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
	if err := genConfigRoot(appPath); err != nil {
		return err
	}
	return nil
}
