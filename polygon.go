// Also added other functions and some tests related to geo based polygons.

package geo

import (
	"math"
)

// A Polygon is carved out of a 2D plane by a set of (possibly disjoint) contours.
// It can thus contain holes, and can be self-intersecting.
type Polygon struct {
	points []Point
}

// NewPolygon: Creates and returns a new pointer to a Polygon
// composed of the passed in points.  Points are
// considered to be in order such that the last point
// forms an edge with the first point.
func NewPolygon(points []Point) Polygon {
	return Polygon{points: points}
}

// Points returns the points of the current Polygon.
func (p Polygon) Points() []Point {
	return p.points
}

// Add: Appends the passed in contour to the current Polygon and returns
// a new polygon.
func (p Polygon) Add(point Point) Polygon {
	p.points = append(p.points, point)
	return p
}

// IsClosed returns whether or not the polygon is closed.
// TODO:  This can obviously be improved, but for now,
//
//	this should be sufficient for detecting if points
//	are contained using the raycast algorithm.
func (p Polygon) IsClosed() bool {
	if len(p.points) < 3 {
		return false
	}

	return true
}

// Contains returns whether or not the current Polygon contains the passed in Point.
func (p Polygon) Contains(point Point) bool {
	if !p.IsClosed() {
		return false
	}
	// Look here for further options: https://github.com/kellydunn/golang-geo/pull/71#discussion_r303040014
	for _, p := range p.points {
		if p.lat == point.lat && p.lng == point.lng {
			return true
		}
	}

	start := len(p.points) - 1
	end := 0

	contains := p.intersectsWithRaycast(point, &p.points[start], &p.points[end])

	for i := 1; i < len(p.points); i++ {
		if p.intersectsWithRaycast(point, &p.points[i-1], &p.points[i]) {
			contains = !contains
		}
	}

	return contains
}

// Using the raycast algorithm, this returns whether or not the passed in point
// Intersects with the edge drawn by the passed in start and end points.
// Original implementation: http://rosettacode.org/wiki/Ray-casting_algorithm#Go although
// this implementation has bugs if the x point is equal to the x of the start.
// As far as I can tell, the ray that is being cast is going straight up.
func (p Polygon) intersectsWithRaycast(point Point, start *Point, end *Point) bool {
	// Always ensure that the the first point
	// has a y coordinate that is less than the second point
	if start.lat > end.lat {

		// Switch the points if otherwise.
		start, end = end, start

	}

	for point.lat == start.lat || point.lat == end.lat {
		newLat := math.Nextafter(point.lat, math.Inf(1))
		point = NewPoint(newLat, point.lng)
	}

	// If we are outside of the polygon, indicate so.
	if point.lat < start.lat || point.lat > end.lat {
		return false
	}

	if start.lng > end.lng {
		if point.lng > start.lng {
			return false
		}
		if point.lng < end.lng {
			return true
		}

	} else {
		if point.lng > end.lng {
			return false
		}
		if point.lng < start.lng {
			return true
		}
	}

	// this is here to prevent points that share the same longitude value as the start.
	// Unfortunately this also breaks points that are on the same line as the start point that are within the shape
	if point.lng == start.lng {
		return false
	}

	raySlope := (point.lat - start.lat) / (point.lng - start.lng)
	diagSlope := (end.lat - start.lat) / (end.lng - start.lng)
	return raySlope >= diagSlope
}
