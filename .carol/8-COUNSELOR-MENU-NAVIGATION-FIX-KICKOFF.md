# Sprint 8 - Menu Navigation & Separator Width Fix

**Date:** 2026-01-29  
**COUNSELOR:** OpenCode (glm-4.7)  
**ENGINEER:** [To Be Assigned]  

---

## Objective

Fix two UI issues in the cake menu:
1. **Navigation:** Selection index should skip hidden items (Open IDE, Clean when not available)
2. **Separator:** Should be drawn at menu content width, not full column width

---

## Issue 1: Selection Index Counts Hidden Items

### Current Behavior
When `Open IDE` or `Clean` are hidden (invisible), the selection index still counts them:

```
Menu structure (indices):
[0] Generator      ← visible
[1] Regenerate     ← visible  
[2] Open IDE       ← HIDDEN (not configured yet)
[3] Separator      ← visible (always)
[4] Configuration  ← visible
[5] Build          ← visible
[6] Clean          ← HIDDEN (no build yet)
```

**Problem:** When navigating with ↑↓, selection jumps over hidden items but index still references the full array. This causes:
- Footer shows wrong item's hint
- Selection highlight can appear on empty rows
- Navigation feels inconsistent

### Expected Behavior
Navigation should use **visible indices only** (0-4 when Open IDE and Clean are hidden):

```
Visible items only:
[0] Generator      
[1] Regenerate     
[2] Separator      
[3] Configuration  
[4] Build          
```

**Navigation:** 0→1→2→3→4, footer shows correct hint for each.

---

## Issue 2: Separator Width Too Wide

### Current Behavior
Separator line spans full column width (`contentWidth`):

```go
// internal/ui/menu.go line 107-109
sepLine := lipgloss.NewStyle().
    Foreground(lipgloss.Color(theme.SeparatorColor)).
    Render(strings.Repeat("─", contentWidth))  // ← full column width
```

**Result:** Separator extends beyond menu content into empty space.

### Expected Behavior
Separator should match menu content width only (`menuBoxWidth`):

```
Menu columns: SHORTCUT(3) | EMOJI(3) | LABEL(18) | VALUE(10) = 34 chars
Separator should be: ────────────────────────────────── (34 chars)
Not: ────────────────────────────────────────────────────────── (full width)
```

---

## Implementation Tasks

### Task 1: Fix Navigation to Use Visible Indices

**File:** `internal/app/app.go`

**Problem:** `GetVisibleRows()` returns ALL rows, not just visible ones. Navigation uses array indices but should use visible indices.

**Solution:** Create a proper `GetVisibleRows()` that filters out hidden rows:

```go
// GetVisibleRows returns only visible menu items (excludes hidden rows)
func (a *Application) GetVisibleRows() []ui.MenuRow {
	var visible []ui.MenuRow
	for _, row := range a.menuItems {
		if row.Visible {
			visible = append(visible, row)
		}
	}
	return visible
}
```

**Update navigation in `handleMenuKeyPress()`:**

Current code (line 504-531) uses `selectedIndex` as array index. Need to:
1. Map between visible index and array index
2. Update `ToggleRowAtIndex()` to work with visible indices
3. Update shortcut handlers (g, o, b, c) to find visible indices

**New helper functions needed:**

```go
// GetVisibleIndex returns the visible index for a given row ID
// Returns -1 if row is hidden
func (a *Application) GetVisibleIndex(rowID string) int {
	visibleIndex := 0
	for _, row := range a.menuItems {
		if row.ID == rowID {
			if row.Visible {
				return visibleIndex
			}
			return -1
		}
		if row.Visible {
			visibleIndex++
		}
	}
	return -1
}

// GetArrayIndex returns the array index for a given visible index
func (a *Application) GetArrayIndex(visibleIdx int) int {
	visibleCount := 0
	for i, row := range a.menuItems {
		if row.Visible {
			if visibleCount == visibleIdx {
				return i
			}
			visibleCount++
		}
	}
	return -1
}

// ToggleRowAtIndex handles menu row toggle/action at given VISIBLE index
func (a *Application) ToggleRowAtIndex(visibleIndex int) (bool, tea.Cmd) {
	arrayIndex := a.GetArrayIndex(visibleIndex)
	if arrayIndex < 0 || arrayIndex >= len(a.menuItems) {
		return false, nil
	}

	row := a.menuItems[arrayIndex]
	return a.executeRowAction(row.ID)
}
```

**Update `handleMenuKeyPress()` navigation:**

