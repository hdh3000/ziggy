package plane

// A point is a thing in a cell on a plane
type Point struct {
	x, y int
}

type PointInterface interface {
	SetLocation(x, y int)
	GetLocation() (int, int)
}

func NewPoint() *Point {
	return &Point{x: -1, y: -1}
}

func (p *Point) SetLocation(x, y int) {
	p.x = x
	p.y = y
}

func (p *Point) GetLocation() (int, int) {
	return p.x, p.y
}
