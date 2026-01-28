# Sprint 6 - TIT Compliance Refactoring Kickoff

**Objective:** Incrementally refactor cake to match TIT architecture, with verification at each step.

**Approach:** 6 phases, each tested independently before proceeding.

---

## Phase 1: SSOT - Extract Build Path Logic

**Goal:** Eliminate duplicated build directory construction (5 copies → 1)

### Changes
**File:** `internal/state/project.go`

Add method to ProjectState:
```go
// GetBuildDirectory returns the build directory path for given generator and config
func (ps *ProjectState) GetBuildDirectory(generatorName string, config string) string {
    for _, gen := range ps.AvailableGenerators {
        if gen.Name == generatorName {
            if gen.IsMultiConfig {
                return filepath.Join(ps.WorkingDirectory, "Builds", generatorName)
            }
            return filepath.Join(ps.WorkingDirectory, "Builds", generatorName, config)
        }
    }
    return ""
}

// IsGeneratorMultiConfig returns true if generator is multi-config
func (ps *ProjectState) IsGeneratorMultiConfig(generatorName string) bool {
    for _, gen := range ps.AvailableGenerators {
        if gen.Name == generatorName {
            return gen.IsMultiConfig
        }
    }
    return false
}
```

**Files to update:**
- `internal/app/app.go:697-730` - cmdGenerateProject
- `internal/app/app.go:809-815` - cmdBuildProject  
- `internal/app/app.go:855-861` - cmdCleanProject
- `internal/ops/setup.go:23-31` - ExecuteSetupProject
- `internal/ops/build.go:17-29` - ExecuteBuildProject
- `internal/ops/clean.go:14-22` - ExecuteCleanProject
- `internal/ops/open.go:16-46` - ExecuteOpenIDE

### Verification Steps
1. Build cake: `./build.sh`
2. Run cake in a CMake project directory
3. Test Generate operation - should create build in correct directory
4. Test Build operation - should find and build correct directory
5. Test Clean operation - should clean correct directory
6. Test Open IDE - should open correct build directory
7. Verify all paths identical to before refactoring

### Success Criteria
- [ ] All 7 locations use `GetBuildDirectory()` method
- [ ] No duplicated filepath.Join logic remains
- [ ] All operations work correctly
- [ ] Build directories created in same locations as before

---

## Phase 2: SSOT - Extract isMultiConfig Detection

**Goal:** Eliminate duplicated isMultiConfig loop (3 copies → 1)

### Changes
**Already done in Phase 1** - `IsGeneratorMultiConfig()` method added.

**Files to update:**
- `internal/app/app.go:709-715` - Replace with `a.projectState.IsGeneratorMultiConfig(generator)`
- `internal/app/app.go:809-815` - Replace with method call
- `internal/app/app.go:855-861` - Replace with method call

### Verification Steps
1. Build cake: `./build.sh`
2. Run cake with Xcode generator selected
3. Verify isMultiConfig=true (no config subdirectory)
4. Switch to Ninja generator
5. Verify isMultiConfig=false (config subdirectory created)
6. Test Generate/Build/Clean with both generator types

### Success Criteria
- [ ] All 3 locations use `IsGeneratorMultiConfig()` method
- [ ] No duplicated for-loop logic remains
- [ ] Multi-config generators (Xcode, VS) work correctly
- [ ] Single-config generators (Ninja) work correctly

---

## Phase 3: TIT Compliance - Create Footer Renderer

**Goal:** Match TIT footer pattern (string field → render function)

### Changes
**New File:** `internal/ui/footer.go`

Copy from TIT: `tit/internal/ui/footer.go`
Adapt for cake's simpler needs (no complex footer states).

```go
package ui

import (
    "strings"
    "github.com/charmbracelet/lipgloss"
)

// RenderFooter renders the footer hint with proper styling
func RenderFooter(hint string, theme Theme, width int) string {
    if hint == "" {
        hint = " "
    }
    
    footerStyle := lipgloss.NewStyle().
        Foreground(lipgloss.Color(theme.FooterTextColor)).
        Background(lipgloss.Color(theme.MainBackgroundColor)).
        Width(width).
        Align(lipgloss.Center)
    
    return footerStyle.Render(hint)
}
```

**Update:** `internal/app/app.go`
- Replace `footerHint string` field usage with `RenderFooter()` call in `View()`
- Keep `footerHint` field for state, but render through function

