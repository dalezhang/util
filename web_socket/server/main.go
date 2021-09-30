package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// We'll need to define an Upgrader
// this will require a Read and Write buffer size
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}
var (
	wait    = make(chan struct{})
	wsCount = 0
)

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Home Page")
}

func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	wsCount += 1

	// upgrade this connection to a WebSocket
	// connection
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	} else {
		go wsHandler(ws)
	}

}

func wsHandler(ws *websocket.Conn) {
	log.Println("New handler: ", wsCount, ", Client Connected")
	defer func() {
		ws.Close()
	}()

	err := ws.WriteMessage(1, []byte("Hi Client!"))
	if err != nil {
		log.Println(err)
	}
	reader(ws)
}

// define a reader which will listen for
// new messages being sent to our WebSocket
// endpoint
func reader(conn *websocket.Conn) {
	for {
		// read in a message
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			// log.Println("err in ReadMessage:\n", err)
			return
		}
		// print out that message for clarity
		fmt.Println("messageType: ", messageType, "pload: ", string(p))

		// if err := conn.WriteMessage(messageType, p); err != nil {
		// 	log.Println("err in WriteMessage:\n", err)
		// 	return
		// }

	}
}

func setupRoutes() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/ws", wsEndpoint)
}

func main() {
	fmt.Println("Hellt World")
	setupRoutes()
	log.Fatal(http.ListenAndServe(":8080", nil))
}
