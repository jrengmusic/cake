# RFC — CAKE-cpp Port to C++/JUCE via `jam::tui`

**Date:** 2026-04-25
**Status:** Ready for COUNSELOR handoff
**Author:** BRAINSTORMER
**Project:** CAKE — CMake Project Manager TUI
**Path:** `~/Documents/Poems/dev/cake/`

---

## 1. Problem Statement

CAKE Go ships and works. Same stack fracture as TIT Go before its port:

- Go toolchain outside ARCHITECT's `jam_*` ecosystem
- Not BLESSED-auditable in ARCHITECT's native language
- Cannot share `jam_*` substrate consumed by END, TIT-cpp, CAROLINE, Kuassa, whatdbg
- Cannot benefit from `jam_tui` primitives already forged and battle-tested during TIT-cpp port
- bubbletea/Elm model is a worse fit than `juce::ValueTree` listener-driven observable state

**Port target:** full feature parity with Go CAKE, rewritten in C++17/JUCE8 against `jam::tui`. Go sources preserved under `___legacy___/` for reference and continued shipping.

**Strategic value:**

1. **Daily driver feedback loop:** ARCHITECT is building jam with cmake constantly. CAKE-cpp becomes the first `jam_tui` app used daily, exposing framework issues under real sustained use before TIT-cpp resumes its complex Sprint 3+.
2. **jam_tui hardening:** CAKE exercises Menu, ConsoleStream, ThemeResolver, Spinner, SplitPane against simple app logic. Battle-tests primitives without compounding complex state machine debugging.
3. **Ecosystem convergence:** Second `jam_tui` consumer validates framework generality. If CAKE needs a primitive change, it surfaces before TIT-cpp's complex views depend on the current shape.
4. **Reduced complexity:** CAKE Go is ~7,700 LOC / 59 files. Single-page preference menu, 8 fixed rows, no FSMs, no protocol state machines, no conflict resolution. 1-2 sprint port vs TIT's 6-sprint arc.

---

## 2. Research Summary

### 2.1 Reference implementations available

- **Go CAKE** (`~/Documents/Poems/dev/cake/`, ~7,700 LOC, 59 files) — complete SPEC.md + ARCHITECTURE.md. Zero design decisions left; port is pure translation.
- **TIT-cpp** (`~/Documents/Poems/dev/tit/`) — established patterns: `configure_app()` CMake, APVTS-mirror state, ValueTree schema, theme XML + hot-reload, `jam_tui` primitive consumption. All TIT patterns apply to CAKE per ARCHITECT decision.
- **END** (`~/Documents/Poems/dev/end/`) — architectural template. `Terminal::State` pattern (atoms + VT + timer flush). Production-proven.

### 2.2 Go CAKE layer inventory (verified 2026-04-25)

| Layer | Package | LOC | Files | CAKE-cpp mapping |
|---|---|---|---|---|
| Entry | `cmd/cake/` | ~20 | 1 | `Source/Main.cpp` + `Source/CakeApp.h/.cpp` |
| App | `internal/app/` | 1,949 | 20 | `Source/app/` — Bubble Tea Model → APVTS-mirror state |
| UI | `internal/ui/` | 2,732 | 18 | `Source/view/` — lipgloss rendering → `jam_tui` primitives |
| State | `internal/state/` | 1,000 | 4 | `Source/state/` — ProjectState → ValueTree |
| Ops | `internal/ops/` | 274 | 4 | `Source/ops/` — cmake subprocess via `jam_subprocess` |
| Utils | `internal/utils/` | 504 | 4 | `Source/utils/` — generator constants, MSVC detection |
| Config | `internal/config/` | 338 | 1 | ValueTree + XML (TIT pattern, no TOML) |
| Banner | `internal/banner/` | 873 | 4 | `jam_tui/braille/` (already implemented) |

### 2.3 jam module readiness (verified 2026-04-25)

