package main

import (
	"os"

	"github.com/gin-gonic/gin"

	_ "github.com/joho/godotenv/autoload" // for .env
)

type errResp struct {
	Msg string `json:"msg"`
}

func main() {
	dev_mode, staging_mode := getEnv()

	// keep debug gin mode when in dev mode or staging
	if !dev_mode && !staging_mode {
		gin.SetMode(gin.ReleaseMode)
	}

	lobby_manager := NewLobbyManager()
	ws_manager := NewWsManager()
	router := setupRouter(lobby_manager, ws_manager)

	// only allow loopback connections on dev
	if !dev_mode || staging_mode {
		router.Run(":8080")
	} else {
		router.Run("127.0.0.1:8080")
	}
}

func getEnv() (bool, bool) {
	dev_mode := os.Getenv("MODE") == "dev"
	staging_mode := os.Getenv("MODE") == "staging"

	return dev_mode, staging_mode
}

func setupRouter(lobby_manager *LobbyManager, ws_manager *WsManager) *gin.Engine {
	router := gin.Default()

	// lobby endpoints
	router.GET("/lobbies", get_lobbies(lobby_manager))
	router.GET("/lobbies/:lobby_uuid", get_lobby(lobby_manager))

	router.POST("/lobbies", create_lobby(lobby_manager))

	// players endpoints
	router.GET("/lobbies/:lobby_uuid/players", get_players(lobby_manager))
	router.GET("/lobbies/:lobby_uuid/players/:player_uuid", get_player(lobby_manager))

	router.POST("/lobbies/:lobby_uuid/players/:name", create_player(lobby_manager))

	router.DELETE("/lobbies/:lobby_uuid/players/:player_uuid", delete_player(lobby_manager))

	// game endpoints
	router.POST("/lobbies/:lobby_uuid/game", create_game(lobby_manager))
	router.POST("/lobbies/:lobby_uuid/suggest", make_suggestion(lobby_manager))
	router.POST("/lobbies/:lobby_uuid/roll/:player_uuid", roll_dice(lobby_manager))

	router.DELETE("/lobbies/:lobby_uuid/game", delete_game(lobby_manager))

	// ws endpoints
	router.GET("/lobbies/:lobby_uuid/ws/:player_uuid", ws_handler(ws_manager, lobby_manager))
	router.GET("/lobbies/:lobby_uuid/broadcast/:message", broadcast_handler(ws_manager, lobby_manager))

	return router
}
