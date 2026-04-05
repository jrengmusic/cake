package banner

import "fmt"

const brailleBrightnessThreshold = 50

// BrailleChar represents a braille character with its color
type BrailleChar struct {
	Char  rune
	Color Color
}


// collectBrailleCell gathers the 2×4 pixel block at (x, y) and returns the
// braille dot-index bitmask and the list of pixel colors.
func collectBrailleCell(canvas [][]PixelColor, x, y, width, height int) (int, []PixelColor) {
	brailleIndex := 0
	var colors []PixelColor

	// Braille dots ordered as: 0,2,4,6 (left col) then 1,3,5,7 (right col)
	for dy := 0; dy < 4 && y+dy < height; dy++ {
		for dx := 0; dx < 2 && x+dx < width; dx++ {
			pixel := canvas[y+dy][x+dx]
			colors = append(colors, pixel)

			brightness := pixel.R + pixel.G + pixel.B
			if brightness > brailleBrightnessThreshold {
				bitIndex := dy
				if dx != 0 {
					bitIndex = dy + 4
				}
				brailleIndex |= (1 << bitIndex)
			}
		}
	}
	return brailleIndex, colors
}

// dominantPixelColor returns the brightest color from a pixel slice.
func dominantPixelColor(colors []PixelColor) Color {
	dominant := Color{0, 0, 0}
	maxBrightness := 0

	for _, c := range colors {
		brightness := c.R + c.G + c.B
		if brightness > maxBrightness {
			dominant = Color{c.R, c.G, c.B}
			maxBrightness = brightness
		}
	}
	return dominant
}

// CanvasToBrailleArray converts pixel canvas to braille character array
// Returns array of rows, each row is an array of BrailleChar
// Braille: 2×4 dots per character
func CanvasToBrailleArray(canvas [][]PixelColor, width, height int) [][]BrailleChar {
	var output [][]BrailleChar

	for y := 0; y < height; y += 4 {
		var row []BrailleChar

		for x := 0; x < width; x += 2 {
			brailleIndex, colors := collectBrailleCell(canvas, x, y, width, height)

			if brailleIndex > 0xFF {
				brailleIndex = 0
			}
			char := rune(0x2800 + brailleIndex)

			row = append(row, BrailleChar{char, dominantPixelColor(colors)})
		}

		output = append(output, row)
	}

	return output
}

// SvgToBrailleArray converts SVG to braille character array
// Each character has RGB color that can be converted to ANSI color codes
func SvgToBrailleArray(svgString string, width, height int) [][]BrailleChar {
	canvas := RenderSvgToBraille(svgString, width, height)
	return CanvasToBrailleArray(canvas, width, height)
}

// RGBToHex converts RGB to hex color string
func RGBToHex(r, g, b int) string {
	return fmt.Sprintf("#%02X%02X%02X", r, g, b)
}

