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

func (b *RandomBot) Play(game *Game, myColor PColor) {
	// Get the valid moves for the current player
	if myColor == WHITE {
		fmt.Println("Bot is playing for White")
	} else {
		fmt.Println("Bot is playing for Black")
	}
	validMoves := game.GetValidMoves(myColor)
	fmt.Println("Valid moves for bot:", validMoves)

	if len(validMoves) == 0 {
		return
	}

	var pieces []uint64
	for piece := range validMoves {
		if validMoves[piece] != 0 {
			pieces = append(pieces, piece)
		}
	}

	if len(pieces) == 0 {
		return
	}

	p := pieces[rand.Intn(len(pieces))]

	moves := GetBits(validMoves[p])
	m := moves[rand.Intn(len(moves))]
	newPos := uint64(1) << m

	game.Move(p, newPos, myColor)
	fmt.Println("Bot played:", p, "to", newPos)
}
