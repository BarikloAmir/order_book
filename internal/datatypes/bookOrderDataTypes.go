package datatypes

import "time"

type Order struct {
	OrderID   int       `gorm:"primaryKey" json:"OrderID"`
	Side      string    `json:"side"`
	Symbol    string    `json:"symbol"`
	Amount    float64   `json:"amount"`
	Price     float64   `json:"price"`
	CreatedAt time.Time `json:"createdAt" gorm:"default:CURRENT_TIMESTAMP"`
}

type Response struct {
	LastUpdateID int          `json:"lastupdateid"`
	Bids         [][2]float64 `json:"bids"`
	Asks         [][2]float64 `json:"asks"`
}

func FormatOrders(orders []Order) [][2]float64 {
	formattedOrders := make([][2]float64, len(orders))

	for i, order := range orders {
		formattedOrders[i] = [2]float64{order.Price, order.Amount}
	}

	return formattedOrders
}
