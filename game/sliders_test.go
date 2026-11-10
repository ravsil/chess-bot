package game

import (
	"testing"
)

func TestRayAttacksAutoInitialized(t *testing.T) {
	// Without manually calling InitRanged(), rayAttacks should be initialized.
	pos := 28 // e4
	attacks := rayAttacks[NORTH][pos]
	if attacks == 0 {
		t.Fatalf("rayAttacks[NORTH][%d] should be initialized automatically", pos)
	}

	occ := uint64(1 << 36) // e5 blocks NORTH ray from e4
	ray := GetPositiveRayAttacks(occ, NORTH, pos)
	if ray&(1<<36) == 0 {
		t.Errorf("expected e5 (bit 36) to be attacked, got %064b", ray)
	}
	if ray&(1<<44) != 0 {
		t.Errorf("expected e6 (bit 44) to be blocked, got %064b", ray)
	}
}

func TestBitScans(t *testing.T) {
	bb := uint64((1 << 12) | (1 << 45))
	fwd := BitScanForward(bb)
	rev := BitScanReverse(bb)
	if fwd != 12 {
		t.Errorf("BitScanForward expected 12, got %d", fwd)
	}
	if rev != 45 {
		t.Errorf("BitScanReverse expected 45, got %d", rev)
	}
}
