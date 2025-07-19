package main

import (
    "fmt"
    "net/http"
    "chess/bot"
)

func main() {
    fmt.Println("Serving on http://0.0.0.0:8080")
    
    http.HandleFunc("/", bot.Root)
    http.HandleFunc("/getBoard", bot.GetBoard)
    http.ListenAndServe("0.0.0.0:8080", nil)
}