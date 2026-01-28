package app

import "cake/internal/ui"

// GenerateMenu returns exactly 7 rows using UI package
func (a *Application) GenerateMenu() []ui.MenuRow {
	buildInfo := a.projectState.GetSelectedBuildInfo()
	canOpenIDE := a.projectState.CanOpenIDE() && buildInfo.Exists
	canClean := buildInfo.Exists

	return ui.GenerateMenuRows(
		a.projectState.GetGeneratorLabel(),
		a.projectState.Configuration,
		canOpenIDE,
		canClean,
	)
}
