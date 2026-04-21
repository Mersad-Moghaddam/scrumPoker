package ws

import (
	"context"
	"encoding/json"
	"log"
	"sync"

	"github.com/gofiber/contrib/websocket"

	"scrum-poker/backend/internal/cache"
	"scrum-poker/backend/internal/models"
	"scrum-poker/backend/internal/service"
)

type Client struct {
	Conn     *websocket.Conn
	RoomCode string
	Token    string
	MemberID int64
}

type Hub struct {
	cache   *cache.RedisCache
	service *service.RoomService
	mu      sync.RWMutex
	rooms   map[string]map[*Client]struct{}
}

func NewHub(cache *cache.RedisCache, service *service.RoomService) *Hub {
	return &Hub{
		cache:   cache,
		service: service,
		rooms:   map[string]map[*Client]struct{}{},
	}
}

func (h *Hub) Run(ctx context.Context) {
	pubsub := h.cache.SubscribeRoomEvents(ctx)
	defer pubsub.Close()

	channel := pubsub.Channel()
	for {
		select {
		case <-ctx.Done():
			return
		case message := <-channel:
			if message == nil {
				return
			}
			var event models.RoomEvent
			if err := json.Unmarshal([]byte(message.Payload), &event); err != nil {
				log.Printf("unmarshal ws event: %v", err)
				continue
			}
			h.broadcast(event)
		}
	}
}

func (h *Hub) Register(roomCode, token string, conn *websocket.Conn) *Client {
	client := &Client{Conn: conn, RoomCode: roomCode, Token: token}
	member, err := h.service.GetRoomSnapshot(context.Background(), roomCode, token)
	if err == nil {
		client.MemberID = member.Me.ID
	}

	h.mu.Lock()
	defer h.mu.Unlock()
	if _, ok := h.rooms[roomCode]; !ok {
		h.rooms[roomCode] = map[*Client]struct{}{}
	}
	h.rooms[roomCode][client] = struct{}{}
	return client
}

func (h *Hub) Unregister(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	clients := h.rooms[client.RoomCode]
	delete(clients, client)
	if len(clients) == 0 {
		delete(h.rooms, client.RoomCode)
	}
}

func (h *Hub) broadcast(event models.RoomEvent) {
	h.mu.RLock()
	clients := h.rooms[event.RoomCode]
	h.mu.RUnlock()

	for client := range clients {
		if err := client.Conn.WriteJSON(event); err != nil {
			log.Printf("ws broadcast failed: %v", err)
		}
	}
}
