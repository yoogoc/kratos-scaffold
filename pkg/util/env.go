package util

import (
	"os"
	"strconv"
)

const (
	DefaultApiVersion = "v1"
)

func EnvOr(name, def string) string {
	if v, ok := os.LookupEnv(name); ok {
		return v
	}
	return def
}

func EnvIntOr(name string, def int) int {
	if name == "" {
		return def
	}
	envVal := EnvOr(name, strconv.Itoa(def))
	ret, err := strconv.Atoi(envVal)
	if err != nil {
		return def
	}
	return ret
}
