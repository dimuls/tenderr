package clickhouse

import (
	"context"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"

	"tednerr/entity"
)

type LogStorage struct {
	Conn driver.Conn
}

func (ls *LogStorage) AddLog(l entity.Log) error {
	return ls.Conn.AsyncInsert(context.Background(), `
		insert into log (time, class_id, id, element_id, message) values (?, ?, ?, ?, ?)
	`, false, l.Time, l.ClassID, l.ID, l.ElementUD, l.Message)
}
