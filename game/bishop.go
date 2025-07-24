package game

type Bishop struct {
	Color PColor
	Pos   uint64
}

func (p *Bishop) GetColor() PColor { return p.Color }
func (p *Bishop) GetType() PType   { return BISHOP }
func (p *Bishop) GetPos() uint64   { return p.Pos }

func (p *Bishop) GetValidMoves(b Board) uint64 {
	ocupied := b.GetOcupiedSquares()
	moves := uint64(0)

	enemyPieces := uint64(0)
	if p.Color == WHITE {
		enemyPieces = b.GetPieces(BLACK)
	} else {
		enemyPieces = b.GetPieces(WHITE)
	}

	pos, err := GetSingleBit(p.Pos)
	if err != nil {
		panic("Error Getting Bishop Moves")
	}
	rightUp := (p.Pos << 7) & ^RIGHT_BORDER
	if (rightUp&ocupied == 0 || rightUp&enemyPieces != 0) && WouldMyKingBeSafeIfIDidThisComicallyLargeFunctionCall(p.Pos, rightUp, b, BISHOP, p.Color) {
		moves |= GetPositiveRayAttacks(ocupied, NORTHEAST, pos)
	}
	leftUp := (p.Pos << 9) & ^LEFT_BORDER
	if (leftUp&ocupied == 0 || leftUp&enemyPieces != 0) && WouldMyKingBeSafeIfIDidThisComicallyLargeFunctionCall(p.Pos, leftUp, b, BISHOP, p.Color) {
		moves |= GetPositiveRayAttacks(ocupied, NORTHWEST, pos)
	}
	rightDown := (p.Pos >> 9) & ^RIGHT_BORDER
	if (rightDown&ocupied == 0 || rightDown&enemyPieces != 0) && WouldMyKingBeSafeIfIDidThisComicallyLargeFunctionCall(p.Pos, rightDown, b, BISHOP, p.Color) {
		moves |= GetNegativeRayAttacks(ocupied, SOUTHEAST, pos)
	}
	leftDown := (p.Pos >> 7) & ^LEFT_BORDER
	if (leftDown&ocupied == 0 || leftDown&enemyPieces != 0) && WouldMyKingBeSafeIfIDidThisComicallyLargeFunctionCall(p.Pos, leftDown, b, BISHOP, p.Color) {
		moves |= GetNegativeRayAttacks(ocupied, SOUTHWEST, pos)
	}

	return moves
}

func (p *Bishop) Move(newPos uint64, b *Board) {
	b.Update(p.Pos, newPos, BISHOP, p.Color)
	p.Pos = newPos
}
