# Sprint 10 - Console TIT Alignment Kickoff

**Role:** ENGINEER  
**Objective:** Make Cake Console exactly match TIT Console  
**Date:** 2026-01-29  

---

## Goal

Implement Cake Console to be 100% identical to TIT Console:
- Same visual structure (title, content, status bar)
- Same color scheme (SSOT from theme)
- Same scroll behavior (auto-scroll, manual scroll)
- Same status bar format

**Only difference:** Cake renders within header/footer layout (content area), TIT renders full terminal height.

---

## Files to Modify

### UI Layer
1. `internal/ui/console.go` - Console rendering
2. `internal/ui/buffer.go` - Output buffer (verify)

### App Layer
3. `internal/app/footer.go` - Console footer
4. `internal/app/app.go` - Console mode integration, confirmation dialogs
5. `internal/app/op_generate.go` - Use typed callback, streaming
6. `internal/app/op_build.go` - Use typed callback, streaming
7. `internal/app/op_clean.go` - Use typed callback, simple flow
8. `internal/app/op_open.go` - Use typed callback
9. `internal/app/op_regenerate.go` - NEW: Clean + Generate sequence
10. `internal/app/messages.go` - Add RegenerateCompleteMsg

### Ops Layer
11. `internal/ops/setup.go` - Change callback + streaming output
12. `internal/ops/build.go` - Change callback + streaming output
13. `internal/ops/clean.go` - Change callback + simple rm -rf
14. `internal/ops/open.go` - Change callback

### State Layer
15. `internal/state/project.go` - Remove `IsMultiConfig` field entirely (all projects are multi-config)

### Menu Layer  
16. `internal/ui/menu.go` - Change "Generator" label to "Project"

---

## Implementation Tasks

**⚠️ CRITICAL: These fixes were previously implemented in Sprint 9 but were REVERTED. They MUST be re-applied.**

### Task 0: CRITICAL FIXES - Menu Index Bug & Generator→Project (HIGHEST PRIORITY)

**Status:** These fixes were implemented in Sprint 9 by SURGEON but were REVERTED by subsequent changes. They are CRITICAL for basic functionality.

**Problem 1:** Menu shows "Generator: Xcode" - WRONG. Should be "Project: Xcode". Xcode/Ninja are the PROJECT type, not "generators".

**Problem 2:** Menu index calculation is BROKEN. `GetVisibleIndex()` and `GetArrayIndex()` only check `row.Visible` but should check `row.Visible && row.IsSelectable`. The separator row (visible but not selectable) causes off-by-one errors - selecting Clean triggers Build, Configuration toggle doesn't work.

**Files to fix:**

**1. `internal/ui/menu.go`:**
```go
// Line 10: Update comment
// Fixed 7 rows: [0]Project [1]Regenerate [2]OpenIDE [3]Separator [4]Configuration [5]Build [6]Clean

// Line 15: Update comment
Label        string // "Project", "Regenerate", "Open IDE", "", "Configuration", "Build", "Clean"

// Line 26-36: Change row ID and Label
{
    ID:           "project",  // Changed from "generator"
    Shortcut:     "",
    Emoji:        "⚙️",
    Label:        "Project",  // Changed from "Generator"
    Value:        projectLabel,
    Visible:      true,
    IsAction:     false,
    IsSelectable: true,
    Hint:         "Select project type (Xcode, Ninja, etc.)",
},
```

**2. `internal/app/app.go` - Fix GetVisibleIndex():**
```go
func (a *Application) GetVisibleIndex(rowID string) int {
    visibleIndex := 0
    for _, row := range a.menuItems {
        if row.ID == rowID {
            if row.Visible && row.IsSelectable {  // <-- ADDED: && row.IsSelectable
                return visibleIndex
            }
            return -1
        }
        if row.Visible && row.IsSelectable {  // <-- ADDED: && row.IsSelectable
            visibleIndex++
        }
    }
    return -1
}
```

**3. `internal/app/app.go` - Fix GetArrayIndex():**
```go
func (a *Application) GetArrayIndex(visibleIdx int) int {
    visibleCount := 0
    for i, row := range a.menuItems {
        if row.Visible && row.IsSelectable {  // <-- ADDED: && row.IsSelectable
            if visibleCount == visibleIdx {
                return i
            }
            visibleCount++
        }
    }
    return -1
}
```

