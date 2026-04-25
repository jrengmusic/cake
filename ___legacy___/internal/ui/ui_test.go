package ui

import (
	"sync"
	"testing"
)

// --- CalculateDynamicSizing ---

func TestCalculateDynamicSizing(t *testing.T) {
	tests := []struct {
		name       string
		w, h       int
		tooSmall   bool
		minContent int
		minInner   int
	}{
		{"80x24 standard", 80, 24, false, MinContentHeight, 20},
		{"120x40 large", 120, 40, false, MinContentHeight, 20},
		{"40x10 small", 40, 10, true, MinContentHeight, 20},
		{"MinWidth exact", MinWidth, MinHeight, false, MinContentHeight, 20},
		{"below min width", MinWidth - 1, MinHeight, true, MinContentHeight, 20},
		{"below min height", MinWidth, MinHeight - 1, true, MinContentHeight, 20},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := CalculateDynamicSizing(tt.w, tt.h)

			if s.TerminalWidth != tt.w {
				t.Errorf("TerminalWidth: got %d want %d", s.TerminalWidth, tt.w)
			}
			if s.TerminalHeight != tt.h {
				t.Errorf("TerminalHeight: got %d want %d", s.TerminalHeight, tt.h)
			}
			if s.IsTooSmall != tt.tooSmall {
				t.Errorf("IsTooSmall: got %v want %v", s.IsTooSmall, tt.tooSmall)
			}
			if s.ContentHeight < tt.minContent {
				t.Errorf("ContentHeight %d below minimum %d", s.ContentHeight, tt.minContent)
			}
			if s.ContentInnerWidth < tt.minInner {
				t.Errorf("ContentInnerWidth %d below minimum %d", s.ContentInnerWidth, tt.minInner)
			}
			if s.HeaderInnerWidth != s.ContentInnerWidth {
				t.Errorf("HeaderInnerWidth %d != ContentInnerWidth %d", s.HeaderInnerWidth, s.ContentInnerWidth)
			}
			if s.FooterInnerWidth != tt.w {
				t.Errorf("FooterInnerWidth: got %d want %d", s.FooterInnerWidth, tt.w)
			}
		})
	}
}

func TestCalculateDynamicSizing_ContentHeightFormula(t *testing.T) {
	s := CalculateDynamicSizing(80, 24)
	expected := 24 - HeaderHeight - FooterHeight
	if s.ContentHeight != expected {
		t.Errorf("ContentHeight: got %d want %d", s.ContentHeight, expected)
	}
}

func TestCalculateDynamicSizing_InnerWidthFormula(t *testing.T) {
	s := CalculateDynamicSizing(80, 24)
	expected := 80 - (HorizontalMargin * 2)
	if s.ContentInnerWidth != expected {
		t.Errorf("ContentInnerWidth: got %d want %d", s.ContentInnerWidth, expected)
	}
}

func TestCalculateDynamicSizing_NarrowTerminalClampsInnerWidth(t *testing.T) {
	s := CalculateDynamicSizing(10, 24)
	if s.ContentInnerWidth < 20 {
		t.Errorf("ContentInnerWidth should be clamped to 20, got %d", s.ContentInnerWidth)
	}
}

// --- NewDynamicSizing ---

func TestNewDynamicSizing(t *testing.T) {
	s := NewDynamicSizing()
	ref := CalculateDynamicSizing(80, 24)

	if s.TerminalWidth != ref.TerminalWidth {
		t.Errorf("TerminalWidth: got %d want %d", s.TerminalWidth, ref.TerminalWidth)
	}
	if s.TerminalHeight != ref.TerminalHeight {
		t.Errorf("TerminalHeight: got %d want %d", s.TerminalHeight, ref.TerminalHeight)
	}
	if s.ContentHeight != ref.ContentHeight {
		t.Errorf("ContentHeight: got %d want %d", s.ContentHeight, ref.ContentHeight)
	}
	if s.ContentInnerWidth != ref.ContentInnerWidth {
		t.Errorf("ContentInnerWidth: got %d want %d", s.ContentInnerWidth, ref.ContentInnerWidth)
	}
	if s.IsTooSmall != ref.IsTooSmall {
		t.Errorf("IsTooSmall: got %v want %v", s.IsTooSmall, ref.IsTooSmall)
	}
}

// --- GenerateMenuRows ---

