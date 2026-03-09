package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	_ "github.com/joho/godotenv/autoload" // for .env
)

type errResp struct {
	Msg string `json:"msg"`
}

var lobbies map[string]Lobby

// all for dev as of now
// START //
func init() {
	lobbies = make(map[string]Lobby)

	lobby_uuid := uuid.New().String()

	new_lobby := NewLobby(lobby_uuid)

	lobbies[lobby_uuid] = new_lobby

	// print out for dev purposes
	log.Println(lobby_uuid)
}

// END //

func main() {
	devMode, stagingMode := getEnv()

	// keep debug gin mode when in dev mode or staging
	if !devMode && !stagingMode {
		gin.SetMode(gin.ReleaseMode)
	}

	router := setupRouter()

	// only allow loopback connections on dev
	if !devMode || stagingMode {
		router.Run(":8080")
	} else {
		router.Run("127.0.0.1:8080")
	}
}

func getEnv() (bool, bool) {
	devMode := os.Getenv("MODE") == "dev"
	stagingMode := os.Getenv("MODE") == "staging"

	return devMode, stagingMode
}

func setupRouter() *gin.Engine {
	router := gin.Default()

	// lobby endpoints
	router.GET("/lobbies", get_lobbies)
	router.GET("/lobbies/:lobby_uuid", get_lobby)

	router.POST("/lobbies", create_lobby)

	// players endpoints
	router.GET("/lobbies/:lobby_uuid/players", get_players)
	router.GET("/lobbies/:lobby_uuid/players/:player_uuid", get_player)

	router.POST("/lobbies/:lobby_uuid/players/:name", create_player)

	router.DELETE("/lobbies/:lobby_uuid/players/:player_uuid", delete_player)

	// game endpoints
	router.POST("/lobbies/:lobby_uuid/game", create_game)

	router.DELETE("/lobbies/:lobby_uuid/game", delete_game)

	return router
}