**4. `internal/app/app.go` - Update executeRowAction case:**
```go
// Change line 417-422 from:
case "generator":
    // Cycle to next project
    a.projectState.CycleToNextGenerator()

// To:
case "project":
    // Cycle to next project
    a.projectState.CycleToNextGenerator()
```

**5. `internal/app/menu.go` - Update function call:**
```go
// Change:
return ui.GenerateMenuRows(
    a.projectState.GetGeneratorLabel(),  // <-- Change to GetProjectLabel()
    ...
)

// To:
return ui.GenerateMenuRows(
    a.projectState.GetProjectLabel(),
    ...
)
```

**6. `internal/state/project.go` - Rename method:**
```go
// Change method name from:
func (ps *ProjectState) GetGeneratorLabel() string

// To:
func (ps *ProjectState) GetProjectLabel() string
```

**7. `internal/ui/menu.go` - Update parameter name:**
```go
// Change function signature from:
func GenerateMenuRows(projectLabel string, ...)

// To:
func GenerateMenuRows(projectLabel string, ...)
```

---



**Action:** Replace Cake's `internal/ui/console.go` with TIT's version.

**Steps:**
1. Delete `/Users/jreng/Documents/Poems/dev/cake/internal/ui/console.go`
2. Copy `/Users/jreng/Documents/Poems/dev/tit/internal/ui/console.go` to `/Users/jreng/Documents/Poems/dev/cake/internal/ui/console.go`
3. Rename function from `RenderConsoleOutputFullScreen` to `RenderConsoleOutput`
4. Change parameters from `(termWidth, termHeight)` to `(maxWidth, totalHeight)`

**DON'T MODIFY ANY OTHER LOGIC.**

---

### Task 2: COPY TIT buffer.go EXACTLY (CRITICAL)

**Action:** Replace Cake's `internal/ui/buffer.go` with TIT's version.

**Steps:**
1. Delete `/Users/jreng/Documents/Poems/dev/cake/internal/ui/buffer.go`
2. Copy `/Users/jreng/Documents/Poems/dev/tit/internal/ui/buffer.go` to `/Users/jreng/Documents/Poems/dev/cake/internal/ui/buffer.go`

**NO CHANGES. COPY EXACTLY.**

---

### Task 3: Update Function Name and Parameters (console.go)

**Change:**
```go
func RenderConsoleOutputFullScreen(
    ...
    termWidth int,
    termHeight int,
    ...
)
```

**To:**
```go
func RenderConsoleOutput(
    ...
    maxWidth int,
    totalHeight int,
    ...
)
```

**ONLY these changes. NOTHING ELSE.**

**Function signature:**
```go
func RenderConsoleOutput(
    state *ConsoleOutState,
    buffer *OutputBuffer,
    palette Theme,
    maxWidth int,
    totalHeight int,
    operationInProgress bool,
    abortConfirmActive bool,
    autoScroll bool,
) string
```

**Structure (totalHeight - 2 for outer border):**
```
title (1) + blank (1) + content (N) + blank (1) + status (1) = totalHeight - 2
contentLines = totalHeight - 4
```

**Color mapping (SSOT from theme):**
```go
getColor := func(lineType OutputLineType) string {
    switch lineType {
    case TypeStdout, TypeCommand:
        return palette.OutputStdoutColor
    case TypeStderr:
        return palette.OutputStderrColor
    case TypeStatus:
        return palette.OutputStatusColor
    case TypeWarning:
        return palette.OutputWarningColor
    case TypeDebug:
        return palette.OutputDebugColor
    case TypeInfo:
        return palette.OutputInfoColor
    default:
        return palette.OutputStdoutColor
    }
}
```

**Title:**
- Text: "OUTPUT"
- Style: Bold, `OutputInfoColor`
- Full width, left-aligned

**Content:**
- If buffer empty: "(no output yet)" in `DimmedTextColor` + italic
- Format: `[HH:MM:SS] <text>`
- Apply color based on line type
- Wrap long lines

**Scroll logic:**
```go
totalOutputLines := len(allOutputLines)
maxScroll := totalOutputLines - contentLines
if maxScroll < 0 { maxScroll = 0 }
state.MaxScroll = maxScroll

if autoScroll {
    state.ScrollOffset = maxScroll
} else {
    // Clamp to valid range
    if state.ScrollOffset > maxScroll { state.ScrollOffset = maxScroll }
    if state.ScrollOffset < 0 { state.ScrollOffset = 0 }
}
```

