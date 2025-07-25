package game

type Knight struct {
	Color PColor
	Pos   uint64
}

func (p *Knight) GetColor() PColor { return p.Color }
func (p *Knight) GetType() PType   { return KNIGHT }
func (p *Knight) GetPos() uint64   { return p.Pos }

func (p *Knight) GetValidMoves(b Board) uint64 {
	moves := uint64(0)
	friendlyPieces := b.GetPieces(p.Color)

	attacks := uint64(0)

	attacks |= (p.Pos << 17) & ^LEFT_BORDER
	attacks |= (p.Pos << 15) & ^RIGHT_BORDER
	attacks |= (p.Pos >> 17) & ^RIGHT_BORDER
	attacks |= (p.Pos >> 15) & ^LEFT_BORDER
	attacks |= (p.Pos << 10) & ^DOUBLE_LEFT_BORDER
	attacks |= (p.Pos << 6) & ^DOUBLE_RIGHT_BORDER
	attacks |= (p.Pos >> 6) & ^DOUBLE_LEFT_BORDER
	attacks |= (p.Pos >> 10) & ^DOUBLE_RIGHT_BORDER
	attacks &= ^friendlyPieces

	for attacks != 0 {
		targetSquare := attacks & -attacks
		attacks &= attacks - 1
		if WouldMyKingBeSafeIfIDidThisComicallyLargeFunctionCall(p.Pos, targetSquare, b, KNIGHT, p.Color) {
			moves |= targetSquare
		}
	}
	return moves
}

func (p *Knight) Move(newPos uint64, b *Board) {
	b.Update(p.Pos, newPos, KNIGHT, p.Color)
	p.Pos = newPos
}
