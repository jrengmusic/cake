# cake Architecture

**Project:** cake - CMake project management TUI
**Framework:** Bubble Tea (Elm architecture)
**Language:** Go 1.21+
**Module:** github.com/jrengmusic/cake
**Last Updated:** 2026-04-05

---

## Module Structure

```
cake/
├── cmd/cake/
│   └── main.go              # Entry point (tea.NewProgram)
├── internal/
│   ├── constants.go         # AppName, AppVersion (ldflags), BuildsDirName, config names
│   ├── app/                 # Bubble Tea Model (application state & updates)
│   │   ├── app.go           # Application struct, Update(), View(), registerKeyHandlers()
│   │   ├── app_actions.go   # GetVisibleRows(), ToggleRowAtIndex(), executeRowAction()
│   │   ├── app_console.go   # Console mode rendering, renderConsoleMode()
│   │   ├── app_handlers.go  # handleMenuKeyPress(), handleAutoScanTick()
│   │   ├── app_keys.go      # handlePreferencesKeyPress(), handleOperationKeyPress(), handleCtrlC()
│   │   ├── app_render.go    # renderMenuWithBanner(), renderPreferencesWithBanner()
│   │   ├── async_state.go   # AsyncState struct (operationActive, operationAborted)
│   │   ├── constants.go     # CacheRefreshInterval, terminal defaults, scan thresholds
│   │   ├── dispatchers.go   # MessageHandler interface, WindowSizeHandler, KeyDispatcher
│   │   ├── footer.go        # GetFooterContent(), getMenuFooter(), getConsoleFooter()
│   │   ├── init.go          # NewApplication(), loadTheme(), captureVSEnvironment()
│   │   ├── menu.go          # GenerateMenu() — delegates to ui.GenerateMenuRows()
│   │   ├── messages.go      # All Msg types, FooterMessageType, FooterHints, FooterHintShortcuts
│   │   ├── modes.go         # AppMode enum (ModeInvalidProject, ModeMenu, ModePreferences, ModeConsole)
│   │   ├── op_build.go      # startBuildOperation()
│   │   ├── op_clean.go      # startCleanOperation()
│   │   ├── op_clean_all.go  # startCleanAllOperation()
│   │   ├── op_generate.go   # startGenerateOperation()
│   │   ├── op_open.go       # startOpenIDEOperation()
│   │   └── op_regenerate.go # startRegenerateOperation()
│   ├── config/              # Configuration persistence
│   │   └── config.go        # TOML config load/save
│   ├── state/               # Domain state (no UI dependencies)
│   │   ├── project.go       # ProjectState struct, lifecycle methods, query methods
│   │   ├── project_paths.go # GetBuildDirectory(), GetProjectLabel(), GetProjectName()
│   │   ├── project_scan.go  # DetectAvailableProjects(), scanBuildDirectories()
│   │   └── state_test.go
│   ├── ui/                  # Rendering layer (pure functions)
│   │   ├── assets/          # Static assets
│   │   ├── box.go           # Box rendering helper
│   │   ├── buffer.go        # OutputBuffer (sync.RWMutex, singleton, GetSnapshot())
│   │   ├── cake_lie.go      # RenderCakeLieBanner() for invalid project mode
│   │   ├── confirmation.go  # ConfirmationDialog, NewConfirmationDialogWithDefault()
│   │   ├── console.go       # ConsoleOutState, RenderConsoleOutput()
│   │   ├── footer.go        # RenderFooter(), RenderFooterHint(), RenderFooterOverride()
│   │   ├── formatters.go    # Text formatting utilities
│   │   ├── header.go        # RenderHeader(), RenderHeaderInfo(), HeaderState
│   │   ├── layout.go        # RenderReactiveLayout()
│   │   ├── menu.go          # MenuRow struct, GenerateMenuRows() — 8 fixed rows
│   │   ├── menu_render.go   # RenderCakeMenu()
│   │   ├── preferences.go   # Preferences panel rendering
│   │   ├── sizing.go        # DynamicSizing, CalculateDynamicSizing(), NewDynamicSizing()
│   │   ├── theme.go         # Theme struct, LoadTheme(), LoadThemeByName(), GetNextTheme()
│   │   ├── theme_defaults.go # GfxTheme, SpringTheme, SummerTheme, AutumnTheme, WinterTheme (TOML literals)
│   │   └── ui_test.go
│   ├── ops/                 # CMake operations (blocking, run in goroutines)
│   │   ├── build.go         # ExecuteBuildProject() — context.Context, streaming callbacks
│   │   ├── clean.go         # Clean build directory
│   │   ├── open.go          # Open IDE or editor
│   │   └── setup.go         # ExecuteSetupProject() — cmake -G -S -B
│   ├── utils/               # Utility functions
│   │   ├── generators.go    # Generator name constants, GetDirectoryName(), IsGeneratorIDE()
│   │   ├── msvc.go          # FindVCVarsAll(), CaptureVSEnv(), DetectInstalledVSVersions() (Windows)
│   │   ├── msvc_stub.go     # Stub implementations for non-Windows builds
│   │   └── stream.go        # StreamCommand() — reads stdout/stderr with \r handling
│   └── banner/
│       ├── braille.go       # Braille banner rendering
│       ├── svg.go           # SVG banner rendering
│       ├── svg_rasterize.go # SVG rasterization helpers
│       └── banner_test.go
```

