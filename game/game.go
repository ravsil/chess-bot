package game

import "fmt"

type Game struct {
	Board           Board
	WhitePieces     []Piece
	BlackPieces     []Piece
	WhiteKing       *King
	BlackKing       *King
	ValidWhiteMoves map[uint64]uint64
	ValidBlackMoves map[uint64]uint64
	CurrentTurn     PColor
}

func (g *Game) NewGame() {
	g.Board.InitBoard()

	for _, pos := range GetBits(g.Board.WhitePawns) {
		g.WhitePieces = append(g.WhitePieces, &Pawn{Color: WHITE, Pos: (1 << pos)})
	}
	for _, pos := range GetBits(g.Board.BlackPawns) {
		g.BlackPieces = append(g.BlackPieces, &Pawn{Color: BLACK, Pos: (1 << pos)})
	}
	for _, pos := range GetBits(g.Board.WhiteRooks) {
		g.WhitePieces = append(g.WhitePieces, &Rook{Color: WHITE, Pos: (1 << pos)})
	}
	for _, pos := range GetBits(g.Board.BlackRooks) {
		g.BlackPieces = append(g.BlackPieces, &Rook{Color: BLACK, Pos: (1 << pos)})
	}
	for _, pos := range GetBits(g.Board.WhiteKnights) {
		g.WhitePieces = append(g.WhitePieces, &Knight{Color: WHITE, Pos: (1 << pos)})
	}
	for _, pos := range GetBits(g.Board.BlackKnights) {
		g.BlackPieces = append(g.BlackPieces, &Knight{Color: BLACK, Pos: (1 << pos)})
	}
	for _, pos := range GetBits(g.Board.WhiteBishops) {
		g.WhitePieces = append(g.WhitePieces, &Bishop{Color: WHITE, Pos: (1 << pos)})
	}
	for _, pos := range GetBits(g.Board.BlackBishops) {
		g.BlackPieces = append(g.BlackPieces, &Bishop{Color: BLACK, Pos: (1 << pos)})
	}
	g.WhitePieces = append(g.WhitePieces, &Queen{Color: WHITE, Pos: g.Board.WhiteQueens})
	g.BlackPieces = append(g.BlackPieces, &Queen{Color: BLACK, Pos: g.Board.BlackQueens})
	whiteKing := &King{Color: WHITE, Pos: g.Board.WhiteKing}
	blackKing := &King{Color: BLACK, Pos: g.Board.BlackKing}
	g.WhitePieces = append(g.WhitePieces, whiteKing)
	g.BlackPieces = append(g.BlackPieces, blackKing)
	g.WhiteKing = whiteKing
	g.BlackKing = blackKing
	g.SaveValidMoves()
}

func (g *Game) SaveValidMoves() {
	g.ValidWhiteMoves = g.GetValidMoves(WHITE)
	g.ValidBlackMoves = g.GetValidMoves(BLACK)
}

func (g *Game) GetValidMoves(color PColor) map[uint64]uint64 {
	moves := make(map[uint64]uint64)
	if color == WHITE {
		for _, piece := range g.WhitePieces {
			moves[piece.GetPos()] = piece.GetValidMoves(g.Board)
		}
		g.ValidWhiteMoves = moves
	} else {
		for _, piece := range g.BlackPieces {
			moves[piece.GetPos()] = piece.GetValidMoves(g.Board)
		}
		g.ValidBlackMoves = moves
	}
	return moves
}
func (g *Game) Move(pos uint64, newPos uint64, color PColor) error {
	if color != g.CurrentTurn {
		return fmt.Errorf("not %v's turn", color)
	}

	if newPos == 0 {
		return fmt.Errorf("null position")
	}

	var it []Piece
	var c PColor
	if color == WHITE {
		it = g.WhitePieces
		c = BLACK
	} else {
		it = g.BlackPieces
		c = WHITE
	}

	for _, p := range it {
		if p.GetPos() == pos {
			p.Move(newPos, &g.Board)
			g.RemovePiece(newPos, c) // go ?
			return nil
		}
	}
	return fmt.Errorf("piece not found at position %d", pos)
}

func (g *Game) RemovePiece(pos uint64, color PColor) error {
	var it *[]Piece
	if color == WHITE {
		it = &g.WhitePieces
	} else {
		it = &g.BlackPieces
	}

	for i, p := range *it {
		if pos == p.GetPos() {
			*it = append((*it)[:i], (*it)[i+1:]...)
			p.Move(uint64(0), &g.Board)
			return nil
		}
	}
	return fmt.Errorf("piece not found at position %d", pos)
}
