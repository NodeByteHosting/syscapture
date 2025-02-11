package sysfs

import (
	"errors"
	"os/exec"
	"strings"
)

// ShellExec executes a shell command and returns the output as a string.
// It prevents the execution of consecutive commands for security reasons.
func ShellExec(c string) (string, error) {
	// Check for forbidden characters to prevent consecutive commands
	if strings.Contains(c, "&&") || strings.Contains(c, "||") || strings.Contains(c, ";") {
		return "", errors.New("it's forbidden to execute consecutive commands")
	}

	// Create the command
	cmd := exec.Command("bash", "-c", c)

	// Run the command and capture the output
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(output), nil
}