---

## Layer Separation Rules

### Layer 1: Entry Point (cmd/cake/main.go)
**Responsibility:** Bootstrap application, start Bubble Tea program
**Dependencies:** None (creates app.Application)
**Rules:**
- Only creates Application and runs tea.Program
- No business logic
- No error handling beyond program.Run() failure

### Layer 2: Application Logic (internal/app/)
**Responsibility:** Bubble Tea Model, message handling, mode switching, operation orchestration
**Dependencies:**
- internal/constants (AppVersion, BuildsDirName)
- internal/state (ProjectState for domain data)
- internal/config (Config for preferences)
- internal/ui (rendering functions)
- internal/ops (CMake operations)
- internal/utils (VS environment capture)
**Forbidden:**
- Direct filesystem access (use state/ops layers)
- Direct CMake execution (use ops/ layer)
**Key Pattern:**
```go
// Update() returns (Model, Cmd)
// Model is new state
// Cmd triggers side effects (tea.Tick, goroutine)
```

### Layer 3: Domain State (internal/state/)
**Responsibility:** Project state, generator detection, build directory scanning, project name extraction
**Dependencies:**
- os, os/exec, runtime (filesystem access)
- internal/utils (generator name constants, VS detection)
- internal/constants (BuildsDirName, config names)
**Forbidden:**
- No UI dependencies (pure data model)
- No CMake execution (detection and scanning only)
**Key Pattern:**
```go
// ProjectState is a plain struct
// ForceRefresh() updates from filesystem
// DetectAvailableProjects() checks system tools
```

### Layer 4: Configuration (internal/config/)
**Responsibility:** Load/save TOML config, apply defaults
**Dependencies:**
- github.com/pelletier/go-toml/v2
**Rules:**
- Auto-creates default config on first run
- Applies defaults for missing values
- Saves on every change (immediate persistence)

### Layer 5: UI Rendering (internal/ui/)
**Responsibility:** Pure rendering functions, no state mutations
**Dependencies:**
- github.com/charmbracelet/lipgloss (styling)
**Forbidden:**
- No state mutation (pure functions)
- No side effects
- No filesystem access (theme loading is the single exception in theme.go)
**Key Pattern:**
```go
// RenderFooter() returns styled string
// Takes width, theme, content as parameters
// No internal state
```

### Layer 6: Operations (internal/ops/)
**Responsibility:** Execute CMake commands, open IDE/editor
**Dependencies:**
- os/exec, context (command execution with cancellation)
- internal/ui (OutputLineType only)
- internal/utils (GetDirectoryName, StreamCommand, FindExecutableInEnv)
**Forbidden:**
- No UI rendering code
- No state mutation
**Key Pattern:**
```go
// ExecuteSetupProject() / ExecuteBuildProject() accept context.Context for cancellation
// Output streamed via appendCallback / replaceCallback
// Returns typed result struct (SetupResult, BuildResult)
```

