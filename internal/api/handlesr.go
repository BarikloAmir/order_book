package api

import (
	"fmt"
	"github.com/gorilla/mux"
	"go.mau.fi/whatsmeow"
	"net/http"
	"time"
	"wconnect/internal/db"
	"wconnect/internal/whatspp"
)

func HandleConnect(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userid"]
	jid, err := db.GetJid(userID)
	fmt.Println("jid", jid, err)

	//if err != nil {
	//	fmt.Println(err)
	//	w.WriteHeader(http.StatusInternalServerError)
	//	return
	//}

	client, qrCode, LinkErr := whatspp.LinkToWhatsApp(jid)

	if LinkErr == nil {
		fmt.Println("in save jid ", userID)

		// if connected user before
		if client.Store.ID != nil {
			w.WriteHeader(200)
			w.Write([]byte("connected to whatsapp successfully !"))
			return
		}
		// but if don't connected yet
		if client.Store.ID == nil {
			w.Write([]byte(*qrCode))
			go waitForClientConnecting(userID, client)
			return
		}

	}
	if LinkErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

func waitForClientConnecting(userID string, client *whatsmeow.Client) {
	println("in the wait function")
	// Get the current time.
	startTime := time.Now()
	// Calculate the time 60 seconds from now.
	finishTime := startTime.Add(time.Second * 1200)

	fmt.Println(client.Store.ID)
	for client.Store.ID == nil {
		// Get the current time.
		now := time.Now()

		if now.After(finishTime) {
			fmt.Println("client with user id of ", userID, "disconnected because of time out")
			client.Disconnect()
			return
		}
	}

	saveErr := db.SaveJid(userID, *client.Store.ID)
	//client.Disconnect()
	if saveErr == nil {
		fmt.Println("saving jid successfully")
		return
	} else {
		fmt.Println("save error")
		return
	}
}

func HandleStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userid"]
	jid, err := db.GetJid(userID)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if &jid != nil {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("It has already tried to connect to WhatsApp jid is:"))
		_, _ = w.Write([]byte(jid.String()))

		return
	} else {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("not connected to Whatsapp"))
		return
	}

}
