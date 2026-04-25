package banner

import (
	"math"
	"regexp"
	"sort"
	"strings"
)

const bezierSegmentCount = 20

// extractPathFill extracts subpaths and color from a single path element attribute string
// Returns ok=false if the path has no fill or is missing required attributes
func extractPathFill(pathAttrs string) (subpaths [][]Point, color Color, ok bool) {
	dRe := regexp.MustCompile(`d="([^"]*)`)
	dMatch := dRe.FindStringSubmatch(pathAttrs)
	if len(dMatch) < 2 {
		return nil, Color{}, false
	}

	styleRe := regexp.MustCompile(`style="([^"]*)`)
	styleMatch := styleRe.FindStringSubmatch(pathAttrs)
	if len(styleMatch) < 2 {
		return nil, Color{}, false
	}

	style := styleMatch[1]
	if !strings.Contains(style, "fill:rgb") || strings.Contains(style, "fill:none") {
		return nil, Color{}, false
	}

	return parsePathDataAsSubpaths(dMatch[1]), parseColor(style), true
}

// parseFillPaths extracts all fill paths with colors from SVG
func parseFillPaths(svgString string) []struct {
	Subpaths [][]Point
	Color    Color
} {
	var fillPaths []struct {
		Subpaths [][]Point
		Color    Color
	}

	pathRe := regexp.MustCompile(`(?s)<path\s+([^>]*)>`)
	matches := pathRe.FindAllStringSubmatch(svgString, -1)

	for _, match := range matches {
		if len(match) < 2 {
			continue
		}
		subpaths, color, ok := extractPathFill(match[1])
		if ok {
			fillPaths = append(fillPaths, struct {
				Subpaths [][]Point
				Color    Color
			}{subpaths, color})
		}
	}

	return fillPaths
}

// approximateCubicBezier approximates a cubic Bezier curve as line segments
func approximateCubicBezier(p0, cp1, cp2, p1 Point, segments int) []Point {
	var points []Point

	for i := 0; i <= segments; i++ {
		t := float64(i) / float64(segments)
		mt := 1 - t

		x := mt*mt*mt*p0.X + 3*mt*mt*t*cp1.X + 3*mt*t*t*cp2.X + t*t*t*p1.X
		y := mt*mt*mt*p0.Y + 3*mt*mt*t*cp1.Y + 3*mt*t*t*cp2.Y + t*t*t*p1.Y

		points = append(points, Point{x, y})
	}

	return points
}

// scaleSubpaths applies uniform scale and offset to all subpath points
func scaleSubpaths(allSubpaths [][]Point, scaleX, scaleY, offsetX, offsetY float64) [][]Point {
	var scaledSubpaths [][]Point
	for _, subpath := range allSubpaths {
		var scaledPath []Point
		for _, pt := range subpath {
			scaledPath = append(scaledPath, Point{pt.X*scaleX + offsetX, pt.Y*scaleY + offsetY})
		}
		scaledSubpaths = append(scaledSubpaths, scaledPath)
	}
	return scaledSubpaths
}

// subpathsBounds returns the Y extent of a set of scaled subpaths
func subpathsBounds(scaledSubpaths [][]Point) (minY, maxY float64) {
	minY = math.Inf(1)
	maxY = math.Inf(-1)
	for _, subpath := range scaledSubpaths {
		for _, pt := range subpath {
			if pt.Y < minY {
				minY = pt.Y
			}
			if pt.Y > maxY {
				maxY = pt.Y
			}
		}
	}
	return minY, maxY
}

