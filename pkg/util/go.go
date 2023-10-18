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

// Exec command
func Exec(command string, arg ...string) error {
	cmd := exec.Command(command, arg...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}
