// gogot the 2dn
// plays with 3 move deep minimax with alpha beta prooning
// takes too much time even with bitmaps, perhaps implementing it beforehand was a waste of time

package bot

import (
	. "chess/game"
	"fmt"
	"math"
	"sync"
)

type GoBot struct {
	Color PColor
}

type moveResult struct {
	move  [2]uint64
	score float64
}

type Cache struct {
	sync.RWMutex
	entries map[string]float64
}

var memo = &Cache{entries: make(map[string]float64)}

func clonePiece(p Piece) Piece {
	switch v := p.(type) {
	case *Pawn:
		c := *v
		return &c
	case *Rook:
		c := *v
		return &c
	case *Knight:
		c := *v
		return &c
	case *Bishop:
		c := *v
		return &c
	case *Queen:
		c := *v
		return &c
	case *King:
		c := *v
		return &c
	default:
		return nil
	}
}

func cloneGame(g *Game) *Game {
	newG := &Game{
		Board:           g.Board,
		CurrentTurn:     g.CurrentTurn,
		ValidWhiteMoves: make(map[uint64]uint64),
		ValidBlackMoves: make(map[uint64]uint64),
	}

	for k, v := range g.ValidWhiteMoves {
		newG.ValidWhiteMoves[k] = v
	}
	for k, v := range g.ValidBlackMoves {
		newG.ValidBlackMoves[k] = v
	}

	newG.WhitePieces = make([]Piece, len(g.WhitePieces))
	for i, p := range g.WhitePieces {
		newG.WhitePieces[i] = clonePiece(p)
	}
	newG.BlackPieces = make([]Piece, len(g.BlackPieces))
	for i, p := range g.BlackPieces {
		newG.BlackPieces[i] = clonePiece(p)
	}

	var whiteRooks, blackRooks []*Rook
	for _, p := range newG.WhitePieces {
		if k, ok := p.(*King); ok {
			newG.WhiteKing = k
		}
		if r, ok := p.(*Rook); ok {
			whiteRooks = append(whiteRooks, r)
		}
	}
	for _, p := range newG.BlackPieces {
		if k, ok := p.(*King); ok {
			newG.BlackKing = k
		}
		if r, ok := p.(*Rook); ok {
			blackRooks = append(blackRooks, r)
		}
	}

	if len(whiteRooks) >= 1 {
		if g.WhiteRooks.Left != nil && whiteRooks[0].GetPos() == g.WhiteRooks.Left.GetPos() {
			newG.WhiteRooks.Left = whiteRooks[0]
			if len(whiteRooks) > 1 {
				newG.WhiteRooks.Right = whiteRooks[1]
			}
		} else {
			newG.WhiteRooks.Right = whiteRooks[0]
			if len(whiteRooks) > 1 {
				newG.WhiteRooks.Left = whiteRooks[1]
			}
		}
	}

	if len(blackRooks) >= 1 {
		if g.BlackRooks.Left != nil && blackRooks[0].GetPos() == g.BlackRooks.Left.GetPos() {
			newG.BlackRooks.Left = blackRooks[0]
			if len(blackRooks) > 1 {
				newG.BlackRooks.Right = blackRooks[1]
			}
		} else {
			newG.BlackRooks.Right = blackRooks[0]
			if len(blackRooks) > 1 {
				newG.BlackRooks.Left = blackRooks[1]
			}
		}
	}

	return newG
}

