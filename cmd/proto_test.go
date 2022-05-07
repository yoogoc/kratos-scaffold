package cmd

import (
	"os"
	"path"
	"testing"
)

func TestProto(t *testing.T) {
	wd, _ := os.Getwd()
	dir := makeTestSingleWorkspace(t)
	t.Cleanup(func() {
		_ = os.RemoveAll(dir)
	})

	runTestCmd(t, "proto -n user user id:int64:eq,in name:string:cont age:int32:gte,lte")
	_ = os.Chdir(wd)

	AssertGoldenFile(t, path.Join(dir, "api/user/v1/user.proto"), "proto.txt")
}
