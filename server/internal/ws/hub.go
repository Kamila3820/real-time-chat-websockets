package ws

import "log"

type Room struct {
	ID      string              `json:"id"`
	Name    string              `json:"name"`
	Clients map[string]*Clients `json:"clients"`
}

type Hub struct {
	Rooms      map[string]*Room
	Register   chan *Clients
	Unregister chan *Clients
	Broadcast  chan *Message
}

func NewHub() *Hub {
	return &Hub{
		Rooms:      make(map[string]*Room),
		Register:   make(chan *Clients),
		Unregister: make(chan *Clients),
		Broadcast:  make(chan *Message, 5),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case cl := <-h.Register:
			// if _, ok := h.Rooms[cl.RoomID]; ok {
			// 	r := h.Rooms[cl.RoomID]

			// 	if _, ok := r.Clients[cl.ID]; ok {
			// 		r.Clients[cl.ID] = cl
			// 	}
			// }
			if room, ok := h.Rooms[cl.RoomID]; ok {
				room.Clients[cl.ID] = cl // Add client to room
			} else {
				// Handle the case where the room does not exist
				log.Printf("Room %s does not exist, cannot register client", cl.RoomID)
			}

		case cl := <-h.Unregister:
			if _, ok := h.Rooms[cl.RoomID]; ok {
				if _, ok := h.Rooms[cl.RoomID].Clients[cl.ID]; ok {
					//broadcast a message saying that the client has left the room
					if len(h.Rooms[cl.RoomID].Clients) != 0 {
						h.Broadcast <- &Message{
							Content:  "user left the chat",
							RoomID:   cl.RoomID,
							Username: cl.Username,
						}
					}

					delete(h.Rooms[cl.RoomID].Clients, cl.ID)
					close(cl.Message)
				}
			}

		case m := <-h.Broadcast:
			// if _, ok := h.Rooms[m.RoomID]; ok {

			// 	for _, cl := range h.Rooms[m.RoomID].Clients {
			// 		cl.Message <- m
			// 	}
			// }
			log.Printf("Broadcasting message to room %s: %s", m.RoomID, m.Content)
			if room, ok := h.Rooms[m.RoomID]; ok {
				for _, cl := range room.Clients {
					select {
					case cl.Message <- m:
						// Successfully sent
					default:
						log.Printf("Client %s message channel is full, dropping message", cl.ID)
					}
				}
			} else {
				log.Printf("Room %s not found for broadcasting", m.RoomID)
			}

		}
	}
}
