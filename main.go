// Colordiff is an updated fork of artyom's 2022 colordiff Go port.
//
// - Introduces stricter linter checks
// - Passes context to the subprocess
// - Revert coloring to match the original kimmel 1991 colordiff
package main

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
)

const (
	colorFileName     = "\x1b[1m" // "\x1b[1;34m"
	colorMeta         = "\x1b[35m"
	colorContextRange = "\x1b[36m"
	colorLineAdded    = "\x1b[32m"
	colorLineRemoved  = "\x1b[31m"
	colorReset        = "\x1b[0m"
	expectedArgs      = 2
	noArgs            = 0
)

var errUsage = errors.New("usage: colordiff file1 file2\nor: diff -u file1 file2 | colordiff")

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	exitCode, err := run(ctx, os.Args[1:])
	cancel()
	if err != nil {
		_, _ = os.Stderr.WriteString(err.Error() + "\n")
	}
	os.Exit(exitCode)
}

func run(ctx context.Context, args []string) (int, error) {
	var sc *bufio.Scanner
	var err error
	switch len(args) {
	case noArgs:
		sc = bufio.NewScanner(os.Stdin)
	case expectedArgs:
		// validate each arg is a real file before passing them to exec
		for _, arg := range args {
			if _, err := os.Stat(arg); err != nil { //nolint:gosec // path traversal acceptable
				return 1, fmt.Errorf("file error: %w", err)
			}
		}
		//nolint:gosec // files are confirmed so safe to call 'diff' with context
		cmd := exec.CommandContext(ctx, "diff", "-u", args[0], args[1])
		cmd.Stderr = os.Stderr
		b, err := cmd.Output()
		if err != nil && len(b) == 0 { // diff exits with 0 only on empty output
			return 1, fmt.Errorf("diff command failed: %w", err)
		}
		sc = bufio.NewScanner(bytes.NewReader(b))
	default:
		return 1, errUsage
	}
	w := os.Stdout
	var text []byte
	for sc.Scan() {
		text = text[:0]
		switch line := sc.Bytes(); {
		case len(line) == 0:
		case bytes.HasPrefix(line, []byte("--- ")):
			fallthrough
		case bytes.HasPrefix(line, []byte("+++ ")):
			text = append(text, colorFileName...)
			text = append(text, line...)
			text = append(text, colorReset...)
		case line[0] == '-':
			text = append(text, colorLineRemoved...)
			text = append(text, line...)
			text = append(text, colorReset...)
		case line[0] == '+':
			text = append(text, colorLineAdded...)
			text = append(text, line...)
			text = append(text, colorReset...)
		case line[0] == '@':
			if i := bytes.Index(line, []byte(" @@")); i != -1 {
				text = append(text, colorContextRange...)
				text = append(text, line[:i+3]...)
				text = append(text, colorReset...)
				text = append(text, line[i+3:]...)
			} else {
				text = append(text, line...)
			}
		case line[0] == ' ':
			text = append(text, line...)
		default:
			text = append(text, colorMeta...)
			text = append(text, line...)
			text = append(text, colorReset...)
		}
		text = append(text, '\n')
		if _, err = w.Write(text); err != nil {
			return 1, fmt.Errorf("write error: %w", err)
		}
	}
	if err = sc.Err(); err != nil {
		return 1, fmt.Errorf("scanner error: %w", err)
	}
	if cap(text) != 0 {
		return 1, nil
	}
	return 0, nil
}
