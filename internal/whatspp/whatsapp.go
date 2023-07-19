package whatspp

import (
	"context"
	"fmt"
	"github.com/mdp/qrterminal/v3"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
	"os"
	"os/signal"
	"syscall"
)

func eventHandler(evt interface{}) {
	switch v := evt.(type) {
	case *events.Message:
		fmt.Println("Received a message!", v.Message.GetConversation())
	}
}

func LinkToWhatsApp(jid *types.JID) (*whatsmeow.Client, *string, error) {

	dbLog := waLog.Stdout("Database", "DEBUG", true)

	// Make sure you add appropriate DB connector imports, e.g. github.com/mattn/go-sqlite3 for SQLite
	//container, err := sqlstore.New("sqlite3", "C:\\sqlite\\test4.db?_foreign_keys=on", dbLog)
	container, err := sqlstore.New("postgres", "user=postgres password=1234 dbname=psql_test1 sslmode=disable", dbLog)
	if err != nil {
		return nil, nil, err
	}

	var deviceStore *store.Device

	if jid != nil {
		// If you want multiple sessions, remember their JIDs and use .GetDevice(jid) or .GetAllDevices() instead.
		device, err := container.GetDevice(*jid)
		if err != nil {
			return nil, nil, err
		}
		if device == nil {
			deviceStore = container.NewDevice()
		} else {
			deviceStore = device
		}
	} else {
		deviceStore = container.NewDevice()
	}

	clientLog := waLog.Stdout("Client", "DEBUG", true)
	client := whatsmeow.NewClient(deviceStore, clientLog)
	client.AddEventHandler(eventHandler)

	if client.Store.ID == nil {
		// No ID stored, new login
		qrChan, _ := client.GetQRChannel(context.Background())
		err = client.Connect()
		if err != nil {
			return nil, nil, err
		}
		for evt := range qrChan {
			var qrCode *string

			if evt.Event == "code" {
				// Render the QR code here
				qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)

				fmt.Println("QR code:", evt.Code)

				qrCode = &evt.Code

			} else {
				fmt.Println("Login event:", evt.Event)
			}
			if qrCode != nil {
				fmt.Println("generated QR code for client !")
				return client, qrCode, nil
			}
		}
	} else {
		// Already logged in, just connect
		err = client.Connect()
		if err != nil {
			return nil, nil, err
		}
		fmt.Println("This user connected to WhatsApp before successfully !")
		return client, nil, nil
	}

	// Listen to Ctrl+C (you can also do something else that prevents the program from exiting)
	// todo: delete below codes
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	//
	return client, nil, nil
}
