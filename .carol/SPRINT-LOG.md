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
