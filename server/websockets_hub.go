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
			if err := json.Unmarshal(m, &parsedMessage); err == nil {
				parsedMessage.CreatedAt = time.Now()
				numReceived := 0
				isMapMessage := parsedMessage.Latitude != 0 && parsedMessage.Longitude != 0
				// Find receivers to send the message to
				for c := range h.connections {
					var receiverMatches bool
					if isMapMessage {
						// Send map changes to both sender and receiver
						receiverMatches = c.user.ID == parsedMessage.ReceiverID || c.user.ID == parsedMessage.SenderID
					} else {
						// Send messages only to the receiver
						receiverMatches = c.user.ID == parsedMessage.ReceiverID
					}

					if receiverMatches {
						select {
						case c.send <- m:
							numReceived++
						default:
							close(c.send)
							delete(h.connections, c)
						}
					}
				}
				if (isMapMessage && numReceived == 2) || (!isMapMessage && numReceived == 1) {
					parsedMessage.Read = true
				}
				db.Create(&parsedMessage)
			}
		}
	}
}