func TestGenerateMenuRows_AlwaysReturns8Rows(t *testing.T) {
	combos := []struct {
		canOpenIDE, canClean, hasBuild, hasBuildsToClean bool
	}{
		{false, false, false, false},
		{true, true, true, true},
		{true, false, true, false},
		{false, true, false, true},
	}
	for _, c := range combos {
		rows := GenerateMenuRows("Xcode", "Debug", c.canOpenIDE, c.canClean, c.hasBuild, c.hasBuildsToClean, false)
		if len(rows) != 8 {
			t.Errorf("expected 8 rows, got %d (combo %+v)", len(rows), c)
		}
	}
}

func TestGenerateMenuRows_AllVisible(t *testing.T) {
	rows := GenerateMenuRows("Ninja", "Release", true, true, true, true, false)
	for _, row := range rows {
		if !row.Visible {
			t.Errorf("row %q should be Visible", row.ID)
		}
	}
}

func TestGenerateMenuRows_SeparatorNotSelectable(t *testing.T) {
	rows := GenerateMenuRows("Xcode", "Debug", true, true, true, true, false)
	sep := rows[3]
	if sep.ID != "separator" {
		t.Fatalf("row[3] expected separator, got %q", sep.ID)
	}
	if sep.IsSelectable {
		t.Error("separator must not be selectable")
	}
}

func TestGenerateMenuRows_SelectabilityByFlags(t *testing.T) {
	tests := []struct {
		name                                    string
		canOpenIDE, canClean, hasBuildsToClean  bool
		wantOpenIDESelectable                   bool
		wantCleanSelectable                     bool
		wantCleanAllSelectable                  bool
	}{
		{"all true", true, true, true, true, true, true},
		{"all false", false, false, false, false, false, false},
		{"openIDE only", true, false, false, true, false, false},
		{"canClean only", false, true, false, false, true, false},
		{"hasBuildsToClean only", false, false, true, false, false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rows := GenerateMenuRows("Xcode", "Debug", tt.canOpenIDE, tt.canClean, false, tt.hasBuildsToClean, false)

			if rows[2].IsSelectable != tt.wantOpenIDESelectable {
				t.Errorf("openIde IsSelectable: got %v want %v", rows[2].IsSelectable, tt.wantOpenIDESelectable)
			}
			if rows[6].IsSelectable != tt.wantCleanSelectable {
				t.Errorf("clean IsSelectable: got %v want %v", rows[6].IsSelectable, tt.wantCleanSelectable)
			}
			if rows[7].IsSelectable != tt.wantCleanAllSelectable {
				t.Errorf("cleanAll IsSelectable: got %v want %v", rows[7].IsSelectable, tt.wantCleanAllSelectable)
			}
		})
	}
}

func TestGenerateMenuRows_RegenerateLabelByHasBuild(t *testing.T) {
	rowsNoBuild := GenerateMenuRows("Xcode", "Debug", false, false, false, false, false)
	if rowsNoBuild[1].Label != "Generate" {
		t.Errorf("hasBuild=false: expected Label 'Generate', got %q", rowsNoBuild[1].Label)
	}

	rowsHasBuild := GenerateMenuRows("Xcode", "Debug", false, false, true, false, false)
	if rowsHasBuild[1].Label != "Regenerate" {
		t.Errorf("hasBuild=true: expected Label 'Regenerate', got %q", rowsHasBuild[1].Label)
	}
}

func TestGenerateMenuRows_RowIDs(t *testing.T) {
	expectedIDs := []string{"project", "regenerate", "openIde", "separator", "configuration", "build", "clean", "cleanAll"}
	rows := GenerateMenuRows("Xcode", "Debug", true, true, true, true, false)

	for i, id := range expectedIDs {
		if rows[i].ID != id {
			t.Errorf("row[%d] ID: got %q want %q", i, rows[i].ID, id)
		}
	}
}

func TestGenerateMenuRows_FixedSelectableRows(t *testing.T) {
	// project, regenerate, configuration, build are always selectable
	rows := GenerateMenuRows("Xcode", "Debug", false, false, false, false, false)

	alwaysSelectable := map[int]string{0: "project", 1: "regenerate", 4: "configuration", 5: "build"}
	for idx, id := range alwaysSelectable {
		if !rows[idx].IsSelectable {
			t.Errorf("row[%d] (%s) should always be selectable", idx, id)
		}
	}
}

// --- OutputBuffer ---

func newTestBuffer() *OutputBuffer {
	return &OutputBuffer{
		maxLines: 10,
		lines:    make([]OutputLine, 0, 10),
	}
}

