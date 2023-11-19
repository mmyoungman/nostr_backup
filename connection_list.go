package main

import (
	"database/sql"
	"fmt"
	"log"
	"github.com/mmyoungman/nostr_backup/internal/json"
	"github.com/mmyoungman/nostr_backup/internal/websocket"
	"time"
)

type ConnectionListMessage struct {
	Connection *Connection
	Message    string
}

type ConnectionList struct {
	Connections [](*Connection)
	MessageChan chan websocket.WSConnectionMessage
}

func CreateConnectionList() *ConnectionList {
	var connList ConnectionList
	connList.MessageChan = make(chan websocket.WSConnectionMessage, 100)
	return &connList
}

func (cl *ConnectionList) AddConnection(server string) {
	newConn := Connect(server, cl.MessageChan)
	newConn.Server = server
	cl.Connections = append(cl.Connections, newConn)
}

func (cl *ConnectionList) Close() {
	for i := range cl.Connections {
		cl.Connections[i].Close()
	}
	close(cl.MessageChan)
}

func (cl *ConnectionList) CloseConnection(server string) {
	for i := range cl.Connections {
		if cl.Connections[i].Server == server {
			cl.Connections[i].Close()

			//assert len(Connections) == len(DoneChans)

			// remove connection from cl arrays
			numConns := len(cl.Connections)
			cl.Connections[i] = cl.Connections[numConns-1]
			cl.Connections = cl.Connections[:numConns-1]
			return
		}
	}
	log.Fatal("Cannot close connection", server, " as not in connection list")
}

func (cl *ConnectionList) CreateSubscriptions(subscriptionId string, filters Filters) {
	for i := range cl.Connections {
		cl.Connections[i].CreateSubscription(subscriptionId, filters)
	}
}

func (cl *ConnectionList) CloseSubscription(server string, subscriptionId string) {
	for i := range cl.Connections {
		if cl.Connections[i].Server == server {
			cl.Connections[i].CloseSubscription(subscriptionId)
			return
		}
	}
	log.Fatal("CloseSubscription fail! Could not find subscriptionId", subscriptionId, " for server", server)
}

func (cl *ConnectionList) EoseSubscription(server string, subscriptionId string) {
	for i := range cl.Connections {
		if cl.Connections[i].Server == server {
			for j := range cl.Connections[i].Subscriptions {
				if cl.Connections[i].Subscriptions[j].Id == subscriptionId {
					cl.Connections[i].Subscriptions[j].Eose = true
					return
				}
			}
			log.Fatal("EoseSubscription fail! Could not find subscription", subscriptionId, " for server", server)
		}
	}
	log.Fatal("EoseSubscription fail! Could not find server", server)
}

func (cl *ConnectionList) HasAllSubsEosed() bool {
	for i := range cl.Connections {
		if !cl.Connections[i].HasAllSubsEosed() {
			return false
		}
	}
	return true
}

func (cl *ConnectionList) ProcessMessages(db *sql.DB) (numOfMessages int, numOfEventMessages int, numOfNewEvents int) {
	for {
		if cl.HasAllSubsEosed() {
			return
		}

		var connListMessage websocket.WSConnectionMessage
		select {
		case connListMessage = <-cl.MessageChan:
		case <-time.After(5 * time.Second):
			fmt.Println("No new message received in 5 seconds")
			return
		}
		server := connListMessage.Server
		label, message := ProcessRelayMessage(connListMessage.Message)
		numOfMessages++

		switch label {
		case "EVENT":
			numOfEventMessages++

			var eventMessage RelayEventMessage
			err := json.UnmarshalJSON(message[0], &eventMessage.SubscriptionId)
			if err != nil {
				log.Fatal("Failed to unmarshal RelayEventMessage.SubscriptionId", err)
			}

			err = json.UnmarshalJSON(message[1], &eventMessage.Event)
			if err != nil {
				log.Fatal("Failed to unmarshal RelayEventMessage.Event", err)
			}
			generatedEventId := eventMessage.Event.GenerateEventId()
			if generatedEventId != eventMessage.Event.Id {
				log.Fatal("Incorrect Id received!")
			}

			eventHasValidSig := eventMessage.Event.IsSigValid()
			if !eventHasValidSig {
				log.Fatal("Event has invalid sig: ",
					eventMessage.Event.ToJson())
			}

			numOfNewEvents += DBInsertEvent(db, eventMessage.Event)

		case "EOSE":
			var eoseMessage RelayEoseMessage
			err := json.UnmarshalJSON(message[0], &eoseMessage.SubscriptionId)
			if err != nil {
				log.Fatal("Failed to unmarshal RelayEoseMessage.SubscriptionId", err)
			}
			cl.EoseSubscription(server, eoseMessage.SubscriptionId)
			cl.CloseConnection(server)

		case "OK":
			var okMessage RelayOkMessage
			err := json.UnmarshalJSON(message[0], &okMessage.EventId)
			if err != nil {
				log.Fatal("Failed to unmarshal RelayOkMessage.EventId", err)
			}

			err = json.UnmarshalJSON(message[1], &okMessage.Status)
			if err != nil {
				log.Fatal("Failed to unmarshal RelayOkMessage.Status", err)
			}

			err = json.UnmarshalJSON(message[2], &okMessage.Message)
			if err != nil {
				log.Fatal("Failed to unmarshal RelayOkMessage.Message", err)
			}
			okJson := okMessage.ToJson()
			fmt.Printf("RelayOkMessage: %s\n", okJson)

		case "NOTICE":
			var noticeMessage RelayNoticeMessage
			err := json.UnmarshalJSON(message[0], &noticeMessage.Message)
			if err != nil {
				log.Fatal("Failed to unmarshal RelayNoticeMessage.Message", err)
			}
			noticeJson := noticeMessage.ToJson()
			fmt.Printf("RelayNoticeMessage: %s\n", noticeJson)

		default:
			log.Fatalf("Unknown Relay Message type: \"%s\"", label)
		}
	}
}
