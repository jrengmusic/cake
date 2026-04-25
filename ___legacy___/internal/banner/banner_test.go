package banner

import (
	"strings"
	"testing"
)

// --- RGBToHex ---

func TestRGBToHex(t *testing.T) {
	cases := []struct {
		r, g, b int
		want    string
	}{
		{0, 0, 0, "#000000"},
		{255, 255, 255, "#FFFFFF"},
		{255, 0, 0, "#FF0000"},
		{0, 255, 0, "#00FF00"},
		{0, 0, 255, "#0000FF"},
		{18, 52, 86, "#123456"},
		{171, 205, 239, "#ABCDEF"},
	}

	for _, tc := range cases {
		got := RGBToHex(tc.r, tc.g, tc.b)
		if got != tc.want {
			t.Errorf("RGBToHex(%d,%d,%d) = %q, want %q", tc.r, tc.g, tc.b, got, tc.want)
		}
	}
}

// --- dominantPixelColor ---

func TestDominantPixelColor_AllBlack(t *testing.T) {
	colors := []PixelColor{{0, 0, 0}, {0, 0, 0}}
	got := dominantPixelColor(colors)
	want := Color{0, 0, 0}
	if got != want {
		t.Errorf("dominantPixelColor(all black) = %v, want %v", got, want)
	}
}

func TestDominantPixelColor_PicksBrightest(t *testing.T) {
	colors := []PixelColor{
		{10, 10, 10},
		{200, 0, 0},
		{50, 50, 50},
	}
	got := dominantPixelColor(colors)
	want := Color{200, 0, 0}
	if got != want {
		t.Errorf("dominantPixelColor = %v, want %v", got, want)
	}
}

func TestDominantPixelColor_Empty(t *testing.T) {
	got := dominantPixelColor(nil)
	want := Color{0, 0, 0}
	if got != want {
		t.Errorf("dominantPixelColor(nil) = %v, want %v", got, want)
	}
}

// --- CanvasToBrailleArray ---

// makeCanvas builds a width x height canvas with all pixels set to the given color.
func makeCanvas(width, height int, pixel PixelColor) [][]PixelColor {
	canvas := make([][]PixelColor, height)
	for y := range canvas {
		canvas[y] = make([]PixelColor, width)
		for x := range canvas[y] {
			canvas[y][x] = pixel
		}
	}
	return canvas
}

func TestCanvasToBrailleArray_Dimensions(t *testing.T) {
	// 4 wide x 8 tall canvas => 2 braille cols, 2 braille rows
	canvas := makeCanvas(4, 8, PixelColor{0, 0, 0})
	result := CanvasToBrailleArray(canvas, 4, 8)

	wantRows := 2
	wantCols := 2

	if len(result) != wantRows {
		t.Errorf("CanvasToBrailleArray rows = %d, want %d", len(result), wantRows)
	}

	for i, row := range result {
		if len(row) != wantCols {
			t.Errorf("CanvasToBrailleArray row[%d] cols = %d, want %d", i, len(row), wantCols)
		}
	}
}

func TestCanvasToBrailleArray_AllBlackIsEmptyBraille(t *testing.T) {
	// All-black pixels: brightness = 0, all braille dots off => 0x2800 (blank braille)
	canvas := makeCanvas(2, 4, PixelColor{0, 0, 0})
	result := CanvasToBrailleArray(canvas, 2, 4)

	if len(result) != 1 || len(result[0]) != 1 {
		t.Fatalf("unexpected dimensions: %dx%d", len(result), len(result[0]))
	}

	got := result[0][0].Char
	want := rune(0x2800)
	if got != want {
		t.Errorf("all-black canvas char = U+%04X, want U+2800", got)
	}
}

func TestCanvasToBrailleArray_AllBrightIsFull(t *testing.T) {
	// All bright pixels (brightness > threshold): all 8 dots set => 0x28FF
	canvas := makeCanvas(2, 4, PixelColor{100, 100, 100})
	result := CanvasToBrailleArray(canvas, 2, 4)

	if len(result) != 1 || len(result[0]) != 1 {
		t.Fatalf("unexpected dimensions: %dx%d", len(result), len(result[0]))
	}

	got := result[0][0].Char
	want := rune(0x28FF)
	if got != want {
		t.Errorf("all-bright canvas char = U+%04X, want U+28FF", got)
	}
}

func TestCanvasToBrailleArray_ColorPropagated(t *testing.T) {
	// Dominant pixel should propagate to BrailleChar.Color
	canvas := makeCanvas(2, 4, PixelColor{200, 100, 50})
	result := CanvasToBrailleArray(canvas, 2, 4)

	if len(result) == 0 || len(result[0]) == 0 {
		t.Fatal("empty result")
	}

	got := result[0][0].Color
	want := Color{200, 100, 50}
	if got != want {
		t.Errorf("BrailleChar.Color = %v, want %v", got, want)
	}
}

func TestCanvasToBrailleArray_OddDimensions(t *testing.T) {
	// 3 wide x 5 tall: cols = ceil(3/2)=2, rows = ceil(5/4)=2
	canvas := makeCanvas(3, 5, PixelColor{0, 0, 0})
	result := CanvasToBrailleArray(canvas, 3, 5)

	wantRows := 2
	wantCols := 2

	if len(result) != wantRows {
		t.Errorf("rows = %d, want %d", len(result), wantRows)
	}

	for i, row := range result {
		if len(row) != wantCols {
			t.Errorf("row[%d] cols = %d, want %d", i, len(row), wantCols)
		}
	}
}

