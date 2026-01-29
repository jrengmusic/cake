# Critical Bug Fix Required

## Bug: Clean Triggers Build Operation

**Root Cause:** Index mismatch between `GetVisibleRows()` and `GetArrayIndex()`

### The Problem

`GetVisibleRows()` filters by BOTH `Visible && IsSelectable`:
```go
func (a *Application) GetVisibleRows() []ui.MenuRow {
    var visible []ui.MenuRow
    for _, row := range a.menuItems {
        if row.Visible && row.IsSelectable {  // <-- Checks BOTH
            visible = append(visible, row)
        }
    }
    return visible
}
```

But `GetArrayIndex()` only checks `Visible`:
```go
func (a *Application) GetArrayIndex(visibleIdx int) int {
    visibleCount := 0
    for i, row := range a.menuItems {
        if row.Visible {  // <-- Only checks Visible, NOT IsSelectable!
            if visibleCount == visibleIdx {
                return i
            }
            visibleCount++
        }
    }
    return -1
}
```

### Impact

The separator row has `Visible=true` but `IsSelectable=false`. This causes an off-by-one error:

Menu array:
- 0: generator (Visible, Selectable)
- 1: regenerate (Visible, Selectable)
- 2: openIde (Visible, Selectable)
- 3: separator (Visible, NOT Selectable) ← Problem!
- 4: configuration (Visible, Selectable)
- 5: build (Visible, Selectable)
- 6: clean (Visible, Selectable)

Visible rows (what user sees):
- 0: generator
- 1: regenerate
- 2: openIde
- 3: configuration
- 4: build
- 5: clean ← User selects this

When user selects Clean (visible index 5):
- `GetArrayIndex(5)` counts: 0,1,2,3(separator),4,5 → returns array index 5
- Array index 5 is "build", not "clean"!

**Result:** Selecting Clean triggers Build operation.

### Fix

Update `GetArrayIndex` to match `GetVisibleRows`:

```go
func (a *Application) GetArrayIndex(visibleIdx int) int {
    visibleCount := 0
    for i, row := range a.menuItems {
        if row.Visible && row.IsSelectable {  // <-- Add row.IsSelectable check
            if visibleCount == visibleIdx {
                return i
            }
            visibleCount++
        }
    }
    return -1
}
```

Also check `GetVisibleIndex`:

```go
func (a *Application) GetVisibleIndex(rowID string) int {
    visibleIndex := 0
    for _, row := range a.menuItems {
        if row.ID == rowID {
            if row.Visible && row.IsSelectable {  // <-- Should also check IsSelectable
                return visibleIndex
            }
            return -1
        }
        if row.Visible && row.IsSelectable {  // <-- Add row.IsSelectable check
            visibleIndex++
        }
    }
    return -1
}
```

### Files to Fix

1. `internal/app/app.go` - Fix `GetArrayIndex()` and `GetVisibleIndex()` to check `row.IsSelectable`

### Verification

After fix:
- Select Clean → Should show "Cleaning build directory..."
- Select Build → Should show "Building: ..."
- Confirmation dialog should appear for Clean
- ESC during operation should abort
