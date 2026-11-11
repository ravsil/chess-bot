package tests

import (
	"bytes"
	. "chess/bot"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRoutesGameplay(t *testing.T) {
	// 1. GET /getBoard
	req := httptest.NewRequest(http.MethodGet, "/getBoard", nil)
	rec := httptest.NewRecorder()
	GetBoard(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("GET /getBoard status = %d, want 200", rec.Code)
	}

	// 2. GET /getMoves
	req = httptest.NewRequest(http.MethodGet, "/getMoves", nil)
	rec = httptest.NewRecorder()
	GetMoves(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("GET /getMoves status = %d, want 200", rec.Code)
	}

	// 3. POST /movePiece (White plays e2 -> e4: 1<<12 -> 1<<28)
	payload, _ := json.Marshal(map[string]uint64{
		"piece": 1 << 12,
		"move":  1 << 28,
	})
	req = httptest.NewRequest(http.MethodPost, "/movePiece", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	MovePiece(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("POST /movePiece valid move status = %d, want 200: %s", rec.Code, rec.Body.String())
	}

	// Verify content-type is application/json
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("Content-Type = %q, want application/json", ct)
	}

	// Verify response contains both PlayerMove/PlayerBoard and BotMove
	var moveResp struct {
		Board       json.RawMessage `json:"Board"`
		PlayerBoard struct {
			WhitePawns uint64 `json:"WhitePawns"`
		} `json:"PlayerBoard"`
		PlayerMove [2]uint64 `json:"PlayerMove"`
		BotMove    [2]uint64 `json:"BotMove"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &moveResp); err != nil {
		t.Fatalf("failed to unmarshal MoveResponse: %v", err)
	}
	if moveResp.PlayerMove != [2]uint64{1 << 12, 1 << 28} {
		t.Errorf("PlayerMove = %v, want [%d %d]", moveResp.PlayerMove, uint64(1<<12), uint64(1<<28))
	}
	if moveResp.PlayerBoard.WhitePawns&(1<<28) == 0 {
		t.Errorf("PlayerBoard should have white pawn on e4 (1<<28)")
	}
	if moveResp.BotMove[0] == 0 || moveResp.BotMove[1] == 0 {
		t.Errorf("expected BotMove to be populated, got %v", moveResp.BotMove)
	}

	// Verify invalid move returns 400
	badPayload, _ := json.Marshal(map[string]uint64{
		"piece": 1 << 12,
		"move":  1 << 36,
	})
	req = httptest.NewRequest(http.MethodPost, "/movePiece", bytes.NewReader(badPayload))
	rec = httptest.NewRecorder()
	MovePiece(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("POST /movePiece invalid move status = %d, want 400", rec.Code)
	}
}