**Status bar:**
- Left: Shortcuts with `AccentTextColor` (bold) + `LabelTextColor`
- Separator: `│` with `DimmedTextColor`
- Right: Scroll position
  - `(at bottom)` when at end
  - `↓ N more lines` when content below
  - `(can scroll up)` when at top

**Return:** Pre-sized panel (no outer border - RenderLayout adds it)

---

### Task 4: Implement Console Footer (footer.go)

**getConsoleFooter:**
```go
func (a *Application) getConsoleFooter(width int) string {
    var hintKey string
    if a.asyncState.IsActive() {
        hintKey = "console_running"
    } else {
        hintKey = "console_complete"
    }
    
    shortcuts := FooterHintShortcuts[hintKey]
    rightContent := a.computeConsoleScrollStatus()
    
    return ui.RenderFooter(shortcuts, width, &a.theme, rightContent)
}
```

**computeConsoleScrollStatus:**
```go
func (a *Application) computeConsoleScrollStatus() string {
    state := &a.consoleState
    
    if state.MaxScroll <= 0 {
        return ""
    }
    
    atBottom := state.ScrollOffset >= state.MaxScroll
    remainingLines := state.MaxScroll - state.ScrollOffset
    
    if atBottom {
        return "(at bottom)"
    }
    if remainingLines > 0 {
        return fmt.Sprintf("↓ %d more", remainingLines)
    }
    return "(can scroll up)"
}
```

---

### Task 5: Implement Console Mode Integration (app.go)

**renderConsoleMode:**
```go
func (a *Application) renderConsoleMode() string {
    return ui.RenderConsoleOutput(
        &a.consoleState,
        a.outputBuffer,
        a.theme,
        a.sizing.ContentInnerWidth,
        a.sizing.ContentHeight,
        a.asyncState.IsActive(),
        false, // abortConfirmActive
        a.asyncState.IsActive(), // autoScroll when active
    )
}
```

**handleOperationKeyPress:**
```go
func (a *Application) handleOperationKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
    switch msg.String() {
    case "up":
        a.consoleState.ScrollUp()
        return a, nil
    case "down":
        a.consoleState.ScrollDown()
        return a, nil
    case "esc":
        if !a.asyncState.IsActive() {
            a.mode = ModeMenu
            a.selectedIndex = 0
            a.menuItems = a.GenerateMenu()
            a.footerHint = FooterHints["menu_navigate"]
        }
        return a, nil
    case "ctrl+c":
        return a.handleCtrlC()
    }
    return a, nil
}
```

---

## Task 7: Remove IsMultiConfig - ALL Projects Are Multi-Config (CRITICAL)

**Problem:** Code still has `IsMultiConfig` field and conditional logic. ALL projects are multi-config.

**Solution:** Remove `IsMultiConfig` field entirely from all structs and functions.

### Files to modify:

**1. `internal/state/project.go`:**
- Remove `IsMultiConfig bool` field from `Generator` struct
- Remove all `IsMultiConfig: true/false` assignments
- Remove `IsGeneratorMultiConfig()` method
- Remove `isMultiConfigGenerator()` method
- Keep only "Ninja" project type (no "Ninja Multi-Config" separate entry)

**2. `internal/ops/setup.go`:**
- Remove `isMultiConfig` parameter from `ExecuteSetupProject()`
- Use simple path: `buildDir := filepath.Join(workingDir, "Builds", project)`

**3. `internal/ops/build.go`:**
- Remove `isMultiConfig` parameter from `ExecuteBuildProject()`
- Always use: `args := []string{"--build", buildDir, "--config", config}`

**4. `internal/ops/clean.go`:**
- Remove `isMultiConfig` parameter from `ExecuteCleanProject()`
- Use simple path: `buildDir := filepath.Join(projectRoot, "Builds", project)`

**5. `internal/app/op_*.go`:**
- Remove all `isMultiConfig` variable declarations and usages
- Remove calls to `a.projectState.IsGeneratorMultiConfig()`

---

## Task 8: Fix Output Color Types (CRITICAL)

**Problem:** All output is currently `TypeStdout` (gray). TIT uses different colors for different output types.

**Solution:** Change callback signature from `func(string)` to `func(string, ui.OutputLineType)`

### Step 1: Update ops/setup.go

Change callback signature:
```go
func ExecuteSetupProject(
    workingDir, project, config string,
    outputCallback func(string, ui.OutputLineType),  // <-- Changed
) SetupResult
```

