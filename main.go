package main

import (
	"fmt"
	"mmyoungman/nostr_backup/internal/uuid"
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

	numOfMessages, numOfEventMessages, numOfNewEvents := ProcessMessages(connList, db)

	fmt.Println("NumOfMessages: ", numOfMessages)
	fmt.Println("NumOfEventMessages: ", numOfEventMessages)
	fmt.Println("NumOfNewEvents: ", numOfNewEvents)
}