| Module | LOC | Status | CAKE needs |
|---|---|---|---|
| `jam_core` | 15,550+ | Implemented | Yes — utilities, file watcher, identifiers |
| `jam_data_structures` | 895+ | Implemented | Yes — ValueTree wrapper (transitive dep of jam_tui) |
| `jam_tui` | 2,028 (primitives) | Implemented | Yes — Menu, ConsoleStream, SplitPane, Spinner, ThemeResolver, TextPane |
| `jam_markdown` | 1,048 | Implemented | Transitive dep of jam_tui at CMake level |
| `jam_subprocess` | 367 | Implemented (macOS) | Yes — cmake execution |

All modules CAKE needs are implemented. Zero module forging required — CAKE is a pure consumer.

### 2.4 Non-novelty

CAKE-cpp is pure translation:

- **Requirements:** new SPEC.md to be written for jam (ARCHITECT decision). Go SPEC.md is the source of requirements truth.
- **Architecture pattern:** locked by TIT-cpp's established patterns (APVTS-mirror, ValueTree schema, configure_app CMake)
- **Framework base:** locked by `jam_tui` (primitives already forged by TIT-cpp)
- **Config format:** XML → `juce::ValueTree` native (TIT pattern)
- **Theme format:** external XML + hot-reload via `jam::File::Watcher` (TIT pattern)

Zero architectural invention required. Zero module forging required.

---

## 3. Principles and Rationale

### 3.1 Daily driver first (per BLESSED — D, Deterministic)

CAKE-cpp will be used daily during jam reorganization. Every cmake generate/build/clean cycle tests `jam_tui` Menu, ConsoleStream, Spinner under real load. Defects surface faster than TIT-cpp's deferred Sprint 3+ demo milestone.

### 3.2 TIT patterns apply wholesale (ARCHITECT decision 2026-04-25)

All established TIT-cpp patterns carry forward unchanged:

- `configure_app()` CMake via `jam/cmake/AppBuilder.cmake`
- APVTS-mirror state: atoms + ValueTree + `juce::Timer` flush
- Theme: external XML at `~/.config/cake/themes/*.xml` + hot-reload via `jam::File::Watcher`
- Config: XML → ValueTree native
- Binary name: `cakec` during port (avoids PATH collision with Go `cake`)
- Go sources → `___legacy___/` (continues shipping)
- Threading: message thread owns state + views, subprocess worker owns `juce::ChildProcess`, `callAsync` only crossing primitive
- Test framework: `juce::UnitTest`
- macOS-first, Windows deferred to later sprint

### 3.3 Simplicity advantage (per BLESSED — L, Lean)

CAKE vs TIT complexity comparison:

| Dimension | TIT | CAKE |
|---|---|---|
| State axes | 5 (WorkingTree × Timeline × Operation × Remote × IsTitTimeTravel) | 2 (ProjectState × BuildState) |
| Menu items | 27 conditional | 8 fixed (always visible, selectability toggles) |
| Protocol FSMs | 4 (DirtyOp, TimeTravel, Conflict, SetupWizard) | 0 |
| Parsers | 4 (Porcelain, Log, Diff, ConflictMarker) | 0 |
| Browser views | 3 (History 2-col, FileHistory 3-pane, ConflictResolver) | 0 |
| Subprocess complexity | git state detection, marker files, streaming parsers | cmake execute, stream stdout/stderr |
| Async operations | git commands + protocol orchestration | cmake generate/build + IDE/editor launch |

CAKE exercises the same `jam_tui` primitives against dramatically simpler application logic. Framework bugs are isolated from app bugs.

### 3.4 Layering — same as TIT (per BLESSED E, Encapsulation)

`jam_tui` is View+Input only. `CakeState` lives in `Source/state/` — domain-specific application code. Views observe ValueTree via listeners. Ops layer executes cmake via `jam_subprocess`. Unidirectional flow: App → State → Ops, App → View (observes State VT).

### 3.5 Config format — XML → ValueTree native (TIT pattern)

Go CAKE uses TOML at `~/.config/cake/config.toml`. C++ port uses XML → ValueTree. Same rationale as TIT RFC §3.5:
- JUCE-native: `juce::XmlDocument::parse()` + `juce::ValueTree::fromXml()`
- Zero new module dependencies
- Config becomes a ValueTree. Hot reload via FileWatcher → parse → VT → Listeners.

---

## 4. Scaffold

