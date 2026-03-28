//go:build windows

package utils

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

const vsWhereStandardPath = `C:\Program Files (x86)\Microsoft Visual Studio\Installer\vswhere.exe`
const vcVarsAllRelativePath = `VC\Auxiliary\Build\vcvarsall.bat`
const vcVarsAllArchitecture = "x64"

// FindVCVarsAll locates vcvarsall.bat via vswhere.exe
func FindVCVarsAll() (string, error) {
	_, statErr := os.Stat(vsWhereStandardPath)
	if statErr != nil {
		return "", fmt.Errorf("vswhere.exe not found at %s: %w", vsWhereStandardPath, statErr)
	}

	cmd := exec.Command(vsWhereStandardPath,
		"-latest",
		"-products", "*",
		"-requires", "Microsoft.VisualStudio.Component.VC.Tools.x86.x64",
		"-property", "installationPath",
	)

	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("vswhere.exe failed: %w", err)
	}

	installPath := strings.TrimSpace(string(output))
	vcVarsAllPath := installPath + `\` + vcVarsAllRelativePath

	_, statErr = os.Stat(vcVarsAllPath)
	if statErr != nil {
		return "", fmt.Errorf("vcvarsall.bat not found at %s: %w", vcVarsAllPath, statErr)
	}

	return vcVarsAllPath, nil
}

// CaptureVSEnv runs vcvarsall.bat and captures the resulting environment
func CaptureVSEnv(vcVarsAllPath string) ([]string, error) {
	if vcVarsAllPath == "" {
		return nil, fmt.Errorf("CaptureVSEnv: vcVarsAllPath must not be empty")
	}

	// Build the command string for cmd.exe
	// Must use SysProcAttr.CmdLine to bypass Go's argument escaping —
	// exec.Command escapes quotes with backslashes, which cmd.exe does not understand
	cmdLine := fmt.Sprintf(`cmd /c call "%s" %s && set`, vcVarsAllPath, vcVarsAllArchitecture)

	cmd := exec.Command("cmd")
	cmd.SysProcAttr = &syscall.SysProcAttr{CmdLine: cmdLine}

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to capture VS environment: %w", err)
	}

	lines := strings.Split(string(output), "\n")
	var env []string
	for _, line := range lines {
		trimmed := strings.TrimRight(line, "\r")
		eqIdx := strings.Index(trimmed, "=")
		if eqIdx > 0 {
			env = append(env, trimmed)
		}
	}

	return env, nil
}

// FindExecutableInEnv searches for an executable in the PATH from a captured environment.
// Returns the full path to the executable, or the original name if not found.
func FindExecutableInEnv(executable string, env []string) string {
	const pathPrefix = "PATH="
	pathValue := ""
	for _, e := range env {
		if strings.HasPrefix(strings.ToUpper(e), pathPrefix) {
			pathValue = e[len(pathPrefix):]
			break
		}
	}

	if pathValue == "" {
		return executable
	}

	for _, dir := range strings.Split(pathValue, ";") {
		dir = strings.TrimSpace(dir)
		if dir == "" {
			continue
		}
		fullPath := dir + `\` + executable + ".exe"
		if _, err := os.Stat(fullPath); err == nil {
			return fullPath
		}
	}

	return executable
}

// DetectInstalledVSVersions returns generator strings for all installed VS versions
func DetectInstalledVSVersions() []string {
	_, statErr := os.Stat(vsWhereStandardPath)
	if statErr != nil {
		return []string{}
	}

	cmd := exec.Command(vsWhereStandardPath,
		"-all",
		"-products", "*",
		"-requires", "Microsoft.VisualStudio.Component.VC.Tools.x86.x64",
		"-property", "installationPath",
	)

	output, err := cmd.Output()
	if err != nil {
		return []string{}
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")

	versionToGenerator := map[string]string{
		"18": GeneratorVS2026,
		"17": GeneratorVS2022,
	}

	seen := map[string]bool{}
	var result []string

	for _, line := range lines {
		path := strings.TrimSpace(strings.TrimRight(line, "\r"))
		for version, generator := range versionToGenerator {
			if strings.Contains(path, version) && !seen[generator] {
				seen[generator] = true
				result = append(result, generator)
			}
		}
	}

	return result
}
