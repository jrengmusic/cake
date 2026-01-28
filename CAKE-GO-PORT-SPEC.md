# cake → Go + Bubble Tea + Lip Gloss Port Specification

**Version:** 1.1.0 (CORRECTED)  
**Date:** 2026-01-22  
**Status:** Analysis & Architecture Planning  
**Author:** OpenCode (ANALYST)

---

## Critical Corrections (v1.1.0)

1. **JUCE is NOT for audio plugins** - It's only used for cross-platform utilities (filesystem, process management). cake is pure CMake management tool.

2. **Console output uses EXACT TIT pattern** - NOT history pane. Reuse `internal/ui/console.go` directly with ConsoleOutState, OutputBuffer, real-time line streaming.

3. **Async operations required** - Operations (Setup, Build, Clean) run in background worker threads using Bubble Tea's Cmd pattern (identical to TIT's git operations).

4. **No quit menu** - Ctrl+C is handled at application level (exactly like TIT), with confirmation timeout. No "Exit" menu item.

---

## Executive Summary

**cake is a pure CMake management tool** that manages project setup, building, and cleanup via CLI. Porting it to Go + Bubble Tea + Lip Gloss is **low-to-medium complexity** due to:

✅ **Advantages:**
- **Simplicity:** No git state tracking, straightforward operations
- **Familiar patterns:** Menu system, async operations identical to TIT
- **Reusable components:** TIT's console output, async handling, banner/braille, state patterns
- **Stateless pages:** No reactive state like git—just linear flows (menu → operation → output → menu)
- **Clear scope:** 3 operations (Setup, Build, Clean)

⚠️ **Challenges:**
- **Cross-platform binary execution:** CMake, platform-specific generators (Xcode, Ninja, VS)
- **Real-time output streaming:** Build logs must scroll while process runs (async worker thread)
- **Project root detection:** Must work without .git
- **Generator discovery:** Platform-specific paths and detection logic

---

## Complexity Assessment

### Codebase Comparison

| Metric | cake (C++) | TIT (Go) | Notes |
|--------|-----------|---------|-------|
| **Total lines** | 3,569 | 15,065 | TIT ~4x larger due to git complexity |
| **Modules** | 10 | 7 internal packages | Similar structure |
| **Operations** | 3 (Setup, Build, Clean) | 7+ git operations | cake simpler: linear flows |
| **Async operations** | Yes (build thread) | Yes (git worker threads) | Identical pattern in both |
| **Menu system** | FTXUI (C++) | Bubble Tea (Go) | Direct port possible |

### Porting Effort Estimate

| Component | Complexity | Effort | Dependencies |
|-----------|-----------|--------|--------------|
| **Core App State** | 🟢 LOW | 2 days | None |
| **Menu System** | 🟢 LOW | 1-2 days | Reuse TIT |
| **Project State/Detection** | 🟡 MEDIUM | 2 days | App state |
| **Async Operation Framework** | 🟡 MEDIUM | 2-3 days | Reuse TIT pattern |
| **Console Output Display** | 🟢 LOW | 1 day | Reuse TIT console.go |
| **Setup Operation** | 🟡 MEDIUM | 2-3 days | Async framework |
| **Build Operation** | 🟡 MEDIUM | 2-3 days | Async framework |
| **Clean Operation** | 🟡 MEDIUM | 1-2 days | Async framework |
| **Cross-platform Generator Logic** | 🟡 MEDIUM | 2-3 days | Platform detection |
| **Banner Rendering** | 🟢 LOW | 1 day | Reuse TIT braille |
| **Testing/Integration** | 🟡 MEDIUM | 2-3 days | All platforms |

**Total Estimate:** 16-20 days (2.5-3 weeks) for complete MVP

---

## Architecture & Design

### Layered Design (Identical to TIT)

