package tests

import (
	. "chess/game"
	"testing"
)

func TestBlackPawnInfluenceNoWrap(t *testing.T) {
	var b Board
	// Black pawn on a7 (sq 48, col 0). Southwest attack (48-9=39, col 7) would wrap without ^RIGHT_BORDER.
	b.BlackPawns = 1 << 48

	inf := b.GetBlackInfluence()
	// Square 39 is h5 (1<<39), which should NOT be influenced by a7 pawn!
	if inf&(1<<39) != 0 {
		t.Fatalf("Black pawn on a7 influence wrapped around to h5 (bit 39), influence = %064b", inf)
	}

	// Black pawn on h7 (sq 55, col 7). Southeast attack (55-7=48, col 0) would wrap without ^LEFT_BORDER.
	b.BlackPawns = 1 << 55
	inf = b.GetBlackInfluence()
	// Square 48 is a6 (1<<48), which should NOT be influenced by h7 pawn!
	if inf&(1<<48) != 0 {
		t.Fatalf("Black pawn on h7 influence wrapped around to a6 (bit 48), influence = %064b", inf)
	}
}
