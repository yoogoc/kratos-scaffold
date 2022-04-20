package util

import (
	"os"
	"path"
)

func GenNullPath(p string) error {
	if err := os.MkdirAll(p, 0o700); err != nil {
		return err
	}
	return os.WriteFile(path.Join(p, ".keep"), []byte{}, 0o644)
}
