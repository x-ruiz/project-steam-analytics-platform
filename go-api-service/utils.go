package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"cloud.google.com/go/bigquery"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

// Utils
func loadEnv() {
	// Load environment variables
	err := godotenv.Load(".env")
	if err != nil {
		println("Env file not present")
	}
}

// Function to check if a string is in a slice
func stringInSlice(str string, list []string) bool {
	for _, item := range list {
		if item == str {
			return true
		}
	}
	return false
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

func postHandler(c *gin.Context) string {
	request := c.Request
	body, _ := io.ReadAll(request.Body)
	defer request.Body.Close()
	return string(body)
}

// Initialize bigquery service
func initBigQueryService() *bigquery.Client {
	ctx := context.Background()

	client, err := bigquery.NewClient(ctx, "steam-analytics-platform", option.WithCredentialsFile("./credentials/steam-analytics-platform-f3bc6b14426b.json"))
	if err != nil {
		log.Fatalln("[ERROR] BigQuery client failed to initialize:", err)
	}

	return client
}

// Execute bigquery query
func executeQuery(query string, client *bigquery.Client) ([]map[string]bigquery.Value, error) {
	ctx := context.Background()

	// Create a new BigQuery query
	q := client.Query(query)
	q.UseStandardSQL = true // Use standard SQL dialect

	// Run the query
	iter, err := q.Read(ctx)
	if err != nil {
		log.Printf("[ERROR] Error executing query: %v", err)
		return nil, err
	}

	// Iterate over query results and collect rows
	// TODO: UNDERSTAND THIS MORE
	var rows []map[string]bigquery.Value
	for {
		var row map[string]bigquery.Value
		err := iter.Next(&row) // BQ loads row into destination data structure, key is column, value is row item
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Printf("[ERROR] Error reading query result: %v", err)
			return nil, err
		}
		rows = append(rows, row)
	}
	log.Printf("ROWS %v", rows)
	return rows, nil
}
