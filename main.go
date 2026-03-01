package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Lobby struct {
	GameState GameState
	Players   map[string]Player
	Options   LobbyOptions
}

type LobbyOptions struct{}

type GameState struct {
	Players  []Player `json:"player"`
	Solution []string `json:"solution"`
}

func NewGameState() GameState {
	return GameState{
		Players:  []Player{},
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
	router.GET("/players", get_players)
	router.POST("/players/:name", create_player)

	// game endpoints
	router.GET("/game/:id", getGameState)

	router.Run(":8080")
}
