# SPRINT-LOG.md

**Project:** cake  
**Repository:** github.com/jrengmusic/cake  
**Started:** 2026-03-28

**Purpose:** Long-term context memory across sessions. Tracks completed work, technical debt, and unresolved issues. Written by PRIMARY agents only when ARCHITECT explicitly requests.

---

## Notation Reference

**[N]** = Sprint Number (e.g., `1`, `2`, `3`...)

**Sprint:** A discrete unit of work completed by one or more agents, ending with ARCHITECT approval ("done", "good", "commit")

---

## CRITICAL RULES

**AGENTS BUILD CODE FOR ARCHITECT TO TEST**
- Agents build/modify code ONLY when ARCHITECT explicitly requests
- ARCHITECT tests and provides feedback
- Agents wait for ARCHITECT approval before proceeding

**AGENTS NEVER RUN GIT COMMANDS**
- Write code changes without running git commands
- Agent runs git ONLY when user explicitly requests
- Never autonomous git operations
- **When committing:** Always stage ALL changes with `git add -A` before commit
  - Don't selectively stage files (agents forget/miss files)
  - Do `git add -A` to capture every modified file

**SPRINT-LOG WRITTEN BY PRIMARY AGENTS ONLY**
- **COUNSELOR** or **SURGEON** write to SPRINT-LOG
- Only when user explicitly says: `"log sprint"`
- No intermediate summary files
- No automatic logging after every task
- Latest sprint at top, keep last 5 entries

**NAMING RULE (CODE VOCABULARY)**
- All identifiers must obey project-specific naming conventions (see carol/NAMES.md)
- Variable names: semantic + precise (not `temp`, `data`, `x`)
- Function names: verb-noun pattern (initRepository, detectCanonBranch)
- Struct fields: domain-specific terminology (not generic `value`, `item`, `entry`)
- Type names: PascalCase, clear intent (CanonBranchConfig, not BranchData)

**BEFORE CODING: ALWAYS SEARCH EXISTING PATTERNS**
- NEVER invent new states, enums, or utility functions without checking if they exist
- Always grep/search the codebase first for existing patterns
- Check types, constants, and error handling patterns before creating new ones
- **Methodology:** Read → Understand → Find SSOT → Use existing pattern

**TRUST THE LIBRARY, DON'T REINVENT**
- NEVER create custom helpers for things the library/framework already does
- Trust the library/framework - it's battle-tested

**FAIL-FAST RULE (CRITICAL)**
- NEVER silently ignore errors (no error suppression)
- NEVER use fallback values that mask failures
- NEVER return empty strings/zero values when operations fail
- NEVER use early returns
- ALWAYS check error returns explicitly
- ALWAYS return errors to caller or log + fail fast

**NEVER REMOVE THESE RULES**
- Rules at top of SPRINT-LOG.md are immutable
- If rules need update: ADD new rules, don't erase old ones

---

## Quick Reference

### For Agents

**When user says:** `"log sprint"`

1. **Check:** Did I (PRIMARY agent) complete work this session?
2. **If YES:** Write sprint block to SPRINT-LOG.md (latest first)
3. **Include:** Files modified, changes made, alignment check, technical debt

### For User

**Activate PRIMARY:**
```
"@CAROL.md COUNSELOR: Rock 'n Roll"
"@CAROL.md SURGEON: Rock 'n Roll"
```

**Log completed work:**
```
"log sprint"
```

**Invoke subagent:**
```
"@oracle analyze this"
"@engineer scaffold that"
"@auditor verify this"
```

**Available Agents:**
- **PRIMARY:** COUNSELOR (domain specific strategic analysis), SURGEON (surgical precision problem solving)
- **Subagents:** Pathfinder, Oracle, Engineer, Auditor, Machinist, Librarian

---

<!-- SPRINT HISTORY STARTS BELOW -->
<!-- Latest sprint at top, oldest at bottom -->
<!-- Keep last 5 sprints, rotate older to git history -->

## SPRINT HISTORY

## Sprint 7: Braille Spinner in Console Header ✅

**Date:** 2026-04-15
**Duration:** ~1.5h

### Agents Participated
- **COUNSELOR** — Discovered existing spinner util in sibling project `tit`; scoped all changes; resolved six iterative ARCHITECT issues (same-color regression, choking spinner, frame-count swap, mid-sprint revert from 24 to 10 frames, two Low audit findings); drove delegation and audit loops
- **Pathfinder** — Four surveys: cake console/build lifecycle, tit spinner implementation, cake palette + async identification, spinner callsite topology
- **Librarian** — Researched cli-spinners catalog, identified dots6/dots7 as the only proven 24-frame braille sets
- **Engineer** — Executed in four delegations: feature scaffold (spinner util copy, OpType enum, SpinnerColor theme field, header render), audit remediation (OpCleanAll/OpRegenerate enums, early-return refactor, magic-string const, whitespace alignment), tick decoupling (SpinnerTickMsg at 80ms), and visible-window render refactor (Pass 1 count / Pass 2 format)
- **Auditor** — Two audit passes: post-scaffold (1 Medium pre-existing early return + 2 Low) and post-window-refactor (1 Low naming, 1 Low magic number, borderline function length)
- **Machinist** — Polish pass: verb-rename `countEntryDisplayLines`, named `consolePanelHorizontalPadding` const

