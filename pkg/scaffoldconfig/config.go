package scaffoldconfig

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

const FileName = ".kratos-scaffold.yaml"

type Config struct {
	Version int
	DI      string
}

func NormalizeDI(value string) (string, error) {
	if value == "" {
		return "do", nil
	}
	if value != "do" && value != "wire" {
		return "", fmt.Errorf("unsupported dependency injection mode %q (want do or wire)", value)
	}
	return value, nil
}

func Write(dir, di string) error {
	di, err := NormalizeDI(di)
	if err != nil {
		return err
	}
	content, err := yaml.Marshal(Config{Version: 1, DI: di})
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, FileName), content, 0o644)
}

func Find(start string) (Config, string, error) {
	dir, err := filepath.Abs(start)
	if err != nil {
		return Config{}, "", err
	}
	for {
		p := filepath.Join(dir, FileName)
		content, readErr := os.ReadFile(p)
		if readErr == nil {
			var cfg Config
			if err := yaml.Unmarshal(content, &cfg); err != nil {
				return Config{}, p, err
			}
			if cfg.DI != "do" && cfg.DI != "wire" {
				return Config{}, p, fmt.Errorf("invalid di %q in %s", cfg.DI, p)
			}
			return cfg, p, nil
		}
		if !errors.Is(readErr, os.ErrNotExist) {
			return Config{}, p, readErr
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return Config{Version: 1, DI: "wire"}, "", nil
		}
		dir = parent
	}
}