func TestOutputBuffer_Append(t *testing.T) {
	b := newTestBuffer()

	b.Append("line one", TypeStdout)
	b.Append("line two", TypeStderr)

	if b.GetLineCount() != 2 {
		t.Fatalf("expected 2 lines, got %d", b.GetLineCount())
	}

	lines := b.GetAllLines()
	if lines[0].Text != "line one" {
		t.Errorf("lines[0].Text: got %q want %q", lines[0].Text, "line one")
	}
	if lines[0].Type != TypeStdout {
		t.Errorf("lines[0].Type: got %q want %q", lines[0].Type, TypeStdout)
	}
	if lines[1].Text != "line two" {
		t.Errorf("lines[1].Text: got %q want %q", lines[1].Text, "line two")
	}
}

func TestOutputBuffer_Append_CircularEviction(t *testing.T) {
	b := newTestBuffer() // maxLines = 10

	for i := 0; i < 12; i++ {
		b.Append("x", TypeStdout)
	}

	if b.GetLineCount() != 10 {
		t.Errorf("buffer should be capped at maxLines=10, got %d", b.GetLineCount())
	}
}

func TestOutputBuffer_ReplaceLast_NonEmpty(t *testing.T) {
	b := newTestBuffer()
	b.Append("original", TypeStdout)
	b.ReplaceLast("replaced", TypeStderr)

	if b.GetLineCount() != 1 {
		t.Fatalf("expected 1 line after ReplaceLast, got %d", b.GetLineCount())
	}
	lines := b.GetAllLines()
	if lines[0].Text != "replaced" {
		t.Errorf("Text: got %q want %q", lines[0].Text, "replaced")
	}
	if lines[0].Type != TypeStderr {
		t.Errorf("Type: got %q want %q", lines[0].Type, TypeStderr)
	}
}

func TestOutputBuffer_ReplaceLast_EmptyBuffer(t *testing.T) {
	b := newTestBuffer()
	b.ReplaceLast("first", TypeStatus)

	if b.GetLineCount() != 1 {
		t.Fatalf("ReplaceLast on empty buffer should append, got count %d", b.GetLineCount())
	}
	lines := b.GetAllLines()
	if lines[0].Text != "first" {
		t.Errorf("Text: got %q want %q", lines[0].Text, "first")
	}
}

func TestOutputBuffer_GetLines(t *testing.T) {
	b := newTestBuffer()
	for i := 0; i < 5; i++ {
		b.Append("line", TypeStdout)
	}

	tests := []struct {
		start, count, wantLen int
	}{
		{0, 3, 3},
		{2, 10, 3}, // clipped to available
		{5, 2, 0},  // out of range
		{-1, 3, 3}, // negative start clamped to 0
	}
	for _, tt := range tests {
		result := b.GetLines(tt.start, tt.count)
		if len(result) != tt.wantLen {
			t.Errorf("GetLines(%d,%d): got len %d want %d", tt.start, tt.count, len(result), tt.wantLen)
		}
	}
}

func TestOutputBuffer_GetAllLines_ReturnsCopy(t *testing.T) {
	b := newTestBuffer()
	b.Append("original", TypeStdout)

	lines := b.GetAllLines()
	lines[0].Text = "mutated"

	fresh := b.GetAllLines()
	if fresh[0].Text != "original" {
		t.Error("GetAllLines should return a copy, not a reference to internal slice")
	}
}

func TestOutputBuffer_GetSnapshot(t *testing.T) {
	b := newTestBuffer()
	b.Append("a", TypeStdout)
	b.Append("b", TypeStdout)

	lines, count := b.GetSnapshot()
	if count != 2 {
		t.Errorf("count: got %d want 2", count)
	}
	if len(lines) != count {
		t.Errorf("len(lines) %d != count %d", len(lines), count)
	}
}

func TestOutputBuffer_GetLineCount(t *testing.T) {
	b := newTestBuffer()
	if b.GetLineCount() != 0 {
		t.Errorf("empty buffer: expected 0, got %d", b.GetLineCount())
	}
	b.Append("x", TypeStdout)
	if b.GetLineCount() != 1 {
		t.Errorf("after 1 append: expected 1, got %d", b.GetLineCount())
	}
}

