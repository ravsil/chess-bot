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

	// TODO: castling
	return moves
}

func (p *King) Move(newPos uint64, b *Board) {
	if !p.Moved {
		p.Moved = true
	}

	b.Update(p.Pos, newPos, KING, p.Color)
	p.Pos = newPos
}
