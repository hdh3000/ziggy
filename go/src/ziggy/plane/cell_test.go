package plane

import (
	"testing"
)

func TestNewPlane(t *testing.T) {
	pl := NewPlane(11, 11)

	pl.GetCell(1, 1).SetAttraction(1)

	pl.AddPoint(NewPoint(), 1, 1)
	pl.AddPoint(NewPoint(), 4, 4)

	if pl.GetCell(0, 0).NumPoints() != 1 {
		t.Errorf("failed to set point in location 0,0")
	}

	pl.Converge(1)

	if pl.GetCell(0, 0).NumPoints() == 1 {
		t.Errorf("failed to remove point from location 0,0")
	}

	if pl.GetCell(1, 1).NumPoints() != 2 {
		t.Errorf("failed to move point to location 1,1")
	}
}