### Files Modified (17 total)
- `internal/ui/spinner.go` — **New.** Classic 10-frame braille set (`⠋⠙⠹⠸⠼⠴⠦⠧⠇⠏`) copied from tit; exports `SpinnerFrames`, `GetSpinnerFrame`, `IsSpinnerFrame`, `SpinnerFrameCount`; `spinnerFrameSet` for O(1) lookup
- `internal/ui/op_type.go` — **New.** `OpType` enum (`OpNone`, `OpBuild`, `OpGenerate`, `OpClean`, `OpCleanAll`, `OpRegenerate`) placed in `ui` package to avoid circular import with `app`
- `internal/ui/theme.go` — Added `SpinnerColor string` to TOML-parsed `ThemeDefinition.Palette`, runtime `Theme`, and `LoadTheme` mapping
- `internal/ui/theme_defaults.go` — `spinnerColor` added to all 5 themes (`#FC704C` GFX/preciousPersimmon, `#FD5B68` Spring/wildWatermelon, `#FF3469` Summer/radicalRed, `#F9C94D` Autumn/saffronMango, `#F6F5FA` Winter/whisper); column alignment matched to neighboring entries
- `internal/ui/console.go` — Added `opLabels` map (BUILDING/CONFIGURING/CLEANING/CLEANING ALL/REGENERATING), `spinnerLabelSeparator`/`fallbackOpLabel`/`minContentHeight`/`consolePanelHorizontalPadding` constants; `assembleConsolePanel` + new `buildConsoleTitle` helper render spinner+label when active, OUTPUT when idle; preserves header width via `lipgloss.Width()`; `RenderConsoleOutput` refactored to positive nested checks (single exit); replaced `formatBufferLines` with Pass 1 `countDisplayLines` + `countEntryDisplayLines` (no rendering) and Pass 2 `formatVisibleLines` + `renderEntry` + `collectVisibleFromEntry` (renders only entries intersecting visible window); `applyScrollState` refactored to take `int totalDisplayLines`, single-exit via `clampScrollOffset` helper
- `internal/ui/ui_test.go` — `TestExtractVisibleWindow` replaced with `TestCountDisplayLines_Empty` and `TestCountDisplayLines_ShortLines`; all `applyScrollState` callsites updated from `[]string` to `int`
- `internal/app/async_state.go` — `currentOp ui.OpType` field; `Start(op ui.OpType)` signature; `End()` resets to `OpNone`; `CurrentOp()` getter
- `internal/app/app.go` — `spinnerFrame int` field (line 64); new `SpinnerTickMsg` case (lines 266-274) with single-exit positive check; `OutputRefreshMsg` case no longer touches spinner
- `internal/app/app_console.go` — `enterConsoleMode(op ui.OpType, footerHint string)` threads op to `asyncState.Start(op)`; resets `spinnerFrame = 0`; new `cmdSpinnerTick()` at `SpinnerTickInterval`; `startAsyncOperation` passes `ui.OpNone`
- `internal/app/app_render.go` — `RenderConsoleOutput` called with `a.asyncState.IsActive()`, `a.spinnerFrame`, `a.asyncState.CurrentOp()`
- `internal/app/constants.go` — Added `SpinnerTickInterval = 80 * time.Millisecond` (decoupled from 100ms `CacheRefreshInterval`)
- `internal/app/messages.go` — New `SpinnerTickMsg struct{}` type
- `internal/app/op_build.go` — `ui.OpBuild` passed to `enterConsoleMode`; `cmdSpinnerTick` added to `tea.Batch`
- `internal/app/op_generate.go` — `ui.OpGenerate` + `cmdSpinnerTick`
- `internal/app/op_clean.go` — `ui.OpClean` + `cmdSpinnerTick`
- `internal/app/op_clean_all.go` — `ui.OpCleanAll` + `cmdSpinnerTick`
- `internal/app/op_regenerate.go` — `ui.OpRegenerate` + `cmdSpinnerTick`

