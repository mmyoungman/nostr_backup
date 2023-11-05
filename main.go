package main

import (
	"fmt"
	"log"
	"mmyoungman/nostr_backup/internal/json"
	"mmyoungman/nostr_backup/internal/uuid"
	"mmyoungman/nostr_backup/internal/websocket"
	"time"
)

func main() {
	//npub := "1f0rwg0z2smrkggypqn7gctscevu22z6thch243365xt0tz8fw9uqupzj2x"
	npubHex := "4bc6e43c4a86c764208104fc8c2e18cb38a50b4bbe2eaac63aa196f588e97178"

	db := DBConnect()
	defer db.Close()

	connList := CreateConnectionList()
	servers := []string{"nos.lol", "nostr.mom"}
	for _, server := range servers {
		connList.AddConnection(server)
	}
	defer connList.Close()

	filters := Filters{{
		Authors: []string{npubHex},
		//Kinds: []int{KindTextNote,KindRepost,KindReaction},
	}}
	connList.CreateSubscriptions(uuid.NewUuid(), filters)

	numOfMessages := 0
	numOfEventMessages := 0
	numOfNewEvents := 0

	for {
		if connList.HasAllSubsEosed() {
			goto end
		}

		var connListMessage websocket.WSConnectionMessage
		select {
		case connListMessage = <-connList.MessageChan:
		case <-time.After(5 * time.Second):
			fmt.Println("No new message received in 5 seconds")
			goto end
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
			connList.EoseSubscription(server, eoseMessage.SubscriptionId)
			connList.CloseConnection(server)

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
end:
	fmt.Println("NumOfMessages: ", numOfMessages)
	fmt.Println("NumOfEventMessages: ", numOfEventMessages)
	fmt.Println("NumOfNewEvents: ", numOfNewEvents)
}
