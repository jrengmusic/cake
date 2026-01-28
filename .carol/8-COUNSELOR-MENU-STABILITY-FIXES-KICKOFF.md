# Sprint 8 - Menu Layout & Stability Fixes

**Date:** 2026-01-29  
**COUNSELOR:** OpenCode (glm-4.7)  
**ENGINEER:** [To Be Assigned]  

---

## Issues to Fix

### Issue 1: Value Shifted by 1 Char When Highlighted

**Problem:** When Generator or Configuration row is selected, the value shifts 1 character to the right.

**Current Behavior:**
```
Not selected:  "Generator  Xcode"
Selected:      "Generator [ Xcode]"  ← Value shifted right by 1
```

**Root Cause:** Space is being added between label and value columns when selected, but this changes the total width.

**Current Code (menu.go lines 197-207):**
```go
// Add space between label and value when selected and has value
space := ""
if row.Value != "" {
    space = " "
}

styledLine = shortcutStyle.Render(shortcutCol) +
    emojiStyle.Render(emojiCol) +
    labelStyle.Render(labelCol) +
    space +  // ← This adds width!
    valueStyle.Render(valueCol)
```

**Fix:** Don't add space - instead, reduce label padding by 1 when selected and has value:

```go
// Column 3: LABEL (left-aligned)
labelStr := row.Label
labelW := lipgloss.Width(labelStr)
if labelW > labelColWidth {
    labelStr = labelStr[:labelColWidth]
    labelW = labelColWidth
}

// Calculate padding
padding := labelColWidth - labelW

// If selected and has value, reduce padding by 1 to create visual space
if isSelected && row.Value != "" && padding > 0 {
    padding--
}

labelCol := labelStr + strings.Repeat(" ", padding)
```

---

### Issue 2: Menu Shifts Up/Down When Items Hidden

**Problem:** When Open IDE and Clean are hidden, the menu recenters with fewer lines, causing it to shift up.

**Current Behavior:**
```
All items visible (7 lines):
  [Generator]
  [Regenerate]
  [Open IDE]
  ───────────
  [Configuration]
  [Build]
  [Clean]

Some hidden (5 lines):
  [Generator]      ← Menu shifted UP
  [Regenerate]
  ───────────
  [Configuration]
  [Build]
```

**Root Cause:** Hidden items are skipped with `continue`, reducing `lines` array length.

**Current Code (menu.go lines 122-125):**
```go
// Skip hidden rows completely
if !row.Visible {
    continue  // ← This removes the line entirely!
}
```

**Fix:** Render hidden items as empty lines instead of skipping:

```go
// Handle hidden rows - render as empty line to maintain fixed height
if !row.Visible {
    lines = append(lines, strings.Repeat(" ", menuBoxWidth))
    continue
}
```

---

### Issue 3: Menu Must Always Have Fixed 7 Items Height

**Requirement:** Menu must ALWAYS render exactly 7 lines, regardless of visibility:
- Visible items: rendered normally
- Hidden items: rendered as empty lines
- Separator: always rendered

**This ensures:**
- Menu never shifts position
- Separator always in same place
- Navigation indices stable

**Implementation:** Already partially done in Issue 2 fix above.

---

## Implementation Tasks

### Task 1: Fix Value Shifting

**File:** `internal/ui/menu.go`

Replace the space-adding logic with padding reduction:

```go
// Column 3: LABEL (left-aligned) - lines 161-168
labelStr := row.Label
labelW := lipgloss.Width(labelStr)
if labelW > labelColWidth {
    labelStr = labelStr[:labelColWidth]
    labelW = labelColWidth
}

// Calculate padding
padding := labelColWidth - labelW

// If selected and has value, reduce padding by 1 to create visual space
// This keeps total width constant while adding visual separation
if isSelected && row.Value != "" && padding > 0 {
    padding--
}

labelCol := labelStr + strings.Repeat(" ", padding)
```

Then REMOVE the space-adding logic in the rendering section (lines 197-207):
```go
// REMOVE THIS:
space := ""
if row.Value != "" {
    space = " "
}

styledLine = shortcutStyle.Render(shortcutCol) +
    emojiStyle.Render(emojiCol) +
    labelStyle.Render(labelCol) +
    space +  // ← REMOVE
    valueStyle.Render(valueCol)

// REPLACE WITH:
styledLine = shortcutStyle.Render(shortcutCol) +
    emojiStyle.Render(emojiCol) +
    labelStyle.Render(labelCol) +
    valueStyle.Render(valueCol)
```

---

### Task 2: Fix Hidden Items Rendering

**File:** `internal/ui/menu.go`

Change hidden row handling (lines 121-125):

```go
// FROM:
// Skip hidden rows completely
if !row.Visible {
    continue
}

// TO:
// Handle hidden rows - render as empty line to maintain fixed height
if !row.Visible {
    lines = append(lines, strings.Repeat(" ", menuBoxWidth))
    visibleIndex++  // Still increment index for selection tracking
    continue
}
```

---

### Task 3: Verify Fixed 7-Line Height

After Task 2, verify the menu always renders 7 lines:
- Line 0: Generator (or empty if hidden)
- Line 1: Regenerate (or empty if hidden)
- Line 2: Open IDE (or empty if hidden)
- Line 3: Separator (always)
- Line 4: Configuration (or empty if hidden)
- Line 5: Build (or empty if hidden)
- Line 6: Clean (or empty if hidden)

**Note:** Currently only Open IDE and Clean can be hidden. Generator, Regenerate, Configuration, Build are always visible.

---

## File Summary

| File | Changes |
|------|---------|
| `internal/ui/menu.go` | Fix label padding logic (lines 161-168), remove space-adding logic (lines 197-207), fix hidden row handling (lines 121-125) |

---

## Acceptance Criteria

1. **No value shifting:**
   - [ ] Generator row selected: "Generator[Xcode]" (no extra space)
   - [ ] Configuration row selected: "Configuration[Debug]" (no extra space)
   - [ ] Value column stays aligned regardless of selection

2. **Fixed menu height:**
   - [ ] Menu always renders exactly 7 lines
   - [ ] Hidden items show as blank lines (not removed)
   - [ ] Menu position stays constant when items hide/show

3. **Stable separator:**
   - [ ] Separator always in same position (line 3)
   - [ ] Doesn't move when Open IDE/Clean hide/show

4. **Build:**
   - [ ] `./build.sh` completes without errors

---

## Notes for ENGINEER

- **Key insight:** Keep total width constant - don't add characters, adjust padding instead
- **Fixed height:** Always render 7 lines, use empty strings for hidden items
- **Test:** Hide/show items by cycling generators and verify menu doesn't shift
- **Visual space:** The "space" between label and value comes from reducing label padding by 1, not adding a character

---

**End of Kickoff Document**
