package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Order struct {
	Price string `json:"price"`
	Size  string `json:"size"`
}

type Orderbook struct {
	Bids []Order `json:"bids"`
	Asks []Order `json:"asks"`
}

func main() {
	// Read input JSON file
	jsonData, err := ioutil.ReadFile("orderbook.json")
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(1)
	}

	// Parse JSON
	var orderbook Orderbook
	err = json.Unmarshal(jsonData, &orderbook)
	if err != nil {
		fmt.Printf("Error parsing JSON: %v\n", err)
		os.Exit(1)
	}

	// Sort asks in descending order
	sort.Slice(orderbook.Asks, func(i, j int) bool {
		priceI, _ := strconv.ParseFloat(orderbook.Asks[i].Price, 64)
		priceJ, _ := strconv.ParseFloat(orderbook.Asks[j].Price, 64)
		return priceI > priceJ
	})

	// Sort bids in descending order
	sort.Slice(orderbook.Bids, func(i, j int) bool {
		priceI, _ := strconv.ParseFloat(orderbook.Bids[i].Price, 64)
		priceJ, _ := strconv.ParseFloat(orderbook.Bids[j].Price, 64)
		return priceI > priceJ
	})

	// Calculate aggregate sizes for asks (from lowest to highest price)
	askAggSizes := make(map[string]float64)
	var totalAskSize float64
	for i := len(orderbook.Asks) - 1; i >= 0; i-- {
		size, _ := strconv.ParseFloat(orderbook.Asks[i].Size, 64)
		totalAskSize += size
		askAggSizes[orderbook.Asks[i].Price] = totalAskSize
	}

	// Calculate aggregate sizes for bids (from highest to lowest price)
	bidAggSizes := make(map[string]float64)
	var totalBidSize float64
	for _, bid := range orderbook.Bids {
		size, _ := strconv.ParseFloat(bid.Size, 64)
		totalBidSize += size
		bidAggSizes[bid.Price] = totalBidSize
	}

	// Generate output
	var output strings.Builder
	output.WriteString("Orderbook Snapshot\n")
	output.WriteString("=================\n\n")

	// Write asks
	output.WriteString("Asks\n")
	output.WriteString("----\n")
	output.WriteString("Price (Subticks) | Size (SOL) | Aggregate Size (SOL)\n")
	output.WriteString("--------------------------------------------------\n")
	for _, ask := range orderbook.Asks {
		price := ask.Price
		size, _ := strconv.ParseFloat(ask.Size, 64)
		aggSize := askAggSizes[price]
		output.WriteString(fmt.Sprintf("%-15s | %-9.7f | %.7f\n", price, size, aggSize))
	}
	output.WriteString("\n")

	// Write bids
	output.WriteString("Bids\n")
	output.WriteString("----\n")
	output.WriteString("Price (Subticks) | Size (SOL) | Aggregate Size (SOL)\n")
	output.WriteString("--------------------------------------------------\n")
	for _, bid := range orderbook.Bids {
		price := bid.Price
		size, _ := strconv.ParseFloat(bid.Size, 64)
		aggSize := bidAggSizes[price]
		output.WriteString(fmt.Sprintf("%-15s | %-9.7f | %.7f\n", price, size, aggSize))
	}

	// Write output to file
	err = ioutil.WriteFile("orderbook_snapshot.txt", []byte(output.String()), 0644)
	if err != nil {
		fmt.Printf("Error writing output file: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Successfully converted orderbook to snapshot format. Check orderbook_snapshot.txt")
}