func (b *GoBot) Play(g *Game, myColor PColor) error {
	validMoves := g.GetValidMoves(myColor)

	if len(validMoves) == 0 {
		return fmt.Errorf("No valid moves available for %v", myColor)
	}

	var bestMove [2]uint64
	var bestScore float64

	if myColor == WHITE {
		bestScore = math.Inf(-1)
	} else {
		bestScore = math.Inf(1)
	}

	var wg sync.WaitGroup
	results := make(chan moveResult)

	for piece, moves := range validMoves {
		if moves == 0 {
			continue
		}
		for _, move := range GetBits(moves) {
			wg.Add(1)
			go func(piece, move uint64) {
				defer wg.Done()
				newPos := uint64(1) << move

				gameCopy := cloneGame(g)
				err := gameCopy.Move(piece, newPos, myColor)
				if err != nil {
					return
				}

				isMaximizing := myColor == WHITE
				score := minimax(gameCopy, 3, !isMaximizing, math.Inf(-1), math.Inf(1))
				results <- moveResult{move: [2]uint64{piece, newPos}, score: score}
			}(piece, uint64(move))
		}
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	foundFirstMove := false
	for result := range results {
		if !foundFirstMove {
			bestScore = result.score
			bestMove = result.move
			foundFirstMove = true
		} else {
			isMaximizing := myColor == WHITE
			if isMaximizing {
				if result.score > bestScore {
					bestScore = result.score
					bestMove = result.move
				}
			} else {
				if result.score < bestScore {
					bestScore = result.score
					bestMove = result.move
				}
			}
		}
	}

	if !foundFirstMove {
		return fmt.Errorf("No valid moves could be found or evaluated for %v", myColor)
	}

	return g.Move(bestMove[0], bestMove[1], myColor)
}

func evaluate(g *Game) float64 {
	score := 0.0
	score += float64(len(GetBits(g.Board.WhitePawns))-len(GetBits(g.Board.BlackPawns))) * 1.0
	score += float64(len(GetBits(g.Board.WhiteKnights))-len(GetBits(g.Board.BlackKnights))) * 3.0
	score += float64(len(GetBits(g.Board.WhiteBishops))-len(GetBits(g.Board.BlackBishops))) * 3.0
	score += float64(len(GetBits(g.Board.WhiteRooks))-len(GetBits(g.Board.BlackRooks))) * 5.0
	score += float64(len(GetBits(g.Board.WhiteQueens))-len(GetBits(g.Board.BlackQueens))) * 9.0

	if len(g.ValidWhiteMoves) == 0 {
		score -= 1000.0
	}
	if len(g.ValidBlackMoves) == 0 {
		score += 1000.0
	}
	return score
}

func getPositionKey(g *Game) string {
	return fmt.Sprintf("%v|%v|%v|%v|%v|%v", g.Board.WhitePawns, g.Board.BlackPawns, g.Board.WhiteKnights, g.Board.BlackKnights, g.CurrentTurn, g.Board.WhiteKing^g.Board.BlackKing)
}

func minimax(g *Game, depth int, maximizingPlayer bool, alpha float64, beta float64) float64 {
	if depth == 0 {
		return evaluate(g)
	}

	key := getPositionKey(g)
	memo.RLock()
	if val, ok := memo.entries[key]; ok {
		memo.RUnlock()
		return val
	}
	memo.RUnlock()

	var moves map[uint64]uint64
	var colorToMove PColor

	if maximizingPlayer {
		colorToMove = WHITE
	} else {
		colorToMove = BLACK
	}
	moves = g.GetValidMoves(colorToMove)

	if len(moves) == 0 {
		return evaluate(g)
	}

	var result float64
	if maximizingPlayer {
		maxEval := math.Inf(-1)
		for piece, pieceMoves := range moves {
			for _, move := range GetBits(pieceMoves) {
				newPos := uint64(1) << move
				gameCopy := cloneGame(g)
				err := gameCopy.Move(piece, newPos, colorToMove)
				if err != nil {
					continue
				}
				eval := minimax(gameCopy, depth-1, false, alpha, beta)
				maxEval = math.Max(maxEval, eval)
				alpha = math.Max(alpha, eval)
				if beta <= alpha {
					break
				}
			}
			if beta <= alpha {
				break
			}
		}
		result = maxEval
	} else {
		minEval := math.Inf(1)
		for piece, pieceMoves := range moves {
			for _, move := range GetBits(pieceMoves) {
				newPos := uint64(1) << move
				gameCopy := cloneGame(g)
				err := gameCopy.Move(piece, newPos, colorToMove)
				if err != nil {
					continue
				}
				eval := minimax(gameCopy, depth-1, true, alpha, beta)
				minEval = math.Min(minEval, eval)
				beta = math.Min(beta, eval)
				if beta <= alpha {
					break
				}
			}
			if beta <= alpha {
				break
			}
		}
		result = minEval
	}

	memo.Lock()
	memo.entries[key] = result
	memo.Unlock()
	return result
}