```
┌─────────────────────────────────────────────────────┐
│  Bubble Tea Application (app/app.go)                │
│  - Main event loop, screen management               │
│  - Ctrl+C handling at app level                     │
├─────────────────────────────────────────────────────┤
│  Modes / Pages (app/app.go)                         │
│  - ModeMain (menu), ModeSetup, ModeBuild, ModeClean │
│  - Console output display during operations         │
├─────────────────────────────────────────────────────┤
│  Async Operations (app/handlers.go)                 │
│  - cmdSetupProject() - CMake configuration          │
│  - cmdBuildProject() - Build execution              │
│  - cmdCleanProject() - Directory cleanup             │
│  - Worker threads via Bubble Tea Cmd pattern        │
├─────────────────────────────────────────────────────┤
│  State Management (app/app.go)                      │
│  - ProjectState: detection + caching                │
│  - asyncOperationActive: UI state during operation  │
│  - OutputBuffer: real-time line streaming           │
├─────────────────────────────────────────────────────┤
│  UI Components (internal/ui/*.go)                   │
│  - Console output (reuse TIT)                       │
│  - Menu rendering (reuse TIT)                       │
│  - Banner/braille (reuse TIT)                       │
└─────────────────────────────────────────────────────┘
```

### Key Design Principles

1. **Async Operation Pattern (Identical to TIT):**
   - Set `asyncOperationActive = true`
   - Return `Cmd` that runs operation in worker thread
   - Operation writes to `OutputBuffer` via callback
   - UI renders `console.RenderConsoleOutput()` each frame
   - On completion, set `asyncOperationActive = false`
   - ESC during operation aborts (not implemented for cake)

2. **Ctrl+C Handling (Application Level):**
   - No menu item for "Exit"
   - Ctrl+C at main menu: activate confirmation
   - Ctrl+C during operation: show "operation in progress" message
   - Second Ctrl+C: quit immediately
   - 2-second confirmation timeout (MessageCtrlCConfirm)

3. **Console Output (Reuse TIT Exactly):**
   - `ConsoleOutState` tracks scroll position
   - `OutputBuffer` holds lines with type+color info
   - `RenderConsoleOutput()` displays with wrapping, coloring
   - Real-time streaming as operation writes lines
   - Arrow up/down to scroll manually (auto-scroll enabled)

4. **No Menu-Based Navigation:**
   - Main menu → Select operation (Setup/Build/Clean)
   - Operation runs, shows output
   - On completion, return to main menu
   - Linear flow: menu → operation → menu

---

## Core Components

### 1. Application State (`app/app.go`)

```go
type Application struct {
    // Current mode/page
    Mode              AppMode
    
    // Project state
    ProjectState      *ProjectState
    
    // Async operation state
    asyncOperationActive  bool
    asyncOperationAborted bool
    
    // Console output
    OutputBuffer      *ui.OutputBuffer
    ConsoleState      ui.ConsoleOutState
    
    // Quit confirmation
    quitConfirmActive bool
    quitConfirmTime   time.Time
    
    // Menu selection
    selectedMenuIdx   int
    
    // UI state
    sizing            *ui.DynamicSizing
    theme             ui.Theme
    footerHint        string
}
```

### 2. AppMode Enumeration

```go
type AppMode int

const (
    ModeMain AppMode = iota
    ModeSetupChooseGen
    ModeBuild
    ModeClean
    ModeError
    // NO ModeExit - use Ctrl+C instead
)
```

### 3. Project State (`internal/state/project.go`)

```go
type ProjectState struct {
    ProjectRoot     string                  // CWD
    HasBuildDir     bool
    AvailableIDEs   map[string]string       // IDE name → build dir path
    SelectedGen     string                  // User's CMake generator choice
    
    LastRefreshTime time.Time
    RefreshTimeout  time.Duration
}

// Methods:
// Refresh() - conditional refresh if stale
// ForceRefresh() - immediate refresh (after operation)
// ShouldRefresh() bool
```

**CRITICAL:** No DAW detection. Ignore `isPluginProject`.

### 4. Async Operations Pattern (Reuse TIT)

**Setup Operation:**
```go
// In handlers.go
func (a *Application) cmdSetupProject(generator string) tea.Cmd {
    return func() tea.Msg {
        // WORKER THREAD
        result := executeSetupProject(generator)
        return SetupCompleteMsg{
            Success: result.Success,
            Error:   result.Error,
        }
    }
}

// In Update()
a.asyncOperationActive = true
return a, a.cmdSetupProject(generator)

// Handle completion
case SetupCompleteMsg:
    a.asyncOperationActive = false
    a.projectState.ForceRefresh() // Re-scan for build dirs
    // Return to main menu
```

**Build Operation:**
```go
// Similar pattern - executes cmake --build or platform-specific
func (a *Application) cmdBuildProject(buildDir string) tea.Cmd {
    return func() tea.Msg {
        // WORKER THREAD
        result := executeBuildProject(buildDir)
        return BuildCompleteMsg{...}
    }
}
```

