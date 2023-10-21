package entity

import (
	"database/sql/driver"
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
