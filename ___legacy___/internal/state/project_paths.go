package state

import (
	"github.com/jrengmusic/cake/internal"
	"github.com/jrengmusic/cake/internal/utils"
	"os"
	"path/filepath"
	"strings"
)

// GetBuildDirectory returns the build directory path for the given generator
func (ps *ProjectState) GetBuildDirectory(generatorName string) string {
	// All projects are multi-config: Builds/<dir>/
	return filepath.Join(ps.WorkingDirectory, internal.BuildsDirName, utils.GetDirectoryName(generatorName))
}

// GetProjectLabel returns a display-friendly project name
func (ps *ProjectState) GetProjectLabel() string {
	name := ps.SelectedProject

	// Truncate long names for display
	switch name {
	case utils.GeneratorVS2026:
		return "VS 2026"
	case utils.GeneratorVS2022:
		return "VS 2022"
	default:
		return name
	}
}

// GetProjectName extracts project name using layered file parsing
// Priority: KANJUT Parameters.xml > set(PROJECT_NAME) > project() literal > directory name
func (ps *ProjectState) GetProjectName() string {
	// Layer 1: KANJUT Parameters.xml
	if name := ps.extractFromParametersXML(); name != "" {
		return name
	}

	// Layer 2: set(PROJECT_NAME "...") in CMakeLists.txt
	if name := ps.extractFromSetProjectName(); name != "" {
		return name
	}

	// Layer 3: project(LiteralName) in CMakeLists.txt
	if name := ps.extractFromProjectCall(); name != "" {
		return name
	}

	// Layer 4: Fallback to project working directory
	return filepath.Base(ps.WorkingDirectory)
}

// extractFromParametersXML parses KANJUT's Parameters.xml for PRODUCT name attribute
func (ps *ProjectState) extractFromParametersXML() string {
	paramsPath := filepath.Join(ps.WorkingDirectory, "Source", "layout", "Parameters.xml")

	data, err := os.ReadFile(paramsPath)
	if err != nil {
		return ""
	}

	content := string(data)

	// Find <PRODUCT ... name="..." ...>
	productIdx := strings.Index(content, "<PRODUCT")
	if productIdx == -1 {
		return ""
	}

	// Find closing > of PRODUCT tag
	closeIdx := strings.Index(content[productIdx:], ">")
	if closeIdx == -1 {
		return ""
	}

	productTag := content[productIdx : productIdx+closeIdx]

	// Extract name="value"
	nameIdx := strings.Index(productTag, `name="`)
	if nameIdx == -1 {
		return ""
	}

	valueStart := nameIdx + 6 // len(`name="`)
	valueEnd := strings.Index(productTag[valueStart:], `"`)
	if valueEnd == -1 {
		return ""
	}

	return productTag[valueStart : valueStart+valueEnd]
}

// extractFromSetProjectName parses set(PROJECT_NAME "...") from CMakeLists.txt
func (ps *ProjectState) extractFromSetProjectName() string {
	cmakePath := filepath.Join(ps.WorkingDirectory, internal.CMakeListsFile)

	data, err := os.ReadFile(cmakePath)
	if err != nil {
		return ""
	}

	content := string(data)
	lines := strings.Split(content, "\n")

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Match: set(PROJECT_NAME "value") or set(PROJECT_NAME "value" ...)
		if strings.HasPrefix(strings.ToLower(trimmed), "set(project_name") {
			// Find first quoted string
			quoteStart := strings.Index(trimmed, `"`)
			if quoteStart == -1 {
				continue
			}
			quoteEnd := strings.Index(trimmed[quoteStart+1:], `"`)
			if quoteEnd == -1 {
				continue
			}
			return trimmed[quoteStart+1 : quoteStart+1+quoteEnd]
		}
	}

	return ""
}

// parseProjectCallName extracts the literal name from the content after "project("
// Returns empty string if the name contains a variable reference
func parseProjectCallName(rest string) string {
	if strings.Contains(rest, "$") {
		return ""
	}

	// Handle quoted: project("Name" ...)
	if strings.HasPrefix(strings.TrimSpace(rest), `"`) {
		quoteStart := strings.Index(rest, `"`)
		quoteEnd := strings.Index(rest[quoteStart+1:], `"`)
		if quoteEnd != -1 {
			return rest[quoteStart+1 : quoteStart+1+quoteEnd]
		}
	}

	// Handle unquoted: project(Name) or project(Name VERSION ...)
	var name strings.Builder
	for _, ch := range strings.TrimSpace(rest) {
		if ch == ' ' || ch == '\t' || ch == ')' || ch == '"' {
			break
		}
		name.WriteRune(ch)
	}

	result := name.String()
	if result != "" && !strings.Contains(result, "$") {
		return result
	}

	return ""
}

// extractFromProjectCall parses project(Name) from CMakeLists.txt (literal names only)
func (ps *ProjectState) extractFromProjectCall() string {
	cmakePath := filepath.Join(ps.WorkingDirectory, internal.CMakeListsFile)

	data, err := os.ReadFile(cmakePath)
	if err != nil {
		return ""
	}

	lines := strings.Split(string(data), "\n")

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if !strings.HasPrefix(strings.ToLower(trimmed), "project(") {
			continue
		}

		parenStart := strings.Index(trimmed, "(")
		if parenStart == -1 {
			continue
		}

		if name := parseProjectCallName(trimmed[parenStart+1:]); name != "" {
			return name
		}
	}

	return ""
}
