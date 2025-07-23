package game

type PType int

const (
	PAWN PType = iota
	ROOK
	KNIGHT
	BISHOP
	QUEEN
	KING
)

type PColor int

const (
	WHITE PColor = iota
	BLACK
)

type Piece interface {
	GetColor() PColor
	GetType() PType
	GetPos() uint64
	GetValidMoves(b Board) uint64
	Move(newPos int, b *Board)
}

func WouldMyKingBeSafeIfIDidThisComicallyLargeFunctionCall(cur uint64, desire uint64, realBoard Board, color PColor) bool {
	b := realBoard
	if desire == 0 {
		return true
	}
	c, err := GetSingleBit(cur)
	if err != nil {
		panic("Current piece position Could not be found")
	}
	d, err := GetSingleBit(desire)
	if err != nil {
		panic("Current piece position Could not be found")
	}

	b.Update(c, d)

	if color == WHITE {
		if b.GetBlackInfluence()&b.WhiteKing != 0 {
			return false
		}
	} else {
		if b.GetWhiteInfluence()&b.BlackKing != 0 {
			return false
		}
	}
	return true
}
