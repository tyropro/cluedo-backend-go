package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Player struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
}

type GamePlayer struct {
	UUID string   `json:"uuid"`
	Hand []string `json:"hand"`
	X    int      `json:"x"`
	Y    int      `json:"y"`
}

type errResp struct {
	Msg string `json:"msg"`
}

func NewPlayer(player_uuid string, name string) Player {
	return Player{
		UUID: player_uuid,
		Name: name,
	}
}

// var players = []Player{}

func get_players(c *gin.Context) {
	lobbyId := c.Param("lobby_uuid")
	lobby := lobbies[lobbyId]
	players := lobby.Players

	c.IndentedJSON(http.StatusOK, players)
}

func create_player(c *gin.Context) {
	// TODO: Change name param to a body
	name := c.Param("name")
	lobby_uuid := c.Param("lobby_uuid")

	lobby := lobbies[lobby_uuid]

	// check if player already exists
	for _, player := range lobby.Players {
		if name == player.Name {
			c.IndentedJSON(http.StatusBadRequest, errResp{Msg: "Player already exists."})
			return
		}
	}

	// create new player
	new_player_uuid := uuid.New().String()
	new_player := NewPlayer(new_player_uuid, name)
	lobby.Players[new_player_uuid] = new_player

	// edit player order to add newly created player to last
	if lobby.FirstPlayer == nil {
		// if no one is added, set the new player to point to the terminator
		lobby.PlayerOrder[new_player_uuid] = nil
		lobby.FirstPlayer = &new_player_uuid
	} else {
		for player_uuid, next_player_uuid := range lobby.PlayerOrder {
			// if player is last,
			//   - set player to new player
			//   - set new player to terminator
			if next_player_uuid == nil {
				lobby.PlayerOrder[player_uuid] = &new_player_uuid
				lobby.PlayerOrder[new_player_uuid] = nil

				break
			}
		}
	}

	lobbies[lobby_uuid] = lobby

	c.IndentedJSON(http.StatusCreated, new_player)
}
