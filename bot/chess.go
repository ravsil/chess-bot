package main

import (
	"fmt"
)


type PieceType int
const (
	ROOK PieceType = iota
	KNIGHT
	BISHOP
	QUEEN
	KING
)

type Color int
const (
	WHITE Color = iota
	BLACK
)

type Pieces interface {
	GetColor() Color
	GetValidMoves(board [][]Pieces, player Color) ValidMoves
}

type Pos struct {
	x int
	y int
}

type Pawn struct {
	pos Pos
	color Color
	firstMove bool
	canBeEnPassanted bool
}

type Piece struct {
	pos Pos
	pieceType PieceType
	color Color
	moved bool
}

type ValidMoves struct {
	p Pieces
	moves []Pos
}

func (p *Pawn) GetColor() Color { return p.color  }
func (p *Piece) GetColor() Color { return p.color  }

func (p *Pawn) GetValidMoves(board [][]Pieces, player Color) ValidMoves {
	validMoves := ValidMoves{ p: p, moves: []Pos{} }
	var direction int
	if p.color == WHITE {
		direction = 1
	} else {
		direction = -1
	}
	
	// forward once
	if board[p.pos.y + direction][p.pos.x] == nil {
		validMoves.moves = append(validMoves.moves, Pos{x: p.pos.x, y: p.pos.y + direction})
	}

	// forward twice
	if p.firstMove && board[p.pos.y + 2*direction][p.pos.x] == nil && board[p.pos.y + direction][p.pos.x] == nil {
		validMoves.moves = append(validMoves.moves, Pos{x: p.pos.x, y: p.pos.y + 2*direction})
	}

	if p.pos.x - 1 >= 0 {
		// capture left
		if board[p.pos.y + direction][p.pos.x - 1] != nil && board[p.pos.y + direction][p.pos.x - 1].GetColor() != p.color {
			validMoves.moves = append(validMoves.moves, Pos{x: p.pos.x - 1, y: p.pos.y + direction})
		}
		// en passant left
		if pawn, ok := board[p.pos.y][p.pos.x - 1].(*Pawn); ok && pawn.canBeEnPassanted && pawn.color != p.color {
			validMoves.moves = append(validMoves.moves, Pos{x: p.pos.x - 1, y: p.pos.y + direction})
		}
	}

	if p.pos.x + 1 < 8 {
		// capture right
		if board[p.pos.y + direction][p.pos.x + 1] != nil && board[p.pos.y + direction][p.pos.x + 1].GetColor() != p.color {
			validMoves.moves = append(validMoves.moves, Pos{x: p.pos.x + 1, y: p.pos.y + direction})
		}
		// en passant right
		if pawn, ok := board[p.pos.y][p.pos.x + 1].(*Pawn); ok && pawn.canBeEnPassanted && pawn.color != p.color {
			validMoves.moves = append(validMoves.moves, Pos{x: p.pos.x + 1, y: p.pos.y + direction})
		}
	}

	return validMoves
}

