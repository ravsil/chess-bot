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
		if (m&ocupied == 0 || (m&enemyPieces != 0)) && WouldMyKingBeSafeIfIDidThisComicallyLargeFunctionCall(p.Pos, m, b, KING, p.Color) {
			moves |= m
		}
	}

	if !p.Moved {
		if p.Color == WHITE {
			lk1 := uint64(0b00000000_00000000_00000000_00000000_00000000_00000000_00000000_00100000)
			lk2 := uint64(0b00000000_00000000_00000000_00000000_00000000_00000000_00000000_01000000)
			if (lk1|lk2)&ocupied == 0 &&
				WouldMyKingBeSafeIfIDidThisComicallyLargeFunctionCall(p.Pos, lk1, b, KING, WHITE) &&
				WouldMyKingBeSafeIfIDidThisComicallyLargeFunctionCall(p.Pos, lk2, b, KING, WHITE) {
				moves |= lk2
			}
			rk1 := uint64(0b00000000_00000000_00000000_00000000_00000000_00000000_00000000_00001000)
			rk2 := uint64(0b00000000_00000000_00000000_00000000_00000000_00000000_00000000_00000100)
			rk3 := uint64(0b00000000_00000000_00000000_00000000_00000000_00000000_00000000_00000010)
			if (rk1|rk2|rk3)&ocupied == 0 &&
				WouldMyKingBeSafeIfIDidThisComicallyLargeFunctionCall(p.Pos, rk1, b, KING, WHITE) &&
				WouldMyKingBeSafeIfIDidThisComicallyLargeFunctionCall(p.Pos, rk2, b, KING, WHITE) &&
				WouldMyKingBeSafeIfIDidThisComicallyLargeFunctionCall(p.Pos, rk3, b, KING, WHITE) {
				moves |= rk2
			}
		} else {
			lk1 := uint64(0b00100000_00000000_00000000_00000000_00000000_00000000_00000000_00000000)
			lk2 := uint64(0b01000000_00000000_00000000_00000000_00000000_00000000_00000000_00000000)
			if (lk1|lk2)&ocupied == 0 &&
				WouldMyKingBeSafeIfIDidThisComicallyLargeFunctionCall(p.Pos, lk1, b, KING, BLACK) &&
				WouldMyKingBeSafeIfIDidThisComicallyLargeFunctionCall(p.Pos, lk2, b, KING, BLACK) {
				moves |= lk2
			}
			rk1 := uint64(0b00001000_00000000_00000000_00000000_00000000_00000000_00000000_00000000)
			rk2 := uint64(0b00000100_00000000_00000000_00000000_00000000_00000000_00000000_00000000)
			rk3 := uint64(0b00000010_00000000_00000000_00000000_00000000_00000000_00000000_00000000)
			if (rk1|rk2|rk3)&ocupied == 0 &&
				WouldMyKingBeSafeIfIDidThisComicallyLargeFunctionCall(p.Pos, rk1, b, KING, BLACK) &&
				WouldMyKingBeSafeIfIDidThisComicallyLargeFunctionCall(p.Pos, rk2, b, KING, BLACK) &&
				WouldMyKingBeSafeIfIDidThisComicallyLargeFunctionCall(p.Pos, rk3, b, KING, BLACK) {
				moves |= rk2
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
