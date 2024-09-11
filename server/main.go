package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// CORS middleware
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Next()
	})

	router.GET("/api/graphite", getGraphiteData)

	router.Run(":8080")
}

func getGraphiteData(c *gin.Context) {
	apiURL := "https://553f5731-6546-4660-bca9-42e8866c0cf9.mock.pstmn.io/v1/vpchelper/?host=vigilant-vino-iamr-02&action=GET&aws_region=us-east-2&service=graphite"
	apiData := map[string]string{
		"path":    "/render/?target=virgil.gov-vigilant-nat-01.host.hostalive.perfdata.pl.value&format=json&from=-90d&to=-30d",
		"payload": "",
	}

	jsonData, err := json.Marshal(apiData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal JSON"})
		return
	}

	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to make API request"})
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response body"})
		return
	}

	var result interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse JSON response"})
		return
	}

	c.JSON(resp.StatusCode, result)
}
