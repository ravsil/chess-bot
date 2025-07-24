package game

type Rook struct {
	Color PColor
	Pos   uint64
	Moved bool
}

func (p *Rook) GetColor() PColor { return p.Color }
func (p *Rook) GetType() PType   { return ROOK }
func (p *Rook) GetPos() uint64   { return p.Pos }

func (p *Rook) GetValidMoves(b Board) uint64 {
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
		panic("Error Getting Rook Moves")
	}
	right := (p.Pos >> 1) & ^RIGHT_BORDER
	if (right&ocupied == 0 || right&enemyPieces != 0) && WouldMyKingBeSafeIfIDidThisComicallyLargeFunctionCall(p.Pos, right, b, ROOK, p.Color) {
		moves |= GetPositiveRayAttacks(ocupied, EAST, pos)
	}
	left := (p.Pos << 1) & ^LEFT_BORDER
	if (left&ocupied == 0 || left&enemyPieces != 0) && WouldMyKingBeSafeIfIDidThisComicallyLargeFunctionCall(p.Pos, left, b, ROOK, p.Color) {
		moves |= GetNegativeRayAttacks(ocupied, WEST, pos)
	}
	up := (p.Pos << 8)
	if (up&ocupied == 0 || up&enemyPieces != 0) && WouldMyKingBeSafeIfIDidThisComicallyLargeFunctionCall(p.Pos, up, b, ROOK, p.Color) {
		moves |= GetPositiveRayAttacks(ocupied, NORTH, pos)
	}
	down := (p.Pos >> 8)
	if (down&ocupied == 0 || down&enemyPieces != 0) && WouldMyKingBeSafeIfIDidThisComicallyLargeFunctionCall(p.Pos, down, b, ROOK, p.Color) {
		moves |= GetNegativeRayAttacks(ocupied, SOUTH, pos)
	}

	return moves
}

func (p *Rook) Move(newPos uint64, b *Board) {
	if !p.Moved {
		p.Moved = true
	}

	b.Update(p.Pos, newPos, ROOK, p.Color)
	p.Pos = newPos
}
