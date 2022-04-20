package util

import (
	"fmt"
	"os"
	"os/exec"
)

// Go exec go command
func Go(command string, path ...string) error {
	for _, p := range path {
		fmt.Printf("go %s %s\n", command, p)
		cmd := exec.Command("go", command, p)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return err
		}
	}
	return nil
}
