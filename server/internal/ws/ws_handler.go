package ws

import (
	"net/http"

	"github.com/gin-gonic/gin"
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
