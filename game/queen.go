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

	pos, err := GetSingleBit(p.Pos)
	if err != nil {
		panic("Error Getting Queen Moves")
	}

	neAttacks := GetPositiveRayAttacks(ocupied, NORTHEAST, pos)
	nwAttacks := GetPositiveRayAttacks(ocupied, NORTHWEST, pos)
	seAttacks := GetNegativeRayAttacks(ocupied, SOUTHEAST, pos)
	swAttacks := GetNegativeRayAttacks(ocupied, SOUTHWEST, pos)

	northAttacks := GetPositiveRayAttacks(ocupied, NORTH, pos)
	southAttacks := GetNegativeRayAttacks(ocupied, SOUTH, pos)
	eastAttacks := GetPositiveRayAttacks(ocupied, EAST, pos)
	westAttacks := GetNegativeRayAttacks(ocupied, WEST, pos)

	friendlyPieces := b.GetPieces(p.Color)
	neAttacks &= ^friendlyPieces
	nwAttacks &= ^friendlyPieces
	seAttacks &= ^friendlyPieces
	swAttacks &= ^friendlyPieces
	northAttacks &= ^friendlyPieces
	southAttacks &= ^friendlyPieces
	eastAttacks &= ^friendlyPieces
	westAttacks &= ^friendlyPieces

	allDirections := []uint64{neAttacks, nwAttacks, seAttacks, swAttacks, northAttacks, southAttacks, eastAttacks, westAttacks}
	for _, directionAttacks := range allDirections {
		if directionAttacks != 0 {
			for _, square := range GetBits(directionAttacks) {
				testPos := uint64(1) << square
				if WouldMyKingBeSafeIfIDidThisComicallyLargeFunctionCall(p.Pos, testPos, b, QUEEN, p.Color) {
					moves |= testPos
				}
			}
		}
	}
	return moves
}

func (p *Queen) Move(newPos uint64, b *Board) {
	b.Update(p.Pos, newPos, QUEEN, p.Color)
	p.Pos = newPos
}
