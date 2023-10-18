package project_generator

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"text/template"

	"github.com/yoogoc/kratos-scaffold/pkg/util"
)

type LogTmpl struct {
	AppPkgPath string
	LogPath    string
	ModName    string
}

type LogTmplOption func(*LogTmpl)

func NewLogTmpl(appPkgPath, logPath string, options ...LogTmplOption) *LogTmpl {
	st := &LogTmpl{
		AppPkgPath: appPkgPath,
		LogPath:    logPath,
		ModName:    util.ModName(),
	}
	for _, option := range options {
		option(st)
	}
	return st
}

func (st LogTmpl) Generate() error {
	fmt.Println("generate log/log.go ...")
	buf := new(bytes.Buffer)
	tmpl, err := template.New("logTmpl").Parse(logTmpl)
	if err != nil {
		return err
	}
	if err := tmpl.Execute(buf, st); err != nil {
		return err
	}
	if err := os.WriteFile(path.Join(st.LogPath, "log.go"), buf.Bytes(), 0o644); err != nil {
		return err
	}
	return nil
}
