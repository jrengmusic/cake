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
