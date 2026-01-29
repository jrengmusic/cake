package utils

import (
	"bufio"
	"cake/internal/ui"
	"io"
	"os/exec"
)

// StreamCommand pipes stdout/stderr from a command to a callback function.
// This eliminates duplication in ops/*.go files.
func StreamCommand(cmd *exec.Cmd, outputCallback func(string, ui.OutputLineType)) error {
	// Get stdout pipe for streaming
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	// Get stderr pipe for streaming
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	// Start command
	if err := cmd.Start(); err != nil {
		return err
	}

	// Stream stdout in goroutine
	go streamPipe(stdout, ui.TypeStdout, outputCallback)

	// Stream stderr in goroutine
	go streamPipe(stderr, ui.TypeStderr, outputCallback)

	return nil
}

// streamPipe reads from a reader and streams lines to the callback
func streamPipe(r io.Reader, lineType ui.OutputLineType, callback func(string, ui.OutputLineType)) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		if line != "" {
			callback(line, lineType)
		}
	}
}
