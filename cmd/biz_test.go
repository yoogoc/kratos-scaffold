package cmd

import (
	"os"
	"path"
	"testing"
)

func TestBiz(t *testing.T) {
	wd, _ := os.Getwd()
	dir := makeTestSingleWorkspace(t)
	t.Cleanup(func() {
		_ = os.RemoveAll(dir)
	})

	runTestCmd(t, "biz user id:int64:eq,in name:string:cont age:int32:gte,lte")
	_ = os.Chdir(wd)

	AssertGoldenFile(t, path.Join(dir, "internal/biz/user.go"), "biz-user.txt")
}
