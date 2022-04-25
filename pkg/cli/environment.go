package cli

import (
	"os"
	"strconv"

	"github.com/spf13/pflag"
)

type PolicyIfExistsTarget int

const (
	RewriteIfExistsTarget PolicyIfExistsTarget = iota
	SkipIfExistsTarget
	ExitIfExistsTarget
)

type EnvSettings struct {
	AppDirName           string
	ApiDirName           string
	Namespace            string // only gen api or gen mono repo scaffold is active
	TemplatePath         string // TODO
	PolicyIfExistsTarget PolicyIfExistsTarget
}

func New() *EnvSettings {
	return &EnvSettings{
		AppDirName:           envOr("KRATOS_APPDIRNAME", "app"),
		ApiDirName:           envOr("KRATOS_APIDIRNAME", "api"),
		Namespace:            envOr("KRATOS_APPDIRNAME", ""),
		PolicyIfExistsTarget: PolicyIfExistsTarget(envIntOr("KRATOS_IF_EXISTS", 0)),
		// TemplatePath:         envOr("KRATOS_TEMPLATE_PATH", ""),
	}
}

func (s *EnvSettings) AddFlags(fs *pflag.FlagSet) {
	fs.StringVarP(&s.AppDirName, "app", "", s.AppDirName, "app dir name of the mono repo, only sub repo of mono repo is active")
	fs.StringVarP(&s.ApiDirName, "api", "", s.ApiDirName, "api dir name")
	fs.StringVarP(&s.Namespace, "namespace", "n", s.Namespace, "target sub-service , only gen api or gen mono repo scaffold is active")
	fs.IntVarP((*int)(&s.PolicyIfExistsTarget), "if-exists", "", int(s.PolicyIfExistsTarget), "only gen api or gen mono repo scaffold is active")
}

func envOr(name, def string) string {
	if v, ok := os.LookupEnv(name); ok {
		return v
	}
	return def
}

func envIntOr(name string, def int) int {
	if name == "" {
		return def
	}
	envVal := envOr(name, strconv.Itoa(def))
	ret, err := strconv.Atoi(envVal)
	if err != nil {
		return def
	}
	return ret
}
