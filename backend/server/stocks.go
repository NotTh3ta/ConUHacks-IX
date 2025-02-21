package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
)

// GetTickerInfo returns information about a stock
//
// GET /api/v1/stocks/tickers/{symbol}
//
// Input:
//   - symbol: the ticker's symbol
//
// Output:
//   - TickerInfo: the ticker information struct
func (server *Server) GetTickerInfo(c *gin.Context) {
	symbol := c.Param("symbol")
	info, err := server.getTickerInfo(c, symbol)
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, info)
	time.Sleep(THROTTLE_TIME * time.Second)
}

func (server *Server) getTickerInfo(c *gin.Context, symbol string) (TickerInfo, error) {
	url := fmt.Sprintf("https://api.polygon.io/v3/reference/tickers?ticker=%s&active=true&limit=100&apiKey=%s", symbol, server.GetPolygonKey())
	method := "GET"

	defaultErrMsg := "Error receiving ticker info"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println("Error generating request for Polygon.io", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": defaultErrMsg})
		return TickerInfo{}, errors.New(defaultErrMsg)
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending/receiving request to Polygon.io", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": defaultErrMsg})
		return TickerInfo{}, errors.New(defaultErrMsg)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error reading Polygon.io response", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": defaultErrMsg})
		return TickerInfo{}, errors.New(defaultErrMsg)
	}
	//fmt.Println(string(body))

	// Unmarshall the unmarshalledBody
	var unmarshalledBody map[string]interface{}
	if err = json.Unmarshal(body, &unmarshalledBody); err != nil {
		fmt.Println("Error unmarshalling response", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": defaultErrMsg})
		return TickerInfo{}, errors.New(defaultErrMsg)
	}

	// Convert to correct data types
	results, ok := unmarshalledBody["results"].([]interface{})
	if !ok || len(results) == 0 {
		if !ok {
			fmt.Println("Error: results not found")
		} else {
			fmt.Println("Error: results empty")
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": defaultErrMsg})
		return TickerInfo{}, errors.New(defaultErrMsg)
	}

	// Convert each element to map[string]interface{}
	var convertedResults []map[string]interface{}
	for _, result := range results {
		resultMap, ok := result.(map[string]interface{})
		if !ok {
			fmt.Println("Error: result element is not a map")
			c.JSON(http.StatusInternalServerError, gin.H{"error": defaultErrMsg})
			return TickerInfo{}, errors.New(defaultErrMsg)
		}
		convertedResults = append(convertedResults, resultMap)
	}

	url = fmt.Sprintf("https://api.polygon.io/v2/aggs/ticker/%s/prev?adjusted=true&apiKey=%s", symbol, server.GetPolygonKey())

	req, err = http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println("Error generating request for Polygon.io", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": defaultErrMsg})
		return TickerInfo{}, errors.New(defaultErrMsg)
	}
	res, err = client.Do(req)
	if err != nil {
		fmt.Println("Error sending/receiving request to Polygon.io", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": defaultErrMsg})
		return TickerInfo{}, errors.New(defaultErrMsg)
	}
	defer res.Body.Close()

	body, err = io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error reading Polygon.io response", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": defaultErrMsg})
		return TickerInfo{}, errors.New(defaultErrMsg)
	}
	//fmt.Println(string(body))

	// Unmarshall the unmarshalledBody
	if err = json.Unmarshal(body, &unmarshalledBody); err != nil {
		fmt.Println("Error unmarshalling response", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": defaultErrMsg})
		return TickerInfo{}, errors.New(defaultErrMsg)
	}

	// Convert to correct data types
	results, ok = unmarshalledBody["results"].([]interface{})
	if !ok || len(results) == 0 {
		if !ok {
			fmt.Println("Error: results not found")
		} else {
			fmt.Println("Error: results empty")
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": defaultErrMsg})
		return TickerInfo{}, errors.New(defaultErrMsg)
	}

	// Convert each element to map[string]interface{}
	var closeResults []map[string]interface{}
	for _, result := range results {
		resultMap, ok := result.(map[string]interface{})
		if !ok {
			fmt.Println("Error: result element is not a map")
			c.JSON(http.StatusInternalServerError, gin.H{"error": defaultErrMsg})
			return TickerInfo{}, errors.New(defaultErrMsg)
		}
		closeResults = append(closeResults, resultMap)
	}

	openPrice, ok := closeResults[0]["o"].(float64)
	if !ok {
		fmt.Println("Error: open price not converted")
		c.JSON(http.StatusInternalServerError, gin.H{"error": defaultErrMsg})
		return TickerInfo{}, errors.New(defaultErrMsg)
	}
	closePrice, ok := closeResults[0]["c"].(float64)
	if !ok {
		fmt.Println("Error: close price not converted")
		c.JSON(http.StatusInternalServerError, gin.H{"error": defaultErrMsg})
		return TickerInfo{}, errors.New(defaultErrMsg)
	}

	// Use convertedResults for further processing
	info := TickerInfo{
		Symbol:          fmt.Sprintf("%v", convertedResults[0]["ticker"]),
		Name:            fmt.Sprintf("%v", convertedResults[0]["name"]),
		Industry:        "Not yet set",
		Locale:          fmt.Sprintf("%v", convertedResults[0]["locale"]),
		PrimaryExchange: fmt.Sprintf("%v", convertedResults[0]["primary_exchange"]),
		OpenPrice:       openPrice,
		ClosePrice:      closePrice,
	}
	return info, nil
}

// GetTickerHistory returns the historical prices of a stock
//
// GET /api/v1/stocks/tickers/{symbol}/history
//
// Input:
//   - symbol: the ticker's symbol
//
// Output:
//   - TickerHistory: the ticker history struct
func (server *Server) GetTickerHistory(c *gin.Context) {
	symbol := c.Param("symbol")
	url := fmt.Sprintf("https://api.polygon.io/v2/aggs/ticker/%s/range/1/day/2024-06-30/2025-02-01?adjusted=true&sort=asc&limit=5000&apiKey=%s", symbol, server.GetPolygonKey())
	method := "GET"

	defaultErrMsg := "Error receiving ticker history"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println("Error generating request", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": defaultErrMsg})
		return
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": defaultErrMsg})
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error reading response body", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": defaultErrMsg})
		return
	}
	//fmt.Println(string(body))

	// Unmarshall the unmarshalledBody
	var unmarshalledBody map[string]interface{}
	if err = json.Unmarshal(body, &unmarshalledBody); err != nil {
		fmt.Println("Error unmarshalling response", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": defaultErrMsg})
		return
	}

	// Check if results exist
	if _, exists := unmarshalledBody["results"]; !exists {
		fmt.Println("Error: results not found")
		c.JSON(http.StatusInternalServerError, gin.H{"error": defaultErrMsg})
		return
	}

	// Convert to correct data types
	results, ok := unmarshalledBody["results"].([]interface{})
	if !ok {
		fmt.Println("Error: results not converted")
		c.JSON(http.StatusInternalServerError, gin.H{"error": defaultErrMsg})
		return
	}

	// Convert each element to map[string]interface{}
	var convertedResults []map[string]interface{}
	for _, result := range results {
		resultMap, ok := result.(map[string]interface{})
		if !ok {
			fmt.Println("Error: result element is not a map")
			c.JSON(http.StatusInternalServerError, gin.H{"error": defaultErrMsg})
			return
		}
		convertedResults = append(convertedResults, resultMap)
	}

	for _, value := range convertedResults {
		if timestamp, exists := value["t"]; exists {
			// Convert to seconds
			value["time"] = timestamp.(float64) / 1000
			delete(value, "t")
		}
		if close, exists := value["c"]; exists {
			value["value"] = close
			delete(value, "c")
		}
		delete(value, "v")
		delete(value, "vw")
		delete(value, "o")
		delete(value, "h")
		delete(value, "l")
		delete(value, "n")
	}

	history := TickerHistory{
		History: convertedResults,
	}

	c.JSON(http.StatusOK, history)
	time.Sleep(THROTTLE_TIME * time.Second)
}

