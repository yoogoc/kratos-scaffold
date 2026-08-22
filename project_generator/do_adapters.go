package project_generator

import (
	"fmt"
	"path"

	"github.com/yoogoc/kratos-scaffold/pkg/util"
)

func genDoAdapters(appPkgPath, internalPath string, isBff bool) error {
	if !isDo() {
		return nil
	}
	bizProviders := `do.Lazy(func(i do.Injector) (*GreeterUsecase, error) {
		return NewGreeterUsecase(do.MustInvoke[GreeterRepo](i), do.MustInvoke[*slog.Logger](i)), nil
	}),`
	if !isBff {
		bizProviders = `do.Lazy(func(i do.Injector) (*TxUsecase, error) {
		return NewTxUsecase(do.MustInvoke[TxRepo](i), do.MustInvoke[*slog.Logger](i)), nil
	}),
	` + bizProviders
	}
	files := map[string]string{
		"biz/do.go": fmt.Sprintf(`package biz
import (
	"log/slog"
	"github.com/samber/do/v2"
)
var Package = do.Package(%s)
`, bizProviders),
		"service/do.go": fmt.Sprintf(`package service
import (
	"log/slog"
	"%s/internal/biz"
	"github.com/samber/do/v2"
)
var Package = do.Package(do.Lazy(func(i do.Injector) (*GreeterService, error) {
	return NewGreeterService(do.MustInvoke[*biz.GreeterUsecase](i), do.MustInvoke[*slog.Logger](i)), nil
}))
`, appPkgPath),
		"log/do.go": `package log
import (
	"log/slog"
	"github.com/samber/do/v2"
)
var Package = do.Package(do.Lazy(func(i do.Injector) (*slog.Logger, error) {
	return NewKratosLogger(do.MustInvoke[NameString](i), do.MustInvoke[VersionString](i)), nil
}))
`,
		"otel/do.go": fmt.Sprintf(`package otel
import (
	"%s/internal/conf"
	"github.com/samber/do/v2"
	"go.opentelemetry.io/otel/trace"
)
var Package = do.Package(do.Lazy(func(i do.Injector) (trace.TracerProvider, error) {
	provider, _, err := NewTracerProvider(do.MustInvoke[*conf.Trace](i))
	return provider, err
}))
`, appPkgPath),
	}
	if isBff {
		files["data/do.go"] = fmt.Sprintf(`package data
import (
	"log/slog"
	"%s/internal/biz"
	"%s/internal/conf"
	"github.com/samber/do/v2"
)
var Package = do.Package(do.Lazy(func(i do.Injector) (biz.GreeterRepo, error) {
	return NewGreeterRepo(do.MustInvoke[*conf.Data](i), do.MustInvoke[*slog.Logger](i)), nil
}))
`, appPkgPath, appPkgPath)
		files["server/do.go"] = fmt.Sprintf(`package server
import (
	"log/slog"
	"%s/internal/conf"
	"%s/internal/service"
	"github.com/go-kratos/kratos/v3/transport/http"
	"github.com/samber/do/v2"
	"go.opentelemetry.io/otel/trace"
)
var Package = do.Package(do.Lazy(func(i do.Injector) (*http.Server, error) {
	return NewHTTPServer(do.MustInvoke[*conf.Server](i), do.MustInvoke[*slog.Logger](i), do.MustInvoke[trace.TracerProvider](i), do.MustInvoke[*service.GreeterService](i)), nil
}))
`, appPkgPath, appPkgPath)
	} else {
		files["data/do.go"] = fmt.Sprintf(`package data
import (
	"log/slog"
	"%s/internal/biz"
	"%s/internal/conf"
	"%s/internal/data/ent"
	"github.com/samber/do/v2"
)
var Package = do.Package(
	do.Lazy(func(i do.Injector) (*ent.Client, error) { return NewEntClient(do.MustInvoke[*conf.Data](i), do.MustInvoke[*slog.Logger](i)), nil }),
	do.Lazy(func(i do.Injector) (*Data, error) { return NewData(do.MustInvoke[*conf.Data](i), do.MustInvoke[*slog.Logger](i), do.MustInvoke[*ent.Client](i)) }),
	do.Lazy(func(i do.Injector) (biz.TxRepo, error) { return NewTxRepo(do.MustInvoke[*Data](i)), nil }),
	do.Lazy(func(i do.Injector) (biz.GreeterRepo, error) { return NewGreeterRepo(do.MustInvoke[*Data](i), do.MustInvoke[*slog.Logger](i)), nil }),
)
`, appPkgPath, appPkgPath, appPkgPath)
		files["server/do.go"] = fmt.Sprintf(`package server
import (
	"log/slog"
	"%s/internal/conf"
	"%s/internal/service"
	"github.com/go-kratos/kratos/v3/transport/grpc"
	"github.com/go-kratos/kratos/v3/transport/http"
	"github.com/samber/do/v2"
	"go.opentelemetry.io/otel/trace"
)
var Package = do.Package(
	do.Lazy(func(i do.Injector) (*grpc.Server, error) { return NewGRPCServer(do.MustInvoke[*conf.Server](i), do.MustInvoke[*slog.Logger](i), do.MustInvoke[trace.TracerProvider](i), do.MustInvoke[*service.GreeterService](i)), nil }),
	do.Lazy(func(i do.Injector) (*http.Server, error) { return NewHTTPServer(do.MustInvoke[*conf.Server](i), do.MustInvoke[*slog.Logger](i), do.MustInvoke[trace.TracerProvider](i), do.MustInvoke[*service.GreeterService](i)), nil }),
)
`, appPkgPath, appPkgPath)
	}
	for name, content := range files {
		if err := util.WhiteGo(path.Join(internalPath, name), []byte(content)); err != nil {
			return err
		}
	}
	return nil
}
