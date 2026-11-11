package tests

import (
	. "chess/game"
	"testing"
)

func TestCannotCastleWhileInCheck(t *testing.T) {
	var g Game
	// White King on e1 (1<<4), White Rook on h1 (1<<7), Black Rook on e8 (1<<60) giving check.
	g.Board.WhiteKing = 1 << 4
	g.Board.WhiteRooks = 1 << 7
	g.Board.BlackKing = 1 << 56
	g.Board.BlackRooks = 1 << 60

	wk := &King{Color: WHITE, Pos: 1 << 4}
	wr := &Rook{Color: WHITE, Pos: 1 << 7}
	bk := &King{Color: BLACK, Pos: 1 << 56}
	br := &Rook{Color: BLACK, Pos: 1 << 60}

	g.WhiteKing = wk
	g.WhiteRooks.Left = wr
	g.BlackKing = bk
	g.WhitePieces = []Piece{wk, wr}
	g.BlackPieces = []Piece{bk, br}
	g.CurrentTurn = WHITE

	moves := g.GetValidMoves(WHITE)
	if moves[wk.Pos]&(1<<6) != 0 {
		t.Fatalf("White King should NOT be able to castle kingside while in check, got %064b", moves[wk.Pos])
	}
}

func TestBlackQueensideCastlingBlockedOnB8(t *testing.T) {
	var g Game
	// Black King on e8 (1<<60), Black Rook on a8 (1<<56), Black Knight on b8 (1<<57).
	g.Board.WhiteKing = 1 << 4
	g.Board.BlackKing = 1 << 60
	g.Board.BlackRooks = 1 << 56
	g.Board.BlackKnights = 1 << 57

	wk := &King{Color: WHITE, Pos: 1 << 4}
	bk := &King{Color: BLACK, Pos: 1 << 60}
	br := &Rook{Color: BLACK, Pos: 1 << 56}
	bn := &Knight{Color: BLACK, Pos: 1 << 57}

	g.WhiteKing = wk
	g.BlackKing = bk
	g.BlackRooks.Right = br
	g.WhitePieces = []Piece{wk}
	g.BlackPieces = []Piece{bk, br, bn}
	g.CurrentTurn = BLACK

	moves := g.GetValidMoves(BLACK)
	if moves[bk.Pos]&(1<<58) != 0 {
		t.Fatalf("Black King should NOT be able to castle queenside when b8 is occupied, got %064b", moves[bk.Pos])
	}
}

func TestWhiteQueensideCastlingAllowedWhenB1Attacked(t *testing.T) {
	var g Game
	// White King on e1 (1<<4), White Rook on a1 (1<<0).
	// Black Rook on b8 (1<<57) attacks b1 (1<<1). b1 is empty, d1 & c1 & e1 are safe.
	// By FIDE rules, White CAN castle queenside!
	g.Board.WhiteKing = 1 << 4
	g.Board.WhiteRooks = 1 << 0
	g.Board.BlackKing = 1 << 60
	g.Board.BlackRooks = 1 << 57

	wk := &King{Color: WHITE, Pos: 1 << 4}
	wr := &Rook{Color: WHITE, Pos: 1 << 0}
	bk := &King{Color: BLACK, Pos: 1 << 60}
	br := &Rook{Color: BLACK, Pos: 1 << 57}

	g.WhiteKing = wk
	g.WhiteRooks.Right = wr
	g.BlackKing = bk
	g.WhitePieces = []Piece{wk, wr}
	g.BlackPieces = []Piece{bk, br}
	g.CurrentTurn = WHITE

	moves := g.GetValidMoves(WHITE)
	if moves[wk.Pos]&(1<<2) == 0 {
		t.Fatalf("White King SHOULD be able to castle queenside when only b1 is attacked, got %064b", moves[wk.Pos])
	}
}

func TestCapturedRookCannotCastle(t *testing.T) {
	var g Game
	// White King on e1 (1<<4), White Rook on h1 (1<<7).
	// Black Bishop on e4 (1<<28) captures White Rook on h1 (1<<7)!
	g.Board.WhiteKing = 1 << 4
	g.Board.WhiteRooks = 1 << 7
	g.Board.BlackKing = 1 << 60
	g.Board.BlackBishops = 1 << 28

	wk := &King{Color: WHITE, Pos: 1 << 4}
	wr := &Rook{Color: WHITE, Pos: 1 << 7}
	bk := &King{Color: BLACK, Pos: 1 << 60}
	bb := &Bishop{Color: BLACK, Pos: 1 << 28}

	g.WhiteKing = wk
	g.WhiteRooks.Left = wr
	g.BlackKing = bk
	g.WhitePieces = []Piece{wk, wr}
	g.BlackPieces = []Piece{bk, bb}
	g.CurrentTurn = BLACK

	// Black captures h1 rook
	err := g.Move(1<<28, 1<<7, BLACK)
	if err != nil {
		t.Fatalf("unexpected error capturing rook: %v", err)
	}
	// Black bishop moves back off h1 so g1/f1 are empty
	// White moves pawn? No, let's check White King moves right after Black captures or after Bishop moves away
	// Even while Bishop is on h1, White shouldn't think h1 rook can castle
	// Let's move White King? Wait: let's move White King or check valid moves after bishop moves away:
	// First check White's valid moves right now: h1 is occupied by Black Bishop, f1/g1 empty.
	moves := g.GetValidMoves(WHITE)
	if moves[wk.Pos]&(1<<6) != 0 {
		t.Fatalf("White King should NOT be able to castle kingside after h1 rook was captured, got %064b", moves[wk.Pos])
	}
	if g.WhiteRooks.Left != nil {
		t.Fatalf("expected WhiteRooks.Left to be nil after capture")
	}
}

func TestCastlingExecution(t *testing.T) {
	var g Game
	g.NewGame()
	// Clear f1, g1, b1, c1, d1 on White side
	g.Board.WhiteKnights = 0
	g.Board.WhiteBishops = 0
	g.Board.WhiteQueens = 0
	var kept []Piece
	for _, p := range g.WhitePieces {
		if p.GetType() == KING || p.GetType() == ROOK || p.GetType() == PAWN {
			kept = append(kept, p)
		}
	}
	g.WhitePieces = kept

	// White castles kingside: e1 (1<<4) -> g1 (1<<6)
	if err := g.Move(1<<4, 1<<6, WHITE); err != nil {
		t.Fatalf("failed white kingside castle: %v", err)
	}
	if g.Board.WhiteKing != (1 << 6) {
		t.Errorf("expected White King at g1, got %064b", g.Board.WhiteKing)
	}
	if g.Board.WhiteRooks&(1<<5) == 0 || g.Board.WhiteRooks&(1<<7) != 0 {
		t.Errorf("expected White Rook moved from h1 to f1, got %064b", g.Board.WhiteRooks)
	}
}
