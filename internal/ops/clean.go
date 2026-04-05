package ops

import (
	"github.com/jrengmusic/cake/internal"
	"github.com/jrengmusic/cake/internal/ui"
	"github.com/jrengmusic/cake/internal/utils"
	"os"
	"path/filepath"
)

type CleanResult struct {
	Success bool
	Error   string
}

func ExecuteCleanProject(generator, config, projectRoot string, outputCallback func(string, ui.OutputLineType)) CleanResult {
	buildDir := filepath.Join(projectRoot, internal.BuildsDirName, utils.GetDirectoryName(generator))

	outputCallback("Cleaning...", ui.TypeInfo)

	if _, err := os.Stat(buildDir); os.IsNotExist(err) {
		outputCallback("Project directory clean.", ui.TypeStatus)
		outputCallback("Press ESC to return to menu", ui.TypeInfo)
		return CleanResult{Success: true}
	}

	err := os.RemoveAll(buildDir)
	if err != nil {
		outputCallback("ERROR: Failed to remove directory", ui.TypeStderr)
		outputCallback(err.Error(), ui.TypeStderr)
		return CleanResult{Success: false, Error: err.Error()}
	}

	outputCallback("ok", ui.TypeStatus)
	outputCallback("Project directory clean.", ui.TypeStatus)
	outputCallback("Press ESC to return to menu", ui.TypeInfo)
	return CleanResult{Success: true}
}
