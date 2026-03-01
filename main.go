package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

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

var gameState GameState

var lobbies map[string]GameState

func init() {
	lobbies = make(map[string]GameState)
}

func getGameState(c *gin.Context) {
	id := c.Param("id")

	game := lobbies[id]

	c.IndentedJSON(http.StatusOK, game)
}

func main() {
	router := gin.Default()
	router.GET("/players", get_players)
	router.POST("/players/:name", create_player)
	router.GET("/game/:id", getGameState)

	router.Run("127.0.0.1:8080")
}