// GetTickerNews returns the news sentiment of a stock
//
// GET /api/v1/stocks/tickers/{symbol}/news
//
// Input:
//   - symbol: the ticker's symbol
//
// Output:
//   - TickerNews: the ticker news struct
func (server *Server) GetTickerNews(c *gin.Context) {
	symbol := c.Param("symbol")
	url := fmt.Sprintf("https://api.polygon.io/v2/reference/news?ticker=%s&order=desc&limit=350&sort=published_utc&apiKey=%s&published_utc.gte=2024-10-11T19:01:33Z", symbol, server.GetPolygonKey())
	method := "GET"

	defaultErrMsg := "Error receiving ticker news"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println("Error generating request", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": defaultErrMsg})
		return
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": defaultErrMsg})
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error reading response body", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": defaultErrMsg})
		return
	}
	//fmt.Println(string(body))

	// Unmarshall the unmarshalledBody
	var unmarshalledBody map[string]interface{}
	if err = json.Unmarshal(body, &unmarshalledBody); err != nil {
		fmt.Println("Error unmarshalling response", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": defaultErrMsg})
		return
	}

	// Convert to correct data types
	results, ok := unmarshalledBody["results"].([]interface{})
	if !ok {
		fmt.Println("Error: results not converted")
		c.JSON(http.StatusInternalServerError, gin.H{"error": defaultErrMsg})
		return
	}

	// Get number of articles
	numArticles, ok := unmarshalledBody["count"].(float64)
	if !ok {
		fmt.Println("Error: count not converted")
		c.JSON(http.StatusInternalServerError, gin.H{"error": defaultErrMsg})
		return
	}
	numArticlesInt := int(numArticles)

	// Convert each element to map[string]interface{}
	var articles []map[string]interface{}
	for _, result := range results {
		resultMap, ok := result.(map[string]interface{})
		if !ok {
			fmt.Println("Error: result element is not a map")
			c.JSON(http.StatusInternalServerError, gin.H{"error": defaultErrMsg})
			return
		}
		articles = append(articles, resultMap)
	}

	// Calculate average sentiment

	var sentiments []float64
	for _, article := range articles {
		if insights, exists := article["insights"]; exists {
			insightsList, ok := insights.([]interface{})
			if !ok {
				fmt.Println("Error: insights not converted")
				c.JSON(http.StatusInternalServerError, gin.H{"error": defaultErrMsg})
				return
			}

			for _, singleTickerInsight := range insightsList {
				convertedSingleTickerInsight, ok := singleTickerInsight.(map[string]interface{})
				if !ok {
					fmt.Println("Error: result element is not a map")
					c.JSON(http.StatusInternalServerError, gin.H{"error": defaultErrMsg})
					return
				}
				if ticker, exists := convertedSingleTickerInsight["ticker"]; exists {
					if ticker == symbol {
						if sentiment, exists := convertedSingleTickerInsight["sentiment"]; exists {
							sentimentString, ok := sentiment.(string)
							if !ok {
								fmt.Println("Error: sentiment not converted")
								c.JSON(http.StatusInternalServerError, gin.H{"error": defaultErrMsg})
								return
							}

							if sentimentString == "positive" {
								sentiments = append(sentiments, 1)
							} else if sentimentString == "negative" {
								sentiments = append(sentiments, -1)
							} else {
								sentiments = append(sentiments, 0)
							}
						}
					}
				}
			}
		}
	}

	// Calculate average sentiment
	var sumSentiment float64
	for _, sentiment := range sentiments {
		sumSentiment += sentiment
	}
	avgSentiment := sumSentiment / float64(len(sentiments))

	// Calculate standard deviation
	var sumSquaredDifferences float64
	for _, sentiment := range sentiments {
		sumSquaredDifferences += (sentiment - avgSentiment) * (sentiment - avgSentiment)
	}
	stdDevSentiment := sumSquaredDifferences / float64(len(sentiments))

	news := TickerNews{
		AverageSentiment: float32(avgSentiment),
		StdDevSentiment:  float32(stdDevSentiment),
		NumArticles:      numArticlesInt,
	}

	c.JSON(http.StatusOK, news)
	time.Sleep(THROTTLE_TIME * time.Second)
}

