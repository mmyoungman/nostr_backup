package nostr_backup

import (
	"fmt"
)

type ClientEventMessage struct {
	Event Event
}

type ClientReqMessage struct {
	SubscriptionId string
	Filters Filters
}

type ClientCloseMessage struct {
	SubscriptionId string
}

func (em ClientEventMessage) ToJson() string {
	result := fmt.Sprintf("[\"EVENT\",%s]", em.Event.ToJson())
	DevBuildValidJson(result)
	return result
}

func (rm ClientReqMessage) ToJson() string {
	result := fmt.Sprintf("[\"REQ\",\"%s\",%s]",
		rm.SubscriptionId, rm.Filters.ToJson())
	DevBuildValidJson(result)
	return result
}

func (cm ClientCloseMessage) ToJson() string {
	result := fmt.Sprintf("[\"CLOSE\",\"%s\"]", cm.SubscriptionId)
	DevBuildValidJson(result)
	return result
}

