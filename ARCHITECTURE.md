# cake Architecture

**Project:** cake - CMake project management TUI
**Framework:** Bubble Tea (Elm architecture)
**Language:** Go 1.21+
**Last Updated:** 2026-01-29

---

## Module Structure

```
cake/
├── cmd/cake/
│   └── main.go              # Entry point (tea.NewProgram)
├── internal/
│   ├── app/                 # Bubble Tea Model (application state & updates)
│   │   ├── app.go           # Application model, Update(), View()
│   │   ├── async_state.go   # Async operation state tracker
│   │   ├── dispatchers.go   # MessageHandler interfaces (WindowSize, Key)
│   │   ├── footer.go        # Footer content manager (mode-specific)
│   │   ├── init.go          # Application initialization
│   │   ├── keyboard.go      # Keyboard shortcuts mapping
│   │   ├── menu.go          # Menu generation from state
│   │   ├── messages.go      # Application messages
│   │   ├── modes.go         # AppMode enum (Menu, Preferences, Console)
│   │   └── op_*.go          # Operation handlers (generate, build, clean, open)
│   ├── config/              # Configuration persistence
│   │   └── config.go        # TOML config load/save
│   ├── state/               # Domain state (no UI dependencies)
│   │   └── project.go       # ProjectState, Generator, BuildInfo
│   ├── ui/                  # Rendering layer (pure functions)
│   │   ├── footer.go        # RenderFooter(), FooterShortcut, FooterStyles
│   │   ├── header.go        # RenderHeader(), HeaderState
│   │   ├── layout.go        # RenderReactiveLayout(), composition
│   │   ├── menu.go          # MenuRow, GenerateMenuRows()
│   │   ├── sizing.go        # DynamicSizing calculations
│   │   ├── theme.go         # Theme management (5 themes)
│   │   ├── confirmation.go  # Confirmation dialog (TIT pattern)
│   │   ├── console.go        # Console output rendering
│   │   ├── spinner.go        # Spinner animation
│   │   ├── formatters.go    # Text formatting utilities
│   │   ├── buffer.go        # Output buffer management
│   │   ├── box.go           # Box rendering helper
│   │   └── statusbar.go     # Status bar component
│   ├── ops/                 # CMake operations (pure functions)
│   │   ├── setup.go         # CMake setup commands
│   │   ├── build.go         # CMake build commands
│   │   ├── clean.go         # Clean build directories
│   │   └── open.go          # Open IDE/editor
 │   └── utils/               # Utility functions
 │       ├── generators.go    # Generator constants, GetDirectoryName()
 │       ├── stream.go        # StreamCommand() helper
 │       ├── generator.go      # Generator detection
 │       ├── platform.go       # Platform detection
 │       ├── project.go        # Project utilities
 │       └── exec.go          # Command execution
└── internal/banner/
    ├── svg.go               # SVG banner rendering
    └── braille.go           # Braille banner rendering
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
**Responsibility:** Bubble Tea Model, message handling, mode switching
**Dependencies:**
- ✓ internal/state (ProjectState for domain data)
- ✓ internal/config (Config for preferences)
- ✓ internal/ui (rendering functions)
- ✓ internal/ops (CMake operations)
**Forbidden:**
- ❌ Direct filesystem access (use state/ops layers)
- ❌ Direct CMake execution (use ops/ layer)
**Key Pattern:**
```go
// Update() returns (Model, Cmd)
// Model is new state (immutable update)
// Cmd triggers side effects (tea.Tick, tea.Exec)
```

### Layer 3: Domain State (internal/state/)
**Responsibility:** Project state, generator detection, build scanning
**Dependencies:**
- ✓ os, os/exec, runtime (filesystem access)
**Forbidden:**
- ❌ No UI dependencies (pure data model)
- ❌ No CMake execution (just detection/scanning)
**Key Pattern:**
```go
// ProjectState is pure Go struct
// ForceRefresh() updates from filesystem
// DetectAvailableGenerators() checks system tools
```

### Layer 4: Configuration (internal/config/)
**Responsibility:** Load/save TOML config, apply defaults
**Dependencies:**
- ✓ github.com/pelletier/go-toml/v2
**Rules:**
- Auto-creates default config on first run
- Applies defaults for missing values
- Handles config file corruption gracefully

### Layer 5: UI Rendering (internal/ui/)
**Responsibility:** Pure rendering functions, no state
**Dependencies:**
- ✓ github.com/charmbracelet/lipgloss (styling)
**Forbidden:**
- ❌ No state mutation (pure functions)
- ❌ No side effects
- ❌ No filesystem access
**Key Pattern:**
```go
// RenderFooter() returns styled string
// Takes width, theme, content as parameters
// No internal state
```

### Layer 6: Operations (internal/ops/)
**Responsibility:** Execute CMake commands, open IDE/editor
**Dependencies:**
- ✓ os/exec (command execution)
**Forbidden:**
- ❌ No UI code
- ❌ No state mutation (pure functions)
**Key Pattern:**
```go
// ExecuteGenerate() returns tea.Cmd
// tea.Exec runs command, returns ProcessDoneMsg
// App handles ProcessDoneMsg in Update()
```

### Layer 7: Utilities (internal/utils/)
**Responsibility:** Platform detection, generator detection, execution helpers
**Dependencies:**
- ✓ os, os/exec, runtime
**Rules:**
- Pure functions where possible
- Shared across app/, state/, ops/

---

## Dependency Graph

```
cmd/cake/main.go
    ↓
