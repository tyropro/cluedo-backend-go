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

	// lobby size checks
	if len(lobby.Players) < 2 {
		c.IndentedJSON(http.StatusBadRequest, errResp{Msg: "Not enough players"})
		return
	}

	// TODO: we can fix this by limiting lobby size instead of checking on game creation
	// if len(lobby.Players) > 6 {
	// 	c.IndentedJSON(http.StatusBadRequest, errResp{Msg: "Too many players"})
	// 	return
	// }

	lobby.CreateGame()

	lobbies[lobby_uuid] = lobby

	c.IndentedJSON(http.StatusCreated, lobby.GameState)
}

func delete_game(c *gin.Context) {
	lobby_uuid := c.Param("lobby_uuid")

	lobby, ok := lobbies[lobby_uuid]

	if !ok {
		c.IndentedJSON(http.StatusNotFound, errResp{Msg: "Lobby not found"})
		return
	}

	if lobby.GameState == nil {
		c.IndentedJSON(http.StatusBadRequest, errResp{Msg: "Game not found"})
		return
	}

	lobby.GameState = nil

	lobbies[lobby_uuid] = lobby

	c.Status(http.StatusNoContent)
}

type SuggestionRequest struct {
	PlayerUUID string `json:"player_uuid"`
	Cards      struct {
		Character string `json:"character"`
		Weapon    string `json:"weapon"`
		Room      string `json:"room"`
	} `json:"cards"`
}

type SuggestionResponse struct {
	Msg        string   `json:"msg"`
	PlayerUUID string   `json:"player_uuid"`
	Cards      []string `json:"card"`
}

func make_suggestion(c *gin.Context) {
	lobby_uuid := c.Param("lobby_uuid")

	lobby, ok := lobbies[lobby_uuid]

	// quit if lobby not found
	if !ok {
		c.IndentedJSON(http.StatusNotFound, errResp{Msg: "Lobby not found"})
		return
	}

	// quit if game not created
	if lobby.GameState == nil {
		c.IndentedJSON(http.StatusNotFound, errResp{Msg: "Game not found"})
		return
	}

	var req_body SuggestionRequest

	err := c.ShouldBindJSON(&req_body)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, errResp{Msg: "Failed to marshall JSON"})
		return
	}

	_, ok = lobby.Players[req_body.PlayerUUID]

	// quit if player requesting suggestion is invalid
	if !ok {
		c.IndentedJSON(http.StatusNotFound, errResp{Msg: "Player not found"})
		return
	}

	// quit if the player requesting suggestion is not the current player
	if req_body.PlayerUUID != lobby.GameState.CurrentPlayerUUID {
		c.IndentedJSON(http.StatusBadRequest, errResp{Msg: "Not this user's turn"})
		return
	}

	player_holding_card_uuid, cards := findCard(req_body, lobby.GameState)

	if player_holding_card_uuid == nil {
		resp := SuggestionResponse{
			Msg:        "Card not found",
			PlayerUUID: "null",
			Cards:      []string{},
		}

		c.IndentedJSON(http.StatusNotFound, resp)
		return
	}

	resp := SuggestionResponse{
		Msg:        "Card Found",
		PlayerUUID: *player_holding_card_uuid,
		Cards:      *cards,
	}

	c.IndentedJSON(http.StatusOK, resp)
}

func findCard(req_body SuggestionRequest, gamestate *GameState) (*string, *[]string) {
	checking_player_uuid := gamestate.PlayerOrder[req_body.PlayerUUID]

	var held_cards []string

	for range len(gamestate.Players) - 1 { // limit number of checks to number of players - 1 (excluding original player)
		for _, card := range gamestate.Players[checking_player_uuid].Hand {
			has_character := card == req_body.Cards.Character
			has_weapon := card == req_body.Cards.Weapon
			has_room := card == req_body.Cards.Room

			if has_character || has_weapon || has_room {
				held_cards = append(held_cards, card)
			}
		}

		if len(held_cards) != 0 {
			return &checking_player_uuid, &held_cards
		}

		checking_player_uuid = gamestate.PlayerOrder[checking_player_uuid]
	}

	return nil, nil
}

type DiceRollResponse struct {
	DiceRoll int `json:"dice_roll"`
}

func roll_dice(c *gin.Context) {
	lobby_uuid := c.Param("lobby_uuid")
	player_uuid := c.Param("player_uuid")

	lobby, ok := lobbies[lobby_uuid]

	// quit if lobby not found
	if !ok {
		c.IndentedJSON(http.StatusNotFound, errResp{Msg: "Lobby not found"})
		return
	}

	// quit if game not created
	if lobby.GameState == nil {
		c.IndentedJSON(http.StatusNotFound, errResp{Msg: "Game not found"})
		return
	}

	_, ok = lobby.Players[player_uuid]

	// quit if player not found
	if !ok {
		c.IndentedJSON(http.StatusNotFound, errResp{Msg: "Player not found"})
		return
	}

	if player_uuid != lobby.GameState.CurrentPlayerUUID {
		c.IndentedJSON(http.StatusBadRequest, errResp{Msg: "Not this player's turn"})
		return
	}

	dice_roll := rand.IntN(lobby.Options.NoDiceFaces) + 1 // generates from [0,n) so + 1 to get [1,n]

	resp := DiceRollResponse{
		DiceRoll: dice_roll,
	}

	c.IndentedJSON(http.StatusOK, resp)
}
