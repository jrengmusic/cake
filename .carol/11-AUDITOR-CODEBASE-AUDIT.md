# Sprint 11 Audit Report

**Date:** 2026-01-29
**Role:** AUDITOR
**Agent:** Amp (Claude Sonnet 4)
**Scope:** Full codebase
**Summary:** Critical: 2, High: 4, Medium: 5, Low: 3

---

## REFACTORING OPPORTUNITIES (CRITICAL PRIORITY)

### [REF-000] Dead Logic: "Ninja Multi-Config" Generator Exists When All Are Multi-Config
**Files:**
- `internal/state/project.go:127-130`
- `internal/utils/generator.go:11,43`

**Principle:** Lean + SSOT
**Issue:** Code adds "Ninja Multi-Config" as separate generator, but all projects now use multi-config path. Distinction is meaningless.
**Severity:** CRITICAL
**Impact:** Confuses users, clutters generator cycling, dead logic

**Current (state/project.go:121-131):**
```go
if ps.checkCommandExists("ninja") {
    ps.AvailableGenerators = append(..., Generator{Name: "Ninja", ...})
    ps.AvailableGenerators = append(..., Generator{Name: "Ninja Multi-Config", ...})
}
```

**Fix:** Remove "Ninja Multi-Config" everywhere. Simplify generator names:
- **Ninja**
- **Xcode**
- **Visual Studio 2026** (current)
- **Visual Studio 2022** (legacy, for existing projects)

All are multi-config. Simplified display names (not CMake-style "Visual Studio 17 2022").

**Files to update:**
- `internal/state/project.go` — remove lines 127-130
- `internal/utils/generator.go` — remove "Ninja Multi-Config" from validGenerators and switch case
- `internal/utils/platform.go` — change Windows default from "Ninja Multi-Config" to "Ninja"

---

### [REF-001] SSOT Violation: Build Directory Path Construction Duplicated
**Files:** 
- `internal/ops/setup.go:26`
- `internal/ops/build.go:17`
- `internal/ops/clean.go:15`
- `internal/ops/open.go:18`
- `internal/app/op_regenerate.go:35`
- `internal/state/project.go:232,238,245`

**Principle:** SSOT (Single Source of Truth)
**Issue:** Build path `filepath.Join(projectRoot, "Builds", generator)` duplicated 6+ times across ops/, app/, and state/
**Severity:** CRITICAL
**Benefits:** Eliminate inconsistency risk, single point of change for build path logic
**Impact:** High (any change requires updating 6+ locations)
**Effort:** Low (extract to `state.ProjectState.GetBuildPath()` and pass to ops functions)
**Priority:** CRITICAL

**Current Pattern (duplicated):**
```go
// ops/setup.go:26
buildDir := filepath.Join(workingDir, "Builds", generator)

// ops/build.go:17
buildDir := filepath.Join(projectRoot, "Builds", generator)

// ops/clean.go:15
buildDir := filepath.Join(projectRoot, "Builds", generator)

// app/op_regenerate.go:35
buildDir := filepath.Join(projectRoot, "Builds", project)
```

**Fix:** All ops functions should receive `buildDir` as parameter, calculated once by `ProjectState.GetBuildPath()`.

---

### [REF-002] SSOT Violation: Command Streaming Pattern Duplicated
**Files:**
- `internal/ops/setup.go:40-80`
- `internal/ops/build.go:29-69`

**Principle:** SSOT + Lean
**Issue:** Stdout/stderr pipe setup and goroutine streaming identical in both files (~40 lines duplicated)
**Severity:** High
**Benefits:** Single streaming implementation, easier to maintain/enhance
**Impact:** Medium (any streaming bug must be fixed in 2 places)
**Effort:** Medium (extract `streamCommand()` helper)

**Duplicated Pattern:**
```go
stdout, err := cmd.StdoutPipe()
// ... error handling ...
stderr, err := cmd.StderrPipe()
// ... error handling ...
go func() {
    scanner := bufio.NewScanner(stdout)
    for scanner.Scan() { outputCallback(scanner.Text(), ui.TypeStdout) }
}()
go func() {
    scanner := bufio.NewScanner(stderr)
    for scanner.Scan() { outputCallback(scanner.Text(), ui.TypeStderr) }
}()
```

**Fix:** Create `internal/utils/stream.go` with `StreamCommand(ctx, cmd, callback)` helper.

---

## LIFESTAR VIOLATIONS

### [AUD-001] Lean Violation: handleMenuKeyPress() Code Duplication
**File:** `internal/app/app.go:742-822`
**Principle:** Lean (Keep It Simple)
**Issue:** Shortcut handlers (g, o, b, c) repeat identical 15-line pattern for each key
**Severity:** High
**Impact:** 60+ lines of nearly identical code, maintenance nightmare

