package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"strings"
	"time"

	"cloud.google.com/go/bigquery"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

const projectID = "steam-analytics-platform"
const datasetID = "main"
const tableID = "t_user_table"

// Routing
func main() {
	loadEnv()
	router := gin.Default()

	// Configure CORS
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE"}
	config.AllowHeaders = []string{"Content-Type", "Authorization"}

	// Apply CORS to the router
	router.Use(cors.New(config))

	router.GET("/health", health)
	router.GET("/getSteamId", getSteamId)
	router.GET("/getData", getData)
	router.GET("/getPlaytime", getPlaytime)
	router.GET("/getLifetimePlaytime", getLifetimePlaytime)

	router.POST("/syncData", syncData)

	router.Run() // default to 8080
}

// Router Functions
func health(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, "healthy")
}

func getData(c *gin.Context) {
	steamId := c.Query("steamid")
	data := getConsolidateData(steamId)
	c.IndentedJSON(http.StatusOK, data)
}

func getSteamId(c *gin.Context) {
	username := c.Query("username")
	url := "https://api.steampowered.com/ISteamUser/ResolveVanityURL/v0001/?key=<API_KEY>&vanityurl=" + username

	body := httpGetRequest(url)
	// Unmarshal response to struct
	var steamIdObj SteamId
	json.Unmarshal(body, &steamIdObj)
	steamIdObj.Response.Username = username

	c.IndentedJSON(http.StatusOK, steamIdObj)
}

func (game *GameData) Modify(index int, img_icon_url string) {
	game.Response.OwnedGames[index].ImageIcon = img_icon_url
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

// Create bigquery table if does not exist already
// Per user when sync function called
func syncData(c *gin.Context) {
	bodyString := postHandler(c)

	var steamidObj SteamId
	err := json.Unmarshal([]byte(bodyString), &steamidObj)
	if err != nil {
		log.Printf("[Error] Error getting SteamID from request: %v", err)
	}
	println("[INFO] Syncing Data for " + steamidObj.Response.Steamid)

	client := initBigQueryService()

	// Call Steam API for up-to-date data
	data := getConsolidateData(steamidObj.Response.Steamid)

	// Insert up-to-date data into table
	var maxRetries = 10
	var delay = 4 * time.Second

	for attempts := 0; attempts < maxRetries; attempts++ {
		log.Printf("[INFO] Attempt %v to insert data into table %v", attempts+1, tableID)
		err := insertData(client, projectID, datasetID, tableID, data)
		if err == nil {
			log.Printf("[SUCCESS] Data successfully inserted into table %v", tableID)
			c.IndentedJSON(http.StatusOK, "Success")
			return
		}

		log.Printf("[Error] Failed to insert data: %v", err)
		if attempts < maxRetries-1 {
			log.Printf("[INFO] Retrying in %v seconds", delay)
			time.Sleep(delay)
		}

		log.Printf("[Error] Failed to insert data after %v retries", maxRetries)
		c.IndentedJSON(http.StatusInternalServerError, "Failed")
	}
}

func getLifetimePlaytime(c *gin.Context) {
	viewID := "v_steam_lifetime_playtime"
	rows := queryLifetimePlaytime(c.Query("steamid"), projectID, datasetID, viewID)
	c.IndentedJSON(http.StatusOK, rows)
}

// Business Logic
func getConsolidateData(steamId string) Data {
	playerData := getPlayerSummary(steamId)
	gameData := getGameData(steamId)

	data := Data{
		SteamId:     steamId,
		PersonaName: playerData.Response.Player[0].PersonaName,
		ProfileUrl:  playerData.Response.Player[0].ProfileUrl,
		Avatar:      playerData.Response.Player[0].Avatar,
		AvatarFull:  playerData.Response.Player[0].AvatarFull,
		GameCount:   gameData.Response.GameCount,
		Games:       gameData.Response.OwnedGames,
	}
	return data
}

func getGameData(steamId string) GameData {
	url := "https://api.steampowered.com/IPlayerService/GetOwnedGames/v0001/?key=<API_KEY>&steamid=" + steamId + "&format=json&include_appinfo=true"
	var gameData GameData
	body := httpGetRequest(url)
	json.Unmarshal(body, &gameData)

	for idx, i := range gameData.Response.OwnedGames {
		urlString := "http://media.steampowered.com/steamcommunity/public/images/apps/{appid}/{hash}.jpg"
		fullUrl := strings.Replace(urlString, "{appid}", fmt.Sprint(i.AppId), 1)
		fullUrl = strings.Replace(fullUrl, "{hash}", i.ImageIcon, 1)
		gameData.Modify(idx, fullUrl)
	}
	return gameData
}

func getPlayerSummary(steamId string) PlayerSummaryResponse {
	url := "https://api.steampowered.com/ISteamUser/GetPlayerSummaries/v0002/?key=<API_KEY>&steamids=" + steamId

	body := httpGetRequest(url)
	var playerSummaryObj PlayerSummaryResponse
	json.Unmarshal(body, &playerSummaryObj)

	return playerSummaryObj
}

// TODO: Understand why this transformation is necessary and why using OwnedGames struct did not work
func transformOwnedGames(games []OwnedGame) []map[string]interface{} {
	transformed := []map[string]interface{}{}
	for _, game := range games {
		transformed = append(transformed, map[string]interface{}{
			"appid":            game.AppId,
			"name":             game.GameName,
			"playtime_forever": game.Playtime,
			"img_icon_url":     game.ImageIcon,
			"img_logo_url":     game.ImageLogo,
		})
	}
	return transformed
}

func insertData(client *bigquery.Client, projectID string, datasetID string, tableID string, data Data) error {
	ctx := context.Background()

	// Prepare data to insert
	type Row struct {
		Timestamp   string                   `bigquery:"timestamp"`
		SteamID     string                   `bigquery:"steam_id"`
		PersonaName string                   `bigquery:"persona_name"`
		GameCount   int                      `bigquery:"game_count"`
		Games       []map[string]interface{} `bigquery:"games"`
	}

	// currentTime := time.Now().Add(-24 * time.Hour)
	currentTime := time.Now()
	row := Row{
		Timestamp:   currentTime.Format(time.RFC3339),
		SteamID:     data.SteamId,
		PersonaName: data.PersonaName,
		GameCount:   data.GameCount,
		Games:       transformOwnedGames(data.Games),
	}

	// Get the table inserter
	inserter := client.Dataset(datasetID).Table(tableID).Inserter()

	// Insert data into BigQuery
	if err := inserter.Put(ctx, []interface{}{row}); err != nil {
		log.Printf("[ERROR] Failed to insert rows into BigQuery: %v", err)
		return err
	}

	log.Printf("[INFO] Data inserted successfully into table: %v", tableID)
	return nil
}

// BQ QUERYING
func queryLifetimePlaytime(steamId string, projectID string, datasetID string, viewID string) []map[string]bigquery.Value {
	client := initBigQueryService()

	query := fmt.Sprintf("SELECT * FROM `%s.%s.%s` "+
		"WHERE steam_id = '%s' LIMIT 5",
		projectID, datasetID, viewID, steamId)

	rows, _ := executeQuery(query, client)
	return rows

}
