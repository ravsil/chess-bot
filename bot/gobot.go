// gobot the 2dn
// plays with 3 move deep minimax with alpha beta pruning

package bot

import (
	. "chess/game"
	"fmt"
	"math"
	"math/bits"
	"sync"
)

type GoBot struct {
	Color PColor
}

type moveResult struct {
	move  [2]uint64
	score float64
}

const (
	ttExact = iota
	ttLower
	ttUpper
)

type ttEntry struct {
	score float64
	depth int
	flag  int
}

type PositionKey struct {
	board       Board
	turn        PColor
	castleFlags uint8
}

type Cache struct {
	sync.RWMutex
	entries map[PositionKey]ttEntry
}

var memo = &Cache{entries: make(map[PositionKey]ttEntry)}

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
		ValidWhiteMoves: make(map[uint64]uint64, len(g.ValidWhiteMoves)),
		ValidBlackMoves: make(map[uint64]uint64, len(g.ValidBlackMoves)),
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

	for _, r := range whiteRooks {
		if g.WhiteRooks.Left != nil && r.GetPos() == g.WhiteRooks.Left.GetPos() {
			newG.WhiteRooks.Left = r
		}
		if g.WhiteRooks.Right != nil && r.GetPos() == g.WhiteRooks.Right.GetPos() {
			newG.WhiteRooks.Right = r
		}
	}

	for _, r := range blackRooks {
		if g.BlackRooks.Left != nil && r.GetPos() == g.BlackRooks.Left.GetPos() {
			newG.BlackRooks.Left = r
		}
		if g.BlackRooks.Right != nil && r.GetPos() == g.BlackRooks.Right.GetPos() {
			newG.BlackRooks.Right = r
		}
	}

	return newG
}

func (b *GoBot) Play(g *Game, myColor PColor) error {
	validMoves := g.GetValidMoves(myColor)

	if len(validMoves) == 0 {
		return fmt.Errorf("No valid moves available for %v", myColor)
	}

	memo.Lock()
	if len(memo.entries) > 500000 {
		memo.entries = make(map[PositionKey]ttEntry)
	}
	memo.Unlock()

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

func Evaluate(g *Game, depth int) float64 {
	if g.CurrentTurn == WHITE && len(g.ValidWhiteMoves) == 0 {
		if g.IsInCheck(WHITE) {
			return -10000.0 - float64(depth)
		}
		return 0.0
	}
	if g.CurrentTurn == BLACK && len(g.ValidBlackMoves) == 0 {
		if g.IsInCheck(BLACK) {
			return 10000.0 + float64(depth)
		}
		return 0.0
	}

	score := 0.0
	score += float64(bits.OnesCount64(g.Board.WhitePawns)-bits.OnesCount64(g.Board.BlackPawns)) * 1.0
	score += float64(bits.OnesCount64(g.Board.WhiteKnights)-bits.OnesCount64(g.Board.BlackKnights)) * 3.0
	score += float64(bits.OnesCount64(g.Board.WhiteBishops)-bits.OnesCount64(g.Board.BlackBishops)) * 3.0
	score += float64(bits.OnesCount64(g.Board.WhiteRooks)-bits.OnesCount64(g.Board.BlackRooks)) * 5.0
	score += float64(bits.OnesCount64(g.Board.WhiteQueens)-bits.OnesCount64(g.Board.BlackQueens)) * 9.0
	return score
}

func GetPositionKey(g *Game) PositionKey {
	var flags uint8
	if g.WhiteKing != nil && !g.WhiteKing.Moved {
		flags |= 1 << 0
	}
	if g.WhiteRooks.Left != nil && !g.WhiteRooks.Left.Moved {
		flags |= 1 << 1
	}
	if g.WhiteRooks.Right != nil && !g.WhiteRooks.Right.Moved {
		flags |= 1 << 2
	}
	if g.BlackKing != nil && !g.BlackKing.Moved {
		flags |= 1 << 3
	}
	if g.BlackRooks.Left != nil && !g.BlackRooks.Left.Moved {
		flags |= 1 << 4
	}
	if g.BlackRooks.Right != nil && !g.BlackRooks.Right.Moved {
		flags |= 1 << 5
	}
	return PositionKey{
		board:       g.Board,
		turn:        g.CurrentTurn,
		castleFlags: flags,
	}
}

func minimax(g *Game, depth int, maximizingPlayer bool, alpha float64, beta float64) float64 {
	var moves map[uint64]uint64
	var colorToMove PColor

	if maximizingPlayer {
		colorToMove = WHITE
		moves = g.ValidWhiteMoves
	} else {
		colorToMove = BLACK
		moves = g.ValidBlackMoves
	}
	if moves == nil {
		moves = g.GetValidMoves(colorToMove)
	}

	if len(moves) == 0 || depth == 0 {
		return Evaluate(g, depth)
	}

	key := GetPositionKey(g)
	origAlpha, origBeta := alpha, beta

	memo.RLock()
	entry, ok := memo.entries[key]
	memo.RUnlock()
	if ok && entry.depth >= depth {
		switch entry.flag {
		case ttExact:
			return entry.score
		case ttLower:
			if entry.score > alpha {
				alpha = entry.score
			}
		case ttUpper:
			if entry.score < beta {
				beta = entry.score
			}
		}
		if alpha >= beta {
			return entry.score
		}
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
				if eval > maxEval {
					maxEval = eval
				}
				if eval > alpha {
					alpha = eval
				}
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
				if eval < minEval {
					minEval = eval
				}
				if eval < beta {
					beta = eval
				}
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

	flag := ttExact
	if result <= origAlpha {
		flag = ttUpper
	} else if result >= origBeta {
		flag = ttLower
	}

	memo.Lock()
	if existing, exists := memo.entries[key]; !exists || depth >= existing.depth {
		memo.entries[key] = ttEntry{score: result, depth: depth, flag: flag}
	}
	memo.Unlock()
	return result
}
