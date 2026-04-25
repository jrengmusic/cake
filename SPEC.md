# CAKE Specification v1.0.0

## Overview

**Purpose:** CMake project management tool with visual preference-style TUI for quick project configuration and building
**Target User:** Developers who want visual CMake control without typing commands
**Core Workflow:** Select generator → Configure → Build/Clean/Open, with persistent settings

## Technology Stack

- **Language:** C++17
- **Framework:** JUCE8 (headless console app), jam_tui (terminal UI primitives)
- **Platform:** macOS first, Windows post-MVP
- **Dependencies:**
  - JUCE modules: juce_core, juce_events, juce_data_structures, juce_graphics, juce_gui_basics
  - jam modules: jam_core, jam_data_structures, jam_markdown, jam_tui, jam_subprocess
- **External Requirements:** CMake installed in PATH
- **Binary:** `cakec`

## Architecture

### BLESSED MVC — Event-Driven, Unidirectional

**Model:** `State` — `juce::ValueTree` is the only SSOT while app lives. Domain truth only. Atoms + timer flush for cross-thread writes.

**View:** jam_tui components — stateless. Read/write State VT. Hold only transient render state (spinnerFrame, scrollOffset, menuIndex). Many thin view files.

**Controller:** `MainComponent` — sole orchestrator. Listens to State VT via `juce::ValueTree::Listener`. Dispatches callbacks event-driven. No manual booleans, no manual lambdas.

**Data flow:** Input → View → State VT mutation → Listener fires → Controller dispatches side effects (cmake ops, auto-scan, config persistence)

### Threading Model

| Thread | Owns | Crossing |
|---|---|---|
| Message (JUCE main) | State VT, all Views, MainComponent | — |
| Subprocess worker (jam_subprocess) | juce::ChildProcess (cmake) | atomic writes → State flush, `callAsync` for console lines |
| Timer (juce::Timer) | State flush, auto-scan tick | message thread (inherited) |
| File watcher (jam::File::Watcher) | theme directory monitoring | `callAsync` → message thread |

Zero locks on hot path. `callAsync` only crossing primitive.

### ValueTree Schema (Domain Truth Only)

```
CAKE                                    # root
├── PROJECT                             # domain state
│   ├── workingDirectory    (string)
│   ├── hasCMakeLists       (bool)
│   ├── selectedGenerator   (string)
│   ├── configuration       (string)
│   └── GENERATOR[]                     # detected generators
│       ├── name            (string)
│       └── isIDE           (bool)
├── BUILDS                              # scan cache (refreshed by auto-scan timer)
│   └── BUILD[]
│       ├── generator       (string)
│       ├── path            (string)
│       ├── exists          (bool)
│       └── isConfigured    (bool)
├── ASYNC                               # operation state
│   ├── isActive            (bool)
│   ├── isAborted           (bool)
│   └── currentOp           (string)
├── CONFIG                              # persisted preferences (grafted from XML on startup)
│   ├── autoScanEnabled     (bool)
│   ├── autoScanInterval    (int, minutes)
│   └── theme               (string)
└── THEME                               # loaded from ~/.config/cake/themes/*.xml
    └── ...                             # color properties per theme
```

**Not in VT (view-owned transient state):** menuIndex, spinnerFrame, scrollOffset, preferencesVisible.

**Mode** is computed from VT state, not stored:
- `ASYNC.isActive` → console
- `not PROJECT.hasCMakeLists` → invalidProject
- `preferencesVisible` (view transient) → preferences
- else → menu

### Config Persistence

`~/.config/cake/config.xml` is bootstrap and persistence layer only. On startup: load XML → graft into State VT CONFIG subtree. Write defaults if missing. While app lives, State VT is the only SSOT. A `ValueTree::Listener` on CONFIG subtree persists changes back to XML automatically.

## Core Principles

1. **Preference-Style Interface**: Single-page menu with toggleable values, no nested submenus
2. **Conditional Selectability**: Menu items are always visible but become unselectable (dimmed) when unavailable based on build state and generator capabilities
3. **System Tool Detection**: Available generators determined by installed tools, not disk scanning
4. **65/35 Split Layout**: Menu on left, braille banner on right, both centered
5. **Build Path Convention**: Strict `Builds/<Generator>/` structure (all generators multi-config)
6. **Persistent Configuration**: Settings saved to `~/.config/cake/config.xml`

## State Model

### Project State Determination

```
State determined by: (workingDirectory, availableGenerators, buildDirectories)
- workingDirectory: Current directory containing CMakeLists.txt
- availableGenerators: System tools detected (Xcode, Ninja, Visual Studio)
- buildDirectories: Scanned Builds/ subdirectories matching generator pattern
```

### Generator Types