### Alignment Check
- [x] BLESSED principles followed (B: `spinnerFrame` is transient view state on `Application`, reset on `enterConsoleMode`; L: every new/modified function within 30-line limit, `console.go` 288 lines under 300, max 3 branches respected; E: `RenderConsoleOutput` refactored to single-exit positive-nested form, all new code follows same pattern, named constants for padding/labels/intervals, no magic strings; S-SSOT: `opLabels` map is single source for op→label mapping, `SpinnerColor` single-sourced per theme, `countDisplayLines`/`formatVisibleLines` share cheap line-count formula via `countEntryDisplayLines`; S-Stateless: `OpType` lives in `AsyncState.currentOp`, not shadowed; view reads via `CurrentOp()`; E-Encapsulation: `CurrentOp()` getter has a single proven caller (`app_render`), justified; spinner tick and output tick decoupled — each owns its cadence; D: deterministic — same frame index + op yields same render)
- [x] carol/NAMES.md adhered (`countEntryDisplayLines`, `countDisplayLines`, `formatVisibleLines`, `renderEntry`, `collectVisibleFromEntry`, `clampScrollOffset`, `cmdSpinnerTick`, `SpinnerTickMsg`, `SpinnerTickInterval`, `SpinnerColor`, `buildConsoleTitle`, `consolePanelHorizontalPadding`, `fallbackOpLabel`, `spinnerLabelSeparator`, `opLabels`, `currentOp` — verbs for functions, nouns for variables, domain-semantic, no type encoding)
- [x] carol/MANIFESTO.md principles applied (no magic values; all new literals promoted to constants; no early returns in any new/modified code; positive nested checks throughout; label map replaces what would otherwise be a switch at the 3-branch limit; pre-existing early return in `RenderConsoleOutput` resolved in-sprint per no-deferral directive)

### Problems Solved
- **No progress indicator during async ops** — static "OUTPUT" header gave no feedback during multi-second builds/configures. Added braille spinner + op-specific label while `asyncState.IsActive()`, reverts to "OUTPUT" on completion.
- **Theme contrast** — added `SpinnerColor` per-theme field so spinner reads visually distinct from the label; colors picked per theme (preciousPersimmon / wildWatermelon / radicalRed / saffronMango / whisper).
- **Op-type identification** — `AsyncState` previously held no op type, only boolean active/aborted/exitAllowed. Added `currentOp` with enum + getter; 5 callsites thread their op type through `enterConsoleMode`.
- **Spinner choking under heavy output** — initial implementation tied `spinnerFrame++` to `OutputRefreshMsg` (100ms); spinner stuttered because `formatBufferLines` was O(buffer_size) per View(), blocking the render loop under heavy builds. Fixed in two moves: (1) decoupled spinner tick into dedicated `SpinnerTickMsg` at 80ms, (2) refactored console render to O(visible_height) via Pass 1 cheap count (no rendering) + Pass 2 format-only-visible (walks buffer, skips entries before `scrollOffset`, stops at `contentHeight`).
- **Circular import risk** — `OpType` originally envisioned in `internal/app` would create an import cycle with `internal/ui` (which needed the enum for label lookup in `buildConsoleTitle`). Placed `OpType` in `internal/ui` instead — it is a view-facing label concept.
- **Pre-existing early return** — `RenderConsoleOutput` had `if maxWidth <= 0 || totalHeight <= 0 { return "" }` from before this sprint. Refactored in-sprint per "auditor findings never ignored" rule: single-exit, positive nested checks, result initialized to `""` and assigned inside positive branch.
- **Magic string `"WORKING"` fallback** — promoted to `fallbackOpLabel` constant.
- **Label ambiguity for clean-all / regenerate** — initial audit flagged reuse of `OpClean`/`OpGenerate` for these ops; ARCHITECT directed distinct enum values + labels "CLEANING ALL" / "REGENERATING".

### Debts Paid
- None

### Debts Deferred
- None

**Status:** ✅ AUDIT PASS — `go build ./...` clean, `go test ./internal/ui/...` clean. Ready for commit on `main`.

