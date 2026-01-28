# cake → Go + Bubble Tea + Lip Gloss: Port Analysis Summary

**Status:** ✅ Feasibility Analysis Complete (CORRECTED v1.1)  
**Complexity:** 🟡 Medium (2.5-3 weeks, ~16-20 days effort)  
**Recommended:** ✅ Proceed with porting

---

## Critical Corrections (v1.1)

1. **JUCE is NOT for audio plugins** - Used only for cross-platform utilities (filesystem, process management). cake is a pure **CMake management tool**, not plugin-specific.

2. **Console output is IDENTICAL to TIT** - Reuse `internal/ui/console.go` directly. Real-time line streaming, ConsoleOutState, OutputBuffer, RenderConsoleOutput().

3. **Async operations required** - Operations (Setup, Build, Clean) run in background worker threads using Bubble Tea's Cmd pattern (identical to TIT's git operations).

4. **No quit menu** - Ctrl+C is handled at application level (exactly like TIT), with 2-second confirmation timeout. No "Exit" menu item.

---

## What is cake?

A **pure CMake management tool** for any CMake project (not limited to JUCE). It provides:
- Interactive terminal UI for CMake configuration
- Build execution with platform-specific generators (Xcode, Ninja, Visual Studio)
- Build cleanup
- Cross-platform support (macOS, Windows, Linux)

**JUCE is only used for:**
- Cross-platform filesystem utilities
- Process/child process management
- Utility functions

**Key: It's NOT a git tool and NOT plugin-specific.**

---

## Why Port It?

1. **Simplicity:** 3,569 lines C++ → ~3,000-4,000 lines Go (similar size)
2. **Complements TIT:** TIT manages git, cake manages builds
3. **Reusable components:** TIT's banner, menu, UI system already solved
4. **Lower complexity:** No git state tracking, no history management, no async mutations
5. **Natural integration:** Both terminal UIs, same tech stack (Bubble Tea + Lip Gloss)

---

## Complexity Breakdown

### Easy (Reuse from TIT)
- ✅ Menu system (navigation, selection, shortcuts)
- ✅ Banner rendering (braille SVG display)
- ✅ Output scrolling (build logs)
- ✅ Footer hints (keyboard legend)
- ✅ Confirmation dialogs (Ctrl+C quit)

### Medium (Adapt from TIT)
- 🟡 Project state detection (cake: find build dirs; TIT: track git)
- 🟡 Operation execution (cake: CMake/build; TIT: git commands)
- 🟡 Cross-platform command handling
- 🟡 Real-time output streaming

### Hard (New)
- 🟠 Platform-specific generator logic (Xcode args, VS paths, etc.)
- 🟠 CMake path detection and validation
- 🟠 IDE directory discovery (Ninja vs Xcode vs Visual Studio)

---

## Key Architectural Decisions

### 1. **Async Operations Pattern (IDENTICAL to TIT)**

**Critical:** cake's build, setup, and clean operations MUST run in background worker threads like TIT's git operations:

```go
// Set flag + start worker
a.asyncOperationActive = true
return a, a.cmdBuildProject(buildDir)

// In worker (Bubble Tea Cmd):
func (a *Application) cmdBuildProject(buildDir string) tea.Cmd {
    return func() tea.Msg {
        // WORKER THREAD - runs cmake --build
        // Calls callback for each output line
        result := executeBuildProject(buildDir, outputCallback)
        return BuildCompleteMsg{Success: result.Success}
    }
}

// Handle completion:
case BuildCompleteMsg:
    a.asyncOperationActive = false
    // Return to menu
```

This prevents UI freezing and enables smooth console scrolling.

### 2. **Console Output (EXACT Reuse from TIT)**

**NOT using history pane pattern.** Use console output pattern:
- `ConsoleOutState` - scroll tracking
- `OutputBuffer` - thread-safe line buffer
- `RenderConsoleOutput()` - real-time rendering with wrapping
- Real-time line streaming from operation callbacks

### 3. **Banner Rendering (DIRECT SVG Reuse)**

