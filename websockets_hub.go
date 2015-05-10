package main

import (
	"encoding/json"
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

type MessageType struct {
	SenderId   int    `json:"senderId"`
	ReceiverId int    `json:"receiverId"`
	Message    string `json:"message"`
}

func (h *hub) run() {
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
			var parsedMessage MessageType
			err := json.Unmarshal(m, &parsedMessage)
			// TODO: create Message database object
			// if the message parsed successfully
			if err == nil {
				for c := range h.connections {
					// only send it if the receiver matches
					if c.user.Id == parsedMessage.ReceiverId {
						select {
						case c.send <- m:
						default:
							close(c.send)
							delete(h.connections, c)
						}
					}
				}
			}
		}
	}
}
