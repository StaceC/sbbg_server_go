package socket

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"networkgaming.co.uk/techtest/pkg/game"
)

type GameWebSocketHandler struct {
	Upgrader    websocket.Upgrader
	Broadcaster *game.Broadcaster
}

func New(broadcaster *game.Broadcaster) *GameWebSocketHandler {
	upgrader := websocket.Upgrader{}
	return &GameWebSocketHandler{
		Upgrader:    upgrader,
		Broadcaster: broadcaster,
	}
}

func (gws *GameWebSocketHandler) Subscribe(w http.ResponseWriter, r *http.Request) {
	//ctx := r.Context()
	gws.Upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	sock, err := gws.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade:", err)
		return
	}
	//defer sock.Close()
	gws.Broadcaster.SubChannel <- sock
	for {
		_, _, err := sock.ReadMessage()
		if err != nil {
			log.Println("Websocket - Read:", err)
			break
		}
		// select {
		// // case <-time.After(5 * time.Second):
		// // 	fmt.Println("Socket - 5 second timeout!")
		// // 	return
		// case <-ctx.Done():
		// 	log.Println("Socket - We're done here!")
		// 	return
		// }
	}
}
