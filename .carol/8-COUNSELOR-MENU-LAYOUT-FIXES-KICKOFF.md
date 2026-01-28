# Sprint 8 - Menu Fixes: Separator, Generate/Regenerate, Layout

**Date:** 2026-01-29  
**COUNSELOR:** OpenCode (glm-4.7)  
**ENGINEER:** [To Be Assigned]  

---

## Issues to Fix

### Issue 1: No Separator Shows

**Root Cause:** `GetVisibleRows()` filters out separator (IsSelectable: false), but then `RenderCakeMenu()` receives filtered rows without separator.

**Current flow:**
```go
// app.go
visibleRows := a.GetVisibleRows()  // ← Excludes separator!
menuContent := ui.RenderCakeMenu(visibleRows, ...)  // ← No separator to render
```

**Fix:** Pass ALL rows to RenderCakeMenu, let it handle visibility:

```go
// app.go renderMenuWithBanner()
// FROM:
visibleRows := a.GetVisibleRows()
menuContent := ui.RenderCakeMenu(visibleRows, a.selectedIndex, ...)

// TO:
menuContent := ui.RenderCakeMenu(a.menuItems, a.selectedIndex, ...)
```

**Update RenderCakeMenu to:**
1. Skip non-visible rows (don't render them)
2. Render separator line
3. Handle selection index properly (map visible index to array index)

---

### Issue 2: Regenerate Shows When No Build Exists

**Current Behavior:**
- Fresh project (no Builds/ directory)
- Menu shows "Regenerate" 
- Should show "Generate" when no build exists

**Expected Behavior:**
- No builds = "Generate" (initial setup)
- Build exists = "Regenerate" (re-run CMake)

**Root Cause:** Menu always shows "Regenerate" row, should switch label based on build state.

**Fix:** Add `hasBuild` parameter to GenerateMenuRows:

```go
// menu.go
func (a *Application) GenerateMenu() []ui.MenuRow {
    buildInfo := a.projectState.GetSelectedBuildInfo()
    canOpenIDE := a.projectState.CanOpenIDE() && buildInfo.Exists
    canClean := buildInfo.Exists
    hasBuild := buildInfo.Exists  // ← NEW

    return ui.GenerateMenuRows(
        a.projectState.GetGeneratorLabel(),
        a.projectState.Configuration,
        canOpenIDE,
        canClean,
        hasBuild,  // ← NEW
    )
}

// ui/menu.go
func GenerateMenuRows(generatorLabel string, configuration string, canOpenIDE bool, canClean bool, hasBuild bool) []MenuRow {
    // ...
    {
        ID:       "regenerate",
        Label:    map[bool]string{true: "Regenerate", false: "Generate"}[hasBuild],  // ← NEW
        // ...
    }
}
```

---

### Issue 3: Value Column Too Narrow

**Current:** `valueColWidth = 10`
**Requested:** `valueColWidth = 12`

**Fix:** Update in `internal/ui/menu.go`:

```go
const (
    shortcutColWidth = 3
    emojiColWidth    = 3
    labelColWidth    = 18
    valueColWidth    = 12  // ← CHANGED from 10
    menuBoxWidth     = shortcutColWidth + emojiColWidth + labelColWidth + valueColWidth  // ← Now 36
)
```

---

### Issue 4: Ninja Multi-Config Display

**Current:** "Ninja Multi"  
**Requested:** "Ninja (multi)"

**Fix:** Update `GetGeneratorLabel()` in `internal/state/project.go`:

```go
func (ps *ProjectState) GetGeneratorLabel() string {
    name := ps.SelectedGenerator
    
    switch name {
    case "Ninja Multi-Config":
        return "Ninja (multi)"  // ← CHANGED from "Ninja Multi"
    case "Visual Studio 17 2022":
        return "VS 2022"
    case "Visual Studio 16 2019":
        return "VS 2019"
    default:
        return name
    }
}
```

---

### Issue 5: Space Between Highlighted Item and Value

**Current:** `Generator[Xcode]` (no space)  
**Requested:** `Generator [Xcode]` (space between)

**Fix:** In `RenderCakeMenu`, add space after label when selected:

```go
// Column 3: LABEL (left-aligned)
labelStr := row.Label

// ADD: If selected and has value, add space after label
if isSelected && row.Value != "" {
    labelStr = labelStr + " "  // ← Add space
}

labelW := lipgloss.Width(labelStr)
if labelW > labelColWidth {
    labelStr = labelStr[:labelColWidth]
    labelW = labelColWidth
}
labelCol := labelStr + strings.Repeat(" ", labelColWidth-labelW)
```

---

## Implementation Tasks

### Task 1: Fix Separator Rendering

**File:** `internal/app/app.go`

Change line 305-308:
```go
// FROM:
visibleRows := a.GetVisibleRows()
menuContent := ui.RenderCakeMenu(visibleRows, a.selectedIndex, a.theme, a.sizing.ContentHeight, leftWidth)

// TO:
menuContent := ui.RenderCakeMenu(a.menuItems, a.selectedIndex, a.theme, a.sizing.ContentHeight, leftWidth)
```

**File:** `internal/ui/menu.go`

Update `RenderCakeMenu` to handle visibility:
```go
func RenderCakeMenu(rows []MenuRow, selectedIndex int, theme Theme, contentHeight int, contentWidth int) string {
    // ... constants ...
    
    var lines []string
    visibleIndex := 0  // Track visible index separately
    
    for i, row := range rows {
        // Skip hidden rows completely
        if !row.Visible {
            continue
        }
        
        // Handle separator
        if row.ID == "separator" {
            sepLine := lipgloss.NewStyle().
                Foreground(lipgloss.Color(theme.SeparatorColor)).
                Render(strings.Repeat("─", menuBoxWidth))
            lines = append(lines, sepLine)
            continue
        }
        
        // Check if this visible row is selected
        isSelected := visibleIndex == selectedIndex
        
        // ... render row ...
        
        visibleIndex++  // Increment only for visible, non-separator rows
    }
    // ... rest of function ...
}
```

---

### Task 2: Generate vs Regenerate

**File:** `internal/ui/menu.go`

Update function signature:
```go
func GenerateMenuRows(generatorLabel string, configuration string, canOpenIDE bool, canClean bool, hasBuild bool) []MenuRow {
```

Update regenerate row:
```go
{
    ID:       "regenerate",
    Shortcut: "g",
    Emoji:    "🚀",
    Label:    map[bool]string{true: "Regenerate", false: "Generate"}[hasBuild],
    Value:    "",
    Visible:  true,
    IsAction: true,
    IsSelectable: true,
    Hint:     map[bool]string{true: "Re-run CMake configuration", false: "Run initial CMake configuration"}[hasBuild],
}
```

**File:** `internal/app/menu.go`

```go
func (a *Application) GenerateMenu() []ui.MenuRow {
    buildInfo := a.projectState.GetSelectedBuildInfo()
    canOpenIDE := a.projectState.CanOpenIDE() && buildInfo.Exists
    canClean := buildInfo.Exists
    hasBuild := buildInfo.Exists  // ← NEW

    return ui.GenerateMenuRows(
        a.projectState.GetGeneratorLabel(),
        a.projectState.Configuration,
        canOpenIDE,
        canClean,
        hasBuild,  // ← NEW
    )
}
```

---

### Task 3: Widen Value Column

**File:** `internal/ui/menu.go`

```go
const (
    shortcutColWidth = 3
    emojiColWidth    = 3
    labelColWidth    = 18
    valueColWidth    = 12  // ← CHANGED from 10
    menuBoxWidth     = shortcutColWidth + emojiColWidth + labelColWidth + valueColWidth  // = 36
)
```

---

### Task 4: Ninja (multi) Display

**File:** `internal/state/project.go`

```go
func (ps *ProjectState) GetGeneratorLabel() string {
    name := ps.SelectedGenerator
    
    switch name {
    case "Ninja Multi-Config":
        return "Ninja (multi)"  // ← CHANGED
    case "Visual Studio 17 2022":
        return "VS 2022"
    case "Visual Studio 16 2019":
        return "VS 2019"
    default:
        return name
    }
}
```

---

### Task 5: Space Between Label and Value

**File:** `internal/ui/menu.go`

```go
// Column 3: LABEL (left-aligned)
labelStr := row.Label

// ADD SPACE if selected and has value
if isSelected && row.Value != "" {
    labelStr = labelStr + " "
}

labelW := lipgloss.Width(labelStr)
if labelW > labelColWidth {
    labelStr = labelStr[:labelColWidth]
    labelW = labelColWidth
}
labelCol := labelStr + strings.Repeat(" ", labelColWidth-labelW)
```

---

## File Summary

| File | Changes |
|------|---------|
| `internal/app/app.go` | Pass a.menuItems instead of visibleRows to RenderCakeMenu |
| `internal/ui/menu.go` | Handle visibility in RenderCakeMenu, update column width, add space logic, update GenerateMenuRows signature |
| `internal/app/menu.go` | Add hasBuild parameter, pass to GenerateMenuRows |
| `internal/state/project.go` | Change "Ninja Multi" to "Ninja (multi)" |

---

## Acceptance Criteria

1. **Separator shows:**
   - [ ] Separator line visible between Open IDE and Configuration
   - [ ] Navigation skips separator (not selectable)

2. **Generate vs Regenerate:**
   - [ ] Fresh project (no Builds/) shows "Generate"
   - [ ] After build exists, shows "Regenerate"

3. **Value column width:**
   - [ ] Value column is 12 characters wide
   - [ ] "Ninja (multi)" fits without truncation

4. **Display names:**
   - [ ] "Ninja Multi-Config" displays as "Ninja (multi)"

5. **Spacing:**
   - [ ] Space between highlighted label and value: "Generator [Xcode]"

6. **Build:**
   - [ ] `./build.sh` completes without errors

---

**End of Kickoff Document**
