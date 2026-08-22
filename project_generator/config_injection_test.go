package project_generator

import (
	"strings"
	"testing"
)

func TestBootstrapFieldsAreInjectedSeparately(t *testing.T) {
	for _, field := range []string{"Server", "Data", "Log"} {
		doProvider := "(*conf." + field + ", error)"
		if !strings.Contains(cmdDoTmpl, doProvider) {
			t.Errorf("do template does not provide conf.%s", field)
		}
	}
	if !strings.Contains(cmdWireTmpl, `wire.FieldsOf(new(*conf.Bootstrap), "Server", "Data", "Log")`) {
		t.Error("wire template does not export Bootstrap fields")
	}
}

func TestRuntimeConstructorsDoNotDependOnBootstrap(t *testing.T) {
	for name, tmpl := range map[string]string{
		"migration": cmdMigrationTmpl,
		"http":      httpServerTmpl,
		"grpc":      grpcServerTmpl,
		"ent":       dataEntTmpl,
		"bff data":  demoGreeterBffDataTmpl,
	} {
		if strings.Contains(tmpl, "*conf.Bootstrap") {
			t.Errorf("%s template still depends on conf.Bootstrap", name)
		}
	}
}
