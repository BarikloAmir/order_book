package api

import (
	"bookorder/internal/datatypes"
	"bookorder/internal/db"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func HandleOrderBook(w http.ResponseWriter, r *http.Request) {
	symbol := r.URL.Query().Get("symbol")

	limitStr := r.URL.Query().Get("limit")
	limit := 100 //for default value

	if limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err != nil {
			// Handle the case where "limit" is not a valid integer.
			http.Error(w, "Invalid limit parameter. It must be an integer.", http.StatusBadRequest)
			return
		}

		// Check if the parsed limit is within the valid range (less than 5000).
		if parsedLimit <= 0 || parsedLimit > 5000 {
			http.Error(w, "Invalid limit parameter. It must be between 1 and 5000.", http.StatusBadRequest)
			return
		}

		// If all checks pass, use the parsed limit.
		limit = parsedLimit
	}

	// Your code to process the "symbol" and "limit" parameters goes here.
	fmt.Println(w, "Symbol: %s, Limit: %d\n", symbol, limit)

	buyOrders, sellOrders, lastUpdateID := db.GetSymbolOrder(symbol, limit)

	response := datatypes.Response{
		LastUpdateID: lastUpdateID,
		Bids:         datatypes.FormatOrders(buyOrders),
		Asks:         datatypes.FormatOrders(sellOrders),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}
