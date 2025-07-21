package bot

type PType int

const (
	ROOK PType = iota
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

type Pieces interface {
	GetColor() PColor
	GetPos() Pos
	GetValidMoves(board [][]Pieces) ValidMoves
	Move(board [][]Pieces, pos Pos)
	LightMove(board [][]Pieces, pos Pos)
}

type Pos struct {
	X int
	Y int
}

type Pawn struct {
	Pos              Pos
	Color            PColor
	FirstMove        bool
	CanBeEnPassanted bool
}

type Piece struct {
	Pos       Pos
	PieceType PType
	Color     PColor
	Moved     bool
}

type ValidMoves struct {
	P     Pieces
	Moves []Pos
}

var blackKing *Piece
var whiteKing *Piece

func (p *Pawn) GetColor() PColor  { return p.Color }
func (p *Piece) GetColor() PColor { return p.Color }

func (p *Pawn) GetPos() Pos  { return p.Pos }
func (p *Piece) GetPos() Pos { return p.Pos }

func (p *Pawn) GetValidMoves(board [][]Pieces) ValidMoves {
	validMoves := ValidMoves{P: p, Moves: []Pos{}}
	var direction int
	if p.Color == WHITE {
		direction = 1
	} else {
		direction = -1
	}

	// forward once
	if board[p.Pos.Y+direction][p.Pos.X] == nil {
		curPos := p.Pos
		p.LightMove(board, Pos{X: p.Pos.X, Y: p.Pos.Y + direction})
		if !wouldThisKillMyKing(board, p.Color) {
			validMoves.Moves = append(validMoves.Moves, Pos{X: curPos.X, Y: curPos.Y + direction})
		}
		p.LightMove(board, curPos) // revert move
	}

	// forward twice
	if p.FirstMove && board[p.Pos.Y+2*direction][p.Pos.X] == nil && board[p.Pos.Y+direction][p.Pos.X] == nil {
		curPos := p.Pos
		p.LightMove(board, Pos{X: p.Pos.X, Y: p.Pos.Y + 2*direction})
		if !wouldThisKillMyKing(board, p.Color) {
			validMoves.Moves = append(validMoves.Moves, Pos{X: curPos.X, Y: curPos.Y + 2*direction})
		}
		p.LightMove(board, curPos) // revert move
	}

	if p.Pos.X-1 >= 0 {
		// capture left
		if board[p.Pos.Y+direction][p.Pos.X-1] != nil && board[p.Pos.Y+direction][p.Pos.X-1].GetColor() != p.Color {
			curPos := p.Pos
			p.LightMove(board, Pos{X: p.Pos.X - 1, Y: p.Pos.Y + direction})
			if !wouldThisKillMyKing(board, p.Color) {
				validMoves.Moves = append(validMoves.Moves, Pos{X: curPos.X - 1, Y: curPos.Y + direction})
			}
			p.LightMove(board, curPos) // revert move
		}
		// en passant left
		if pawn, ok := board[p.Pos.Y][p.Pos.X-1].(*Pawn); ok && pawn.CanBeEnPassanted && pawn.Color != p.Color {
			curPos := p.Pos
			p.LightMove(board, Pos{X: p.Pos.X - 1, Y: p.Pos.Y + direction})
			if !wouldThisKillMyKing(board, p.Color) {
				validMoves.Moves = append(validMoves.Moves, Pos{X: curPos.X - 1, Y: curPos.Y + direction})
			}
			p.LightMove(board, curPos) // revert move
		}
	}

	if p.Pos.X+1 < 8 {
		// capture right
		if board[p.Pos.Y+direction][p.Pos.X+1] != nil && board[p.Pos.Y+direction][p.Pos.X+1].GetColor() != p.Color {
			curPos := p.Pos
			p.LightMove(board, Pos{X: p.Pos.X + 1, Y: p.Pos.Y + direction})
			if !wouldThisKillMyKing(board, p.Color) {
				validMoves.Moves = append(validMoves.Moves, Pos{X: curPos.X + 1, Y: curPos.Y + direction})
			}
			p.LightMove(board, curPos) // revert move
		}
		// en passant right
		if pawn, ok := board[p.Pos.Y][p.Pos.X+1].(*Pawn); ok && pawn.CanBeEnPassanted && pawn.Color != p.Color {
			curPos := p.Pos
			p.LightMove(board, Pos{X: p.Pos.X + 1, Y: p.Pos.Y + direction})
			if !wouldThisKillMyKing(board, p.Color) {
				validMoves.Moves = append(validMoves.Moves, Pos{X: curPos.X + 1, Y: curPos.Y + direction})
			}
			p.LightMove(board, curPos) // revert move
		}
	}

	return validMoves
}

func (p *Piece) GetValidMoves(board [][]Pieces) ValidMoves {
	validMoves := ValidMoves{P: p, Moves: []Pos{}}
	if p.PieceType == KNIGHT {
		moves := []Pos{{2, 1}, {2, -1}, {-2, 1}, {-2, -1}, {1, 2}, {1, -2}, {-1, 2}, {-1, -2}}
		for _, move := range moves {
			if p.Pos.X+move.X >= 0 && p.Pos.X+move.X < 8 && p.Pos.Y+move.Y >= 0 && p.Pos.Y+move.Y < 8 {
				if board[p.Pos.Y+move.Y][p.Pos.X+move.X] == nil || board[p.Pos.Y+move.Y][p.Pos.X+move.X].GetColor() != p.Color {
					m := Pos{p.Pos.X + move.X, p.Pos.Y + move.Y}
					curPos := p.Pos
					p.LightMove(board, m)
					if !wouldThisKillMyKing(board, p.Color) {
						validMoves.Moves = append(validMoves.Moves, m)
					}
					p.LightMove(board, curPos) // revert move
				}
			}
		}
	} else if p.PieceType == KING {
		moves := []Pos{{1, 0}, {-1, 0}, {0, 1}, {0, -1}, {1, 1}, {1, -1}, {-1, 1}, {-1, -1}}
		for _, move := range moves {
			if p.Pos.X+move.X >= 0 && p.Pos.X+move.X < 8 && p.Pos.Y+move.Y >= 0 && p.Pos.Y+move.Y < 8 {
				if board[p.Pos.Y+move.Y][p.Pos.X+move.X] == nil || board[p.Pos.Y+move.Y][p.Pos.X+move.X].GetColor() != p.Color {
					m := Pos{p.Pos.X + move.X, p.Pos.Y + move.Y}
					curPos := p.Pos
					p.LightMove(board, m)
					if !wouldThisKillMyKing(board, p.Color) {
						validMoves.Moves = append(validMoves.Moves, m)
					}
					p.LightMove(board, curPos) // revert move
				}
			}
		}
	} else {
		var dirs []struct{ dx, dy int }
		if p.PieceType == ROOK || p.PieceType == QUEEN {
			dirs = append(dirs, []struct{ dx, dy int }{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}...)
		}
		if p.PieceType == BISHOP || p.PieceType == QUEEN {
			dirs = append(dirs, []struct{ dx, dy int }{{1, 1}, {-1, 1}, {1, -1}, {-1, -1}}...)
		}

		for _, d := range dirs {
			for i := 1; i < 8; i++ {
				X, Y := p.Pos.X+d.dx*i, p.Pos.Y+d.dy*i
				if X < 0 || X >= 8 || Y < 0 || Y >= 8 {
					break
				}
				if board[Y][X] == nil {
					curPos := p.Pos
					p.LightMove(board, Pos{X: X, Y: Y})
					if !wouldThisKillMyKing(board, p.Color) {
						validMoves.Moves = append(validMoves.Moves, Pos{X: X, Y: Y})
					}
					p.LightMove(board, curPos) // revert move
				} else {
					if board[Y][X].GetColor() != p.Color {
						curPos := p.Pos
						p.LightMove(board, Pos{X: X, Y: Y})
						if !wouldThisKillMyKing(board, p.Color) {
							validMoves.Moves = append(validMoves.Moves, Pos{X: X, Y: Y})
						}
						p.LightMove(board, curPos) // revert move
					}
					break
				}
			}
		}

	}
	return validMoves
}

func (p *Pawn) LightMove(board [][]Pieces, destination Pos) {
	board[p.Pos.Y][p.Pos.X] = nil
	p.Pos = destination
	board[p.Pos.Y][p.Pos.X] = p
}

func (p *Piece) LightMove(board [][]Pieces, destination Pos) {
	board[p.Pos.Y][p.Pos.X] = nil
	p.Pos = destination
	board[p.Pos.Y][p.Pos.X] = p
}

func (p *Pawn) Move(board [][]Pieces, destination Pos) {
	p.LightMove(board, destination)
	if destination.Y == p.Pos.Y+2 || destination.Y == p.Pos.Y-2 {
		p.FirstMove = false
		p.CanBeEnPassanted = true
	} else {
		p.FirstMove = false
		p.CanBeEnPassanted = false
	}
}

func (p *Piece) Move(board [][]Pieces, destination Pos) {
	p.LightMove(board, destination)
	p.Moved = true
}

func wouldThisKillMyKing(board [][]Pieces, player PColor) bool {
	var king *Piece
	if player == WHITE {
		king = whiteKing
	} else {
		king = blackKing
	}
	if king == nil {
		panic("King not found")
	}

	// ↑
	for i := king.Pos.Y + 1; i < 8; i++ {
		if board[i][king.Pos.X] != nil {
			if board[i][king.Pos.X].GetColor() == player {
				break
			}
			if piece, ok := board[i][king.Pos.X].(*Piece); ok && (piece.PieceType == ROOK || piece.PieceType == QUEEN) {
				return true
			}
		}
	}
	// ↓
	for i := king.Pos.Y - 1; i >= 0; i-- {
		if board[i][king.Pos.X] != nil {
			if board[i][king.Pos.X].GetColor() == player {
				break
			}
			if piece, ok := board[i][king.Pos.X].(*Piece); ok && (piece.PieceType == ROOK || piece.PieceType == QUEEN) {
				return true
			}
		}
	}
	// →
	for i := king.Pos.X + 1; i < 8; i++ {
		if board[king.Pos.Y][i] != nil {
			if board[king.Pos.Y][i].GetColor() == player {
				break
			}
			if piece, ok := board[king.Pos.Y][i].(*Piece); ok &&
				(piece.PieceType == ROOK || piece.PieceType == QUEEN) {
				return true
			}
		}
	}
	// ←
	for i := king.Pos.X - 1; i >= 0; i-- {
		if board[king.Pos.Y][i] != nil {
			if board[king.Pos.Y][i].GetColor() == player {
				break
			}
			if piece, ok := board[king.Pos.Y][i].(*Piece); ok &&
				(piece.PieceType == ROOK || piece.PieceType == QUEEN) {
				return true
			}
		}
	}
	// ↗
	for i, j := king.Pos.X+1, king.Pos.Y+1; i < 8 && j < 8; i, j = i+1, j+1 {
		if board[j][i] != nil {
			if board[j][i].GetColor() == player {
				break
			}
			if piece, ok := board[j][i].(*Piece); ok && (piece.PieceType == BISHOP || piece.PieceType == QUEEN) {
				return true
			}
			if pawn, ok := board[j][i].(*Pawn); ok && pawn.Color != player && j == king.Pos.Y+1 && i == king.Pos.X+1 {
				return true
			}
		}
	}
	// ↖
	for i, j := king.Pos.X-1, king.Pos.Y+1; i >= 0 && j < 8; i, j = i-1, j+1 {
		if board[j][i] != nil {
			if board[j][i].GetColor() == player {
				break
			}
			if piece, ok := board[j][i].(*Piece); ok && (piece.PieceType == BISHOP || piece.PieceType == QUEEN) {
				return true
			}
			if pawn, ok := board[j][i].(*Pawn); ok && pawn.Color != player && j == king.Pos.Y+1 && i == king.Pos.X-1 {
				return true
			}
		}
	}
	// ↘
	for i, j := king.Pos.X+1, king.Pos.Y-1; i < 8 && j >= 0; i, j = i+1, j-1 {
		if board[j][i] != nil {
			if board[j][i].GetColor() == player {
				break
			}
			if piece, ok := board[j][i].(*Piece); ok && (piece.PieceType == BISHOP || piece.PieceType == QUEEN) {
				return true
			}
		}
	}
	// ↙
	for i, j := king.Pos.X-1, king.Pos.Y-1; i >= 0 && j >= 0; i, j = i-1, j-1 {
		if board[j][i] != nil {
			if board[j][i].GetColor() == player {
				break
			}
			if piece, ok := board[j][i].(*Piece); ok && (piece.PieceType == BISHOP || piece.PieceType == QUEEN) {
				return true
			}
		}
	}
	return false
}

func genBoard() [][]Pieces {
	board := make([][]Pieces, 8)
	for i := range board {
		board[i] = make([]Pieces, 8)
	}

	for i := 0; i < 8; i++ {
		board[1][i] = &Pawn{Pos: Pos{X: i, Y: 1}, Color: WHITE, FirstMove: true, CanBeEnPassanted: false}
		board[6][i] = &Pawn{Pos: Pos{X: i, Y: 6}, Color: BLACK, FirstMove: true, CanBeEnPassanted: false}
	}

	pieceTypes := []PType{ROOK, KNIGHT, BISHOP, QUEEN, KING, BISHOP, KNIGHT, ROOK}
	for i, pt := range pieceTypes {
		board[0][i] = &Piece{Pos: Pos{X: i, Y: 0}, PieceType: pt, Color: WHITE, Moved: false}
		board[7][i] = &Piece{Pos: Pos{X: i, Y: 7}, PieceType: pt, Color: BLACK, Moved: false}
	}
	whiteKing = board[0][4].(*Piece)
	blackKing = board[7][4].(*Piece)
	return board
}

// this function will return the list of valid moves for a given player
func getValidMoves(board [][]Pieces, player PColor) []ValidMoves {
	validMoves := []ValidMoves{}

	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			if board[i][j] == nil {
				continue
			}

			if piece, ok := board[i][j].(*Piece); ok && piece.Color == player {
				pieceMoves := piece.GetValidMoves(board)
				if len(pieceMoves.Moves) > 0 {
					validMoves = append(validMoves, pieceMoves)
				}
			} else if pawn, ok := board[i][j].(*Pawn); ok && pawn.Color == player {
				pawnMoves := pawn.GetValidMoves(board)
				if len(pawnMoves.Moves) > 0 {
					validMoves = append(validMoves, pawnMoves)
				}
			}
		}
	}

	return validMoves
}
