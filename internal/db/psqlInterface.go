package db

import (
	"bookorder/internal/datatypes"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitDatabase() {

	var err error

	db, err = gorm.Open(postgres.Open("user=postgres password=1234 dbname=order_book sslmode=disable"), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&datatypes.Order{})

}

func SaveOrder(order *datatypes.Order) error {
	fmt.Println("save order : ", order)
	if err := db.Create(order).Error; err != nil {
		return err
	}

	return nil
}

func GetSymbolOrder(symbol string, limit int) ([]datatypes.Order, []datatypes.Order, int) {

	// Query buy orders
	var buyOrders []datatypes.Order
	db.Where("symbol = ? AND side = ?", symbol, "buy").Limit(limit).Find(&buyOrders)

	// Query sell orders
	var sellOrders []datatypes.Order
	db.Where("symbol = ? AND side = ?", symbol, "sell").Limit(limit).Find(&sellOrders)

	// Query last update order
	var lastUpdatedOrder datatypes.Order
	result := db.Order("created_at desc").Where("symbol = ?", symbol).First(&lastUpdatedOrder)
	if result.Error != nil {
		log.Fatal(result.Error)
	}
	lastUpdatedOrderId := lastUpdatedOrder.OrderID

	return buyOrders, sellOrders, lastUpdatedOrderId
}
