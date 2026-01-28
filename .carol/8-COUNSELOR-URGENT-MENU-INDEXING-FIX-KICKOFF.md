# Sprint 8 - URGENT: Fix Broken Menu Indexing

**Date:** 2026-01-29  
**COUNSELOR:** OpenCode (glm-4.7)  
**ENGINEER:** [To Be Assigned]  

---

## Critical Issues Found

### Issue 1: Indexing Completely Broken

**Symptoms:**
- Separator gets selected
- Footer shows wrong description
- Can't reach Clean (capped at Build)
- Value column shifts left when highlighted

**Root Cause:** The ENGINEER mixed up array indices vs visible indices.

**Current Broken Code:**
```go
var lines []string
visibleIndex := 0 // This tracks VISIBLE SELECTABLE items

for _, row := range rows {
    if !row.Visible {
        lines = append(lines, strings.Repeat(" ", menuBoxWidth))
        visibleIndex++ // ← WRONG! Hidden rows shouldn't increment
        continue
    }
    
    if row.ID == "separator" {
        // ... render separator ...
        continue // ← WRONG! Doesn't increment visibleIndex
    }
    
    isSelected := visibleIndex == selectedIndex // ← WRONG! 
    // selectedIndex is VISIBLE index, but we're comparing to array position
```

**The Problem:**
- `selectedIndex` from app.go is the **visible selectable index** (0-4)
- `visibleIndex` in the loop is being incremented for ALL rows including hidden/separator
- This causes complete misalignment

---

## Correct Architecture

**Two Index Systems:**
1. **Array Index** (0-6): Position in the 7-row array
2. **Visible Selectable Index** (0-4): Position among visible, selectable items

**Mapping:**
```
Array Index | Row ID        | Visible | Selectable | Visible Selectable Index
------------|---------------|---------|------------|------------------------
0           | generator     | Yes     | Yes        | 0
1           | regenerate    | Yes     | Yes        | 1
2           | openIde       | Maybe   | Yes        | 2 (if visible)
3           | separator     | Yes     | No         | - (skip)
4           | configuration | Yes     | Yes        | 3
5           | build         | Yes     | Yes        | 4
6           | clean         | Maybe   | Yes        | 5 (if visible)
```

**Key Insight:**
- Navigation uses **Visible Selectable Index** (0-4 or 0-5)
- Rendering needs to map this back to array position
- Separator is visible but NOT selectable (skipped in navigation)
- Hidden items are NOT visible and NOT selectable

---

## Implementation Tasks

### Task 1: Fix RenderCakeMenu Function

**File:** `internal/ui/menu.go`

**Complete Rewrite of Rendering Logic:**

```go
func RenderCakeMenu(rows []MenuRow, selectedVisibleIndex int, theme Theme, contentHeight int, contentWidth int) string {
    const (
        shortcutColWidth = 3
        emojiColWidth    = 3
        labelColWidth    = 18
        valueColWidth    = 14
        menuBoxWidth     = shortcutColWidth + emojiColWidth + labelColWidth + valueColWidth
    )

    var lines []string
    visibleSelectableIndex := 0 // Tracks only visible AND selectable items
    
    for _, row := range rows {
        // Handle hidden rows - render empty line, don't count in navigation
        if !row.Visible {
            lines = append(lines, strings.Repeat(" ", menuBoxWidth))
            continue // Don't increment visibleSelectableIndex
        }
        
        // Handle separator - render it, but don't count in navigation
        if row.ID == "separator" {
            sepLine := lipgloss.NewStyle().
                Foreground(lipgloss.Color(theme.SeparatorColor)).
                Render(strings.Repeat("─", menuBoxWidth))
            lines = append(lines, sepLine)
            continue // Don't increment visibleSelectableIndex
        }
        
        // This row is visible and selectable
        // Check if it's selected
        isSelected := visibleSelectableIndex == selectedVisibleIndex
        
        // Column 1: SHORTCUT
        shortcutStr := row.Shortcut
        shortcutW := lipgloss.Width(shortcutStr)
        shortcutPad := shortcutColWidth - shortcutW
        if shortcutPad < 0 {
            shortcutPad = 0
        }
        shortcutCol := shortcutStr + strings.Repeat(" ", shortcutPad)
        
        // Column 2: EMOJI
        emojiStr := row.Emoji
        emojiW := lipgloss.Width(emojiStr)
        emojiLeftPad := (emojiColWidth - emojiW) / 2
        emojiRightPad := emojiColWidth - emojiW - emojiLeftPad
        if emojiLeftPad < 0 {
            emojiLeftPad = 0
        }
        if emojiRightPad < 0 {
            emojiRightPad = 0
        }
        emojiCol := strings.Repeat(" ", emojiLeftPad) + emojiStr + strings.Repeat(" ", emojiRightPad)
        
        // Column 3: LABEL
        labelStr := row.Label
        labelW := lipgloss.Width(labelStr)
        if labelW > labelColWidth {
            labelStr = labelStr[:labelColWidth]
            labelW = labelColWidth
        }
        
        // Calculate padding - reduce by 1 when selected to create visual space
        padding := labelColWidth - labelW
        if isSelected && row.Value != "" && padding > 0 {
            padding--
        }
        labelCol := labelStr + strings.Repeat(" ", padding)
        
        // Column 4: VALUE
        valueStr := row.Value
        valueW := lipgloss.Width(valueStr)
        if valueW > valueColWidth {
            valueStr = valueStr[:valueColWidth]
            valueW = valueColWidth
        }
        valueCol := strings.Repeat(" ", valueColWidth-valueW) + valueStr
        
        // Build styled line
        var styledLine string
        
        if isSelected {
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
            
            styledLine = shortcutStyle.Render(shortcutCol) +
                emojiStyle.Render(emojiCol) +
                labelStyle.Render(labelCol) +
                valueStyle.Render(valueCol)
        } else {
            shortcutStyle := lipgloss.NewStyle().
                Foreground(lipgloss.Color(theme.LabelTextColor))
            emojiStyle := lipgloss.NewStyle().
                Foreground(lipgloss.Color(theme.LabelTextColor))
            labelStyle := lipgloss.NewStyle().
                Foreground(lipgloss.Color(theme.LabelTextColor))
            valueStyle := lipgloss.NewStyle().
                Foreground(lipgloss.Color(theme.ContentTextColor))
            
            styledLine = shortcutStyle.Render(shortcutCol) +
                emojiStyle.Render(emojiCol) +
                labelStyle.Render(labelCol) +
                valueStyle.Render(valueCol)
        }
        
        lines = append(lines, styledLine)
        visibleSelectableIndex++ // Only increment for visible, selectable rows
    }
    
    // ... rest of function (centering) stays same ...
}
```

