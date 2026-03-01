package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Player struct {
	UUID uuid.UUID `json:"uuid"`
	Name string    `json:"name"`
}

type GamePlayer struct {
	UUID uuid.UUID `json:"uuid"`
	Hand []string  `json:"hand"`
	X    int       `json:"x"`
	Y    int       `json:"y"`
}

type errResp struct {
	Msg string `json:"msg"`
}

func NewPlayer(uuid uuid.UUID, name string) Player {
	return Player{
		UUID: uuid,
		Name: name,
	}
}

// var players = []Player{}

func get_players(c *gin.Context) {
	lobbyId := c.Param("lobby_id")
	gameState := lobbies[lobbyId]

	c.IndentedJSON(http.StatusOK, gameState.Players)
}

func create_player(c *gin.Context) {
	// TODO: Change name param to a body
	name := c.Param("name")
	lobby_uuid := c.Param("lobby_uuid")

	lobby := lobbies[lobby_uuid]

	for _, player := range lobby.Players {
		if name == player.Name {
			c.IndentedJSON(http.StatusBadRequest, errResp{Msg: "Player already exists."})
			return
		}
	}

	new_player_uuid := uuid.New()

	new_player := NewPlayer(new_player_uuid, name)

	lobby.Players[new_player_uuid] = new_player

	lobbies[lobby_uuid] = lobby

	c.IndentedJSON(http.StatusCreated, new_player)
}
