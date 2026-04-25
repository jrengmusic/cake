# PLAN: CAKE-cpp Port

**RFC:** RFC.md
**Date:** 2026-04-25
**BLESSED Compliance:** verified
**Language Constraints:** C++17 / JUCE8 / jam_tui (JRENG-CODING-STANDARD.md applies)

## Overview

Port Go CAKE to C++17/JUCE8 against `jam::tui`, achieving full feature parity. Go sources archived to `___legacy___/`. macOS first.

## Contract

- **jam_tui is the product.** Modules must be generic, universal, reusable for ANY TUI app. CAKE is the first consumer that proves it.
- **Current jam_tui state is unusable.** Modules may be rewritten, removed, added. Any jam codebase is touchable — jam_tui, jam_core, jam_subprocess, whatever it takes.
- **No debt. No deferral.** If a primitive is broken, wrong, or missing — fix it now, in jam, in this sprint. No "jam PR later."
- **CAKE Go is the reference.** Port to jam C++ with BLESSED compliance. No new patterns invented. Go behavior is the specification.
- **TIT-cpp established patterns apply** for project structure, state architecture (atoms + VT + flush), threading model, CMake setup.

## Language / Framework Constraints

- JRENG-CODING-STANDARD.md enforced: allman braces, `not`/`and`/`or` tokens, `.at()` container access, brace initialization, no anonymous namespaces, no `namespace detail`
- MANIFESTO.md enforced fully (C++ is reference implementation, no LANGUAGE.md overrides)
- No early returns. Positive nested checks. `jassert` at boundaries.
- `DONT_SET_USING_JUCE_NAMESPACE=1` — all JUCE types fully qualified
- Threading: message thread owns state + views, subprocess worker owns ChildProcess, `callAsync` only crossing

## Validation Gate

Each step validated by @Auditor against:
- MANIFESTO.md (BLESSED principles)
- NAMES.md (naming philosophy — Rule -1: no improvised names)
- JRENG-CODING-STANDARD.md (C++ coding standards)
- Locked PLAN decisions (no deviation, no scope drift)
- `cmake --build` compiles clean with zero warnings

## Steps

### Step 0: Write SPEC.md

**Scope:** `SPEC.md` at project root
**Action:** COUNSELOR writes SPEC for C++/JUCE/jam. Go SPEC is requirements source — UI/UX identical, feature set unchanged. Architecture is BLESSED MVC: State (VT) is Model, jam_tui components are View, MainComponent is Controller. No Elm patterns. VT schema is domain truth only:

```
CAKE                                    # root
├── PROJECT                             # domain state
│   ├── workingDirectory    (string)
│   ├── hasCMakeLists       (bool)
│   ├── selectedGenerator   (string)
│   ├─��� configuration       (string)
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
├── CONFIG                              # persisted preferences
│   ├── autoScanEnabled     (bool)
│   ├── autoScanInterval    (int)
│   └── theme               (string)
└── THEME                               # loaded from ~/.config/cake/themes/*.xml
    └── ...                             # color properties
```

View-owned transient state (NOT in VT): menuIndex, spinnerFrame, scrollOffset, preferencesVisible. Mode computed from VT (isActive → console, not hasCMakeLists → invalidProject, preferencesVisible → preferences, else → menu).

**Validation:** SPEC.md exists, covers all Go SPEC features, zero Go/Elm terminology, VT schema is domain-only.

### Step 1: Scaffold + Legacy Archive

**Scope:** `___legacy___/`, `CMakeLists.txt`, `Assets/`, `Source/Main.cpp`, `Source/MainComponent.h/.cpp`

**Structure (TIT-cpp verbatim + CAKE-specific subdirs):**
```
Source/
  Main.cpp                  # juce::JUCEApplication entry (TIT pattern)
  MainComponent.h/.cpp      # root jam::tui::Component, owns layout
  Identifier.h              # namespace ID — all juce::Identifier constants
  state/
    Axis.h                  # enums + toString/parse
    State.h/.cpp            # APVTS-mirror: atoms + VT + timer flush (~10 atoms, single file)
  component/                # all view components (TIT: component/)
    Banner.h/.cpp
    CakeLie.h/.cpp
    Footer.h/.cpp
    Header.h/.cpp           # includes spinner during active ops
    Menu.h/.cpp
    Preferences.h/.cpp
  menu/                     # menu building (TIT: menu/)
    MenuBuilder.h/.cpp
    MenuItems.h
  ops/                      # cmake execution
    CmakeRunner.h/.cpp      # single file, methods per op (generate/build/clean/open)
  detect/                   # generator detection
    GeneratorDetector.h/.cpp
    MsvcEnvironment.h/.cpp
  theme/                    # theme loading
    ThemeLoader.h/.cpp
  tests/                    # inside Source/ (TIT pattern)
    TestRunner.cpp
    fixtures/
```

