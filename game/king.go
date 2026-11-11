package game

type King struct {
	Color PColor
	Pos   uint64
	Moved bool
}

func (p *King) GetColor() PColor { return p.Color }
func (p *King) GetType() PType   { return KING }
func (p *King) GetPos() uint64   { return p.Pos }

func (p *King) GetValidMoves(b Board) uint64 {
	ocupied := b.GetOcupiedSquares()
	moves := uint64(0)

	enemyPieces := uint64(0)
	if p.Color == WHITE {
		enemyPieces = b.GetPieces(BLACK)
	} else {
		enemyPieces = b.GetPieces(WHITE)
	}

	ms := []uint64{
		p.Pos << 8, p.Pos >> 8, (p.Pos << 1) & ^LEFT_BORDER, (p.Pos >> 1) & ^RIGHT_BORDER,
		(p.Pos << 9) & ^LEFT_BORDER, (p.Pos << 7) & ^RIGHT_BORDER, (p.Pos >> 7) & ^LEFT_BORDER,
		(p.Pos >> 9) & ^RIGHT_BORDER}

	for _, m := range ms {
		if m != 0 && (m&ocupied == 0 || (m&enemyPieces != 0)) && WouldMyKingBeSafeIfIDidThisComicallyLargeFunctionCall(p.Pos, m, b, KING, p.Color) {
			moves |= m
		}
	}

	inCheck := false
	if p.Color == WHITE {
		inCheck = b.GetBlackInfluence()&p.Pos != 0
	} else {
		inCheck = b.GetWhiteInfluence()&p.Pos != 0
	}

	if !p.Moved && !inCheck {
		if p.Color == WHITE && p.Pos == (1<<4) {
			f1, g1, h1 := uint64(1<<5), uint64(1<<6), uint64(1<<7)
			if b.WhiteRooks&h1 != 0 && (f1|g1)&ocupied == 0 &&
				WouldMyKingBeSafeIfIDidThisComicallyLargeFunctionCall(p.Pos, f1, b, KING, WHITE) &&
				WouldMyKingBeSafeIfIDidThisComicallyLargeFunctionCall(p.Pos, g1, b, KING, WHITE) {
				moves |= g1
			}
			a1, b1, c1, d1 := uint64(1<<0), uint64(1<<1), uint64(1<<2), uint64(1<<3)
			if b.WhiteRooks&a1 != 0 && (b1|c1|d1)&ocupied == 0 &&
				WouldMyKingBeSafeIfIDidThisComicallyLargeFunctionCall(p.Pos, d1, b, KING, WHITE) &&
				WouldMyKingBeSafeIfIDidThisComicallyLargeFunctionCall(p.Pos, c1, b, KING, WHITE) {
				moves |= c1
			}
		} else if p.Color == BLACK && p.Pos == (1<<60) {
			f8, g8, h8 := uint64(1<<61), uint64(1<<62), uint64(1<<63)
			if b.BlackRooks&h8 != 0 && (f8|g8)&ocupied == 0 &&
				WouldMyKingBeSafeIfIDidThisComicallyLargeFunctionCall(p.Pos, f8, b, KING, BLACK) &&
				WouldMyKingBeSafeIfIDidThisComicallyLargeFunctionCall(p.Pos, g8, b, KING, BLACK) {
				moves |= g8
			}
			a8, b8, c8, d8 := uint64(1<<56), uint64(1<<57), uint64(1<<58), uint64(1<<59)
			if b.BlackRooks&a8 != 0 && (b8|c8|d8)&ocupied == 0 &&
				WouldMyKingBeSafeIfIDidThisComicallyLargeFunctionCall(p.Pos, d8, b, KING, BLACK) &&
				WouldMyKingBeSafeIfIDidThisComicallyLargeFunctionCall(p.Pos, c8, b, KING, BLACK) {
				moves |= c8
			}
		}
	}
	return moves
}

func (p *King) Move(newPos uint64, b *Board) {
	if !p.Moved {
		p.Moved = true
	}

	b.Update(p.Pos, newPos, KING, p.Color)
	p.Pos = newPos
}
