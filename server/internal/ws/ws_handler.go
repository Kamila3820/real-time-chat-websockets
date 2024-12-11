package ws

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Handler struct {
	hub *Hub
}

func NewHandler(hub *Hub) *Handler {
	return &Handler{
		hub: hub,
	}
}

type (
	CreateRoomReq struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}

	RoomRes struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}

	ClientRes struct {
		ID       string `json:"id"`
		Username string `json:"username"`
	}
)

func (h *Handler) CreateRoom(c *gin.Context) {
	var req CreateRoomReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.hub.Rooms[req.ID] = &Room{
		ID:      req.ID,
		Name:    req.Name,
		Clients: make(map[string]*Clients),
	}

	c.JSON(http.StatusOK, req)
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *Handler) JoinRoom(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	roomID := c.Param("roomId")
	clientID := c.Query("userId")
	username := c.Query("username")

	if _, ok := h.hub.Rooms[roomID]; !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "Room not found"})
		conn.Close()
		return
	}

	cl := &Clients{
		Conn:     conn,
		Message:  make(chan *Message, 10),
		ID:       clientID,
		RoomID:   roomID,
		Username: username,
	}

	m := &Message{
		Content:  "A new user has joined the room",
		RoomID:   roomID,
		Username: username,
	}

	//Register a new client through the register channel
	h.hub.Register <- cl
	//Broadcast the message
	h.hub.Broadcast <- m

	// writeMessage() for client
	// readMessage() for client

	go cl.writeMessage()
	cl.readMessage(h.hub)

}

func (h *Handler) GetRooms(c *gin.Context) {
	rooms := make([]RoomRes, 0)

	for _, r := range h.hub.Rooms {
		rooms = append(rooms, RoomRes{
			ID:   r.ID,
			Name: r.Name,
		})
	}

	c.JSON(http.StatusOK, rooms)
}

func (h *Handler) GetClients(c *gin.Context) {
	var clients []ClientRes
	roomId := c.Param("roomId")

	// if no room
	if _, ok := h.hub.Rooms[roomId]; !ok {
		clients := make([]ClientRes, 0)
		c.JSON(http.StatusOK, clients)
	}

	for _, cl := range h.hub.Rooms[roomId].Clients {
		clients = append(clients, ClientRes{
			ID:       cl.ID,
			Username: cl.Username,
		})
	}

	c.JSON(http.StatusOK, clients)

}