```
All generators multi-config (build contains all configurations):
- Xcode (macOS only, IDE)
- Ninja (cross-platform, CLI)
- Visual Studio (Windows only, IDE)
Path: Builds/<Generator>/
```

### Menu Item Selectability Rules

```
All rows are always visible. Unavailable rows are shown dimmed and are not navigable.

Project: ALWAYS selectable
Generate/Regenerate: ALWAYS selectable (requires CMakeLists.txt to execute)
Open IDE/Editor: selectable if any generator is selected AND build exists
Configuration: ALWAYS selectable
Build: selectable if build exists AND configured
Clean: selectable if build exists for selected project
Clean All: selectable if any build directory exists
```

Selectability is computed by MenuBuilder as a pure function of PROJECT + BUILDS VT subtrees. No stored menu state.

## Feature Specifications

### Feature: Main Menu Navigation

#### User Flow (Happy Path)
1. User launches `cakec` in directory with CMakeLists.txt
2. System shows preference menu with available options
3. User presses up/down or j/k to navigate rows
4. System highlights selected row with inverted colors
5. User presses Enter or Space on a row
6. System executes toggle (for settings) or action (for operations)

**UI Display:**
```
Project directory:
/Users/username/myproject

  Project                Xcode
  Generate
  Open IDE / Open Editor  (label is dynamic, depends on generator)
────────────────────────────────────────
  Configuration          Debug
  Build
  Clean
  Clean All

up/down navigate | Enter select | Ctrl+C quit | / preferences
```

**User Input:**
- up/k: Move selection up (skip unselectable rows)
- down/j: Move selection down (skip unselectable rows)
- Enter/Space: Toggle value or execute action
- /: Open preferences screen
- Ctrl+C: Quit (confirm if pressed twice)

**System Response:**
- Selection moves with smooth highlight
- Separators cannot be selected (auto-skip)
- Actions trigger immediately on Enter

#### Edge Cases

##### Edge Case 1: No CMakeLists.txt
**Scenario:** User runs cakec in directory without CMakeLists.txt
**Expected Behavior:** Show CakeLie banner (full-screen "the cake is a lie" braille art)

##### Edge Case 2: Empty Builds Directory
**Scenario:** No builds exist yet
**Expected Behavior:** Show only Project selection and Generate action as selectable
**Menu State:** Build, Clean, Open IDE/Editor shown but unselectable (dimmed)

##### Edge Case 3: Navigation at Boundaries
**Scenario:** User presses up at first item or down at last
**Expected Behavior:** Selection stays at boundary (no wrap)

### Feature: Generator Selection

#### User Flow (Happy Path)
1. User selects Project row
2. User presses Enter or Space
3. View writes next generator to PROJECT.selectedGenerator in State VT
4. Display updates via VT listener

**Available Generators Detection:**
- macOS: Xcode (if xcodebuild exists), Ninja (if ninja exists)
- Linux: Ninja (if ninja exists)
- Windows: Visual Studio (if vswhere.exe exists), Ninja (if ninja exists)

**Cycle Order:**
```
Xcode → Ninja → Xcode (loop)
```

#### Error Handling

| Error Condition | User Sees | System Action |
|---|---|---|
| No generators available | "No CMake generators found" | Disable generate action |
| Only one generator | Generator name shown (no cycling) | Enter does nothing |

### Feature: Configuration Toggle

#### User Flow
1. User selects Configuration row
2. User presses Enter or Space
3. View writes toggled value to PROJECT.configuration in State VT
4. Display updates via VT listener

**Values:** Debug / Release (bidirectional toggle)

### Feature: Generate/Regenerate Operation

#### User Flow (Generate - First Time)
1. User selects "Generate" row and presses Enter
2. View writes to State VT → Controller listener fires → dispatches CmakeRunner.generate()
3. CmakeRunner sets ASYNC atoms (isActive, currentOp) → flushed to VT
4. Mode computes to console → MainComponent shows jam::tui::Console
5. CmakeRunner streams stdout/stderr via `appendLine()` through `callAsync`
6. On completion, ASYNC.isActive → false → mode computes back to menu

**CMake Command:**
```bash
cmake -S . -B Builds/Xcode -G Xcode
```

#### User Flow (Regenerate - Build Exists)
Label changes to "Regenerate" when build exists. Same flow, overwrites existing build.

**UI During Operation:**
```
[Console with braille spinner in header]
GENERATING ⠋
Setting up CMake...
-- The C compiler identification is AppleClang 14.0.0
-- Detecting C compiler ABI info
[... cmake output streams ...]

ESC to abort
```

#### Edge Cases

##### Edge Case 1: CMake Not Installed
**Scenario:** cmake command not found in PATH
**Expected Behavior:** Show error immediately

##### Edge Case 2: Invalid CMakeLists.txt
**Scenario:** CMakeLists.txt has syntax errors
**Expected Behavior:** Show CMake error output in console

