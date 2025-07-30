// gobot the 1st
// plays with random moves

package bot

import (
	. "chess/game"
	"fmt"
	"math/rand"
)

type RandomBot struct {
	Color PColor
}

func (b *RandomBot) Play(game *Game, myColor PColor) error {
	validMoves := game.GetValidMoves(myColor)

	if len(validMoves) == 0 {
		return fmt.Errorf("No valid moves available for %s", myColor)
	}

	var pieces []uint64
	for piece := range validMoves {
		if validMoves[piece] != 0 {
			pieces = append(pieces, piece)
		}
	}

	if len(pieces) == 0 {
		return fmt.Errorf("No pieces available for %s", myColor)
	}

	p := pieces[rand.Intn(len(pieces))]

	moves := GetBits(validMoves[p])
	m := moves[rand.Intn(len(moves))]
	newPos := uint64(1) << m

	game.Move(p, newPos, myColor)
	return nil
}
