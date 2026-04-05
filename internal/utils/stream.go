package utils

import (
	"github.com/jrengmusic/cake/internal/ui"
	"fmt"
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
		return fmt.Errorf("stream: stdout pipe: %w", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("stream: stderr pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("stream: start command: %w", err)
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

// flushLine emits a non-empty line via the appropriate callback based on progress state
func flushLine(line string, lineType ui.OutputLineType, isProgress bool, appendCallback func(string, ui.OutputLineType), replaceCallback func(string, ui.OutputLineType)) {
	if line == "" {
		return
	}
	if isProgress {
		replaceCallback(line, lineType)
	} else {
		appendCallback(line, lineType)
	}
}

// streamPipe reads byte-by-byte from a reader, handling \n and \r line terminators.
// Tracks isProgressLine state:
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
				flushLine(strings.TrimSpace(currentLine.String()), lineType, isProgressLine, appendCallback, replaceCallback)
				currentLine.Reset()
				isProgressLine = false
			case '\r':
				line := strings.TrimSpace(currentLine.String())
				if line != "" && !isProgressLine {
					appendCallback(line, lineType)
					isProgressLine = true
				} else {
					flushLine(line, lineType, isProgressLine, appendCallback, replaceCallback)
				}
				currentLine.Reset()
			default:
				currentLine.WriteByte(ch)
			}
		}
		if err == io.EOF {
			flushLine(strings.TrimSpace(currentLine.String()), lineType, isProgressLine, appendCallback, replaceCallback)
			break
		}
		if err != nil {
			break
		}
	}
}
