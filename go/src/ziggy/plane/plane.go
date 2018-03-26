package plane

type Plane interface {
	Converge(radius int)
	GetCell(x, y int) Cell
	GetWidth() int
	GetHeight() int
	AddPoint(PointInterface, int, int)
}

// A plane is a collection of cells
type plane struct {
	cells  [][]Cell
	width  int
	height int
}

func NewPlane(width, height int) Plane {
	// No Locks here as the plane doesn't exist yet...
	p := &plane{
		width:  width,
		height: height,
	}
	for x := 0; x < width; x++ {
		p.cells = append(p.cells, make([]Cell, height))
		for y := range p.cells[x] {
			p.cells[x][y] = NewCell(p, x, y)
		}
	}
	return p
}

// Converge moves the points on the plane to their best neighbors
// The radius sets the number of cells a point can move on the plane in one iteration
// of convergence (see cell.GetBestNeighbor)
func (p *plane) Converge(radius int) {
	var adjusted bool
	for x := range p.cells {
		for y := range p.cells[x] {
			oldCell := p.GetCell(x, y)

			newX, newY := oldCell.GetBestNeighbor(radius)

			if newX == -1 && newY == -1 {
				continue // there is no better neighbor
			}

			if oldCell.NumPoints() == 0 {
				continue
			}

			// The oldCell has a better neighbor, AND it has points in it that need moving!
			adjusted = true
			newCell := p.GetCell(newX, newY)
			movingPoints := oldCell.DrainPoints()
			for i := range movingPoints {
				// Update the locations on the points
				movingPoints[i].SetLocation(newX, newY)
			}
			newCell.AddPoints(movingPoints...)
		}
	}

	if adjusted {
		// Recurse until we are no longer adjusting
		p.Converge(radius)
	}
}

func (p *plane) GetCell(x, y int) Cell {
	return p.cells[x][y]
}

func (p *plane) GetWidth() int {
	return p.width
}

func (p *plane) GetHeight() int {
	return p.height
}

func (p *plane) AddPoint(point PointInterface, x, y int) {
	p.GetCell(x, y).AddPoints(point)
}
