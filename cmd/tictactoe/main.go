package main

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
)

var (
	moves    = make(chan string, 10) // e.g. "O:1"
	clients  = make(map[chan string]bool)
	mu       sync.Mutex
	occupied = make(map[string]bool)
	gaming   []string
)

func main() {
	fs := http.FileServer(http.Dir("./cmd/tictactoe/static"))
	http.Handle("/", fs)
	http.HandleFunc("/move", handleMove)
	http.HandleFunc("/events", sseHandler)
	http.HandleFunc("/reset", resetGame)

	fmt.Println("Server started at :9527")
	if err := http.ListenAndServe(":9527", nil); err != nil {
		fmt.Println("Error:", err)
	}
}

func handleMove(w http.ResponseWriter, r *http.Request) {
	move := r.URL.Query().Get("move") // ex: "O:1"
	if move == "" {
		http.Error(w, "missing move", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()
	if occupied[move] {
		mu.Unlock()
		http.Error(w, "duplicate move", http.StatusConflict)
		return
	}
	occupied[move] = true
	gaming = append(gaming, move)
	moves <- move

	aiResp, err := callAI(gaming)
	if err != nil {
		http.Error(w, "AI error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	aiResp = strings.TrimSpace(aiResp)
	if aiResp == "X win" || aiResp == "O win" || aiResp == "Draw" {
		moves <- aiResp
		fmt.Fprintf(w, "game over: %s", aiResp)
		return
	}
	occupied[aiResp] = true
	gaming = append(gaming, aiResp)
	mu.Unlock()

	moves <- aiResp
	fmt.Fprint(w, "move done")
}

func callAI(gaming []string) (string, error) {
	return "X:5", nil // Mock AI response
}

func sseHandler(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	clientChan := make(chan string)
	mu.Lock()
	clients[clientChan] = true
	mu.Unlock()
	defer func() {
		mu.Lock()
		delete(clients, clientChan)
		mu.Unlock()
		close(clientChan)
	}()

	go func() {
		for move := range moves {
			mu.Lock()
			for c := range clients {
				select {
				case c <- move:
				default:
				}
			}
			mu.Unlock()
		}
	}()

	for msg := range clientChan {
		fmt.Fprintf(w, "data: %s\n\n", strings.TrimSpace(msg))
		flusher.Flush()
	}
}

func resetGame(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()
	occupied = make(map[string]bool)
	gaming = []string{}
	fmt.Fprint(w, "reset done")
}