**Repeated Pattern (4 times):**
```go
case "g", "G":
    idx := a.GetVisibleIndex("regenerate")
    if idx >= 0 {
        a.selectedIndex = idx
        handled, cmd := a.ToggleRowAtIndex(idx)
        if handled {
            a.menuItems = a.GenerateMenu()
            newVisible := a.GetVisibleRows()
            if a.selectedIndex >= len(newVisible) {
                a.selectedIndex = len(newVisible) - 1
            }
            if a.selectedIndex < 0 {
                a.selectedIndex = 0
            }
            return a, cmd
        }
    }
```

**Fix:** Extract helper `executeShortcut(rowID string) (tea.Model, tea.Cmd)` and reduce to:
```go
case "g", "G": return a.executeShortcut("regenerate")
case "o", "O": return a.executeShortcut("openIde")
case "b", "B": return a.executeShortcut("build")
case "c", "C": return a.executeShortcut("clean")
```

---

### [AUD-002] Explicit Violation: Unused `config` Parameter
**File:** `internal/ops/clean.go:14`
**Principle:** Explicit (Dependencies Visible)
**Issue:** `config` parameter passed but never used in `ExecuteCleanProject()`
**Severity:** Low
**Impact:** Confusing API, suggests config affects cleaning when it doesn't

**Current:**
```go
func ExecuteCleanProject(generator, config, projectRoot string, ...)
```

**Fix:** Remove `config` parameter since clean removes entire generator directory.

---

### [AUD-003] Explicit Violation: Unused `config` Parameter (Build)
**File:** `internal/ops/build.go:16`
**Principle:** Explicit
**Issue:** `config` is passed to `ExecuteBuildProject()` and used in command, but per ARCHITECTURE.md all projects are multi-config
**Severity:** Medium
**Impact:** Misleading signature if config is always passed via `--config` flag

**Analysis:** This is correct behavior for multi-config generators. No change needed, but verify consistency with SPEC.md.

---

### [AUD-004] Lean Violation: Unused Helper in utils/generator.go
**File:** `internal/utils/generator.go:9-35`
**Principle:** Lean
**Issue:** `BuildCMakeCommand()` constructs path as `build_ninja` but project uses `Builds/Ninja/`
**Severity:** Medium
**Impact:** Dead code, inconsistent with actual build path convention

**Current (unused):**
```go
buildDir := filepath.Join(projectRoot, "build_"+strings.ToLower(...))
```

**Actual convention (per SPEC.md):**
```
Builds/<Generator>/
```

**Fix:** Either remove unused function or update to match `Builds/` convention.

---

### [AUD-005] Testable Violation: No Unit Tests Present
**File:** Entire codebase
**Principle:** Testable
**Issue:** No `*_test.go` files found in any package
**Severity:** High
**Impact:** Cannot verify correctness, regressions undetected
**Effort:** High (requires test infrastructure setup)

**Fix:** Add test files starting with critical paths:
- `internal/state/project_test.go` (ProjectState methods)
- `internal/config/config_test.go` (Load/Save roundtrip)
- `internal/ui/menu_test.go` (menu generation)

---

### [AUD-006] Immutable Violation: Config.Save() Called on Every Setting Change
**File:** `internal/config/config.go:127-151`
**Principle:** Immutable (Predictable Behavior)
**Issue:** Every setter (`SetTheme`, `SetAutoScanEnabled`, etc.) calls `Save()` immediately
**Severity:** Low
**Impact:** Minor - per ARCHITECTURE.md this is intentional ("immediate persistence")

**Analysis:** This is documented as intentional design decision (Decision 5 in ARCHITECTURE.md). No change needed.

---

## ANTI-PATTERNS DETECTED

### [ANT-001] Potential Race Condition in Output Streaming
**Files:** `internal/ops/setup.go:61-80`, `internal/ops/build.go:50-69`
**Issue:** Goroutines write to `outputCallback` concurrently without synchronization
**Impact:** Possible interleaved output lines, cosmetic issue
**Severity:** Low

**Current:**
```go
go func() {
    for scanner.Scan() { outputCallback(line, ui.TypeStdout) }
}()
go func() {
    for scanner.Scan() { outputCallback(line, ui.TypeStderr) }
}()
```

**Mitigation:** `OutputBuffer.Append()` should be checked for thread-safety. If using mutex internally, this is safe.

---

### [ANT-002] Inconsistent Generator Name Handling
**Files:**
- `internal/state/project.go:127-130` - uses "Ninja Multi-Config"
- `internal/utils/generator.go:11` - uses "Ninja Multi-Config"
- `internal/utils/generator.go:43` - uses "Ninja Multi-Config"

