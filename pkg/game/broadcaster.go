package game

import (
	"context"
	"fmt"

	"github.com/gorilla/websocket"
)

type Broadcaster struct {
	Subscribers  []*websocket.Conn
	SubChannel   chan *websocket.Conn
	EventChannel chan *Event
}

func NewBroadcaster(eventChannel chan *Event) *Broadcaster {
	return &Broadcaster{
		EventChannel: eventChannel,
	}
}

func (gb *Broadcaster) Start(ctx context.Context) error {

	gb.SubChannel = make(chan *websocket.Conn)

	go func() {
		fmt.Println("Starting Broadcaster...")
		//message := []byte("Word to the purd!")
		for {
			select {
			case socket := <-gb.SubChannel:
				// log.Println("Broadcaster - Adding Subscriber:")
				gb.Subscribers = append(gb.Subscribers, socket)
			case event := <-gb.EventChannel:

				// case <-time.After(1 * time.Second):
				for i, subscriber := range gb.Subscribers {
					err := subscriber.WriteJSON(event)
					if err != nil {
						// log.Println("Broadcaster - Removing Subscriber:", err)
						// Unsubscribe
						gb.Subscribers[i] = gb.Subscribers[len(gb.Subscribers)-1] // Copy last element to index i.
						gb.Subscribers[len(gb.Subscribers)-1] = nil               // Erase last element (write zero value).
						gb.Subscribers = gb.Subscribers[:len(gb.Subscribers)-1]   // Truncate slice.
						// TODO: Should close the socket here first?
						subscriber.Close()
						break
					}
				}
			case <-ctx.Done():
				fmt.Println("Broadcaster Exited.")
				return
			}
		}
	}()

	return nil
}