```go
func (a *Application) handleMenuKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	visibleRows := a.GetVisibleRows()
	visibleCount := len(visibleRows)

	switch msg.String() {
	case "up", "k":
		if a.selectedIndex > 0 {
			a.selectedIndex--
		}
		return a, nil
	case "down", "j":
		if a.selectedIndex < visibleCount-1 {
			a.selectedIndex++
		}
		return a, nil
	case "enter", " ":
		if a.selectedIndex >= 0 && a.selectedIndex < visibleCount {
			handled, cmd := a.ToggleRowAtIndex(a.selectedIndex)
			if handled {
				a.menuItems = a.GenerateMenu()
				// Clamp selectedIndex if visible count changed
				newVisible := a.GetVisibleRows()
				if a.selectedIndex >= len(newVisible) {
					a.selectedIndex = len(newVisible) - 1
					if a.selectedIndex < 0 {
						a.selectedIndex = 0
					}
				}
				return a, cmd
			}
		}
		return a, nil
	// ... shortcut handlers use GetVisibleIndex()
	case "o", "O":
		idx := a.GetVisibleIndex("openIde")
		if idx >= 0 {
			a.selectedIndex = idx
			handled, cmd := a.ToggleRowAtIndex(idx)
			// ...
		}
		return a, nil
	// ... etc
	}
}
```

**Update `footer.go` to use visible indices:**

```go
// getMenuFooter returns footer for menu mode (selected item's hint)
func (a *Application) getMenuFooter(width int) string {
	visibleRows := a.GetVisibleRows()
	if a.selectedIndex < 0 || a.selectedIndex >= len(visibleRows) {
		return ""
	}

	selectedRow := visibleRows[a.selectedIndex]
	
	if selectedRow.Hint != "" {
		return ui.RenderFooterHint(selectedRow.Hint, width, &a.theme)
	}

	return ""
}
```

**Update `menu.go` rendering:**

The `RenderCakeMenu()` function needs to receive the visible index, not array index:

```go
// In app.go renderMenuWithBanner()
visibleRows := a.GetVisibleRows()
menuContent := ui.RenderCakeMenu(visibleRows, a.selectedIndex, a.theme, a.sizing.ContentHeight, leftWidth)
```

---

### Task 2: Fix Separator Width

**File:** `internal/ui/menu.go`

**Current (line 107-109):**
```go
if row.ID == "separator" {
    sepLine := lipgloss.NewStyle().
        Foreground(lipgloss.Color(theme.SeparatorColor)).
        Render(strings.Repeat("─", contentWidth))  // ← WRONG: full column width
    lines = append(lines, sepLine)
    continue
}
```

**Fix:** Use `menuBoxWidth` instead of `contentWidth`:

```go
if row.ID == "separator" {
    sepLine := lipgloss.NewStyle().
        Foreground(lipgloss.Color(theme.SeparatorColor)).
        Render(strings.Repeat("─", menuBoxWidth))  // ← CORRECT: menu content width
    lines = append(lines, sepLine)
    continue
}
```

---

## Files to Modify

| File | Changes |
|------|---------|
| `internal/app/app.go` | Rewrite `GetVisibleRows()` to filter hidden, add `GetVisibleIndex()`, `GetArrayIndex()`, update `ToggleRowAtIndex()`, update `handleMenuKeyPress()` navigation, update shortcut handlers |
| `internal/app/footer.go` | Update `getMenuFooter()` to use visible rows |
| `internal/app/app.go` | Update `renderMenuWithBanner()` to pass visible rows to renderer |
| `internal/ui/menu.go` | Fix separator width: `contentWidth` → `menuBoxWidth` (line 109) |

---

## Acceptance Criteria

1. **Navigation with hidden items:**
   - [ ] When Open IDE is hidden, navigation skips it (0→1→2→3→4)
   - [ ] When Clean is hidden, it's not in navigation path
   - [ ] Footer shows correct hint for selected visible item
   - [ ] Selection highlight appears on correct row

2. **Separator width:**
   - [ ] Separator line matches menu content width (34 chars)
   - [ ] Does not extend into empty column space

3. **Shortcuts still work:**
   - [ ] `o` key jumps to Open IDE only if visible
   - [ ] `c` key jumps to Clean only if visible
   - [ ] `g`, `b` work as before

4. **Build:**
   - [ ] `./build.sh` completes without errors
   - [ ] No regression in menu functionality

---

## Notes for ENGINEER

- **Key insight:** Navigation should work on VISIBLE indices (0 to N-1), not array indices
- **Mapping:** Need bidirectional mapping between visible index ↔ array index
- **Separator:** Simple fix - change one variable in menu.go
- **Test scenario:** Fresh project (no build) should show only: Generator, Regenerate, Separator, Configuration, Build

---

**End of Kickoff Document**
