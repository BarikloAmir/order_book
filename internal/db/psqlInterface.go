package db

import (
	"database/sql"
	"fmt"
	"go.mau.fi/whatsmeow/types"
	"log"
)

var db *sql.DB

func InitDatabase() {
	var err error
	connStr := "user=postgres password=1234 dbname=psql_test_db sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
}

func SaveJid(userID string, jid types.JID) error {
	fmt.Println("in save", userID, jid)
	query := "INSERT INTO jid_table (user_id, user_name, agent, device, server, ad) VALUES ($1, $2, $3, $4, $5, $6)"
	_, err := db.Exec(query, userID, jid.User, jid.Agent, jid.Device, jid.Server, jid.AD)
	if err != nil {
		return err
	}
	return nil
}

func GetJid(userID string) (*types.JID, error) {
	var jid types.JID
	query := "SELECT user_name, agent, device, server, ad FROM jid_table WHERE user_id = $1"
	err := db.QueryRow(query, userID).Scan(&jid.User, &jid.Agent, &jid.Device, &jid.Server, &jid.AD)
	if err != nil {
		return nil, err
	}
	return &jid, nil
}