##### Edge Case 3: Disk Full
**Scenario:** No space to create build directory
**Expected Behavior:** Show system error in console

### Feature: Build Operation

#### User Flow
1. User selects "Build" row and presses Enter
2. Controller dispatches CmakeRunner.build()
3. Same console mode flow as Generate

**Build Command:**
```bash
cmake --build Builds/Xcode --config Debug
```

#### Edge Cases

##### Edge Case 1: Build Directory Deleted
**Scenario:** User deletes build directory externally
**Expected Behavior:** Auto-scan detects missing directory, updates BUILDS VT, menu selectability recomputes

##### Edge Case 2: Compilation Errors
**Scenario:** Source code has errors
**Expected Behavior:** Show full compiler output in console

### Feature: Clean Operation

#### User Flow
1. User selects "Clean" row and presses Enter
2. Controller dispatches CmakeRunner.clean()
3. Deletes entire build directory
4. BUILDS VT updated → menu selectability recomputes

**Clean Action:**
```
Delete Builds/<Generator>/ recursively
```

### Feature: Open IDE/Editor

Single menu row. Label and behavior depend on selected generator:

- **IDE generators (Xcode, Visual Studio):** Row shows "Open IDE". Launches IDE project file.
- **CLI generators (Ninja):** Row shows "Open Editor". Opens nvim in build directory.

The `o` shortcut works for both. Selectable when any generator is selected and build exists.

**Xcode:** `open Builds/Xcode/*.xcodeproj`
**Ninja:** `nvim Builds/Ninja/`
**Visual Studio:** `start Builds/VS2026/*.sln` (post-MVP)

### Feature: Preferences Screen

#### User Flow
1. User presses `/` from main menu
2. preferencesVisible (view transient) → true → mode computes to preferences
3. MainComponent shows Preferences view
4. User navigates and toggles settings (writes directly to CONFIG in State VT)
5. User presses `/` or Esc → preferencesVisible → false → mode computes back to menu

**Preferences Display:**
```
  Auto-update            ON
  Update Interval        10 min
────────────────────────────────────
  Theme                  gfx
────────────────────────────────────
  Back to Menu

up/down navigate | Enter change | / back
```

#### Settings Behavior

##### Auto-update Toggle
- Values: ON / OFF
- Writes CONFIG.autoScanEnabled in State VT
- Controller listens → starts/stops auto-scan timer

##### Update Interval Adjustment
- +/= decreases by 1 min; -/_ increases by 1 min; Shift+= increases by 10 min; Shift+- decreases by 10 min
- Range: 1 min to 60 min
- Writes CONFIG.autoScanInterval in State VT
- Only selectable when auto-scan is ON

##### Theme Cycle
- Values: gfx → spring → summer → autumn → winter → gfx
- Writes CONFIG.theme in State VT
- Controller listens → ThemeLoader updates THEME VT subtree → views repaint

### Feature: Auto-Scan

#### Behavior
- Controller starts `juce::Timer` based on CONFIG.autoScanEnabled / autoScanInterval
- Timer fires → GeneratorDetector + Builds/ scan → updates PROJECT + BUILDS VT subtrees
- Views recompute from VT listeners (menu selectability, labels)
- Skips scan when ASYNC.isActive

#### Edge Cases

##### Edge Case 1: Build Created Externally
**Scenario:** User runs `cmake` manually, creates Builds/Xcode
**Expected Behavior:** Next scan detects it, BUILDS VT updated, Build/Clean/Open become selectable

##### Edge Case 2: Generator Changed Externally
**Scenario:** User deletes Xcode build, creates Ninja build
**Expected Behavior:** BUILDS VT updated, menu reflects available builds

### Feature: Console Mode

#### Behavior
- Active when ASYNC.isActive is true (Mode computes to console)
- `jam::tui::Console` renders streaming output from CmakeRunner
- `jam::tui::Spinner` in Header shows operation type label + braille animation
- ESC aborts: CmakeRunner kills subprocess → ASYNC.isAborted → true, ASYNC.isActive → false → mode returns to menu

### Feature: Braille Banners

Two SVG banners rendered as braille art via `jam::tui` braille primitives:

- **cake-logo.svg** — main banner, rendered in 35% right pane during menu/preferences mode
- **cake-lie.svg** — "the cake is a lie" banner, rendered full-screen when no CMakeLists.txt found (invalidProject mode)

Both embedded as binary data via CMake `BINARY_FILES`.

## UI Specifications

### Layout Structure
```
+------------------------------------------+
| Header (Project directory + spinner)     |
+------------------------------------------+
|         |                                |
|  Menu   |         Braille Banner         |
| (65%)   |            (35%)               |
|         |                                |
+------------------------------------------+
| Footer (Hints/Status)                    |
+------------------------------------------+
```

### Keyboard Shortcuts

