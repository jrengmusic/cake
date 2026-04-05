# LANGUAGE.md Handoff for CAKE
## Go/Bubbletea BLESSED Compliance

**From:** COUNSELOR (TIT project)
**Date:** 2026-04-05
**For:** COUNSELOR (CAKE project)

---

## Context

ARCHITECT formalized `~/.carol/LANGUAGE.md` (v0.1, April 2026) — a multi-language addendum to MANIFESTO.md that defines how BLESSED principles map to each language's idioms. This was written after a full production audit of TIT (Go/bubbletea/lipgloss) revealed systemic friction between BLESSED (designed for C++/JUCE) and Go's design constraints.

CAKE shares the same stack (Go/bubbletea/lipgloss) and the same architectural patterns (single-model TEA, dispatcher maps, state clusters). Everything below applies directly.

**Read `~/.carol/LANGUAGE.md` first.** This handoff explains the rationale and what to look for — LANGUAGE.md is the contract.

---

## The Three Go-BLESSED Adaptations

### 1. Early Returns — Permitted for Error Guards Only

MANIFESTO forbids early returns. Go's error handling is built on them.

**Resolution:** LANGUAGE.md permits early returns for error handling (`if err != nil { return }`) and precondition guards (`if input == nil { return }`). Returns inside business logic remain violations.

**What this means for CAKE:**
- `if err != nil { return ..., fmt.Errorf("context: %w", err) }` — compliant
- Bare `return err` without context wrapping — violation (E: Explicit)
- `if score > threshold { return result, nil }` inside business logic — violation
- Guard returns at the top of a function, happy path below — compliant pattern

**Audit action:** Triage all return statements. Classify as error-guard (keep) or business-logic (refactor). Do not blindly refactor all returns.

### 2. Encapsulation — Package Boundary, Not Struct Boundary

MANIFESTO says private by default, no poking internals. Go has no struct-private — visibility is package-scoped.

**Resolution:** LANGUAGE.md redefines the encapsulation boundary as the package. Within a package, direct field access IS the API. Accessors (getters/setters) used only within the same package are dead code.

**What this means for CAKE:**
- If `app.Application` has state structs with accessor methods only called from within `internal/app/` — those accessors are dead code. Delete them.
- Direct field access within the same package is correct, not a violation.
- Accessors matter only at package boundaries — e.g., `internal/state/` exporting methods for `internal/app/` to call.
- CAKE's layer separation (app / state / config / ui / ops) already creates real encapsulation boundaries. Fields exported across these packages should use proper exported methods. Fields accessed within the same package should be accessed directly.

**Audit action:** Check for same-package accessors. Delete them. Verify cross-package boundaries use exported methods correctly.

### 3. Bubbletea God Object — Accepted Framework Constraint

MANIFESTO's L principle forbids god objects. Bubbletea's Elm Architecture mandates a single root Model.

**Resolution:** LANGUAGE.md declares the single root model an accepted framework constraint. It is mitigated, not eliminated.

**Mitigation contracts (enforce these):**
- **State clusters** — decompose root model into embedded state structs by domain. Each cluster owns its domain. Audit that state is not scattered across the root struct.
- **File decomposition** — 300-line limit applies per file. The logical model spans files, but each file owns one concern.
- **Handler maps over switch chains** — action dispatch uses `map[Type]Handler`, not `switch`. Adding a case is data, not code. Respects 3-branch limit and SSOT.
- **View is pure** — `View()` reads state, produces string, mutates nothing.
- **`tea.Cmd` is the only side effect channel** — no goroutines spawned from `Update()` directly.

**What this means for CAKE:**
- `app.Application` being large is not a violation if state is clustered and files are decomposed.
- `Update()` dispatching to many handlers across files is the accepted pattern.
- Check that no `switch` chain exceeds 3 branches — use handler maps instead.
- Check that `View()` has no side effects.

---

## BLESSED Principles Unchanged for Go

These apply exactly as MANIFESTO.md states — no override:

| Principle | Note for Go |
|-----------|-------------|
| **L (Lean)** | 300/30/3 applies. Go culture agrees. |
| **S (SSOT)** | No excuse. No shadow state, no duplicated constants, no duplicated logic. |
| **S (Stateless)** | Objects are dumb workers. Bubbletea's MVU naturally enforces this — View is pure, Update is the only mutation point. |
| **D (Deterministic)** | Emergent. Same input, same output. Tests prove it. |
| **B (Bound)** | Adapted for GC — ownership via convention, `defer` for cleanup, `context.Context` for goroutine lifecycle. No goroutine without cancellation. |

---

## Go-Specific Anti-Patterns to Audit

| Anti-Pattern | Violation | What to look for |
|---|---|---|
| Bare `return err` | E (Explicit) | Error returns without `fmt.Errorf` context wrapping |
| `_ = Func()` undocumented | E (Explicit) | Silently discarded errors without comment explaining why |
| Goroutine without context | B (Bound) | `go func()` without cancellation path |
| Same-package accessor | E (Encapsulation) | Getter/setter only called within its own package |
| `interface{}` / `any` | E (Explicit) | Loses type safety where concrete type exists |
| Package-level `var` mutable | S (Stateless) + B | Global mutable state outside the model |
| `init()` with side effects | E (Explicit) | Hidden initialization |
| Magic numbers/strings | E + S (SSOT) | Hardcoded values that appear more than once |
| Map rebuilt per call | S (SSOT) + perf | Maps that should be package-level vars |

---

## Lessons from TIT Audit

What we found in TIT that CAKE likely shares (same stack, same patterns):

1. **Dead accessors** — state clusters had getters/setters used only within `package app`. Pure ceremony. Delete them.
2. **Dual access paths** — same state reachable via accessor AND direct field access. Pick one (direct, per LANGUAGE.md).
3. **Magic numbers** — timing constants, scroll sizes hardcoded in multiple places. Extract to named constants.
4. **Console transition duplication** — the pattern of entering console mode (clear buffer, set autoscroll, switch mode, set footer) repeated 15+ times. Extract into a single function.
5. **Handler maps defined twice** — if shortcuts or handlers are registered in multiple places, consolidate to SSOT.
6. **Global singleton ambiguity** — if a global (like an output buffer) is wrapped by a state struct AND accessed directly, pick one path.

---

## How to Use This

1. Read `~/.carol/LANGUAGE.md` — it is the contract
2. Run a production audit against LANGUAGE.md (not raw MANIFESTO.md)
3. Use the anti-patterns table as a checklist
4. When in doubt: LANGUAGE.md overrides MANIFESTO.md examples, but not principles

---

**JRENG!**
