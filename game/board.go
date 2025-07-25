package game

type Board struct {
	WhitePawns   uint64
	BlackPawns   uint64
	WhiteRooks   uint64
	BlackRooks   uint64
	WhiteKnights uint64
	BlackKnights uint64
	WhiteBishops uint64
	BlackBishops uint64
	WhiteQueens  uint64
	BlackQueens  uint64
	WhiteKing    uint64
	BlackKing    uint64
}

func (b *Board) InitBoard() {
	b.WhitePawns = 0b00000000_00000000_00000000_00000000_00000000_00000000_11111111_00000000
	b.WhiteRooks = 0b00000000_00000000_00000000_00000000_00000000_00000000_00000000_10000001
	b.WhiteKnights = 0b00000000_00000000_00000000_00000000_00000000_00000000_00000000_01000010
	b.WhiteBishops = 0b00000000_00000000_00000000_00000000_00000000_00000000_00000000_00100100
	b.WhiteQueens = 0b00000000_00000000_00000000_00000000_00000000_00000000_00000000_00001000
	b.WhiteKing = 0b00000000_00000000_00000000_00000000_00000000_00000000_00000000_00010000

	b.BlackPawns = 0b00000000_11111111_00000000_00000000_00000000_00000000_00000000_00000000
	b.BlackRooks = 0b10000001_00000000_00000000_00000000_00000000_00000000_00000000_00000000
	b.BlackKnights = 0b01000010_00000000_00000000_00000000_00000000_00000000_00000000_00000000
	b.BlackBishops = 0b00100100_00000000_00000000_00000000_00000000_00000000_00000000_00000000
	b.BlackQueens = 0b00001000_00000000_00000000_00000000_00000000_00000000_00000000_00000000
	b.BlackKing = 0b00010000_00000000_00000000_00000000_00000000_00000000_00000000_00000000
}

func (b *Board) GetPieces(c PColor) uint64 {
	if c == WHITE {
		return b.WhitePawns | b.WhiteRooks | b.WhiteKnights | b.WhiteBishops | b.WhiteQueens | b.WhiteKing
	}
	return b.BlackPawns | b.BlackRooks | b.BlackKnights | b.BlackBishops | b.BlackQueens | b.BlackKing
}

func (b *Board) GetOcupiedSquares() uint64 {
	return b.WhitePawns | b.BlackPawns | b.WhiteRooks | b.BlackRooks |
		b.WhiteKnights | b.BlackKnights | b.WhiteBishops | b.BlackBishops |
		b.WhiteQueens | b.BlackQueens | b.WhiteKing | b.BlackKing
}

// returns the squares black pieces are attacking
func (b *Board) GetWhiteInfluence() uint64 {
	influence := uint64(0)
	allPieces := b.GetOcupiedSquares()

	influence |= (b.WhitePawns << 9) & ^LEFT_BORDER  // ↖
	influence |= (b.WhitePawns << 7) & ^RIGHT_BORDER // ↗

	influence |= (b.WhiteKing << 8) | // ↑
		(b.WhiteKing >> 8) | // ↓
		((b.WhiteKing << 1) & ^LEFT_BORDER) | // ←
		((b.WhiteKing >> 1) & ^RIGHT_BORDER) | // →
		((b.WhiteKing << 9) & ^LEFT_BORDER) | // ↖
		((b.WhiteKing << 7) & ^RIGHT_BORDER) | // ↗
		((b.WhiteKing >> 7) & ^LEFT_BORDER) | // ↙
		((b.WhiteKing >> 9) & ^RIGHT_BORDER) // ↘

	influence |= (b.WhiteKnights << 17) & ^LEFT_BORDER         // ↑↑→
	influence |= (b.WhiteKnights << 15) & ^RIGHT_BORDER        // ↑↑←
	influence |= (b.WhiteKnights >> 17) & ^RIGHT_BORDER        // ↓↓←
	influence |= (b.WhiteKnights >> 15) & ^LEFT_BORDER         // ↓↓→
	influence |= (b.WhiteKnights << 10) & ^DOUBLE_LEFT_BORDER  // ↑→→
	influence |= (b.WhiteKnights << 6) & ^DOUBLE_RIGHT_BORDER  // ↑←←
	influence |= (b.WhiteKnights >> 6) & ^DOUBLE_LEFT_BORDER   // ↓→→
	influence |= (b.WhiteKnights >> 10) & ^DOUBLE_RIGHT_BORDER // ↓←←

	for _, bishop := range GetBits(b.WhiteBishops) {
		influence |= BishopAttacks(allPieces, bishop)
	}
	for _, rook := range GetBits(b.WhiteRooks) {
		influence |= RookAttacks(allPieces, rook)
	}
	for _, queen := range GetBits(b.WhiteQueens) {
		influence |= BishopAttacks(allPieces, queen)
		influence |= RookAttacks(allPieces, queen)
	}
	return influence
}

