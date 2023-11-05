package main

import (
	"log"
	"mmyoungman/nostr_backup/internal/websocket"
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

func (cp *ConnectionList) Close() {
	for i := range cp.Connections {
		cp.Connections[i].Close()
	}
	close(cp.MessageChan)
}

func (cp *ConnectionList) CloseConnection(server string) {
	for i := range cp.Connections {
		if cp.Connections[i].Server == server {
			cp.Connections[i].Close()

			//assert len(Connections) == len(DoneChans)

			// remove connection from connList arrays
			numConns := len(cp.Connections)
			cp.Connections[i] = cp.Connections[numConns-1]
			cp.Connections = cp.Connections[:numConns-1]
			return
		}
	}
	log.Fatal("Cannot close connection", server, " as not in connection list")
}

func (cp *ConnectionList) CreateSubscriptions(subscriptionId string, filters Filters) {
	for i := range cp.Connections {
		cp.Connections[i].CreateSubscription(subscriptionId, filters)
	}
}

func (cp *ConnectionList) CloseSubscription(server string, subscriptionId string) {
	for i := range cp.Connections {
		if cp.Connections[i].Server == server {
			cp.Connections[i].CloseSubscription(subscriptionId)
			return
		}
	}
	log.Fatal("CloseSubscription fail! Could not find subscriptionId", subscriptionId, " for server", server)
}

func (cp *ConnectionList) EoseSubscription(server string, subscriptionId string) {
	for i := range cp.Connections {
		if cp.Connections[i].Server == server {
			for j := range cp.Connections[i].Subscriptions {
				if cp.Connections[i].Subscriptions[j].Id == subscriptionId {
					cp.Connections[i].Subscriptions[j].Eose = true
					return
				}
			}
			log.Fatal("EoseSubscription fail! Could not find subscription", subscriptionId, " for server", server)
		}
	}
	log.Fatal("EoseSubscription fail! Could not find server", server)
}

func (cp *ConnectionList) HasAllSubsEosed() bool {
	for i := range cp.Connections {
		if !cp.Connections[i].HasAllSubsEosed() {
			return false
		}
	}
	return true
}
