package main

import (
	"bookorder/internal/kafka"
	"fmt"
	"log"
	"net/http"

	"bookorder/internal/api"
	"bookorder/internal/db"

	"github.com/gorilla/mux"
)

func main() {
	fmt.Println("start server ...")

	db.InitDatabase()
	go kafka.InitKafka()

	r := mux.NewRouter()

	r.HandleFunc("/OrderBook", api.HandleOrderBook).Methods("GET")

	log.Fatal(http.ListenAndServe("localhost:8080", r))
}
