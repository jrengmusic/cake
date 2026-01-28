# Sprint 8 - Critical Menu Bug Fixes

**Date:** 2026-01-29  
**COUNSELOR:** OpenCode (glm-4.7)  
**ENGINEER:** [To Be Assigned]  

---

## Critical Issues Found

### Issue 1: Clean Appears When No Builds Exist

**Root Cause:** Menu is not regenerated when cycling generators.

**Current flow:**
```go
// executeRowAction for "generator"
case "generator":
    a.projectState.CycleToNextGenerator()  // ← Changes generator
    return true, nil  // ← Menu NOT regenerated!
```

The menu still shows old visibility state from previous generator.

**Fix:** Regenerate menu immediately after cycling generator:
```go
case "generator":
    a.projectState.CycleToNextGenerator()
    a.menuItems = a.GenerateMenu()  // ← ADD THIS
    return true, nil
```

---

### Issue 2: Regenerate Shows Instead of Generate

**Same root cause:** Menu not regenerated, shows stale state.

**Fix:** Same as Issue 1 - regenerate menu after generator change.

---

### Issue 3: Text Still Cutoff - "Ninja (multi)" Too Long

**Root Cause:** "Ninja (multi)" = 13 characters, value column = 12 characters

```
Ninja (multi)
1234567890123  ← 13 chars

valueColWidth = 12  ← TOO SMALL
```

**Fix:** Increase value column to 14:
```go
const (
    shortcutColWidth = 3
    emojiColWidth    = 3
    labelColWidth    = 18
    valueColWidth    = 14  // ← INCREASE from 12 to 14
    menuBoxWidth     = shortcutColWidth + emojiColWidth + labelColWidth + valueColWidth
)
```

---

### Issue 4: No Space Between Label and Value

**Root Cause:** Space IS being added but then truncated by label column width.

Current code:
```go
labelStr := row.Label  // "Generator" = 9 chars
if isSelected && row.Value != "" {
    labelStr = labelStr + " "  // "Generator " = 10 chars
}
labelW := lipgloss.Width(labelStr)
if labelW > labelColWidth {  // 10 > 18, false
    labelStr = labelStr[:labelColWidth]
}
labelCol := labelStr + strings.Repeat(" ", labelColWidth-labelW)  // "Generator " + 8 spaces
```

The space IS there, but it's INSIDE the label column. The issue is the value column starts immediately after.

**Real Fix:** The space needs to be BETWEEN the label column and value column, not inside label.

**Better approach:** Reduce label width by 1 when selected and has value:
```go
// Column 3: LABEL (left-aligned)
labelStr := row.Label
labelW := lipgloss.Width(labelStr)

// If selected and has value, reduce label width by 1 to create space
if isSelected && row.Value != "" {
    if labelW < labelColWidth {
        // Label fits, just add padding normally
        labelCol := labelStr + strings.Repeat(" ", labelColWidth-labelW)
    } else {
        // Label would be truncated, truncate early to make room for space
        labelStr = labelStr[:labelColWidth-1]
        labelCol := labelStr + " "
    }
} else {
    labelCol := labelStr + strings.Repeat(" ", labelColWidth-labelW)
}
```

**Actually, simpler fix:** Just add 1 to the padding between columns:

```go
// After building all columns, add space between label and value when selected
if isSelected && row.Value != "" {
    styledLine = shortcutStyle.Render(shortcutCol) +
        emojiStyle.Render(emojiCol) +
        labelStyle.Render(labelCol) +
        " " +  // ← ADD SPACE HERE
        valueStyle.Render(valueCol)
} else {
    styledLine = shortcutStyle.Render(shortcutCol) +
        emojiStyle.Render(emojiCol) +
        labelStyle.Render(labelCol) +
        valueStyle.Render(valueCol)
}
```

---

## Implementation Tasks

### Task 1: Regenerate Menu After Generator Change

**File:** `internal/app/app.go`

