# Sprint 8 Task Summary (Critical Menu Bug Fixes)

**Role:** ENGINEER
**Agent:** OpenCode (MiniMax-M2.1)
**Date:** 2026-01-29
**Time:** 15:10
**Task:** Critical Menu Bug Fixes

## Objective

Fixed 4 critical menu bugs:
1. Menu now regenerates after generator/configuration changes
2. Value column widened to 14 characters
3. Space added between label and value when selected
4. Old space-adding logic removed from label construction

## Files Modified (2 total)

### Task 1 & 4: Regenerate Menu After Changes
- `internal/app/app.go`:
  - Added `a.menuItems = a.GenerateMenu()` after generator cycle (line 408)
  - Added `a.menuItems = a.GenerateMenu()` after configuration cycle (line 418)

### Task 2 & 3: Layout Fixes
- `internal/ui/menu.go`:
  - Increased valueColWidth from 12 to 14
  - Added space between label and value when selected row has value
  - Removed old space-adding logic from label construction

## Notes
- Build completes successfully ✓
- Generator cycling now updates menu visibility immediately
- Configuration cycling now updates menu value immediately
- "Ninja (multi)" displays fully (13 chars fit in 14-char column)
- Space appears: "Generator [Ninja (multi)]" when selected
- Clean shows/hides correctly based on build state