### 5. Console Output (Reuse TIT Exactly)

```go
// TIT's OutputBuffer pattern - lines with type/color
type OutputLineType int
const (
    TypeStdout
    TypeStderr
    TypeCommand
    TypeStatus
    TypeWarning
    TypeInfo
)

// In operation handler:
outputCallback := func(line string) {
    a.OutputBuffer.Append(ui.OutputLine{
        Text: line,
        Type: ui.TypeStdout,
    })
}

// In Update() render:
return a.View() // which calls RenderConsoleOutput()
```

### 6. Menu System (Reuse TIT)

**Main Menu Items:**
1. Setup Project (choose generator)
2. Build Project (compile)
3. Clean Project (remove artifacts)

No Exit menu - use Ctrl+C instead.

---

## UI Components

### Reusable from TIT ♻️

| Component | Source | Usage | Changes |
|-----------|--------|-------|---------|
| **Console output** | `internal/ui/console.go` | Build/Setup/Clean output | No changes - exact reuse |
| **Menu rendering** | `internal/ui/menu.go` | Main menu, Setup generator menu | No changes - exact reuse |
| **Braille SVG** | `internal/banner/` | Convert cake SVG to braille | **DIRECT REUSE** - Pass cake's SVG to `SvgToBrailleArray()` |
| **Status bar** | `internal/ui/statusbar.go` | Footer hints | No changes - exact reuse |
| **Layout** | `internal/ui/layout.go` | Header/content/footer | **COPY pattern only** - Create `RenderBannerDynamic()` for cake |
| **Info rows** | `internal/ui/inforow.go` | Project info (CWD, Generator) | No changes - exact reuse |
| **Theme/colors** | `internal/app/theme.go` | Color palette | Reuse existing palette |

### Layout Structure

```
┌─────────────────────────────────────────────────────┐
│ HEADER: Project info + Banner                       │
│ ┌ CWD: /path/to/project                            │
│ └ Generator: Xcode                                  │
├─────────────────────────────────────────────────────┤
│ CONTENT:                                            │
│ - Menu mode: List of operations                     │
│ - Operation mode: Live console output               │
├─────────────────────────────────────────────────────┤
│ FOOTER: Keyboard hints + Ctrl+C confirmation        │
└─────────────────────────────────────────────────────┘
```

---

## Operation Execution Details

### Setup Operation

**Flow:**
1. Main menu → "Setup Project"
2. Display generator selection menu (platform-specific)
3. User selects generator (1-3 options)
4. Start async operation: `cmdSetupProject(generator)`
5. Run: `cmake -G "Generator" -S . -B build_dir`
6. Stream output to console in real-time
7. On completion: refresh ProjectState (re-detect build dirs)
8. Return to main menu

**Generators:**
- macOS: Xcode, Ninja
- Windows: Visual Studio, Ninja Multi-Config
- Linux: Ninja, Unix Makefiles

### Build Operation

**Flow:**
1. Main menu → "Build Project"
2. Check if build directory exists (from cached ProjectState)
3. If no build dir: show error, return to menu
4. Otherwise: start async operation `cmdBuildProject(buildDir)`
5. Run platform-specific build command:
   - Xcode: `xcodebuild -scheme KeepJUCEUpdated`
   - Ninja: `ninja -C build_dir`
   - Visual Studio: `cmake --build build_dir --config Release`
6. Stream output to console
7. On completion: show exit code + success/error message
8. Return to main menu (ESC or after confirmation)

### Clean Operation

**Flow:**
1. Main menu → "Clean Project"
2. List available build directories from cached IDEs
3. User selects directory to clean
4. Start async operation: `cmdCleanProject(buildDir)`
5. Run: `rm -rf buildDir` (cross-platform)
6. Stream output to console
7. On completion: refresh ProjectState
8. Return to main menu

---

## Cross-Platform Execution

### CMake Command Building

```go
func BuildCMakeCommand(generator string, projectRoot string) []string {
    buildDir := filepath.Join(projectRoot, "build_" + strings.ToLower(generator))
    
    args := []string{
        "-G", generator,
        "-S", projectRoot,
        "-B", buildDir,
    }
    return args
    // Full: cmake -G "Xcode" -S . -B build_xcode
}
```

### Build Command Building