TIT reads SVG files directly using Go's `embed` package:

```go
//go:embed assets/cake-logo.svg
var logoFS embed.FS

// In RenderBannerDynamic():
logoData, _ := logoFS.ReadFile("assets/cake-logo.svg")
svgString := string(logoData)

// Pass to banner rendering (no modifications)
brailleArray := banner.SvgToBrailleArray(svgString, width, height)
```

**For cake:**
- Embed cake's existing SVG (`Source/banner/cake.svg`)
- Pass it to `SvgToBrailleArray()` (already in TIT package)
- No code changes to `braille.go` or `svg.go` needed
- Just use the functions as-is with cake's SVG

### 4. **Ctrl+C Handling (Application Level)**

**NO menu item for "Exit":**
- Ctrl+C during menu → show "Press Ctrl+C again to confirm"
- Ctrl+C during operation → show "operation in progress"
- Second Ctrl+C or 2-second timeout → quit
- Identical to TIT's logic

---

## Reusable from TIT (50-60% of code)

| Component | Reuse Type | How It Works |
|-----------|-----------|-------------|
| Console output (`console.go`) | **EXACT REUSE** | Copy entire file, no changes. Stream CMake output same as git operations. |
| Menu system (`menu.go`) | **EXACT REUSE** | Copy entire file. Main menu + generator selection, same pattern as git menus. |
| Braille SVG (`braille.go` + `svg.go`) | **DIRECT USE** | Copy both files. Call `SvgToBrailleArray(cakeSvgString)` with cake's SVG - no modifications. |
| Ctrl+C logic (`handlers.go`) | **COPY PATTERN** | Copy quit confirmation logic (2-sec timeout). Handle in app-level Update(). |
| Msg types (`messages.go`) | **ADAPT** | Copy pattern, create SetupCompleteMsg, BuildCompleteMsg, CleanCompleteMsg. |
| Layout (`layout.go`) | **COPY + ADAPT** | Copy embedding pattern, create cake-specific RenderBannerDynamic(). |
| Status bar (`statusbar.go`) | **EXACT REUSE** | Copy entire file. Same footer hints + keyboard legend. |
| Theme/colors (`theme.go`) | **REUSE** | Copy color definitions. Same palette for both tools. |
| Info rows (`inforow.go`) | **OPTIONAL** | Nice-to-have for project info display. |

---

## Effort Estimation

| Phase | Task | Complexity | Time | Owner |
|-------|------|-----------|------|-------|
| 1 | Core app state + mode (copy from TIT pattern) | 🟢 LOW | 2 days | SCAFFOLDER |
| 2 | Menu system (reuse from TIT) | 🟢 LOW | 1 day | SCAFFOLDER |
| 3 | Console output (exact reuse from TIT) | 🟢 LOW | 1 day | SCAFFOLDER |
| 4 | Project detection logic | 🟡 MEDIUM | 2 days | CARETAKER |
| 5 | Async operation framework (copy TIT pattern) | 🟡 MEDIUM | 2-3 days | SCAFFOLDER |
| 6 | Setup operation (CMake execution) | 🟡 MEDIUM | 2-3 days | CARETAKER |
| 7 | Build operation (cmake --build) | 🟡 MEDIUM | 2-3 days | CARETAKER |
| 8 | Clean operation (rm -rf) | 🟡 MEDIUM | 1-2 days | CARETAKER |
| 9 | Cross-platform generator logic | 🟡 MEDIUM | 2-3 days | CARETAKER |
| 10 | Ctrl+C handling (copy from TIT) | 🟢 LOW | 1 day | SCAFFOLDER |
| 11 | Banner/braille (adapt from TIT) | 🟢 LOW | 1 day | SCAFFOLDER |
| 12 | Testing all platforms | 🟡 MEDIUM | 2-3 days | User |

**Total:** 16-20 days (~2.5-3 weeks) for complete MVP

---

## Key Risks & Mitigations

