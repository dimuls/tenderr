package mock

import (
	"go.uber.org/zap"

	"tednerr/entity"
)

type MessageSender struct {
	Logger *zap.Logger
}

func (m *MessageSender) SendMessage(c entity.Contact, message string) {
	m.Logger.Info("message send", zap.String("contact_type", c.Type), zap.String("contact_data", c.Data))
}