**Action:**
- Move Go sources (`cmd/`, `internal/`, `go.mod`, `go.sum`, `build.sh`, `release.sh`, `scripts/`, `ARCHITECTURE.md`, `SPEC.md`, `README.md`, `RELEASE_NOTES.md`, `.goreleaser.yaml`) to `___legacy___/`
- Copy SVGs from `internal/ui/assets/` to `Assets/` (cake-logo.svg, cake-lie.svg)
- Create `CMakeLists.txt` from RFC section 4.2 (`configure_app` with `cakec` target)
- Create `Source/Main.cpp` — `juce::JUCEApplication` entry (TIT pattern)
- Create `Source/MainComponent.h/.cpp` — root `jam::tui::Component` stub
- Builds green: `cmake -S . -B Builds/Ninja -G Ninja && cmake --build Builds/Ninja`
**Validation:** Compiles clean. `cakec` binary runs and exits. SVGs in `Assets/`. Go sources untouched in `___legacy___/`.

### Step 2: State Layer

**Scope:** `Source/Identifier.h`, `Source/state/Axis.h`, `Source/state/State.h/.cpp`
**Action:**
- `Identifier.h` — all `juce::Identifier` constants for VT schema (namespace `ID`, sectioned by node type, TIT pattern)
- `Axis.h` — `enum class Mode { menu, preferences, console, invalidProject }` (computed from VT state, not stored), `enum class Generator { xcode, ninja, vs2026, vs2022 }`, `enum class Configuration { debug, release }`, `enum class OpType { none, build, generate, clean, cleanAll, regenerate }` + `toString()`/`parse()` free functions
- `State.h/.cpp` — APVTS-mirror: ~10 atoms + ValueTree root (domain truth only: PROJECT, BUILDS, ASYNC, CONFIG, THEME), `juce::Timer` flush at 16ms. Single file — split only if L-violated.
- `Main.cpp` updated to own `State` via `std::unique_ptr`
**Validation:** Compiles clean. VT schema matches `Identifier.h`. Flush timer fires. Atom → VT round-trip works (unit test).

### Step 3: Generator Detection

**Scope:** `Source/detect/GeneratorDetector.h/.cpp`, `Source/detect/MsvcEnvironment.h/.cpp` (stub)
**Action:**
- `GeneratorDetector` — detect available generators via system tool presence: `xcodebuild` (macOS), `ninja`, `vswhere.exe` (Windows). Returns `juce::Array<Generator>`.
- `MsvcEnvironment.h/.cpp` — macOS stub (no-op). Windows implementation deferred to post-MVP sprint.
- Wire into `State` — detected generators populate `PROJECT.GENERATORS[]` subtree on construction.
**Validation:** Compiles clean. On macOS: detects Xcode (if installed) and Ninja (if installed). Generator list in VT matches system reality.

### Step 4: Views — Menu + Layout

**Scope:** `Source/MainComponent.h/.cpp`, `Source/component/Menu.h/.cpp`, `Source/component/Header.h/.cpp`, `Source/component/Footer.h/.cpp`, `Source/component/Banner.h/.cpp`, `Source/menu/MenuBuilder.h/.cpp`, `Source/menu/MenuItems.h`
**Action:**
- `MainComponent` — root `jam::tui::Component`, owns `jam::tui::SplitPane` (65/35), dispatches input via `jam::tui::Input`
- `Menu` component — 8 fixed rows via `jam::tui::Menu`, selectability computed from VT PROJECT + BUILDS listeners
- `MenuBuilder` — builds menu rows from VT state. `MenuItems.h` — row ID constants.
- `Header` — project directory display
- `Footer` — context hints + status line
- `Banner` — cake-logo.svg via `jam::tui/braille/`, rendered in 35% right pane
- `Main.cpp` updated to own `MainComponent`, wire `State` VT to views
**Validation:** Compiles clean. `cakec` launches, renders menu with 8 rows, banner in right pane, header shows cwd, footer shows hints. Navigation works (j/k/arrows). Unselectable rows dimmed and skipped.

### Step 5: Views — Preferences + CakeLie

**Scope:** `Source/component/Preferences.h/.cpp`, `Source/component/CakeLie.h/.cpp`
**Action:**
- Console mode uses `jam::tui::Console` directly in `MainComponent` layout. `CmakeRunner` calls `appendLine()`. Spinner rendered by `Header` when ASYNC.isActive.
- `Preferences` — stateless View. Reads/writes CONFIG properties in State VT (autoScanEnabled, autoScanInterval, theme). Same pattern as Menu — pure View, no logic.
- `CakeLie` — stateless View. Full-screen cake-lie.svg banner via braille.
- Mode routing: `MainComponent` computes Mode from VT state and shows the correct view.
**Validation:** Compiles clean. `/` toggles preferences. Mode switching works. CakeLie banner shows when launched outside cmake project.

### Step 6: Ops — cmake execution

