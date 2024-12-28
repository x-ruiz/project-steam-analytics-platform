package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// Data structures for api request responses

//{"response":{"steamid":"76561198305662842","success":1}}
type SteamUser struct {
	Response struct {
		Steamid string `json:"steamid"`
		Success int `json:"success"`
		Username string `json:"username"`
	} `json:"response"`
}


// Routing
func main() {
	loadEnv()
	router := gin.Default()
	router.GET("/health", health)
	router.GET("/getSteamUser", getSteamUser)

	router.Run() // default to 8080
}


func loadEnv() {
	// Load environment variables
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading environment variables")
	}
}


// Business Logic
func health(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, "healthy")
}


func getSteamUser(c *gin.Context) {
	steam_api_key := os.Getenv("STEAM_API_KEY")
	username := "Kodiris"
	url := "https://api.steampowered.com/ISteamUser/ResolveVanityURL/v0001/?key=" + steam_api_key + "&vanityurl=" + username

	// Make request
	res, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Unmarshal response to struct
	var steamUser SteamUser
	json.Unmarshal(body, &steamUser)
	steamUser.Response.Username = username
	

	c.IndentedJSON(http.StatusOK, steamUser)
}
// func getAlbums(c *gin.Context) {
	// c.IndentedJSON(http.StatusOK, albums)
// }
 