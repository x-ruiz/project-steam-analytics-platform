package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
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
