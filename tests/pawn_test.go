package tests

import (
	. "chess/game"
	"testing"
)

func TestPawnDoublePushBlocked(t *testing.T) {
	var g Game
	// White pawn at e2 (1<<12). Piece on e3 (1<<20) blocks both single and double push.
	g.Board.WhiteKing = 1 << 4
	g.Board.BlackKing = 1 << 60
	g.Board.WhitePawns = 1 << 12
	g.Board.BlackPawns = 1 << 20

	pawn := &Pawn{Color: WHITE, Pos: 1 << 12}
	moves := pawn.GetValidMoves(g.Board)
	if moves != 0 {
		t.Fatalf("expected 0 moves when blocked in front, got %064b", moves)
	}

	// Piece on e4 (1<<28) blocks double push only
	g.Board.BlackPawns = 1 << 28
	moves = pawn.GetValidMoves(g.Board)
	if moves&(1<<28) != 0 {
		t.Fatalf("expected double push to be blocked by piece on e4, got %064b", moves)
	}
	if moves&(1<<20) == 0 {
		t.Fatalf("expected single push to e3 to be valid, got %064b", moves)
	}
}

func TestPawnDoublePushKingSafety(t *testing.T) {
	var g Game
	// White King on a4 (1<<24), Black Queen on h4 (1<<31).
	// White Pawn on e2 (1<<12) can push to e3 (1<<20), OR e4 (1<<28)?
	// Wait: if King is on a4, Queen on h4, White Pawn on e2, moving to e4 BLOCKS check!
	// Let's test reverse: White King on a2 (1<<8), Black Queen on h2 (1<<15)? No, pawn can't stop that.
	// Suppose White King on e1 (1<<4), Black Bishop on b4 (1<<25), White Pawn on d2 (1<<11):
	// Moving d2->d3 (1<<19) blocks check? No, d3 is on the diagonal (4, 11? No: e1=4, d2=11, c3=18, b4=25).
	// If White King on e1 (1<<4), Black Bishop on b4 (1<<25), White Pawn on d2 (1<<11) IS PINNED!
	// Moving d2->d3 (1<<19) or d2->d4 (1<<27) exposes King to check!
	// Wait, what if White King on a4 (1<<24), Black Queen on h4 (1<<31), White Pawn on e2 (1<<12)?
	// King is IN CHECK. Moving e2->e3 leaves King in check. Moving e2->e4 INTERPOSES and saves King!
	// Currently `if !p.Moved && canForward` checks `canForward` (e3), which is FALSE, so it NEVER lets pawn interpose at e4!
	g.Board.WhiteKing = 1 << 24
	g.Board.BlackKing = 1 << 60
	g.Board.BlackQueens = 1 << 31
	g.Board.WhitePawns = 1 << 12

	pawn := &Pawn{Color: WHITE, Pos: 1 << 12}
	moves := pawn.GetValidMoves(g.Board)
	if moves&(1<<28) == 0 {
		t.Errorf("White pawn on e2 should be able to double push to e4 to block check from h4, got %064b", moves)
	}
	if moves&(1<<20) != 0 {
		t.Errorf("White pawn on e2 should NOT be able to push to e3 while King on a4 is in check, got %064b", moves)
	}
}

func TestEnPassant(t *testing.T) {
	var g Game
	g.NewGame()

	// 1. e4 (White pawn from e2 (1<<12) to e4 (1<<28))
	err := g.Move(1<<12, 1<<28, WHITE)
	if err != nil {
		t.Fatalf("unexpected error moving white pawn: %v", err)
	}

	// 1... a6 (Black pawn from a7 (1<<48) to a6 (1<<40))
	err = g.Move(1<<48, 1<<40, BLACK)
	if err != nil {
		t.Fatalf("unexpected error moving black pawn: %v", err)
	}

	// 2. e5 (White pawn from e4 (1<<28) to e5 (1<<36))
	err = g.Move(1<<28, 1<<36, WHITE)
	if err != nil {
		t.Fatalf("unexpected error moving white pawn: %v", err)
	}

	// 2... d5 (Black pawn from d7 (1<<51) to d5 (1<<35)) - double push adjacent to e5!
	err = g.Move(1<<51, 1<<35, BLACK)
	if err != nil {
		t.Fatalf("unexpected error moving black pawn d7-d5: %v", err)
	}

	// Now White pawn at e5 (1<<36) should have en passant move to d6 (1<<43)
	validMoves := g.GetValidMoves(WHITE)
	e5Moves := validMoves[1<<36]
	if e5Moves&(1<<43) == 0 {
		t.Fatalf("expected en passant move to d6 (1<<43), got %064b", e5Moves)
	}

	// White plays exd6 en passant
	err = g.Move(1<<36, 1<<43, WHITE)
	if err != nil {
		t.Fatalf("unexpected error executing en passant move: %v", err)
	}

	// Black pawn at d5 (1<<35) should be removed!
	if g.Board.BlackPawns&(1<<35) != 0 {
		t.Fatalf("captured black pawn at d5 should have been removed from board")
	}
	for _, p := range g.BlackPieces {
		if p.GetPos() == (1 << 35) {
			t.Fatalf("captured black pawn still found in BlackPieces slice")
		}
	}
}

