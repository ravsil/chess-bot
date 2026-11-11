package bot

import (
	. "chess/game"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

var (
	game        Game
	gameMu      sync.Mutex
	playerColor PColor = WHITE
	bot         GoBot  = GoBot{Color: BLACK}
)

// Print the board with Unicode chess symbols
func printBoard() {
	// Unicode symbols for chess pieces
	movelist := game.GetValidMoves(playerColor)
	pieceSymbols := map[[2]int]string{
		// White
		{int(WHITE), int(ROOK)}:   "♖",
		{int(WHITE), int(KNIGHT)}: "♘",
		{int(WHITE), int(BISHOP)}: "♗",
		{int(WHITE), int(QUEEN)}:  "♕",
		{int(WHITE), int(KING)}:   "♔",
		{int(WHITE), int(PAWN)}:   "♙",
		// Black
		{int(BLACK), int(ROOK)}:   "♜",
		{int(BLACK), int(KNIGHT)}: "♞",
		{int(BLACK), int(BISHOP)}: "♝",
		{int(BLACK), int(QUEEN)}:  "♛",
		{int(BLACK), int(KING)}:   "♚",
		{int(BLACK), int(PAWN)}:   "♟",
	}

	fmt.Print("   a b c d e f g h\n")
	for rank := 0; rank < 8; rank++ {
		fmt.Printf("%d  ", 8-rank)
		for file := 0; file < 8; file++ {
			// Convert rank/file to bit position (0-63)
			bitPos := (7-rank)*8 + file
			square := uint64(1) << bitPos

			// Check each piece type and color
			symbol := ". "
			if game.Board.WhitePawns&square != 0 {
				symbol = pieceSymbols[[2]int{int(WHITE), int(PAWN)}] + " "
			} else if game.Board.BlackPawns&square != 0 {
				symbol = pieceSymbols[[2]int{int(BLACK), int(PAWN)}] + " "
			} else if game.Board.WhiteRooks&square != 0 {
				symbol = pieceSymbols[[2]int{int(WHITE), int(ROOK)}] + " "
			} else if game.Board.BlackRooks&square != 0 {
				symbol = pieceSymbols[[2]int{int(BLACK), int(ROOK)}] + " "
			} else if game.Board.WhiteKnights&square != 0 {
				symbol = pieceSymbols[[2]int{int(WHITE), int(KNIGHT)}] + " "
			} else if game.Board.BlackKnights&square != 0 {
				symbol = pieceSymbols[[2]int{int(BLACK), int(KNIGHT)}] + " "
			} else if game.Board.WhiteBishops&square != 0 {
				symbol = pieceSymbols[[2]int{int(WHITE), int(BISHOP)}] + " "
			} else if game.Board.BlackBishops&square != 0 {
				symbol = pieceSymbols[[2]int{int(BLACK), int(BISHOP)}] + " "
			} else if game.Board.WhiteQueens&square != 0 {
				symbol = pieceSymbols[[2]int{int(WHITE), int(QUEEN)}] + " "
			} else if game.Board.BlackQueens&square != 0 {
				symbol = pieceSymbols[[2]int{int(BLACK), int(QUEEN)}] + " "
			} else if game.Board.WhiteKing&square != 0 {
				symbol = pieceSymbols[[2]int{int(WHITE), int(KING)}] + " "
			} else if game.Board.BlackKing&square != 0 {
				symbol = pieceSymbols[[2]int{int(BLACK), int(KING)}] + " "
			} else {
				// Check if this square is a valid move destination
				for _, validMoves := range movelist {
					if validMoves&square != 0 {
						symbol = "v "
						break
					}
				}
			}

			fmt.Print(symbol)
		}
		fmt.Printf(" %d\n", 8-rank)
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

	gameMu.Lock()
	defer gameMu.Unlock()

	game.NewGame()
	printBoard()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(game.Board)
}

func GetMoves(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	gameMu.Lock()
	defer gameMu.Unlock()

	validMoves := game.GetValidMoves(playerColor)
	response := struct {
		ValidMoves  map[uint64]uint64 `json:"validMoves"`
		PlayerColor PColor            `json:"playerColor"`
	}{
		ValidMoves:  validMoves,
		PlayerColor: playerColor,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

type MoveResponse struct {
	Game
	PlayerBoard Board     `json:"PlayerBoard"`
	PlayerMove  [2]uint64 `json:"PlayerMove"`
	BotMove     [2]uint64 `json:"BotMove"`
	Status      string    `json:"Status,omitempty"`
}

func MovePiece(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Piece uint64 `json:"piece"`
		Move  uint64 `json:"move"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	gameMu.Lock()
	defer gameMu.Unlock()

	validMoves := game.GetValidMoves(playerColor)
	if v, ok := validMoves[req.Piece]; ok && v&req.Move != 0 {
		err := game.Move(req.Piece, req.Move, playerColor)
		if err != nil {
			http.Error(w, "Invalid move", http.StatusBadRequest)
			return
		}
	} else {
		http.Error(w, "Invalid move", http.StatusBadRequest)
		return
	}

	playerBoard := game.Board
	playerMove := [2]uint64{req.Piece, req.Move}

	err := bot.Play(&game, bot.Color)
	if err != nil {
		status := "you won"
		if game.IsStalemate(bot.Color) {
			status = "stalemate"
		} else {
			fmt.Println("Bot error:", err)
		}
		printBoard()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(MoveResponse{
			Game:        game,
			PlayerBoard: playerBoard,
			PlayerMove:  playerMove,
			Status:      status,
		})
		return
	}

	status := ""
	if game.IsCheckmate(playerColor) {
		status = "you lost"
	} else if game.IsStalemate(playerColor) {
		status = "stalemate"
	}

	printBoard()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(MoveResponse{
		Game:        game,
		PlayerBoard: playerBoard,
		PlayerMove:  playerMove,
		BotMove:     bot.LastMove,
		Status:      status,
	})
}
