package vo

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/valdinei-santos/cpf-backend/internal/infra/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ID uuid.UUID

func NewUUID(log logger.ILogger) (ID, error) {
	u, err := uuid.NewRandom()
	if err != nil {
		log.Error("Falha ao gerar UUID", err)
		return ID{}, err
	}
	return ID(u), nil
}

func (id ID) String() string {
	return uuid.UUID(id).String()
}

// MarshalBSONValue implementa a interface bson.ValueMarshaler
func (id ID) MarshalBSONValue() (bsontype.Type, []byte, error) {
	u := uuid.UUID(id)
	bin := primitive.Binary{
		Subtype: 4,
		Data:    u[:],
	}
	return bson.MarshalValue(bin)
}

// UnmarshalBSONValue implementa a interface bson.ValueUnmarshaler
func (id *ID) UnmarshalBSONValue(t bsontype.Type, data []byte) error {
	var bin primitive.Binary
	if err := bson.UnmarshalValue(t, data, &bin); err != nil {
		return err
	}
	if bin.Subtype != 4 || len(bin.Data) != 16 {
		return fmt.Errorf("invalid UUID binary format")
	}
	copy((*id)[:], bin.Data)
	return nil
}

func (id ID) Bytes() []byte {
	u := uuid.UUID(id)
	return u[:]
}

func FromUUID(u uuid.UUID) ID {
	return ID(u)
}
