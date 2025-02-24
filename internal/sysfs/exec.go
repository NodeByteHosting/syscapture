package sysfs

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

var (
	// ErrConsecutiveCommands is returned when attempting to execute multiple commands
	ErrConsecutiveCommands = errors.New("consecutive commands are forbidden")
	// ErrCommandTimeout is returned when command execution exceeds timeout
	ErrCommandTimeout = errors.New("command execution timed out")
	// ErrEmptyCommand is returned when attempting to execute an empty command
	ErrEmptyCommand = errors.New("empty command")
)

// CommandOptions configures command execution
type CommandOptions struct {
	Timeout     time.Duration
	WorkingDir  string
	Environment []string
}

// DefaultCommandOptions returns default command options
func DefaultCommandOptions() *CommandOptions {
	return &CommandOptions{
		Timeout: 30 * time.Second,
	}
}

// ShellExec executes a shell command and returns the output as a string.
// It prevents the execution of consecutive commands for security reasons.
func ShellExec(command string, opts *CommandOptions) (string, error) {
	if command = strings.TrimSpace(command); command == "" {
		return "", ErrEmptyCommand
	}

	// Use default options if none provided
	if opts == nil {
		opts = DefaultCommandOptions()
	}

	// Check for forbidden characters to prevent command injection
	forbiddenChars := []string{"&&", "||", ";", "|", ">", "<", "`", "$", "(", ")"}
	for _, char := range forbiddenChars {
		if strings.Contains(command, char) {
			return "", fmt.Errorf("%w: forbidden character %q", ErrConsecutiveCommands, char)
		}
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), opts.Timeout)
	defer cancel()

	// Determine shell and arguments based on OS
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.CommandContext(ctx, "cmd", "/C", command)
	} else {
		cmd = exec.CommandContext(ctx, "bash", "-c", command)
	}

	// Set working directory if specified
	if opts.WorkingDir != "" {
		cmd.Dir = opts.WorkingDir
	}

	// Set environment variables if specified
	if len(opts.Environment) > 0 {
		cmd.Env = append(cmd.Env, opts.Environment...)
	}

	// Run command and capture both stdout and stderr
	output, err := cmd.CombinedOutput()
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("%w after %v", ErrCommandTimeout, opts.Timeout)
		}
		return "", fmt.Errorf("command failed: %w: %s", err, string(output))
	}

	return strings.TrimSpace(string(output)), nil
}
