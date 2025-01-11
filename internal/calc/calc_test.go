package calc_test

import (
	"testing"
	"trucks/internal/calc"
	"trucks/internal/models"
)

func TestDistanceZero(t *testing.T) {
	p1 := models.NewPoint(0, 0)
	if calc.Distance(p1, p1) != 0.0 {
		t.Fatalf("Distance not zero between same point")
	}
	p2 := models.NewPoint(5, 5)
	if calc.Distance(p2, p2) != 0.0 {
		t.Fatalf("Distance not zero between same point")
	}

}

func TestDistance345(t *testing.T) {
	d := calc.Distance(models.NewPoint(0, 0), models.NewPoint(3, 4))
	if d != 5.0 {
		t.Fatalf("Failed to make 3 4 5 triangle")
	}
}
