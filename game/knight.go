package game

type Knight struct {
	Color PColor
	Pos   uint64
}

func (p *Knight) GetColor() PColor { return p.Color }
func (p *Knight) GetType() PType   { return KNIGHT }
func (p *Knight) GetPos() uint64   { return p.Pos }

func (p *Knight) GetValidMoves(b Board) uint64 {
	ocupied := b.GetOcupiedSquares()
	moves := uint64(0)

	enemyPieces := uint64(0)
	if p.Color == WHITE {
		enemyPieces = b.GetPieces(BLACK)
	} else {
		enemyPieces = b.GetPieces(WHITE)
	}

	m1 := (p.Pos << 15) & ^LEFT_BORDER // ↓↓←
	notPinned := WouldMyKingBeSafeIfIDidThisComicallyLargeFunctionCall(p.Pos, m1, b, KNIGHT, p.Color)
	if notPinned {
		m2 := (p.Pos << 17) & ^RIGHT_BORDER        // ↓↓→
		m3 := (p.Pos >> 15) & ^RIGHT_BORDER        // ↑↑→
		m4 := (p.Pos >> 17) & ^LEFT_BORDER         // ↑↑←
		m5 := (p.Pos << 6) & ^DOUBLE_LEFT_BORDER   // ←←
		m6 := (p.Pos << 10) & ^DOUBLE_RIGHT_BORDER // →→↓
		m7 := (p.Pos >> 6) & ^DOUBLE_RIGHT_BORDER  // →→↑
		m8 := (p.Pos >> 10) & ^DOUBLE_LEFT_BORDER  // ←
		if m1&ocupied == 0 || m1&enemyPieces != 0 {
			moves |= m1
		}
		if m2&ocupied == 0 || m2&enemyPieces != 0 {
			moves |= m2
		}
		if m3&ocupied == 0 || m3&enemyPieces != 0 {
			moves |= m3
		}
		if m4&ocupied == 0 || m4&enemyPieces != 0 {
			moves |= m4
		}
		if m5&ocupied == 0 || m5&enemyPieces != 0 {
			moves |= m5
		}
		if m6&ocupied == 0 || m6&enemyPieces != 0 {
			moves |= m6
		}
		if m7&ocupied == 0 || m7&enemyPieces != 0 {
			moves |= m7
		}
		if m8&ocupied == 0 || m8&enemyPieces != 0 {
			moves |= m8
		}
	}
	return moves
}

func (p *Knight) Move(newPos uint64, b *Board) {
	b.Update(p.Pos, newPos, KNIGHT, p.Color)
	p.Pos = newPos
}
