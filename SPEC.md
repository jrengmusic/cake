# cake Specification v2.0

## Overview

**Purpose:** CMake project management tool with visual preference-style TUI for quick project configuration and building  
**Target User:** Developers who want visual CMake control without typing commands  
**Core Workflow:** Select generator → Configure → Build/Clean/Open, with persistent settings

## Technology Stack

- **Language:** Go
- **Framework:** Bubble Tea (TUI), Lip Gloss (styling)
- **Platform:** Cross-platform (macOS, Linux, Windows)
- **Dependencies:** 
  - github.com/charmbracelet/bubbletea
  - github.com/charmbracelet/lipgloss
  - github.com/pelletier/go-toml/v2 (config persistence)
- **External Requirements:** CMake installed in PATH

## Core Principles

1. **Preference-Style Interface**: Single-page menu with toggleable values (like TIT), no nested submenus
2. **Dynamic Visibility**: Menu items appear/disappear based on build state and generator capabilities
3. **System Tool Detection**: Available generators determined by installed tools, not disk scanning
4. **50/50 Split Layout**: Menu on left half, ASCII banner on right half, both centered
5. **Build Path Convention**: Strict `Builds/<Generator>/<Config?>` structure (multi-config vs single-config)
6. **Persistent Configuration**: Settings saved to `~/.config/cake/config.toml`

## State Model

### Project State Determination
```
State determined by: (WorkingDirectory, AvailableGenerators, BuildDirectories)
- WorkingDirectory: Current directory containing CMakeLists.txt
- AvailableGenerators: System tools detected (Xcode, Ninja, Unix Makefiles, Visual Studio)
- BuildDirectories: Scanned `Builds/` subdirectories matching generator pattern
```

### Generator Types
```
Multi-Config Generators (build contains all configurations):
- Xcode (macOS only, IDE)
- Visual Studio (Windows only, IDE)
Path: Builds/<Generator>/

Single-Config Generators (separate build per configuration):
- Ninja (cross-platform, CLI)
- Unix Makefiles (always available, CLI)  
Path: Builds/<Generator>/<Configuration>/
```

### Menu Item Visibility Rules
```
Generator: ALWAYS visible
Generate/Regenerate: visible if CMakeLists.txt exists
Open IDE: visible if IDE generator AND build exists
Open Editor: visible if CLI generator AND build exists
Configuration: ALWAYS visible
Build: visible if build exists AND configured
Clean: visible if build exists
```

## Feature Specifications

### Feature: Main Menu Navigation

#### User Flow (Happy Path)
1. User launches `cake` in directory with CMakeLists.txt
2. System shows preference menu with available options
3. User presses ↑/↓ or j/k to navigate rows
4. System highlights selected row with inverted colors
5. User presses Enter or Space on a row
6. System executes toggle (for settings) or action (for operations)

**UI Display:**
```
Project directory:
/Users/username/myproject

⚙️  Generator              Xcode           
🚀  Generate                               
📂  Open IDE                               
────────────────────────────────────────
🏗️  Configuration          Debug           
🔨  Build                                  
🧹  Clean                                  

↑↓ navigate │ Enter select │ Ctrl+C quit │ / preferences
```

**User Input:**
- ↑/k: Move selection up (skip separators)
- ↓/j: Move selection down (skip separators)
- Enter/Space: Toggle value or execute action
- /: Open preferences screen
- Ctrl+C: Quit (confirm if pressed twice)

**System Response:**
- Selection moves with smooth highlight
- Separators cannot be selected (auto-skip)
- Actions trigger immediately on Enter

#### Edge Cases

##### Edge Case 1: No CMakeLists.txt
**Scenario:** User runs cake in directory without CMakeLists.txt
**Expected Behavior:** Show error in footer, limited menu
**Error Message:** "No CMakeLists.txt found in current directory"

##### Edge Case 2: Empty Builds Directory
**Scenario:** No builds exist yet
**Expected Behavior:** Show only Generator selection and Generate action
**Menu State:** Build, Clean, Open IDE/Editor hidden

##### Edge Case 3: Navigation at Boundaries
**Scenario:** User presses up at first item or down at last
**Expected Behavior:** Selection stays at boundary (no wrap)

### Feature: Generator Selection

#### User Flow (Happy Path)
1. User selects Generator row
2. User presses Enter or Space
3. System cycles to next available generator
4. Display updates immediately with new generator name

**Available Generators Detection:**
- macOS: Xcode (if xcodebuild exists), Ninja (if ninja exists), Unix Makefiles (always)
- Linux: Ninja (if ninja exists), Unix Makefiles (always)
- Windows: Visual Studio (if vswhere.exe exists), Ninja (if ninja exists), Unix Makefiles (always)

**Cycle Order:**
```
Xcode → Ninja → Unix Makefiles → Xcode (loop)
```

#### Error Handling

