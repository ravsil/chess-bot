package game

import "fmt"

type Game struct {
	Board       Board
	WhitePieces []Piece
	BlackPieces []Piece
	WhiteKing   *King
	BlackKing   *King
	WhiteRooks  struct {
		Left  *Rook
		Right *Rook
	}
	BlackRooks struct {
		Left  *Rook
		Right *Rook
	}
	ValidWhiteMoves map[uint64]uint64
	ValidBlackMoves map[uint64]uint64
	CurrentTurn     PColor
}

func (g *Game) NewGame() {
	g.Board.InitBoard()
	g.CurrentTurn = WHITE
	g.WhitePieces = nil
	g.BlackPieces = nil
	g.WhiteKing = nil
	g.BlackKing = nil
	g.WhiteRooks.Left = nil
	g.WhiteRooks.Right = nil
	g.BlackRooks.Left = nil
	g.BlackRooks.Right = nil
	g.ValidWhiteMoves = nil
	g.ValidBlackMoves = nil

	for _, pos := range GetBits(g.Board.WhitePawns) {
		g.WhitePieces = append(g.WhitePieces, &Pawn{Color: WHITE, Pos: (1 << pos)})
	}
	for _, pos := range GetBits(g.Board.BlackPawns) {
		g.BlackPieces = append(g.BlackPieces, &Pawn{Color: BLACK, Pos: (1 << pos)})
	}
	for i, pos := range GetBits(g.Board.WhiteRooks) {
		rook := &Rook{Color: WHITE, Pos: (1 << pos)}
		g.WhitePieces = append(g.WhitePieces, rook)
		switch i {
		case 0:
			g.WhiteRooks.Right = rook
		case 1:
			g.WhiteRooks.Left = rook
		}
	}
	for i, pos := range GetBits(g.Board.BlackRooks) {
		rook := &Rook{Color: BLACK, Pos: (1 << pos)}
		g.BlackPieces = append(g.BlackPieces, rook)
		switch i {
		case 0:
			g.BlackRooks.Right = rook
		case 1:
			g.BlackRooks.Left = rook
		}
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
		if !g.WhiteKing.Moved {
			l := uint64(0b00000000_00000000_00000000_00000000_00000000_00000000_00000000_01000000)
			// removes white's short castling if it should'nt be available
			if moves[g.WhiteKing.Pos]&l != 0 && (g.WhiteRooks.Left == nil || g.WhiteRooks.Left.Moved) {
				moves[g.WhiteKing.Pos] &^= l
			}
			r := uint64(0b00000000_00000000_00000000_00000000_00000000_00000000_00000000_00000100)
			// removes white's long castling if it should'nt be available
			if moves[g.WhiteKing.Pos]&r != 0 && (g.WhiteRooks.Right == nil || g.WhiteRooks.Right.Moved) {
				moves[g.WhiteKing.Pos] &^= r
			}
		}
	} else {
		for _, piece := range g.BlackPieces {
			moves[piece.GetPos()] = piece.GetValidMoves(g.Board)
		}
		if !g.BlackKing.Moved {
			l := uint64(0b01000000_00000000_00000000_00000000_00000000_00000000_00000000_00000000)
			// removes black's short castling if it should'nt be available
			if moves[g.BlackKing.Pos]&l != 0 && (g.BlackRooks.Left == nil || g.BlackRooks.Left.Moved) {
				moves[g.BlackKing.Pos] &^= l
			}
			r := uint64(0b00000100_00000000_00000000_00000000_00000000_00000000_00000000_00000000)
			// removes black's long castling if it should'nt be available
			if moves[g.BlackKing.Pos]&r != 0 && (g.BlackRooks.Right == nil || g.BlackRooks.Right.Moved) {
				moves[g.BlackKing.Pos] &^= r
			}
		}
	}
	return moves
}

func (g *Game) handleCastling(newPos uint64, color PColor) {
	if color == WHITE {
		// moves the white rook to make short castles
		if newPos == uint64(0b00000000_00000000_00000000_00000000_00000000_00000000_00000000_01000000) {
			if g.WhiteRooks.Left != nil && !g.WhiteRooks.Left.Moved {
				g.WhiteRooks.Left.Move(0b00000000_00000000_00000000_00000000_00000000_00000000_00000000_00100000, &g.Board)
			}
		}
		// moves the white rook to make long castles
		if newPos == uint64(0b00000000_00000000_00000000_00000000_00000000_00000000_00000000_00000100) {
			if g.WhiteRooks.Right != nil && !g.WhiteRooks.Right.Moved {
				g.WhiteRooks.Right.Move(0b00000000_00000000_00000000_00000000_00000000_00000000_00000000_00001000, &g.Board)
			}
		}
	} else {
		// moves the black rook to make short castles
		if newPos == uint64(0b01000000_00000000_00000000_00000000_00000000_00000000_00000000_00000000) {
			if g.BlackRooks.Left != nil && !g.BlackRooks.Left.Moved {
				g.BlackRooks.Left.Move(0b00100000_00000000_00000000_00000000_00000000_00000000_00000000_00000000, &g.Board)
			}
		}
		// moves the black rook to make long castles
		if newPos == uint64(0b00000100_00000000_00000000_00000000_00000000_00000000_00000000_00000000) {
			if g.BlackRooks.Right != nil && !g.BlackRooks.Right.Moved {
				g.BlackRooks.Right.Move(0b00001000_00000000_00000000_00000000_00000000_00000000_00000000_00000000, &g.Board)
			}
		}
	}
}

func (g *Game) handleEnPassant(newPos uint64, color PColor) {
	if color == WHITE {
		dest := newPos >> 8
		for i, p := range g.BlackPieces {
			if p.GetPos() == dest && p.GetType() == PAWN && p.(*Pawn).CanBeEnPassanted {
				g.BlackPieces = append(g.BlackPieces[:i], g.BlackPieces[i+1:]...)
				p.Move(uint64(0), &g.Board)
			}
		}
	} else {
		dest := newPos << 8
		for i, p := range g.WhitePieces {
			if p.GetPos() == dest && p.GetType() == PAWN && p.(*Pawn).CanBeEnPassanted {
				g.WhitePieces = append(g.WhitePieces[:i], g.WhitePieces[i+1:]...)
				p.Move(uint64(0), &g.Board)
			}
		}
	}
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
			if p.GetType() == KING && !p.(*King).Moved {
				g.handleCastling(newPos, color)
			} else if p.GetType() == PAWN && p.(*Pawn).WantsToEnPassant {
				g.handleEnPassant(newPos, color)
			}
			p.Move(newPos, &g.Board)
			g.ChangeTurn()
			g.RemovePiece(newPos, c) // go ?
			g.SaveValidMoves()
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

func (g *Game) ChangeTurn() {
	var it *[]Piece

	if g.CurrentTurn == WHITE {
		it = &g.BlackPieces
		g.CurrentTurn = BLACK
	} else {
		it = &g.WhitePieces
		g.CurrentTurn = WHITE
	}

	for _, p := range *it {
		if p.GetType() == PAWN {
			p.(*Pawn).CanBeEnPassanted = false
		}
	}
}

func (g *Game) GetTurn() PColor {
	return g.CurrentTurn
}
