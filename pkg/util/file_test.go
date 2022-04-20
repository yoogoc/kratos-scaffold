package util

import (
	"os"
	"path"
	"testing"
)

func TestGenNullPath(t *testing.T) {
	tmpDir := os.TempDir()
	targetDir := path.Join(tmpDir, "kratos-scaffold/gen-null-path-test")
	keepPath := path.Join(targetDir, ".keep")
	err := GenNullPath(targetDir)
	defer os.RemoveAll(targetDir)
	if err != nil {
		t.Error(err)
	}

	if _, err := os.Stat(keepPath); err != nil {
		t.Error(err)
	}
}