### Layer 7: Utilities (internal/utils/)
**Responsibility:** Generator constants, directory name mapping, stream command, VS environment
**Dependencies:**
- os, os/exec, runtime
**Rules:**
- Pure functions where possible
- Shared across app/, state/, ops/
- Platform-specific code (msvc.go / msvc_stub.go) guarded by build constraints

---

## Dependency Graph

```
cmd/cake/main.go
    |
    v
internal/app/  (Application, Update, View)
    |
    +---> internal/state/   (ProjectState)
    |         |
    |         +---> internal/utils/  (generator names, VS detection)
    |         +---> internal/constants
    |
    +---> internal/config/  (Config)
    |
    +---> internal/ui/      (rendering)
    |
    +---> internal/ops/     (CMake operations)
              |
              +---> internal/utils/  (StreamCommand, GetDirectoryName)
              +---> internal/ui/     (OutputLineType)
```

**No circular dependencies.**
**Direction: app -> state -> utils, app -> ops -> utils**

---

## Interface Contracts

### ProjectState — Methods Called by App

```go
// Lifecycle
ps.ForceRefresh()                        // Immediate filesystem refresh
ps.Refresh()                             // Rate-limited refresh (skips if within RefreshInterval)
ps.ShouldRefresh() bool                  // Check if refresh interval has elapsed

// Generator / project selection
ps.DetectAvailableProjects()             // Populate AvailableProjects from system tools
ps.CycleToNextProject()                  // Advance SelectedProject forward
ps.CycleToPrevProject()                  // Advance SelectedProject backward
ps.SetSelectedProject(generator string)  // Set directly (restoring from config)
ps.GetProjectLabel() string              // Display-friendly selected project name
ps.GetProjectName() string               // Project name from CMakeLists / Parameters.xml

// Build info
ps.GetBuildPath() string                 // Builds/<dir>/ for selected project
ps.GetBuildDirectory(name string) string // Builds/<dir>/ for arbitrary generator name
ps.GetSelectedBuildInfo() BuildInfo      // BuildInfo for SelectedProject

// Configuration
ps.CycleConfiguration()                  // Toggle Debug <-> Release
ps.SetConfiguration(cfg string)          // Set directly (restoring from config)
ps.Configuration string                  // "Debug" or "Release" (accessed directly)

// Predicates
ps.CanGenerate() bool   // SelectedProject != "" && HasCMakeLists
ps.CanBuild() bool      // build dir exists and IsConfigured
ps.CanOpenIDE() bool    // SelectedProject is an IDE generator
ps.CanOpenEditor() bool // build dir exists
ps.HasBuildsToClean() bool // Builds/ directory is non-empty

// Direct field reads (within app package)
ps.WorkingDirectory string
ps.HasCMakeLists bool
ps.AvailableProjects []Generator
ps.SelectedProject string
ps.Builds map[string]BuildInfo
```

**Contract:**
- ProjectState owns all filesystem scanning
- App calls ForceRefresh() after operations that change build state
- ProjectState is not thread-safe; all access from the Bubble Tea main goroutine

### App Layer ↔ UI Layer

```go
ui.CalculateDynamicSizing(width, height int) DynamicSizing
ui.NewDynamicSizing() DynamicSizing

ui.RenderHeaderInfo(sizing DynamicSizing, theme Theme, state HeaderState) string
ui.RenderHeader(sizing DynamicSizing, theme Theme, info string) string

ui.RenderFooter(shortcuts []FooterShortcut, width int, theme *Theme, rightContent string) string
ui.RenderFooterHint(hint string, width int, theme *Theme) string
ui.RenderFooterOverride(msg string, width int, theme *Theme) string

ui.RenderCakeMenu(rows []MenuRow, selectedIndex int, theme Theme, sizing DynamicSizing) string
ui.GenerateMenuRows(projectLabel, configuration string, canOpenIDE, canClean, hasBuild, hasBuildsToClean bool) []MenuRow

ui.RenderReactiveLayout(sizing DynamicSizing, theme Theme, header, content, footer string) string
ui.RenderConsoleOutput(state *ConsoleOutState, buffer *OutputBuffer, palette Theme, ...) string
ui.RenderCakeLieBanner(contentInnerWidth, contentHeight int) string

ui.LoadThemeByName(name string) (Theme, error)
ui.LoadDefaultTheme() (Theme, error)
ui.GetNextTheme(currentTheme string) (string, error)
ui.CreateDefaultThemeIfMissing() (string, error)
ui.GetBuffer() *OutputBuffer
```

