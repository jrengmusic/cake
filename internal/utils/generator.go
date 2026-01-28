package utils

import (
	"fmt"
	"path/filepath"
	"strings"
)

func BuildCMakeCommand(generator string, projectRoot string) (string, []string, error) {
	validGenerators := []string{
		"Xcode", "Ninja", "Visual Studio", "Ninja Multi-Config",
	}

	isValid := false
	for _, valid := range validGenerators {
		if generator == valid {
			isValid = true
			break
		}
	}

	if !isValid {
		return "", nil, fmt.Errorf("invalid generator: %s", generator)
	}

	buildDir := filepath.Join(projectRoot, "build_"+strings.ToLower(strings.ReplaceAll(generator, " ", "_")))

	args := []string{
		"-G", generator,
		"-S", projectRoot,
		"-B", buildDir,
	}

	return buildDir, args, nil
}

func BuildBuildCommand(generator string, buildDir string) (string, []string, error) {
	switch generator {
	case "Xcode":
		return "xcodebuild", []string{"-scheme", "cake"}, nil
	case "Ninja":
		return "ninja", []string{"-C", buildDir}, nil
	case "Ninja Multi-Config":
		return "cmake", []string{"--build", buildDir, "--config", "Release"}, nil
	case "Visual Studio":
		return "cmake", []string{"--build", buildDir, "--config", "Release"}, nil
	default:
		return "cmake", []string{"--build", buildDir}, nil
	}
}
