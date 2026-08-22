package project_generator

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"text/template"
)

type LogTmpl struct {
	AppPkgPath string
	LogPath    string
}

type LogTmplOption func(*LogTmpl)

func NewLogTmpl(appPkgPath, logPath string, options ...LogTmplOption) *LogTmpl {
	st := &LogTmpl{
		AppPkgPath: appPkgPath,
		LogPath:    logPath,
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
	content := buf.Bytes()
	if isDo() {
		content = stripWire(content)
	}
	if err := os.WriteFile(path.Join(st.LogPath, "log.go"), content, 0o644); err != nil {
		return err
	}
	return nil
}
