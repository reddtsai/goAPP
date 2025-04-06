ADDR?=:8080

.PHONY: chatbot-http
chatbot-http: 
	. ./cmd/chatbot/env.sh && go run ./cmd/chatbot http --addr=$(ADDR)