**Contract:**
- UI functions are pure (no side effects beyond theme file I/O at startup)
- App passes all state as parameters; no internal state in rendering functions
- UI returns styled strings; App composes them in View()

### App Layer ↔ Ops Layer

```go
ops.ExecuteSetupProject(
    ctx context.Context,
    workingDir, generator, config string,
    vsEnv []string,
    appendCallback func(string, ui.OutputLineType),
    replaceCallback func(string, ui.OutputLineType),
) SetupResult

ops.ExecuteBuildProject(
    ctx context.Context,
    generator, config, projectRoot string,
    vsEnv []string,
    appendCallback func(string, ui.OutputLineType),
    replaceCallback func(string, ui.OutputLineType),
) BuildResult
```

**Contract:**
- Ops functions are blocking; app wraps them in goroutines via tea.Cmd
- Cancellation via context.Context (cancelContext stored on Application)
- Output streamed live via callbacks (appendCallback / replaceCallback)
- Return typed result structs; app converts to completion Msg types

### App Layer ↔ Config Layer

```go
cfg, _ := config.Load()          // Load config (creates default if missing)
cfg.IsAutoScanEnabled() bool
cfg.AutoScanInterval() int       // minutes
cfg.Theme() string
cfg.LastProject() string
cfg.LastConfiguration() string
cfg.SetAutoScanEnabled(bool) error
cfg.SetAutoScanInterval(int) error
cfg.SetTheme(string) error
```

**Contract:**
- Config layer handles persistence only
- App owns the Config pointer (stored on Application struct)
- Every setter persists immediately (no buffering)

---

## Data Flow

### Startup Flow

```
main.go
  -> app.NewApplication()
      -> ui.CreateDefaultThemeIfMissing()    // ensure 5 themes exist
      -> config.Load()
      -> loadTheme(cfg)
      -> state.NewProjectState()
          -> DetectAvailableProjects()
          -> ForceRefresh()
      -> initialModeAndHint()                // ModeInvalidProject if no CMakeLists.txt
      -> captureVSEnvironment()              // Windows: run vcvarsall, cache env
  -> tea.NewProgram(application)
      -> application.Init()
          -> projectState.ForceRefresh()
          -> GenerateMenu()
          -> NewAsyncState()
          -> NewKeyDispatcher()
          -> registerKeyHandlers()
          -> cmdAutoScanTick() if enabled
  -> program.Run() (event loop)
```

### Menu Navigation Flow

```
User presses down arrow
  -> tea.KeyMsg sent to Update()
      -> windowSize.CanHandle() = false
      -> confirmDialog not active
      -> keyDispatcher.CanHandle() = true
          -> handler[ModeMenu] = handleMenuKeyPress()
              -> moveSelectionDown()
              -> Return (model, nil)
  -> View() called
      -> GetFooterContent() -> getMenuFooter()
          -> visibleRows[selectedIndex].Hint
      -> renderModeContent() -> renderMenuWithBanner()
      -> RenderReactiveLayout(...)
```

### Generate Operation Flow

```
User presses Enter on Regenerate row (or shortcut "g")
  -> handleMenuKeyPress() -> executeRowAction("regenerate")
      -> executeRowActionRegenerate()
          -> if build exists: showConfirmationDialog("regenerate")
          -> else: startGenerateOperation()
              -> asyncState.operationActive = true
              -> mode = ModeConsole
              -> context.WithCancel() -> cancelContext stored
              -> tea.Cmd wraps goroutine:
                  -> ops.ExecuteSetupProject(ctx, ...)
                      -> StreamCommand() streams output via callbacks
                  -> returns GenerateCompleteMsg
  -> GenerateCompleteMsg received in Update()
      -> asyncState.operationActive = false
      -> projectState.ForceRefresh()
      -> GenerateMenu()
      -> if buildAfterGenerate: chain startBuildOperation()
```

