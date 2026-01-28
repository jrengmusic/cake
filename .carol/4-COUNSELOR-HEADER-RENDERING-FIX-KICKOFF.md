# Sprint 4 - COUNSELOR - Header Rendering Fix - KICKOFF

**Role:** COUNSELOR
**Agent:** OpenCode (glm-4.7)
**Date:** 2026-01-28

## Objective

Fix inconsistent header rendering between menu mode (2 lines) and console mode (3 lines) by removing double placement pattern in `cake/internal/ui/layout.go` to match TIT's architecture exactly.

## Root Cause Analysis

**Problem:** Header renders inconsistently across modes
- Menu mode: Shows only 2 lines (CWD + separator), missing project name line
- Console mode: Correctly shows 3 lines (project name + CWD + separator)

**Root Cause:** Double placement pattern in `cake/internal/ui/layout.go`

Current CAKE implementation has a nested placement wrapper that doesn't exist in TIT:

```go
// CAKE (WRONG) - lines 83-92
combined := lipgloss.JoinVertical(lipgloss.Left, headerSection, contentSection, footerSection)

// DOUBLE PLACEMENT - wraps already-placed content AGAIN
return lipgloss.Place(
    sizing.TerminalWidth,
    sizing.TerminalHeight,
    lipgloss.Left,
    lipgloss.Bottom,
    combined,
)
```

```go
// TIT (CORRECT) - line 102
// Join sections vertically - no wrapping Place
return lipgloss.JoinVertical(lipgloss.Left, headerSection, contentSection, footerSection)
```

**Why this causes the issue:**
1. Outer `lipgloss.Place()` re-evaluates entire layout constraints
2. When content is short (menu mode), lipgloss may clip header to fit layout
3. Console mode has full-height content, so header renders correctly
4. TIT's simple join pattern avoids this re-layout behavior

## Design Decisions

### 1. Match TIT Architecture Exactly
- Remove outer `lipgloss.Place()` wrapper from `RenderReactiveLayout()`
- Keep header/footer `lipgloss.Place()` for fixed positioning
- Return joined sections directly (no outer wrapper)

### 2. No API Changes
- All existing functions keep same signatures
- No changes to header construction logic
- No changes to sizing calculations

### 3. Reference: TIT/internal/ui/layout.go
Compare lines 64-103 in TIT with CAKE's implementation
TIT uses simple `lipgloss.JoinVertical()` without outer placement

## Implementation Plan

### Phase 1: Modify RenderReactiveLayout (Primary Fix)

**File:** `cake/internal/ui/layout.go`

**Changes:**
1. Remove lines 86-92 (double placement wrapper)
2. Replace with single return matching TIT pattern

**Before (lines 83-92):**
```go
	// Join sections
	combined := lipgloss.JoinVertical(lipgloss.Left, headerSection, contentSection, footerSection)

	// Place in exact terminal dimensions
	return lipgloss.Place(
		sizing.TerminalWidth,
		sizing.TerminalHeight,
		lipgloss.Left,
		lipgloss.Bottom,
		combined,
	)
```

**After:**
```go
	// Join sections vertically - no wrapping Place (TIT pattern)
	return lipgloss.JoinVertical(lipgloss.Left, headerSection, contentSection, footerSection)
```

### Phase 2: Verify RenderConfirmDialog (Secondary Check)

**File:** `cake/internal/ui/layout.go`

**Check:** `RenderConfirmDialog()` function (lines 112-183)

This function may have similar double placement pattern. Verify if it needs the same fix.

**If yes:** Remove outer placement wrapper, match TIT pattern
**If no:** No changes needed (this function works correctly)

### Phase 3: Verification

**Build and Test:**
1. Run `./build.sh` to compile CAKE
2. Run CAKE and switch between modes
3. Verify header shows 3 lines in both menu and console modes:
   - Line 1: Project name (uppercase, 4-space indent)
   - Line 2: CWD with folder emoji
   - Line 3: Separator

**Expected Result:**
- Menu mode header: 3 lines (fixed - matches console)
- Console mode header: 3 lines (unchanged)
- Preferences mode header: 3 lines (should also fix)
- All headers render identically across all modes

## Acceptance Criteria

1. ✅ `RenderReactiveLayout()` matches TIT pattern (no double placement)
2. ✅ Header renders 3 lines in menu mode (project name visible)
3. ✅ Header renders 3 lines in console mode (unchanged)
4. ✅ Header renders 3 lines in preferences mode (also fixed)
5. ✅ Build succeeds with no errors
6. ✅ No visual artifacts or layout shifts after fix

## Edge Cases

**Terminal Resizing:**
- Test header remains 3 lines at different terminal sizes
- Verify header doesn't clip at minimum dimensions (69×19)

**Empty Project Name:**
- If `GetProjectName()` returns empty string, header should show directory name fallback
- Verify fallback rendering works correctly

**Long Project Names:**
- Verify project name truncation doesn't affect other header lines
- Ensure separator remains full width

## Integration Points

**No integration changes needed:**
- `RenderReactiveLayout()` is UI rendering function only
- No changes to state management or business logic
- All callers in `app.go` work the same way

**Files Modified:**
- `cake/internal/ui/layout.go` - Remove double placement wrapper

## Dependencies

None - this is a pure UI rendering fix with no external dependencies

## Testing Instructions

1. Build: `./build.sh`
2. Run CAKE
3. Observe header in default menu mode:
   - Should see: "    PROJECTNAME" (uppercase, 4-space indent)
   - Should see: "📂  /path/to/directory"
   - Should see: "─────────────────────"
4. Press any action key (e.g., Enter) to switch to console mode
5. Verify console header shows same 3 lines
6. Press "/" to open preferences, verify preferences header shows 3 lines
7. Test terminal resize: resize window, verify all headers remain 3 lines

## Notes

**Why this fix works:**
- Removing outer `lipgloss.Place()` eliminates re-layout behavior
- Simple vertical join guarantees header/footer render at exact heights
- TIT has been running this pattern successfully, proving it works

**Why not other options:**
- Option 2 (debug header construction): Header construction is correct, issue is layout
- Option 3 (content height): Height calculation matches TIT exactly, not the issue
- Double placement is the obvious architectural deviation from TIT

**Potential side effects:**
- None expected - this brings CAKE in line with proven TIT pattern
- All content rendering logic remains unchanged
- Only positioning wrapper is removed

## References

- TIT reference: `/Users/jreng/Documents/Poems/dev/tit/internal/ui/layout.go` lines 64-103
- Current CAKE: `/Users/jreng/Documents/Poems/dev/cake/internal/ui/layout.go` lines 45-93
- SPRINT-LOG.md: Sprint 3 - Header Rendering Fixes and Project Name Detection

---

**Engineer, execute this plan to fix header rendering inconsistency.**
