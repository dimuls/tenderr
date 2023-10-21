package entity

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Rules []string

func (rs *Rules) Scan(src interface{}) error {
	return pq.Array((*[]string)(rs)).Scan(src)
}

func (rs Rules) Value() (driver.Value, error) {
	return pq.Array(rs), nil
}

type Class struct {
	ID    uuid.UUID `json:"id" db:"id"`
	Name  string    `json:"name" db:"name"`
	Rules Rules     `json:"rules" db:"rules"`
}

type Log struct {
	Time    time.Time `json:"time" db:"time"`
	ID      string    `json:"id" db:"id"`
	Message string    `json:"message" db:"message"`
	ClassID uuid.UUID `json:"class" db:"class"`
}

func (c *Contact) Scan(src interface{}) error {
	var j []byte

	switch tSrc := src.(type) {
	case []byte:
		j = tSrc
	case string:
		j = []byte(tSrc)
	default:
		return fmt.Errorf("unsupported source type %T", src)
	}

	return json.Unmarshal(j, c)
}

func (c Contact) Value() (driver.Value, error) {
	j, err := json.Marshal(c)
	return string(j), err
}

type Contact struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

type UserError struct {
	ID      uuid.UUID `json:"id" db:"id"`
	URL     string    `json:"url" db:"url"`
	Message string    `json:"message" db:"message"`
	Contact Contact   `json:"contact" db:"contact"`
}

type ErrorNotification struct {
	ID       uuid.UUID `json:"id" db:"id"`
	URL      *string   `json:"url" db:"url"`
	Message  string    `json:"message" db:"message"`
	Resolved bool      `json:"resolved" db:"resolved"`
}

type ErrorSolveWaiter struct {
	ID                  uuid.UUID `json:"id" db:"id"`
	ErrorNotificationID uuid.UUID `json:"errorNotificationId" db:"error_notification_id"`
	Contact             Contact   `json:"contact" db:"contact"`
}
