package cliente

import (
	"github.com/valdinei-santos/cpf-backend/internal/infra/logger"
	"github.com/valdinei-santos/cpf-backend/internal/modules/cliente/infra/controller"
	"github.com/valdinei-santos/cpf-backend/internal/modules/cliente/infra/repository"
	"github.com/valdinei-santos/cpf-backend/internal/modules/cliente/usecases/create"
	"github.com/valdinei-santos/cpf-backend/internal/modules/cliente/usecases/delete"
	"github.com/valdinei-santos/cpf-backend/internal/modules/cliente/usecases/get"
	"github.com/valdinei-santos/cpf-backend/internal/modules/cliente/usecases/getall"
	"github.com/valdinei-santos/cpf-backend/internal/modules/cliente/usecases/update"
	"go.mongodb.org/mongo-driver/mongo"
)

// ModuleCliente é uma estrutura que contém os controladores (handlers). Mudei para usar só um controller por módulo
// para as rotas. Ela será criada uma única vez.
type ModuleCliente struct {
	Controller *controller.ClienteController
}

// NewModuleCliente - Inicializa TODAS as dependências do módulo uma única vez.
func NewModuleCliente(log logger.ILogger, db *mongo.Database) *ModuleCliente {
	log.Info("Inicializando Módulo Cliente...")

	repo := repository.NewRepoClienteMongoDB(db, "cliente", log)
	createUC := create.NewUseCase(repo, log)
	deleteUC := delete.NewUseCase(repo, log)
	getUC := get.NewUseCase(repo, log)
	getAllUC := getall.NewUseCase(repo, log)
	updateUC := update.NewUseCase(repo, log)

	clienteController := controller.NewClienteController(
		log,
		createUC,
		deleteUC,
		getUC,
		getAllUC,
		updateUC,
	)

	return &ModuleCliente{
		Controller: clienteController,
	}
}
