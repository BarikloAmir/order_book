package main

import (
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	_ "go.mau.fi/whatsmeow/types"
	"log"
	"net/http"
	"wconnect/internal/api"
	"wconnect/internal/db"
)

func main() {
	fmt.Println("start")

	db.InitDatabase()

	r := mux.NewRouter()
	fmt.Println(r)

	r.HandleFunc("/connect/{userid}", api.HandleConnect).Methods("POST")
	r.HandleFunc("/status/{userid}", api.HandleStatus).Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", r))
}
