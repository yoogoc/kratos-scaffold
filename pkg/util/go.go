package util

import (
	"os"
	"os/exec"
)

// Go exec go command
func Go(command string, path ...string) error {
	for _, p := range path {
		cmd := exec.Command("go", command, p)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return err
		}
	}
	return nil
}
