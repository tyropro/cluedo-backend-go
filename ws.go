package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// websocket manager
type WsManager struct {
	clients map[string]*websocket.Conn // player uuid links to a ws connection (this leads to only 1 ws connection being valid at once)
	mu      sync.RWMutex               // read/write lock
}

func NewWsManager() *WsManager {
	return &WsManager{
		clients: make(map[string]*websocket.Conn),
	}
}

func (m *WsManager) Add(player_uuid string, conn *websocket.Conn) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.clients[player_uuid] = conn

	log.Printf("Client '%v' connected. Total: %d\n", player_uuid, len(m.clients))
}

func (m *WsManager) Remove(conn *websocket.Conn) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for player_uuid, client := range m.clients {
		if client == conn {
			delete(m.clients, player_uuid)
		}
	}

	conn.Close()

	log.Printf("Client disconnected. Total: %d\n", len(m.clients))
}

func (m *WsManager) Broadcast(lobby_manager *LobbyManager, lobby_uuid string, message string) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	lobby, _ := lobby_manager.CheckLobbyExists(lobby_uuid)

	for player_uuid := range lobby.Players {
		conn, ok := m.clients[player_uuid]
		if !ok {
			continue
		}

		err := conn.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			log.Printf("Write error: %v\n", err)
		}
	}
}

func (m *WsManager) CheckClientExists(player_uuid string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	_, ok := m.clients[player_uuid]

	return ok
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func ws_handler(ws_manager *WsManager, lobby_manager *LobbyManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		// link websocket connection to player in lobby
		lobby_uuid := c.Param("lobby_uuid")
		player_uuid := c.Param("player_uuid")

		_, err_resp := lobby_manager.CheckPlayerExists(lobby_uuid, player_uuid)

		if err_resp != nil {
			c.IndentedJSON(http.StatusNotFound, err_resp)
			return
		}

		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Printf("Upgrade error: %v\n", err)
			return
		}

		client_exists := ws_manager.CheckClientExists(player_uuid)
		if client_exists {
			c.IndentedJSON(http.StatusBadRequest, errResp{Msg: fmt.Sprintf("Client for player '%v' already connected.", player_uuid)})
			return
		}

		ws_manager.Add(player_uuid, conn)
		defer ws_manager.Remove(conn)

		// read loop - keeps connection alive
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				break
			}
		}
	}
}

func broadcast_handler(ws_manager *WsManager, lobby_manager *LobbyManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		lobby_uuid := c.Param("lobby_uuid")
		// TODO: change 'message' param to body from path
		message := c.Param("message")
		ws_manager.Broadcast(lobby_manager, lobby_uuid, message)
		c.JSON(http.StatusOK, gin.H{"status": "sent", "message": message})
	}
}
