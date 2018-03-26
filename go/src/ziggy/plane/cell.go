package plane

import "sync"

type Cell interface {
	GetBestNeighbor(radius int) (x, y int)

	GetAttraction() float64
	SetAttraction(a float64)

	DrainPoints() []PointInterface
	AddPoints(points ...PointInterface)
	NumPoints() int

	GetXY() (int, int)
}

func NewCell(p Plane, x, y int) Cell {
	return &cell{
		plane: p,
		x:     x,
		y:     y,
	}
}

// A cell is a location on a plane
// It holds points.
type cell struct {
	x, y         int
	attraction   float64
	points       []PointInterface
	plane        Plane
	bestNeighbor *struct{ x, y int }
	lock         sync.Mutex
}

// GetBestNeighbor finds the most attractive neighbor
// It looks at all the neighbours (out to a given radius)
// + + +
// + ! +
// + + +
//
// If the cell is the most attractive thing in the neighborhood (or is tied)
// GetBestNeighbor will return -1, -1
func (c *cell) GetBestNeighbor(radius int) (x, y int) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.bestNeighbor != nil {
		// Don't compute it twice
		// TODO: may want some more sophisticated caching so you can change the radius.
		return c.bestNeighbor.x, c.bestNeighbor.y
	}

	bestX, bestY := -1, -1
	mostAttractive := c.attraction

	for x := c.x - radius; x <= c.x+radius; x++ {
		if x < 0 || x >= c.GetPlane().GetWidth() {
			// If we are off the grid, ignore
			continue
		}
		for y := c.y - radius; y <= c.y+radius; y++ {
			if x == c.x && y == c.y {
				// Don't check yourself, or we will lock
				continue
			}

			if y < 0 || y >= c.GetPlane().GetHeight() {
				continue
			}

			attraction := c.GetPlane().GetCell(x, y).GetAttraction()
			if attraction > mostAttractive {
				bestX, bestY = x, y
				mostAttractive = attraction
			}
		}
	}

	c.bestNeighbor = &struct {
		x, y int
	}{
		x: bestX,
		y: bestY,
	}

	return bestX, bestY

}

func (c *cell) DrainPoints() []PointInterface {
	c.lock.Lock()
	defer c.lock.Unlock()

	points := c.points
	c.points = nil
	return points
}

func (c *cell) AddPoints(points ...PointInterface) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.points = append(c.points, points...)
}

func (c *cell) NumPoints() int {
	c.lock.Lock()
	defer c.lock.Unlock()

	return len(c.points)
}

func (c *cell) GetAttraction() float64 {
	c.lock.Lock()
	defer c.lock.Unlock()

	return c.attraction
}

func (c *cell) SetAttraction(a float64) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.attraction = a
}

func (c *cell) GetXY() (int, int) {
	return c.x, c.y
}

func (c *cell) GetPlane() Plane {
	return c.plane
}
