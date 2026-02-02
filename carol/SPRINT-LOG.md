# SPRINT-LOG.md Template

**Project:** cake  
**Repository:** /Users/jreng/Documents/Poems/dev/cake  
**Started:** 2026-01-28

**Purpose:** Track agent registrations, sprint work, and completion reports. This file is mutable and rotates old entries (keeps last 5 sprints).

---

## 📖 Notation Reference

**[N]** = Sprint Number (e.g., `1`, `2`, `3`...)

**File Naming Convention:**
- `[N]-[ROLE]-[OBJECTIVE].md` — Task summary files written by agents
- `[N]-COUNSELOR-[OBJECTIVE]-KICKOFF.md` — Phase kickoff plans (COUNSELOR)
- `[N]-AUDITOR-[OBJECTIVE]-AUDIT.md` — Audit reports (AUDITOR)

**Example Filenames:**
- `1-COUNSELOR-INITIAL-PLANNING-KICKOFF.md` — COUNSELOR's plan for sprint 1
- `1-ENGINEER-MODULE-SCAFFOLD.md` — ENGINEER's task in sprint 1
- `2-AUDITOR-QUALITY-CHECK-AUDIT.md` — AUDITOR's audit after sprint 2

---

## ⚠️ CRITICAL AGENT RULES

**ALWAYS BUILD WITH SCRIPT**
- ✅ ALWAYS use `./build.sh` for all builds
- ❌ NEVER use direct compilation commands (make, cmake, gcc, etc.)
- The build script ensures consistent build configuration and environment

**AGENTS BUILD CODE FOR USER TO TEST**
- Agents build/modify code ONLY when user explicitly requests
- USER tests and provides feedback
- Agents wait for user approval before proceeding

**AGENTS CAN RUN GIT ONLY IF USER EXPLICITLY ASKS**
- Write code changes without running git commands
- Agent runs git ONLY when user explicitly requests
- Never autonomous git operations
- **When committing:** Always stage ALL changes with `git add -A` before commit
  - ❌ DON'T selectively stage files (agents forget/miss files)
  - ✅ DO `git add -A` to capture every modified file
  - This ensures complete commits with nothing accidentally left unstaged

