package main

import (
	"chat-bot/pkg/database/redis"
	"chat-bot/pkg/server"
	"chat-bot/pkg/web"
	"net/http"
)

func main() {
	redis.Initialize()
	web.ServeStaticFiles()
	http.HandleFunc("/chatbot", web.HandleWebSocketConnections)
	server.StartServer()
}
