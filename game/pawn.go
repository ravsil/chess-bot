package game

type Pawn struct {
	Color            PColor
	Pos              uint64
	Moved            bool
	CanBeEnPassanted bool
}

func (p *Pawn) GetColor() PColor { return p.Color }
func (p *Pawn) GetType() PType   { return PAWN }
func (p *Pawn) GetPos() uint64   { return p.Pos }

func (p *Pawn) GetValidMoves(b Board) uint64 {
	ocupied := b.GetOcupiedSquares()
	moves := uint64(0)
	direction := 0
	enemyPieces := uint64(0)
	if p.Color == WHITE {
		enemyPieces = b.GetPieces(BLACK)
	} else {
		direction = -16
		enemyPieces = b.GetPieces(WHITE)
	}

	forward := p.Pos << (8 + direction)
	canForward := WouldMyKingBeSafeIfIDidThisComicallyLargeFunctionCall(p.Pos, forward, b, PAWN, p.Color)
	if forward&ocupied == 0 && canForward {
		moves |= forward
	}

	if !p.Moved && canForward {
		doubleForward := forward << (8 + direction)
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
	// TODO: en-passant
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
