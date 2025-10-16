package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/valdinei-santos/cpf-backend/internal/infra/logger"
	"github.com/valdinei-santos/cpf-backend/internal/modules/cliente/domain/domainerr"
	"github.com/valdinei-santos/cpf-backend/internal/modules/cliente/domain/entities"
)

// RepoClienteMongoDB é um repositório para gerenciar clientes no MongoDB
type RepoClienteMongoDB struct {
	collection *mongo.Collection
	log        logger.ILogger
}

// NewRepoClienteMongoDB - cria uma nova instância do repositório
func NewRepoClienteMongoDB(db *mongo.Database, collectionName string, l logger.ILogger) *RepoClienteMongoDB {
	collection := db.Collection(collectionName)
	return &RepoClienteMongoDB{
		collection: collection,
		log:        l,
	}
}

// AddCliente - adiciona um novo cliente ao repositório
func (r *RepoClienteMongoDB) AddCliente(p *entities.Cliente) error {
	ctx := context.Background()

	_, err := r.collection.InsertOne(ctx, p)
	if err != nil {
		return err
	}
	return nil
}

// GetClienteByID - busca um cliente por ID
func (r *RepoClienteMongoDB) GetClienteByID(id string) (*entities.Cliente, error) {
	ctx := context.Background()
	var cliente entities.Cliente

	idUUID, err := uuid.Parse(id)
	if err != nil {
		return nil, domainerr.ErrClienteIDInvalid
	}

	// Converte o UUID do Go para o formato Binário do Mongo (UUID)
	idBSON := primitive.Binary{
		Subtype: 4,
		Data:    idUUID[:],
	}

	// Consulta o MongoDB pelo campo _id (que deve ser o ID UUID)
	filter := bson.M{"id": idBSON}
	err = r.collection.FindOne(ctx, filter).Decode(&cliente)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domainerr.ErrClienteNotFound
		}
		return nil, err
	}

	return &cliente, nil
}

// GetAllClientes - retorna todos os clientes com paginação
func (r *RepoClienteMongoDB) GetAllClientes(offset int64, limit int64) ([]*entities.Cliente, int64, error) {
	ctx := context.Background()

	// Contagem total dos documentos
	total, err := r.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, 0, err
	}

	// Se o offset for maior ou igual ao total, retorna vazio
	if offset >= total {
		return []*entities.Cliente{}, total, nil
	}

	// Configura as opções de busca para paginação
	findOptions := options.Find()
	findOptions.SetSkip(offset)
	findOptions.SetLimit(limit)

	// Busca os documentos paginados
	cursor, err := r.collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var clientes []*entities.Cliente
	if err = cursor.All(ctx, &clientes); err != nil {
		return nil, 0, err
	}

	return clientes, total, nil
}

// UpdateCliente - atualiza os dados de um cliente existente
func (r *RepoClienteMongoDB) UpdateCliente(id string, p *entities.Cliente) error {
	ctx := context.Background()

	idUUID, err := uuid.Parse(id)
	if err != nil {
		return domainerr.ErrClienteIDInvalid
	}

	idBSON := primitive.Binary{
		Subtype: 4,
		Data:    idUUID[:],
	}

	filter := bson.M{"id": idBSON}
	updateDoc := bson.M{"$set": p}

	result, err := r.collection.UpdateOne(ctx, filter, updateDoc)
	if err != nil {
		return err
	}

	if result.ModifiedCount == 0 {
		// Se não modificou, pode ser porque não encontrou o documento.
		// É importante verificar se result.MatchedCount é 0 para ter certeza.
		if result.MatchedCount == 0 {
			return domainerr.ErrClienteNotFound
		}
		// Se MatchedCount > 0 mas ModifiedCount == 0, o documento existe,
		// mas os dados eram os mesmos. Neste caso, não é um erro.
	}

	return nil
}

// DeleteCliente - remove um cliente do repositório
func (r *RepoClienteMongoDB) DeleteCliente(id string) error {
	ctx := context.Background()

	idUUID, err := uuid.Parse(id)
	if err != nil {
		return domainerr.ErrClienteIDInvalid
	}

	idBSON := primitive.Binary{
		Subtype: 4,
		Data:    idUUID[:],
	}

	filter := bson.M{"id": idBSON}

	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return domainerr.ErrClienteNotFound
	}

	return nil
}

// Count - retorna a contagem total de clientes
func (r *RepoClienteMongoDB) Count() (int64, error) {
	ctx := context.Background()
	total, err := r.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return 0, err
	}
	return total, nil
}
