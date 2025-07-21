package bot

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
)

var board [][]Pieces = nil
var playerColor PColor = WHITE

// Print the board with Unicode chess symbols
func printBoard() {
	// Unicode symbols for chess pieces
	pieceSymbols := map[[2]int]string{
		// White
		{int(WHITE), int(ROOK)}:   "♖",
		{int(WHITE), int(KNIGHT)}: "♘",
		{int(WHITE), int(BISHOP)}: "♗",
		{int(WHITE), int(QUEEN)}:  "♕",
		{int(WHITE), int(KING)}:   "♔",
		// Black
		{int(BLACK), int(ROOK)}:   "♜",
		{int(BLACK), int(KNIGHT)}: "♞",
		{int(BLACK), int(BISHOP)}: "♝",
		{int(BLACK), int(QUEEN)}:  "♛",
		{int(BLACK), int(KING)}:   "♚",
	}
	pawnSymbols := map[PColor]string{
		WHITE: "♙",
		BLACK: "♟",
	}
	fmt.Print("   a b c d e f g h\n")
	for y := 0; y < 8; y++ {
		fmt.Printf("%d  ", 8-y)
		for x := 0; x < 8; x++ {
			piece := board[y][x]
			if piece == nil {
				fmt.Print(". ")
				continue
			}
			switch p := piece.(type) {
			case *Piece:
				symbol, ok := pieceSymbols[[2]int{int(p.Color), int(p.PieceType)}]
				if ok {
					fmt.Printf("%s ", symbol)
				} else {
					fmt.Print("? ")
				}
			case *Pawn:
				fmt.Printf("%s ", pawnSymbols[p.Color])
			default:
				fmt.Print("? ")
			}
		}
		fmt.Printf(" %d\n", 8-y)
	}
	fmt.Print("   a b c d e f g h\n")
}

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
	printBoard()
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

func MovePiece(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Piece Pos `json:"piece"`
		Move  Pos `json:"move"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	if board == nil {
		http.Error(w, "Board not initialized", http.StatusInternalServerError)
		return
	}

	validMoves := getValidMoves(board, playerColor)
	found := false
	for _, validMove := range validMoves {
		if validMove.P.GetPos() == req.Piece {
			for _, move := range validMove.Moves {
				if move == req.Move {
					validMove.P.Move(board, req.Move)
					found = true
				}
			}
		}
	}
	if !found {
		http.Error(w, "Invalid move", http.StatusBadRequest)
		return
	}
	printBoard()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success", "message": "Move received"})
}