// GetHoldings returns the holdings of a user
//
// GET /api/v1/stocks/holdings
//
// Output:
//   - TickerHoldings: the ticker holdings struct
func (server *Server) GetHoldings(c *gin.Context) {
	holdingsInfo := getUniqueHoldings(testTickerPurchases)

	// Convert to Holding struct
	holdings := []Holding{}
	for _, holding := range holdingsInfo {
		holdings = append(holdings, Holding{Symbol: holding.Symbol, CurrentShares: holding.CurrentShares})
	}

	c.JSON(http.StatusOK, holdings)
}

// GetHoldingInfo returns the holdings of a stock
//
// GET /api/v1/stocks/holdings/:symbol
//
// Input:
//   - symbol: the ticker's symbol
//
// Output:
//   - TickerHoldings: the ticker holdings struct
func (server *Server) GetHoldingInfo(c *gin.Context) {
	// Get the user's holdings
	holdingsInfo := getUniqueHoldings(testTickerPurchases)

	// Get symbol parameter
	symbol := c.Param("symbol")

	for i, holding := range holdingsInfo {
		if holding.Symbol == symbol {
			fmt.Println("Processing holding", i, holding)
			holdingData, err := server.getTickerInfo(c, symbol)
			if err != nil {
				return
			}
			// Get the transactions for the holding
			transactions := getTransactionsByHolding(testTickerPurchases, holding)
			holdingsInfo[i].CurrentShares = transactions[len(transactions)-1].TotalShares

			// Get the history for the holding
			history, err := server.getTickerHistory(holding.Symbol)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting ticker history"})
				return
			}
			//fmt.Print("History: ", history)

			currentTransaction := 0
			// Adjust the value of the holding by the price and the number of shares, according to date
			for _, record := range history {
				//fmt.Println("Old Record: ", i, record)
				if currentTransaction+1 < len(transactions) && record["time"].(float64) > float64(transactions[currentTransaction+1].Date) {
					record["value"] = float64(transactions[currentTransaction+1].TotalShares) * record["value"].(float64)
					currentTransaction++
					record["shares"] = transactions[currentTransaction].TotalShares
				} else if record["time"].(float64) > float64(transactions[currentTransaction].Date) {
					record["value"] = float64(transactions[currentTransaction].TotalShares) * record["value"].(float64)
					record["shares"] = transactions[currentTransaction].TotalShares
				} else {
					record["value"] = 0
					record["shares"] = 0
				}
				//fmt.Println("New Record: ", i, record)
			}

			holding.History = history
			holding.ShareInfo = holdingData
			c.JSON(http.StatusOK, holding)
			return
		}
	}

	c.JSON(http.StatusBadRequest, gin.H{"error": "the requested holding does not exist in the portfolio"})
}

