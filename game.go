package main

import (
	"math/rand/v2"
	"net/http"

	"github.com/gin-gonic/gin"
)

var STARTING_POINT [2]int = [...]int{0, 0}

type GameState struct {
	Players           map[string]GamePlayer `json:"players"`
	PlayerOrder       map[string]string     `json:"player_order"`
	Solution          Solution              `json:"solution"`
	CurrentPlayerUUID string                `json:"current_player_uuid"`
}

type GamePlayer struct {
	UUID string   `json:"uuid"`
	Hand []string `json:"hand"`
	X    int      `json:"x"`
	Y    int      `json:"y"`
}

type Solution struct {
	Character string `json:"character"`
	Weapon    string `json:"weapon"`
	Room      string `json:"room"`
}

func NewGamePlayer(uuid string) GamePlayer {
	return GamePlayer{
		UUID: uuid,
	}
}

func NewSolution(character string, weapon string, room string) Solution {
	return Solution{
		Character: character,
		Weapon:    weapon,
		Room:      room,
	}
}

func (l *Lobby) CreateGame() {
	gamestate := GameState{}

	// move the player order into the game state
	gamestate_player_order := make(map[string]string)

	for player_uuid, next_player_uuid := range l.PlayerOrder {
		// make the loop cyclic instead of terminate at the last player
		if next_player_uuid == nil {
			next_player_uuid = l.FirstPlayer
		}

		gamestate_player_order[player_uuid] = *next_player_uuid
	}

	gamestate.PlayerOrder = gamestate_player_order

	// make the solution
	solution, remaining_cards := l.generateSolution()
	gamestate.Solution = solution

	// randomise the hands and distribute them
	no_players := len(l.Players)

	for index := range remaining_cards { // shuffle the cards
		new_index := rand.IntN(index + 1)                                                                       // generate a new index
		remaining_cards[index], remaining_cards[new_index] = remaining_cards[new_index], remaining_cards[index] // and swap the 2 cards
	}

	var player_cards = make([][]string, no_players)

	for index, card := range remaining_cards {
		player_cards[index%no_players] = append(player_cards[index%no_players], card)
	}

	// move the players into the gamestate
	gamestate_players := make(map[string]GamePlayer)

	var index int

	for player_uuid := range l.Players {
		gamestate_player := GamePlayer{
			UUID: player_uuid,
			Hand: player_cards[index],
			X:    STARTING_POINT[0],
			Y:    STARTING_POINT[1],
		}

		gamestate_players[player_uuid] = gamestate_player

		index++
	}

	gamestate.Players = gamestate_players

	gamestate.CurrentPlayerUUID = *l.FirstPlayer

	l.GameState = &gamestate
}

func (l *Lobby) generateSolution() (Solution, []string) {
	character_cards := l.Options.Characters
	weapon_cards := l.Options.Weapons
	room_cards := l.Options.Rooms

	// generate a random solution

	solution_character_index := rand.IntN(MAX_CHARACTERS) // rand.IntN is exclusive of the upper bound so adding 1 is necessary
	solution_weapon_index := rand.IntN(MAX_WEAPONS)
	solution_room_index := rand.IntN(MAX_ROOMS)

	solution_character := character_cards[solution_character_index]
	solution_weapon := weapon_cards[solution_weapon_index]
	solution_room := room_cards[solution_room_index]

	solution := NewSolution(solution_character, solution_weapon, solution_room)

	// filter out the solution cards so the players can't have
	// the card that is in the centre, in their hand

	var remaining_cards []string

	for index, card := range character_cards {
		if index == solution_character_index {
			continue
		}

		remaining_cards = append(remaining_cards, card)
	}

	for index, card := range weapon_cards {
		if index == solution_weapon_index {
			continue
		}

		remaining_cards = append(remaining_cards, card)
	}

	for index, card := range room_cards {
		if index == solution_room_index {
			continue
		}

		remaining_cards = append(remaining_cards, card)
	}

	return solution, remaining_cards
}

func create_game(c *gin.Context) {
	lobby_uuid := c.Param("lobby_uuid")

	lobby, ok := lobbies[lobby_uuid]

	if !ok {
		c.IndentedJSON(http.StatusNotFound, errResp{Msg: "Lobby not found"})
		return
	}

	if lobby.GameState != nil {
		c.IndentedJSON(http.StatusBadRequest, errResp{Msg: "Game already exists"})
		return
	}

	lobby.CreateGame()

	c.IndentedJSON(http.StatusCreated, lobby.GameState)
}
