package cmd

import (
	"os"
	"path"
	"testing"
)

func TestDataEnt(t *testing.T) {
	wd, _ := os.Getwd()
	dir := makeTestSingleWorkspace(t)
	t.Cleanup(func() {
		_ = os.RemoveAll(dir)
	})

	runTestCmd(t, "data user id:int64:eq name:string:cont age:int32:gte,lte")
	_ = os.Chdir(wd)

	AssertGoldenFile(t, path.Join(dir, "internal/data/user.go"), "data-ent-user.txt")
	AssertGoldenFile(t, path.Join(dir, "internal/data/user_transfer.go"), "data-ent-user_transfer.txt")
}