### 4.1 Project structure

```
~/Documents/Poems/dev/cake/
├── ___legacy___/                    # Go sources archived verbatim
│   ├── cmd/
│   ├── internal/
│   ├── go.mod, go.sum
│   ├── ARCHITECTURE.md
│   ├── SPEC.md                      # Go SPEC (requirements source)
│   └── ...
├── CMakeLists.txt                   # jam configure_app()
├── Assets/
│   ├── cake-logo.svg                # main banner
│   └── cake-lie.svg                 # "the cake is a lie" banner
├── Source/
│   ├── Main.cpp                     # juce::JUCEApplication entry
│   ├── CakeApp.h/.cpp               # headless app, owns CakeState + CakeScreen + CmakeRunner
│   ├── CakeIdentifier.h             # all juce::Identifier constants (SSOT)
│   ├── state/
│   │   ├── CakeState.h/.cpp         # APVTS-mirror: atoms + ValueTree + timer flush
│   │   └── CakeAxis.h               # enum: AppMode, Generator, Configuration
│   ├── ops/
│   │   ├── CmakeRunner.h/.cpp       # subprocess orchestrator via jam_subprocess
│   │   ├── GenerateOp.h/.cpp        # cmake -S . -B Builds/<Gen> -G <Gen>
│   │   ├── BuildOp.h/.cpp           # cmake --build Builds/<Gen> --config <Cfg>
│   │   ├── CleanOp.h/.cpp           # rm -rf Builds/<Gen>
│   │   └── OpenOp.h/.cpp            # open IDE/editor
│   ├── detect/
│   │   ├── GeneratorDetector.h/.cpp  # system tool detection (xcodebuild, ninja, vswhere)
│   │   └── MsvcEnvironment.h/.cpp   # vcvarsall capture (Windows)
│   ├── view/
│   │   ├── CakeScreen.h/.cpp        # root tui::Component, owns layout
│   │   ├── Header.h/.cpp            # project directory display
│   │   ├── Footer.h/.cpp            # context hints + status
│   │   ├── MenuView.h/.cpp          # 8-row preference menu (jam_tui::Menu)
│   │   ├── BannerView.h/.cpp        # braille SVG banner (jam_tui braille)
│   │   ├── ConsoleView.h/.cpp       # streaming cmake output (jam_tui::ConsoleStream)
│   │   ├── PreferencesView.h/.cpp   # auto-scan, interval, theme settings
│   │   └── CakeLieView.h/.cpp       # "the cake is a lie" invalid-project banner
│   └── theme/
│       └── ThemeLoader.h/.cpp        # ~/.config/cake/themes/*.xml → ValueTree
├── tests/
│   └── fixtures/                     # ValueTree snapshots for state permutations
├── SPEC.md                           # NEW — C++/jam specification
├── RFC.md                            # this document
├── CLAUDE.md -> ~/.carol/CAROL.md
└── carol/
```

### 4.2 CMakeLists.txt (from TIT pattern)

```cmake
cmake_minimum_required(VERSION 4.2.0)

set(JAM_ROOT "$ENV{HOME}/Documents/Poems/dev/jam")
list(APPEND CMAKE_MODULE_PATH "${JAM_ROOT}/cmake")
include(BuildSetup)

if(APPLE)
    set(_EXTRA_LANGS OBJC OBJCXX)
else()
    set(_EXTRA_LANGS "")
endif()

project(cake-cpp
    VERSION 0.0.0
    LANGUAGES C CXX ${_EXTRA_LANGS}
    DESCRIPTION "CAKE-cpp: CMake Project Manager TUI"
)

configure_app(
    TARGET_NAME      cakec
    PRODUCT_NAME     cakec
    VERSION          "${PROJECT_VERSION}"
    COMPANY          jrengmusic
    COMPANY_WEBSITE  "https://jrengmusic.com"
    BUNDLE_ID        "com.jrengmusic.cakec"
    MODULES
        juce_core
        juce_events
        juce_data_structures
        juce_graphics
        juce_gui_basics
    JAM_MODULES
        jam_core
        jam_data_structures
        jam_markdown
        jam_tui
        jam_subprocess
    BINARY_FILES
        "${CMAKE_CURRENT_SOURCE_DIR}/Assets/cake-logo.svg"
        "${CMAKE_CURRENT_SOURCE_DIR}/Assets/cake-lie.svg"
    BINARY_NAMESPACE "BinaryData"
    EXTRA_DEFINES
        JUCE_WEB_BROWSER=0
        JUCE_USE_CURL=0
)
```

