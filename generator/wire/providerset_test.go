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
	if !strings.Contains(string(content), "NewUserUsecase") {
		t.Fatalf("expected NewUserUsecase in ProviderSet, got:\n%s", content)
	}
	if !strings.Contains(string(content), "NewTxUsecase") {
		t.Fatalf("expected NewTxUsecase still in ProviderSet, got:\n%s", content)
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
	if !strings.Contains(string(content), "NewGreeterService") {
		t.Fatalf("expected NewGreeterService in ProviderSet, got:\n%s", content)
	}
}
