package app

import "cake/internal/ui"

// GenerateMenu returns exactly 8 rows using UI package
func (a *Application) GenerateMenu() []ui.MenuRow {
	buildInfo := a.projectState.GetSelectedBuildInfo()
	canOpenIDE := a.projectState.CanOpenIDE() && buildInfo.Exists
	canClean := buildInfo.Exists
	hasBuild := buildInfo.Exists
	hasBuildsToClean := a.projectState.HasBuildsToClean()

	return ui.GenerateMenuRows(
		a.projectState.GetProjectLabel(),
		a.projectState.Configuration,
		canOpenIDE,
		canClean,
		hasBuild,
		hasBuildsToClean,
	)
}