### Abort Flow

```
User presses ESC during active operation
  -> handleOperationKeyPress()
      -> asyncState.operationActive = true -> abortActiveOperation()
          -> cancelContext()              // cancels context.Context
          -> asyncState.operationAborted = true
          -> outputBuffer.Append("Operation aborted by user", TypeStderr)
  -> ops function: ctx.Err() == context.Canceled
      -> returns result with Error: "aborted"
  -> CompleteMsg received in Update()
      -> asyncState.operationAborted = true -> reset to false
      -> footerHint = "Operation aborted"
```

### Auto-Scan Flow

```
AutoScanTickMsg received (every N minutes, idle guard: IdleScanThreshold)
  -> handleAutoScanTick()
      -> skip if asyncState.operationActive
      -> skip if time.Since(lastActivityTime) < IdleScanThreshold
      -> projectState.ForceRefresh()
      -> if build state changed:
          -> GenerateMenu()
      -> return cmdAutoScanTick() (schedule next tick)
```

### Console Output Refresh Flow

```
During async operation:
  -> goroutine writes to outputBuffer via appendCallback/replaceCallback
  -> OutputRefreshMsg sent via tea.Tick(CacheRefreshInterval)
      -> forces View() re-render
      -> RenderConsoleOutput() reads buffer via GetSnapshot()
      -> if asyncState.operationActive: schedule next OutputRefreshMsg
      -> else: stop (operation complete, final render already triggered by CompleteMsg)
```

---

## Design Patterns in Use

### Pattern 1: Elm Architecture (Bubble Tea)

**Used for:** Application state management and UI updates

**Structure:**
```go
type Application struct {
    width, height int
    mode          AppMode
    projectState  *state.ProjectState
    config        *config.Config
    // ...
}

func (a *Application) Init() tea.Cmd
func (a *Application) Update(msg tea.Msg) (tea.Model, tea.Cmd)
func (a *Application) View() string
```

**Key Insight:**
- Update() returns new state (Model) and side effects (Cmd)
- View() is a pure function of state (no side effects)
- Model drives View, not the other way around

---

### Pattern 2: Message Dispatcher

**Used for:** Routing messages to mode-specific handlers without a giant switch in Update()

**Implementation location:** `internal/app/dispatchers.go`

**Structure:**
```go
type MessageHandler interface {
    CanHandle(msg tea.Msg) bool
    Handle(a *Application, msg tea.Msg) (tea.Model, tea.Cmd)
}

type WindowSizeHandler struct{}
type KeyDispatcher struct {
    handlers map[AppMode]func(*Application, tea.KeyMsg) (tea.Model, tea.Cmd)
}

// Registration in registerKeyHandlers():
a.keyDispatcher.Register(ModeMenu, app.handleMenuKeyPress)
a.keyDispatcher.Register(ModePreferences, app.handlePreferencesKeyPress)
a.keyDispatcher.Register(ModeConsole, app.handleOperationKeyPress)
a.keyDispatcher.Register(ModeInvalidProject, app.handleInvalidProjectKeyPress)

// In Update():
if a.windowSize.CanHandle(msg) {
    return a.windowSize.Handle(a, msg)
}
if a.keyDispatcher.CanHandle(msg) {
    return a.keyDispatcher.Handle(a, msg)
}
```

**Key Insight:**
- Delegates handling to specialized components
- Avoids giant switch statements in Update()
- New modes require only a Register() call

---

### Pattern 3: Mode-Specific Footer Content

**Used for:** Displaying context-aware hints in footer

**Implementation location:** `internal/app/footer.go`

**Structure:**
```go
func (a *Application) GetFooterContent() string {
    // Priority 1: Quit confirmation (global override)
    if a.quitConfirmActive { ... }

    switch a.mode {
    case ModeInvalidProject: // "The cake is a lie"
    case ModeMenu:           // selected row's Hint
    case ModeConsole:        // scroll shortcuts + scroll status
    case ModePreferences:    // navigation shortcuts
    }
}
```

**Key Insight:**
- Global overrides take priority
- Mode-specific handlers for normal flow
- FooterHintShortcuts map (SSOT) drives all shortcut lists