func (b *Board) GetBlackInfluence() uint64 {
	influence := uint64(0)

	allPieces := b.GetOcupiedSquares()

	influence |= (b.BlackPawns >> 9) & ^LEFT_BORDER  // ↖
	influence |= (b.BlackPawns >> 7) & ^RIGHT_BORDER // ↗

	influence |= (b.BlackKing << 8) | // ↑
		(b.BlackKing >> 8) | // ↓
		((b.BlackKing << 1) & ^LEFT_BORDER) | // ←
		((b.BlackKing >> 1) & ^RIGHT_BORDER) | // →
		((b.BlackKing << 9) & ^LEFT_BORDER) | // ↖
		((b.BlackKing << 7) & ^RIGHT_BORDER) | // ↗
		((b.BlackKing >> 7) & ^LEFT_BORDER) | // ↙
		((b.BlackKing >> 9) & ^RIGHT_BORDER) // ↘

	influence |= (b.BlackKnights << 17) & ^LEFT_BORDER         // ↑↑→
	influence |= (b.BlackKnights << 15) & ^RIGHT_BORDER        // ↑↑←
	influence |= (b.BlackKnights >> 17) & ^RIGHT_BORDER        // ↓↓←
	influence |= (b.BlackKnights >> 15) & ^LEFT_BORDER         // ↓↓→
	influence |= (b.BlackKnights << 10) & ^DOUBLE_LEFT_BORDER  // ↑→→
	influence |= (b.BlackKnights << 6) & ^DOUBLE_RIGHT_BORDER  // ↑←←
	influence |= (b.BlackKnights >> 6) & ^DOUBLE_LEFT_BORDER   // ↓→→
	influence |= (b.BlackKnights >> 10) & ^DOUBLE_RIGHT_BORDER // ↓←←

	for _, bishop := range GetBits(b.BlackBishops) {
		influence |= BishopAttacks(allPieces, bishop)
	}
	for _, rook := range GetBits(b.BlackRooks) {
		influence |= RookAttacks(allPieces, rook)
	}
	for _, queen := range GetBits(b.BlackQueens) {
		influence |= BishopAttacks(allPieces, queen)
		influence |= RookAttacks(allPieces, queen)
	}
	return influence
}

func (b *Board) Update(pos uint64, newPos uint64, piece PType, color PColor) {
	var target *uint64
	if color == WHITE {
		switch piece {
		case KING:
			target = &b.WhiteKing
		case QUEEN:
			target = &b.WhiteQueens
		case ROOK:
			target = &b.WhiteRooks
		case BISHOP:
			target = &b.WhiteBishops
		case KNIGHT:
			target = &b.WhiteKnights
		case PAWN:
			target = &b.WhitePawns
		}
	} else {
		switch piece {
		case KING:
			target = &b.BlackKing
		case QUEEN:
			target = &b.BlackQueens
		case ROOK:
			target = &b.BlackRooks
		case BISHOP:
			target = &b.BlackBishops
		case KNIGHT:
			target = &b.BlackKnights
		case PAWN:
			target = &b.BlackPawns
		}
	}
	*target &^= pos
	*target |= newPos
}
