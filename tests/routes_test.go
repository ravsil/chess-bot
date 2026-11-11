package tests

import (
	"bytes"
	. "chess/bot"
	. "chess/game"
	"encoding/json"
	"math/bits"
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
	var boardAfterInit Board
	if err := json.Unmarshal(rec.Body.Bytes(), &boardAfterInit); err != nil {
		t.Fatalf("failed to unmarshal /getBoard response: %v", err)
	}

	// 2. GET /getMoves
	req = httptest.NewRequest(http.MethodGet, "/getMoves", nil)
	rec = httptest.NewRecorder()
	GetMoves(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("GET /getMoves status = %d, want 200", rec.Code)
	}

	var movesResp struct {
		ValidMoves  map[uint64]uint64 `json:"validMoves"`
		PlayerColor PColor            `json:"playerColor"`
		StartBoard  Board             `json:"startBoard"`
		OpeningMove [2]uint64         `json:"openingMove"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &movesResp); err != nil {
		t.Fatalf("failed to unmarshal /getMoves response: %v", err)
	}

	if movesResp.PlayerColor == BLACK {
		if movesResp.OpeningMove[0] == 0 || movesResp.OpeningMove[1] == 0 {
			t.Fatalf("expected non-zero OpeningMove when player is Black, got %v", movesResp.OpeningMove)
		}
		if movesResp.StartBoard == boardAfterInit {
			t.Fatalf("expected boardAfterInit to reflect White's opening move when player is Black")
		}
	} else {
		if movesResp.OpeningMove != [2]uint64{0, 0} {
			t.Fatalf("expected zero OpeningMove when player is White, got %v", movesResp.OpeningMove)
		}
		if movesResp.StartBoard != boardAfterInit {
			t.Fatalf("expected StartBoard to match boardAfterInit when player is White")
		}
	}

	var pieceBit, moveBit uint64
	for p, targets := range movesResp.ValidMoves {
		if targets != 0 {
			pieceBit = p
			moveBit = uint64(1) << bits.TrailingZeros64(targets)
			break
		}
	}
	if pieceBit == 0 || moveBit == 0 {
		t.Fatalf("expected at least one valid move in /getMoves")
	}

	// 3. POST /movePiece with a valid move for assigned color
	payload, _ := json.Marshal(map[string]uint64{
		"piece": pieceBit,
		"move":  moveBit,
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
		PlayerBoard Board           `json:"PlayerBoard"`
		PlayerMove  [2]uint64       `json:"PlayerMove"`
		BotMove     [2]uint64       `json:"BotMove"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &moveResp); err != nil {
		t.Fatalf("failed to unmarshal MoveResponse: %v", err)
	}
	if moveResp.PlayerMove != [2]uint64{pieceBit, moveBit} {
		t.Errorf("PlayerMove = %v, want [%d %d]", moveResp.PlayerMove, pieceBit, moveBit)
	}
	if moveResp.BotMove[0] == 0 || moveResp.BotMove[1] == 0 {
		t.Errorf("expected BotMove to be populated, got %v", moveResp.BotMove)
	}

	// Verify invalid move returns 400
	badPayload, _ := json.Marshal(map[string]uint64{
		"piece": pieceBit,
		"move":  pieceBit,
	})
	req = httptest.NewRequest(http.MethodPost, "/movePiece", bytes.NewReader(badPayload))
	rec = httptest.NewRecorder()
	MovePiece(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("POST /movePiece invalid move status = %d, want 400", rec.Code)
	}
}

func TestRandomPlayerColor5050(t *testing.T) {
	whiteCount := 0
	blackCount := 0

	for i := 0; i < 60; i++ {
		req := httptest.NewRequest(http.MethodGet, "/getBoard", nil)
		rec := httptest.NewRecorder()
		GetBoard(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("GET /getBoard status = %d, want 200", rec.Code)
		}

		req = httptest.NewRequest(http.MethodGet, "/getMoves", nil)
		rec = httptest.NewRecorder()
		GetMoves(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("GET /getMoves status = %d, want 200", rec.Code)
		}

		var movesResp struct {
			ValidMoves  map[uint64]uint64 `json:"validMoves"`
			PlayerColor PColor            `json:"playerColor"`
			OpeningMove [2]uint64         `json:"openingMove"`
		}
		if err := json.Unmarshal(rec.Body.Bytes(), &movesResp); err != nil {
			t.Fatalf("failed to unmarshal /getMoves response: %v", err)
		}

		switch movesResp.PlayerColor {
		case WHITE:
			whiteCount++
		case BLACK:
			blackCount++
			if movesResp.OpeningMove[0] == 0 || movesResp.OpeningMove[1] == 0 {
				t.Fatalf("expected non-zero OpeningMove when player is Black")
			}
		default:
			t.Fatalf("unexpected PlayerColor: %v", movesResp.PlayerColor)
		}

		if len(movesResp.ValidMoves) == 0 {
			t.Fatalf("expected non-empty ValidMoves for color %v", movesResp.PlayerColor)
		}
	}

	if whiteCount == 0 || blackCount == 0 {
		t.Fatalf("expected both White and Black assignments across 60 games, got white=%d black=%d", whiteCount, blackCount)
	}
}
