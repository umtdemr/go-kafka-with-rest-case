package server

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/umtdemr/go-kafka-with-rest-case/pkg/store"
	"log"
	"net/http"
)

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan string)
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type websocketServer struct {
	store *store.Store
}

type WebsocketMessage struct {
	Action string `json:"action"`
	Data   any    `json:"data"`
}

func newWebsocketServer(store *store.Store) *websocketServer {
	return &websocketServer{store}
}

func (wsServer *websocketServer) HandleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatalf("Failed to upgrade connection: %v", err)
	}
	defer ws.Close()

	clients[ws] = true

	for {
		msgType, msg, err := ws.ReadMessage()
		if msgType == 1 {
			userMsg := string(msg)
			switch userMsg {
			case "getData":
				// get the logs and send them to the user
				logData, getAllLogErr := wsServer.store.GetAllLogs()
				if getAllLogErr != nil {
					continue
				}
				val, _ := json.Marshal(&WebsocketMessage{
					Action: "allLogs",
					Data:   logData,
				})
				NotifyClients(string(val))
			}
		}
		if err != nil {
			delete(clients, ws)
			break
		}
	}
}

func HandleMessages() {
	for {
		msg := <-broadcast
		for client := range clients {
			err := client.WriteMessage(websocket.TextMessage, []byte(msg))
			if err != nil {
				client.Close()
				delete(clients, client)
			}
		}
	}
}

func NotifyClients(message string) {
	broadcast <- message
}

func StartWebSocketServer(store *store.Store) {
	wsServer := newWebsocketServer(store)
	http.HandleFunc("/ws", wsServer.HandleConnections)
	go HandleMessages()

	log.Println("WebSocket server started on :8081")
	go func() {
		log.Fatal(http.ListenAndServe(":8081", nil))
	}()
}