### 4.3 ValueTree schema (CakeState root)

```
CAKE
├── PROJECT                          # domain state
│   ├── workingDirectory     (string)
│   ├── hasCMakeLists        (bool)
│   ├── selectedGenerator    (string: Xcode | Ninja | VS2026 | VS2022)
│   ├── configuration        (string: Debug | Release)
│   └── GENERATORS[]                 # detected generators
│       └── name             (string)
│           isIDE            (bool)
├── BUILDS                           # scanned build directories
│   └── BUILD[]
│       ├── generator        (string)
│       ├── path             (string)
│       ├── exists           (bool)
│       └── isConfigured     (bool)
├── MENU                             # 8 fixed rows, selectability computed
│   └── ITEM[]
│       ├── id               (string)
│       ├── label            (string)
│       ├── value            (string)
│       ├── shortcut         (string)
│       ├── isSelectable     (bool)
│       ├── isAction         (bool)
│       └── hint             (string)
├── CONSOLE                          # streaming cmake output
│   └── LINE[] (text, stream: stdout|stderr|info)
├── SELECTION                        # UI-layer transient state
│   ├── menuIndex            (int)
│   └── mode                 (string: Menu | Preferences | Console | InvalidProject)
├── ASYNC                            # operation state
│   ├── isActive             (bool)
│   ├── isAborted            (bool)
│   ├── currentOp            (string: None | Build | Generate | Clean | CleanAll | Regenerate)
│   └── spinnerFrame         (int)
├── CONFIG                           # persisted preferences
│   ├── autoScanEnabled      (bool)
│   ├── autoScanInterval     (int, minutes)
│   ├── theme                (string)
│   ├── lastGenerator        (string)
│   └── lastConfiguration    (string)
└── THEME                            # loaded from ~/.config/cake/themes/*.xml
    └── ...                          # all theme color properties
```

Menu regeneration is a ValueTree::Listener on `PROJECT` + `BUILDS` subtrees. No explicit `rebuildMenu()` calls.

### 4.4 Module consumption (zero forging)

| Module | Action | Notes |
|---|---|---|
| `jam_core` | Consume from jam | File::Watcher for theme hot-reload, identifiers, utilities |
| `jam_data_structures` | Consume from jam | ValueTree wrapper (transitive dep) |
| `jam_markdown` | Consume from jam | Transitive dep of jam_tui (unused by CAKE) |
| `jam_tui` | Consume from jam | Menu, ConsoleStream, SplitPane, Spinner, ThemeResolver, braille |
| `jam_subprocess` | Consume from jam | cmake execution, streaming stdout/stderr |

No module needs implementation. CAKE is a pure consumer of the existing jam ecosystem.

### 4.5 Threading model (mirror TIT/END)

| Thread | Owns | Crossing |
|---|---|---|
| Message (JUCE main) | `CakeState` ValueTree, `CakeScreen`, all views | — |
| Subprocess worker (`jam_subprocess`) | `juce::ChildProcess` (cmake) | atomic writes to `CakeState` + `callAsync` |
| Timer (`juce::Timer`) | `CakeState::flush()`, auto-scan tick | message thread (inherited) |
| File watcher (`jam::File::Watcher`) | theme directory monitoring | `callAsync` → message thread |

Zero locks on hot path. `callAsync` only crossing primitive. Identical to TIT/END.

### 4.6 Phase sequence