---

### Pattern 4: Fixed 8-Item Menu with Conditional Selectability

**Used for:** Stable layout with availability-driven interactivity

**Implementation location:** `internal/ui/menu.go`

**Structure:**
```go
// Always returns exactly 8 rows
// Fixed order: Project, Regenerate, OpenIDE, Separator, Configuration, Build, Clean, CleanAll
// Unavailable items: Visible=true, IsSelectable=false (dimmed, not navigable)
func GenerateMenuRows(projectLabel, configuration string, canOpenIDE, canClean, hasBuild, hasBuildsToClean bool) []MenuRow

// Navigation skips non-selectable rows
func (a *Application) GetVisibleRows() []MenuRow {
    // Returns rows where Visible && IsSelectable
}
```

**Key Insight:**
- Fixed row count (always 8) simplifies layout
- Selectability (not visibility) gates navigation
- Separator row: Visible=true, IsSelectable=false (always skipped by navigation)

---

### Pattern 5: Async State Tracker

**Used for:** Managing async operation lifecycle within the Bubble Tea main goroutine

**Implementation location:** `internal/app/async_state.go`

**Structure:**
```go
type AsyncState struct {
    operationActive  bool // unexported — accessed directly within package
    operationAborted bool // unexported — accessed directly within package
}

func NewAsyncState() *AsyncState

// Access pattern (within app package):
a.asyncState.operationActive = true   // start
a.asyncState.operationActive = false  // end (in CompleteMsg handler)
a.asyncState.operationAborted = true  // on ESC abort
if a.asyncState.operationActive { ... }
if a.asyncState.operationAborted { ... }
```

**Key Insight:**
- No public methods — fields accessed directly within package
- Auto-scan checks operationActive to skip during builds
- Ctrl+C blocked while operationActive (blocks quit, shows "operation in progress")

---

### Pattern 6: Context-Based Operation Cancellation

**Used for:** Aborting running cmake commands

**Implementation location:** `internal/app/op_*.go`, `internal/ops/*.go`

**Structure:**
```go
// In app (op_generate.go / op_build.go):
ctx, cancel := context.WithCancel(context.Background())
a.cancelContext = cancel

// In goroutine -> ops.ExecuteSetupProject(ctx, ...)
//   -> exec.CommandContext(ctx, "cmake", ...)
//   -> StreamCommand reads output
//   -> if ctx cancelled: cmd is killed by Go runtime

// On ESC (app_keys.go):
a.cancelContext()              // signal cancellation
a.asyncState.operationAborted = true

// In ops (setup.go / build.go):
if ctx.Err() == context.Canceled {
    return SetupResult{Success: false, Error: "aborted"}
}
```

**Key Insight:**
- context.Context replaces manual process.Kill() calls
- exec.CommandContext kills the subprocess automatically on cancel
- Abort flag on AsyncState distinguishes user-abort from real failure

---

### Pattern 7: Layered Build Path Logic (SSOT)

**Used for:** Consistent build directory path construction

**Implementation location:** `internal/utils/generators.go`, `internal/state/project_paths.go`

**Structure:**
```go
// SSOT constants in utils/generators.go:
const (
    GeneratorXcode  = "Xcode"
    GeneratorNinja  = "Ninja"
    GeneratorVS2026 = "Visual Studio 18 2026"
    GeneratorVS2022 = "Visual Studio 17 2022"
)

func GetDirectoryName(generator string) string {
    // VS2026 -> "VS2026", VS2022 -> "VS2022", others unchanged
}

func GetGeneratorNameFromDirectory(dirName string) string {
    // Reverse: "VS2026" -> GeneratorVS2026, etc.
}

// Build path in state/project_paths.go:
func (ps *ProjectState) GetBuildDirectory(generatorName string) string {
    return filepath.Join(ps.WorkingDirectory, internal.BuildsDirName, utils.GetDirectoryName(generatorName))
}
```

**Key Insight:**
- CMake generator name constants live in utils (shared by state and ops)
- Short directory names (VS2026, VS2022) from GetDirectoryName
- Reverse mapping (GetGeneratorNameFromDirectory) used when scanning existing build dirs

