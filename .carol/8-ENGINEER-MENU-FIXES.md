# Sprint 8 Task Summary (Menu Fixes)

**Role:** ENGINEER
**Agent:** OpenCode (MiniMax-M2.1)
**Date:** 2026-01-29
**Time:** 14:50
**Task:** Menu Fixes: Unix Makefiles, Separator Navigation, Ninja Multi-Config

## Objective

Fixed 3 menu issues:
1. Removed Unix Makefiles completely from generator options
2. Separator is now skipped in navigation (not selectable)
3. Fixed menu height with long generator names ("Ninja Multi-Config" → "Ninja Multi")

## Files Modified (5 total)

### Task 1: Remove Unix Makefiles
- `internal/state/project.go` - Removed Unix Makefiles fallback (lines 137-142)
- `internal/utils/generator.go` - Removed from validGenerators list AND switch case
- `internal/utils/platform.go` - Removed from Linux/default returns

### Task 2: Skip Separator in Navigation
- `internal/ui/menu.go`:
  - Added `IsSelectable` field to MenuRow struct
  - Set `IsSelectable: false` for separator row
  - Set `IsSelectable: true` for all other rows

- `internal/app/app.go`:
  - Updated `GetVisibleRows()` to check `row.Visible && row.IsSelectable`

### Task 3: Fix Menu Height
- `internal/state/project.go`:
  - Updated `GetGeneratorLabel()` with truncation:
    - "Ninja Multi-Config" → "Ninja Multi"
    - "Visual Studio 17 2022" → "VS 2022"
    - "Visual Studio 16 2019" → "VS 2019"

## Notes
- Build completes successfully ✓
- Generator cycling: Xcode → Ninja → Ninja Multi-Config → Xcode (on macOS)
- Navigation skips separator (Generator → Regenerate → Open IDE → Configuration → Build → Clean)
- Menu height stable regardless of generator name length
- Footer shows hints for actual menu items only
