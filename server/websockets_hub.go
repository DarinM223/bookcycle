package server

import (
	"encoding/json"
	"time"

	"github.com/jinzhu/gorm"
)

// hub maintains the set of active connections and broadcasts messages to the
// connections.
type hub struct {
	// Registered connections.
	connections map[*connection]bool

	// Inbound messages from the connections.
	broadcast chan []byte

	// Register requests from the connections.
	register chan *connection

	// Unregister requests from connections.
	unregister chan *connection
}

var h = hub{
	broadcast:   make(chan []byte),
	register:    make(chan *connection),
	unregister:  make(chan *connection),
	connections: make(map[*connection]bool),
}

func (h *hub) run(db gorm.DB) {
	for {
		select {
		case c := <-h.register:
			h.connections[c] = true
		case c := <-h.unregister:
			if _, ok := h.connections[c]; ok {
				delete(h.connections, c)
				close(c.send)
			}
		case m := <-h.broadcast:
			var parsedMessage Message
			err := json.Unmarshal(m, &parsedMessage)
			if err == nil {
				parsedMessage.CreatedAt = time.Now()
				received := false
				if parsedMessage.Latitude != 0 && parsedMessage.Longitude != 0 {
					for c := range h.connections {
						// only send it if the sender or receiver matches
						if c.user.ID == parsedMessage.ReceiverID || c.user.ID == parsedMessage.SenderID {
							select {
							case c.send <- m:
								received = true
							default:
								close(c.send)
								delete(h.connections, c)
							}
						}
					}
				} else {
					for c := range h.connections {
						// only send it if the receiver matches
						if c.user.ID == parsedMessage.ReceiverID {
							select {
							case c.send <- m:
								received = true
							default:
								close(c.send)
								delete(h.connections, c)
							}
						}
					}
				}
				if received {
					parsedMessage.Read = true
				}
				db.Create(&parsedMessage)
			}
		}
	}
}
