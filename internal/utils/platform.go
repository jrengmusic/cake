package utils

import (
	"os/exec"
	"runtime"
	"strings"
)

func GetPlatformGenerators() []string {
	switch runtime.GOOS {
	case "darwin":
		return []string{"Xcode", "Ninja"}
	case "windows":
		return []string{"Visual Studio", "Ninja Multi-Config"}
	case "linux":
		return []string{"Ninja"}
	default:
		return []string{"Ninja"}
	}
}

func GetDefaultGenerator() string {
	switch runtime.GOOS {
	case "darwin":
		return "Xcode"
	case "windows":
		return "Visual Studio"
	default:
		return "Ninja"
	}
}

func DetectCMakeVersion() (string, error) {
	cmd := exec.Command("cmake", "--version")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	lines := strings.Split(string(output), "\n")
	if len(lines) > 0 {
		parts := strings.Fields(lines[0])
		if len(parts) >= 3 {
			return parts[2], nil
		}
	}

	return "unknown", nil
}
