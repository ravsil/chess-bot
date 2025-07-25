// Complicated Boolean Algebra to calculate ranged pieces attacking squares
// Here I gave up on trying to do this myself and read the chess programming wiki
// https://www.chessprogramming.org/Classical_Approach

package game

import "fmt"

type Direction int

const (
	NORTH Direction = iota
	NORTHEAST
	EAST
	SOUTHEAST
	SOUTH
	SOUTHWEST
	WEST
	NORTHWEST
)
const LEFT_BORDER uint64 = 0b00000001_00000001_00000001_00000001_00000001_00000001_00000001_00000001
const RIGHT_BORDER uint64 = 0b10000000_10000000_10000000_10000000_10000000_10000000_10000000_10000000
const DOUBLE_LEFT_BORDER uint64 = (LEFT_BORDER | (LEFT_BORDER << 1))
const DOUBLE_RIGHT_BORDER uint64 = (RIGHT_BORDER | (RIGHT_BORDER >> 1))

var rayAttacks [8][64]uint64
var ms1bTable [256]int

func SetBit(bitboard *uint64, square int) {
	*bitboard |= (1 << square)
}

func GetSingleBit(bb uint64) (int, error) {
	for i := 0; i < 64; i++ {
		if (bb>>i)&1 != 0 {
			return i, nil
		}
	}
	return -1, fmt.Errorf("no bit set in bitboard")
}

func GetBits(bb uint64) []int {
	var bits []int
	for i := 0; i < 64; i++ {
		if (bb>>i)&1 != 0 {
			bits = append(bits, i)
		}
	}
	return bits
}

func InitRayAttacks() {
	for sq := 0; sq < 64; sq++ {

		var northAttack uint64 = 0
		for i := sq + 8; i < 64; i += 8 {
			SetBit(&northAttack, i)
		}
		rayAttacks[NORTH][sq] = northAttack

		var southAttack uint64 = 0
		for i := sq - 8; i >= 0; i -= 8 {
			SetBit(&southAttack, i)
		}
		rayAttacks[SOUTH][sq] = southAttack

		var eastAttack uint64 = 0
		for i := sq + 1; (i % 8) != 0; i++ {
			SetBit(&eastAttack, i)
		}
		rayAttacks[EAST][sq] = eastAttack

		var westAttack uint64 = 0
		for i := sq - 1; (i%8) != 7 && i >= 0; i-- {
			SetBit(&westAttack, i)
		}
		rayAttacks[WEST][sq] = westAttack

		var neAttack uint64 = 0
		for i := sq + 9; i < 64 && (i%8) != 0; i += 9 {
			SetBit(&neAttack, i)
		}
		rayAttacks[NORTHEAST][sq] = neAttack

		var nwAttack uint64 = 0
		for i := sq + 7; i < 64 && (i%8) != 7; i += 7 {
			SetBit(&nwAttack, i)
		}
		rayAttacks[NORTHWEST][sq] = nwAttack

		var seAttack uint64 = 0
		for i := sq - 7; i >= 0 && (i%8) != 0; i -= 7 {
			SetBit(&seAttack, i)
		}
		rayAttacks[SOUTHEAST][sq] = seAttack

		var swAttack uint64 = 0
		for i := sq - 9; i >= 0 && (i%8) != 7; i -= 9 {
			SetBit(&swAttack, i)
		}
		rayAttacks[SOUTHWEST][sq] = swAttack
	}
}

func InitMs1bTable() {
	for i := 0; i < 256; i++ {
		if i > 127 {
			ms1bTable[i] = 7
		} else if i > 63 {
			ms1bTable[i] = 6
		} else if i > 31 {
			ms1bTable[i] = 5
		} else if i > 15 {
			ms1bTable[i] = 4
		} else if i > 7 {
			ms1bTable[i] = 3
		} else if i > 3 {
			ms1bTable[i] = 2
		} else if i > 1 {
			ms1bTable[i] = 1
		} else {
			ms1bTable[i] = 0
		}
	}
}

func flipVertical(x uint64) uint64 {
	return (x << 56) |
		((x << 40) & uint64(0x00ff000000000000)) |
		((x << 24) & uint64(0x0000ff0000000000)) |
		((x << 8) & uint64(0x000000ff00000000)) |
		((x >> 8) & uint64(0x00000000ff000000)) |
		((x >> 24) & uint64(0x0000000000ff0000)) |
		((x >> 40) & uint64(0x000000000000ff00)) |
		(x >> 56)
}
func flipHorizontal(x uint64) uint64 {
	return ((x & 0x0101010101010101) << 7) |
		((x & 0x0202020202020202) << 5) |
		((x & 0x0404040404040404) << 3) |
		((x & 0x0808080808080808) << 1) |
		((x & 0x1010101010101010) >> 1) |
		((x & 0x2020202020202020) >> 3) |
		((x & 0x4040404040404040) >> 5) |
		((x & 0x8080808080808080) >> 7)
}

func Rotate180(x uint64) uint64 {
	return flipHorizontal(flipVertical(x))
}

func BitScanReverse(bb uint64) int {
	result := 0
	if bb > 0xFFFFFFFF {
		bb >>= 32
		result = 32
	}
	if bb > 0xFFFF {
		bb >>= 16
		result += 16
	}
	if bb > 0xFF {
		bb >>= 8
		result += 8
	}
	return result + ms1bTable[bb]
}

func BitScan(bb uint64, reverse bool) int {
	var rMask uint64
	if reverse {
		rMask = ^uint64(0)
	} else {
		rMask = 0
	}
	bb &= (^bb + 1) | rMask
	return BitScanReverse(bb)
}

func BitScanForward(bb uint64) int {
	bb &= ^bb + 1
	return BitScanReverse(bb)
}

func GetPositiveRayAttacks(occupied uint64, dir Direction, square int) uint64 {
	attacks := rayAttacks[dir][square]
	blocker := attacks & occupied
	if blocker != 0 {
		square = BitScanForward(blocker)
		attacks ^= rayAttacks[dir][square]
	}
	return attacks
}

func GetNegativeRayAttacks(occupied uint64, dir Direction, square int) uint64 {
	attacks := rayAttacks[dir][square]
	blocker := attacks & occupied
	if blocker != 0 {
		square = BitScanReverse(blocker)
		attacks ^= rayAttacks[dir][square]
	}
	return attacks
}

func BishopAttacks(occ uint64, sq int) uint64 {
	return GetPositiveRayAttacks(occ, NORTHEAST, sq) |
		GetNegativeRayAttacks(occ, SOUTHEAST, sq) | GetPositiveRayAttacks(occ, NORTHWEST, sq) |
		GetNegativeRayAttacks(occ, SOUTHWEST, sq)
}

func RookAttacks(occ uint64, sq int) uint64 {
	return GetPositiveRayAttacks(occ, NORTH, sq) |
		GetNegativeRayAttacks(occ, SOUTH, sq) | GetPositiveRayAttacks(occ, EAST, sq) |
		GetNegativeRayAttacks(occ, WEST, sq)
}

func InitRanged() {
	InitMs1bTable()
	InitRayAttacks()
}

func printBits(occ uint64) {
	for i := 0; i < 64; i++ {
		if (occ>>i)&1 != 0 {
			fmt.Print("1")
		} else {
			fmt.Print("0")
		}
		if (i+1)%8 == 0 {
			fmt.Println()
		}
	}
	fmt.Println()
}
