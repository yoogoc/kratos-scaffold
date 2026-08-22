package scaffoldconfig

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWriteAndFind(t *testing.T) {
	root := t.TempDir()
	if err := Write(root, "do"); err != nil {
		t.Fatal(err)
	}
	child := filepath.Join(root, "a", "b")
	if err := os.MkdirAll(child, 0o700); err != nil {
		t.Fatal(err)
	}
	cfg, file, err := Find(child)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.DI != "do" || file != filepath.Join(root, FileName) {
		t.Fatalf("cfg=%#v file=%q", cfg, file)
	}
}

func TestNormalizeDI(t *testing.T) {
	if got, _ := NormalizeDI(""); got != "do" {
		t.Fatalf("default = %q", got)
	}
	if _, err := NormalizeDI("bad"); err == nil {
		t.Fatal("expected error")
	}
}
