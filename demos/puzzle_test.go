package demos

import (
	"testing"

	"github.com/tamnd/vigo/event"
)

func TestPuzzleSolvedStartsSolved(t *testing.T) {
	p := NewPuzzleSolved(PuzzleDefaultBounds)
	hostFor(p)
	if !p.Solved() {
		t.Fatalf("solved fixture should be solved; tiles=%v", p.board.tiles)
	}
}

func TestPuzzleOneAwayNeedsOneMove(t *testing.T) {
	p := NewPuzzleOneAway(PuzzleDefaultBounds)
	hostFor(p)

	if p.Solved() {
		t.Fatal("one-away board should not be solved")
	}
	// Blank is at (2,3); arrow Right moves blank to (3,3), bringing 15 left.
	p.HandleEvent(keyOf(event.KeyArrowRight))
	if !p.Solved() {
		t.Fatalf("expected solved after ArrowRight; tiles=%v", p.board.tiles)
	}
	if p.Moves() != 1 {
		t.Fatalf("Moves=%d want 1", p.Moves())
	}
}

func TestPuzzleSlideOffBoardIgnored(t *testing.T) {
	p := NewPuzzleSolved(PuzzleDefaultBounds)
	hostFor(p)

	// Blank is at (3,3): Right and Down would leave the board.
	p.HandleEvent(keyOf(event.KeyArrowRight))
	p.HandleEvent(keyOf(event.KeyArrowDown))
	if p.Moves() != 0 {
		t.Fatalf("off-board slides should not count: moves=%d", p.Moves())
	}
}

func TestPuzzleResetClearsMoves(t *testing.T) {
	p := NewPuzzleOneAway(PuzzleDefaultBounds)
	hostFor(p)

	p.HandleEvent(keyOf(event.KeyArrowRight))
	if p.Moves() != 1 {
		t.Fatalf("setup: moves=%d", p.Moves())
	}
	p.HandleEvent(keyRune('R'))
	if p.Moves() != 0 {
		t.Fatalf("after R moves=%d want 0", p.Moves())
	}
}

func TestPuzzleStatusReflectsSolved(t *testing.T) {
	p := NewPuzzleOneAway(PuzzleDefaultBounds)
	hostFor(p)

	p.HandleEvent(keyOf(event.KeyArrowRight))
	if p.status.Text != "Moves:   1  SOLVED" {
		t.Fatalf("status=%q", p.status.Text)
	}
}

func TestPuzzleRandomShuffleIsValid(t *testing.T) {
	tiles := randomShuffle()
	seen := map[int]bool{}
	for _, v := range tiles {
		if seen[v] {
			t.Fatalf("duplicate tile %d", v)
		}
		seen[v] = true
		if v < 0 || v >= PuzzleSize*PuzzleSize {
			t.Fatalf("tile out of range: %d", v)
		}
	}
}

func TestPuzzleNonArrowKeyIgnored(t *testing.T) {
	p := NewPuzzleSolved(PuzzleDefaultBounds)
	hostFor(p)

	before := p.board.tiles
	p.HandleEvent(keyRune('q'))
	if p.board.tiles != before {
		t.Fatal("unrelated key changed the board")
	}
}

func TestPuzzleClickAdjacentTileSlides(t *testing.T) {
	p := NewPuzzleOneAway(PuzzleDefaultBounds)
	hostFor(p)

	// Blank sits at (2,3); tile 15 is the right-side neighbor at (3,3).
	x := p.board.Bounds.X + 3*4 + 1
	y := p.board.Bounds.Y + 3
	p.HandleEvent(&event.Event{
		What:  event.ClassMouseDown,
		Mouse: event.MouseEvent{X: x, Y: y, Buttons: event.MouseLeft},
	})
	if !p.Solved() {
		t.Fatalf("click on adjacent tile should solve; tiles=%v", p.board.tiles)
	}
	if p.Moves() != 1 {
		t.Fatalf("moves=%d want 1", p.Moves())
	}
}

func TestPuzzleClickNonAdjacentTileIgnored(t *testing.T) {
	p := NewPuzzleSolved(PuzzleDefaultBounds)
	hostFor(p)

	// Blank is at (3,3); click tile (0,0) which is not adjacent.
	x := p.board.Bounds.X + 1
	y := p.board.Bounds.Y
	p.HandleEvent(&event.Event{
		What:  event.ClassMouseDown,
		Mouse: event.MouseEvent{X: x, Y: y, Buttons: event.MouseLeft},
	})
	if p.Moves() != 0 {
		t.Fatalf("non-adjacent click should not slide; moves=%d", p.Moves())
	}
}
