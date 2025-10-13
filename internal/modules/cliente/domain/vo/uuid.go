package vo

import (
	"github.com/google/uuid"
	"github.com/valdinei-santos/cpf-backend/internal/infra/logger"
)

type ID uuid.UUID

func NewUUID(log logger.ILogger) (ID, error) {
	uuid, err := uuid.NewRandom()
	if err != nil {
		log.Error("Falha ao gerar UUID", err)
		return ID(uuid), err
	}
	return ID(uuid), nil
}

func (id ID) String() string {
	return uuid.UUID(id).String()
}