| Error Condition | User Sees | System Action |
|-----------------|-----------|---------------|
| No generators available | "No CMake generators found" | Disable generate action |
| Only one generator | Generator name (no cycling) | Enter does nothing |

### Feature: Configuration Toggle

#### User Flow
1. User selects Configuration row
2. User presses Enter or Space
3. System toggles between Debug ↔ Release
4. Display updates immediately

**Value Display:**
- Shows: "Debug" or "Release"
- Toggle is bidirectional (Debug → Release → Debug)

### Feature: Generate/Regenerate Operation

#### User Flow (Generate - First Time)
1. User selects "Generate" row
2. User presses Enter
3. System switches to console mode
4. CMake executes with selected generator and configuration
5. Output streams in real-time
6. On completion, returns to menu with status message

**CMake Command (Single-Config):**
```bash
cmake -S . -B Builds/Ninja/Debug -G Ninja -DCMAKE_BUILD_TYPE=Debug
```

**CMake Command (Multi-Config):**
```bash
cmake -S . -B Builds/Xcode -G Xcode
```

#### User Flow (Regenerate - Build Exists)
1. User selects "Regenerate" row (label changes when build exists)
2. User presses Enter
3. Same as Generate flow but overwrites existing build

**UI During Operation:**
```
[Console Output Mode]
Setting up CMake...
-- The C compiler identification is AppleClang 14.0.0
-- Detecting C compiler ABI info
[... cmake output streams ...]

Setting up CMake... (ESC to abort)
```

#### Edge Cases

##### Edge Case 1: CMake Not Installed
**Scenario:** cmake command not found in PATH
**Expected Behavior:** Show error immediately
**Error Message:** "CMake not found in PATH"

##### Edge Case 2: Invalid CMakeLists.txt
**Scenario:** CMakeLists.txt has syntax errors
**Expected Behavior:** Show CMake error output
**Footer Message:** "Generate failed: CMake configuration error"

##### Edge Case 3: Disk Full
**Scenario:** No space to create build directory
**Expected Behavior:** Show system error
**Error Message:** "Generate failed: [OS error message]"

### Feature: Build Operation

#### User Flow
1. User selects "Build" row (only visible if build exists)
2. User presses Enter
3. System switches to console mode
4. CMake build executes for selected configuration
5. Compiler output streams in real-time
6. On completion, returns to menu with status message

**Build Command (Single-Config):**
```bash
cmake --build Builds/Ninja/Debug
```

**Build Command (Multi-Config):**
```bash
cmake --build Builds/Xcode --config Debug
```

**Success Message:** "Operation completed. Press ESC to return."
**Failure Message:** "Build failed: [error summary]"

#### Edge Cases

##### Edge Case 1: Build Directory Deleted
**Scenario:** User deletes build directory externally
**Expected Behavior:** Auto-scan detects missing directory, hides Build option
**Auto Recovery:** Menu updates within scan interval

##### Edge Case 2: Compilation Errors
**Scenario:** Source code has errors
**Expected Behavior:** Show full compiler output, remain in console
**Footer:** "Build failed: compilation errors"

### Feature: Clean Operation

#### User Flow
1. User selects "Clean" row
2. User presses Enter
3. System deletes entire build directory
4. Menu regenerates (Clean/Build/Open options disappear)
5. Footer shows completion message

**Clean Command:**
```bash
rm -rf Builds/<Generator>/<Config?>
```

**Success Message:** "Operation completed. Press ESC to return."

### Feature: Open IDE/Editor

#### IDE Flow (Xcode, Visual Studio)
1. User selects "Open IDE" row (only for IDE generators)
2. User presses Enter
3. System launches IDE with project file
4. cake continues running, shows status

**Xcode Command:**
```bash
open Builds/Xcode/*.xcodeproj
```

**Visual Studio Command:**
```bash
start Builds/VisualStudio/*.sln
```

#### Editor Flow (Ninja, Makefiles)
1. User selects "Open Editor" row (only for CLI generators)
2. User presses Enter
3. System launches Neovim in build directory
4. cake waits for editor to close
5. Returns to menu with "Editor closed" message

**Neovim Command:**
```bash
nvim Builds/<Generator>/<Config>/
```

### Feature: Preferences Screen

#### User Flow
1. User presses `/` from main menu
2. System shows preferences screen (replaces menu)
3. User navigates and toggles settings
4. User presses `/` or Esc to return to main menu

**Preferences Display:**
```
🔄  Auto-update            ON
⏱️  Update Interval        10 min
────────────────────────────────────
🎨  Theme                  gfx
────────────────────────────────────
←  Back to Menu

↑↓ navigate │ Enter change │ / back
```

#### Settings Behavior

##### Auto-update Toggle
- Values: ON ↔ OFF
- When ON: Project state refreshes at interval
- When OFF: No automatic scanning

##### Update Interval Cycle
- Values: 5 min → 10 min → 15 min → 30 min → 5 min
- Only visible when Auto-update is ON

