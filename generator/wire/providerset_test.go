package wire

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAddToProviderSet(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "biz.go")

	initial := `package biz

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewTxUsecase)
`
	if err := os.WriteFile(p, []byte(initial), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := AddToProviderSet(p, "NewUserUsecase"); err != nil {
		t.Fatal(err)
	}

	content, _ := os.ReadFile(p)
	s := string(content)
	if !strings.Contains(s, "NewUserUsecase") {
		t.Fatalf("expected NewUserUsecase in ProviderSet, got:\n%s", s)
	}
	if !strings.Contains(s, "NewTxUsecase") {
		t.Fatalf("expected NewTxUsecase still in ProviderSet, got:\n%s", s)
	}
	// verify separate lines
	if !strings.Contains(s, "\tNewTxUsecase,\n") {
		t.Fatalf("expected NewTxUsecase on its own line, got:\n%s", s)
	}
	if !strings.Contains(s, "\tNewUserUsecase,\n") {
		t.Fatalf("expected NewUserUsecase on its own line, got:\n%s", s)
	}

	// idempotent
	if err := AddToProviderSet(p, "NewUserUsecase"); err != nil {
		t.Fatal(err)
	}
	content2, _ := os.ReadFile(p)
	if strings.Count(string(content2), "NewUserUsecase") != 1 {
		t.Fatalf("expected exactly one NewUserUsecase, got:\n%s", content2)
	}
}

func TestAddToProviderSet_Empty(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "service.go")

	initial := `package service

import "github.com/google/wire"

var ProviderSet = wire.NewSet()
`
	if err := os.WriteFile(p, []byte(initial), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := AddToProviderSet(p, "NewGreeterService"); err != nil {
		t.Fatal(err)
	}

	content, _ := os.ReadFile(p)
	s := string(content)
	if !strings.Contains(s, "NewGreeterService") {
		t.Fatalf("expected NewGreeterService in ProviderSet, got:\n%s", s)
	}
	// verify separate line format
	if !strings.Contains(s, "\tNewGreeterService,\n") {
		t.Fatalf("expected NewGreeterService on its own line, got:\n%s", s)
	}
}

func TestAddToPackage(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "do.go")
	src := "package biz\n\nimport \"github.com/samber/do/v2\"\n\nvar Package = do.Package(do.Lazy(NewOld))\n"
	if err := os.WriteFile(p, []byte(src), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := AddToPackage(p, "NewUser", "NewUser"); err != nil {
		t.Fatal(err)
	}
	content, err := os.ReadFile(p)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Count(string(content), "do.Lazy(NewUser)") != 1 {
		t.Fatalf("not idempotent:\n%s", content)
	}
}
