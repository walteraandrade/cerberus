package clipboard

import (
	"crypto/sha256"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

type Backend int

const (
	BackendNone Backend = iota
	BackendWlClipboard
	BackendXclip
	BackendXsel
	BackendPbcopy
)

func Detect() Backend {
	switch runtime.GOOS {
	case "darwin":
		return BackendPbcopy
	case "linux":
		if os.Getenv("WAYLAND_DISPLAY") != "" {
			if _, err := exec.LookPath("wl-copy"); err == nil {
				return BackendWlClipboard
			}
		}
		if os.Getenv("DISPLAY") != "" {
			if _, err := exec.LookPath("xclip"); err == nil {
				return BackendXclip
			}
			if _, err := exec.LookPath("xsel"); err == nil {
				return BackendXsel
			}
		}
	}
	return BackendNone
}

func Copy(text string) error {
	b := Detect()
	switch b {
	case BackendWlClipboard:
		return run("wl-copy", text)
	case BackendXclip:
		return run("xclip", "-selection", "clipboard", text)
	case BackendXsel:
		return run("xsel", "--clipboard", "--input", text)
	case BackendPbcopy:
		return run("pbcopy", text)
	default:
		return fmt.Errorf("no clipboard backend available")
	}
}

func Read() (string, error) {
	b := Detect()
	var cmd *exec.Cmd
	switch b {
	case BackendWlClipboard:
		cmd = exec.Command("wl-paste", "--no-newline")
	case BackendXclip:
		cmd = exec.Command("xclip", "-selection", "clipboard", "-o")
	case BackendXsel:
		cmd = exec.Command("xsel", "--clipboard", "--output")
	case BackendPbcopy:
		cmd = exec.Command("pbpaste")
	default:
		return "", fmt.Errorf("no clipboard backend available")
	}
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func Clear() error {
	return Copy("")
}

func CopyWithAutoClear(text string, timeout time.Duration) (<-chan struct{}, error) {
	if err := Copy(text); err != nil {
		return nil, err
	}

	checksum := sha256.Sum256([]byte(text))
	done := make(chan struct{})

	go func() {
		defer close(done)
		timer := time.NewTimer(timeout)
		defer timer.Stop()
		<-timer.C

		current, err := Read()
		if err != nil {
			return
		}
		currentHash := sha256.Sum256([]byte(current))
		if currentHash == checksum {
			Clear()
		}
	}()

	return done, nil
}

func run(name string, args ...string) error {
	input := args[len(args)-1]
	cmdArgs := args[:len(args)-1]
	cmd := exec.Command(name, cmdArgs...)
	cmd.Stdin = strings.NewReader(input)
	return cmd.Run()
}
