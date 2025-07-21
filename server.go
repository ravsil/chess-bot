package main

import (
	"chess/bot"
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("Serving on http://0.0.0.0:8080")

	http.HandleFunc("/", bot.Root)
	http.HandleFunc("/getBoard", bot.GetBoard)
	http.HandleFunc("/getMoves", bot.GetMoves)
	http.HandleFunc("/movePiece", bot.MovePiece)
	http.ListenAndServe("0.0.0.0:8080", nil)
}
