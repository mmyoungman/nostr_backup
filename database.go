package nostr_backup

import (
	"database/sql"
	"log"
	"github.com/mmyoungman/nostr_backup/internal/json"

	_ "github.com/mattn/go-sqlite3"
)

func DBConnect() *sql.DB {
	db, err := sql.Open("sqlite3", "./nostr_backup.db")
	if err != nil {
		log.Fatal("Failed to open db", err)
	}

	stm, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS Events(
		id TEXT UNIQUE,
		pubkey TEXT,
		created_at UNSIGNED INT(2),
		kind INT,
		tags TEXT,
		content TEXT,
		sig TEXT
	);`)
	if err != nil {
		log.Fatal("Invalid SQL prepared to create table", err)
	}
	defer stm.Close()

	_, err = stm.Exec()
	if err != nil {
		log.Fatal("Failed to execute SQL to create table", err)
	}

	return db
}

func DBInsertEvent(db *sql.DB, event Event) (eventsAdded int) {
	stm, err := db.Prepare(`
	INSERT OR IGNORE INTO Events (id, pubkey, created_at, kind, tags, content, sig)
		VALUES (?, ?, ?, ?, ?, ?, ?);`)
	if err != nil {
		log.Fatal("Invalid SQL to insert new event", err)
	}
	defer stm.Close()

	result, err := stm.Exec(event.Id,
		event.PubKey, event.CreatedAt, event.Kind,
		event.Tags.ToJson(), DecorateJsonStr(event.Content), event.Sig)
	if err != nil {
		log.Fatal("Failed to execute query to insert new event", err)
	}

	n, _ := result.RowsAffected()
	return int(n)
}

func DBGetEvents(db *sql.DB) []Event {
	stm, err := db.Prepare(`
	SELECT id, pubkey, created_at, kind, tags, content, sig
	FROM Events
	ORDER BY created_at;`)
	if err != nil {
		log.Fatal("Invalid SQL to select all events", err)
	}
	defer stm.Close()

	rows, err := stm.Query()
	if err != nil {
		log.Fatal("Failed to execute query to fetch all events", err)
	}
	defer rows.Close()

	var events []Event = make([]Event, 0)
	for rows.Next() {
		var event Event
		var tags string

		err = rows.Scan(&event.Id, &event.PubKey, &event.CreatedAt,
		&event.Kind, &tags, &event.Content, &event.Sig)
		if err != nil {
			log.Fatal("Failed to read row of event from DB", err)
		}
		json.UnmarshalJSON([]byte(tags), event.Tags)

		events = append(events, event)
	}

	return events
}
