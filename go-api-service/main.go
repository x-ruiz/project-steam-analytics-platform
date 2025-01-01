// This serves as a first POC on basic api requests with steam.
// Ideally there will be a sync endpoint called to sync all data with bigquery
// so that api requests to the steam api are limited.

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"google.golang.org/api/bigquery/v2"
	"google.golang.org/api/option"
)

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

	router.POST("/syncData", syncData)

	router.Run() // default to 8080
}

// Router Functions
func health(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, "healthy")
}

// TODO: Refactor to pull getdata logic outside of the route for reusability in sync function
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

	// Initialize BQ service
	ctx := context.Background()
	srv, err := bigquery.NewService(ctx, option.WithCredentialsFile("./credentials/steam-analytics-platform-f3bc6b14426b.json"))

	if err != nil {
		log.Fatalf("Unable to initialize Bigquery service: %v", err)
	}
	projectID := "steam-analytics-platform"
	datasetID := "main"
	tableID := "user_table"

	// Create BQ Table if does not exist for this user
	// createBQTable(srv, ctx, projectID, datasetID, tableID)

	// Call Steam API for up-to-date data
	data := getConsolidateData(steamidObj.Response.Steamid)

	// Insert up-to-date data into table
	var maxRetries = 10
	var delay = 4 * time.Second

	for attempts := 0; attempts < maxRetries; attempts++ {
		log.Printf("[INFO] Attempt %v to insert data into table %v", attempts+1, tableID)
		err := insertData(srv, projectID, datasetID, tableID, data)
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

func checkBQTableExists(srv bigquery.Service, ctx context.Context, projectID string, datasetID string, tableID string) bool {
	ts, err := srv.Tables.List(projectID, datasetID).Context(ctx).Do()

	if err != nil {
		log.Printf("[ERROR] Failed to list tables in dataset %v: %v", datasetID, err)
		return false
	}

	for _, i := range ts.Tables {
		if i.TableReference.TableId == tableID {
			return true
		}
	}
	return false
}

// func createBQTable(srv *bigquery.Service, ctx context.Context, projectID string, datasetID string, tableID string) {
// 	tableExists := checkBQTableExists(*srv, ctx, projectID, datasetID, tableID)
// 	if !tableExists {
// 		log.Printf("[INFO] Table %v does not exist in dataset %v, creating...", tableID, datasetID)
// 		// Table Schema
// 		schema := []*bigquery.TableFieldSchema{
// 			{
// 				Name:        "timestamp",
// 				Type:        "TIMESTAMP",
// 				Mode:        "REQUIRED",
// 				Description: "Timestamp of data insert",
// 			},
// 			{
// 				Name:        "steam_id",
// 				Type:        "STRING",
// 				Mode:        "REQUIRED",
// 				Description: "Steamid of the user",
// 			},
// 			{
// 				Name:        "persona_name",
// 				Type:        "STRING",
// 				Mode:        "REQUIRED",
// 				Description: "Persona Name of the user",
// 			},
// 			{
// 				Name:        "game_count",
// 				Type:        "INTEGER",
// 				Mode:        "REQUIRED",
// 				Description: "Total number of games owned by the user account",
// 			},
// 			{
// 				Name:        "games",
// 				Type:        "RECORD",
// 				Mode:        "REPEATED",
// 				Description: "List of games and metadata of games",
// 				Fields: []*bigquery.TableFieldSchema{
// 					{
// 						Name: "appid",
// 						Type: "INTEGER",
// 						Mode: "REQUIRED",
// 					},
// 					{
// 						Name: "name",
// 						Type: "STRING",
// 						Mode: "REQUIRED",
// 					},
// 					{
// 						Name: "playtime_forever",
// 						Type: "INTEGER",
// 						Mode: "REQUIRED",
// 					},
// 					{
// 						Name: "img_icon_url",
// 						Type: "STRING",
// 					},
// 					{
// 						Name: "img_logo_url",
// 						Type: "STRING",
// 					},
// 				},
// 			},
// 		}

// 		// Table metadata
// 		table := &bigquery.Table{
// 			TableReference: &bigquery.TableReference{
// 				ProjectId: projectID,
// 				DatasetId: datasetID,
// 				TableId:   tableID,
// 			},
// 			Schema: &bigquery.TableSchema{
// 				Fields: schema,
// 			},
// 			Description: "Table to store user data",
// 		}

// 		_, err := srv.Tables.Insert(projectID, datasetID, table).Context(ctx).Do()
// 		if err != nil {
// 			log.Fatalf("[ERROR] Unable to create table: %v", err)
// 		}

// 		log.Printf("[SUCCESS] Table %v created successfully", tableID)

// 	} else {
// 		log.Printf("[INFO] Table %v in Dataset %v already exists, skipping creation...", tableID, datasetID)
// 	}
// }

func insertData(srv *bigquery.Service, projectID string, datasetID string, tableID string, data Data) error {
	currentTime := time.Now().Add(-24 * time.Hour)
	bigqueryTimestamp := currentTime.Format(time.RFC3339)

	rows := []*bigquery.TableDataInsertAllRequestRows{
		{
			Json: map[string]bigquery.JsonValue{
				"timestamp":    bigqueryTimestamp,
				"steam_id":     data.SteamId,
				"persona_name": data.PersonaName,
				"game_count":   data.GameCount,
				"games":        data.Games,
			},
		},
	}

	// Create Insert Request
	insertRequest := &bigquery.TableDataInsertAllRequest{Rows: rows}

	insertCall := srv.Tabledata.InsertAll(projectID, datasetID, tableID, insertRequest)
	response, err := insertCall.Do()

	if err != nil {
		log.Printf("[ERROR] Failed to insert rows: %v", err)
		return err
	}

	if len(response.InsertErrors) > 0 {
		// Loop through each insert error and log details
		for i, insertError := range response.InsertErrors {
			log.Printf("Insert error for row %d:", i+1)

			// Loop through the errors for the specific row
			for _, errDetail := range insertError.Errors {
				log.Printf("    - Reason: %s", errDetail.Reason)
				log.Printf("    - Message: %s", errDetail.Message)
				log.Printf("    - Location: %s", errDetail.Location)
			}
		}
		return err
	} else {
		log.Printf("[INFO] Data inserted successfully to table: %v", tableID)
		return nil
	}
}
