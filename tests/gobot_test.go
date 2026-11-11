package tests

import (
	. "chess/bot"
	. "chess/game"
	"testing"
)

func TestBotFindsMateInOne(t *testing.T) {
	var g Game
	// Back rank mate in 1 for Black:
	// White King at g1 (1<<6), pawns at f2, g2, h2 (bits 13, 14, 15).
	// Black King at g8 (1<<62), Black Rook at a8 (1<<56).
	// Black Rook a8 -> a1 (1<<0) is checkmate!
	g.Board.WhiteKing = 1 << 6
	g.Board.WhitePawns = (1 << 13) | (1 << 14) | (1 << 15)
	g.Board.BlackKing = 1 << 62
	g.Board.BlackRooks = 1 << 56

	wk := &King{Color: WHITE, Pos: 1 << 6, Moved: true}
	wp1 := &Pawn{Color: WHITE, Pos: 1 << 13}
	wp2 := &Pawn{Color: WHITE, Pos: 1 << 14}
	wp3 := &Pawn{Color: WHITE, Pos: 1 << 15}
	bk := &King{Color: BLACK, Pos: 1 << 62, Moved: true}
	br := &Rook{Color: BLACK, Pos: 1 << 56, Moved: true}

	g.WhiteKing = wk
	g.BlackKing = bk
	g.WhitePieces = []Piece{wk, wp1, wp2, wp3}
	g.BlackPieces = []Piece{bk, br}
	g.CurrentTurn = BLACK
	g.SaveValidMoves()

	bot := GoBot{Color: BLACK}
	if err := bot.Play(&g, BLACK); err != nil {
		t.Fatalf("bot failed to play: %v", err)
	}

	if !g.IsCheckmate(WHITE) {
		t.Fatalf("expected bot to find back-rank checkmate on a1, board WhiteKing=%064b BlackRooks=%064b", g.Board.WhiteKing, g.Board.BlackRooks)
	}
}

func TestBotAvoidsStalemateWhenWinning(t *testing.T) {
	var g Game
	// White King at a8 (1<<56), Black King at c6 (1<<42), Black Queen at d7 (1<<51).
	// If Black moves Qd7->b6 (1<<41), it's stalemate (draw, 0)!
	// Black can move Qd7->c7 (1<<50) or Qb7? Wait: Qd7->b7 (1<<49) is checkmate! (protected by Kc6? No, c6 attacks b7! Yes! c6=42, b7=49, diff=+7 diagonal! Qb7# is checkmate!)
	// What if Black Queen is at d5 (1<<35)?
	// From d5, Queen can go to b7? No, d5->b7 is diagonal (35+7+7=49)! Wait: Kc6 blocks? c6=42 is between d5(35) and b7(49)!
	// Let's set Black Queen at h6 (1<<47).
	// From h6 (47): h6->b6 (41) is stalemate! h6->a6? (check? No, a8->a6 is check? Wait: a8 is on file a, a6 is on file a, a7 is empty -> Qa6 is check, Ka8 takes? No, Kc6 doesn't protect a6).
	// Better: Black Queen at g6 (1<<46), Black King at c6 (1<<42).
	// Black can play Qg6->b6 (stalemate!) or Qg6->c7 (1<<50)? No, g6->g8 is check!
	// Even simpler: White King at h1 (1<<7), Black King at f2 (1<<13), Black Queen at g4 (1<<30).
	// If Black plays Qg4->g2 (1<<14) -> checkmate!
	// If Black plays Qg4->g3 (1<<22) -> White King at h1 has NO moves (g1 attacked by Qg3, g2 attacked by Kf2 & Qg3, h2 attacked by Qg3 & Kf2? No, f2 doesn't attack h2, Qg3 attacks h2! So Qg3 is stalemate!).
	// Black should NOT play Qg3 (stalemate) and SHOULD play Qg2 (mate) or another winning move!
	g.Board.WhiteKing = 1 << 7
	g.Board.BlackKing = 1 << 13
	g.Board.BlackQueens = 1 << 30

	wk := &King{Color: WHITE, Pos: 1 << 7, Moved: true}
	bk := &King{Color: BLACK, Pos: 1 << 13, Moved: true}
	bq := &Queen{Color: BLACK, Pos: 1 << 30}

	g.WhiteKing = wk
	g.BlackKing = bk
	g.WhitePieces = []Piece{wk}
	g.BlackPieces = []Piece{bk, bq}
	g.CurrentTurn = BLACK
	g.SaveValidMoves()

	bot := GoBot{Color: BLACK}
	if err := bot.Play(&g, BLACK); err != nil {
		t.Fatalf("bot failed to play: %v", err)
	}

	if g.IsStalemate(WHITE) {
		t.Fatalf("bot walked into stalemate instead of winning!")
	}
	if !g.IsCheckmate(WHITE) {
		t.Fatalf("expected bot to find Qg2# checkmate, BlackQueens=%064b", g.Board.BlackQueens)
	}
}

func TestEvaluateStalemateIsZero(t *testing.T) {
	var g Game
	// White King on a8 (1<<56), Black Queen on b6 (1<<41), Black King on c6 (1<<42).
	// White to move, White has 0 moves and is NOT in check => Stalemate.
	g.Board.WhiteKing = 1 << 56
	g.Board.BlackQueens = 1 << 41
	g.Board.BlackKing = 1 << 42
	g.CurrentTurn = WHITE
	g.ValidWhiteMoves = make(map[uint64]uint64) // 0 legal moves

	score := Evaluate(&g, 0)
	if score != 0.0 {
		t.Fatalf("expected stalemate evaluation to be 0.0, got %f", score)
	}
}

func TestPositionKeyIncludesAllPieces(t *testing.T) {
	var g1, g2 Game
	g1.NewGame()
	g2.NewGame()

	// Remove a bishop from g2
	g2.Board.WhiteBishops = 0
	if GetPositionKey(&g1) == GetPositionKey(&g2) {
		t.Errorf("position keys collided despite different bishops")
	}

	// Remove a queen from g2
	g2.NewGame()
	g2.Board.BlackQueens = 0
	if GetPositionKey(&g1) == GetPositionKey(&g2) {
		t.Errorf("position keys collided despite different queens")
	}

	// Swap kings
	g2.NewGame()
	g2.Board.WhiteKing, g2.Board.BlackKing = g2.Board.BlackKing, g2.Board.WhiteKing
	if GetPositionKey(&g1) == GetPositionKey(&g2) {
		t.Errorf("position keys collided despite swapped kings")
	}
}
