package main

import (
	"fmt"

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

// var gameState GameState

var lobbies map[string]Lobby

func init() {
	lobbies = make(map[string]Lobby)

	lobby_uuid := uuid.New().String()

	new_lobby := NewLobby(lobby_uuid)

	lobbies[lobby_uuid] = new_lobby

	// print out for dev purposes
	fmt.Println(lobby_uuid)
}

// func getGameState(c *gin.Context) {
// 	id := c.Param("id")

// 	lobby_uuid

// 	game := lobbies[id]

// 	c.IndentedJSON(http.StatusOK, game)
// }

func main() {
	router := gin.Default()

	// players endpoints
	router.GET("/:lobby_uuid/players", get_players)

	router.POST("/:lobby_uuid/players/:name", create_player)

	// game endpoints
	// router.GET("/game/:id", getGameState)

	router.Run("127.0.0.1:8080")
}
