package game

type Pawn struct {
	Color            PColor
	Pos              uint64
	Moved            bool
	CanBeEnPassanted bool
	WantsToEnPassant bool
}

func (p *Pawn) GetColor() PColor { return p.Color }
func (p *Pawn) GetType() PType   { return PAWN }
func (p *Pawn) GetPos() uint64   { return p.Pos }

func (p *Pawn) GetValidMoves(b Board) uint64 {
	ocupied := b.GetOcupiedSquares()
	moves := uint64(0)
	enemyPieces := uint64(0)
	forward := uint64(0)
	if p.Color == WHITE {
		enemyPieces = b.GetPieces(BLACK)
		forward = p.Pos << 8
	} else {
		forward = p.Pos >> 8
		enemyPieces = b.GetPieces(WHITE)
	}

	canForward := WouldMyKingBeSafeIfIDidThisComicallyLargeFunctionCall(p.Pos, forward, b, PAWN, p.Color)
	if forward&ocupied == 0 && canForward {
		moves |= forward
	}

	if !p.Moved && canForward {
		doubleForward := uint64(0)
		if p.Color == WHITE {
			doubleForward = (forward << 8) | forward

		} else {
			doubleForward = (forward >> 8) | forward
		}

		if doubleForward&ocupied == 0 {
			moves |= doubleForward
		}
	}

	left := (forward << 1) & ^LEFT_BORDER
	if left&enemyPieces != 0 && WouldMyKingBeSafeIfIDidThisComicallyLargeFunctionCall(p.Pos, left, b, PAWN, p.Color) {
		moves |= left
	}
	right := (forward >> 1) & ^RIGHT_BORDER
	if right&enemyPieces != 0 && WouldMyKingBeSafeIfIDidThisComicallyLargeFunctionCall(p.Pos, right, b, PAWN, p.Color) {
		moves |= right
	}

	if ((p.Pos<<1) & ^LEFT_BORDER)&enemyPieces != 0 && WouldMyKingBeSafeIfIAnPassanted(p.Pos, left, b, p.Color) {
		p.WantsToEnPassant = true
		moves |= left
	} else if ((p.Pos>>1) & ^RIGHT_BORDER)&enemyPieces != 0 && WouldMyKingBeSafeIfIAnPassanted(p.Pos, right, b, p.Color) {
		p.WantsToEnPassant = true
		moves |= right
	} else {
		p.WantsToEnPassant = false
	}
	return moves
}

func (p *Pawn) Move(newPos uint64, b *Board) {
	if !p.Moved {
		p.Moved = true
	}
	if (p.Pos<<16)&newPos != 0 || (p.Pos>>16)&newPos != 0 {
		p.CanBeEnPassanted = true
	}

	b.Update(p.Pos, newPos, PAWN, p.Color)
	p.Pos = newPos
}
