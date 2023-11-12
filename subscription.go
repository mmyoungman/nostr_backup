package nostr_backup

type Subscription struct {
	Id      string
	Filters Filters
	Eose    bool
}