**Scope:** `Source/ops/CmakeRunner.h/.cpp`
**Action:**
- `CmakeRunner` — subprocess orchestrator via `jam_subprocess`. Owns `juce::ChildProcess`. Streams stdout/stderr lines to `jam::tui::Console` via `callAsync`. Updates ASYNC atoms (isActive, currentOp) → flushed to VT. Supports abort via process termination.
- Methods: `generate()`, `build()`, `clean()`, `open()` — each assembles the correct cmake command
- `generate`: `cmake -S . -B Builds/<Gen> -G <Gen>`
- `build`: `cmake --build Builds/<Gen> --config <Cfg>`
- `clean`: `juce::File::deleteRecursively()` on `Builds/<Gen>/`
- `open`: `open *.xcodeproj` (Xcode), `nvim Builds/Ninja/` (Ninja), `start *.sln` (VS, deferred)
- Wire into `MainComponent` — menu actions dispatch ops
**Validation:** Compiles clean. Generate creates `Builds/<Gen>/`. Build compiles a cmake project. Clean removes build dir. Open launches IDE/editor. Console streams output in real-time. Abort (ESC) kills subprocess.

### Step 7: Integration — Config Bootstrap + Auto-scan + Theme

**Scope:** `Source/theme/ThemeLoader.h/.cpp`, config bootstrap in `Main.cpp`
**Action:**
- Config bootstrap: on startup, load `~/.config/cake/config.xml` → graft into State VT CONFIG subtree. Write defaults if missing. While app lives, State VT is only SSOT. `ValueTree::Listener` on CONFIG subtree persists changes back to XML. No separate config logic.
- Auto-scan: `MainComponent` (Controller) listens to CONFIG.autoScanEnabled/autoScanInterval in VT. Drives `juce::Timer` that calls `GeneratorDetector` + scans `Builds/`, updates PROJECT + BUILDS VT subtrees. Skips when ASYNC.isActive.
- `ThemeLoader` — loads `~/.config/cake/themes/*.xml`, hot-reload via `jam::File::Watcher`. Updates THEME VT subtree. Defaults generated if missing (5 themes).
- Generator cycling / Configuration toggle — Menu View writes directly to State VT on user input. `MainComponent` listens and dispatches side effects.
**Validation:** Compiles clean. Config persists across restarts. Auto-scan detects external build changes. Theme switching works. Hot-reload works.

### Step 8: Test Suite

**Scope:** `tests/`
**Action:**
- State tests: VT schema validity, atom → flush round-trip, selectability computation from state permutations
- Generator detection tests: mock tool paths
- Menu tests: row generation, selectability rules match SPEC
- Fixture-based: ValueTree snapshots for known state permutations
**Validation:** All `juce::UnitTest` tests pass. Coverage matches Go test suite scope (state, config, UI, menu).

### Step 9: macOS Release

**Scope:** Build/release infrastructure
**Action:**
- Install target or release script for `cakec` binary
- Codesign + notarize (reuse entitlements.plist from Go release)
- Verify `cakec` runs from install location
**Validation:** Signed binary runs on clean macOS. No Gatekeeper warnings.

## Sprint Grouping (suggested)

| Sprint | Steps | Deliverable |
|---|---|---|
| Sprint 1 | Steps 0-3 | SPEC + scaffold + state + detection — builds green, VT schema populated |
| Sprint 2 | Steps 4-6 | Views + ops — functional CAKE with cmake execution |
| Sprint 3 | Steps 7-9 | Config + themes + auto-scan + tests + macOS release |

## BLESSED Alignment

- **B (Bound):** `Main.cpp` owns `State` + `MainComponent` via `std::unique_ptr`. Subprocess owned by `CmakeRunner`. Timer owned by `State`. FileWatcher owned by `ThemeLoader`. Thread ownership explicit.
- **L (Lean):** 300/30/3 enforced. 8 fixed menu rows from data, not switch chains. Generator detection is lookup. Views are thin — many small files, each stateless.
- **E (Explicit):** Zero early returns. All parameters visible. Constants in `Identifier.h`. `jassert` at boundaries. No silent fails. Event-driven via `ValueTree::Listener` — no manual booleans, no manual lambdas.
- **S (SSOT):** State VT is the only SSOT while app lives. Config XML is bootstrap/persistence layer only. No shadow state between config and runtime.
- **S (Stateless):** All Views are stateless — read/write State VT, hold only transient render state (spinnerFrame, scrollOffset, menuIndex). No View remembers anything for the Controller.
- **E (Encapsulation):** `jam_tui` imports no CAKE headers. `ops/` imports no `component/`. Unidirectional: Input → View → State VT → Listener → Controller dispatches. `MainComponent` (Controller) is the sole orchestrator — listens to VT, dispatches callbacks. Views never call ops directly.
- **D (Deterministic):** Same VT state → same render. Same cmake command → same output. Emergent from BLESSE.

## Risks

- **Windows:** macOS first. MSVC environment capture, vswhere detection, VS generator — post-MVP sprint (ARCHITECT decision).
