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

	pos, err := GetSingleBit(p.Pos)
	if err != nil {
		panic("Error Getting Bishop Moves")
	}

	neAttacks := GetPositiveRayAttacks(ocupied, NORTHEAST, pos)
	nwAttacks := GetPositiveRayAttacks(ocupied, NORTHWEST, pos)
	seAttacks := GetNegativeRayAttacks(ocupied, SOUTHEAST, pos)
	swAttacks := GetNegativeRayAttacks(ocupied, SOUTHWEST, pos)

	friendlyPieces := b.GetPieces(p.Color)
	neAttacks &= ^friendlyPieces
	nwAttacks &= ^friendlyPieces
	seAttacks &= ^friendlyPieces
	swAttacks &= ^friendlyPieces

	if neAttacks != 0 {
		for _, square := range GetBits(neAttacks) {
			testPos := uint64(1) << square
			if WouldMyKingBeSafeIfIDidThisComicallyLargeFunctionCall(p.Pos, testPos, b, BISHOP, p.Color) {
				moves |= testPos
			}
		}
	}
	if nwAttacks != 0 {
		for _, square := range GetBits(nwAttacks) {
			testPos := uint64(1) << square
			if WouldMyKingBeSafeIfIDidThisComicallyLargeFunctionCall(p.Pos, testPos, b, BISHOP, p.Color) {
				moves |= testPos
			}
		}
	}
	if seAttacks != 0 {
		for _, square := range GetBits(seAttacks) {
			testPos := uint64(1) << square
			if WouldMyKingBeSafeIfIDidThisComicallyLargeFunctionCall(p.Pos, testPos, b, BISHOP, p.Color) {
				moves |= testPos
			}
		}
	}
	if swAttacks != 0 {
		for _, square := range GetBits(swAttacks) {
			testPos := uint64(1) << square
			if WouldMyKingBeSafeIfIDidThisComicallyLargeFunctionCall(p.Pos, testPos, b, BISHOP, p.Color) {
				moves |= testPos
			}
		}
	}
	return moves
}

func (p *Bishop) Move(newPos uint64, b *Board) {
	b.Update(p.Pos, newPos, BISHOP, p.Color)
	p.Pos = newPos
}
