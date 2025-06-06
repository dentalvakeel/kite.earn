package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// Struct to represent the entire JSON response
type BulkDealsResponse struct {
	AsOnDate       string          `json:"as_on_date"`
	BulkDealsData  []BulkDealData  `json:"BULK_DEALS_DATA"`
	BulkDeals      string          `json:"BULK_DEALS"`
	ShortDeals     string          `json:"SHORT_DEALS"`
	BlockDeals     string          `json:"BLOCK_DEALS"`
	ShortDealsData []ShortDealData `json:"SHORT_DEALS_DATA"`
	BlockDealsData []BlockDealData `json:"BLOCK_DEALS_DATA"`
}

// Struct for bulk deals data
type BulkDealData struct {
	Date       string `json:"date"`
	Symbol     string `json:"symbol"`
	Name       string `json:"name"`
	ClientName string `json:"clientName"`
	BuySell    string `json:"buySell"`
	Qty        string `json:"qty"`
	WATP       string `json:"watp"`
	Remarks    string `json:"remarks"`
}

// Struct for short deals data
type ShortDealData struct {
	Date       string  `json:"date"`
	Symbol     string  `json:"symbol"`
	Name       string  `json:"name"`
	ClientName *string `json:"clientName"`
	BuySell    *string `json:"buySell"`
	Qty        string  `json:"qty"`
	WATP       *string `json:"watp"`
	Remarks    *string `json:"remarks"`
}

// Struct for block deals data
type BlockDealData struct {
	Date       string  `json:"date"`
	Symbol     string  `json:"symbol"`
	Name       string  `json:"name"`
	ClientName string  `json:"clientName"`
	BuySell    string  `json:"buySell"`
	Qty        string  `json:"qty"`
	WATP       string  `json:"watp"`
	Remarks    *string `json:"remarks"`
}

func fetchBulkDeals() {
	url := "https://www.nseindia.com/api/snapshot-capital-market-largedeal"

	// Create an HTTP client
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Create a new HTTP GET request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// Add necessary headers
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Cache-Control", "no-cache")
	req.Header.Add("Cookie", os.Getenv("Cookie")) // Use a valid NSE cookie if required

	// Send the request
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer res.Body.Close()

	// Read the response body
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	// Parse the JSON response
	var bulkDealsResponse BulkDealsResponse
	err = json.Unmarshal(body, &bulkDealsResponse)
	if err != nil {
		fmt.Println("Error unmarshalling response:", err)
		fmt.Println("Response Body:", string(body)) // Print raw response for debugging
		return
	}
	stocks := []string{} // Example stock symbols
	for _, val := range instruments {
		stocks = append(stocks, val)
	}
	// Print the bulk deals
	for _, deal := range bulkDealsResponse.BulkDealsData {
		// if symbol in stocks array
		// Print the details of each bulk deal
		if contains(stocks, deal.Symbol) {
			fmt.Printf("\033[31mDate: %s, Symbol: %s, Name: %s, Client: %s, Buy/Sell: %s, Quantity: %s, WATP: %s, Remarks: %s\033[0m\n",
				deal.Date, deal.Symbol, deal.Name, deal.ClientName, deal.BuySell, deal.Qty, deal.WATP, deal.Remarks)
		}
	}
}

func contains(slice []string, value string) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}
