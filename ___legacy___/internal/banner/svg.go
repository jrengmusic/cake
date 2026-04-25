package banner

import (
	"regexp"
	"strconv"
	"strings"
)

// Point represents a 2D coordinate
type Point struct {
	X float64
	Y float64
}

// Color represents an RGB color
type Color struct {
	R, G, B int
}

// PixelColor represents a pixel with RGB color
type PixelColor struct {
	R, G, B int
}

// ScanlineRange represents a filled horizontal range
type ScanlineRange struct {
	ScanY  int
	StartX int
	EndX   int
	Color  Color
}

// Intersection represents where an edge crosses a scanline
type Intersection struct {
	X         int
	Direction int
}

// extractSvgDimensions extracts width/height from SVG string
// SVG numeric parsing: malformed values default to zero, producing degraded but safe render
func extractSvgDimensions(svgString string) (width, height float64) {
	// Try explicit width/height attributes
	widthRe := regexp.MustCompile(`width\s*=\s*["']?(\d+(?:\.\d+)?)`)
	heightRe := regexp.MustCompile(`height\s*=\s*["']?(\d+(?:\.\d+)?)`)

	widthMatch := widthRe.FindStringSubmatch(svgString)
	heightMatch := heightRe.FindStringSubmatch(svgString)

	if len(widthMatch) > 1 && len(heightMatch) > 1 {
		w, _ := strconv.ParseFloat(widthMatch[1], 64)
		h, _ := strconv.ParseFloat(heightMatch[1], 64)
		return w, h
	}

	// Fall back to viewBox
	viewBoxRe := regexp.MustCompile(`viewBox\s*=\s*["']([^"']*)["']`)
	viewBoxMatch := viewBoxRe.FindStringSubmatch(svgString)

	if len(viewBoxMatch) > 1 {
		parts := strings.FieldsFunc(viewBoxMatch[1], func(r rune) bool {
			return r == ' ' || r == ','
		})
		if len(parts) >= 4 {
			w, _ := strconv.ParseFloat(parts[2], 64)
			h, _ := strconv.ParseFloat(parts[3], 64)
			return w, h
		}
	}

	// Default fallback
	return 400, 320
}

// handleMoveToCommand handles M/m SVG path commands
func handleMoveToCommand(cmd string, tokens []string, i int, currentPath []Point, currentPos Point, subpaths [][]Point) ([][]Point, []Point, Point, *Point, int) {
	if len(currentPath) > 0 {
		subpaths = append(subpaths, currentPath)
		currentPath = nil
	}
	if i+2 < len(tokens) {
		x, _ := strconv.ParseFloat(tokens[i+1], 64)
		y, _ := strconv.ParseFloat(tokens[i+2], 64)
		if cmd == "M" {
			currentPos = Point{x, y}
		} else {
			currentPos = Point{currentPos.X + x, currentPos.Y + y}
		}
		currentPath = append(currentPath, currentPos)
		return subpaths, currentPath, currentPos, nil, i + 3
	}
	return subpaths, currentPath, currentPos, nil, i + 1
}

// handleLineToCommand handles L/l SVG path commands
func handleLineToCommand(cmd string, tokens []string, i int, currentPath []Point, currentPos Point) ([]Point, Point, *Point, int) {
	if i+2 < len(tokens) {
		x, _ := strconv.ParseFloat(tokens[i+1], 64)
		y, _ := strconv.ParseFloat(tokens[i+2], 64)
		if cmd == "L" {
			currentPos = Point{x, y}
		} else {
			currentPos = Point{currentPos.X + x, currentPos.Y + y}
		}
		currentPath = append(currentPath, currentPos)
		return currentPath, currentPos, nil, i + 3
	}
	return currentPath, currentPos, nil, i + 1
}

// handleCubicBezierCommand handles C/c SVG path commands
func handleCubicBezierCommand(cmd string, tokens []string, i int, currentPath []Point, currentPos Point) ([]Point, Point, *Point, int) {
	if i+6 < len(tokens) {
		cp1x, _ := strconv.ParseFloat(tokens[i+1], 64)
		cp1y, _ := strconv.ParseFloat(tokens[i+2], 64)
		cp2x, _ := strconv.ParseFloat(tokens[i+3], 64)
		cp2y, _ := strconv.ParseFloat(tokens[i+4], 64)
		x, _ := strconv.ParseFloat(tokens[i+5], 64)
		y, _ := strconv.ParseFloat(tokens[i+6], 64)
		var cp1, cp2, endPoint Point
		if cmd == "C" {
			cp1 = Point{cp1x, cp1y}
			cp2 = Point{cp2x, cp2y}
			endPoint = Point{x, y}
		} else {
			cp1 = Point{currentPos.X + cp1x, currentPos.Y + cp1y}
			cp2 = Point{currentPos.X + cp2x, currentPos.Y + cp2y}
			endPoint = Point{currentPos.X + x, currentPos.Y + y}
		}
		curvePoints := approximateCubicBezier(currentPos, cp1, cp2, endPoint, bezierSegmentCount)
		if len(curvePoints) > 1 {
			currentPath = append(currentPath, curvePoints[1:]...)
		}
		return currentPath, endPoint, &cp2, i + 7
	}
	return currentPath, currentPos, nil, i + 1
}