| Key | Action | Context |
|---|---|---|
| up/k | Navigate up | Menu/Preferences |
| down/j | Navigate down | Menu/Preferences |
| Enter | Execute/Toggle | All menus |
| Space | Execute/Toggle | All menus |
| g | Generate/Regenerate | Menu |
| b | Build | Menu |
| o | Open IDE / Open Editor | Menu |
| c | Clean | Menu |
| x | Clean All | Menu |
| +/= | Increase interval by 1 min | Preferences (interval row) |
| -/_ | Decrease interval by 1 min | Preferences (interval row) |
| Shift+= | Increase interval by 10 min | Preferences (interval row) |
| Shift+- | Decrease interval by 10 min | Preferences (interval row) |
| / | Toggle preferences | Menu / Preferences |
| Esc | Exit console / Close preferences | Console → Menu, Prefs → Menu |
| Ctrl+C | Quit (2x to confirm) | All modes |

### Color Themes

Five themes loaded from `~/.config/cake/themes/*.xml`. Defaults generated on first run.

#### gfx (Default)
- Background: #0d1117
- Text: #c9d1d9
- Selection: #1f6feb
- Separator: #30363d

#### spring
- Background: #fef5e7
- Text: #2e7d32
- Selection: #81c784
- Separator: #c8e6c9

#### summer
- Background: #e3f2fd
- Text: #0277bd
- Selection: #29b6f6
- Separator: #b3e5fc

#### autumn
- Background: #fff3e0
- Text: #e65100
- Selection: #ff9800
- Separator: #ffe0b2

#### winter
- Background: #f3e5f5
- Text: #4a148c
- Selection: #ab47bc
- Separator: #e1bee7

## Error Handling Matrix

### Global Errors

| Error Condition | User Sees | System Action |
|---|---|---|
| No CMakeLists.txt | CakeLie banner (full-screen) | invalidProject mode |
| CMake not in PATH | "CMake not found in PATH" | Disable all operations |
| Config file corrupt | Config loads with defaults | Auto-recreate config |

### Operation Errors

| Operation | Error Condition | Message | Recovery |
|---|---|---|---|
| Generate | CMake fails | cmake error in console | Stay in console, ESC to return |
| Build | Compilation fails | compiler output in console | Stay in console, ESC to return |
| Clean | Permission denied | "Clean failed: [OS error]" | Return to menu |
| Open IDE | Project not found | "Failed to open IDE: project file not found" | Return to menu |

## Success Criteria

A user can:
- [ ] Select any available CMake generator without typing commands
- [ ] Toggle between Debug and Release configurations visually
- [ ] Generate/regenerate CMake builds with one keypress
- [ ] Build projects with real-time output streaming
- [ ] Clean builds with one keypress
- [ ] Open IDE projects (Xcode/VS) or editor (Ninja) directly from menu
- [ ] Access preferences via `/` key
- [ ] Enable/disable auto-scanning
- [ ] Change themes instantly
- [ ] Navigate with both arrow keys and vim keys (j/k)
- [ ] Exit cleanly with Ctrl+C (twice to confirm)

The system:
- [ ] Detects available generators based on installed tools
- [ ] Maintains correct build path structure (Builds/<Generator>/)
- [ ] Dims/enables menu items based on build state (selectability computed from VT)
- [ ] Persists configuration to ~/.config/cake/config.xml (grafted to State VT on startup)
- [ ] Auto-scans at configured intervals when enabled
- [ ] Handles CMake operations asynchronously with live console output
- [ ] Prevents invalid operations (selectability model)
- [ ] Provides clear error messages for all failure modes
- [ ] Maintains 65/35 split layout with centered content
- [ ] Skips separator rows during navigation automatically

## Architecture Constraints

### Required Patterns
- BLESSED MVC: State VT (Model), jam_tui components (View), MainComponent (Controller)
- Event-driven via ValueTree::Listener — no manual booleans, no manual lambdas
- Views are stateless — read/write State VT, hold only transient render state
- MainComponent is sole orchestrator — listens and dispatches
- Config XML is bootstrap/persistence only — State VT is SSOT while app lives
- Single-page preference menu (no submenus/navigation)
- Conditional selectability based on VT state (all rows always visible)
- System tool detection (not disk scanning for generators)
- Strict build path convention (Builds/<Generator>/)
- Real-time console output via jam::tui::Console
- Persistent configuration in XML → ValueTree

### Forbidden Patterns
- No Elm/single-model patterns — MVC separation enforced
- No shadow state — VT schema is domain truth only, no UI state in VT
- No nested menus or screens (except preferences via `/`)
- No manual build path entry
- No generator options not available on system
- No caching of CMake state (always read fresh via auto-scan)
- No blocking operations on message thread
- No custom build directories outside Builds/ structure
- No manual booleans or lambdas for event dispatch — listener pattern only