func (p *Piece) GetValidMoves(board [][]Pieces, player Color) ValidMoves {
	validMoves := ValidMoves{ p: p, moves: []Pos{} }
	if p.pieceType == ROOK || p.pieceType == QUEEN {
		for i := p.pos.x + 1; i < 8; i++ {
			if board[p.pos.y][i] == nil {
				validMoves.moves = append(validMoves.moves, Pos{x: i, y: p.pos.y})
			} else if board[p.pos.y][i].GetColor() != p.color {
				validMoves.moves = append(validMoves.moves, Pos{x: i, y: p.pos.y})
				break
			} else {
				break
			}
		}
		for i := p.pos.x - 1; i >= 0; i-- {
			if board[p.pos.y][i] == nil {
				validMoves.moves = append(validMoves.moves, Pos{x: i, y: p.pos.y})
			} else if board[p.pos.y][i].GetColor() != p.color {
				validMoves.moves = append(validMoves.moves, Pos{x: i, y: p.pos.y})
				break
			} else {
				break
			}
		}
		for i := p.pos.y + 1; i < 8; i++ {
			if board[i][p.pos.x] == nil {
				validMoves.moves = append(validMoves.moves, Pos{x: p.pos.x, y: i})
			} else if board[i][p.pos.x].GetColor() != p.color {
				validMoves.moves = append(validMoves.moves, Pos{x: p.pos.x, y: i})
				break
			} else {
				break
			}
		}
		for i := p.pos.y - 1; i >= 0; i-- {
			if board[i][p.pos.x] == nil {
				validMoves.moves = append(validMoves.moves, Pos{x: p.pos.x, y: i})
			} else if board[i][p.pos.x].GetColor() != p.color {
				validMoves.moves = append(validMoves.moves, Pos{x: p.pos.x, y: i})
				break
			} else {
				break
			}
		}
	}
	if p.pieceType == BISHOP || p.pieceType == QUEEN {
		for i, j := p.pos.x + 1, p.pos.y + 1; i < 8 && j < 8; i, j = i + 1, j + 1 {
			if board[j][i] == nil {
				validMoves.moves = append(validMoves.moves, Pos{x: i, y: j})
			} else if board[j][i].GetColor() != p.color {
				validMoves.moves = append(validMoves.moves, Pos{x: i, y: j})
				break
			} else {
				break
			}
		}
		for i, j := p.pos.x - 1, p.pos.y + 1; i >= 0 && j < 8; i, j = i - 1, j + 1 {
			if board[j][i] == nil {
				validMoves.moves = append(validMoves.moves, Pos{x: i, y: j})
			} else if board[j][i].GetColor() != p.color {
				validMoves.moves = append(validMoves.moves, Pos{x: i, y: j})
				break
			} else {
				break
			}
		}
		for i, j := p.pos.x + 1, p.pos.y - 1; i < 8 && j >= 0; i, j = i + 1, j - 1 {
			if board[j][i] == nil {
				validMoves.moves = append(validMoves.moves, Pos{x: i, y: j})
			} else if board[j][i].GetColor() != p.color {
				validMoves.moves = append(validMoves.moves, Pos{x: i, y: j})
				break
			} else {
				break
			}
		}
		for i, j := p.pos.x - 1, p.pos.y - 1; i >= 0 && j >= 0; i, j = i - 1, j - 1 {
			if board[j][i] == nil {
				validMoves.moves = append(validMoves.moves, Pos{x: i, y: j})
			} else if board[j][i].GetColor() != p.color {
				validMoves.moves = append(validMoves.moves, Pos{x: i, y: j})
				break
			} else {
				break
			}
		}
	}
	if p.pieceType == KNIGHT {
		knightMoves := []Pos{
			{x: p.pos.x + 2, y: p.pos.y + 1},
			{x: p.pos.x + 2, y: p.pos.y - 1},
			{x: p.pos.x - 2, y: p.pos.y + 1},
			{x: p.pos.x - 2, y: p.pos.y - 1},
			{x: p.pos.x + 1, y: p.pos.y + 2},
			{x: p.pos.x + 1, y: p.pos.y - 2},
			{x: p.pos.x - 1, y: p.pos.y + 2},
			{x: p.pos.x - 1, y: p.pos.y - 2},
		}
		for _, move := range knightMoves {
			if move.x >= 0 && move.x < 8 && move.y >= 0 && move.y < 8 {
				if board[move.y][move.x] == nil {
					validMoves.moves = append(validMoves.moves, move)
				} else if board[move.y][move.x].GetColor() != p.color {
					validMoves.moves = append(validMoves.moves, move)
				}
			}
		}
	}
	if p.pieceType == KING {
		kingMoves := []Pos{
			{x: p.pos.x + 1, y: p.pos.y},
			{x: p.pos.x - 1, y: p.pos.y},
			{x: p.pos.x, y: p.pos.y + 1},
			{x: p.pos.x, y: p.pos.y - 1},
			{x: p.pos.x + 1, y: p.pos.y + 1},
			{x: p.pos.x + 1, y: p.pos.y - 1},
			{x: p.pos.x - 1, y: p.pos.y + 1},
			{x: p.pos.x - 1, y: p.pos.y - 1},
		}
		// for now king can walk into check
		for _, move := range kingMoves {
			if move.x >= 0 && move.x < 8 && move.y >= 0 && move.y < 8 {
				if board[move.y][move.x] == nil {
					validMoves.moves = append(validMoves.moves, move)
				} else if board[move.y][move.x].GetColor() != p.color {
					validMoves.moves = append(validMoves.moves, move)
				}
			}
		}
	}
	return validMoves
}

func genBoard() [][]Pieces {
	board := make([][]Pieces, 8)
	for i := range board {
		board[i] = make([]Pieces, 8)
	}

	for i := 0; i < 8; i++ {
		board[1][i] = &Pawn{pos: Pos{x: i, y: 1}, color: WHITE, firstMove: true, canBeEnPassanted: false}
		board[6][i] = &Pawn{pos: Pos{x: i, y: 6}, color: BLACK, firstMove: true, canBeEnPassanted: false}
	}

	pieceTypes := []PieceType{ROOK, KNIGHT, BISHOP, QUEEN, KING, BISHOP, KNIGHT, ROOK}
	for i, pt := range pieceTypes {
		board[0][i] = &Piece{pos: Pos{x: i, y: 0}, pieceType: pt, color: WHITE, moved: false}
		board[7][i] = &Piece{pos: Pos{x: i, y: 7}, pieceType: pt, color: BLACK, moved: false}
	}

	return board
}

// this function will return the list of valid moves for a given player
func getValidMoves(board [][]Pieces, player Color) []ValidMoves {
	validMoves := []ValidMoves{}

	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			if board[i][j] == nil {
				continue
			}
			
			if piece, ok := board[i][j].(*Piece); ok && piece.color == player {
				pieceMoves := piece.GetValidMoves(board, player)
				if len(pieceMoves.moves) > 0 {
					validMoves = append(validMoves, pieceMoves)
				}
			} else if pawn, ok := board[i][j].(*Pawn); ok && pawn.color == player {
				pawnMoves := pawn.GetValidMoves(board, player)
				if len(pawnMoves.moves) > 0 {
					validMoves = append(validMoves, pawnMoves)
				}
			}
		}
	}

	return validMoves
}

func getPieceName(pieceType PieceType) string {
	switch pieceType {
	case ROOK:
		return "Rook"
	case KNIGHT:
		return "Knight"
	case BISHOP:
		return "Bishop"
	case QUEEN:
		return "Queen"
	case KING:
		return "King"
	default:
		return "Unknown"
	}
}

func main() {
	board := genBoard()
	moves := getValidMoves(board, WHITE)
	for _, move := range moves {
		if piece, ok := move.p.(*Piece); ok {
			fmt.Println("Valid moves for", getPieceName(piece.pieceType), "at position:", piece.pos.x, piece.pos.y)
		} else if pawn, ok := move.p.(*Pawn); ok {
			fmt.Println("Valid moves for pawn at position:", pawn.pos.x, pawn.pos.y)
		}
		for _, m := range move.moves {
			fmt.Println("  Move to:", m.x, m.y)
		}
	}
}