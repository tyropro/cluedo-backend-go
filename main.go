package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Lobby struct {
	UUID      string
	GameState *GameState
	Players   map[uuid.UUID]Player
	Options   LobbyOptions
}

type LobbyOptions struct{}

type GameState struct {
	Players  map[uuid.UUID]GamePlayer `json:"players"`
	Solution []string                 `json:"solution"`
}

func NewGameState() GameState {
	return GameState{
		Players:  make(map[uuid.UUID]GamePlayer),
		Solution: []string{},
	}
}

// var gameState GameState

var lobbies map[string]Lobby

func init() {
	lobbies = make(map[string]Lobby)
}

func getGameState(c *gin.Context) {
	id := c.Param("id")

	game := lobbies[id]

	c.IndentedJSON(http.StatusOK, game)
}

func main() {
	router := gin.Default()

	// players endpoints
	router.GET("/:lobby_id/players", get_players)
	router.POST("/:lobby_id/players/:name", create_player)

	// game endpoints
	router.GET("/game/:id", getGameState)

	router.Run(":8080")
}
