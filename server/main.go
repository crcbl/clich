package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var (
	clients   = make(map[*websocket.Conn]bool)
	broadcast = make(chan []byte)
	mutex     = &sync.Mutex{}
)

func main() {
	port := os.Getenv("PORT")
	fmt.Println("starting clich server...")
	http.HandleFunc("/ws", wsHandler)
	go handleMessages()
	fmt.Println("clich server started on port " + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	println("new connection")
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading new connection:", err)
		return
	}
    defer conn.Close()

	mutex.Lock()
	clients[conn] = true
	mutex.Unlock()

	for {
		typ, msg, err := conn.ReadMessage()
		if err != nil {
			println("Error reading message: ", err)
			mutex.Lock()
			delete(clients, conn)
			mutex.Unlock()
		}
		if len(msg) != 0 {
			if typ == websocket.CloseMessage {
                println("Client disconnected. Closing...")
				mutex.Lock()
				delete(clients, conn)
				mutex.Unlock()
                return
			} else {
				broadcast <- msg
			}
		}
	}
}

func handleMessages() {
	for {
		msg := <-broadcast

		mutex.Lock()
		for client := range clients {
			println("Writing message to clients")
			err := client.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				println("Error writing message: ", err)
				client.Close()
				delete(clients, client)
			}
		}
		mutex.Unlock()
	}
}
