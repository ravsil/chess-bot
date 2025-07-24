package game

type Queen struct {
	Color PColor
	Pos   uint64
}

func (p *Queen) GetColor() PColor { return p.Color }
func (p *Queen) GetType() PType   { return QUEEN }
func (p *Queen) GetPos() uint64   { return p.Pos }

func (p *Queen) GetValidMoves(b Board) uint64 {
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
		panic("Error Getting Queen Moves")
	}
	rightUp := (p.Pos << 7) & ^RIGHT_BORDER
	if (rightUp&ocupied == 0 || rightUp&enemyPieces != 0) && WouldMyKingBeSafeIfIDidThisComicallyLargeFunctionCall(p.Pos, rightUp, b, QUEEN, p.Color) {
		moves |= GetPositiveRayAttacks(ocupied, NORTHEAST, pos)
	}
	leftUp := (p.Pos << 9) & ^LEFT_BORDER
	if (leftUp&ocupied == 0 || leftUp&enemyPieces != 0) && WouldMyKingBeSafeIfIDidThisComicallyLargeFunctionCall(p.Pos, leftUp, b, QUEEN, p.Color) {
		moves |= GetPositiveRayAttacks(ocupied, NORTHWEST, pos)
	}
	rightDown := (p.Pos >> 9) & ^RIGHT_BORDER
	if (rightDown&ocupied == 0 || rightDown&enemyPieces != 0) && WouldMyKingBeSafeIfIDidThisComicallyLargeFunctionCall(p.Pos, rightDown, b, QUEEN, p.Color) {
		moves |= GetNegativeRayAttacks(ocupied, SOUTHEAST, pos)
	}
	leftDown := (p.Pos >> 7) & ^LEFT_BORDER
	if (leftDown&ocupied == 0 || leftDown&enemyPieces != 0) && WouldMyKingBeSafeIfIDidThisComicallyLargeFunctionCall(p.Pos, leftDown, b, QUEEN, p.Color) {
		moves |= GetNegativeRayAttacks(ocupied, SOUTHWEST, pos)
	}
	right := (p.Pos >> 1) & ^RIGHT_BORDER
	if (right&ocupied == 0 || right&enemyPieces != 0) && WouldMyKingBeSafeIfIDidThisComicallyLargeFunctionCall(p.Pos, right, b, QUEEN, p.Color) {
		moves |= GetPositiveRayAttacks(ocupied, EAST, pos)
	}
	left := (p.Pos << 1) & ^LEFT_BORDER
	if (left&ocupied == 0 || left&enemyPieces != 0) && WouldMyKingBeSafeIfIDidThisComicallyLargeFunctionCall(p.Pos, left, b, QUEEN, p.Color) {
		moves |= GetNegativeRayAttacks(ocupied, WEST, pos)
	}
	up := (p.Pos << 8)
	if (up&ocupied == 0 || up&enemyPieces != 0) && WouldMyKingBeSafeIfIDidThisComicallyLargeFunctionCall(p.Pos, up, b, QUEEN, p.Color) {
		moves |= GetPositiveRayAttacks(ocupied, NORTH, pos)
	}
	down := (p.Pos >> 8)
	if (down&ocupied == 0 || down&enemyPieces != 0) && WouldMyKingBeSafeIfIDidThisComicallyLargeFunctionCall(p.Pos, down, b, QUEEN, p.Color) {
		moves |= GetNegativeRayAttacks(ocupied, SOUTH, pos)
	}
	return moves
}

func (p *Queen) Move(newPos uint64, b *Board) {
	b.Update(p.Pos, newPos, QUEEN, p.Color)
	p.Pos = newPos
}