```go
func BuildBuildCommand(generator string, buildDir string) (string, []string) {
    switch generator {
    case "Xcode":
        return "xcodebuild", []string{"-scheme", "ProjectName"}
    case "Ninja":
        return "ninja", []string{"-C", buildDir}
    case "Visual Studio":
        return "cmake", []string{"--build", buildDir, "--config", "Release"}
    default:
        return "make", []string{"-C", buildDir}
    }
}
```

### Platform Detection

```go
func GetPlatformGenerator() string {
    switch runtime.GOOS {
    case "darwin":
        return "Xcode"
    case "windows":
        return "Visual Studio"  // Detect version via detectVisualStudioGenerator()
    default:
        return "Ninja"
    }
}
```

---

## Error Handling & Validation

### Fail-Fast Rule (per SESSION-LOG.md CRITICAL)

**NEVER silent failures:**
- ❌ No `_ = cmd.Output()` (suppresses errors)
- ❌ No empty fallback strings
- ❌ No swallowing stderr
- ✅ Return explicit errors
- ✅ Log exact failure reason
- ✅ Show error on screen immediately

**Pattern:**
```go
if err != nil {
    return BuildCompleteMsg{
        Success: false,
        Error:   fmt.Sprintf("CMake failed: %v", err),
    }
}
```

### Validation

1. **Project root:** Must exist and be readable
2. **CMake:** Must be in PATH (or detect version on startup)
3. **Generator:** Validate against platform-specific list
4. **Build dir:** Check before running build
5. **Command execution:** Timeout after 30 minutes (configurable)

---

## File Structure

```
cake-go/
├── cmd/cake/
│   └── main.go                      # Entry point
├── internal/
│   ├── app/
│   │   ├── app.go                   # Application state + Update/View
│   │   ├── handlers.go              # Ctrl+C, menu selection, async cmds
│   │   ├── messages.go              # Msg types (SetupCompleteMsg, etc.)
│   │   └── theme.go                 # Colors (reuse from TIT)
│   ├── state/
│   │   └── project.go               # ProjectState detection + caching
│   ├── ops/
│   │   ├── setup.go                 # executeSetupProject()
│   │   ├── build.go                 # executeBuildProject()
│   │   └── clean.go                 # executeCleanProject()
│   ├── utils/
│   │   ├── project.go               # Project detection (no DAW)
│   │   ├── exec.go                  # Cross-platform command execution
│   │   ├── generator.go             # Generator detection + args
│   │   └── platform.go              # OS detection
│   ├── ui/                          # (REUSE FROM TIT)
│   │   ├── console.go               # RenderConsoleOutput()
│   │   ├── menu.go                  # Menu rendering
│   │   ├── styles.go                # Lip Gloss styles
│   │   ├── layout.go                # Overall layout
│   │   ├── statusbar.go             # Footer hints
│   │   ├── inforow.go               # Info display
│   │   └── types.go                 # OutputBuffer, ConsoleOutState
│   └── banner/                      # (REUSE FROM TIT)
│       └── render.go                # SVG banner rendering
├── go.mod
├── go.sum
└── build.sh
```

---

## Operation Messages (Bubble Tea Msg Types)

```go
// SetupCompleteMsg - after CMake setup finishes
type SetupCompleteMsg struct {
    Success bool
    Error   string
}

// BuildCompleteMsg - after build finishes
type BuildCompleteMsg struct {
    Success  bool
    ExitCode int
    Error    string
}

// CleanCompleteMsg - after clean finishes
type CleanCompleteMsg struct {
    Success bool
    Error   string
}

// TickMsg - for confirmation timeout
type TickMsg time.Time
```

---

## Success Criteria

### Functional
- ✅ Main menu displays all 3 operations
- ✅ Setup generates CMake build directory
- ✅ Build executes and shows real-time output
- ✅ Clean removes build artifacts
- ✅ Ctrl+C shows confirmation (first press)
- ✅ Ctrl+C quits (second press or timeout)
- ✅ All keyboard shortcuts work (arrows, numbers)
- ✅ Async operations don't block UI (smooth scrolling)
- ✅ ESC returns to menu from operation

### Code Quality
- ✅ No silent failures (all errors logged + displayed)
- ✅ No SSOT violations (single source of truth for state)
- ✅ Reusable components documented
- ✅ Cross-platform compatibility verified