Update all outputCallback calls:
```go
outputCallback("Running: cmake " + strings.Join(args, " "), ui.TypeInfo)
outputCallback("", ui.TypeStdout)  // Empty lines as stdout

// Command output
for _, line := range strings.Split(string(output), "\n") {
    if line != "" {
        outputCallback(line, ui.TypeStdout)
    }
}

// Error
if err != nil {
    outputCallback("", ui.TypeStdout)
    outputCallback("ERROR: " + err.Error(), ui.TypeStderr)
    return SetupResult{Success: false, Error: err.Error()}
}

// Success
outputCallback("", ui.TypeStdout)
outputCallback("Setup completed successfully: " + buildDir, ui.TypeStatus)
```

### Step 2: Update ops/build.go

Same pattern - change callback signature and use appropriate types:
- Info messages: `ui.TypeInfo`
- Command output: `ui.TypeStdout`
- Errors: `ui.TypeStderr`
- Success: `ui.TypeStatus`

### Step 3: Update ops/clean.go

Same pattern.

### Step 4: Update ops/open.go

Same pattern.

### Step 5: Update Operation Handlers (op_*.go)

Update the callback creation:
```go
// OLD (all gray):
outputCallback := func(line string) {
    a.outputBuffer.Append(line, ui.TypeStdout)
}

// NEW (color-coded):
outputCallback := func(line string, lineType ui.OutputLineType) {
    a.outputBuffer.Append(line, lineType)
}
```

---

## Task 9: Fix Real-Time Console Streaming (CRITICAL)

**Problem:** Console updates are slow because `cmd.CombinedOutput()` waits for entire command to complete.

**Solution:** Stream output line-by-line using `StdoutPipe()` and `bufio.Scanner` (TIT pattern).

### Pattern for Real-Time Streaming:

```go
func ExecuteBuildProject(...) BuildResult {
    // ... setup code ...
    
    outputCallback("Building: " + buildDir, ui.TypeInfo)
    outputCallback("", ui.TypeStdout)
    
    cmd := exec.Command("cmake", args...)
    cmd.Dir = projectRoot
    
    // Get stdout pipe for streaming
    stdout, err := cmd.StdoutPipe()
    if err != nil {
        outputCallback("ERROR: Failed to create stdout pipe", ui.TypeStderr)
        return BuildResult{Success: false, Error: err.Error()}
    }
    
    // Get stderr pipe for streaming
    stderr, err := cmd.StderrPipe()
    if err != nil {
        outputCallback("ERROR: Failed to create stderr pipe", ui.TypeStderr)
        return BuildResult{Success: false, Error: err.Error()}
    }
    
    // Start command
    if err := cmd.Start(); err != nil {
        outputCallback("ERROR: Failed to start command", ui.TypeStderr)
        return BuildResult{Success: false, Error: err.Error()}
    }
    
    // Stream stdout in goroutine
    go func() {
        scanner := bufio.NewScanner(stdout)
        for scanner.Scan() {
            line := scanner.Text()
            if line != "" {
                outputCallback(line, ui.TypeStdout)
            }
        }
    }()
    
    // Stream stderr in goroutine
    go func() {
        scanner := bufio.NewScanner(stderr)
        for scanner.Scan() {
            line := scanner.Text()
            if line != "" {
                outputCallback(line, ui.TypeStderr)
            }
        }
    }()
    
    // Wait for command to complete
    err = cmd.Wait()
    
    if err != nil {
        outputCallback("", ui.TypeStdout)
        outputCallback("ERROR: Build failed", ui.TypeStderr)
        return BuildResult{Success: false, Error: err.Error()}
    }
    
    outputCallback("", ui.TypeStdout)
    outputCallback("Build completed successfully", ui.TypeStatus)
    return BuildResult{Success: true}
}
```

**Apply this pattern to:**
- `ops/setup.go` - ExecuteSetupProject
- `ops/build.go` - ExecuteBuildProject
- `ops/clean.go` - ExecuteCleanProject (simpler - just rm -rf, no streaming needed)

---

## Task 10: Fix Clean Operation Flow (CRITICAL)

**Problem:** Clean is showing "Building" messages and taking forever.

**Root Cause:** Clean is synchronous and doesn't have confirmation dialog.

**Solution:** 

### Step 1: Clean should be simple and fast

