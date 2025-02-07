package sysfs

import (
	"errors"
	"os/exec"
	"strings"
)

// ShellExec executes a shell command and returns the output as a string.
func ShellExec(command string) (string, error) {
	// Validate the command to prevent execution of consecutive commands
	if strings.Contains(command, "&&") || strings.Contains(command, "||") || strings.Contains(command, ";") {
		return "", errors.New("it's forbidden to execute consecutive commands")
	}

	// Create the command
	cmd := exec.Command("bash", "-c", command)

	// Run the command and capture the output
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(output)), nil
}
