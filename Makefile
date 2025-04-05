run-http-chatbot:
	./cmd/chatbot/env.sh
	go run ./cmd/chatbot/main.go http --addr=:8080