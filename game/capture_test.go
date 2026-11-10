package game

import (
	"testing"
)

func TestCaptureCheckingPiece(t *testing.T) {
	InitRanged()
	var g Game
	// Set up board where White King is on e1 (1<<4), Black Bishop on b4 (1<<25) checking King.
	// White Pawn on c3 (1<<18) can capture Bishop on b4.
	g.Board.WhiteKing = 1 << 4
	g.Board.BlackKing = 1 << 60
	g.Board.BlackBishops = 1 << 25
	g.Board.WhitePawns = 1 << 18

	whiteKing := &King{Color: WHITE, Pos: 1 << 4}
	whitePawn := &Pawn{Color: WHITE, Pos: 1 << 18, Moved: true}
	blackKing := &King{Color: BLACK, Pos: 1 << 60}
	blackBishop := &Bishop{Color: BLACK, Pos: 1 << 25}

	g.WhiteKing = whiteKing
	g.BlackKing = blackKing
	g.WhitePieces = []Piece{whiteKing, whitePawn}
	g.BlackPieces = []Piece{blackKing, blackBishop}
	g.CurrentTurn = WHITE

	// Pawn at c3 should have valid move to b4 (capturing the bishop)
	moves := whitePawn.GetValidMoves(g.Board)
	target := uint64(1 << 25)
	if moves&target == 0 {
		t.Fatalf("Pawn at c3 should be able to capture checking Bishop at b4, but moves bitmask is %064b", moves)
	}
}

func TestKingCaptureCheckingPiece(t *testing.T) {
	InitRanged()
	var g Game
	// White King at e1 (1<<4), unprotected Black Queen at e2 (1<<12) checking King.
	g.Board.WhiteKing = 1 << 4
	g.Board.BlackKing = 1 << 60
	g.Board.BlackQueens = 1 << 12

	whiteKing := &King{Color: WHITE, Pos: 1 << 4}
	blackKing := &King{Color: BLACK, Pos: 1 << 60}
	blackQueen := &Queen{Color: BLACK, Pos: 1 << 12}

	g.WhiteKing = whiteKing
	g.BlackKing = blackKing
	g.WhitePieces = []Piece{whiteKing}
	g.BlackPieces = []Piece{blackKing, blackQueen}
	g.CurrentTurn = WHITE

	moves := whiteKing.GetValidMoves(g.Board)
	target := uint64(1 << 12)
	if moves&target == 0 {
		t.Fatalf("King at e1 should be able to capture unprotected checking Queen at e2, but got %064b", moves)
	}
}

func TestKingCannotCaptureProtectedCheckingPiece(t *testing.T) {
	InitRanged()
	var g Game
	// White King at e1 (1<<4), Black Queen at e2 (1<<12) checking King, protected by Black Bishop at c4 (1<<26).
	g.Board.WhiteKing = 1 << 4
	g.Board.BlackKing = 1 << 60
	g.Board.BlackQueens = 1 << 12
	g.Board.BlackBishops = 1 << 26

	whiteKing := &King{Color: WHITE, Pos: 1 << 4}
	blackKing := &King{Color: BLACK, Pos: 1 << 60}
	blackQueen := &Queen{Color: BLACK, Pos: 1 << 12}
	blackBishop := &Bishop{Color: BLACK, Pos: 1 << 26}

	g.WhiteKing = whiteKing
	g.BlackKing = blackKing
	g.WhitePieces = []Piece{whiteKing}
	g.BlackPieces = []Piece{blackKing, blackQueen, blackBishop}
	g.CurrentTurn = WHITE

	moves := whiteKing.GetValidMoves(g.Board)
	target := uint64(1 << 12)
	if moves&target != 0 {
		t.Fatalf("King at e1 should NOT be able to capture protected Queen at e2, but got move %064b", moves)
	}
}
