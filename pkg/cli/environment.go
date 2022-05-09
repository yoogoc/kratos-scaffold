package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/YoogoC/kratos-scaffold/internal/merge"
	"github.com/YoogoC/kratos-scaffold/pkg/field"
	"github.com/YoogoC/kratos-scaffold/pkg/util"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/pflag"
	"gopkg.in/yaml.v2"
)

type PolicyIfExistsTarget int

const (
	RewriteIfExistsTarget PolicyIfExistsTarget = iota
	SkipIfExistsTarget
	ExitIfExistsTarget
)

type EnvSettings struct {
	AppDirName           string               `yaml:"app_dir_name"`
	ApiDirName           string               `yaml:"api_dir_name"`
	Namespace            string               `yaml:"namespace"`     // only gen api or gen mono repo scaffold is active
	TemplatePath         string               `yaml:"template_path"` // TODO
	PolicyIfExistsTarget PolicyIfExistsTarget `yaml:"policy_if_exists_target"`
	FieldStyle           string               `yaml:"field_style"`
	PrimaryKey           string               `yaml:"primary_key"`
}

func New() *EnvSettings {
	home, _ := homedir.Dir()

	yamlSources := make([][]byte, 0, 2)
	yamlPaths := make([]string, 0, 2)
	configPath := util.EnvOr("KRATOS_CONFIG", home+"/.kratos/config")
	if _, err := os.Stat(configPath); err == nil {
		yamlPaths = append(yamlPaths, configPath)
	}
	if wd, err := os.Getwd(); err == nil {
		configPath := wd + "/.kratos-scaffold.yaml"
		if _, err := os.Stat(configPath); err == nil {
			yamlPaths = append(yamlPaths, configPath)
		}
	}
	for _, path := range yamlPaths {
		contents, err := os.ReadFile(path)
		if err != nil {
			fmt.Printf("[kratos-scaffold] file read error: %s\n", err.Error())
			os.Exit(1)
		}
		yamlSources = append(yamlSources, contents)
	}

	buffer, err := merge.YAML(yamlSources, false)
	if err != nil {
		fmt.Printf("[kratos-scaffold] file read error: %s\n", err.Error())
		os.Exit(1)
	}

	env := &EnvSettings{
		AppDirName: "app",
		ApiDirName: "api",
		PrimaryKey: "id",
		FieldStyle: field.DefaultStyleField,
	}

	if err := yaml.Unmarshal(buffer.Bytes(), env); err != nil {
		fmt.Printf("[kratos-scaffold] yaml unmarshal error: %s\n", err.Error())
		os.Exit(1)
	}
	env.AppDirName = util.EnvOr("KRATOS_APPDIRNAME", env.AppDirName)
	env.ApiDirName = util.EnvOr("KRATOS_APIDIRNAME", env.ApiDirName)
	env.Namespace = util.EnvOr("KRATOS_NAMESPACE", env.Namespace)
	env.FieldStyle = util.EnvOr("KRATOS_FIELD_STYLE", env.FieldStyle)
	env.PrimaryKey = util.EnvOr("KRATOS_PRIMARY_KEY", env.PrimaryKey)
	env.PolicyIfExistsTarget = PolicyIfExistsTarget(util.EnvIntOr("KRATOS_IF_EXISTS", int(env.PolicyIfExistsTarget)))

	return env
}

func (s *EnvSettings) AddFlags(fs *pflag.FlagSet) {
	fs.StringVarP(&s.AppDirName, "app", "", s.AppDirName, "app dir name of the mono repo, only sub repo of mono repo is active")
	fs.StringVarP(&s.ApiDirName, "api", "", s.ApiDirName, "api dir name")
	fs.StringVarP(&s.Namespace, "namespace", "n", s.Namespace, "target sub-service , only gen api or gen mono repo scaffold is active")
	fs.IntVarP((*int)(&s.PolicyIfExistsTarget), "if-exists", "", int(s.PolicyIfExistsTarget), "only gen api or gen mono repo scaffold is active")
	fs.StringVarP(&s.FieldStyle, "field-style", "", s.FieldStyle,
		fmt.Sprintf("proto field style. allowed values: %s. defult value: %s", strings.Join(field.StyleFields, ","), s.FieldStyle))
}