internal/app/app.go (Model)
    ↓
    ├── internal/state/project.go (ProjectState)
    ├── internal/config/config.go (Config)
    ├── internal/ui/* (rendering)
    └── internal/ops/* (CMake operations)
         ↓
         └── internal/utils/* (shared utilities)
```

**No circular dependencies.**
**Direction: app → state → utils, app → ops → utils**

---

## Interface Contracts

### ProjectState ↔ App Layer

**ProjectState Methods Called by App:**
```go
ps.Refresh()              // Refresh state (with rate limiting)
ps.ForceRefresh()         // Force immediate refresh
ps.GetGeneratorLabel()     // Selected generator name
ps.Configuration          // Debug/Release string
ps.GetSelectedBuildInfo() // BuildInfo for selected generator
ps.CanOpenIDE()           // bool (IDE available)
ps.CycleGenerator()       // Switch to next generator
ps.CycleConfiguration()   // Toggle Debug/Release
```

**Contract:**
- ProjectState owns filesystem scanning
- App calls Refresh(), not direct access
- ProjectState is immutable from App's perspective (methods return new state or modify internal fields)

### App Layer ↔ UI Layer

**UI Functions Called by App:**
```go
ui.RenderHeader(width, theme, HeaderState{...})         // returns styled string
ui.RenderFooter(shortcuts, width, theme, rightContent)  // returns styled string
ui.RenderCakeMenu(menuItems, selectedIndex, theme)       // returns styled string
ui.RenderReactiveLayout(header, content, footer, sizing) // returns styled string
ui.CalculateDynamicSizing(width, height)                // returns DynamicSizing
```

**Contract:**
- UI functions are pure (no side effects)
- App passes all state as parameters (no internal state in UI layer)
- UI returns styled strings, App composes them

### App Layer ↔ Ops Layer

**Ops Functions Called by App:**
```go
ops.ExecuteGenerate(projectPath, buildDir, generator, config)      // returns tea.Cmd
ops.ExecuteBuild(buildDir, configuration)                            // returns tea.Cmd
ops.ExecuteClean(buildDir)                                            // returns tea.Cmd
ops.ExecuteOpenIDE(buildPath, generator)                              // returns tea.Cmd
ops.ExecuteOpenEditor(buildPath)                                      // returns tea.Cmd
```

**Contract:**
- Ops functions return tea.Cmd (async command)
- App handles ProcessDoneMsg in Update()
- Ops layer has no UI dependencies

### App Layer ↔ Config Layer

**Config Functions Called by App:**
```go
cfg := config.Load()        // Load config (creates default if missing)
config.Save(cfg)            // Save config
cfg.IsAutoScanEnabled()     // bool
cfg.AutoScanInterval()      // int (minutes)
cfg.Theme                  // string
```

**Contract:**
- Config layer handles persistence only
- App owns Config instance (in Application struct)
- Config saves on every change (immediate persistence)

---

## Data Flow

### Startup Flow

```
main.go
  → app.NewApplication()
      → config.Load() (creates default if missing)
      → state.NewProjectState()
          → DetectAvailableGenerators()
          → ForceRefresh()
      → Build menu items
  → tea.NewProgram(application)
      → application.Init()
          → ForceRefresh project state
          → Generate menu
          → Start auto-scan ticker (if enabled)
  → program.Run() (event loop)
```

### Menu Navigation Flow

```
User presses ↓ key
  → tea.KeyMsg sent to Update()
      → KeyDispatcher.Handle()
          → handler[a.mode] (ModeMenu)
              → handleMenuKeyPress()
                  → Increment selectedIndex (skip separators)
                  → Return (model, nil)
  → View() called with new model
      → GetFooterContent() (for selected row's hint)
          → getMenuFooter()
              → RenderFooterHint(selectedRow.Hint, ...)
      → RenderReactiveLayout()
          → RenderHeader(...)
          → RenderCakeMenu(menuItems, selectedIndex, ...)
          → GetFooterContent()
```

### Generate Operation Flow

```
User selects Generate, presses Enter
  → handleMenuKeyPress()
      → executeRowAction("generate")
          → cmdGenerateProject()
              → ops.ExecuteGenerate(...)
                  → tea.Exec("cmake -S . -B Builds/...")
  → ProcessDoneMsg received
      → asyncState.End()
      → ExitAllowed = true
      → ProjectState.ForceRefresh() (detect new build)
      → Return to menu mode
```

### Auto-Scan Flow

```
AutoScanTickMsg received (every N minutes)
  → handleAutoScanTick()
      → Skip if asyncState.IsActive()
      → ProjectState.ForceRefresh()
          → Scan build directories
          → Update Builds map
      → If build state changed:
          → menuItems = GenerateMenu() (update visibility)
          → Show "Builds detected" in footer
      → Return cmdAutoScanTick() (schedule next tick)
```

---

## Design Patterns in Use

### Pattern 1: Elm Architecture (Bubble Tea)

**Used for:** Application state management and UI updates

**Structure:**
```go
type Application struct {
    // State fields
    width, height int
    mode          AppMode
    projectState  *state.ProjectState
    config        *config.Config
    // ... other state
}

func (a *Application) Init() tea.Cmd
func (a *Application) Update(msg tea.Msg) (tea.Model, tea.Cmd)
func (a *Application) View() string
```

**Key Insight:**
- Update() returns new state (Model) and side effects (Cmd)
- View() is pure function of state (no side effects)
- Model drives View, not the other way around

---

### Pattern 2: Message Dispatcher

**Used for:** Routing messages to mode-specific handlers

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

// In Update():
for _, handler := range []MessageHandler{a.windowSize, a.keyDispatcher} {
    if handler.CanHandle(msg) {
        return handler.Handle(a, msg)
    }
}
```

**Key Insight:**
- Delegates handling to specialized components
- Avoids giant switch statements in Update()
- Easy to add new handlers (open/closed principle)

---

### Pattern 3: Mode-Specific Footer Content

**Used for:** Displaying context-aware hints in footer

**Implementation location:** `internal/app/footer.go`

**Structure:**
```go
func (a *Application) GetFooterContent() string {
    // Priority 1: Quit confirmation (global override)
    if a.quitConfirmActive {
        return ui.RenderFooterOverride(GetFooterMessageText(MessageCtrlCConfirm), ...)
    }

    // Priority 2: Mode-specific footer
    switch a.mode {
    case ModeMenu:
        return a.getMenuFooter(width) // selected row's hint
    case ModeConsole:
        return a.getConsoleFooter(width) // scroll shortcuts + status
    case ModePreferences:
        return ui.RenderFooter(shortcuts, ...) // navigation shortcuts
    }
}
```

**Key Insight:**
- Global overrides take priority (Ctrl+C confirm)
- Mode-specific handlers for normal flow
- Separates content from rendering

---

### Pattern 4: Fixed 7-Item Menu with Conditional Visibility

**Used for:** Preference-style menu matching TIT

**Implementation location:** `internal/ui/menu.go`

**Structure:**
```go
// Always returns 7 rows, visibility determined by Visible field
func GenerateMenuRows(generator, config string, canOpenIDE, canClean bool) []MenuRow {
    return []MenuRow{
        {ID: "generator", Visible: true, ...},
        {ID: "regenerate", Visible: true, ...},
        {ID: "openIde", Visible: canOpenIDE, ...},
        {ID: "separator", Visible: true, ...},
        {ID: "configuration", Visible: true, ...},
        {ID: "build", Visible: buildExists, ...},
        {ID: "clean", Visible: buildExists, ...},
    }
}

// Navigation skips non-visible rows
func (a *Application) GetVisibleRows() []MenuRow {
    var visible []MenuRow
    for _, row := range a.menuItems {
        if row.Visible {
            visible = append(visible, row)
        }
    }
    return visible
}
```

**Key Insight:**
- Fixed array size (always 7 items) simplifies layout
- Visibility flag controls display (no dynamic array resizing)
- Navigation uses visible indices (0-5) instead of array indices (0-6)

---

### Pattern 5: Async State Tracker

**Used for:** Managing async operation lifecycle

**Implementation location:** `internal/app/async_state.go`

**Structure:**
```go
type AsyncState struct {
    OperationActive  bool
    OperationAborted bool
    ExitAllowed      bool
}

// Lifecycle methods
func (as *AsyncState) Start()  // Active=true, Aborted=false, ExitAllowed=false
func (as *AsyncState) End()     // Active=false
func (as *AsyncState) Abort()   // Aborted=true
func (as *AsyncState) SetExitAllowed(bool)

// Queries
func (as *AsyncState) IsActive() bool
func (as *AsyncState) IsAborted() bool
func (as *AsyncState) CanExit() bool
```

**Key Insight:**
- Encapsulates async operation state in dedicated struct
- Prevents invalid state transitions (e.g., can't abort if not active)
- Used by auto-scan to skip during operations

---

### Pattern 6: Layered Build Path Logic

**Used for:** Generating correct build paths using SSOT constants and directory mapping

**Implementation location:** `internal/utils/generators.go`, `internal/state/project.go`

**Structure:**
```go
// In utils/generators.go:
const (
    GeneratorXcode      = "Xcode"
    GeneratorNinja      = "Ninja"
    GeneratorVS2026     = "Visual Studio 18 2026"
    GeneratorVS2022     = "Visual Studio 17 2022"
)

func GetDirectoryName(generator string) string {
    // Maps generator names to directory names
    switch generator {
    case GeneratorXcode:
        return "Xcode"
    case GeneratorNinja:
        return "Ninja"
    case GeneratorVS2026:
        return "VS2026"
    case GeneratorVS2022:
        return "VS2022"
    default:
        return "Build"
    }
}

// In state/project.go:
func (ps *ProjectState) GetBuildPath() string {
    buildDir := filepath.Join(ps.WorkingDirectory, "Builds")
    return filepath.Join(buildDir, GetDirectoryName(ps.SelectedGenerator))
}
```

**Key Insight:**
- Single source of truth for build path logic
- Eliminates code duplication (removed 5+ duplicate build path constructions in Sprint 6)
- All generators use multi-config structure (Builds/<Generator>/)
- CMake names in constants, shortened directory names via GetDirectoryName()

---

## Key Design Decisions

### Decision 1: Elm Architecture (Bubble Tea)

**Why:**
- Predictable state management (pure functions)
- Easy to test (no hidden state)
- Real-time UI updates (tea.Cmd for async operations)

**Trade-offs:**
- Learning curve for Elm architecture concepts
- More verbose than direct DOM manipulation (but safer)

---

### Decision 2: Fixed 7-Item Menu Structure

**Why:**
- Matches TIT's preference menu pattern
- Simplifies layout calculations (no dynamic resizing)
- Easy to add/remove items later (just change visibility)

**Trade-offs:**
- Fixed positions (can't reorder items easily)
- Requires placeholder rows for hidden items

---

### Decision 3: Generator Detection via System Tools

**Why:**
- Detects what's actually available (not what might exist)
- Fast (checks xcodebuild, ninja, vswhere.exe once)
- No disk scanning for generator existence

**Trade-offs:**
- Doesn't detect build directories until generated
- Requires manual refresh after external changes (auto-scan mitigates)

---

### Decision 4: Async Operations with tea.Exec

**Why:**
- Non-blocking UI (terminal doesn't freeze during CMake)
- Real-time output streaming (ProcessDoneMsg streams stdout)
- Easy cancellation (tea.Exec supports kill)

**Trade-offs:**
- Requires AsyncState tracker to manage lifecycle
- More complex than blocking operations

---

### Decision 5: Immediate Config Persistence

**Why:**
- User never loses changes
- No "unsaved changes" state to manage
- Simple mental model (change = saved)

**Trade-offs:**
- More disk I/O (negligible for small config)
- No undo (but settings are simple toggles/cycles)

---

## Threading Model

**Single-threaded Bubble Tea event loop**

All operations run in the main goroutine:
- UI updates: Immediate (no thread safety issues)
- Async operations: tea.Exec spawns subprocess, but callbacks run in main loop
- Auto-scan: tea.Tick sends messages, handled in main loop

**No mutexes or locks needed.**

---

## Error Handling Strategy

### Fail Fast at Boundaries

**Pattern:** Check errors immediately, don't suppress

**Example:**
```go
// In ops/setup.go
cmd := exec.Command("cmake", ...)
if err := cmd.Run(); err != nil {
    return fmt.Errorf("cmake setup failed: %w", err)
}
```

**Why:**
- Silent failures waste hours debugging
- User needs to know what went wrong immediately

---

### UI Layer Shows Errors, Doesn't Handle Them

**Pattern:** App layer shows errors in footer/console, doesn't fix

**Example:**
```go
// In handleMenuKeyPress()
if err := ops.ExecuteGenerate(...); err != nil {
    footerHint = "Generate failed: " + err.Error()
    return
}
```

**Why:**
- Separates error detection (ops/) from error display (ui/)
- User can retry with different settings

---

## Anti-Patterns Avoided

### ❌ Anti-Pattern: Global State

**Problem:** Implicit dependencies, hard to test

**Solution:** Pass all state as parameters
```go
// GOOD (pure function)
ui.RenderFooter(shortcuts, width, theme, rightContent)

// BAD (hidden state)
ui.RenderFooter() // reads globalWidth, globalTheme
```

---

### ❌ Anti-Pattern: Silent Failure

**Problem:** Wastes hours debugging hidden bugs

**Solution:** Explicit error handling
```go
// GOOD (fail fast)
if err != nil {
    return fmt.Errorf("operation failed: %w", err)
}

// BAD (silent failure)
if err != nil {
    return  // What happened?
}
```

---

### ❌ Anti-Pattern: Direct DOM Manipulation

**Problem:** Breaks Elm architecture, unpredictable state

**Solution:** Return new model, let View() render
```go
// GOOD (Elm pattern)
func (a *Application) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    a.selectedIndex++  // Modify state
    return a, nil      // Return new model
}

// BAD (imperative UI)
func (a *Application) Update(msg tea.Msg) {
    renderMenu()        // Render directly
    updateDOM()         // Mutable state
}
```

---

### ❌ Anti-Pattern: Code Duplication

**Problem:** Maintenance nightmare, inconsistencies

**Solution:** Single Source of Truth (SSOT)
```go
// GOOD (SSOT)
func (ps *ProjectState) GetBuildPath() string { ... }

// BAD (duplicate logic)
path1 := filepath.Join("Builds", generator)
path2 := filepath.Join("Builds", generator, config) // Which is correct?
```

---

## Glossary

**TIT**: "The Interactive Terminal" - Reference TUI framework (architectural pattern source)

**Bubble Tea**: Go framework for building terminal user interfaces using Elm architecture

**Elm Architecture**: UI pattern: State → View → Messages → Update State

**tea.Cmd**: Bubble Tea command (async operation, returns messages)

**tea.Msg**: Bubble Tea message (key press, window resize, process done, etc.)

**MenuRow**: Single menu item with ID, Shortcut, Emoji, Label, Value, Visible, IsAction, Hint

**AsyncState**: Tracker for async operation lifecycle (Active, Aborted, ExitAllowed)

**ProjectState**: Domain state (WorkingDirectory, AvailableGenerators, Builds, Configuration)

**Generator**: CMake generator with metadata (Name, IsIDE, IsMultiConfig)

**BuildInfo**: Build directory state (Generator, Path, Exists, IsConfigured, Configs)

**AppMode**: Application mode (ModeMenu, ModePreferences, ModeConsole)

---

**End of ARCHITECTURE.md**

Rock 'n Roll!
JRENG!