func TestOutputBuffer_Clear(t *testing.T) {
	b := newTestBuffer()
	b.Append("a", TypeStdout)
	b.Append("b", TypeStdout)
	b.Clear()

	if b.GetLineCount() != 0 {
		t.Errorf("after Clear: expected 0, got %d", b.GetLineCount())
	}
}

func TestOutputBuffer_ConcurrentAppendGetSnapshot(t *testing.T) {
	b := newTestBuffer()
	var wg sync.WaitGroup
	writers := 5
	reads := 5

	for i := 0; i < writers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 20; j++ {
				b.Append("concurrent", TypeStdout)
			}
		}()
	}

	for i := 0; i < reads; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 20; j++ {
				lines, count := b.GetSnapshot()
				if len(lines) != count {
					t.Errorf("snapshot inconsistency: len(lines)=%d count=%d", len(lines), count)
				}
			}
		}()
	}

	wg.Wait()
}

// --- ConsoleOutState ---

func TestConsoleOutState_Reset(t *testing.T) {
	s := NewConsoleOutState()
	s.ScrollOffset = 5
	s.MaxScroll = 10
	s.Reset()

	if s.ScrollOffset != 0 {
		t.Errorf("ScrollOffset: got %d want 0", s.ScrollOffset)
	}
	if s.MaxScroll != 0 {
		t.Errorf("MaxScroll: got %d want 0", s.MaxScroll)
	}
}

func TestConsoleOutState_ScrollUp(t *testing.T) {
	s := NewConsoleOutState()
	s.ScrollOffset = 3
	s.ScrollUp()

	if s.ScrollOffset != 2 {
		t.Errorf("ScrollOffset after ScrollUp: got %d want 2", s.ScrollOffset)
	}
}

func TestConsoleOutState_ScrollUp_AtZero(t *testing.T) {
	s := NewConsoleOutState()
	s.ScrollOffset = 0
	s.ScrollUp()

	if s.ScrollOffset != 0 {
		t.Errorf("ScrollOffset should not go below 0, got %d", s.ScrollOffset)
	}
}

func TestConsoleOutState_ScrollDown(t *testing.T) {
	s := NewConsoleOutState()
	s.MaxScroll = 5
	s.ScrollOffset = 3
	s.ScrollDown()

	if s.ScrollOffset != 4 {
		t.Errorf("ScrollOffset after ScrollDown: got %d want 4", s.ScrollOffset)
	}
}

func TestConsoleOutState_ScrollDown_AtMax(t *testing.T) {
	s := NewConsoleOutState()
	s.MaxScroll = 5
	s.ScrollOffset = 5
	s.ScrollDown()

	if s.ScrollOffset != 5 {
		t.Errorf("ScrollOffset should not exceed MaxScroll, got %d", s.ScrollOffset)
	}
}

func TestConsoleOutState_ScrollBounds(t *testing.T) {
	s := NewConsoleOutState()
	s.MaxScroll = 3

	for i := 0; i < 10; i++ {
		s.ScrollDown()
	}
	if s.ScrollOffset != 3 {
		t.Errorf("ScrollOffset capped at MaxScroll=3, got %d", s.ScrollOffset)
	}

	for i := 0; i < 10; i++ {
		s.ScrollUp()
	}
	if s.ScrollOffset != 0 {
		t.Errorf("ScrollOffset floored at 0, got %d", s.ScrollOffset)
	}
}

// --- ConfirmationDialog ---

func testConfig() ConfirmationConfig {
	return ConfirmationConfig{
		Title:       "Confirm",
		Explanation: "Are you sure?",
		YesLabel:    "Yes",
		NoLabel:     "No",
		ActionID:    "test_action",
	}
}

func TestConfirmationDialog_DefaultSelectedButton(t *testing.T) {
	d := NewConfirmationDialog(testConfig(), 80, nil)
	if d.SelectedButton != ButtonYes {
		t.Errorf("default button: got %q want %q", d.SelectedButton, ButtonYes)
	}
}

func TestConfirmationDialog_SelectYes(t *testing.T) {
	d := NewConfirmationDialog(testConfig(), 80, nil)
	d.SelectNo()
	d.SelectYes()

	if d.GetSelectedButton() != ButtonYes {
		t.Errorf("GetSelectedButton: got %q want %q", d.GetSelectedButton(), ButtonYes)
	}
}

func TestConfirmationDialog_SelectNo(t *testing.T) {
	d := NewConfirmationDialog(testConfig(), 80, nil)
	d.SelectNo()

	if d.GetSelectedButton() != ButtonNo {
		t.Errorf("GetSelectedButton: got %q want %q", d.GetSelectedButton(), ButtonNo)
	}
}