### Performance
- ✅ Sub-100ms menu render time
- ✅ Smooth output scrolling during build
- ✅ No memory leaks on long operations

---

## Risks & Mitigations

| Risk | Impact | Mitigation |
|------|--------|-----------|
| **CMake path detection** | Setup fails if not found | Detect cmake in PATH on startup, error if missing |
| **Generator args quoting** | Build cmd syntax error | Test all platforms before release |
| **Output encoding on Windows** | Garbled logs | UTF-8 normalization, ANSI stripping |
| **Process hanging** | Ctrl+C doesn't work | Use context.Context with timeout |
| **Async operation blocking** | UI freezes during build | Ensure streaming callback is non-blocking |
| **IDE path discovery** | Can't find build dirs | Cache results, user can specify manually |

---

## Notes for Implementation

### CRITICAL: Async Pattern (Identical to TIT)

```go
// This is how git operations work in TIT - cake follows same pattern:

// 1. Start operation
a.asyncOperationActive = true
a.OutputBuffer.Reset()
return a, a.cmdBuildProject(buildDir)

// 2. In worker thread (Cmd)
func (a *Application) cmdBuildProject(buildDir string) tea.Cmd {
    return func() tea.Msg {
        // WORKER THREAD - No UI access
        outputCallback := func(line string) {
            // This is called from worker thread
            // Safe to call: a.OutputBuffer.Append() is thread-safe
        }
        result := executeBuildProject(buildDir, outputCallback)
        return BuildCompleteMsg{Success: result.Success, Error: result.Error}
    }
}

// 3. Handle completion in Update()
case BuildCompleteMsg:
    a.asyncOperationActive = false
    if msg.Success {
        a.footerHint = "Build completed successfully"
    } else {
        a.footerHint = "Build failed: " + msg.Error
    }
    return a, nil
```

### CRITICAL: Console Output (Reuse TIT Exactly)

```go
// OutputBuffer must be thread-safe (already is in TIT)
// During operation:
// 1. Operation writes lines via callback
// 2. UI re-renders every frame
// 3. Console shows newest at bottom, scrollable

// Reuse these from TIT:
// - ConsoleOutState (scroll tracking)
// - OutputBuffer (thread-safe buffer)
// - RenderConsoleOutput() (rendering logic)
// - OutputLineType + color mapping
```

### CRITICAL: No Menu-Based Exit

```go
// WRONG:
"Exit" menu item → calls tea.Quit()

// CORRECT (TIT pattern):
// Ctrl+C handling at app level:
case tea.KeyCtrlC:
    return a.handleKeyCtrlC()

// In handleKeyCtrlC():
if a.quitConfirmActive {
    return tea.Quit()
}
a.quitConfirmActive = true
a.footerHint = GetFooterMessageText(MessageCtrlCConfirm)
// Wait for second Ctrl+C or timeout
```

---

## Deliverables Checklist

- [ ] `cmd/cake/main.go` - Entry point
- [ ] `internal/app/app.go` - Application state + Update/View
- [ ] `internal/app/handlers.go` - Ctrl+C, menu selection, async commands
- [ ] `internal/app/messages.go` - Msg types
- [ ] `internal/state/project.go` - ProjectState
- [ ] `internal/ops/setup.go` - Setup execution
- [ ] `internal/ops/build.go` - Build execution
- [ ] `internal/ops/clean.go` - Clean execution
- [ ] `internal/utils/project.go` - Project detection
- [ ] `internal/utils/exec.go` - Command execution
- [ ] `internal/utils/generator.go` - Generator logic
- [ ] `internal/utils/platform.go` - Platform detection
- [ ] `internal/ui/*.go` - (COPY FROM TIT, NO CHANGES)
- [ ] `internal/banner/*.go` - (COPY FROM TIT, ADAPT SVG)
- [ ] `go.mod` - Dependency declarations
- [ ] `build.sh` - Build script
- [ ] `README.md` - User documentation

---

## Next Steps for SCAFFOLDER

1. Create literal scaffold (no complexity yet)
2. Copy menu/console/banner files directly from TIT
3. Set up async operation framework (copy pattern exactly)
4. Implement three operation handlers (Setup, Build, Clean)
5. Wire Ctrl+C handling (copy from TIT)

**Key Mandate:** DO NOT add error handling or polish yet. Just scaffolding.

---

**END OF SPECIFICATION**
