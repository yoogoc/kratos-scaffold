package cmd

import (
	"os"
	"path"
	"testing"
)

func TestService(t *testing.T) {
	wd, _ := os.Getwd()
	dir := makeTestSingleWorkspace(t)
	t.Cleanup(func() {
		_ = os.RemoveAll(dir)
	})

	runTestCmd(t, "service user id:int64:eq name:string:cont age:int32:gte,lte")
	_ = os.Chdir(wd)

	AssertGoldenFile(t, path.Join(dir, "internal/service/user.go"), "service-user.txt")
	AssertGoldenFile(t, path.Join(dir, "internal/service/user_transfer.go"), "service-user_transfer.txt")
}