func TestNewConfirmationDialogWithDefault_Yes(t *testing.T) {
	d := NewConfirmationDialogWithDefault(testConfig(), 80, nil, ButtonYes)
	if d.SelectedButton != ButtonYes {
		t.Errorf("expected ButtonYes, got %q", d.SelectedButton)
	}
}

func TestNewConfirmationDialogWithDefault_No(t *testing.T) {
	d := NewConfirmationDialogWithDefault(testConfig(), 80, nil, ButtonNo)
	if d.SelectedButton != ButtonNo {
		t.Errorf("expected ButtonNo, got %q", d.SelectedButton)
	}
}

func TestConfirmationDialog_ActiveDefaultsFalse(t *testing.T) {
	d := NewConfirmationDialog(testConfig(), 80, nil)
	if d.Active {
		t.Error("Active should default to false")
	}
}

func TestConfirmationDialog_ContextInitialized(t *testing.T) {
	d := NewConfirmationDialog(testConfig(), 80, nil)
	if d.Context == nil {
		t.Error("Context map should be initialized, not nil")
	}
}

// --- countDisplayLines (white-box) ---

func makeOutputLines(n int) []OutputLine {
	lines := make([]OutputLine, n)
	for i := range lines {
		lines[i] = OutputLine{Time: "00:00:00", Type: TypeStdout, Text: "line"}
	}
	return lines
}

func TestCountDisplayLines_Empty(t *testing.T) {
	result := countDisplayLines([]OutputLine{}, 80)
	if result != 0 {
		t.Errorf("empty buffer: got %d want 0", result)
	}
}

func TestCountDisplayLines_ShortLines(t *testing.T) {
	// Each entry "[00:00:00] line" is well under 80 chars → 1 display line each
	lines := makeOutputLines(5)
	result := countDisplayLines(lines, 80)
	if result != 5 {
		t.Errorf("5 short lines: got %d want 5", result)
	}
}

// --- applyScrollState (white-box) ---

func TestApplyScrollState_AutoScroll(t *testing.T) {
	s := &ConsoleOutState{}
	applyScrollState(s, 20, 10, true)

	if s.ScrollOffset != s.MaxScroll {
		t.Errorf("autoScroll: ScrollOffset %d should equal MaxScroll %d", s.ScrollOffset, s.MaxScroll)
	}
	if s.MaxScroll != 10 {
		t.Errorf("MaxScroll: got %d want 10", s.MaxScroll)
	}
}

func TestApplyScrollState_NoAutoScroll_OffsetClamped(t *testing.T) {
	s := &ConsoleOutState{ScrollOffset: 100}
	applyScrollState(s, 20, 10, false)

	if s.ScrollOffset > s.MaxScroll {
		t.Errorf("ScrollOffset %d should be clamped to MaxScroll %d", s.ScrollOffset, s.MaxScroll)
	}
}

func TestApplyScrollState_NoAutoScroll_NegativeOffsetClamped(t *testing.T) {
	s := &ConsoleOutState{ScrollOffset: -5}
	applyScrollState(s, 20, 10, false)

	if s.ScrollOffset < 0 {
		t.Errorf("ScrollOffset should be >= 0, got %d", s.ScrollOffset)
	}
}

func TestApplyScrollState_MaxScrollZeroWhenFewLines(t *testing.T) {
	s := &ConsoleOutState{}
	applyScrollState(s, 3, 10, false)

	if s.MaxScroll != 0 {
		t.Errorf("MaxScroll should be 0 when totalLines < contentHeight, got %d", s.MaxScroll)
	}
}

func TestApplyScrollState_MaxScrollCalculation(t *testing.T) {
	tests := []struct {
		totalDisplayLines int
		contentHeight     int
		wantMaxScroll     int
	}{
		{20, 10, 10},
		{10, 10, 0},
		{5, 10, 0},
		{1, 10, 0},
		{11, 10, 1},
	}
	for _, tt := range tests {
		s := &ConsoleOutState{}
		applyScrollState(s, tt.totalDisplayLines, tt.contentHeight, false)
		if s.MaxScroll != tt.wantMaxScroll {
			t.Errorf("total=%d height=%d: MaxScroll got %d want %d",
				tt.totalDisplayLines, tt.contentHeight, s.MaxScroll, tt.wantMaxScroll)
		}
	}
}
