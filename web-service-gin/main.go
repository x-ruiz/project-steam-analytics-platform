// This serves as a first POC on basic api requests with steam.
// Ideally there will be a sync endpoint called to sync all data with bigquery
// so that api requests to the steam api are limited.

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// Models
//{"response":{"steamid":"76561198305662842","success":1}}
type UserData struct {
	Response struct {
		Steamid string `json:"steamid"`
		Success int `json:"success"`
		Username string `json:"username"`
	} `json:"response"`
}

type OwnedGame struct {
	AppId int `json:"appid"`
	GameName string `json:"name"`
	Playtime int `json:"playtime_forever"` // Total # of minutes played
	ImageIcon string `json:"img_icon_url"`
	ImageLogo string `json:"img_logo_url"`
}

type GameData struct {
	Response struct {
		GameCount int `json:"game_count"`
		OwnedGames []OwnedGame `json:"games"`
	} `json:"response"`
}


type Data struct {
	Username string `json:"username"`
	SteamId string `json:"steamid"`
	GameCount int `json:"game_count"`
	Games []OwnedGame `json:"games"`
}

type Playtime struct {
	TotalPlaytime int `json:"total_playtime"` // in minutes
	MostPlayedGame OwnedGame `json:"most_played_game"`
}


// Routing
func main() {
	loadEnv()
	router := gin.Default()
	router.GET("/health", health)
	router.GET("/getData", getData)
	router.GET("/getPlaytime", getPlaytime)

	router.Run() // default to 8080
}

// Utils
func loadEnv() {
	// Load environment variables
	err := godotenv.Load(".env")
	if err != nil {
		println("Env file not present")
	}
}


func httpGetRequest(url string) []byte {
	steam_api_key := os.Getenv("STEAM_API_KEY")
	url = strings.Replace(url, "<API_KEY>", steam_api_key, -1)
	println("URL: " + url)

	// Make request
	res, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return body

}


// Business Logic
func health(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, "healthy")
}

func getData(c *gin.Context) {
	username := c.Query("username")

	userData := getSteamUser(username)
	gameData := getGameData(userData.Response.Steamid)

	data := Data{
		Username: userData.Response.Username, 
		SteamId: userData.Response.Steamid,
		GameCount: gameData.Response.GameCount,
		Games: gameData.Response.OwnedGames,
	}

	c.IndentedJSON(http.StatusOK, data)
}

func getSteamUser(username string) UserData {
	url := "https://api.steampowered.com/ISteamUser/ResolveVanityURL/v0001/?key=<API_KEY>&vanityurl=" + username

	body := httpGetRequest(url)
	// Unmarshal response to struct
	var userData UserData
	json.Unmarshal(body, &userData)
	userData.Response.Username = username
	
	return userData
}

func (game *GameData) Modify(index int, img_icon_url string) {
  game.Response.OwnedGames[index].ImageIcon = img_icon_url
}

func getGameData(steamId string) GameData {
	url := "https://api.steampowered.com/IPlayerService/GetOwnedGames/v0001/?key=<API_KEY>&steamid=" + steamId + "&format=json&include_appinfo=true"
	var gameData GameData
	body := httpGetRequest(url)
	json.Unmarshal(body, &gameData)

		
	for idx, i := range gameData.Response.OwnedGames {
		println(i.ImageIcon)
		urlString := "http://media.steampowered.com/steamcommunity/public/images/apps/{appid}/{hash}.jpg"
		fullUrl := strings.Replace(urlString, "{appid}", fmt.Sprint(i.AppId), 1)
		fullUrl = strings.Replace(fullUrl, "{hash}", i.ImageIcon, 1)
		gameData.Modify(idx, fullUrl)
	}
	return gameData
}

// Temporary get playtime function
// TODO: store all data in bigquery and utilize bigquery for insights
// Can also do adjust playtime days (divide by 8 and divide by 24)
// Idea: games percentage of total playtime
func getPlaytime(c *gin.Context) {
	var playtimeObj Playtime
	steamId := c.Query("steamid")
	gameData := getGameData(steamId) 

	// Loop and add total playtime
	// Store most played game
	sum := 0
	mostPlayedGame := OwnedGame{Playtime: 0}

	println(gameData.Response.OwnedGames)
	for _, i := range gameData.Response.OwnedGames {
		sum += i.Playtime
		if i.Playtime > mostPlayedGame.Playtime {
			mostPlayedGame = i
		}
	}

	playtimeObj.TotalPlaytime = sum
	playtimeObj.MostPlayedGame = mostPlayedGame

	c.IndentedJSON(http.StatusOK, playtimeObj)
}

