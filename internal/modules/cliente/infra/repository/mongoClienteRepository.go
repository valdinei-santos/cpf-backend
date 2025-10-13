package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
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
func NewRepoClienteMongoDB(client *mongo.Client, databaseName, collectionName string, l logger.ILogger) *RepoClienteMongoDB {
	collection := client.Database(databaseName).Collection(collectionName)
	return &RepoClienteMongoDB{
		collection: collection,
		log:        l,
	}
}

// AddCliente - adiciona um novo cliente ao repositório
func (r *RepoClienteMongoDB) AddCliente(p *entities.Cliente) error {
	// Cria um contexto para a operação, geralmente com timeout
	ctx := context.Background()

	// O MongoDB usa o campo _id como identificador único.
	// Assumindo que seu struct entities.Cliente já está configurado
	// para mapear o campo ID (uuid.UUID) para _id no BSON.
	_, err := r.collection.InsertOne(ctx, p)
	if err != nil {
		// Você pode adicionar tratamento de erro para duplicidade de chave aqui
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

	// Consulta o MongoDB pelo campo _id (que deve ser o ID UUID)
	filter := bson.M{"_id": idUUID}
	err = r.collection.FindOne(ctx, filter).Decode(&cliente)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domainerr.ErrClienteNotFound
		}
		return nil, err
	}

	return &cliente, nil
}

// GetManyClienteByIDs - busca vários clientes por ID
func (r *RepoClienteMongoDB) GetManyClienteByIDs(ids []string) ([]*entities.Cliente, error) {
	ctx := context.Background()

	// Converte os IDs de string para uuid.UUID
	var uuidList []uuid.UUID
	for _, idStr := range ids {
		idUUID, err := uuid.Parse(idStr)
		if err != nil {
			return nil, domainerr.ErrClienteIDInvalid
		}
		uuidList = append(uuidList, idUUID)
	}

	// Filtro usando $in para buscar múltiplos IDs
	filter := bson.M{"_id": bson.M{"$in": uuidList}}

	// Busca os documentos
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var clientes []*entities.Cliente
	if err = cursor.All(ctx, &clientes); err != nil {
		return nil, err
	}

	if len(clientes) == 0 {
		return nil, domainerr.ErrClienteNotFoundMany
	}

	return clientes, nil
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

	// O filtro é o ID (que é o _id)
	filter := bson.M{"_id": idUUID}

	// O update usa $set para atualizar apenas os campos fornecidos no objeto
	// Assumindo que o struct 'p' contém os dados a serem atualizados,
	// você pode querer criar um DTO (Data Transfer Object) para atualização
	// ou usar o próprio 'p', mas garantir que os campos sejam marcados corretamente no struct
	// para BSON, ou construir o documento de atualização manualmente.
	// Neste exemplo, vou usar o documento completo, mas excluindo o _id da atualização.
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

	filter := bson.M{"_id": idUUID}

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