// collectIntersections finds all edge/scanline intersections for a given Y
func collectIntersections(scaledSubpaths [][]Point, scanYFloat float64) []Intersection {
	var intersections []Intersection
	for _, scaledPoints := range scaledSubpaths {
		for i := 0; i < len(scaledPoints)-1; i++ {
			y1 := scaledPoints[i].Y
			y2 := scaledPoints[i+1].Y
			if (y1 < scanYFloat && y2 >= scanYFloat) || (y2 < scanYFloat && y1 >= scanYFloat) {
				x1 := scaledPoints[i].X
				x2 := scaledPoints[i+1].X
				intersectX := x1 + (scanYFloat-y1)*(x2-x1)/(y2-y1)
				direction := 1
				if y1 >= y2 {
					direction = -1
				}
				intersections = append(intersections, Intersection{int(intersectX), direction})
			}
		}
	}
	sort.Slice(intersections, func(i, j int) bool { return intersections[i].X < intersections[j].X })
	return intersections
}

// applyWindingRule converts sorted intersections to filled ranges using non-zero winding rule
func applyWindingRule(intersections []Intersection, scanY int, color Color) []ScanlineRange {
	var ranges []ScanlineRange
	windingCount := 0
	fillStartX := -1

	for _, inter := range intersections {
		if windingCount != 0 && windingCount+inter.Direction == 0 {
			ranges = append(ranges, ScanlineRange{ScanY: scanY, StartX: fillStartX, EndX: inter.X, Color: color})
		}
		if windingCount == 0 && windingCount+inter.Direction != 0 {
			fillStartX = inter.X
		}
		windingCount += inter.Direction
	}

	return ranges
}

// computeScanlineRanges computes filled pixel ranges using non-zero winding rule
func computeScanlineRanges(allSubpaths [][]Point, color Color, scaleX, scaleY, offsetX, offsetY float64) []ScanlineRange {
	var ranges []ScanlineRange

	if len(allSubpaths) == 0 {
		return ranges
	}

	scaledSubpaths := scaleSubpaths(allSubpaths, scaleX, scaleY, offsetX, offsetY)
	minY, maxY := subpathsBounds(scaledSubpaths)

	for scanY := int(minY); scanY < int(maxY)+1; scanY++ {
		intersections := collectIntersections(scaledSubpaths, float64(scanY))
		ranges = append(ranges, applyWindingRule(intersections, scanY, color)...)
	}

	return ranges
}

// computeScaleAndOffset calculates uniform scale and centering offsets for SVG-to-canvas mapping
func computeScaleAndOffset(svgWidth, svgHeight float64, canvasWidth, canvasHeight int) (scale, offsetX, offsetY float64) {
	scaleX := float64(canvasWidth) / svgWidth
	scaleY := float64(canvasHeight) / svgHeight
	scale = scaleX
	if scaleY < scaleX {
		scale = scaleY
	}
	offsetX = (float64(canvasWidth) - svgWidth*scale) / 2
	offsetY = (float64(canvasHeight) - svgHeight*scale) / 2
	return scale, offsetX, offsetY
}

// drawRangesToCanvas paints scanline ranges onto a pixel canvas
func drawRangesToCanvas(canvas [][]PixelColor, ranges []ScanlineRange, width, height int) {
	for _, r := range ranges {
		for x := r.StartX; x <= r.EndX; x++ {
			if r.ScanY >= 0 && r.ScanY < height && x >= 0 && x < width {
				canvas[r.ScanY][x] = PixelColor{r.Color.R, r.Color.G, r.Color.B}
			}
		}
	}
}

// RenderSvgToBraille renders SVG to a pixel canvas
func RenderSvgToBraille(svgString string, width, height int) [][]PixelColor {
	canvas := make([][]PixelColor, height)
	for i := range canvas {
		canvas[i] = make([]PixelColor, width)
	}

	svgWidth, svgHeight := extractSvgDimensions(svgString)
	scale, offsetX, offsetY := computeScaleAndOffset(svgWidth, svgHeight, width, height)

	for _, fp := range parseFillPaths(svgString) {
		ranges := computeScanlineRanges(fp.Subpaths, fp.Color, scale, scale, offsetX, offsetY)
		drawRangesToCanvas(canvas, ranges, width, height)
	}

	return canvas
}