### Verification Steps
1. Build cake: `./build.sh`
2. Run cake
3. Navigate through menu items
4. Verify footer hints display correctly at bottom
5. Verify footer styling matches theme
6. Test all modes (Menu, Console, Preferences)

### Success Criteria
- [ ] Footer renders through `RenderFooter()` function
- [ ] Footer displays at bottom of screen
- [ ] Footer shows correct hints for each mode
- [ ] Footer styling matches theme colors

---

## Phase 4: TIT Compliance - Extract Async State

**Goal:** Split Application god object - extract AsyncState

### Changes
**New File:** `internal/app/async_state.go`

```go
package app

// AsyncState tracks async operation state
type AsyncState struct {
    OperationActive  bool
    OperationAborted bool
    ExitAllowed      bool
}

func NewAsyncState() *AsyncState {
    return &AsyncState{
        OperationActive:  false,
        OperationAborted: false,
        ExitAllowed:      false,
    }
}

func (as *AsyncState) Start() {
    as.OperationActive = true
    as.OperationAborted = false
    as.ExitAllowed = false
}

func (as *AsyncState) End() {
    as.OperationActive = false
}

func (as *AsyncState) Abort() {
    as.OperationAborted = true
}

func (as *AsyncState) ClearAborted() {
    as.OperationAborted = false
}

func (as *AsyncState) IsActive() bool {
    return as.OperationActive
}

func (as *AsyncState) IsAborted() bool {
    return as.OperationAborted
}

func (as *AsyncState) CanExit() bool {
    return as.ExitAllowed
}

func (as *AsyncState) SetExitAllowed(allowed bool) {
    as.ExitAllowed = allowed
}
```

**Update:** `internal/app/app.go`
- Replace fields: `asyncOperationActive`, `asyncOperationAborted`
- Add field: `asyncState *AsyncState`
- Update all references to use `a.asyncState.IsActive()` etc.

### Verification Steps
1. Build cake: `./build.sh`
2. Run cake
3. Start Generate operation
4. Verify async state tracking works (footer shows "Setting up...")
5. Try to quit during operation (should show "wait" message)
6. Let operation complete
7. Start Build operation
8. Press ESC during build
9. Verify abort works correctly

### Success Criteria
- [ ] AsyncState struct extracted to separate file
- [ ] Application uses AsyncState instead of raw bools
- [ ] Generate/Build/Clean operations work
- [ ] Quit during operation shows wait message
- [ ] ESC during operation works

---

## Phase 5: TIT Compliance - Add Message Dispatchers

**Goal:** Replace giant Update() switch with dispatcher pattern

### Changes
**New File:** `internal/app/dispatchers.go`

```go
package app

import tea "github.com/charmbracelet/bubbletea"

// MessageHandler handles specific message types
type MessageHandler interface {
    CanHandle(msg tea.Msg) bool
    Handle(a *Application, msg tea.Msg) (tea.Model, tea.Cmd)
}

// WindowSizeHandler handles terminal resize
type WindowSizeHandler struct{}

func (h WindowSizeHandler) CanHandle(msg tea.Msg) bool {
    _, ok := msg.(tea.WindowSizeMsg)
    return ok
}

func (h WindowSizeHandler) Handle(a *Application, msg tea.Msg) (tea.Model, tea.Cmd) {
    m := msg.(tea.WindowSizeMsg)
    a.width = m.Width
    a.height = m.Height
    a.sizing = ui.CalculateDynamicSizing(m.Width, m.Height)
    return a, nil
}

// KeyHandler routes keyboard input
type KeyHandler struct {
    handlers map[AppMode]func(*Application, tea.KeyMsg) (tea.Model, tea.Cmd)
}

func NewKeyHandler() *KeyHandler {
    return &KeyHandler{
        handlers: make(map[AppMode]func(*Application, tea.KeyMsg) (tea.Model, tea.Cmd)),
    }
}

func (h *KeyHandler) Register(mode AppMode, handler func(*Application, tea.KeyMsg) (tea.Model, tea.Cmd)) {
    h.handlers[mode] = handler
}

func (h *KeyHandler) CanHandle(msg tea.Msg) bool {
    _, ok := msg.(tea.KeyMsg)
    return ok
}

func (h *KeyHandler) Handle(a *Application, msg tea.Msg) (tea.Model, tea.Cmd) {
    keyMsg := msg.(tea.KeyMsg)
    if handler, ok := h.handlers[a.mode]; ok {
        return handler(a, keyMsg)
    }
    return a, nil
}
```

