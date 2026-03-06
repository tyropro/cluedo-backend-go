package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
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

type LobbyResponse struct {
	UUID        string `json:"uuid"`
	PlayerCount int    `json:"player_count"`
}

func get_lobbies(c *gin.Context) {
	var response []LobbyResponse

	for uuid, lobby := range lobbies {
		player_count := len(lobby.Players)

		lobby_response := LobbyResponse{
			UUID:        uuid,
			PlayerCount: player_count,
		}

		response = append(response, lobby_response)
	}

	c.IndentedJSON(http.StatusOK, response)
}