| Phase | Deliverable | Effort |
|---|---|---|
| **0** | Scaffold — `___legacy___/` archive, CMake, `Source/Main.cpp` + `CakeApp` stub builds green | <1 day |
| **1** | `CakeState` + `CakeIdentifier` + `CakeAxis` + ValueTree schema + fixture framework | <1 day |
| **2** | Views — MenuView (8 fixed rows), Header, Footer, BannerView, CakeLieView, PreferencesView, ConsoleView | 1-2 days |
| **3** | Ops — `CmakeRunner` + `GeneratorDetector` + Generate/Build/Clean/Open operations | 1 day |
| **4** | Integration — wire ops to state, auto-scan, config persistence, theme hot-reload | 1 day |
| **5** | macOS release — codesign, notarize, install target | <1 day |
| **Total MVP** | Feature parity with Go CAKE on macOS | **3-5 days CAROL walltime** |

**Windows MSYS2 parity post-MVP** — estimated +1-2 days (MSVC environment capture via `jam_subprocess`, `vswhere` detection, path normalization). Simpler than TIT's Windows sprint because CAKE has no ConPTY byte handling or process-tree kill requirements — cmake is a single subprocess.

### 4.7 CAKE-specific patterns (not in TIT)

**Fixed menu with selectability model:** Go CAKE has 8 fixed rows (Project, Regenerate, Open, Separator, Configuration, Build, Clean, CleanAll). All are always visible. Unavailable items are dimmed and not navigable. This is simpler than TIT's conditional menu — `jam_tui::Menu` already supports `isSelectable` per row.

**Dual banner:** Main banner (cake-logo.svg) shown in 35% right pane during normal operation. "Cake is a lie" banner (cake-lie.svg) shown full-screen when no CMakeLists.txt found. Both rendered via `jam_tui/braille/`.

**Generator cycling:** Project row cycles through detected generators on Enter/Space. Pure toggle — no submenu, no picker. Cycle order: Xcode → Ninja → (loop), or VS → Ninja → (loop) on Windows.

**Configuration toggle:** Debug ↔ Release, bidirectional.

**Console mode:** cmake output streams in real-time. Spinner in header during active operation (reuses TIT Sprint 7 pattern — braille spinner + op label). ESC aborts via context cancellation. Same `jam_subprocess` streaming pattern as TIT.

---

## 5. BLESSED Compliance Checklist

- [x] **Bounds** — `CakeApp` owns `CakeState` + `CakeScreen` + `CmakeRunner` via `std::unique_ptr`. `CmakeRunner` owns subprocess via `jam_subprocess`. Timer owned by `CakeState`. FileWatcher owned by ThemeLoader. Thread ownership: message thread owns VT+views, subprocess worker owns ChildProcess.
- [x] **Lean** — 300/30/3 per file. 8 fixed menu rows generated from data, not switch chain. Generator detection is lookup, not conditional ladder. YAGNI: no cmake-project-file parser, no build-system abstraction layer, no plugin system.
- [x] **Explicit** — Zero early returns. All parameters visible. Magic values → named constants in `CakeIdentifier.h`. `jassert` on invariants. No silent fails — every cmake error streams to ConsoleView.
- [x] **Single Source of Truth** — `CakeState` ValueTree is SSOT for all application state. `CakeIdentifier.h` is SSOT for all VT keys. Theme XML is SSOT for colors. Generator constants are SSOT (single definition site).
- [x] **Stateless** — View components hold transient render state only (scroll offset, focus). All persistent state in `CakeState`. Orchestrator tells, never asks.
- [x] **Encapsulation** — `jam_tui` imports no CAKE application header. `Source/ops/` imports no `Source/view/`. Unidirectional layer flow. Views observe VT via listeners — never poke ops layer.
- [x] **Deterministic** — Same ProjectState VT → same menu row selectability. Same cmake command + same cwd → same subprocess output. Emergent from BLESSE.

---

## 6. Open Questions

None. All decisions resolved:

- Project location: same dir, Go to `___legacy___/` (ARCHITECT 2026-04-25)
- Binary name: `cakec` during port (ARCHITECT 2026-04-25)
- Go legacy: continues shipping from `___legacy___/` (ARCHITECT 2026-04-25)
- SPEC: new SPEC for jam (ARCHITECT 2026-04-25)
- All patterns: TIT-cpp established patterns apply (ARCHITECT 2026-04-25)
- Font restructuring: not a concern for CAKE port (ARCHITECT 2026-04-25)
- Config: XML → ValueTree (TIT pattern)
- Theme: external XML + hot-reload (TIT pattern)
- State: APVTS-mirror (TIT pattern)
- Windows: macOS-first, deferred sprint (TIT pattern)
- Modules: all already implemented, zero forging (verified 2026-04-25)