**Update:** `internal/app/app.go`
- Replace giant `Update()` switch with dispatcher registration
- Register handlers in `NewApplication()` or `Init()`

### Verification Steps
1. Build cake: `./build.sh`
2. Run cake
3. Test all keyboard shortcuts (↑↓, Enter, g, o, b, c, /, ESC, Ctrl+C)
4. Resize terminal window
5. Verify layout adjusts correctly
6. Test mode switching (Menu ↔ Preferences ↔ Console)
7. Test confirmation dialogs (Y/N/Enter/ESC)

### Success Criteria
- [ ] Message dispatchers file created
- [ ] Update() method simplified to dispatch calls
- [ ] All keyboard input works
- [ ] Window resize works
- [ ] All modes work
- [ ] Confirmation dialogs work

---

## Phase 6: TIT Compliance - Extract Operation Handlers

**Goal:** Move operations from app.go to separate op_*.go files

### Changes
**New Files:**
- `internal/app/op_generate.go` - Generate/Regenerate operations
- `internal/app/op_build.go` - Build operations  
- `internal/app/op_clean.go` - Clean operations
- `internal/app/op_open.go` - Open IDE/Editor operations

**Example:** `internal/app/op_generate.go`
```go
package app

import (
    "cake/internal/ops"
    "cake/internal/ui"
    tea "github.com/charmbracelet/bubbletea"
)

// startGenerateOperation begins the generate/regenerate operation
func (a *Application) startGenerateOperation() (tea.Model, tea.Cmd) {
    a.mode = ModeConsole
    a.asyncState.Start()
    a.outputBuffer.Clear()
    a.footerHint = GetFooterMessageText(MessageSetupInProgress)
    return a, a.cmdGenerateProject()
}

// cmdGenerateProject executes the generate/regenerate command
func (a *Application) cmdGenerateProject() tea.Cmd {
    return func() tea.Msg {
        outputCallback := func(line string) {
            a.outputBuffer.Append(line, ui.TypeStdout)
        }
        
        generator := a.projectState.SelectedGenerator
        config := a.projectState.Configuration
        projectRoot := a.projectState.WorkingDirectory
        isMultiConfig := a.projectState.IsGeneratorMultiConfig(generator)
        
        result := ops.ExecuteSetupProject(
            projectRoot,
            generator,
            config,
            isMultiConfig,
            outputCallback,
        )
        
        return GenerateCompleteMsg{
            Success: result.Success,
            Error:   result.Error,
        }
    }
}
```

**Update:** `internal/app/app.go`
- Remove operation methods (moved to new files)
- Keep only coordination logic

### Verification Steps
1. Build cake: `./build.sh`
2. Run cake in CMake project
3. Test Generate operation - should complete successfully
4. Test Build operation - should build project
5. Test Clean operation - should clean build directory
6. Test Open IDE - should open Xcode/VS
7. Test Open Editor - should open Neovim
8. Test all operations with different generators (Xcode, Ninja)

### Success Criteria
- [ ] All operation files created
- [ ] app.go reduced in size (operations moved out)
- [ ] Generate works
- [ ] Build works
- [ ] Clean works
- [ ] Open IDE works
- [ ] Open Editor works

---

## Verification Summary

### Phase 1-2: SSOT Fixes
- [ ] Build path logic centralized
- [ ] isMultiConfig detection centralized
- [ ] All operations work correctly
- [ ] No code duplication remains

### Phase 3: Footer Renderer
- [ ] Footer renders through function
- [ ] Matches TIT pattern
- [ ] All modes show correct footer

### Phase 4: Async State
- [ ] AsyncState extracted
- [ ] Application uses new struct
- [ ] Operations track state correctly
- [ ] Quit/ESC during operation works

### Phase 5: Message Dispatchers
- [ ] Dispatchers file created
- [ ] Update() simplified
- [ ] All input works
- [ ] Window resize works

### Phase 6: Operation Handlers
- [ ] All op_*.go files created
- [ ] app.go cleaned up
- [ ] All operations work
- [ ] Code organized by function

---

## Success Metrics

**Before:**
- TIT Compliance: 67%
- Code duplication: 5+ locations
- Application struct: 877 lines, 24 fields
- Update() method: 73 lines

**After:**
- TIT Compliance: 90%+
- Code duplication: 0 locations
- Application struct: <500 lines, <15 fields
- Update() method: <20 lines

---

**COUNSELOR:** OpenCode (glm-4.7)  
**Date:** 2026-01-28  
**Kickoff for:** ENGINEER