---

### Pattern 8: OutputBuffer Thread Safety

**Used for:** Sharing console output between the goroutine writing output and the Bubble Tea render loop

**Implementation location:** `internal/ui/buffer.go`

**Structure:**
```go
type OutputBuffer struct {
    mu       sync.RWMutex
    maxLines int
    lines    []OutputLine
}

// Writers (goroutine context):
b.Append(text string, lineType OutputLineType)
b.ReplaceLast(text string, lineType OutputLineType) // for \r progress lines

// Readers (Bubble Tea main goroutine):
b.GetSnapshot() ([]OutputLine, int)  // atomic read of lines + count
b.GetAllLines() []OutputLine
b.GetLines(startIdx, count int) []OutputLine
b.GetLineCount() int

// Global singleton:
var globalBuffer = &OutputBuffer{maxLines: 1000, ...}
func GetBuffer() *OutputBuffer
```

**Key Insight:**
- sync.RWMutex: multiple concurrent readers, exclusive writers
- GetSnapshot() reads lines and count under a single lock (no TOCTOU)
- Global singleton; app holds pointer via ui.GetBuffer() at init

---

## Threading Model

**Bubble Tea runs a single event loop on the main goroutine.**

- All Update() and View() calls: main goroutine (no synchronization needed for Application fields)
- Async operations (cmake): spawned as goroutines via tea.Cmd
  - Write to outputBuffer (Append/ReplaceLast) — mutex protected
  - Return completion Msg via channel — delivered to main goroutine by Bubble Tea
- OutputBuffer: sync.RWMutex protects concurrent writes (goroutine) and reads (main goroutine via RenderConsoleOutput)
- cancelContext: written once at operation start, read from ESC handler (both on main goroutine via Update)

**No Application struct fields require mutex protection — all reads/writes on main goroutine.**

---

## Error Handling Strategy

### Fail Fast at Operation Boundaries

Ops functions return typed result structs with Success bool and Error string. The app converts these to footer hints or console output. Errors are never silently dropped.

### Non-Fatal Initialization Failures

config.Load(), LoadTheme(), and captureVSEnvironment() are non-fatal at startup — app proceeds with defaults if they fail (explicit comment in init.go per function).

### UI Shows Errors, Does Not Handle Them

App layer displays error text in footerHint or console output. Ops layer does not know about UI.

---

## Anti-Patterns Avoided

### No Global State

All state passed as parameters to rendering functions. No package-level variables in ui/ except the OutputBuffer singleton (justified: shared between goroutines).

### No Silent Failures

Errors propagated explicitly as return values. Non-fatal cases documented at call sites.

### No Manual Process Kill

context.Context passed to exec.CommandContext — Go runtime kills subprocess on cancel. No process.Kill() calls.

### No Elm-Breaking Imperative Rendering

View() is pure. No side effects, no direct output writes in View().

---

## Glossary

**AppMode:** Application mode enum — ModeInvalidProject, ModeMenu, ModePreferences, ModeConsole

**AsyncState:** Tracks active operation and abort flag (unexported fields, package-local access)

**BuildInfo:** Build directory state — Generator, Path, Exists, IsConfigured, Configs

**DynamicSizing:** Terminal dimension calculations — ContentHeight, ContentInnerWidth, etc.

**Generator:** CMake generator with metadata — Name string, IsIDE bool

**MenuRow:** Single menu row — ID, Shortcut, ShortcutLabel, Emoji, Label, Value, Visible, IsAction, IsSelectable, Hint

**OutputBuffer:** Thread-safe circular buffer (sync.RWMutex, 1000 lines, global singleton)

**ProjectState:** Domain state — WorkingDirectory, AvailableProjects, SelectedProject, Builds, Configuration

**Theme:** Semantic color set loaded from TOML file (~/.config/cake/themes/<name>.toml)

**tea.Cmd:** Bubble Tea command — async operation returning a Msg

**tea.Msg:** Bubble Tea message — key press, window resize, completion event, tick

**TIT:** "Terminal Interface for git" — reference TUI project (architectural pattern source)

---

**End of ARCHITECTURE.md**

Rock 'n Roll!
JRENG!
