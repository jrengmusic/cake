# Sprint 8 - Menu Fixes: Unix Makefiles, Separator Navigation, Ninja Multi-Config

**Date:** 2026-01-29  
**COUNSELOR:** OpenCode (glm-4.7)  
**ENGINEER:** [To Be Assigned]  

---

## Issues to Fix

### Issue 1: Unix Makefiles Still Present

**Status:** NOT REMOVED - Still appears in generator cycling

**Locations still containing Unix Makefiles:**

1. **`internal/state/project.go` lines 137-142** - Hardcoded fallback:
```go
// Unix Makefiles is always available as fallback
ps.AvailableGenerators = append(ps.AvailableGenerators, Generator{
    Name:          "Unix Makefiles",
    IsIDE:         false,
    IsMultiConfig: false,
})
```

2. **`internal/utils/generator.go` line 12** - In validGenerators list:
```go
validGenerators := []string{
    "Xcode", "Ninja", "Visual Studio",
    "Unix Makefiles", "Ninja Multi-Config",  // ← REMOVE "Unix Makefiles"
}
```

3. **`internal/utils/platform.go` lines 16, 18** - Platform defaults:
```go
case "linux":
    return []string{"Ninja", "Unix Makefiles"}  // ← REMOVE
default:
    return []string{"Unix Makefiles"}  // ← REMOVE
```

4. **`internal/utils/generator.go` line 48** - Dead code in switch:
```go
case "Unix Makefiles":
    return "make", []string{"-C", buildDir}, nil  // ← REMOVE
```

5. **Comments throughout codebase** - References to "Makefiles" in comments

**Fix:** Remove ALL Unix Makefiles references from the codebase.

---

### Issue 2: Separator Not Skipped in Navigation

**Current Behavior:**
- Navigation with ↑↓ moves through separator row
- User can "select" the separator (empty row)
- Footer shows empty hint when separator is selected

**Expected Behavior:**
- Separator should be skipped during navigation (like TIT)
- Navigation: Generator → Regenerate → Open IDE → Configuration → Build → Clean
- Separator is visible but not selectable

**Root Cause:**
`GetVisibleRows()` includes separator because `Visible: true`, but separator should be treated as non-selectable.

**Fix Options:**

**Option A (Recommended):** Add `IsSelectable` field to MenuRow:
```go
type MenuRow struct {
    // ... existing fields ...
    Visible     bool   // true if row should be shown
    IsAction    bool   // true if row triggers action (not toggle)
    IsSelectable bool  // NEW: false for separator
    Hint        string // Footer hint/description for this row
}

// In GenerateMenuRows, set separator:
{
    ID:           "separator",
    Visible:      true,
    IsSelectable: false,  // ← NEW
    // ...
}

// In GetVisibleRows, filter non-selectable:
func (a *Application) GetVisibleRows() []ui.MenuRow {
    var visible []ui.MenuRow
    for _, row := range a.menuItems {
        if row.Visible && row.IsSelectable {  // ← ADD IsSelectable check
            visible = append(visible, row)
        }
    }
    return visible
}
```

**Option B:** Check ID in GetVisibleRows:
```go
func (a *Application) GetVisibleRows() []ui.MenuRow {
    var visible []ui.MenuRow
    for _, row := range a.menuItems {
        if row.Visible && row.ID != "separator" {  // ← Skip separator
            visible = append(visible, row)
        }
    }
    return visible
}
```

---

### Issue 3: "Ninja mult" Display / Menu Height Changes

**Problem:** When selecting "Ninja Multi-Config", the menu height changes.

**Root Cause:** "Ninja Multi-Config" is longer than "Ninja" and may be wrapping or affecting layout calculations.

**Investigation Needed:**

Check `internal/ui/menu.go` - column width calculations:
```go
const (
    shortcutColWidth = 3
    emojiColWidth    = 3
    labelColWidth    = 18  // ← Is this enough for "Ninja Multi-Config"?
    valueColWidth    = 10
    menuBoxWidth     = shortcutColWidth + emojiColWidth + labelColWidth + valueColWidth
)
```

"Ninja Multi-Config" = 18 characters exactly. If the value column also has content, it may overflow.

**Fix:**

1. **Truncate long generator names in Value field:**
```go
// In GenerateMenuRows, truncate long generator names:
generatorLabel := a.projectState.GetGeneratorLabel()
if len(generatorLabel) > 10 {  // Max 10 chars for value column
    generatorLabel = generatorLabel[:10]
}
```

2. **OR increase value column width:**
```go
valueColWidth = 12  // or more
```

3. **OR display shorter name:**
```go
// In project.go GetGeneratorLabel()
func (ps *ProjectState) GetGeneratorLabel() string {
    switch ps.SelectedGenerator {
    case "Ninja Multi-Config":
        return "Ninja Multi"  // Shorter display name
    default:
        return ps.SelectedGenerator
    }
}
```

---

## Implementation Tasks

### Task 1: Remove Unix Makefiles Completely

**Files to modify:**

1. **`internal/state/project.go`** - Remove lines 137-142:
```go
// DELETE:
// Unix Makefiles is always available as fallback
ps.AvailableGenerators = append(ps.AvailableGenerators, Generator{
    Name:          "Unix Makefiles",
    IsIDE:         false,
    IsMultiConfig: false,
})
```

