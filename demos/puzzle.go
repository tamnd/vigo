package demos

import (
	"fmt"
	"math/rand/v2"

	"github.com/tamnd/vigo/event"
	"github.com/tamnd/vigo/view"
	"github.com/tamnd/vigo/vio"
	"github.com/tamnd/vigo/widget"
	"github.com/tamnd/vigo/window"
)

// PuzzleSize is the side length of the puzzle grid.
const PuzzleSize = 4

// Puzzle is the classic 4x4 sliding tile puzzle. Tiles are numbered
// 1..15 with one empty slot (0). Arrow keys slide the neighboring
// tile into the empty slot. R resets to a freshly shuffled board.
type Puzzle struct {
	*window.Window

	board   *puzzleBoard
	status  *widget.StaticText
	moves   int
	shuffle func() [PuzzleSize * PuzzleSize]int
}

// PuzzleDefaultBounds is sized to fit a 4x4 grid of 4-wide cells
// plus the status row.
//
//nolint:gochecknoglobals // immutable default
var PuzzleDefaultBounds = vio.R(2, 2, PuzzleSize*4+4, PuzzleSize+5)

// NewPuzzle returns a Puzzle window with a randomly shuffled board.
func NewPuzzle(bounds vio.Rect) *Puzzle {
	return newPuzzleWith(bounds, randomShuffle)
}

// NewPuzzleSolved returns a Puzzle in the solved 1..15,blank state.
// Used by tests that want a deterministic starting point.
func NewPuzzleSolved(bounds vio.Rect) *Puzzle {
	return newPuzzleWith(bounds, solvedTiles)
}

// NewPuzzleOneAway returns a Puzzle one move away from solved: 15
// and the blank are swapped, so a single ArrowRight slide solves it.
func NewPuzzleOneAway(bounds vio.Rect) *Puzzle {
	return newPuzzleWith(bounds, almostSolvedTiles)
}

func newPuzzleWith(bounds vio.Rect, shuffle func() [PuzzleSize * PuzzleSize]int) *Puzzle {
	w := window.New(bounds, "Puzzle", window.FlagMove|window.FlagClose)
	p := &Puzzle{Window: w, shuffle: shuffle}

	client := w.ClientRect()
	p.board = newPuzzleBoard(vio.R(client.X+1, client.Y+1, PuzzleSize*4, PuzzleSize), shuffle())
	w.Insert(p.board)

	p.status = widget.NewStaticText(
		vio.R(client.X+1, client.Y+PuzzleSize+2, client.W-2, 1),
		"",
	)
	w.Insert(p.status)
	p.refresh()
	return p
}

// Tile returns the value at (col, row). 0 is the blank slot.
func (p *Puzzle) Tile(col, row int) int { return p.board.tiles[row*PuzzleSize+col] }

// Moves returns the number of successful slides since the last reset.
func (p *Puzzle) Moves() int { return p.moves }

// Solved reports whether the tiles are in the canonical 1..15,blank
// arrangement.
func (p *Puzzle) Solved() bool { return p.board.solved() }

// HandleEvent intercepts arrow keys (slide) and 'R' (reset).
func (p *Puzzle) HandleEvent(e *event.Event) {
	if e.What == event.ClassKey && p.handleKey(e) {
		p.refresh()
		e.Clear()
		return
	}
	p.Window.HandleEvent(e)
}

func (p *Puzzle) handleKey(e *event.Event) bool {
	if e.Key.Key == event.KeyRune && (e.Key.Rune == 'r' || e.Key.Rune == 'R') {
		p.reset()
		return true
	}
	dx, dy, ok := arrowDelta(e.Key.Key)
	if !ok {
		return false
	}
	if p.board.slide(dx, dy) {
		p.moves++
	}
	return true
}

// arrowDelta maps an arrow key to the (dx, dy) the BLANK should move.
// The neighbor tile in the opposite direction slides into the blank.
func arrowDelta(k event.Key) (int, int, bool) {
	switch k {
	case event.KeyArrowLeft:
		return -1, 0, true
	case event.KeyArrowRight:
		return 1, 0, true
	case event.KeyArrowUp:
		return 0, -1, true
	case event.KeyArrowDown:
		return 0, 1, true
	}
	return 0, 0, false
}

