package utils

import (
	"fmt"
	"os"
	"os/exec"
)

// RunShellCommand runs a shell command and streams its output to the terminal
func RunShellCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	fmt.Printf("\n▶️ Running command: %s %v\n", name, args)
	return cmd.Run()
}