| Risk | Impact | Solution |
|------|--------|----------|
| **CMake detection fails** | App won't start | Check PATH, error if missing, document requirement |
| **Generator quoting wrong** | Build command syntax fails | Test all platforms exhaustively before release |
| **Output encoding issues** | Garbled logs on Windows | Normalize UTF-8, strip ANSI codes |
| **Process hanging** | Ctrl+C doesn't work | Use context.Context with timeout |
| **Cross-platform paths** | Hardcoded paths break on other OS | Use `filepath` package, test on all platforms |

---

## Critical Requirements (Per Your Feedback)

⚠️ **Watch out for CARETAKER refactoring:**

Code may change during polishing. When adding error handling:
- ✅ Use explicit error returns (never `_`)
- ✅ Wrap errors with context (`fmt.Errorf(...%w...", err)`)
- ✅ Validate all user inputs before operation
- ✅ Log exact failure reasons

**FAIL-FAST Rule** (from SESSION-LOG.md):
- ❌ NO silent failures
- ❌ NO empty fallback strings
- ❌ NO swallowed stderr
- ✅ Return explicit errors
- ✅ User sees error immediately

---

## File Structure

```
cake-go/
├── cmd/cake/main.go                    # Entry point
├── internal/
│   ├── app/
│   │   ├── app.go                      # Application + Update/View
│   │   ├── handlers.go                 # Ctrl+C, menu, async commands
│   │   ├── messages.go                 # SetupCompleteMsg, etc.
│   │   └── theme.go                    # (COPY FROM TIT)
│   ├── state/
│   │   └── project.go                  # ProjectState detection
│   ├── ops/
│   │   ├── setup.go                    # executeSetupProject()
│   │   ├── build.go                    # executeBuildProject()
│   │   └── clean.go                    # executeCleanProject()
│   ├── utils/
│   │   ├── project.go, exec.go, generator.go, platform.go
│   ├── ui/                             # (COPY FROM TIT - NO CHANGES)
│   │   ├── console.go, menu.go, statusbar.go, etc.
│   └── banner/                         # (COPY FROM TIT)
│       └── render.go (adapt SVG source)
├── go.mod, go.sum, build.sh, README.md
```

---

## Next Steps

1. **Write SPEC.md** ✅ DONE (saved to cake repo)

2. **SCAFFOLDER creates literal scaffold**
   - Follow SPEC exactly (no improvements yet)
   - Create file structure, function stubs, imports
   - Make it compile with TODO markers

3. **CARETAKER polishes implementation**
   - Add error handling (fail-fast rule)
   - Add platform detection
   - Wire components together
   - Test on macOS first

4. **User tests on all platforms**
   - macOS (Xcode + Ninja)
   - Windows (Visual Studio + Ninja Multi-Config)
   - Linux (Ninja + Unix Makefiles)

5. **SURGEON fixes bugs** (if discovered)
   - Process execution issues
   - Output streaming problems
   - Platform-specific failures

---

## Success Metrics

✅ **Functional:**
- Main menu displays 3 operations
- Setup generates CMake build directory
- Build executes with live output
- Clean removes artifacts
- Ctrl+C works
- All keyboard shortcuts work

✅ **Code Quality:**
- No silent failures
- Reusable components documented
- Cross-platform compatibility

✅ **Performance:**
- Menu render < 100ms
- Smooth output scrolling
- No memory leaks

---

## Final Verdict

**✅ PORT IS FEASIBLE AND HIGHLY RECOMMENDED**

cake is simpler than TIT and can reuse 50-60% of code directly:
1. Console output display (exact copy)
2. Async operation framework (exact copy of pattern)
3. Ctrl+C handling (exact copy)
4. Menu system (with minimal changes)
5. UI components (layouts, theming)

**The main work is platform-specific logic:**
- Generator detection and argument building
- Build command variants (Xcode vs Ninja vs Visual Studio)
- IDE directory discovery

**Estimated ROI:** 2.5-3 weeks of engineering for a production-grade CLI tool that complements TIT perfectly.

---

**JRENG! Analysis complete. Ready for SCAFFOLDER to begin literal scaffold phase.**
