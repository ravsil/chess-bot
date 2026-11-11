package tests

import (
	. "chess/game"
	"testing"
)

func TestCheckmateDetection(t *testing.T) {
	var g Game
	// Fool's mate position:
	// White: f2->f3, g2->g4; Black: e7->e5, Qd8->h4#
	g.NewGame()
	_ = g.Move(1<<13, 1<<21, WHITE) // f3
	_ = g.Move(1<<52, 1<<36, BLACK) // e5
	_ = g.Move(1<<14, 1<<30, WHITE) // g4
	_ = g.Move(1<<59, 1<<31, BLACK) // Qh4#

	moves := g.GetValidMoves(WHITE)
	if len(moves) != 0 {
		t.Errorf("expected 0 valid moves for checkmated White, got %d entries: %v", len(moves), moves)
	}
	if !g.IsInCheck(WHITE) {
		t.Errorf("expected White to be in check")
	}
	if !g.IsCheckmate(WHITE) {
		t.Errorf("expected White to be in checkmate")
	}
	if g.IsStalemate(WHITE) {
		t.Errorf("expected White NOT to be in stalemate")
	}
}

func TestStalemateDetection(t *testing.T) {
	var g Game
	// Classic stalemate: White King on a8 (1<<56), Black King on c7 (1<<50), Black Queen on b6 (1<<41).
	// Wait: a8=56, b8=57, a7=48, b7=49. Queen on c7 (50) and King on c6 (42)?
	// Queen on b6 (1<<41) attacks a7 (48), b7 (49), b8 (57), a6, etc. Wait: b6 to a7 is -1+8=+7 (diagonal). b6 to b7 is +8. b6 to b8 is +16.
	// What attacks b8 (57)? Queen on b6 attacks b8! What attacks a7 (48)? Queen on b6 attacks a7!
	// What attacks a8 (56)? Nothing! So King on a8 is NOT in check, and a7, b7, b8 are all attacked!
	// Black King on c8 (1<<58) or h8 (1<<63) protects Queen? King on c6 (1<<42) protects b6 (1<<41)!
	// Check: White King on a8 (56). Can it move to:
	// - b8 (57): attacked by Qb6 (41 + 8 + 8 = 57)
	// - b7 (49): attacked by Qb6 (41 + 8 = 49)
	// - a7 (48): attacked by Qb6 (41 + 7 = 48)
	// Is a8 attacked? Qb6 is at rank 5 col 1. a8 is rank 7 col 0. Not on file, rank, or diagonal!
	// And White has no other pieces. Black King at c6 (42) protects Qb6 (41) and doesn't interfere.
	g.Board.WhiteKing = 1 << 56
	g.Board.BlackQueens = 1 << 41
	g.Board.BlackKing = 1 << 42

	wk := &King{Color: WHITE, Pos: 1 << 56, Moved: true}
	bq := &Queen{Color: BLACK, Pos: 1 << 41}
	bk := &King{Color: BLACK, Pos: 1 << 42, Moved: true}

	g.WhiteKing = wk
	g.BlackKing = bk
	g.WhitePieces = []Piece{wk}
	g.BlackPieces = []Piece{bk, bq}
	g.CurrentTurn = WHITE

	moves := g.GetValidMoves(WHITE)
	if len(moves) != 0 {
		t.Errorf("expected 0 valid moves in stalemate, got %d entries: %v", len(moves), moves)
	}
	if g.IsInCheck(WHITE) {
		t.Errorf("expected White NOT to be in check during stalemate")
	}
	if g.IsCheckmate(WHITE) {
		t.Errorf("expected White NOT to be in checkmate during stalemate")
	}
	if !g.IsStalemate(WHITE) {
		t.Errorf("expected White to be in stalemate")
	}
}
