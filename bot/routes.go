package bot

import (
	"encoding/json"
	"math/rand"
	"net/http"
)

var board [][]Pieces = nil
var playerColor PColor = WHITE

func Root(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		http.ServeFile(w, r, "src/index.html")
	} else {
		fs := http.FileServer(http.Dir("src"))
		fs.ServeHTTP(w, r)
	}
}

func GetBoard(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	board = genBoard()
	if rand.Intn(2) == 0 {
		playerColor = WHITE
	} else {
		playerColor = BLACK
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(board)
}

func GetMoves(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if board == nil {
		http.Error(w, "Board not initialized", http.StatusInternalServerError)
		return
	}
	validMoves := getValidMoves(board, playerColor)

	response := struct {
		ValidMoves  []ValidMoves `json:"validMoves"`
		PlayerColor PColor       `json:"playerColor"`
	}{
		ValidMoves:  validMoves,
		PlayerColor: playerColor,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
