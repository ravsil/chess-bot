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
	Move(newPos uint64, b *Board)
}

func WouldMyKingBeSafeIfIDidThisComicallyLargeFunctionCall(cur uint64, desire uint64, realBoard Board, piece PType, color PColor) bool {
	b := realBoard

	b.Update(cur, desire, piece, color)

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
