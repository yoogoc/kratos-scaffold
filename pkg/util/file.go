package util

import (
	"os"
	"path"

	"golang.org/x/mod/modfile"
)

func GenNullPath(p string) error {
	if err := os.MkdirAll(p, 0o700); err != nil {
		return err
	}
	return os.WriteFile(path.Join(p, ".keep"), []byte{}, 0o644)
}

func ModName() string {
	modBytes, err := os.ReadFile("go.mod")
	if err != nil {
		if modBytes, err = os.ReadFile("../go.mod"); err != nil {
			return ""
		}
	}
	return modfile.ModulePath(modBytes)
}