2. **`internal/utils/generator.go`** - Line 12:
```go
// FROM:
validGenerators := []string{
    "Xcode", "Ninja", "Visual Studio",
    "Unix Makefiles", "Ninja Multi-Config",
}

// TO:
validGenerators := []string{
    "Xcode", "Ninja", "Visual Studio", "Ninja Multi-Config",
}
```

3. **`internal/utils/generator.go`** - Remove lines 48-49:
```go
// DELETE:
case "Unix Makefiles":
    return "make", []string{"-C", buildDir}, nil
```

4. **`internal/utils/platform.go`** - Lines 15-19:
```go
// FROM:
case "linux":
    return []string{"Ninja", "Unix Makefiles"}
default:
    return []string{"Unix Makefiles"}

// TO:
case "linux":
    return []string{"Ninja"}
default:
    return []string{"Ninja"}
```

5. **Update comments** - Remove "Makefiles" references in:
   - `internal/state/project.go` lines 15-16, 245
   - `internal/ops/build.go` line 19
   - `internal/ops/clean.go` line 16

6. **`SPEC.md`** - Update documentation (lines 35, 48, 128-130, 134, 278)

---

### Task 2: Skip Separator in Navigation

**File:** `internal/ui/menu.go`

Add `IsSelectable` field:
```go
type MenuRow struct {
    ID           string
    Shortcut     string
    Emoji        string
    Label        string
    Value        string
    Visible      bool
    IsAction     bool
    IsSelectable bool  // NEW
    Hint         string
}
```

Update `GenerateMenuRows()`:
```go
// Separator:
{
    ID:           "separator",
    Visible:      true,
    IsSelectable: false,  // ← NEW
    // ...
}

// All other rows:
{
    ID:           "generator",
    IsSelectable: true,  // ← NEW
    // ...
}
```

**File:** `internal/app/app.go`

Update `GetVisibleRows()`:
```go
func (a *Application) GetVisibleRows() []ui.MenuRow {
    var visible []ui.MenuRow
    for _, row := range a.menuItems {
        if row.Visible && row.IsSelectable {  // ← ADD IsSelectable
            visible = append(visible, row)
        }
    }
    return visible
}
```

---

### Task 3: Fix Menu Height with Long Generator Names

**File:** `internal/state/project.go`

Add `GetGeneratorLabel()` method with truncation:
```go
// GetGeneratorLabel returns a display-friendly generator name
func (ps *ProjectState) GetGeneratorLabel() string {
    name := ps.SelectedGenerator
    
    // Truncate long names for display
    switch name {
    case "Ninja Multi-Config":
        return "Ninja Multi"
    case "Visual Studio 17 2022":
        return "VS 2022"
    case "Visual Studio 16 2019":
        return "VS 2019"
    default:
        return name
    }
}
```

**File:** `internal/app/menu.go`

Update to use `GetGeneratorLabel()`:
```go
return ui.GenerateMenuRows(
    a.projectState.GetGeneratorLabel(),  // ← Use truncated label
    a.projectState.Configuration,
    canOpenIDE,
    canClean,
)
```

---

## File Summary

| File | Changes |
|------|---------|
| `internal/state/project.go` | Remove Unix Makefiles fallback (lines 137-142), add GetGeneratorLabel() method |
| `internal/utils/generator.go` | Remove from validGenerators list, remove switch case |
| `internal/utils/platform.go` | Remove from Linux/default returns |
| `internal/ui/menu.go` | Add IsSelectable field to MenuRow, update GenerateMenuRows() |
| `internal/app/app.go` | Update GetVisibleRows() to check IsSelectable |
| `internal/app/menu.go` | Update to use GetGeneratorLabel() |
| `internal/ops/build.go` | Update comments |
| `internal/ops/clean.go` | Update comments |
| `SPEC.md` | Update documentation |

---

## Acceptance Criteria

1. **Unix Makefiles removed:**
   - [ ] `internal/state/project.go` - no Unix Makefiles fallback
   - [ ] `internal/utils/generator.go` - not in validGenerators
   - [ ] `internal/utils/platform.go` - not in platform lists
   - [ ] Generator cycling: Xcode → Ninja → Ninja Multi-Config → Xcode (on macOS)

2. **Separator skipped:**
   - [ ] Navigation with ↑↓ skips separator
   - [ ] Cannot "select" separator row
   - [ ] Footer shows hint for actual menu items only

3. **Menu height stable:**
   - [ ] Selecting "Ninja Multi-Config" doesn't change menu height
   - [ ] Display shows "Ninja Multi" (truncated)
   - [ ] No text overflow or wrapping

4. **Build:**
   - [ ] `./build.sh` completes without errors
   - [ ] All existing functionality preserved

---

## Notes for ENGINEER

- **Unix Makefiles:** Search entire codebase for "Unix Makefiles" and "Makefiles" to ensure complete removal
- **Separator:** The `IsSelectable` approach is cleaner than hardcoding ID checks
- **Ninja Multi-Config:** Truncation in `GetGeneratorLabel()` is the safest fix - keeps menu layout stable
- **Test:** Verify generator cycling on macOS shows: Xcode → Ninja → Ninja Multi-Config → (back to Xcode)

---

**End of Kickoff Document**