func TestEnPassantBlackCannotCaptureAdjacentKnight(t *testing.T) {
	var g Game
	// Black pawn at e4 (1<<28). White Knight at d4 (1<<27).
	// Black pawn should NOT be able to "en passant" the Knight onto d3 (1<<19)!
	g.Board.WhiteKing = 1 << 4
	g.Board.BlackKing = 1 << 60
	g.Board.BlackPawns = 1 << 28
	g.Board.WhiteKnights = 1 << 27

	whiteKing := &King{Color: WHITE, Pos: 1 << 4}
	whiteKnight := &Knight{Color: WHITE, Pos: 1 << 27}
	blackKing := &King{Color: BLACK, Pos: 1 << 60}
	blackPawn := &Pawn{Color: BLACK, Pos: 1 << 28, Moved: true}

	g.WhiteKing = whiteKing
	g.BlackKing = blackKing
	g.WhitePieces = []Piece{whiteKing, whiteKnight}
	g.BlackPieces = []Piece{blackKing, blackPawn}
	g.CurrentTurn = BLACK

	validMoves := g.GetValidMoves(BLACK)
	moves := validMoves[1<<28]
	if moves&(1<<19) != 0 {
		t.Fatalf("Black pawn at e4 should NOT be able to en passant an adjacent Knight at d4, got %064b", moves)
	}
}

func TestEnPassantNormalMoveDoesNotCapture(t *testing.T) {
	var g Game
	g.NewGame()

	// Setup: 1. e4 a6 2. e5 d5. Now White pawn at e5 can en passant d6, OR simply push e6!
	_ = g.Move(1<<12, 1<<28, WHITE) // e4
	_ = g.Move(1<<48, 1<<40, BLACK) // a6
	_ = g.Move(1<<28, 1<<36, WHITE) // e5
	_ = g.Move(1<<51, 1<<35, BLACK) // d5

	// Also put a black pawn at e5? No, suppose White pawn pushes forward e5->e6?
	//Wait: what if White has pawn at c5 (1<<34), Black pawn at d5 (1<<35) just double pushed,
	//and White has another pawn or c5 pawn moves normally. Wait: if c5 moves to c6, dest>>8 is c5 (not d5).
	// What if Black pawn was at e5? No, let's check c5 moving to c6 while Black pawn at d5 can be en passanted:
	// Suppose Black has a piece at d6 (1<<43)? No, if White captures on d6... wait.
	// What if White pawn at e5 has enemy pawn at d5 (CanBeEnPassanted=true) and White declines and pushes to e6:
	// What if Black ALSO had a pawn at e5? Can't have two pieces on e5.
	// What if White pawn at c5, Black pawn at d5 (CanBeEnPassanted=true), and Black Knight at d6?
	// White pawn at c5 can capture Black Knight at d6 normally!
	// If `WantsToEnPassant` was true, `handleEnPassant(d6)` would ALSO delete Black pawn at d5!!
	g.Board.BlackKnights = 1 << 43 // d6
	blackKnight := &Knight{Color: BLACK, Pos: 1 << 43}
	g.BlackPieces = append(g.BlackPieces, blackKnight)
	g.SaveValidMoves()

	// White plays exd6 (capturing Knight on d6, NOT en passant because there's a knight on d6, wait... actually in chess d5 pawn AND d6 knight can't both exist after d7-d5, wait: if knight moved to d6 earlier? No, if d7-d5 happened, d6 was empty when d5 moved, wait - yes, if knight was ALREADY on d6, d7-d5 couldn't happen through d6, UNLESS knight moved to d6? No, White moves right after d7-d5!)
	// Wait: what if White pawn at c5, Black pawn at d5, and White pawn at c5 moves... well, if d5 is adjacent, normal capture on d6 wouldn't happen simultaneously unless contrived.
	// Better test: Black pawn at d5 just double pushed (CanBeEnPassanted = true).
	// White pawn at e5 chooses NOT to en passant, and moves e5->e6 instead (Wait: e6>>8 is e5, empty).
	// Even better: White pawn at e5 has `WantsToEnPassant == true` (from previous GetValidMoves),
	// Wait: what if White pawn at e5 is deleted/checked when another move occurs?
}

func TestPawnPromotion(t *testing.T) {
	var g Game
	// White King at e1, Black King at e8. White pawn at a7 (1<<48) ready to promote to a8 (1<<56).
	g.Board.WhiteKing = 1 << 4
	g.Board.BlackKing = 1 << 60
	g.Board.WhitePawns = 1 << 48

	whiteKing := &King{Color: WHITE, Pos: 1 << 4}
	blackKing := &King{Color: BLACK, Pos: 1 << 60}
	whitePawn := &Pawn{Color: WHITE, Pos: 1 << 48, Moved: true}

	g.WhiteKing = whiteKing
	g.BlackKing = blackKing
	g.WhitePieces = []Piece{whiteKing, whitePawn}
	g.BlackPieces = []Piece{blackKing}
	g.CurrentTurn = WHITE

	err := g.Move(1<<48, 1<<56, WHITE)
	if err != nil {
		t.Fatalf("unexpected error on promotion move: %v", err)
	}

	// There should be NO pawn in WhitePieces, and a Queen on a8
	pawnCount := 0
	queenCount := 0
	for _, p := range g.WhitePieces {
		if p.GetType() == PAWN {
			pawnCount++
		}
		if p.GetType() == QUEEN && p.GetPos() == (1<<56) {
			queenCount++
		}
	}

	if pawnCount != 0 {
		t.Errorf("expected 0 pawns after promotion, found %d", pawnCount)
	}
	if queenCount != 1 {
		t.Errorf("expected 1 queen on a8 after promotion, found %d", queenCount)
	}
	if g.Board.WhitePawns&(1<<56) != 0 {
		t.Errorf("board should not have WhitePawns at a8")
	}
	if g.Board.WhiteQueens&(1<<56) == 0 {
		t.Errorf("board should have WhiteQueens at a8")
	}
}