```go
func ExecuteCleanProject(project, config string, projectRoot string, outputCallback func(string, ui.OutputLineType)) CleanResult {
    // ALL projects are multi-config - use simple path
    buildDir := filepath.Join(projectRoot, "Builds", project)
    
    outputCallback("Cleaning build directory...", ui.TypeInfo)
    outputCallback("Target: " + buildDir, ui.TypeStdout)
    outputCallback("", ui.TypeStdout)
    
    // Check if directory exists
    if _, err := os.Stat(buildDir); os.IsNotExist(err) {
        outputCallback("Build directory does not exist", ui.TypeWarning)
        outputCallback("Nothing to clean", ui.TypeStatus)
        return CleanResult{Success: true}
    }
    
    // Simple, fast delete
    err := os.RemoveAll(buildDir)
    if err != nil {
        outputCallback("ERROR: Failed to remove directory", ui.TypeStderr)
        outputCallback(err.Error(), ui.TypeStderr)
        return CleanResult{Success: false, Error: err.Error()}
    }
    
    outputCallback("Successfully removed: " + buildDir, ui.TypeStatus)
    return CleanResult{Success: true}
}
```

### Step 2: Add Confirmation Dialog for Clean

In `app.go` `executeRowAction`:

```go
case "clean":
    // Show confirmation dialog first
    a.confirmDialog = ui.NewConfirmationDialog(
        "Clean Build Directory",
        "Remove all build artifacts?",
        ui.ButtonNo,  // Default to No for safety
    )
    a.confirmDialog.Active = true
    a.pendingOperation = "clean"
    return true, nil
```

In `handleConfirmDialogKeyPress`:

```go
case "clean":
    a.confirmDialog.Active = false
    a.confirmDialog = nil
    a.pendingOperation = ""
    _, cmd := a.startCleanOperation()
    return a, cmd
```

---

## Task 10: Fix Regenerate Flow (CRITICAL)

**Problem:** Regenerate takes forever.

**Expected Behavior:** Regenerate = Clean + Generate (fast)

**Solution:** Regenerate should:
1. Clean first (fast rm -rf)
2. Then run Generate

In `executeRowAction`:

```go
case "regenerate":
    // Show confirmation dialog
    a.confirmDialog = ui.NewConfirmationDialog(
        "Regenerate Project",
        "Clean and re-run CMake configuration?",
        ui.ButtonNo,
    )
    a.confirmDialog.Active = true
    a.pendingOperation = "regenerate"
    return true, nil
```

In `handleConfirmDialogKeyPress`:

```go
case "regenerate":
    a.confirmDialog.Active = false
    a.confirmDialog = nil
    a.pendingOperation = ""
    // Start regenerate sequence: Clean then Generate
    return a.startRegenerateOperation()
```

Create `op_regenerate.go`:

```go
func (a *Application) startRegenerateOperation() (tea.Model, tea.Cmd) {
    a.mode = ModeConsole
    a.asyncState.Start()
    a.outputBuffer.Clear()
    a.footerHint = "Regenerating project..."
    return a, a.cmdRegenerateProject()
}

func (a *Application) cmdRegenerateProject() tea.Cmd {
    return func() tea.Msg {
        outputCallback := func(line string, lineType ui.OutputLineType) {
            a.outputBuffer.Append(line, lineType)
        }
        
        project := a.projectState.SelectedGenerator
        config := a.projectState.Configuration
        projectRoot := a.projectState.WorkingDirectory
        // ALL projects are multi-config - no need to check
        
        // Step 1: Clean
        outputCallback("=== Step 1: Clean ===", ui.TypeInfo)
        // ALL projects are multi-config - use simple path
        buildDir := filepath.Join(projectRoot, "Builds", project)
        
        if _, err := os.Stat(buildDir); err == nil {
            outputCallback("Removing: " + buildDir, ui.TypeStdout)
            if err := os.RemoveAll(buildDir); err != nil {
                outputCallback("ERROR: Clean failed", ui.TypeStderr)
                return RegenerateCompleteMsg{Success: false, Error: err.Error()}
            }
            outputCallback("Clean completed", ui.TypeStatus)
        } else {
            outputCallback("No build directory to clean", ui.TypeWarning)
        }
        
        outputCallback("", ui.TypeStdout)
        outputCallback("=== Step 2: Generate ===", ui.TypeInfo)
        
        // Step 2: Generate
        result := ops.ExecuteSetupProject(
            projectRoot,
            project,
            config,
            outputCallback,
        )
        
        return RegenerateCompleteMsg{
            Success: result.Success,
            Error:   result.Error,
        }
    }
}
```