##### Theme Cycle
- Values: gfx → spring → summer → autumn → winter → gfx
- Changes colors immediately on selection

### Feature: Auto-Scan

#### Behavior
- Triggers every N minutes (based on Update Interval)
- Refreshes project state (detects new/deleted builds)
- Updates menu if build state changed
- Shows "[Scanning...]" in footer during scan
- Skips scan if async operation is active

#### Edge Cases

##### Edge Case 1: Build Created Externally
**Scenario:** User runs `cmake` manually, creates Builds/Xcode
**Expected Behavior:** Next scan detects it, shows Build/Clean/Open options

##### Edge Case 2: Generator Changed Externally
**Scenario:** User deletes Xcode build, creates Ninja build
**Expected Behavior:** Menu updates to reflect available builds

### Feature: Configuration Persistence

#### Storage Location
```
~/.config/cake/config.toml
```

#### File Format
```toml
[preferences]
auto_scan_enabled = true
auto_scan_interval = 10

[appearance]
theme = "gfx"

[project]
last_generator = "Xcode"
last_configuration = "Debug"
```

#### Persistence Rules
- Settings save immediately on change
- Missing config file created on first run
- Invalid config reverts to defaults

## UI Specifications

### Layout Structure
```
┌─────────────────────────────────────────┐
│ Header (Project directory)              │
├─────────────────────────────────────────┤
│         │                               │
│  Menu   │         ASCII Banner          │
│ (50%)   │            (50%)              │
│         │                               │
├─────────────────────────────────────────┤
│ Footer (Hints/Status)                   │
└─────────────────────────────────────────┘
```

### Keyboard Shortcuts

| Key | Action | Context |
|-----|--------|---------|
| ↑/k | Navigate up | Menu/Preferences |
| ↓/j | Navigate down | Menu/Preferences |
| Enter | Execute/Toggle | All menus |
| Space | Execute/Toggle | All menus |
| / | Toggle preferences | Menu ↔ Preferences |
| Esc | Exit/Cancel | Console → Menu, Prefs → Menu |
| Ctrl+C | Quit (2x to confirm) | All modes |

### Color Themes

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
|-----------------|-----------|---------------|
| No CMakeLists.txt | "No CMakeLists.txt found in current directory" | Show limited menu |
| CMake not in PATH | "CMake not found in PATH" | Disable all operations |
| Config file corrupt | Config loads with defaults | Auto-recreate config |
| Terminal too small | Content truncated | Adjust layout dynamically |

### Operation Errors

| Operation | Error Condition | Message | Recovery |
|-----------|----------------|---------|----------|
| Generate | CMake fails | "Generate failed: [cmake error]" | Stay in console, ESC to return |
| Build | Compilation fails | "Build failed: [error summary]" | Show full output, ESC to return |
| Clean | Permission denied | "Clean failed: [OS error]" | Return to menu |
| Open IDE | Project not found | "Failed to open IDE: project file not found" | Return to menu |
| Open Editor | Neovim not found | "Failed to open editor: nvim not found" | Return to menu |

## Success Criteria

A user can:
- [x] Select any available CMake generator without typing commands
- [x] Toggle between Debug and Release configurations visually
- [x] Generate/regenerate CMake builds with one keypress
- [x] Build projects with real-time output streaming
- [x] Clean builds completely with confirmation
- [x] Open IDE projects (Xcode/VS) directly from menu
- [x] Open editor in build directory for CLI generators
- [x] Access preferences via `/` key
- [x] Enable/disable auto-scanning
- [x] Change themes instantly
- [x] Navigate with both arrow keys and vim keys (j/k)
- [x] Exit cleanly with Ctrl+C (twice to confirm)

The system:
- [x] Detects available generators based on installed tools
- [x] Maintains correct build path structure (Builds/<Generator>/<Config?>)
- [x] Shows/hides menu items based on build state
- [x] Persists configuration to ~/.config/cake/config.toml
- [x] Auto-scans at configured intervals when enabled
- [x] Handles CMake operations asynchronously with live output
- [x] Prevents invalid operations (can't build without generate)
- [x] Provides clear error messages for all failure modes
- [x] Maintains 50/50 split layout with centered content
- [x] Skips separator rows during navigation automatically

## Architecture Constraints

### Required Patterns
- ✅ Single-page preference menu (no submenus/navigation)
- ✅ Dynamic visibility based on state
- ✅ System tool detection (not disk scanning for generators)
- ✅ Strict build path convention (Builds/<Generator>/<Config?>)
- ✅ Real-time output streaming for operations
- ✅ Persistent configuration in TOML

### Forbidden Patterns
- ❌ No nested menus or screens (except preferences via `/`)
- ❌ No manual build path entry
- ❌ No generator options not available on system
- ❌ No caching of CMake state (always read fresh)
- ❌ No blocking operations in UI thread
- ❌ No custom build directories outside Builds/ structure