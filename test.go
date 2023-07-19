package main

import (
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"go.mau.fi/whatsmeow"
	_ "go.mau.fi/whatsmeow/types"
	"log"
	"net/http"
	"time"
)

//
//func eventHandler(evt interface{}) {
//	switch v := evt.(type) {
//	case *events.Message:
//		fmt.Println("Received a message!", v.Message.GetConversation())
//	}
//}
//
//func linkToWhatsApp(jid *types.JID) (*whatsmeow.Client, *string, error) {
//
//	dbLog := waLog.Stdout("Database", "DEBUG", true)
//
//	// Make sure you add appropriate DB connector imports, e.g. github.com/mattn/go-sqlite3 for SQLite
//	//container, err := sqlstore.New("sqlite3", "C:\\sqlite\\test4.db?_foreign_keys=on", dbLog)
//	container, err := sqlstore.New("postgres", "user=postgres password=1234 dbname=psql_test1 sslmode=disable", dbLog)
//	if err != nil {
//		return nil, nil, err
//	}
//
//	var deviceStore *store.Device
//
//	if jid != nil {
//		// If you want multiple sessions, remember their JIDs and use .GetDevice(jid) or .GetAllDevices() instead.
//		device, err := container.GetDevice(*jid)
//		if err != nil {
//			return nil, nil, err
//		}
//		if device == nil {
//			deviceStore = container.NewDevice()
//		} else {
//			deviceStore = device
//		}
//	} else {
//		deviceStore = container.NewDevice()
//	}
//
//	clientLog := waLog.Stdout("Client", "DEBUG", true)
//	client := whatsmeow.NewClient(deviceStore, clientLog)
//	client.AddEventHandler(eventHandler)
//
//	if client.Store.ID == nil {
//		// No ID stored, new login
//		qrChan, _ := client.GetQRChannel(context.Background())
//		err = client.Connect()
//		if err != nil {
//			return nil, nil, err
//		}
//		for evt := range qrChan {
//			var qrCode *string
//
//			if evt.Event == "code" {
//				// Render the QR code here
//				qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
//
//				fmt.Println("QR code:", evt.Code)
//
//				qrCode = &evt.Code
//
//			} else {
//				fmt.Println("Login event:", evt.Event)
//			}
//			if qrCode != nil {
//				fmt.Println("generated QR code for client !")
//				return client, qrCode, nil
//			}
//		}
//	} else {
//		// Already logged in, just connect
//		err = client.Connect()
//		if err != nil {
//			return nil, nil, err
//		}
//		fmt.Println("This user connected to WhatsApp before successfully !")
//		return client, nil, nil
//	}
//
//	// Listen to Ctrl+C (you can also do something else that prevents the program from exiting)
//	// todo: delete below codes
//	c := make(chan os.Signal)
//	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
//	<-c
//
//	//
//	return client, nil, nil
//}
//
//var db *sql.DB
//
//func initDatabase() {
//	var err error
//	connStr := "user=postgres password=1234 dbname=psql_test_db sslmode=disable"
//	db, err = sql.Open("postgres", connStr)
//	if err != nil {
//		log.Fatal(err)
//	}
//}
//
//func saveJid(userID string, jid types.JID) error {
//	fmt.Println("in save", userID, jid)
//	query := "INSERT INTO jid_table (user_id, user_name, agent, device, server, ad) VALUES ($1, $2, $3, $4, $5, $6)"
//	_, err := db.Exec(query, userID, jid.User, jid.Agent, jid.Device, jid.Server, jid.AD)
//	if err != nil {
//		return err
//	}
//	return nil
//}
//
//func getJid(userID string) (*types.JID, error) {
//	var jid types.JID
//	query := "SELECT user_name, agent, device, server, ad FROM jid_table WHERE user_id = $1"
//	err := db.QueryRow(query, userID).Scan(&jid.User, &jid.Agent, &jid.Device, &jid.Server, &jid.AD)
//	if err != nil {
//		return nil, err
//	}
//	return &jid, nil
//}

func main() {
	fmt.Println("start")

	initDatabase()

	r := mux.NewRouter()
	fmt.Println(r)

	r.HandleFunc("/connect/{userid}", connectHandler).Methods("POST")
	r.HandleFunc("/status/{userid}", statusHandler).Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", r))
}

func connectHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userid"]
	jid, err := getJid(userID)
	fmt.Println("jid", jid, err)

	//if err != nil {
	//	fmt.Println(err)
	//	w.WriteHeader(http.StatusInternalServerError)
	//	return
	//}

	client, qrCode, LinkErr := linkToWhatsApp(jid)

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

	saveErr := saveJid(userID, *client.Store.ID)
	//client.Disconnect()
	if saveErr == nil {
		fmt.Println("saving jid successfully")
		return
	} else {
		fmt.Println("save error")
		return
	}
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userid"]
	jid, err := getJid(userID)

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