Add to `messages.go`:

```go
type RegenerateCompleteMsg struct {
    Success bool
    Error   string
}
```

---

## Task 11: Fix Auto-Scroll Behavior (CRITICAL)

**Problem:** Console doesn't always scroll to bottom during operations.

**Solution:** Ensure `autoScroll=true` when operation is active.

In `renderConsoleMode` (already correct):

```go
func (a *Application) renderConsoleMode() string {
    return ui.RenderConsoleOutput(
        &a.consoleState,
        a.outputBuffer,
        a.theme,
        a.sizing.ContentInnerWidth,
        a.sizing.ContentHeight,
        a.asyncState.IsActive(),
        false, // abortConfirmActive
        a.asyncState.IsActive(), // autoScroll when active
    )
}
```

The `RenderConsoleOutput` function must:
1. When `autoScroll=true`: Set `state.ScrollOffset = maxScroll`
2. When `autoScroll=false`: Allow manual scroll

This is already in the spec - verify it's working correctly.

---

## Color Reference (Theme SSOT)

| Element | Theme Field | GFX Default |
|---------|-------------|-------------|
| Stdout text | OutputStdoutColor | #999999 |
| Stderr text | OutputStderrColor | #FC704C |
| Status text | OutputStatusColor | #01C2D2 |
| Warning text | OutputWarningColor | #F2AB53 |
| Debug text | OutputDebugColor | #33535B |
| Info text | OutputInfoColor | #01C2D2 |
| Title | OutputInfoColor | #01C2D2 |
| Empty hint | DimmedTextColor | #33535B |
| Shortcut keys | AccentTextColor | #01C2D2 |
| Descriptions | LabelTextColor | #8CC9D9 |
| Separators | DimmedTextColor | #33535B |

---

## Footer Hint Keys (messages.go)

Ensure these keys exist in FooterHintShortcuts:
```go
"console_running": []FooterHint{
    {Key: "↑↓", Description: "scroll"},
    {Key: "ESC", Description: "abort"},
},
"console_complete": []FooterHint{
    {Key: "↑↓", Description: "scroll"},
    {Key: "ESC", Description: "back to menu"},
},
```

---

## Build & Test

```bash
./build.sh
```

**Test scenarios:**
1. Run Generate operation - console should auto-scroll
2. Press ↑ during operation - should scroll up, auto-scroll pauses
3. Press ↓ - should scroll down
4. After operation completes, ESC returns to menu
5. All line types render with correct colors
6. Status bar shows correct scroll position
7. Empty buffer shows "(no output yet)"

---

## Success Criteria

### Visual
- [ ] Console renders with TIT-identical structure
- [ ] All colors from theme SSOT (no hardcoded)
- [ ] Commands/info appear in `OutputInfoColor` (cyan)
- [ ] Errors appear in `OutputStderrColor` (coral/red)
- [ ] Success messages appear in `OutputStatusColor` (cyan)
- [ ] Regular output appears in `OutputStdoutColor` (gray)

### Behavior
- [ ] Console updates in REAL-TIME (streaming, not batch)
- [ ] Auto-scroll keeps viewport at bottom during operations
- [ ] Manual scroll with ↑↓ keys works when auto-scroll paused
- [ ] Status bar shows shortcuts + scroll position
- [ ] ESC returns to menu when idle

### Menu
- [ ] **Label**: Menu shows "Project: Xcode" (not "Generator: Xcode")
- [ ] **Index Fix**: GetVisibleIndex/GetArrayIndex check `IsSelectable` (fixes Clean triggering Build)
- [ ] **Toggle**: Configuration row toggle works with Space/Enter

### Console
- [ ] **Line Limit**: Console doesn't flood with 8k+ lines (respects 1000 line buffer limit)
- [ ] **ESC Works**: ESC returns to menu when operation idle, aborts when running

### Operations
- [ ] **Clean**: Shows confirmation dialog → fast rm -rf → meaningful message
- [ ] **Regenerate**: Shows confirmation dialog → Clean → Generate (fast)
- [ ] **Generate**: Streams cmake output in real-time
- [ ] **Build**: Streams build output in real-time
- [ ] Clean does NOT show "Building" messages
- [ ] Regenerate completes quickly (not taking forever)

### Quality
- [ ] All 5 themes render correctly
- [ ] Build passes: `./build.sh`

---

**JRENG!**
