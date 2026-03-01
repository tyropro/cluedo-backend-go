package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Player struct {
	UUID  string    `json:"uuid"`
	Name  string    `json:"name"`
	Cards *[]string `json:"cards"`
}

// type GamePlayer struct {
// 	UUID  string   `json:"uuid"`
// 	Cards []string `json:"cards"`
// }

type errResp struct {
	Msg string `json:"msg"`
}

func NewPlayer(name string) Player {
	return Player{
		Name:  name,
		Cards: &[]string{},
	}
}

// var players = []Player{}

func get_players(c *gin.Context) {
	lobbyId := c.Param("lobby_id")
	gameState := lobbies[lobbyId]

	c.IndentedJSON(http.StatusOK, gameState.Players)
}

func create_player(c *gin.Context) {
	name := c.Param("name")

	new_player := NewPlayer(name)

	for _, player := range gameState.Players {
		if new_player.Name == player.Name {
			c.IndentedJSON(http.StatusBadRequest, errResp{Msg: "Player already exists."})
			return
		}
	}

	gameState.Players = append(gameState.Players, new_player)

	c.IndentedJSON(http.StatusCreated, new_player)
}