**LOG MAINTENANCE RULE**
- **All sprint logs must be written from latest to earliest (top to bottom), BELOW this rules section**
- **Only the last 5 sprints are kept in active log**
- **All agent roles except JOURNALIST write [N]-[ROLE]-[OBJECTIVE].md for each completed task**
- **JOURNALIST compiles all task summaries with same sprint number, updates SPRINT-LOG.md as new entry**
- **Only JOURNALIST can add new sprint entry to SPRINT HISTORY**
- **Sprints can be executed in parallel with multiple agents**
- Remove older sprints from active log (git history serves as permanent archive)
- This keeps log focused on recent work
- **JOURNALIST NEVER updates log without explicit user request**
- **During active sprints, only user decides whether to log**
- **All changes must be tested/verified by user, or marked UNTESTED**
- If rule not in this section, agent must ADD it (don't erase old rules)

**NAMING RULE (CODE VOCABULARY)**
- All identifiers must obey project-specific naming conventions (see project docs)
- Variable names: semantic + precise (not `temp`, `data`, `x`)
- Function names: verb-noun pattern (initRepository, detectCanonBranch)
- Struct fields: domain-specific terminology (not generic `value`, `item`, `entry`)
- Type names: PascalCase, clear intent (CanonBranchConfig, not BranchData)

**BEFORE CODING: ALWAYS SEARCH EXISTING PATTERNS**
- ❌ NEVER invent new states, enums, or utility functions without checking if they exist
- ✅ Always grep/search the codebase first for existing patterns
- ✅ Check types, constants, and error handling patterns before creating new ones
- **Methodology:** Read → Understand → Find SSOT → Use existing pattern
- Overcomplications usually mean you missed an existing solution

**TRUST THE LIBRARY, DON'T REINVENT**
- ❌ NEVER create custom helpers for things the library/framework already does
- ✅ Trust the library/framework - it's battle-tested
- **Philosophy:** Libraries are battle-tested. Your custom code is not.
- If you find yourself writing 10+ lines of utility code, stop—the library probably does it

**FAIL-FAST RULE (CRITICAL)**
- ❌ NEVER silently ignore errors (no error suppression)
- ❌ NEVER use fallback values that mask failures
- ❌ NEVER return empty strings/zero values when operations fail
- ✅ ALWAYS check error return values explicitly
- ✅ ALWAYS return errors to caller or log + fail fast
- Better to panic/error early than debug silent failure for hours

**META-PATTERN RULE (CRITICAL)**
- ❌ NEVER start complex task without reading PATTERNS.md
- ✅ ALWAYS use Problem Decomposition Framework for multi-step tasks
- ✅ ALWAYS use Debug Methodology checklist when investigating bugs
- ✅ ALWAYS run Self-Validation Checklist before responding
- ✅ Follow role-specific patterns (COUNSELOR, ENGINEER, SURGEON, MACHINIST, AUDITOR)
- Better to pause and read patterns than repeat documented failures

**SCRIPT USAGE RULE**
- ✅ ALWAYS use scripts from SCRIPTS.md for code editing (when available)
- ✅ Scripts have dry-run mode - use it before actual edit
- ✅ Scripts create backups - verify before committing
- ❌ NEVER use raw sed/awk without safe-edit.sh wrapper (when script available)
- Scripts prevent common mistakes and enforce safety

**⚠️ NEVER EVER REMOVE THESE RULES**
- Rules at top of SPRINT-LOG.md are immutable
- If rules need update: ADD new rules, don't erase old ones
- Any agent that removes or modifies these rules has failed
- Rules protect the integrity of the development log

---

## Quick Reference

### For Agents Starting New Sprint

1. **Check:** Do I see my registration in ROLE ASSIGNMENT REGISTRATION?
2. **If YES:** Proceed with role constraints, include `[Acting as: ROLE]` in responses
3. **If NO:** STOP and ask: "What is my role in this sprint?"

### For Human Orchestrator

**Register agent:**
```
"Read CAROL.md. You are assigned as [ROLE], register yourself in SPRINT-LOG.md"
```

**Verify registration:**
```
"What is your current role?"
```

**Reassign role:**
```
"You are now reassigned as [NEW_ROLE], register yourself in SPRINT-LOG.md"
```

**Complete sprint (call JOURNALIST):**
```
"Read CAROL, act as JOURNALIST. Log sprint [N] to SPRINT-LOG.md"
```

---

## ROLE ASSIGNMENT REGISTRATION

COUNSELOR: OpenCode (glm-4.7)
ENGINEER: OpenCode (MiniMax-M2.1)
SURGEON: OpenCode (glm-4.6)
AUDITOR: [Agent (Model)] or [Not Assigned]
MACHINIST: Gemini (Gemini 2.0 Flash)
JOURNALIST: OpenCode (zai-coding-plan/glm-4.7)

---

<!-- SPRINT HISTORY STARTS BELOW -->
<!-- Only JOURNALIST writes entries here -->
<!-- Latest sprint at top, oldest at bottom -->
<!-- Keep last 5 sprints, rotate older to git history -->

## SPRINT HISTORY

## Sprint 12 - Auto Update, Preferences, and Project Naming
**Date:** 2026-02-02
**Agents:** SURGEON (OpenCode - k2p5)

### Summary
Implemented auto-update with lazy scanning, full TIT-style preferences UI, and comprehensive Generator→Project rename throughout codebase. All user-facing terminology now consistently uses "Project" while internal CMake API retains "Generator" terminology.

### Tasks Completed
- **SURGEON**: Auto Update with Activity Tracking - Added `lastActivityTime` field, tracks user activity in all key handlers (menu, preferences, console), lazy auto-scan skips if user active within 30 seconds (TIT pattern)
- **SURGEON**: Preferences UI (TIT Copy) - Created `internal/ui/preferences.go` with exact TIT implementation: `PreferenceRow` struct, `BuildPreferenceRows()`, `RenderPreferencesMenu()` with 3-column layout (emoji | label | value), `RenderPreferencesWithBanner()` for 50/50 split
- **SURGEON**: Preference Functionality - 3 preference items: Auto-scan toggle (ON/OFF), Scan Interval (+/- 1min, =/_ 10min), Theme cycling with immediate reload; Enter/Space toggles, +/- adjusts interval, `/` or ESC returns to menu
- **SURGEON**: Save/Restore Last Options - Added `BuildConfig` to config with `last_project` and `last_configuration`, restore on app init, save when changed
- **SURGEON**: Generator→Project Rename - All user-facing comments, hints, labels, error messages changed from "Generator" to "Project": `SelectedProject`, `AvailableProjects`, `CycleToNextProject()`, menu labels, footer hints

### Files Modified
- `internal/app/app.go` - Added `lastActivityTime` field, lazy auto-scan logic, activity tracking in key handlers, `GetVisiblePreferenceRows()`, `TogglePreferenceAtIndex()`, interval adjustment keys (+/-)
- `internal/app/init.go` - Restore `last_project` and `last_configuration` from config on startup
- `internal/app/messages.go` - Changed "generator" to "project" in hints and messages
- `internal/config/config.go` - Added `BuildConfig` struct with `last_project` and `last_configuration`, getter/setter methods
- `internal/state/project.go` - Renamed `SelectedGenerator`→`SelectedProject`, `AvailableGenerators`→`AvailableProjects`, `CycleToNextGenerator`→`CycleToNextProject`, added `SetSelectedProject()`, `SetConfiguration()`
- `internal/ui/menu.go` - Updated comments and hint text
- `internal/ui/preferences.go` - NEW FILE: Complete TIT-style preferences implementation

### Alignment Check
- [x] LIFESTAR principles followed (lazy updates, immediate persistence)
- [x] NAMING-CONVENTION.md adhered (Project for user-facing, Generator for CMake API)
- [x] ARCHITECTURAL-MANIFESTO.md principles applied (TIT pattern matching)

### Build Status
- Build completes successfully ✓
- All changes tested and verified

## Sprint 11 - Codebase Audit and Cleanup
**Date:** 2026-01-29
**Agents:** AUDITOR (Amp - Claude Sonnet 4), MACHINIST (OpenCode - glm-4.7)

### Summary
Full codebase audit identified 15 issues (3 CRITICAL, 4 HIGH, 5 MEDIUM, 3 LOW). Fixed all CRITICAL and HIGH priority issues: removed dead "Ninja Multi-Config" logic, implemented CMake name approach with directory mapping, extracted command streaming helper, eliminated code duplication in menu handlers, and achieved SSOT for all generator references.

### Tasks Completed
- **AUDITOR**: Codebase Audit - Full audit revealing 15 issues: SSOT violations (build paths, streaming), code duplication (menu handlers, command streaming), dead code (Ninja Multi-Config), inconsistent generator naming, missing unit tests, god object tendency
- **MACHINIST**: Cleanup (6 FIX items) - Removed "Ninja Multi-Config" from all locations, updated Visual Studio versions to CMake names (VS2026, VS2022), added GetDirectoryName() and GetGeneratorNameFromDirectory() helpers, updated all 7 build path constructions to use GetDirectoryName()
- **MACHINIST**: Code Polish (3 CRITICAL + 1 HIGH) - Created StreamCommand() helper (removed 120+ lines duplication), created executeShortcut() helper (eliminated handler duplication), created generator name constants (SSOT in utils/generators.go)

### Files Modified
- `internal/utils/generators.go` - NEW: Generator name constants (SSOT), GetDirectoryName(), GetGeneratorNameFromDirectory(), IsGeneratorIDE()
- `internal/utils/stream.go` - NEW: StreamCommand() helper for unified command streaming
- `internal/utils/generator.go` - Updated validGenerators slice and switch cases, removed Ninja Multi-Config
- `internal/utils/platform.go` - Updated GetPlatformGenerators() and GetDefaultGenerator(), changed Windows default from Ninja Multi-Config to Ninja
- `internal/state/project.go` - Updated VS names to CMake format, removed NinjaMulti, integrated GetDirectoryName()
- `internal/ops/setup.go` - Updated buildDir to use GetDirectoryName(), refactored to use StreamCommand()
- `internal/ops/build.go` - Updated buildDir to use GetDirectoryName(), refactored to use StreamCommand()
- `internal/ops/clean.go` - Updated buildDir to use GetDirectoryName()
- `internal/ops/open.go` - Updated VS switch cases and buildDir to use GetDirectoryName()
- `internal/app/op_regenerate.go` - Updated buildDir to use GetDirectoryName()
- `internal/app/app.go` - Added executeShortcut() helper

### Notes
- Build completes successfully ✓
- All grep verifications pass
- All 6 FIX items verified complete
- Removed 120+ lines of duplicated code
- SSOT achieved: generator names in utils/generators.go, directory names via GetDirectoryName()
- All generators (Xcode, Ninja, VS2022, VS2019) are multi-config
- CMake-style generator names: "Visual Studio 18 2026", "Visual Studio 17 2022"
- Directory names: "VS2026", "VS2022"
- Documentation updates pending: SPEC.md and ARCHITECTURE.md require updates to remove Unix Makefiles and single-config references

## Sprint 10 - Console TIT Alignment
**Date:** 2026-01-29
**Agents:** COUNSELOR (OpenCode - glm-4.7), ENGINEER (OpenCode - MiniMax-M2.1, glm-4.7), SURGEON (OpenCode - glm-4.6)

### Summary
Achieved 100% visual and behavioral parity between Cake Console and TIT Console. Fixed critical menu index bug causing wrong operations, implemented real-time streaming output with color-coded types, added process cancellation via ESC, created confirmation dialogs for destructive operations, and replaced entire theme system with TIT's exact implementation.

### Tasks Completed
- **COUNSELOR**: Console TIT Alignment Kickoff - Comprehensive 11-task plan for console alignment, identified critical issues (menu Generator label, index calculation bug, IsMultiConfig removal, console flooding, operation completion issues)
- **COUNSELOR**: Critical Bug Fix - Fixed menu index mismatch: GetArrayIndex/GetVisibleIndex only checked Visible, not IsSelectable, causing Clean to trigger Build operation
- **COUNSELOR**: Console Bugs - Identified 3 critical bugs: ESC doesn't work during operations (no process kill), console flooded with lines (9k+), operations never complete (goroutine leaks)
- **ENGINEER**: Console TIT Alignment (Tasks 1-5) - Updated menu labels (Generator→Project), fixed GetVisibleIndex/GetArrayIndex to check both Visible and IsSelectable, removed IsMultiConfig field and methods, renamed GetGeneratorLabel→GetProjectLabel
- **ENGINEER**: Console TIT Alignment Updated (Tasks 6-10) - Fixed output color types (Info/cyan, Stdout/gray, Stderr/coral/red, Status/cyan, Warning/orange), implemented real-time streaming via StdoutPipe/StderrPipe+bufio.Scanner, simplified Clean to fast rm -rf, created op_regenerate.go with Clean→Generate sequence, added confirmation dialogs for clean/regenerate
- **SURGEON**: Console TIT Alignment - Complete console system overhaul: process cancellation (ESC calls cancelContext), console auto-scroll (every 100ms), confirmation dialog key routing fix (arrow keys: left→Yes, right→No), removed console centering, copied complete TIT theme system (5 themes), added OutputRefreshMsg for console refresh ticks

### Files Modified
- `internal/ui/console.go` - Copied from TIT, removed blank line after OUTPUT title, left-aligned wrapped lines
- `internal/ui/theme.go` - Complete TIT theme system (5 themes: gfx, spring, summer, autumn, winter), config path to `.config/cake`
- `internal/ui/buffer.go` - Verified identical to TIT
- `internal/ui/menu.go` - Changed ID from "generator" to "project", updated labels
- `internal/state/project.go` - Removed IsMultiConfig field and methods, renamed GetGeneratorLabel→GetProjectLabel
- `internal/ops/setup.go` - Complete rewrite with context support, uses exec.CommandContext for cancellable processes, streams stdout/stderr in real-time via goroutines, changed callback signature to include lineType
- `internal/ops/build.go` - Changed callback signature, added real-time streaming, removed isMultiConfig
- `internal/ops/clean.go` - Simplified to rm -rf, added config parameter, updated messages
- `internal/ops/open.go` - Added config and projectRoot parameters, changed callback signature
- `internal/app/op_generate.go` - Updated callback, removed isMultiConfig usage, added context creation for cancellation
- `internal/app/op_build.go` - Updated callback, removed isMultiConfig usage, added consoleAutoScroll = true on operation start
- `internal/app/op_clean.go` - Updated callback, added consoleAutoScroll = true on operation start
- `internal/app/op_open.go` - Updated callback
- `internal/app/op_regenerate.go` - NEW FILE with Clean then Generate sequence, context creation for cancellation
- `internal/app/messages.go` - Added RegenerateCompleteMsg, OutputRefreshMsg
- `internal/app/menu.go` - Updated to use GetProjectLabel()
- `internal/app/app.go` - Added consoleAutoScroll field, added runningCmd and cancelContext for process cancellation, fixed confirmation dialog key routing (check BEFORE key dispatcher), fixed arrow key mapping (left→Yes, right→No), added spacebar support, ctrl+c handling, fixed "g" shortcut to use "regenerate" row ID, console mode bypasses RenderReactiveLayout to avoid centering, updated all completion handlers to check for abort, ESC handler prints abort message to console, added confirmation dialogs for clean/regenerate

### Notes
- Build completes successfully ✓
- Menu index bug fixed: Clean now triggers Clean, not Build
- Real-time streaming: Output appears line-by-line as generated
- Color-coded output: Info/cyan, Stdout/gray, Stderr/coral/red, Status/cyan, Warning/orange
- Auto-scroll: Enabled during operations, disabled when user scrolls manually, re-enabled by new operations
- Process cancellation: ESC kills cmake process, prints abort message, shows "Press ESC to return to menu"
- Confirmation dialogs: Clean and Regenerate default to "No", arrow keys fixed (left→Yes, right→No), Y always confirms, N always cancels, ESC always cancels
- Console bypasses RenderReactiveLayout (was centering), renders full-screen with internal OUTPUT title
- Theme system: 5 themes generated in `~/.config/cake/themes/` (gfx, spring, summer, autumn, winter)
- IsMultiConfig removed: All projects now use multi-config path structure

## Sprint 8 - Menu Navigation, Layout, and Indexing Fixes
**Date:** 2026-01-29
**Agents:** COUNSELOR (OpenCode - glm-4.7), ENGINEER (OpenCode - MiniMax-M2.1)

### Summary
Fixed comprehensive menu issues through 5 iterations: navigation with hidden items, separator handling, Unix Makefiles removal, Generate/Regenerate state, layout improvements (value width, generator display, spacing), and critical indexing bug. Menu now uses visible selectable indices (0-4 or 0-5), separator is visible but not selectable, and menu always renders exactly 7 lines regardless of visibility.

### Tasks Completed
- **ENGINEER**: Menu Navigation Fix - Rewrote GetVisibleRows(), added GetVisibleIndex() and GetArrayIndex(), updated navigation to use visible indices, fixed separator width to match menu content
- **ENGINEER**: Menu Fixes - Removed Unix Makefiles from all locations, added IsSelectable field to MenuRow, implemented "Ninja Multi" display truncation
- **ENGINEER**: Menu Layout Fixes - Pass all menu items to RenderCakeMenu, implemented Generate vs Regenerate label switching, widened value column to 12 chars, changed "Ninja Multi" to "Ninja (multi)", added space between label and value when selected
- **ENGINEER**: Critical Menu Fixes - Regenerate menu after generator/configuration changes, widened value column to 14 chars, added space between label and value, removed old padding logic
- **ENGINEER**: Menu Stability Fixes - Fixed label padding calculation (reduce by 1 when selected and has value), render hidden rows as empty lines to maintain fixed 7-line height
- **ENGINEER**: Urgent Menu Indexing Fix - Complete rewrite of RenderCakeMenu indexing: renamed visibleIndex to visibleSelectableIndex, only increment for visible AND selectable rows, separator never increments, hidden rows don't increment, removed padding manipulation logic

### Files Modified
- `internal/app/app.go` — Rewrote GetVisibleRows() to filter hidden rows, added GetVisibleIndex() and GetArrayIndex() methods, updated ToggleRowAtIndex() to work with visible indices, updated handleMenuKeyPress() navigation (↑↓), updated shortcut handlers (g, o, b, c) to use visible indices, added menu regeneration after generator/configuration changes, pass a.menuItems to RenderCakeMenu instead of visibleRows
- `internal/app/footer.go` — Updated getMenuFooter() to use visible rows correctly
- `internal/ui/menu.go` — Fixed separator width (contentWidth → menuBoxWidth), added IsSelectable field to MenuRow struct, set separator IsSelectable=false, updated GenerateMenuRows() with hasBuild parameter for Generate/Regenerate switching, increased valueColWidth from 10→12→14, added space between label and value when selected, fixed label padding calculation (reduce by 1 when selected), render hidden rows as empty lines, complete rewrite of RenderCakeMenu indexing logic (visibleSelectableIndex, only increment for visible+selectable)
- `internal/app/menu.go` — Pass hasBuild parameter to GenerateMenuRows, update to use GetGeneratorLabel()
- `internal/state/project.go` — Removed Unix Makefiles fallback (lines 137-142), updated GetGeneratorLabel() with truncation: "Ninja Multi-Config" → "Ninja (multi)", "Visual Studio 17 2022" → "VS 2022", "Visual Studio 16 2019" → "VS 2019"
- `internal/utils/generator.go` — Removed Unix Makefiles from validGenerators list, removed switch case for Unix Makefiles
- `internal/utils/platform.go` — Removed Unix Makefiles from Linux/default returns
- `internal/ops/build.go` — Updated comments to remove Makefiles references
- `internal/ops/clean.go` — Updated comments to remove Makefiles references
- `SPEC.md` — Updated generator documentation (5 locations)

### Notes
- Build completes successfully ✓
- Navigation now uses visible selectable indices (0-4 or 0-5 depending on visible items)
- Separator is visible but not selectable (skipped in navigation, always at line 3)
- Unix Makefiles completely removed from generator cycling
- Generator cycling: Xcode → Ninja → Ninja Multi-Config → Xcode (macOS)
- Generate vs Regenerate label switches based on build state (fresh project = "Generate", after build = "Regenerate")
- Value column 14 chars wide (fits "Ninja (multi)" = 13 chars)
- "Ninja (multi)" displays correctly in value column
- Space appears between label and value when selected via reduced padding, not added space
- Menu always renders exactly 7 lines (hidden items as empty lines)
- Menu position stable when items hide/show (no shifting)
- Footer shows correct hint for selected item
- Can navigate to Clean when visible (navigation not capped at Build)
- Value column stays aligned regardless of selection
- Menu regenerates immediately after generator or configuration changes

## Sprint 7 - TIT Footer Alignment and Unix Makefiles Removal
**Date:** 2026-01-28
**Agents:** COUNSELOR (OpenCode - glm-4.7), ENGINEER (OpenCode - MiniMax-M2.1)

### Summary
Implemented complete footer system aligning CAKE with TIT's exact architecture. Replaced scattered footerHint management with structured footer manager (internal/app/footer.go), implemented mode-specific footer content (menu hints vs console scroll status), and removed Unix Makefiles from generator options across entire codebase.

### Tasks Completed
- **COUNSELOR**: TIT Footer Alignment Planning - Designed 7-phase implementation plan to align footer with TIT structure, defined menu hints and console scroll status patterns
- **ENGINEER**: Unix Makefiles Removal - Removed Unix Makefiles from generator options across 6 files (project.go, generator.go, platform.go, build.go, clean.go, SPEC.md)
- **ENGINEER**: TIT Footer Implementation (Phases 2-7) - Complete rewrite of internal/ui/footer.go with FooterShortcut struct and renderers, created internal/app/footer.go manager with getMenuFooter()/getConsoleFooter(), added FooterHintShortcuts map to messages.go, added Hint field to MenuRow, removed 50+ lines of manual footerHint management from app.go and operation files

### Files Modified
- `internal/ui/footer.go` — Complete rewrite: FooterShortcut struct, FooterStyles, RenderFooter(), RenderFooterOverride(), RenderFooterHint()
- `internal/app/footer.go` — NEW: Footer manager with GetFooterContent(), getMenuFooter(), getConsoleFooter(), computeConsoleScrollStatus()
- `internal/app/messages.go` — Added FooterHintShortcuts map with mode-specific shortcuts
- `internal/ui/menu.go` — Added Hint field to MenuRow struct, populated all 7 rows with descriptions
- `internal/app/app.go` — Removed footerHint field, updated View() to use GetFooterContent(), removed 20+ footerHint assignments
- `internal/app/init.go` — Removed footerHint initialization
- `internal/app/op_generate.go` — Removed footerHint assignment
- `internal/app/op_build.go` — Removed footerHint assignment
- `internal/app/op_clean.go` — Removed footerHint assignment
- `internal/app/op_open.go` — Removed 2 footerHint assignments
- `internal/state/project.go` — Removed Unix Makefiles fallback
- `internal/utils/generator.go` — Removed "Unix Makefiles" from validGenerators
- `internal/utils/platform.go` — Removed "Unix Makefiles" from Linux/default returns
- `internal/ops/build.go` — Updated comments to remove Makefiles reference
- `internal/ops/clean.go` — Updated comments to remove Makefiles reference
- `SPEC.md` — Updated generator documentation (5 locations)

### Notes
- Build completes successfully ✓
- Footer now shows menu item hints in menu mode (from MenuRow.Hint field)
- Footer shows scroll shortcuts + status in console mode (left/right split layout)
- Ctrl+C timeout works as global override in both modes
- Unix Makefiles no longer appears in generator cycling
- All operation files cleaned up - no more manual footerHint management

## Sprint 6 - TIT Compliance Refactoring Complete
**Date:** 2026-01-28
**Agents:** COUNSELOR (OpenCode - glm-4.7), ENGINEER (OpenCode - glm-4.7), AUDITOR (OpenCode - glm-4.7)

### Summary
Completed all 6 phases of TIT compliance refactoring, significantly improving code organization and eliminating architectural debt. Eliminated 8 locations of code duplication, improved TIT compliance from 67% to 90%+, extracted operation handlers from 877-line app.go into focused op_*.go files, and introduced structured async state and message dispatcher patterns.

### Tasks Completed
- **COUNSELOR**: TIT Refactor Planning - Identified fix requirements for 6-phase compliance refactoring, created comprehensive kickoff documents
- **ENGINEER**: Phase 1 - SSOT Extract Build Path - Eliminated 5+ duplicated build directory constructions, created GetBuildDirectory() method, updated all 7 operation commands
- **ENGINEER**: Phase 2 - SSOT Extract isMultiConfig Detection - Eliminated 3 duplicated for-loop patterns, replaced all loops with IsGeneratorMultiConfig() method
- **ENGINEER**: Phase 3 - TIT Footer Renderer - Created internal/ui/footer.go with RenderFooter() function matching TIT pattern
- **ENGINEER**: Phase 4 - TIT Async State Extraction - Created internal/app/async_state.go with AsyncState struct, replaced 16 bool references in app.go
- **ENGINEER**: Phase 5 - TIT Message Dispatchers - Created internal/app/dispatchers.go with MessageHandler interface, integrated WindowSizeHandler and KeyDispatcher into app.go Update()
- **ENGINEER**: Phase 6 - TIT Operation Handlers - Extracted all operation methods from app.go into op_generate.go, op_build.go, op_clean.go, op_open.go, removed 179 lines from app.go (20% reduction)
- **AUDITOR**: Sprint 5 Compliance Audit - Provided reference for improvements, identified compliance gaps

### Files Modified
- `internal/ui/footer.go` — NEW: RenderFooter() function for TIT-style footer rendering
- `internal/app/async_state.go` — NEW: AsyncState struct with methods (Start, End, Abort, ClearAborted, IsActive, IsAborted, CanExit, SetExitAllowed)
- `internal/app/dispatchers.go` — NEW: MessageHandler interface, WindowSizeHandler, KeyDispatcher
- `internal/app/op_generate.go` — NEW: startGenerateOperation(), cmdGenerateProject() extracted
- `internal/app/op_build.go` — NEW: startBuildOperation(), cmdBuildProject() extracted
- `internal/app/op_clean.go` — NEW: startCleanOperation(), cmdCleanProject() extracted
- `internal/app/op_open.go` — NEW: startOpenIDEOperation(), cmdOpenIDE(), startOpenEditorOperation(), cmdOpenEditor() extracted
- `internal/state/project.go` — Added GetBuildDirectory() and IsGeneratorMultiConfig() methods
- `internal/ops/setup.go` — Initially added then removed duplicate GetBuildDirectory(), inlined build path logic
- `internal/ops/build.go` — Updated to use inline build path logic with filepath import
- `internal/ops/clean.go` — Updated to use inline build path logic with filepath import
- `internal/app/app.go` — 179 lines removed, integrated dispatchers, replaced bools with asyncState (698 lines, -20%)

### Notes
- Code duplication eliminated from 5+ locations to 0 for build path logic and isMultiConfig detection
- TIT compliance improved from 67% to 90%+ through proper state management, message dispatchers, and operation separation
- app.go reduced from 877 to 698 lines (-179 lines, -20%)
- Build status: ✓ Built successfully to ~/.cake/bin/cake_x64, ✓ Symlinked to ~/.local/bin/cake
- Runtime testing deferred to MACHINIST phase - need to test all keyboard shortcuts, window resize, and different generators
- Key fix: Removed duplicate GetBuildDirectory() from ops/setup.go to avoid new duplication
- All phases build successfully, verified metrics before reporting

## Sprint 5 - Menu Layout and Alignment Fix
**Date:** 2026-01-28
**Agents:** COUNSELOR (OpenCode - glm-4.7), ENGINEER (OpenCode - MiniMax-M2.1)

### Summary
Fixed critical menu alignment issues and implemented comprehensive menu refactoring matching TIT's preference layout. Resolved emoji width handling using lipgloss.Width(), fixed selection index bug, added keyboard shortcuts, implemented conditional visibility with placeholder rows, and established fixed 7-item menu structure.

### Tasks Completed
- **COUNSELOR**: Menu Alignment Analysis - Identified emoji width issue (2 display columns but 4 UTF-8 bytes), designed solution using lipgloss.Width() for proper padding calculation
- **ENGINEER**: Menu Refactoring Implementation - Rewrote renderPreferenceMenu() with lipgloss.Width() for proper emoji handling, fixed 7-item structure with conditional visibility, added keyboard shortcuts (g, o, b, c), fixed critical selection index bug where navigation used array indices but rendering used visible indices
- **ENGINEER**: Menu Architecture Rewrite - Complete rewrite of app.go menu handling (65 lines deleted, 50 lines added), removed renderPreferenceMenu() in favor of ui.RenderCakeMenu(), added 6 new functions (GetVisibleRows, GetVisiblePreferenceRows, ToggleRowAtIndex, RowIndexByID, TogglePreferenceAtIndex, executeRowAction), changed to row ID-based navigation instead of index calculations

### Files Modified
- `internal/app/menu.go` — Added Shortcut field to PreferenceRow struct, rewrote GenerateMenu() to return fixed 7 rows with conditional visibility (placeholder empty rows for hidden items), removed shortcuts from toggle rows
- `internal/app/app.go` — Complete rewrite of menu handling logic: removed renderPreferenceMenu() (65 lines deleted), added 6 new functions (GetVisibleRows, GetVisiblePreferenceRows, ToggleRowAtIndex, RowIndexByID, TogglePreferenceAtIndex, executeRowAction), updated renderMenuWithBanner() to call ui.RenderCakeMenu() with correct argument order, changed all PreferenceRow → ui.MenuRow type references, changed row.Separator → row.ID == "separator" checks, simplified handleMenuKeyPress() to use ID-based navigation, removed unused fmt import
- `internal/app/init.go` — Changed menuItems type from []PreferenceRow to []ui.MenuRow

### Notes
- Fixed column widths: shortcut 3ch, emoji 3ch (centered), label 20ch, value 12ch (total 48ch)
- Emojis display as 2 terminal columns but are 4 UTF-8 bytes - lipgloss.Width() correctly measures display width
- Highlight pattern: Only label column gets background highlight when selected
- Critical bug fix: selectedIndex now uses visible indices (0-5) instead of array indices (0-6)
- Navigation skips Separator: true and Visible: false rows
- Actions (Generate, Open IDE, Build, Clean) use shortcuts g, o, b, c
- Toggles (Generator, Configuration) use Enter/Space only (no shortcuts)
- Menu architecture complete rewrite: Removed 65 lines of inline rendering, added 50 lines of structured functions
- executeRowAction handles all 7 menu actions: generator (cycle), regenerate, openIde, configuration (cycle), build, clean
- Changed from index-based to row ID-based navigation for simpler logic
- Removed inline menu rendering in favor of ui.RenderCakeMenu() with 4-column layout (shortcut | emoji | label | value)
- Confirmation dialogs for regenerate/clean not yet re-implemented (requires ui.ConfirmationDialogConfig pattern)


## Sprint 4 - Header Rendering Fix
**Date:** 2026-01-28
**Agents:** COUNSELOR (OpenCode - glm-4.7), ENGINEER (OpenCode - MiniMax-M2.1)

### Summary
Fixed inconsistent header rendering between menu mode (2 lines) and console mode (3 lines) by removing double placement pattern in RenderReactiveLayout().

### Tasks Completed
- **COUNSELOR**: Header Rendering Inconsistency Investigation - Identified root cause of double placement pattern causing header clipping in menu mode
- **ENGINEER**: Layout Architecture Fix - Removed outer lipgloss.Place() wrapper from RenderReactiveLayout() to match TIT architecture exactly

### Files Modified
- `cake/internal/ui/layout.go` — Removed outer lipgloss.Place() wrapper from RenderReactiveLayout() (lines 82-92 → lines 82-83)

### Notes
- Double placement pattern was causing lipgloss to re-layout entire view, clipping header when content was short
- Fix aligns CAKE with proven TIT architecture pattern (simple JoinVertical without outer Place wrapper)
- Header now renders consistently 3 lines across all modes (menu, console, preferences)

## Sprint 3 - Header Rendering Fixes and Project Name Detection
**Date:** 2026-01-28
**Agents:** COUNSELOR (Amp - Claude Sonnet 4), ENGINEER (OpenCode - MiniMax-M2.1), MACHINIST (Gemini - Gemini 2.0 Flash, OpenCode - MiniMax-M2.1)

### Summary
Fixed critical header rendering issues and implemented robust project name detection. Replaced broken CMake-based project name extraction with layered file parsing, resolved initialization timing problems, and aligned layout rendering with TIT architecture for consistent headers across all modes.

### Tasks Completed
- **COUNSELOR**: Project Name Detection Planning - Designed 4-layer detection strategy: KANJUT Parameters.xml > set(PROJECT_NAME) > project() literal > directory name
- **ENGINEER**: Project Name Detection Implementation - Replaced CMake-based GetProjectName() with 137 lines of pure Go parsing across 3 helper functions, removing external dependencies
- **MACHINIST**: Header Rendering Fixes - Fixed ProjectState initialization to resolve cwd immediately, updated RenderReactiveLayout() to match TIT pattern using lipgloss.Place for consistent header positioning

### Files Modified
- `internal/state/project.go` - Complete rewrite of GetProjectName() with 4-layer detection, added extractFromParametersXML(), extractFromSetProjectName(), extractFromProjectCall(), fixed NewProjectState() to resolve cwd immediately
- `internal/ui/layout.go` - Updated RenderReactiveLayout() to match TIT pattern exactly (lipgloss.Place for header/footer, adjusted content height for terminal quirks)

### Notes
- Cross-platform project name detection without CMake execution
- Layered detection: KANJUT > set() > project() literal > directory fallback
- Header rendering now consistent across Menu, Preferences, and Console modes
- ProjectState initialization resolves WorkingDirectory immediately instead of placeholder "."

## Sprint 2 - Critical Fixes and Theme System Alignment
**Date:** 2026-01-28
**Agents:** COUNSELOR (OpenCode - Claude Sonnet 4.5), ENGINEER (OpenCode - MiniMax-M2.1), MACHINIST (OpenCode - MiniMax-M2.1)

### Summary
Addressed critical functionality and visual consistency issues. Fixed core operations that weren't executing, implemented universal project name resolution via CMake, corrected header format to consistent 3-line layout, aligned menu alignment, styled confirmation dialogs, and synchronized entire theme system with TIT architecture exactly.

### Tasks Completed
- **COUNSELOR**: Critical Fixes Planning - Identified and documented 6 core issues: non-executing operations, wrong project names, inconsistent headers, menu alignment, plain dialogs, console inconsistencies
- **ENGINEER**: Core Functionality Fixes - Fixed header sizing, robust navigation logic, implemented TIT header architecture pattern, resolved command propagation for Generate/Build/Clean operations, corrected console sizing
- **MACHINIST**: Theme System Alignment - Complete theme system overhaul to match TIT exactly: fixed confirmation dialog colors, implemented TIT header rendering with HeaderState pattern, updated theme loading to use "gfx.toml", ensured 5-theme system consistency

### Files Modified
- `internal/ui/sizing.go` - Fixed HeaderInnerWidth from termWidth to innerWidth
- `internal/ui/header.go` - Complete rewrite matching TIT architecture (HeaderState, RenderHeaderInfo, RenderHeader)
- `internal/ui/theme.go` - Removed duplicate DefaultThemeTOML, added EnsureFiveThemesExist(), fixed theme loading, updated all 5 theme confirmation dialog colors
- `internal/ui/layout.go` - Fixed RenderConfirmDialog() to use theme colors instead of hardcoded values
- `internal/app/app.go` - Updated to TIT header pattern, fixed handleMenuKeyPress with robust navigation, fixed handlePreferencesKeyPress, fixed console sizing
- `internal/state/project.go` - Enhanced GetProjectName() using CMake inline script for universal project name resolution
- `internal/app/menu.go` - Fixed ToggleRowAtIndex signature to return tea.Cmd, implemented confirmation dialogs from TIT

### Notes
- Universal project name resolution uses inline CMake script to handle all variables and includes
- TIT header architecture fully adopted: HeaderState → RenderHeaderInfo → RenderHeader
- Theme system now matches TIT exactly with proper confirmation dialog backgrounds per theme
- Command propagation fixed - Generate/Build/Clean operations now execute properly
- Robust navigation prevents getting stuck on separators at boundaries

## Sprint 1 - Preference-Style Dynamic Menu Implementation
**Date:** 2026-01-28
**Agents:** COUNSELOR (OpenCode - Claude Sonnet 4.5), ENGINEER (OpenCode - MiniMax-M2.1)

### Summary
Complete re-implementation of cake menu system following TIT's preference-style pattern. Transformed from static 3-item menu to dynamic preference menu with generator detection, configuration cycling, and conditional action visibility. Fixed critical UI layout issues including 50/50 split, centering, separator navigation, and positioning.

### Tasks Completed
- **COUNSELOR**: Dynamic Menu Planning - Created comprehensive requirements, design decisions, and 10-phase implementation plan with kickoff document
- **ENGINEER**: Dynamic Menu Implementation - Complete re-implementation (Phases 1-9): Build directory structure, generator detection (system tools), preference-style menu, cycling logic, enhanced operations, config system (TOML), preferences screen, auto-scan ticker, confirmation dialogs
- **ENGINEER**: UI Fixes - Restored 50/50 split layout, centered menu/banner in columns, fixed separator row navigation, corrected vertical positioning

### Files Modified
- `internal/state/project.go` - Complete rewrite with Generator struct, BuildInfo, DetectAvailableGenerators(), scanBuildDirectories()
- `internal/app/menu.go` - Complete rewrite with PreferenceRow struct, GenerateMenu(), GetVisibleRows()
- `internal/app/app.go` - Complete rewrite for preference menu, new operation handlers, renderMenuWithBanner(), separator skipping
- `internal/app/modes.go` - Simplified to ModeMenu and ModePreferences
- `internal/app/messages.go` - Added GenerateCompleteMsg
- `internal/ui/theme.go` - Added SeparatorColor field
- `internal/ops/setup.go` - Fixed paths, added config/isMultiConfig parameters
- `internal/ops/build.go` - Fixed paths, added config/isMultiConfig parameters
- `internal/ops/clean.go` - Fixed paths, added config/isMultiConfig parameters
- `internal/ops/open.go` - NEW: ExecuteOpenIDE() for Xcode/VS, ExecuteOpenEditor() for Neovim
- `internal/config/config.go` - NEW: TOML persistence for build, auto_update, and appearance configs

### Notes
- Build structure: `Builds/<Generator>/<Config?>` (multi-config vs single-config)
- Generator detection checks available tools (not disk state)
- Preferences screen accessible via `/` shortcut with theme cycling
- Auto-scan ticker with configurable interval (5, 10, 15, 30 minutes)
- Confirmation dialogs for destructive actions (Clean, Regenerate)
- Menu layout: 50/50 split (menu left, banner right) matching TIT pattern


<!-- Actual sprint entries go here, written by JOURNALIST -->

---

## [N]-[ROLE]-[OBJECTIVE].md Format Reference

**File naming:** `[N]-[ROLE]-[OBJECTIVE].md`  
**Examples:**
- `[N]-ENGINEER-MERMAID-MODULE.md`
- `[N]-MACHINIST-ERROR-HANDLING.md`
- `[N]-SURGEON-COMPILE-FIX.md`

**Content format:**
```markdown
# Sprint [N] Task Summary

**Role:** [ROLE NAME]
**Agent:** [CLI Tool (Model)]
**Date:** 2026-01-28
**Time:** [HH:MM]
**Task:** [Brief task description]

## Objective
[What was accomplished in 1-2 sentences]

## Files Modified ([X] total)
- `path/to/file.ext` — [brief description of changes]
- `path/to/file2.ext` — [brief description of changes]

## Notes
- [Important learnings, blockers, or decisions]
- [Any warnings or follow-up needed]
```

**Lifecycle:**
1. Agent completes task
2. Agent writes [N]-[ROLE]-[OBJECTIVE].md
3. JOURNALIST compiles all [N]-[ROLE]-[OBJECTIVE].md files into SPRINT-LOG.md entry
4. JOURNALIST deletes all [N]-[ROLE]-[OBJECTIVE].md files after compilation

---

**End of SPRINT-LOG.md Template**

Copy this template to your project root as `SPRINT-LOG.md` and customize:
- Project name
- Repository URL/path
- Starting date
- Add project-specific rules to CRITICAL AGENT RULES section

Rock 'n Roll!  
JRENG!
