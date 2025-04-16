package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/reddtsai/goAPP/pkg/genai"
)

var (
	CONFIG_PATH = "./cmd/chatbot"

	genAiClient genai.Client
	moves       = make(chan string, 10) // e.g. "O:1"
	clients     = make(map[chan string]bool)
	mu          sync.Mutex
	occupied    = make(map[string]bool)
	gaming      []string
)

func systemPrompt() string {
	return `你是井字棋遊戲的 AI 對手，只能扮演 "X"（AI 落子）。  
玩家總是先手，使用 "O" 落子。

你每次必須根據目前的落子紀錄，自行還原棋盤，並先檢查是否已有一方獲勝或平手。若已結束，請回應 "X win"、"O win" 或 "Draw"，不得再進行落子。

請你回應一個單一的結果，格式為以下四種其中之一：

1. "X:pos" → 代表你要下在位置 pos（1~9 的數字）  
2. "X win" → 如果你已經贏了  
3. "O win" → 如果對手玩家贏了  
4. "Draw" → 如果遊戲平手結束  

你只能回應這一行結果，不要說明理由、不輸出棋盤、不輸出多行或任何多餘文字。

以下是棋盤位置編號：
1 | 2 | 3  
4 | 5 | 6  
7 | 8 | 9`
}

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

	key := strings.Split(move, ":")
	if len(key) != 2 {
		http.Error(w, "invalid move format", http.StatusBadRequest)
		return
	}
	if key[0] != "O" {
		http.Error(w, "invalid player", http.StatusBadRequest)
		return
	}
	pos := key[1]
	if occupied[pos] {
		http.Error(w, "duplicate move", http.StatusConflict)
		return
	}
	occupied[pos] = true
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
	key = strings.Split(aiResp, ":")
	pos = key[1]
	occupied[pos] = true
	gaming = append(gaming, aiResp)
	moves <- aiResp

	fmt.Fprint(w, "move done")
}

func callAI(gaming []string) (string, error) {
	history := fmt.Sprintf("落子紀錄:[%s]", strings.Join(gaming, ","))
	fmt.Println("History:", history)
	completion, err := genAiClient.Completion(context.Background(), genai.CompletionParams{
		Messages: []genai.Message{
			{
				Role:    genai.ROLE_USER,
				Content: history,
			},
		},
	})
	if err != nil {
		return "", err
	}

	fmt.Println("AI response:", completion.Message.Content)
	return strings.TrimSpace(completion.Message.Content), nil
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
