//go:build !windows

package utils

import "fmt"

func FindVCVarsAll() (string, error) {
	return "", fmt.Errorf("MSVC detection not available on this platform")
}

func CaptureVSEnv(vcVarsAllPath string) ([]string, error) {
	return nil, fmt.Errorf("MSVC environment not available on this platform")
}

func DetectInstalledVSVersions() []string {
	return nil
}

func FindExecutableInEnv(executable string, env []string) string {
	return executable
}