// Gets the unique holdings from a list of transactions
func getUniqueHoldings(transactions []StockTransaction) []HoldingInfo {
	sort.Slice(transactions, func(i, j int) bool {
		return transactions[i].Date < transactions[j].Date
	})
	uniqueHoldings := []HoldingInfo{}
	seen := map[string]bool{}

	for _, transaction := range transactions {
		if _, ok := seen[transaction.Symbol]; !ok {
			seen[transaction.Symbol] = true
			uniqueHoldings = append(uniqueHoldings, HoldingInfo{Symbol: transaction.Symbol, CurrentShares: transaction.TotalShares})
		} else {
			for i, holding := range uniqueHoldings {
				if holding.Symbol == transaction.Symbol {
					uniqueHoldings[i].CurrentShares = transaction.TotalShares
				}
			}
		}
	}

	return uniqueHoldings
}

// Gets the transactions for a specific holding, and sort them by date
func getTransactionsByHolding(transactions []StockTransaction, holding HoldingInfo) []StockTransaction {
	transactionsByHolding := []StockTransaction{}

	for _, transaction := range transactions {
		if transaction.Symbol == holding.Symbol {
			transactionsByHolding = append(transactionsByHolding, transaction)
		}
	}

	sort.Slice(transactionsByHolding, func(i, j int) bool {
		return transactionsByHolding[i].Date < transactionsByHolding[j].Date
	})

	return transactionsByHolding
}

func (server *Server) getTickerHistory(symbol string) ([]map[string]interface{}, error) {
	url := fmt.Sprintf("https://api.polygon.io/v2/aggs/ticker/%s/range/1/day/2024-06-30/2025-02-01?adjusted=true&sort=asc&limit=5000&apiKey=%s", symbol, server.GetPolygonKey())
	method := "GET"

	defaultErrMsg := "Error receiving ticker history"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println("Error generating request", err)
		return nil, err
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request", err)
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error reading response body", err)
		return nil, err
	}
	//fmt.Println(string(body))

	// Unmarshall the unmarshalledBody
	var unmarshalledBody map[string]interface{}
	if err = json.Unmarshal(body, &unmarshalledBody); err != nil {
		fmt.Println("Error unmarshalling response", err)
		return nil, err
	}

	// Check if results exist
	if _, exists := unmarshalledBody["results"]; !exists {
		fmt.Println("Error: results not found")
		return nil, fmt.Errorf(defaultErrMsg)
	}

	// Convert to correct data types
	results, ok := unmarshalledBody["results"].([]interface{})
	if !ok {
		fmt.Println("Error: results not converted")
		return nil, fmt.Errorf(defaultErrMsg)
	}

	// Convert each element to map[string]interface{}
	var convertedResults []map[string]interface{}
	for _, result := range results {
		resultMap, ok := result.(map[string]interface{})
		if !ok {
			fmt.Println("Error: result element is not a map")
			return nil, errors.New("error converting result element")
		}
		convertedResults = append(convertedResults, resultMap)
	}

	for _, value := range convertedResults {
		if timestamp, exists := value["t"]; exists {
			// Convert to seconds
			value["time"] = timestamp.(float64) / 1000
			delete(value, "t")
		}
		if close, exists := value["c"]; exists {
			value["value"] = close
			delete(value, "c")
		}
		delete(value, "v")
		delete(value, "vw")
		delete(value, "o")
		delete(value, "h")
		delete(value, "l")
		delete(value, "n")
	}

	history := TickerHistory{
		History: convertedResults,
	}

	time.Sleep(THROTTLE_TIME * time.Second)
	return history.History, nil
}