// --- parseColor (white-box, same package) ---

func TestParseColor_RGB(t *testing.T) {
	cases := []struct {
		input string
		want  Color
	}{
		{"fill:rgb(255,0,0)", Color{255, 0, 0}},
		{"fill:rgb(0, 128, 64)", Color{0, 128, 64}},
		{"fill:rgb(0,0,0)", Color{0, 0, 0}},
		{"fill:rgb(255,255,255)", Color{255, 255, 255}},
	}

	for _, tc := range cases {
		got := parseColor(tc.input)
		if got != tc.want {
			t.Errorf("parseColor(%q) = %v, want %v", tc.input, got, tc.want)
		}
	}
}

func TestParseColor_Hex(t *testing.T) {
	cases := []struct {
		input string
		want  Color
	}{
		{"fill:#FF0000", Color{255, 0, 0}},
		{"fill:#000000", Color{0, 0, 0}},
		{"fill:#FFFFFF", Color{255, 255, 255}},
		{"fill:#1a2b3c", Color{26, 43, 60}},
	}

	for _, tc := range cases {
		got := parseColor(tc.input)
		if got != tc.want {
			t.Errorf("parseColor(%q) = %v, want %v", tc.input, got, tc.want)
		}
	}
}

func TestParseColor_DefaultWhite(t *testing.T) {
	got := parseColor("fill:none")
	want := Color{255, 255, 255}
	if got != want {
		t.Errorf("parseColor(unrecognized) = %v, want %v (default white)", got, want)
	}
}

// --- normalizePathData (white-box, same package) ---

func TestNormalizePathData_Commands(t *testing.T) {
	cases := []struct {
		input    string
		wantHas  []string
	}{
		{
			"M10,20L30,40Z",
			[]string{"M", "10", "20", "L", "30", "40", "Z"},
		},
		{
			"M 10 20 C 1 2 3 4 5 6 Z",
			[]string{"M", "10", "20", "C", "1", "2", "3", "4", "5", "6", "Z"},
		},
	}

	for _, tc := range cases {
		tokens := normalizePathData(tc.input)
		tokenSet := make(map[string]bool, len(tokens))
		for _, tok := range tokens {
			tokenSet[tok] = true
		}
		for _, expected := range tc.wantHas {
			if !tokenSet[expected] {
				t.Errorf("normalizePathData(%q): missing token %q in %v", tc.input, expected, tokens)
			}
		}
	}
}

func TestNormalizePathData_NoExtraWhitespace(t *testing.T) {
	tokens := normalizePathData("  M  10  ,  20  ")
	for _, tok := range tokens {
		if strings.TrimSpace(tok) != tok || tok == "" {
			t.Errorf("normalizePathData produced whitespace or empty token: %q", tok)
		}
	}
}

// --- SvgToBrailleArray ---

const minimalSVG = `<svg xmlns="http://www.w3.org/2000/svg" width="100" height="100">
<path d="M 10 10 L 90 10 L 90 90 L 10 90 Z" style="fill:rgb(255,0,0);"/>
</svg>`

func TestSvgToBrailleArray_NonEmpty(t *testing.T) {
	result := SvgToBrailleArray(minimalSVG, 8, 8)

	if len(result) == 0 {
		t.Fatal("SvgToBrailleArray returned empty result")
	}

	for i, row := range result {
		if len(row) == 0 {
			t.Errorf("SvgToBrailleArray row[%d] is empty", i)
		}
	}
}

func TestSvgToBrailleArray_Dimensions(t *testing.T) {
	// 8x8 canvas => 4 braille cols, 2 braille rows
	result := SvgToBrailleArray(minimalSVG, 8, 8)

	wantRows := 2
	wantCols := 4

	if len(result) != wantRows {
		t.Errorf("SvgToBrailleArray rows = %d, want %d", len(result), wantRows)
	}

	for i, row := range result {
		if len(row) != wantCols {
			t.Errorf("SvgToBrailleArray row[%d] cols = %d, want %d", i, len(row), wantCols)
		}
	}
}

func TestSvgToBrailleArray_ValidBrailleRange(t *testing.T) {
	result := SvgToBrailleArray(minimalSVG, 8, 8)

	for y, row := range result {
		for x, cell := range row {
			if cell.Char < 0x2800 || cell.Char > 0x28FF {
				t.Errorf("cell[%d][%d].Char = U+%04X outside braille range [U+2800..U+28FF]", y, x, cell.Char)
			}
		}
	}
}

func TestSvgToBrailleArray_EmptySVG(t *testing.T) {
	// SVG with no paths: canvas stays all-black, output should still have correct dimensions
	emptySVG := `<svg xmlns="http://www.w3.org/2000/svg" width="100" height="100"></svg>`
	result := SvgToBrailleArray(emptySVG, 4, 4)

	wantRows := 1
	wantCols := 2

	if len(result) != wantRows {
		t.Errorf("empty SVG rows = %d, want %d", len(result), wantRows)
	}

	if len(result) > 0 && len(result[0]) != wantCols {
		t.Errorf("empty SVG cols = %d, want %d", len(result[0]), wantCols)
	}
}
