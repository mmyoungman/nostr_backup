package websocket

import (
	"log"
	"net/url"

	"github.com/gorilla/websocket"
)

type WSConnectionMessage struct {
	Server  string
	Message string
}

type WSConnection = *websocket.Conn

func Connect(server string, messageChan chan WSConnectionMessage, doneChan chan error) WSConnection {
	URL := url.URL{Scheme: "wss", Host: server}
	conn, _, err := websocket.DefaultDialer.Dial(URL.String(), nil)
	if err != nil {
		log.Fatal("Failed to connect to websocket", server, err)
	}

	go receiveMessages(server, conn, messageChan, doneChan)

	return conn
}

func receiveMessages(server string, conn WSConnection, messageChan chan WSConnectionMessage, doneChan chan error) {
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				doneChan <- nil
				return
			}
			doneChan <- err
			return
		}

		messageChan <- WSConnectionMessage{Server: server, Message: string(message)}
	}
}

func WriteMessage(conn WSConnection, message string) {
	err := conn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		log.Fatal("Failed to write websocket message!", err)
	}
}

func WriteCloseMessage(conn WSConnection) {
	err := conn.WriteMessage(
		websocket.CloseMessage,
		websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		log.Fatal("Writing websocket close message failed!", err)
	}
}
