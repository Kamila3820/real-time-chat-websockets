package ws

type Room struct {
	ID      string              `json:"id"`
	Name    string              `json:"name"`
	Clients map[string]*Clients `json:"clients"`
}

type Hub struct {
	Rooms map[string]*Room
}

func NewHub() *Hub {
	return &Hub{
		Rooms: make(map[string]*Room),
	}
}
