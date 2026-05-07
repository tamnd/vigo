package vio

import "testing"

func TestPointAddSub(t *testing.T) {
	p := Point{X: 3, Y: 4}.Add(Point{X: 1, Y: 2})
	if p != (Point{X: 4, Y: 6}) {
		t.Fatalf("Add: got %v", p)
	}
	q := Point{X: 3, Y: 4}.Sub(Point{X: 1, Y: 2})
	if q != (Point{X: 2, Y: 2}) {
		t.Fatalf("Sub: got %v", q)
	}
}

func TestRectAccessors(t *testing.T) {
	r := R(2, 3, 5, 7)
	if r.Right() != 7 || r.Bottom() != 10 {
		t.Fatalf("right/bottom: %d %d", r.Right(), r.Bottom())
	}
	if r.IsEmpty() {
		t.Fatal("non-empty rect reported empty")
	}
	for _, r := range []Rect{R(0, 0, 0, 5), R(0, 0, 5, 0), R(0, 0, -1, 1)} {
		if !r.IsEmpty() {
			t.Fatalf("zero/negative rect should be empty: %v", r)
		}
	}
}

func TestRectContains(t *testing.T) {
	r := R(1, 1, 3, 3) // x in [1,4), y in [1,4)
	cases := []struct {
		p    Point
		want bool
	}{
		{Point{1, 1}, true},
		{Point{3, 3}, true},
		{Point{4, 3}, false},
		{Point{3, 4}, false},
		{Point{0, 0}, false},
	}
	for _, c := range cases {
		if got := r.Contains(c.p); got != c.want {
			t.Errorf("Contains(%v) = %v, want %v", c.p, got, c.want)
		}
	}
}

func TestRectIntersect(t *testing.T) {
	a := R(0, 0, 4, 4)
	b := R(2, 2, 4, 4)
	got := a.Intersect(b)
	want := R(2, 2, 2, 2)
	if got != want {
		t.Fatalf("intersect: got %v want %v", got, want)
	}
	if got := a.Intersect(R(10, 10, 1, 1)); !got.IsEmpty() {
		t.Fatalf("disjoint intersect should be empty, got %v", got)
	}
}

func TestRectTranslateInset(t *testing.T) {
	r := R(1, 2, 5, 5).Translate(2, 3)
	if r != R(3, 5, 5, 5) {
		t.Fatalf("translate: %v", r)
	}
	r2 := R(0, 0, 10, 10).Inset(2)
	if r2 != R(2, 2, 6, 6) {
		t.Fatalf("inset: %v", r2)
	}
}
