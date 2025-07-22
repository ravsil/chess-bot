#include <stdio.h>

#define BYTE_TO_BINARY_PATTERN "%c%c%c%c%c%c%c%c"
#define bool int
#define true 1
#define false 0
#define BYTE_TO_BINARY(byte)  \
  ((byte) & 0x80 ? '1' : '0'), \
  ((byte) & 0x40 ? '1' : '0'), \
  ((byte) & 0x20 ? '1' : '0'), \
  ((byte) & 0x10 ? '1' : '0'), \
  ((byte) & 0x08 ? '1' : '0'), \
  ((byte) & 0x04 ? '1' : '0'), \
  ((byte) & 0x02 ? '1' : '0'), \
  ((byte) & 0x01 ? '1' : '0') 
#define U64 unsigned long long int


typedef enum eenumSquare {
  A1, B1, C1, D1, E1, F1, G1, H1,
  A2, B2, C2, D2, E2, F2, G2, H2,
  A3, B3, C3, D3, E3, F3, G3, H3,
  A4, B4, C4, D4, E4, F4, G4, H4,
  A5, B5, C5, D5, E5, F5, G5, H5,
  A6, B6, C6, D6, E6, F6, G6, H6,
  A7, B7, C7, D7, E7, F7, G7, H7,
  A8, B8, C8, D8, E8, F8, G8, H8,
  NO_SQ
} enumSquare;

typedef enum eenumDir {
  NORTH,
  NORTHEAST,
  EAST,
  SOUTHEAST,
  SOUTH,
  SOUTHWEST,
  WEST,
  NORTHWEST
} enumDir;

//////////////////////////////

U64 rayAttacks[8][64];

#define set_bit(bitboard, square) ((bitboard) |= (1ULL << (square)))

void initRayAttacks() {
  for (int sq = 0; sq < 64; sq++) {
    // North
    U64 north_attack = 0ULL;
    for (int i = sq + 8; i < 64; i += 8) set_bit(north_attack, i);
    rayAttacks[NORTH][sq] = north_attack;

    // South
    U64 south_attack = 0ULL;
    for (int i = sq - 8; i >= 0; i -= 8) set_bit(south_attack, i);
    rayAttacks[SOUTH][sq] = south_attack;

    // East
    U64 east_attack = 0ULL;
    for (int i = sq + 1; (i % 8) != 0; i++) set_bit(east_attack, i);
    rayAttacks[EAST][sq] = east_attack;

    // West
    U64 west_attack = 0ULL;
    for (int i = sq - 1; (i % 8) != 7 && i >= 0; i--) set_bit(west_attack, i);
    rayAttacks[WEST][sq] = west_attack;

    // Northeast
    U64 ne_attack = 0ULL;
    for (int i = sq + 9; i < 64 && (i % 8) != 0; i += 9) set_bit(ne_attack, i);
    rayAttacks[NORTHEAST][sq] = ne_attack;

    // Northwest
    U64 nw_attack = 0ULL;
    for (int i = sq + 7; i < 64 && (i % 8) != 7; i += 7) set_bit(nw_attack, i);
    rayAttacks[NORTHWEST][sq] = nw_attack;

    // Southeast
    U64 se_attack = 0ULL;
    for (int i = sq - 7; i >= 0 && (i % 8) != 0; i -= 7) set_bit(se_attack, i);
    rayAttacks[SOUTHEAST][sq] = se_attack;

    // Southwest
    U64 sw_attack = 0ULL;
    for (int i = sq - 9; i >= 0 && (i % 8) != 7; i -= 9) set_bit(sw_attack, i);
    rayAttacks[SOUTHWEST][sq] = sw_attack;
  }
}

///////////////////////////////

int ms1bTable[256];

void init_ms1bTable() {
  int i;
  for (i = 0; i < 256; i++) {
    ms1bTable[i] = (
      (i > 127) ? 7 :
      (i > 63) ? 6 :
      (i > 31) ? 5 :
      (i > 15) ? 4 :
      (i > 7) ? 3 :
      (i > 3) ? 2 :
      (i > 1) ? 1 :
      0
      );
  }
}

