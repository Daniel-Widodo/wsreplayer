package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func websocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	for {
		// read message from client
		mt, p, err := conn.ReadMessage()

		if err != nil {
			log.Println(err)
			break
		}
		// show message
		log.Printf("Received type: %d | message: %s, ", mt, p)

		//send message to client
		err = conn.WriteMessage(mt, []byte("sdsds"))
		if err != nil {
			log.Println(err)
			break
		}
	}
}

// giving a ticker that send time every 1 second
func wsHandlerTicker(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	tick := make(chan string)

	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	go func() {
		for {
			tick <- time.Now().Format("15:04:05")
			time.Sleep(1 * time.Second)
		}
	}()

	for {
		now := <-tick
		err = conn.WriteMessage(websocket.TextMessage, []byte(now))
		if err != nil {
			log.Println(err)
			break
		}
	}
}

func main() {
	http.HandleFunc("/websocket/replyer", websocketHandler)
	http.HandleFunc("/websocket/ticker", wsHandlerTicker)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