func (p *Puzzle) reset() {
	p.board.tiles = p.shuffle()
	p.board.findBlank()
	p.moves = 0
}

func (p *Puzzle) refresh() {
	state := "in play"
	if p.Solved() {
		state = "SOLVED"
	}
	p.status.Text = fmt.Sprintf("Moves: %3d  %s", p.moves, state)
}

// puzzleBoard is the inner view that paints the grid.
type puzzleBoard struct {
	*view.View

	tiles    [PuzzleSize * PuzzleSize]int
	blankCol int
	blankRow int
}

func newPuzzleBoard(bounds vio.Rect, tiles [PuzzleSize * PuzzleSize]int) *puzzleBoard {
	v := view.NewView(bounds)
	v.GrowMode = view.GrowFixed
	v.EventMask = 0
	v.PaletteIndex = widget.SlotDialogText
	b := &puzzleBoard{View: v, tiles: tiles}
	b.findBlank()
	return b
}

func (b *puzzleBoard) findBlank() {
	for i, t := range b.tiles {
		if t == 0 {
			b.blankCol = i % PuzzleSize
			b.blankRow = i / PuzzleSize
			return
		}
	}
}

// slide moves the blank by (dx, dy), bringing the neighbor into the
// vacated slot. Returns false if the move would leave the board.
func (b *puzzleBoard) slide(dx, dy int) bool {
	nc, nr := b.blankCol+dx, b.blankRow+dy
	if nc < 0 || nc >= PuzzleSize || nr < 0 || nr >= PuzzleSize {
		return false
	}
	bi := b.blankRow*PuzzleSize + b.blankCol
	ni := nr*PuzzleSize + nc
	b.tiles[bi], b.tiles[ni] = b.tiles[ni], b.tiles[bi]
	b.blankCol, b.blankRow = nc, nr
	return true
}

func (b *puzzleBoard) solved() bool {
	for i := range PuzzleSize*PuzzleSize - 1 {
		if b.tiles[i] != i+1 {
			return false
		}
	}
	return b.tiles[PuzzleSize*PuzzleSize-1] == 0
}

func (b *puzzleBoard) Draw(s *vio.Surface) {
	normal := b.MapColor(b.PaletteIndex)
	highlight := b.MapColor(widget.SlotListSelected)
	s.FillRect(b.Bounds, ' ', normal)

	for row := range PuzzleSize {
		for col := range PuzzleSize {
			t := b.tiles[row*PuzzleSize+col]
			cell := vio.R(b.Bounds.X+col*4, b.Bounds.Y+row, 3, 1)
			if t == 0 {
				s.FillRect(cell, ' ', normal)
				continue
			}
			s.FillRect(cell, ' ', highlight)
			s.DrawString(cell.X, cell.Y, fmt.Sprintf("%2d", t), highlight)
		}
	}
}

// solvedTiles returns the canonical 1..15,blank board.
func solvedTiles() [PuzzleSize * PuzzleSize]int {
	var t [PuzzleSize * PuzzleSize]int
	for i := range PuzzleSize*PuzzleSize - 1 {
		t[i] = i + 1
	}
	return t
}

// almostSolvedTiles returns a board where 15 and the blank are
// swapped: one slide-up away from solved.
func almostSolvedTiles() [PuzzleSize * PuzzleSize]int {
	t := solvedTiles()
	t[len(t)-2], t[len(t)-1] = t[len(t)-1], t[len(t)-2]
	return t
}

// randomShuffle returns a board produced by N legal random moves
// from the solved state, guaranteeing solvability.
func randomShuffle() [PuzzleSize * PuzzleSize]int {
	t := solvedTiles()
	b := &puzzleBoard{tiles: t}
	b.findBlank()
	deltas := [4][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}
	for range 200 {
		d := deltas[rand.IntN(4)] //nolint:gosec // shuffling a toy board, not a security boundary
		b.slide(d[0], d[1])
	}
	return b.tiles
}