int bitScanReverse(U64 bb)
{
  int result = 0;
  if (bb > 0xFFFFFFFF) {
    bb >>= 32;
    result = 32;
  }
  if (bb > 0xFFFF) {
    bb >>= 16;
    result += 16;
  }
  if (bb > 0xFF) {
    bb >>= 8;
    result += 8;
  }
  return result + ms1bTable[bb];
}

int bitScan(U64 bb, bool reverse) {
  U64 rMask;
  //assert (bb != 0);
  rMask = -(U64)reverse;
  bb &= -bb | rMask;
  return bitScanReverse(bb);
}

///////////////////////////////

int bitScanForward(U64 bb) {
  bb &= -bb;
  return bitScanReverse(bb);
}

///////////////////////////////

U64 getPositiveRayAttacks(U64 occupied, enumDir dir8, enumSquare square) {
  U64 attacks = rayAttacks[dir8][square];
  U64 blocker = attacks & occupied;
  if (blocker) {
    square = bitScanForward(blocker);
    attacks ^= rayAttacks[dir8][square];
  }
  return attacks;
}

U64 getNegativeRayAttacks(U64 occupied, enumDir dir8, enumSquare square) {
  U64 attacks = rayAttacks[dir8][square];
  U64 blocker = attacks & occupied;
  if (blocker) {
    square = bitScanReverse(blocker);
    attacks ^= rayAttacks[dir8][square];
  }
  return attacks;
}

U64 diagonalAttacks(U64 occ, enumSquare sq) {
  return getPositiveRayAttacks(occ, NORTHEAST, sq)
    | getNegativeRayAttacks(occ, SOUTHEAST, sq); // ^ +
}

U64 antiDiagAttacks(U64 occ, enumSquare sq) {
  return getPositiveRayAttacks(occ, NORTHWEST, sq)
    | getNegativeRayAttacks(occ, SOUTHWEST, sq); // ^ +
}

U64 fileAttacks(U64 occ, enumSquare sq) {
  return getPositiveRayAttacks(occ, NORTH, sq)
    | getNegativeRayAttacks(occ, SOUTH, sq); // ^ +
}

U64 rankAttacks(U64 occ, enumSquare sq) {
  return getPositiveRayAttacks(occ, EAST, sq)
    | getNegativeRayAttacks(occ, WEST, sq); // ^ +
}

int main()
{
  init_ms1bTable();
  initRayAttacks();
  U64 m = 0b1111111100000011111000000000000000000111111110000000000000000000; // test board setup
  U64 n = diagonalAttacks(m, D4); // attacking squares
  printf(""BYTE_TO_BINARY_PATTERN"\n"BYTE_TO_BINARY_PATTERN"\n"BYTE_TO_BINARY_PATTERN"\n"BYTE_TO_BINARY_PATTERN"\n"BYTE_TO_BINARY_PATTERN"\n"BYTE_TO_BINARY_PATTERN"\n"BYTE_TO_BINARY_PATTERN"\n"BYTE_TO_BINARY_PATTERN"\n",
    BYTE_TO_BINARY(m >> 56), BYTE_TO_BINARY(m >> 48), BYTE_TO_BINARY(m >> 40), BYTE_TO_BINARY(m >> 32), BYTE_TO_BINARY(m >> 24), BYTE_TO_BINARY(m >> 16), BYTE_TO_BINARY(m >> 8), BYTE_TO_BINARY(m));

  printf("\n\n");
  printf(""BYTE_TO_BINARY_PATTERN"\n"BYTE_TO_BINARY_PATTERN"\n"BYTE_TO_BINARY_PATTERN"\n"BYTE_TO_BINARY_PATTERN"\n"BYTE_TO_BINARY_PATTERN"\n"BYTE_TO_BINARY_PATTERN"\n"BYTE_TO_BINARY_PATTERN"\n"BYTE_TO_BINARY_PATTERN"\n",
    BYTE_TO_BINARY(n >> 56), BYTE_TO_BINARY(n >> 48), BYTE_TO_BINARY(n >> 40), BYTE_TO_BINARY(n >> 32), BYTE_TO_BINARY(n >> 24), BYTE_TO_BINARY(n >> 16), BYTE_TO_BINARY(n >> 8), BYTE_TO_BINARY(n));
  return 0;
}