package text

import "code.google.com/p/freetype-go/freetype/truetype"

// lRSweeper represents a left-to-right contour point sweeper. It sweeps point
// by point after sorting the input points in left-to-right order.
type lRSweeper struct {
	points []truetype.Point
}