```go
func (a *Application) executeRowAction(rowID string) (bool, tea.Cmd) {
    switch rowID {
    case "generator":
        a.projectState.CycleToNextGenerator()
        a.menuItems = a.GenerateMenu()  // ← ADD THIS LINE
        return true, nil
    // ... other cases ...
    }
}
```

---

### Task 2: Increase Value Column Width

**File:** `internal/ui/menu.go`

```go
const (
    shortcutColWidth = 3
    emojiColWidth    = 3
    labelColWidth    = 18
    valueColWidth    = 14  // ← CHANGED from 12
    menuBoxWidth     = shortcutColWidth + emojiColWidth + labelColWidth + valueColWidth  // = 38
)
```

---

### Task 3: Add Space Between Label and Value

**File:** `internal/ui/menu.go`

In the rendering section, add space when selected:

```go
if isSelected {
    // Selected: highlight label+value
    shortcutStyle := lipgloss.NewStyle().
        Foreground(lipgloss.Color(theme.LabelTextColor))
    emojiStyle := lipgloss.NewStyle().
        Foreground(lipgloss.Color(theme.LabelTextColor))
    labelStyle := lipgloss.NewStyle().
        Foreground(lipgloss.Color(theme.MainBackgroundColor)).
        Background(lipgloss.Color(theme.MenuSelectionBackground)).
        Bold(true)
    valueStyle := lipgloss.NewStyle().
        Foreground(lipgloss.Color(theme.AccentTextColor)).
        Bold(true)
    
    // ADD SPACE between label and value when selected
    space := ""
    if row.Value != "" {
        space = " "
    }
    
    styledLine = shortcutStyle.Render(shortcutCol) +
        emojiStyle.Render(emojiCol) +
        labelStyle.Render(labelCol) +
        space +  // ← ADDED SPACE
        valueStyle.Render(valueCol)
} else {
    // Normal (no space needed)
    // ... existing code ...
}
```

Also REMOVE the old space-adding logic:
```go
// REMOVE THIS SECTION:
// Check if this visible row is selected
isSelected = visibleIndex == selectedIndex

// Column 3: LABEL (left-aligned)
labelStr := row.Label

// Add space if selected and has value
if isSelected && row.Value != "" {
    labelStr = labelStr + " "
}
```

---

### Task 4: Also Regenerate Menu After Configuration Change

**File:** `internal/app/app.go`

```go
func (a *Application) executeRowAction(rowID string) (bool, tea.Cmd) {
    switch rowID {
    // ...
    case "configuration":
        a.projectState.CycleConfiguration()
        a.menuItems = a.GenerateMenu()  // ← ADD THIS LINE
        return true, nil
    // ...
    }
}
```

---

## File Summary

| File | Changes |
|------|---------|
| `internal/app/app.go` | Regenerate menu after generator change (line ~407), regenerate after config change (line ~417) |
| `internal/ui/menu.go` | Increase valueColWidth to 14, add space between label and value when selected, remove old space logic |

---

## Acceptance Criteria

1. **Generator cycling updates menu:**
   - [ ] Switching from Xcode → Ninja updates menu visibility immediately
   - [ ] Clean option shows/hides correctly based on selected generator's build state
   - [ ] Generate/Regenerate label updates correctly

2. **No text cutoff:**
   - [ ] "Ninja (multi)" displays fully (13 chars fit in 14-char column)

3. **Space between label and value:**
   - [ ] When Generator row selected: "Generator [Ninja (multi)]" (space before value)
   - [ ] When Configuration selected: "Configuration [Debug]" (space before value)

4. **Build:**
   - [ ] `./build.sh` completes without errors

---

## Notes for ENGINEER

- **Critical bug:** Menu not regenerating was causing ALL the visibility issues
- **Space logic:** Add space at render time, not during label construction
- **Test:** Delete Builds/ directory, run cake, cycle through generators - Clean should only appear for generators with builds

---

**End of Kickoff Document**