---

## 7. Handoff Notes

### 7.1 SPEC.md — COUNSELOR writes new

ARCHITECT directed: write new SPEC.md for jam (not carry Go SPEC). Go SPEC.md at `___legacy___/SPEC.md` is the requirements source. COUNSELOR translates to C++/JUCE/jam terms — ValueTree state model, `jam_tui` primitives, `jam_subprocess` ops, XML config/themes. Feature set unchanged.

### 7.2 No module forging

Unlike TIT-cpp which forged 8 `jam_tui` primitives + `jam_subprocess` + braille renderer, CAKE-cpp consumes everything from jam as-is. This is the validation: if CAKE-cpp can port Go CAKE's full feature set using only existing `jam_tui` primitives without modification, the framework is general enough.

If a primitive needs adjustment, that surfaces as a jam PR — not a CAKE-local fork.

### 7.3 Go CAKE ARCHITECTURE.md — reference only

Go CAKE's `ARCHITECTURE.md` (786 lines) documents the Bubble Tea / Elm architecture. Useful as layer-mapping reference. Not carried forward — COUNSELOR writes new ARCHITECTURE.md for C++ port after implementation.

### 7.4 SPRINT-LOG context

Go CAKE has 7 sprints logged in `carol/SPRINT-LOG.md`. Sprint 7 (braille spinner) and Sprint 6 (Ninja via VS env + process-tree abort) are directly relevant — their patterns already exist in `jam_tui` and `jam_subprocess`. COUNSELOR should read Sprint 6-7 for MSVC environment capture and console streaming patterns that need C++ equivalents.

### 7.5 Asset migration

Two SVG files move from `internal/ui/assets/` to `Assets/` at project root:
- `cake-logo.svg` — main banner (35% right pane)
- `cake-lie.svg` — "the cake is a lie" banner (invalid project mode)

Both embedded via `juce_add_binary_data` in CMakeLists.txt, consumed by `jam_tui/braille/` at runtime.

### 7.6 Reference document precedence

1. This RFC — port scope, architecture, phase sequence
2. Go `SPEC.md` (in `___legacy___/`) — requirements source for new SPEC
3. Go `ARCHITECTURE.md` (in `___legacy___/`) — layer mapping reference
4. TIT-cpp `PLAN-tit-cpp-port.md` — established patterns reference
5. TIT-cpp `CMakeLists.txt` — `configure_app()` template

### 7.7 Estimate confidence

**3-5 day MVP** is defensible. Confidence rationale:

- Zero architectural invention
- Zero module forging — all jam modules already implemented
- Full reference implementation in Go (~7,700 LOC)
- All TIT-cpp patterns established and proven
- Dramatically simpler app logic than TIT (no FSMs, no parsers, no multi-axis state)
- `configure_app()` CMake infrastructure already handles all build boilerplate

**Estimate floor** (3 days): sustained CAROL cadence, no primitive adjustments needed, macOS only.

**Estimate ceiling** (5 days): one jam_tui primitive needs minor adjustment, MSVC detection requires iteration on first build.

### 7.8 Sprint ownership (suggested)

- **Sprint 1:** Phase 0 + Phase 1 (scaffold + state layer — builds green, fixtures render)
- **Sprint 2:** Phase 2 + Phase 3 (views + ops — functional CAKE with cmake execution)
- **Sprint 3:** Phase 4 + Phase 5 (integration + macOS release)

Each sprint ends with `log sprint` per protocol.

### 7.9 Post-MVP queue

- Windows MSYS2 parity sprint (+1-2 days)
- Go `cake` binary retirement (ARCHITECT call after Windows parity)
- Theme sharing across TIT-cpp and CAKE-cpp (common `jam_tui` theme schema)

---

*RFC complete. Status: Ready for COUNSELOR handoff.*

**Rock 'n Roll!**
**JRENG!**
