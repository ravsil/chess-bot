package bot

import (
    "encoding/json"
    "net/http"
)

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

    board := []string{
        "8", "rook", "knight", "bishop", "queen", "king", "bishop", "knight", "rook",
        "pawn", "pawn", "pawn", "pawn", "pawn", "pawn", "pawn", "pawn",
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(board)
}