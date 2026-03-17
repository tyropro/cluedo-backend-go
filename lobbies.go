package main

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const MAX_CHARACTERS int = 6
const MAX_WEAPONS int = 6
const MAX_ROOMS int = 9

const MAX_LOBBY_SIZE int = MAX_CHARACTERS

type LobbyManager struct {
	Lobbies map[string]Lobby `json:"lobbies"`
	mu      sync.RWMutex     `json:"-"` // for read/write locking
}

func (m *LobbyManager) checkLobbyExists(lobby_uuid string) (*Lobby, *errResp) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	lobby, ok := m.Lobbies[lobby_uuid]

	if !ok {
		return nil, &errResp{Msg: "Lobby not found"}
	}

	return &lobby, nil
}

type Lobby struct {
	UUID string `json:"uuid"`

	Players     map[string]Player  `json:"players"`
	PlayerOrder map[string]*string `json:"player_order"`
	FirstPlayer *string            `json:"first_player"`

	Options   LobbyOptions `json:"options"`
	GameState *GameState   `json:"game_state"`
}

func (l *Lobby) checkPlayerExists(lobby_manager *LobbyManager, player_uuid string) (*Player, *errResp) {
	lobby_manager.mu.RLock()
	defer lobby_manager.mu.RUnlock()

	player, ok := l.Players[player_uuid]

	if !ok {
		return nil, &errResp{Msg: "Player not found"}
	}

	return &player, nil
}

type LobbyOptions struct {
	Characters [MAX_CHARACTERS]string `json:"character"`
	Weapons    [MAX_WEAPONS]string    `json:"weapons"`
	Rooms      [MAX_ROOMS]string      `json:"rooms"`

	NoDiceFaces int `json:"no_dice_faces"`
}

func NewLobbyManager() *LobbyManager {
	return &LobbyManager{
		Lobbies: make(map[string]Lobby),
	}
}

func NewLobby(lobby_uuid string) Lobby {
	return Lobby{
		UUID:        lobby_uuid,
		Players:     make(map[string]Player),
		PlayerOrder: make(map[string]*string),
		FirstPlayer: nil,
		Options:     DefaultLobbyOptions(),
		GameState:   nil,
	}
}

func DefaultLobbyOptions() LobbyOptions {
	return LobbyOptions{
		Characters: [...]string{
			"Rev. Green",
			"Col. Mustard",
			"Dr. Orchid",
			"Mrs. Peacock",
			"Prof. Plum",
			"Miss Scarlett",
		},

		Weapons: [...]string{
			"Candlestick",
			"Dagger",
			"Lead Pipe",
			"Revolver",
			"Rope",
			"Wrench",
		},

		Rooms: [...]string{
			"Ballroom",
			"Billiard Room",
			"Conservatory",
			"Dining Room",
			"Hall",
			"Kitchen",
			"Library",
			"Lounge",
			"Study",
		},
		NoDiceFaces: 6,
	}
}

type GetLobbyResponse struct {
	UUID        string `json:"uuid"`
	PlayerCount int    `json:"player_count"`
}

func get_lobbies(lobby_manager *LobbyManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		lobby_manager.mu.RLock()
		defer lobby_manager.mu.RUnlock()

		var response []GetLobbyResponse

		for uuid, lobby := range lobby_manager.Lobbies {
			player_count := len(lobby.Players)

			lobby_response := GetLobbyResponse{
				UUID:        uuid,
				PlayerCount: player_count,
			}

			response = append(response, lobby_response)
		}

		c.IndentedJSON(http.StatusOK, response)
	}
}

func get_lobby(lobby_manager *LobbyManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		lobby_manager.mu.RLock()
		defer lobby_manager.mu.RUnlock()

		lobby_uuid := c.Param("lobby_uuid")

		lobby, ok := lobby_manager.Lobbies[lobby_uuid]

		if !ok {
			c.IndentedJSON(http.StatusNotFound, errResp{Msg: "Lobby not found"})
			return
		}

		player_count := len(lobby.Players)
		response := GetLobbyResponse{
			UUID:        lobby_uuid,
			PlayerCount: player_count,
		}

		c.IndentedJSON(http.StatusOK, response)
	}
}

type CreateLobbyResponse struct {
	UUID string `json:"uuid"`
}

func create_lobby(lobby_manager *LobbyManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		lobby_manager.mu.Lock()
		defer lobby_manager.mu.Unlock()

		lobby_uuid := uuid.New().String()

		new_lobby := NewLobby(lobby_uuid)

		lobby_manager.Lobbies[lobby_uuid] = new_lobby

		response := CreateLobbyResponse{
			UUID: lobby_uuid,
		}

		c.IndentedJSON(http.StatusCreated, response)
	}
}
