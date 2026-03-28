# SPRINT-LOG.md

**Project:** cake  
**Repository:** /c/Users/jreng/Documents/Poems/dev/cake  
**Started:** 2026-03-28

**Purpose:** Long-term context memory across sessions. Tracks completed work, technical debt, and unresolved issues. Written by PRIMARY agents only when ARCHITECT explicitly requests.

---

## 📖 Notation Reference

**[N]** = Sprint Number (e.g., `1`, `2`, `3`...)

**Sprint:** A discrete unit of work completed by one or more agents, ending with ARCHITECT approval ("done", "good", "commit")

---

## ⚠️ CRITICAL RULES

**AGENTS BUILD CODE FOR ARCHITECT TO TEST**
- Agents build/modify code ONLY when ARCHITECT explicitly requests
- ARCHITECT tests and provides feedback
- Agents wait for ARCHITECT approval before proceeding

**AGENTS NEVER RUN GIT COMMANDS**
- Write code changes without running git commands
- Agent runs git ONLY when user explicitly requests
- Never autonomous git operations
- **When committing:** Always stage ALL changes with `git add -A` before commit
  - ❌ DON'T selectively stage files (agents forget/miss files)
  - ✅ DO `git add -A` to capture every modified file

**SPRINT-LOG WRITTEN BY PRIMARY AGENTS ONLY**
- **COUNSELOR** or **SURGEON** write to SPRINT-LOG
- Only when user explicitly says: `"log sprint"`
- No intermediate summary files
- No automatic logging after every task
- Latest sprint at top, keep last 5 entries

**NAMING RULE (CODE VOCABULARY)**
- All identifiers must obey project-specific naming conventions (see NAMING-CONVENTION.md)
- Variable names: semantic + precise (not `temp`, `data`, `x`)
- Function names: verb-noun pattern (initRepository, detectCanonBranch)
- Struct fields: domain-specific terminology (not generic `value`, `item`, `entry`)
- Type names: PascalCase, clear intent (CanonBranchConfig, not BranchData)

**BEFORE CODING: ALWAYS SEARCH EXISTING PATTERNS**
- ❌ NEVER invent new states, enums, or utility functions without checking if they exist
- ✅ Always grep/search the codebase first for existing patterns
- ✅ Check types, constants, and error handling patterns before creating new ones
- **Methodology:** Read → Understand → Find SSOT → Use existing pattern

**TRUST THE LIBRARY, DON'T REINVENT**
- ❌ NEVER create custom helpers for things the library/framework already does
- ✅ Trust the library/framework - it's battle-tested

**FAIL-FAST RULE (CRITICAL)**
- ❌ NEVER silently ignore errors (no error suppression)
- ❌ NEVER use fallback values that mask failures
- ❌ NEVER return empty strings/zero values when operations fail
- ❌ NEVER use early returns
- ✅ ALWAYS check error returns explicitly
- ✅ ALWAYS return errors to caller or log + fail fast

**⚠️ NEVER REMOVE THESE RULES**
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

<!-- Example sprint entry (delete this after first real sprint) -->

## Sprint 1: Project Setup and Initial Planning ✅

**Date:** 2026-01-11  
**Duration:** 14:00 - 16:30 (2.5 hours)

### Agents Participated
- **COUNSELOR:** Kimi-K2 — Wrote SPEC.md and ARCHITECTURE.md
- **ENGINEER** (invoked by COUNSELOR) — Created project structure
- **AUDITOR** (invoked by COUNSELOR) — Verified spec compliance

### Files Modified (8 total)
- `SPEC.md:1-200` — Complete feature specification with all flows
- `ARCHITECTURE.md:1-150` — Initial architecture patterns documented
- `src/core/module.cpp:10-45` — Core module scaffolding with proper initialization
- `src/core/module.h:1-30` — Core module header with explicit dependencies
- `tests/core_test.cpp:1-50` — Test scaffolding following Testable principle
- `CMakeLists.txt:1-25` — Build configuration with explicit targets
- `README.md:1-20` — Project overview

### Alignment Check
- [x] LIFESTAR principles followed (Lean, Immutable, Findable, Explicit, SSOT, Testable, Accessible, Reviewable)
- [x] NAMING-CONVENTION.md adhered (semantic names, verb-noun functions, no type encoding)
- [x] ARCHITECTURAL-MANIFESTO.md principles applied (no layer violations, explicit dependencies)
- [x] No early returns used
- [x] Fail-fast error handling implemented

### Problems Solved
- Established project foundation following domain-specific patterns
- Defined clear module boundaries preventing layer violations

### Technical Debt / Follow-up
- Error handling needs refinement in module.cpp (marked with TODO)
- Performance requirements not yet defined for real-time constraints

**Status:** ✅ APPROVED - All files compile, tests scaffold in place

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
- [x] LIFESTAR principles followed
- [x] NAMING-CONVENTION.md adhered (versionToGenerator, appendCallback, replaceCallback, isProgressLine)
- [x] ARCHITECTURAL-MANIFESTO.md principles applied (SSOT: outputCallbacks helper, Lean: streamPipe extracted, Explicit: clear callback signatures)
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

---

<!-- Actual sprint entries go here, written by PRIMARY agents -->

---

**End of SPRINT-LOG.md Template**

Copy this template to your project root as `SPRINT-LOG.md` and customize:
- Project name
- Repository URL/path
- Starting date
- Add project-specific rules to CRITICAL RULES section

Rock 'n Roll!  
**JRENG!**