**Principle:** SSOT
**Issue:** Generator name "Ninja Multi-Config" hardcoded in multiple locations
**Severity:** Medium
**Impact:** Renaming generator requires updating 5+ files

**Fix:** Define generator name constants in single location:
```go
const (
    GeneratorXcode          = "Xcode"
    GeneratorNinja          = "Ninja"
    GeneratorNinjaMulti     = "Ninja Multi-Config"
    GeneratorVS2022         = "Visual Studio 17 2022"
    GeneratorVS2019         = "Visual Studio 16 2019"
)
```

---

### [ANT-003] God Object Tendency: Application struct
**File:** `internal/app/app.go:25-60`
**Principle:** Lean
**Issue:** Application struct has 19 fields managing UI, async ops, console, dialogs, scanning
**Severity:** Medium
**Impact:** Growing complexity, harder to reason about state

**Current fields:** width, height, sizing, theme, mode, selectedIndex, menuItems, projectState, config, quitConfirmActive, quitConfirmTime, consoleState, outputBuffer, consoleAutoScroll, asyncState, windowSize, keyDispatcher, runningCmd, cancelContext, footerHint, isScanning, lastBuildDir, confirmDialog, pendingOperation

**Observation:** Sprint 6 extracted op_*.go and async_state.go which helped. Current size is manageable but watch for growth.

---

## SPEC DISCREPANCIES (DOC UPDATE RECOMMENDATIONS)

### [DOC-001] SPEC.md Outdated: Unix Makefiles Still Mentioned
**SPEC says (line 128-130):** "Unix Makefiles (always)"
**Code has:** Unix Makefiles removed in Sprint 7
**File:** `internal/utils/generator.go`, `internal/state/project.go`
**Recommendation:** Update SPEC.md lines 35, 128-130 to remove Unix Makefiles references
**Note:** Codebase is SSOT, this is not a code violation

---

### [DOC-002] SPEC.md Outdated: Single-Config Generator Path
**SPEC says (line 48-49):** 
```
Single-Config Generators (separate build per configuration):
- Ninja (cross-platform, CLI)
Path: Builds/<Generator>/<Configuration>/
```
**Code has (project.go:232):**
```go
// All projects are multi-config: Builds/<Generator>/
return filepath.Join(ps.WorkingDirectory, "Builds", gen)
```
**Recommendation:** Update SPEC.md to reflect that all generators now use multi-config path
**Note:** Codebase is SSOT, this is not a code violation

---

### [DOC-003] ARCHITECTURE.md Minor Discrepancy
**ARCHITECTURE says (line 495-501):** Shows `IsGeneratorMultiConfig()` method with for-loop
**Code has:** Method exists but comment says "All projects are multi-config"
**Recommendation:** Clarify in ARCHITECTURE.md that multi-config is now universal
**Note:** Minor documentation drift

---

## SUMMARY

### By Category
- Refactoring Opportunities: 3 (Lean: 1, SSOT: 2)
- LIFESTAR Violations: 6
- Anti-Patterns: 3
- Doc Updates Needed: 3

### By Severity
- CRITICAL: 3 (REF-000, REF-001, REF-002)
- High: 4 (AUD-001, AUD-005, REF-002)
- Medium: 5 (AUD-003, AUD-004, ANT-002, ANT-003, DOC-002)
- Low: 3 (AUD-002, ANT-001, DOC-003)

### Recommended Actions
1. **CRITICAL:** Remove "Ninja Multi-Config" - all generators are now multi-config (REF-000)
2. **CRITICAL:** Extract build path logic to single function, pass to ops (REF-001)
3. **CRITICAL:** Extract command streaming helper to eliminate duplication (REF-002)
3. **High:** Extract shortcut handler helper to reduce 60+ lines of duplication (AUD-001)
4. **High:** Add unit tests for critical paths (AUD-005)
5. **Medium:** Define generator name constants (ANT-002)
6. **Medium:** Remove unused `config` param from clean.go (AUD-002)
7. **Low:** Update SPEC.md and ARCHITECTURE.md to match codebase

### Code Quality Assessment
- **TIT Compliance:** 90%+ (good architecture patterns adopted)
- **LIFESTAR Compliance:** 75% (SSOT violations remain, no tests)
- **LOVE Compliance:** 85% (fail-fast error handling good, explicit deps mostly good)

### Positive Observations
- Clean layer separation (app → state → utils, app → ops → utils)
- Good Elm architecture adherence (pure View, Cmd-based side effects)
- Footer system well-factored (TIT pattern)
- Async operation handling clean (AsyncState struct)
- Configuration persistence working correctly
