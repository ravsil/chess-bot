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

// returns the squares black pieces are attacking
func (b *Board) GetWhiteInfluence() uint64 {
	influence := uint64(0)
	const leftBorder uint64 = 0b00000001_00000001_00000001_00000001_00000001_00000001_00000001_00000001
	const rightBorder uint64 = 0b10000000_10000000_10000000_10000000_10000000_10000000_10000000_10000000
	whitePieces := b.WhitePawns | b.WhiteRooks | b.WhiteKnights | b.WhiteBishops | b.WhiteQueens | b.WhiteKing
	blackPieces := b.BlackPawns | b.BlackRooks | b.BlackKnights | b.BlackBishops | b.BlackQueens | b.BlackKing
	allPieces := whitePieces | blackPieces

	influence |= (b.WhitePawns << 9) & ^leftBorder  // ↖
	influence |= (b.WhitePawns << 7) & ^rightBorder // ↗

	influence |= (b.WhiteKing << 8) | // ↑
		(b.WhiteKing >> 8) | // ↓
		((b.WhiteKing << 1) & ^leftBorder) | // ←
		((b.WhiteKing >> 1) & ^rightBorder) | // →
		((b.WhiteKing << 9) & ^leftBorder) | // ↖
		((b.WhiteKing << 7) & ^rightBorder) | // ↗
		((b.WhiteKing >> 7) & ^leftBorder) | // ↙
		((b.WhiteKing >> 9) & ^rightBorder) // ↘

	influence |= (b.WhiteKnights << 15) & ^leftBorder  // ↓↓←
	influence |= (b.WhiteKnights << 17) & ^rightBorder // ↓↓→
	influence |= (b.WhiteKnights >> 15) & ^rightBorder // ↑↑→
	influence |= (b.WhiteKnights >> 17) & ^leftBorder  // ↑↑←
	influence |= (b.WhiteKnights << 6) & ^leftBorder   // ←←↓
	influence |= (b.WhiteKnights << 10) & ^rightBorder // →→
	influence |= (b.WhiteKnights >> 6) & ^rightBorder  // →→↑
	influence |= (b.WhiteKnights >> 10) & ^leftBorder  // ←←

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
	// Bishop and Rook attacks treats all pieces as capturable,
	// so I need to remove the same color pieces to make it work
	influence = influence & ^whitePieces
	return influence
}

func (b *Board) GetBlackInfluence() uint64 {
	influence := uint64(0)
	const leftBorder uint64 = 0b00000001_00000001_00000001_00000001_00000001_00000001_00000001_00000001
	const rightBorder uint64 = 0b10000000_10000000_10000000_10000000_10000000_10000000_10000000_10000000
	whitePieces := b.WhitePawns | b.WhiteRooks | b.WhiteKnights | b.WhiteBishops | b.WhiteQueens | b.WhiteKing
	blackPieces := b.BlackPawns | b.BlackRooks | b.BlackKnights | b.BlackBishops | b.BlackQueens | b.BlackKing
	allPieces := whitePieces | blackPieces

	influence |= (b.BlackPawns << 9) & ^leftBorder  // ↖
	influence |= (b.BlackPawns << 7) & ^rightBorder // ↗

	influence |= (b.BlackKing << 8) | // ↑
		(b.BlackKing >> 8) | // ↓
		((b.BlackKing << 1) & ^leftBorder) | // ←
		((b.BlackKing >> 1) & ^rightBorder) | // →
		((b.BlackKing << 9) & ^leftBorder) | // ↖
		((b.BlackKing << 7) & ^rightBorder) | // ↗
		((b.BlackKing >> 7) & ^leftBorder) | // ↙
		((b.BlackKing >> 9) & ^rightBorder) // ↘

	influence |= (b.BlackKnights << 15) & ^leftBorder  // ↓↓←
	influence |= (b.BlackKnights << 17) & ^rightBorder // ↓↓→
	influence |= (b.BlackKnights >> 15) & ^rightBorder // ↑↑→
	influence |= (b.BlackKnights >> 17) & ^leftBorder  // ↑↑←
	influence |= (b.BlackKnights << 6) & ^leftBorder   // ←←↓
	influence |= (b.BlackKnights << 10) & ^rightBorder // →→
	influence |= (b.BlackKnights >> 6) & ^rightBorder  // →→↑
	influence |= (b.BlackKnights >> 10) & ^leftBorder  // ←←

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
	// Bishop and Rook attacks treats all pieces as capturable,
	// so I need to remove the same color pieces to make it work
	influence = influence & ^blackPieces
	return influence
}

func (b *Board) Update(pos int, newPos int) {
	if pos < 0 || pos > 63 || newPos < 0 || newPos > 63 {
		return
	}

	switch {
	case b.WhitePawns&(1<<pos) != 0:
		b.WhitePawns &^= (1 << pos)
		b.WhitePawns |= (1 << newPos)
	case b.BlackPawns&(1<<pos) != 0:
		b.BlackPawns &^= (1 << pos)
		b.BlackPawns |= (1 << newPos)
	case b.WhiteRooks&(1<<pos) != 0:
		b.WhiteRooks &^= (1 << pos)
		b.WhiteRooks |= (1 << newPos)
	case b.BlackRooks&(1<<pos) != 0:
		b.BlackRooks &^= (1 << pos)
		b.BlackRooks |= (1 << newPos)
	case b.WhiteKnights&(1<<pos) != 0:
		b.WhiteKnights &^= (1 << pos)
		b.WhiteKnights |= (1 << newPos)
	case b.BlackKnights&(1<<pos) != 0:
		b.BlackKnights &^= (1 << pos)
		b.BlackKnights |= (1 << newPos)
	case b.WhiteBishops&(1<<pos) != 0:
		b.WhiteBishops &^= (1 << pos)
		b.WhiteBishops |= (1 << newPos)
	case b.BlackBishops&(1<<pos) != 0:
		b.BlackBishops &^= (1 << pos)
		b.BlackBishops |= (1 << newPos)
	case b.WhiteQueens&(1<<pos) != 0:
		b.WhiteQueens &^= (1 << pos)
		b.WhiteQueens |= (1 << newPos)
	case b.BlackQueens&(1<<pos) != 0:
		b.BlackQueens &^= (1 << pos)
		b.BlackQueens |= (1 << newPos)
	case b.WhiteKing&(1<<pos) != 0:
		b.WhiteKing &^= (1 << pos)
		b.WhiteKing |= (1 << newPos)
	case b.BlackKing&(1<<pos) != 0:
		b.BlackKing &^= (1 << pos)
		b.BlackKing |= (1 << newPos)
	default:
		return
	}
}
