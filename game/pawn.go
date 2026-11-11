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

	if forward != 0 && forward&ocupied == 0 {
		if WouldMyKingBeSafeIfIDidThisComicallyLargeFunctionCall(p.Pos, forward, b, PAWN, p.Color) {
			moves |= forward
		}

		if !p.Moved {
			doubleForward := uint64(0)
			if p.Color == WHITE {
				doubleForward = forward << 8
			} else {
				doubleForward = forward >> 8
			}

			if doubleForward != 0 && doubleForward&ocupied == 0 {
				if WouldMyKingBeSafeIfIDidThisComicallyLargeFunctionCall(p.Pos, doubleForward, b, PAWN, p.Color) {
					moves |= doubleForward
				}
			}
		}
	}

	var leftCapture, rightCapture uint64
	if p.Color == WHITE {
		leftCapture = (p.Pos << 7) & ^RIGHT_BORDER
		rightCapture = (p.Pos << 9) & ^LEFT_BORDER
	} else {
		leftCapture = (p.Pos >> 9) & ^RIGHT_BORDER
		rightCapture = (p.Pos >> 7) & ^LEFT_BORDER
	}

	if leftCapture&enemyPieces != 0 && WouldMyKingBeSafeIfIDidThisComicallyLargeFunctionCall(p.Pos, leftCapture, b, PAWN, p.Color) {
		moves |= leftCapture
	}
	if rightCapture&enemyPieces != 0 && WouldMyKingBeSafeIfIDidThisComicallyLargeFunctionCall(p.Pos, rightCapture, b, PAWN, p.Color) {
		moves |= rightCapture
	}

	if b.EnPassant != 0 {
		if leftCapture == b.EnPassant && WouldMyKingBeSafeIfIAnPassanted(p.Pos, leftCapture, b, p.Color) {
			moves |= leftCapture
		}
		if rightCapture == b.EnPassant && WouldMyKingBeSafeIfIAnPassanted(p.Pos, rightCapture, b, p.Color) {
			moves |= rightCapture
		}
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
