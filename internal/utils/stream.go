package utils

import (
	"cake/internal/ui"
	"io"
	"os/exec"
	"strings"
	"sync"
)

// StreamCommand pipes stdout/stderr from a command to callback functions.
// Uses byte-by-byte reading to handle \r (progress lines) and \n (complete lines).
// Waits for pipes to fully drain before returning — callers call cmd.Wait() after.
func StreamCommand(cmd *exec.Cmd, appendCallback func(string, ui.OutputLineType), replaceCallback func(string, ui.OutputLineType)) error {
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		streamPipe(stdout, ui.TypeStdout, appendCallback, replaceCallback)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		streamPipe(stderr, ui.TypeStderr, appendCallback, replaceCallback)
	}()

	// Drain pipes before returning — prevents race with cmd.Wait()
	wg.Wait()

	return nil
}

// streamPipe reads byte-by-byte from a reader, handling \n and \r line terminators.
// Tracks isProgressLine state per TIT pattern:
//   - First \r on a line: append normally, mark as progress
//   - Subsequent \r while progress: replace last line
//   - \n: flush (replace if progress, append if not), reset progress flag
//   - EOF: flush remaining content
func streamPipe(r io.Reader, lineType ui.OutputLineType, appendCallback func(string, ui.OutputLineType), replaceCallback func(string, ui.OutputLineType)) {
	var currentLine strings.Builder
	isProgressLine := false
	oneByte := make([]byte, 1)

	for {
		n, err := r.Read(oneByte)
		if n > 0 {
			ch := oneByte[0]
			switch ch {
			case '\n':
				line := strings.TrimSpace(currentLine.String())
				if line != "" {
					if isProgressLine {
						replaceCallback(line, lineType)
					} else {
						appendCallback(line, lineType)
					}
				}
				currentLine.Reset()
				isProgressLine = false
			case '\r':
				line := strings.TrimSpace(currentLine.String())
				if line != "" {
					if isProgressLine {
						replaceCallback(line, lineType)
					} else {
						appendCallback(line, lineType)
						isProgressLine = true
					}
				}
				currentLine.Reset()
			default:
				currentLine.WriteByte(ch)
			}
		}
		if err == io.EOF {
			line := strings.TrimSpace(currentLine.String())
			if line != "" {
				if isProgressLine {
					replaceCallback(line, lineType)
				} else {
					appendCallback(line, lineType)
				}
			}
			break
		}
		if err != nil {
			break
		}
	}
}