**Note:** On-disk theme files in `~/.config/cake/themes/` will not contain the new `spinnerColor` key until regenerated (ARCHITECT will regen manually per this sprint's discussion). Code regenerates only when theme files are missing.

---

## Sprint 6: Ninja via VS Env + Process-Tree Abort + AsyncState SSOT ✅

**Date:** 2026-04-14
**Duration:** ~3h (two sessions — pre-remote-sync WIP + post-reapply onto new main layout)

### Agents Participated
- **COUNSELOR** — Diagnosed root causes (ninja not on PATH pre-vcvarsall; build abort no-op; ninja/cl.exe orphans on Windows; misleading menu shortcut); scoped each fix; gated execution; drove iterative audit-remediation loop; planned WIP-branch reconciliation after remote-restructure collision; re-mapped WIP onto new main file layout
- **Pathfinder** — Surveyed pre-WIP generator detection, build command assembly, ESC abort handler, stream buffering, menu shortcut bindings, git sync state, then mapped WIP file-for-file onto the remote-restructured main (app.go → app_*.go, project.go merge, menu.go → menu_render.go)
- **Librarian** *(via Engineer)* — Confirmed `golang.org/x/sys/windows` Job Object API names against v0.12.0
- **Engineer** — Executed in 8 delegations: Ninja VS-env detection; two pre-existing LANGUAGE.md §E violations; ninja progress collapse; abort + tree-kill + bufio bundle; audit remediation (error wrapping, single-exit demux, context leak); menu shortcut fix (later confirmed already-on-main); full reapply onto new main structure with accessor introduction; audit-driven remediation (startAsyncOperation helper, accessor consistency, op_regenerate restructure, computeConsoleScrollStatus single-exit, dead branch + import hygiene)
- **Auditor** — Four audit passes: pre-WIP, post-WIP, post-reapply, post-remediation. Surfaced BLESSED-E violations driving each remediation loop (undocumented CloseHandle discards, bare error returns, mid-logic returns, context leak, direct field access bypassing accessors, SSOT duplication, mid-logic returns in regenerate + scroll status, dead branch, import grouping)

### Files Modified (20 total)
- `internal/utils/process_windows.go` — **New.** `StartProcessTree(cmd)` creates Windows Job Object with `JOB_OBJECT_LIMIT_KILL_ON_JOB_CLOSE`, starts command, assigns process to job post-start; `ProcessTree.Close()` idempotent via `jobHandle != 0` guard; every `CloseHandle` discard documented with rollback/teardown reasoning
- `internal/utils/process_stub.go` — **New.** Non-Windows stub: `StartProcessTree` delegates to `cmd.Start()`, `Close()` no-op
- `internal/utils/msvc.go` — Added `IsExecutableInVSEnv(executable, env) bool` scanning captured PATH for the exe
- `internal/utils/msvc_stub.go` — Non-Windows `IsExecutableInVSEnv` stub for API parity
- `internal/utils/stream.go` — `bufio.NewReader` wrapping of stdout/stderr (throughput); new `ninjaProgressPattern` regex + `lastEmittedWasProgress` state routes `\n`-terminated `[N/M]` lines through `replaceCallback` (progress collapse without TTY); `StreamCommand` now returns `(*ProcessTree, error)` and accepts `onProcessTreeStarted` callback; every error return wrapped with `fmt.Errorf("StreamCommand: site: %w", err)`
- `internal/ops/build.go` — Accept `ctx context.Context` (first param) + `onProcessTreeStarted` callback (last); use `exec.CommandContext`; single-exit demux — streamErr branch sets `result.Error`, success branch has `defer tree.Close()` and nested `abortedByUser || waitErr != nil || success` demux; one `return result` at function end; `fmt.Errorf` context on all errors
- `internal/ops/setup.go` — `onProcessTreeStarted` callback param; matching single-exit demux shape; `fmt.Errorf` context; top-of-scope precondition guards preserved
- `internal/state/project.go` — Added `vsEnv []string` field + `SetVSEnv` setter; `DetectAvailableProjects()` moved out of `NewProjectState` constructor into `ForceRefresh` so it runs after `SetVSEnv`; `SetSelectedProject` rewritten with `found` boolean + single exit (removed mid-loop return)
- `internal/state/project_scan.go` — `checkNinjaAvailable()` method: tries `exec.LookPath("ninja")` first, falls back to `utils.IsExecutableInVSEnv("ninja", ps.vsEnv)`; `DetectAvailableProjects` calls it instead of direct `checkCommandExists("ninja")`
- `internal/app/init.go` — Reordered `NewApplication`: `captureVSEnvironment()` → `NewProjectState` → `SetVSEnv(capturedVSEnv)` → `ForceRefresh` (detection now runs with vsEnv populated); `_ = SomeFunc()` discards documented with reason comments
- `internal/app/async_state.go` — Rewrote with accessor API: `Start`, `End`, `Abort`, `ClearAborted`, `IsActive`, `IsAborted`, `CanExit`, `SetExitAllowed`; fields lowercased (`operationActive`, `operationAborted`, `exitAllowed`) to enforce Tell-Don't-Ask; rationale comment on accessor group
- `internal/app/app.go` — Added `killTree func()` field; `BuildCompleteMsg` / `GenerateCompleteMsg` / `RegenerateCompleteMsg` handlers clear `cancelContext` + `killTree` (call if non-nil, then nil) at top of handler before remaining logic; accessor calls (`asyncState.End()`, `IsAborted()`, `ClearAborted()`) replace direct field mutations
- `internal/app/app_console.go` — New `startAsyncOperation(hint)` helper (Start + Clear + footerHint, no mode switch); `enterConsoleMode` now delegates to it then adds `mode = ModeConsole` + `consoleAutoScroll = true`; `isAutoScanIdle` uses `asyncState.IsActive()`
- `internal/app/app_keys.go` — `abortActiveOperation` calls `killTree` (nil-guarded) then `asyncState.Abort()`; ESC and `handleCtrlC` use `asyncState.IsActive()`
- `internal/app/app_render.go` — `asyncState.IsActive()` replaces direct field read (line 54)
- `internal/app/footer.go` — `asyncState.IsActive()` replaces two direct field reads; `computeConsoleScrollStatus` single-exit via `status` local + positive nesting, unreachable else-branch removed, imports reordered to canonical stdlib-then-third-party grouping
- `internal/app/op_build.go` — `ctx, cancel := context.WithCancel(background)`, store `a.cancelContext = cancel`, pass ctx to `ExecuteBuildProject`, closure sets `a.killTree = tree.Close` via callback
- `internal/app/op_generate.go` — Same pattern for generate; `killTree` wired into `ExecuteSetupProject` call
- `internal/app/op_regenerate.go` — `killTree` wired into `ExecuteSetupProject`; clean-step restructured to positive nesting with `cleanSucceeded` bool + `cleanErr`; single return at function end; `fmt.Errorf("cmdRegenerateProject: remove build dir: %w", err)` wraps the remove error
- `internal/app/op_open.go` — `startOpenIDEOperation` and `startOpenEditorOperation` delegate to `a.startAsyncOperation(hint)` (SSOT — no inlined state setup); IDE/editor launches deliberately do NOT switch to console mode
- `go.mod` — Promoted `golang.org/x/sys v0.12.0` from indirect to direct

### Alignment Check
- [x] BLESSED principles followed (B: Job Object binds process-tree lifetime to job handle, RAII-style via idempotent Close; L: functions within 30-line limit; E: single-exit demux, positive nesting, `fmt.Errorf` wrapping, documented error discards, no magic values; S-SSOT: `startAsyncOperation` helper eliminates duplicated state setup across op_open / enterConsoleMode; S-Stateless: `vsEnv` is a captured input, not mutable machinery state; E-Encapsulation: `AsyncState` accessor API prevents bypass; D: emerges from B+E+S)
- [x] carol/NAMES.md adhered (`SetVSEnv`, `checkNinjaAvailable`, `IsExecutableInVSEnv`, `ninjaProgressPattern`, `matchesNinjaProgress`, `shouldReplace`, `lastEmittedWasProgress`, `ProcessTree`, `StartProcessTree`, `startAsyncOperation`, `cleanSucceeded`, `cleanErr` — full words, verb-noun where applicable, domain-specific)
- [x] carol/MANIFESTO.md principles applied (no magic values; guard-topology respected — returns only at top of scope or function end; no shadow state; no dead code — `runningCmd` field removed; `(can scroll up)` dead branch removed; three-branch limit observed throughout)
- [x] carol/LANGUAGE.md Go overrides respected (`fmt.Errorf("site: action: %w", err)` wrapping on all error paths; every `_ = Func()` has reason comment; guard-topology — returns either top-of-scope precondition/error guards or single exit at function end; mid-logic returns eliminated)
- [x] Accessor consistency — grep across `internal/app/` confirms `operationActive|operationAborted|exitAllowed` appear only inside `async_state.go`

### Problems Solved
- **Ninja invisible on Windows** — `ninja.exe` ships only inside VS install, discoverable only after `vcvarsall.bat`. Detection ran before env capture. Fixed by reordering startup + VS-env fallback scan via `IsExecutableInVSEnv`.
- **Build abort was a no-op** — `runningCmd` field never assigned, `ExecuteBuildProject` accepted no context. Fixed by context propagation end-to-end (startOp → closure → ctx → `exec.CommandContext`).
- **ninja/cl.exe orphans on abort** — `Process.Kill` on Windows kills only the direct child. Fixed by Windows Job Object with `JOB_OBJECT_LIMIT_KILL_ON_JOB_CLOSE`.
- **Output throughput bottleneck** — byte-per-syscall reads on busy build output. Fixed by `bufio.NewReader` wrap.
- **Ninja progress scrolled line-by-line** — when stdout is not a TTY, ninja emits `\n`-terminated lines, no `\r` overwrites. Pattern-based detection routes `[N/M]` progress through `ReplaceLast`.
- **Misleading menu shortcut** — `k` displayed for Clean but bound to up-nav. Already resolved on remote main as `c`/`x`; verified no-op.
- **Context leak on normal completion** — `cancelContext` never called on happy path; next op overwrote silently. Fixed in completion message handlers.
- **Accessor bypass in async state** — 5 sites across app_render/footer/op_open read/wrote `operationActive`/`operationAborted` directly, bypassing the declared `AsyncState` API. Fixed by introducing accessors and replacing all direct access.
- **SSOT duplication** — op_open inlined Start + Clear + footerHint independent of `enterConsoleMode`. Fixed by extracting `startAsyncOperation(hint)` helper; `enterConsoleMode` wraps it with mode + autoScroll.
- **Mid-logic returns** — eliminated in `SetSelectedProject` (pre-existing), `ExecuteBuildProject` / `ExecuteSetupProject` post-Wait demux, `streamErr` demux (defer scoped to success branch), `op_regenerate` clean step, `computeConsoleScrollStatus`.
- **Bare error returns / silent discards** — every error return now wrapped via `fmt.Errorf %w`; every `_ = Func()` documented with reason.
- **Dead code** — `runningCmd` field + `os/exec` import removed from app.go; unreachable `(can scroll up)` branch removed from footer scroll status; import grouping normalized.
- **Remote restructure reconciliation** — after remote decomposed `app.go` → `app_*.go`, `state/project.go` → `project_paths.go`/`project_scan.go` (later re-merged back), etc., WIP was committed to `wip-ninja` branch, main was hard-reset to `origin/main`, then WIP was re-applied file-by-file onto the new layout via Pathfinder mapping table. Zero content loss.

### Technical Debt / Follow-up
- None — all findings resolved in-sprint per ARCHITECT's no-deferral directive. Zero Auditor findings outstanding.

**Status:** ✅ AUDIT PASS — `go build ./...` clean. Ready for commit on `main`.

---

## Sprint 5: Open Editor for Ninja + Release Polish

**Date:** 2026-04-05

### Agents Participated
- **COUNSELOR** — Plan, coordination, goreleaser fix
- **Pathfinder** — Code discovery for nvim wiring (state, ops, menu, actions)
- **Engineer** — Implemented Open IDE/Editor dispatch, doc updates (SPEC, ARCHITECTURE, README)
- **Machinist** — Clean sweep of audit findings (Sprint 4 carryover)

### Files Modified (10 total)
- `internal/state/project.go` — `CanOpenIDE()` simplified to `len(ps.AvailableProjects) > 0`; `CanOpenEditor()` deleted (dead code)
- `internal/ui/menu.go` — `GenerateMenuRows` adds `isIDEGenerator bool` param; `openIde` row label/hint dynamic via `openIdeLabel()`/`openIdeHint()` helpers
- `internal/app/menu.go` — Passes `utils.IsGeneratorIDE()` to `GenerateMenuRows`
- `internal/app/app_actions.go` — `"openIde"` case dispatches to `startOpenIDEOperation()` or `startOpenEditorOperation()` based on generator type
- `internal/ui/ui_test.go` — Updated `GenerateMenuRows` call sites with new `isIDEGenerator` param
- `SPEC.md` — Open IDE/Editor feature documented (dynamic label, both IDE and CLI generators), selectability rules updated, corrected to match codebase
- `ARCHITECTURE.md` — `CanOpenIDE()` contract updated, Pattern 4 annotated with dynamic label behavior
- `README.md` — `o` shortcut updated to "Open IDE / Editor", workflow section updated
- `.goreleaser.yaml` — `format: binary` (deprecated) to `formats: [binary]`
- `release.sh` — `gh release delete --cleanup-tag` replaces manual tag deletion

### Alignment Check
- [x] BLESSED principles followed
- [x] carol/NAMES.md adhered
- [x] carol/MANIFESTO.md principles applied
- [x] No new dead code (CanOpenEditor removed, startOpenEditorOperation now wired)

### Problems Solved
- Ninja generator had no "open" action — now opens nvim via same `openIde` row with dynamic label
- goreleaser `archives.format` deprecation warning — migrated to `formats` (plural, list)
- Release re-run left orphan GitHub releases — `gh release delete --cleanup-tag` cleans up both release and tag

### Technical Debt / Follow-up
- None

**Status:** Complete, ready for v0.0.2 release

---

## Sprint 4: v0.0.1 Release and Documentation Audit

**Date:** 2026-04-05

### Agents Participated
- **COUNSELOR** — Release infrastructure planning, goreleaser debugging, doc audit coordination
- **Pathfinder** — Git remotes, doc inventory, release infrastructure discovery
- **Researcher** — goreleaser OSS capabilities (signing hooks, cross-compilation, token auth)
- **Librarian** — goreleaser OSS vs PRO features, gh CLI token bridging
- **Engineer** (x5) — README, SPEC, ARCHITECTURE, SPRINT-LOG updates, post-build script, token file setup
- **Auditor** (x3) — Release readiness, test suite verification, comprehensive docs audit
- **Machinist** — Clean sweep of all audit findings across SPEC, ARCHITECTURE, SPRINT-LOG

### Files Modified (12 total)
- `.goreleaser.yaml` — goreleaser v2 config with post-build hook (replaced PRO-only `if` filters and `release_notes` field)
- `release.sh` — One-command release: commit, tag, push, gh auth token bridge, release notes via CLI flag
- `scripts/post-build.sh` — macOS sign+notarize wrapper (concurrent-safe with mktemp -d)
- `RELEASE_NOTES.md` — GitHub release description for v0.0.1
- `.gitignore` — Added `/dist/` (goreleaser output)
- `internal/constants.go` — AppVersion changed from const to var, default "dev", ldflags injectable
- `build.sh` — Added git describe version detection and ldflags injection
- `internal/app/messages.go` — Footer hint: added `[x] Clean All`
- `internal/ui/menu.go` — Clean shortcut `k` to `c`, Clean All `ctrl+k` to `x`
- `internal/app/app_handlers.go` — menuShortcutMap: `ctrl+k` to `x`/`X`
- `README.md` — Public release install instructions, shortcuts, removed TIT reference
- `SPEC.md` — v0.0.2, selectability model, corrected menu labels/shortcuts/interval/config format, removed Open Editor as separate row
- `ARCHITECTURE.md` — Removed dead CanOpenEditor() from interface contract
- `carol/SPRINT-LOG.md` — Fixed obsolete doc references (LIFESTAR to BLESSED, NAMING-CONVENTION to NAMES, ARCHITECTURAL-MANIFESTO to MANIFESTO)

### Alignment Check
- [x] BLESSED principles followed
- [x] carol/NAMES.md adhered
- [x] carol/MANIFESTO.md principles applied
- [x] SSOT: SPEC.md corrected to match codebase (codebase is SSOT)
- [x] All audit findings resolved (0 deferred)

### Problems Solved
- goreleaser OSS does not support `if` on build hooks — solved with wrapper script checking $HOOK_TARGET
- goreleaser OSS does not support `release_notes` YAML field — solved with `--release-notes` CLI flag
- notarytool requires zip not bare binary — solved with ditto zip, submit, cleanup
- Concurrent hooks collided on temp file — solved with mktemp -d (unique dir per invocation)
- dist/ committed by git add -A — solved with .gitignore entry
- GITHUB_TOKEN missing — solved with gh auth token bridge
- SPEC.md diverged from implementation (menu labels, shortcuts, visibility model, interval control, config format) — corrected all

### Technical Debt / Follow-up
- Dead code: CanOpenEditor(), startOpenEditorOperation(), cmdOpenEditor() — will be wired in Sprint 5 (nvim as Ninja IDE)
- `archives.format: binary` deprecation warning from goreleaser v2 — cosmetic, works fine
- TIT release RFC written at ~/Documents/Poems/dev/tit/RFC.md — not yet executed

**Status:** v0.0.1 released, signed, notarized, all docs audit-clean

---

## Sprint 3: Production Refactoring, Test Coverage, and Release Infrastructure

**Date:** 2026-04-05

### Agents Participated
- **COUNSELOR** — Led all phases: refactoring plan, test strategy, release infrastructure planning, doc cleanup
- **Pathfinder** — Codebase discovery (testable functions, release infrastructure, doc inventory, git remotes)
- **Engineer** (x8) — Module path migration, test suites (4 packages), doc updates (README, SPEC, ARCHITECTURE, SPRINT-LOG)
- **Auditor** (x3) — Release readiness audit, full test suite verification, final doc verification
- **Researcher** — goreleaser research (signing hooks, cross-compilation, local release workflow)

### Files Modified (88 files changed, 1741 insertions, 15180 deletions)

**Release Infrastructure (new)**
- `.goreleaser.yaml` — goreleaser v2 config: 6 targets (darwin/linux/windows x amd64/arm64), macOS codesign+notarize hooks, version injection via ldflags
- `release.sh` — One-command release: commit, tag, push, goreleaser
- `entitlements.plist` — macOS codesign entitlements (allow-unsigned-executable-memory, disable-library-validation)

**Module Path Migration**
- `go.mod` — `cake` to `github.com/jrengmusic/cake`
- 26 `.go` files — all internal imports updated to new module path

**Version Injection**
- `internal/constants.go` — `AppVersion` changed from `const` to `var`, default `"dev"`, injected via `-ldflags -X`
- `build.sh` — Added `git describe --tags` version detection and ldflags injection

**Test Coverage (121 tests, all pass)**
- `internal/banner/banner_test.go` — 18 tests: RGBToHex, CanvasToBrailleArray, SvgToBrailleArray, parseColor, normalizePathData, dominantPixelColor
- `internal/config/config_test.go` — 14 tests: DefaultConfig, accessors, LastConfiguration fallback
- `internal/state/state_test.go` — 39 tests: cycling, configuration, build info, paths, parseProjectCallName
- `internal/ui/ui_test.go` — 50 tests: sizing, menu generation, OutputBuffer (incl. concurrent), ConsoleOutState, ConfirmationDialog, scroll helpers

**Shortcut Fixes**
- `internal/ui/menu.go` — Clean shortcut: `k` to `c`, Clean All: `ctrl+k` to `x`
- `internal/app/app_handlers.go` — menuShortcutMap: `ctrl+k` to `x`/`X` for cleanAll

**Documentation**
- `README.md` — Public release install instructions (`go install`), `x` shortcut in nav table
- `SPEC.md` — Version v0.0.1, added `c`/`x` shortcuts, 65/35 layout split
- `ARCHITECTURE.md` — Full rewrite matching current codebase (decomposed files, context cancellation, handler maps, atomic buffer reads)
- `carol/SPRINT-LOG.md` — Cleaned boilerplate, removed template sprint, updated repo URL
- `REFACTORING-PLAN.md` — Deleted (completed work, lives in git history)

**Prior Refactoring (from earlier in session, pre-compaction)**
- Bug fixes: context-based build cancellation, atomic OutputBuffer.GetSnapshot()
- Dead code removal: 4 files deleted, 15+ dead functions/methods removed
- SSOT constants extracted to `internal/constants.go`
- File decomposition: app.go (1117 to 282 lines), project.go, theme.go, menu.go, svg.go all split
- Handler maps: menuShortcutMap, consoleLineColorMap, executePendingOperation
- 60+ helper methods extracted to meet 30-line function limit

### Alignment Check
- [x] BLESSED principles followed
- [x] NAMES.md adhered
- [x] MANIFESTO.md principles applied
- [x] All files under 300 lines
- [x] All functions under 30 lines (TEA Update exception accepted)
- [x] No switches over 3 branches (type dispatches and data declarations accepted)

### Problems Solved
- Zero test coverage to 121 tests across 4 packages
- No release infrastructure to one-command `bash release.sh v0.0.1 "msg"` with cross-compilation and macOS signing
- Bare module path `cake` blocked public `go install` — migrated to `github.com/jrengmusic/cake`
- Hardcoded version string — now injected at build time via ldflags
- Stale ARCHITECTURE.md referenced deleted files and dead methods — full rewrite
- Clean shortcut `k` conflicted with vi navigation — changed to `c`
- Clean All `ctrl+k` awkward — changed to `x`

### Technical Debt / Follow-up
- No tests for `internal/app` (bubbletea model — no viable unit test path) or `internal/ops` (subprocess-dependent)
- `extractVisibleWindow` has minor edge case with negative offsets (not reachable in practice, documented in test)
- TIT RFC written at `~/Documents/Poems/dev/tit/RFC.md` for identical release setup

**Status:** Ready for release as v0.0.1

---

## Sprint 2: MSVC Environment Detection + Console Output Streaming

**Date:** 2026-03-28

### Agents Participated
- **COUNSELOR** — Diagnosed root cause (no cmake on PATH in MSYS2), planned MSVC detection and console streaming fixes
- **Pathfinder** — Explored cake/TIT codebases, verified vswhere paths, checked VS installation structure
- **Librarian** — Researched vswhere.exe usage, vcvarsall.bat env capture patterns in Go
- **Engineer** — Implemented MSVC detection, console streaming rewrite, dead code cleanup
- **Auditor** — Validated all changes (2 passes: MSVC detection + console streaming)

### Files Modified (15 total)
- `internal/utils/msvc.go:1-112` — New: vswhere detection, vcvarsall env capture, cmake path resolution (windows build tag)
- `internal/utils/msvc_stub.go:1-20` — New: non-windows stubs for cross-compilation
- `internal/utils/stream.go:1-106` — Rewrite: byte-by-byte pipe reading, \r progress handling, WaitGroup drain
- `internal/utils/exec.go` — Deleted: unused stub
- `internal/utils/generator.go` — Deleted: dead code, wrong path scheme
- `internal/utils/platform.go` — Deleted: dead code, unused by app
- `internal/ui/buffer.go:70-92` — Added ReplaceLast() for progress line updates
- `internal/app/constants.go:1-8` — New: CacheRefreshInterval constant
- `internal/app/app.go:64,1108-1116` — Added vsEnv field, outputCallbacks() helper, use CacheRefreshInterval
- `internal/app/init.go:36-39` — VS env capture at startup
- `internal/ops/setup.go:17` — Dual callback signature, cmake path resolution from vsEnv
- `internal/ops/build.go:16` — Dual callback signature, cmake path resolution from vsEnv
- `internal/app/op_generate.go:27-41` — Uses outputCallbacks(), passes vsEnv
- `internal/app/op_build.go:21-35` — Uses outputCallbacks(), passes vsEnv
- `internal/app/op_regenerate.go:30-44` — Uses outputCallbacks(), passes vsEnv
- `internal/app/op_clean.go:21` — Uses outputCallbacks() (append only)
- `internal/app/op_clean_all.go:25` — Uses outputCallbacks() (append only)
- `internal/app/op_open.go:20,52` — Uses outputCallbacks() (append only)
- `internal/state/project.go:130-150` — VS detection via vswhere, version-aware (only installed versions)

### Alignment Check
- [x] BLESSED principles followed
- [x] carol/NAMES.md adhered (versionToGenerator, appendCallback, replaceCallback, isProgressLine)
- [x] carol/MANIFESTO.md principles applied (SSOT: outputCallbacks helper, Lean: streamPipe extracted, Explicit: clear callback signatures)
- [ ] No early returns — Go exempt per ARCHITECT decision (Go idiom accepted)

### Problems Solved
- cake could not invoke cmake on Windows/MSYS2 (cmake not on PATH)
- VS generator detection was version-blind (added both VS2022+VS2026 regardless of what's installed)
- VS install path uses version number ("18") not marketing year ("2026") — fixed detection
- Go's exec.Command escapes quotes with backslashes, cmd.exe doesn't understand — fixed with SysProcAttr.CmdLine
- Go's exec.LookPath resolves executable against process PATH, not cmd.Env — fixed with FindExecutableInEnv
- Console output race condition: pipes not drained before cmd.Wait() — fixed with WaitGroup
- No \r progress line handling — added byte-by-byte reading matching TIT pattern

### Technical Debt / Follow-up
- CaptureVSEnv error silently discarded in init.go — acceptable for startup (no UI available yet) but could log to file
- Debug file C:\Users\jreng\cake-debug.txt may still exist on disk — delete manually
- Ninja generator on Windows also needs vcvarsall env (for cl.exe) — currently works because vsEnv is passed to all cmake calls regardless of generator

**Status:** APPROVED - VS 2026 detected, cmake generates successfully from cake