---

### Task 2: Fix Value Column Width Issue

**Problem:** Value shifts left when selected because total width changes.

**Root Cause:** Reducing label padding by 1 makes label column narrower, but value column stays same. This causes misalignment.

**Better Fix:** Don't change padding at all. The "space" should come from the value column's left padding.

```go
// Column 4: VALUE (right-aligned)
valueStr := row.Value
valueW := lipgloss.Width(valueStr)
if valueW > valueColWidth {
    valueStr = valueStr[:valueColWidth]
    valueW = valueColWidth
}

// Calculate left padding for value
valueLeftPad := valueColWidth - valueW

// If selected and has value, add 1 space before value for visual separation
// BUT keep total width constant by reducing left padding
if isSelected && row.Value != "" && valueLeftPad > 0 {
    valueLeftPad--
    valueCol := strings.Repeat(" ", valueLeftPad) + " " + valueStr
} else {
    valueCol := strings.Repeat(" ", valueLeftPad) + valueStr
}
```

Actually, even simpler: **Don't add any space at all.** Just render normally:

```go
// Column 3: LABEL (left-aligned) - NO PADDING CHANGE
labelStr := row.Label
labelW := lipgloss.Width(labelStr)
if labelW > labelColWidth {
    labelStr = labelStr[:labelColWidth]
    labelW = labelColWidth
}
labelCol := labelStr + strings.Repeat(" ", labelColWidth-labelW)

// Column 4: VALUE (right-aligned) - NO SPACE ADDED
valueStr := row.Value
valueW := lipgloss.Width(valueStr)
if valueW > valueColWidth {
    valueStr = valueStr[:valueColWidth]
    valueW = valueColWidth
}
valueCol := strings.Repeat(" ", valueColWidth-valueW) + valueStr
```

The "space" between label and value comes naturally from the column widths. Don't try to add extra space.

---

## File Summary

| File | Changes |
|------|---------|
| `internal/ui/menu.go` | Complete rewrite of RenderCakeMenu indexing logic |

---

## Acceptance Criteria

1. **Correct indexing:**
   - [ ] Navigation: 0=Generator, 1=Regenerate, 2=Open IDE (if visible), 3=Configuration, 4=Build, 5=Clean (if visible)
   - [ ] Separator never selected
   - [ ] Footer shows correct hint for selected item
   - [ ] Can navigate to Clean when visible

2. **Fixed height:**
   - [ ] Always renders exactly 7 lines
   - [ ] Hidden items show as empty lines
   - [ ] Menu doesn't shift position

3. **No value shifting:**
   - [ ] Value column stays aligned regardless of selection
   - [ ] No extra space added between label and value

4. **Build:**
   - [ ] `./build.sh` completes without errors

---

## Critical Notes for ENGINEER

- **DO NOT increment visibleSelectableIndex for hidden or separator rows**
- **selectedIndex from app.go is the VISIBLE SELECTABLE index (0-4 or 0-5)**
- **Separator is visible but NOT selectable - skip it in index counting**
- **Hidden items are NOT visible - skip them in index counting**
- **The 7 rows are: Generator, Regenerate, Open IDE, Separator, Configuration, Build, Clean**

---

**End of Kickoff Document**
