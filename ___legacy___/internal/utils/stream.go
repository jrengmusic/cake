package utils

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"regexp"
	"strings"
	"sync"

	"github.com/jrengmusic/cake/internal/ui"
)

// ninjaProgressPattern matches ninja's "[N/M]" progress prefix (e.g. "[42/150] Building ...").
// When stdout is not a TTY, ninja emits each progress update as a \n-terminated line
// instead of \r-overwriting. CAKE collapses consecutive matches to a single line.
var ninjaProgressPattern = regexp.MustCompile(`^\[\d+/\d+\]`)

// StreamCommand pipes stdout/stderr from a command to callback functions.
// Uses byte-by-byte reading to handle \r (progress lines) and \n (complete lines).
// Waits for pipes to fully drain before returning — callers call cmd.Wait() after.
// onProcessTreeStarted is called immediately after the process tree is started;
// callers use it to store tree.Close for abort-path termination.
func StreamCommand(cmd *exec.Cmd, appendCallback func(string, ui.OutputLineType), replaceCallback func(string, ui.OutputLineType), onProcessTreeStarted func(*ProcessTree)) (*ProcessTree, error) {
	stdout, stdoutErr := cmd.StdoutPipe()
	if stdoutErr != nil {
		return nil, fmt.Errorf("StreamCommand: StdoutPipe: %w", stdoutErr)
	}

	stderr, stderrErr := cmd.StderrPipe()
	if stderrErr != nil {
		return nil, fmt.Errorf("StreamCommand: StderrPipe: %w", stderrErr)
	}

	tree, startErr := StartProcessTree(cmd)
	if startErr != nil {
		return nil, fmt.Errorf("StreamCommand: StartProcessTree: %w", startErr)
	}

	if onProcessTreeStarted != nil {
		onProcessTreeStarted(tree)
	}

	stdoutBuf := bufio.NewReader(stdout)
	stderrBuf := bufio.NewReader(stderr)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		streamPipe(stdoutBuf, ui.TypeStdout, appendCallback, replaceCallback)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		streamPipe(stderrBuf, ui.TypeStderr, appendCallback, replaceCallback)
	}()

	// Drain pipes before returning — prevents race with cmd.Wait()
	wg.Wait()

	return tree, nil
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
	lastEmittedWasProgress := false
	oneByte := make([]byte, 1)

	for {
		n, err := r.Read(oneByte)
		if n > 0 {
			ch := oneByte[0]
			switch ch {
			case '\n':
				line := strings.TrimSpace(currentLine.String())
				if line != "" {
					matchesNinjaProgress := ninjaProgressPattern.MatchString(line)

					shouldReplace := isProgressLine || (matchesNinjaProgress && lastEmittedWasProgress)
					if shouldReplace {
						replaceCallback(line, lineType)
					} else {
						appendCallback(line, lineType)
					}

					lastEmittedWasProgress = matchesNinjaProgress
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
