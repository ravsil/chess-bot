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

	pos, err := GetSingleBit(p.Pos)
	if err != nil {
		panic("Error Getting Rook Moves")
	}

	northAttacks := GetPositiveRayAttacks(ocupied, NORTH, pos)
	southAttacks := GetNegativeRayAttacks(ocupied, SOUTH, pos)
	eastAttacks := GetPositiveRayAttacks(ocupied, EAST, pos)
	westAttacks := GetNegativeRayAttacks(ocupied, WEST, pos)

	friendlyPieces := b.GetPieces(p.Color)
	northAttacks &= ^friendlyPieces
	southAttacks &= ^friendlyPieces
	eastAttacks &= ^friendlyPieces
	westAttacks &= ^friendlyPieces

	if northAttacks != 0 {
		for _, square := range GetBits(northAttacks) {
			testPos := uint64(1) << square
			if WouldMyKingBeSafeIfIDidThisComicallyLargeFunctionCall(p.Pos, testPos, b, ROOK, p.Color) {
				moves |= testPos
			}
		}
	}
	if southAttacks != 0 {
		for _, square := range GetBits(southAttacks) {
			testPos := uint64(1) << square
			if WouldMyKingBeSafeIfIDidThisComicallyLargeFunctionCall(p.Pos, testPos, b, ROOK, p.Color) {
				moves |= testPos
			}
		}
	}
	if eastAttacks != 0 {
		for _, square := range GetBits(eastAttacks) {
			testPos := uint64(1) << square
			if WouldMyKingBeSafeIfIDidThisComicallyLargeFunctionCall(p.Pos, testPos, b, ROOK, p.Color) {
				moves |= testPos
			}
		}
	}
	if westAttacks != 0 {
		for _, square := range GetBits(westAttacks) {
			testPos := uint64(1) << square
			if WouldMyKingBeSafeIfIDidThisComicallyLargeFunctionCall(p.Pos, testPos, b, ROOK, p.Color) {
				moves |= testPos
			}
		}
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
