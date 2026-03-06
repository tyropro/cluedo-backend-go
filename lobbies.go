package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Lobby struct {
	UUID string `json:"uuid"`

	Players     map[string]Player  `json:"players"`
	PlayerOrder map[string]*string `json:"player_order"`
	FirstPlayer *string            `json:"first_player"`

	Options   LobbyOptions `json:"options"`
	GameState *GameState   `json:"game_state"`
}

type LobbyOptions struct{}

type GameState struct {
	Players           map[string]GamePlayer `json:"players"`
	Solution          []string              `json:"solution"`
	CurrentPlayerUUID string                `json:"current_player_uuid"`
}

func NewLobby(lobby_uuid string) Lobby {
	return Lobby{
		UUID:        lobby_uuid,
		Players:     make(map[string]Player),
		PlayerOrder: make(map[string]*string),
		FirstPlayer: nil,
		Options:     LobbyOptions{},
		GameState:   nil,
	}
}

func NewGameState(first_player string) GameState {
	return GameState{
		Players:           make(map[string]GamePlayer),
		Solution:          []string{},
		CurrentPlayerUUID: first_player,
	}
}

type GetLobbyResponse struct {
	UUID        string `json:"uuid"`
	PlayerCount int    `json:"player_count"`
}

func get_lobbies(c *gin.Context) {
	var response []GetLobbyResponse

	for uuid, lobby := range lobbies {
		player_count := len(lobby.Players)

		lobby_response := GetLobbyResponse{
			UUID:        uuid,
			PlayerCount: player_count,
		}

		response = append(response, lobby_response)
	}

	c.IndentedJSON(http.StatusOK, response)
}

func get_lobby(c *gin.Context) {
	lobby_uuid := c.Param("lobby_uuid")

	lobby, ok := lobbies[lobby_uuid]

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

type CreateLobbyResponse struct {
	UUID string `json:"uuid"`
}

func create_lobby(c *gin.Context) {
	lobby_uuid := uuid.New().String()

	new_lobby := NewLobby(lobby_uuid)

	lobbies[lobby_uuid] = new_lobby

	response := CreateLobbyResponse{
		UUID: lobby_uuid,
	}

	c.IndentedJSON(http.StatusCreated, response)
}