// handleSmoothCubicBezierCommand handles S/s SVG path commands
func handleSmoothCubicBezierCommand(cmd string, tokens []string, i int, currentPath []Point, currentPos Point, lastControlPoint *Point) ([]Point, Point, *Point, int) {
	if i+4 < len(tokens) {
		cp2x, _ := strconv.ParseFloat(tokens[i+1], 64)
		cp2y, _ := strconv.ParseFloat(tokens[i+2], 64)
		x, _ := strconv.ParseFloat(tokens[i+3], 64)
		y, _ := strconv.ParseFloat(tokens[i+4], 64)
		var cp1 Point
		if lastControlPoint != nil {
			cp1 = Point{2*currentPos.X - lastControlPoint.X, 2*currentPos.Y - lastControlPoint.Y}
		} else {
			cp1 = currentPos
		}
		var cp2, endPoint Point
		if cmd == "S" {
			cp2 = Point{cp2x, cp2y}
			endPoint = Point{x, y}
		} else {
			cp2 = Point{currentPos.X + cp2x, currentPos.Y + cp2y}
			endPoint = Point{currentPos.X + x, currentPos.Y + y}
		}
		curvePoints := approximateCubicBezier(currentPos, cp1, cp2, endPoint, bezierSegmentCount)
		if len(curvePoints) > 1 {
			currentPath = append(currentPath, curvePoints[1:]...)
		}
		return currentPath, endPoint, &cp2, i + 5
	}
	return currentPath, currentPos, nil, i + 1
}

// handleClosePathCommand handles Z/z SVG path commands
func handleClosePathCommand(currentPath []Point) []Point {
	if len(currentPath) > 0 {
		first := currentPath[0]
		last := currentPath[len(currentPath)-1]
		if first.X != last.X || first.Y != last.Y {
			currentPath = append(currentPath, first)
		}
	}
	return currentPath
}

// normalizePathData tokenizes and normalizes SVG path data string
func normalizePathData(pathData string) []string {
	normalized := regexp.MustCompile(`([MLCSZmlcsz])`).ReplaceAllString(pathData, " $1 ")
	normalized = strings.ReplaceAll(normalized, ",", " ")
	normalized = regexp.MustCompile(`\s+`).ReplaceAllString(normalized, " ")
	normalized = strings.TrimSpace(normalized)
	return strings.Fields(normalized)
}

// parsePathDataAsSubpaths parses SVG path "d" attribute into subpaths
// SVG numeric parsing: malformed values default to zero, producing degraded but safe render
func parsePathDataAsSubpaths(pathData string) [][]Point {
	var subpaths [][]Point
	tokens := normalizePathData(pathData)

	var currentPath []Point
	currentPos := Point{0, 0}
	var lastControlPoint *Point
	i := 0

	for i < len(tokens) {
		cmd := tokens[i]
		switch cmd {
		case "M", "m":
			subpaths, currentPath, currentPos, lastControlPoint, i = handleMoveToCommand(cmd, tokens, i, currentPath, currentPos, subpaths)
		case "L", "l":
			currentPath, currentPos, lastControlPoint, i = handleLineToCommand(cmd, tokens, i, currentPath, currentPos)
		case "C", "c":
			currentPath, currentPos, lastControlPoint, i = handleCubicBezierCommand(cmd, tokens, i, currentPath, currentPos)
		case "S", "s":
			currentPath, currentPos, lastControlPoint, i = handleSmoothCubicBezierCommand(cmd, tokens, i, currentPath, currentPos, lastControlPoint)
		case "Z", "z":
			currentPath = handleClosePathCommand(currentPath)
			lastControlPoint = nil
			i++
		default:
			i++
		}
	}

	if len(currentPath) > 0 {
		subpaths = append(subpaths, currentPath)
	}

	return subpaths
}

// parseColor extracts RGB color from SVG style attribute
// SVG numeric parsing: malformed values default to zero, producing degraded but safe render
func parseColor(styleAttr string) Color {
	// Try RGB format: rgb(255, 0, 0)
	rgbRe := regexp.MustCompile(`rgb\((\d+),\s*(\d+),\s*(\d+)\)`)
	rgbMatch := rgbRe.FindStringSubmatch(styleAttr)
	if len(rgbMatch) == 4 {
		r, _ := strconv.Atoi(rgbMatch[1])
		g, _ := strconv.Atoi(rgbMatch[2])
		b, _ := strconv.Atoi(rgbMatch[3])
		return Color{r, g, b}
	}

	// Try HEX format: #RRGGBB
	hexRe := regexp.MustCompile(`#([0-9A-Fa-f]{6})`)
	hexMatch := hexRe.FindStringSubmatch(styleAttr)
	if len(hexMatch) == 2 {
		hex := hexMatch[1]
		r, _ := strconv.ParseInt(hex[0:2], 16, 64)
		g, _ := strconv.ParseInt(hex[2:4], 16, 64)
		b, _ := strconv.ParseInt(hex[4:6], 16, 64)
		return Color{int(r), int(g), int(b)}
	}

	// Default white
	return Color{255, 255, 255}
}
