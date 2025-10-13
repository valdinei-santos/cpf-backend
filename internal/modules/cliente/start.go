package cliente

import (
	"github.com/gin-gonic/gin"
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

// StartCreate - Metodo onde instanciamentos as dependencias e chamamos o controller
func StartCreate(log logger.ILogger, ctx *gin.Context, db *mongo.Database) {
	log.Debug("Entrou product.StartCreate")
	repo := repository.NewRepoClienteMongoDB(db.Client(), "cpf_management", "cliente", log)
	u := create.NewUseCase(repo, log)
	controller.Create(log, ctx, u)
}

// StartDelete - Metodo onde instanciamentos as dependencias e chamamos o controller
func StartDelete(log logger.ILogger, ctx *gin.Context, db *mongo.Database) {
	log.Debug("Entrou product.StartDelete")
	repo := repository.NewRepoClienteMongoDB(db.Client(), "cpf_management", "cliente", log)
	u := delete.NewUseCase(repo, log)
	controller.Delete(log, ctx, u)
}

// StartGet - Metodo onde instanciamentos as dependencias e chamamos o controller
func StartGet(log logger.ILogger, ctx *gin.Context, db *mongo.Database) {
	log.Debug("Entrou product.StartGet")
	repo := repository.NewRepoClienteMongoDB(db.Client(), "cpf_management", "cliente", log)
	u := get.NewUseCase(repo, log)
	controller.Get(log, ctx, u)
}

// StartGetAll - Metodo onde instanciamentos as dependencias e chamamos o controller
func StartGetAll(log logger.ILogger, ctx *gin.Context, db *mongo.Database) {
	log.Debug("Entrou usecases.StartGetAll")
	repo := repository.NewRepoClienteMongoDB(db.Client(), "cpf_management", "cliente", log)
	u := getall.NewUseCase(repo, log)
	controller.GetAll(log, ctx, u)
}

// StartUpdate - Metodo onde instanciamentos as dependencias e chamamos o controller
func StartUpdate(log logger.ILogger, ctx *gin.Context, db *mongo.Database) {
	log.Debug("Entrou product.StartUpdate")
	repo := repository.NewRepoClienteMongoDB(db.Client(), "cpf_management", "cliente", log)
	u := update.NewUseCase(repo, log)
	controller.Update(log, ctx, u)
}
