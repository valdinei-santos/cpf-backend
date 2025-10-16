package controller

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/valdinei-santos/cpf-backend/internal/domain/globalerr"
	"github.com/valdinei-santos/cpf-backend/internal/infra/logger"
	"github.com/valdinei-santos/cpf-backend/internal/modules/cliente/domain/domainerr"
	"github.com/valdinei-santos/cpf-backend/internal/modules/cliente/infra/dto"
	"github.com/valdinei-santos/cpf-backend/internal/modules/cliente/usecases/create"
	"github.com/valdinei-santos/cpf-backend/internal/modules/cliente/usecases/delete"
	"github.com/valdinei-santos/cpf-backend/internal/modules/cliente/usecases/get"
	"github.com/valdinei-santos/cpf-backend/internal/modules/cliente/usecases/getall"
	"github.com/valdinei-santos/cpf-backend/internal/modules/cliente/usecases/update"
)

// ClienteController orquestra todas as ações da entidade Cliente.
type ClienteController struct {
	log      logger.ILogger
	createUC create.IUsecase
	deleteUC delete.IUsecase
	getUC    get.IUsecase
	getAllUC getall.IUsecase
	updateUC update.IUsecase
}

// NewClienteController é o construtor que injeta todas as dependências.
func NewClienteController(
	log logger.ILogger,
	c create.IUsecase,
	d delete.IUsecase,
	g get.IUsecase,
	ga getall.IUsecase,
	u update.IUsecase,
) *ClienteController {
	return &ClienteController{
		log:      log,
		createUC: c,
		deleteUC: d,
		getUC:    g,
		getAllUC: ga,
		updateUC: u,
	}
}

// Handler específico de Criação
func (c *ClienteController) Create(ctx *gin.Context) {
	c.log.Debug("Entrou controller.Create")
	var input *dto.Request
	err := json.NewDecoder(ctx.Request.Body).Decode(&input)
	if err != nil {
		outputError(c.log, ctx, globalerr.ErrBadRequest, "Create/json.Decode")
		return
	}
	resp, err := c.createUC.Execute(input)
	if err != nil {
		outputError(c.log, ctx, err, "Create/usecase.Execute")
		return
	}

	ctx.JSON(http.StatusCreated, resp)
	c.log.Info("### Finished OK", "status_code", http.StatusCreated)
}

// Handler específico de Deleção
func (c *ClienteController) Delete(ctx *gin.Context) {
	c.log.Debug("Entrou controller.Delete")
	id, err := getIdParam(c.log, ctx)
	if err != nil {
		return
	}
	c.log.Debug("ID: " + id)
	err = c.deleteUC.Execute(id)
	if err != nil {
		outputError(c.log, ctx, err, "Delete/usecase.Execute")
		return
	}

	// Retorna a resposta padrão
	result := &dto.OutputDefault{}

	ctx.JSON(http.StatusNoContent, result)
	c.log.Info("### Finished OK", "status_code", http.StatusNoContent)
}

// Handler específico de Obtenção por ID
func (c *ClienteController) Get(ctx *gin.Context) {
	c.log.Debug("Entrou controller.Get")
	id, err := getIdParam(c.log, ctx)
	if err != nil {
		outputError(c.log, ctx, globalerr.ErrBadRequest, "Get/getIdParam")
		return
	}
	c.log.Debug("ID: " + id)
	resp, err := c.getUC.Execute(id)
	if err != nil {
		outputError(c.log, ctx, err, "Get/usecase.Execute")
		return
	}
	ctx.JSON(http.StatusOK, resp)
	c.log.Info("### Finished OK", "status_code", http.StatusOK)
}

// Handler específico de Obtenção de todos os registros (com paginação)
func (c *ClienteController) GetAll(ctx *gin.Context) {
	c.log.Debug("Entrou controller.GetAll")

	// Pega os parâmetros de paginação (page) da query string
	page, err := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		outputError(c.log, ctx, globalerr.ErrBadRequest, "GetAll/parse_size")
		return
	}

	// Pega os parâmetros de paginação (size) da query string
	size, err := strconv.Atoi(ctx.DefaultQuery("size", "10"))
	if err != nil || size < 1 {
		outputError(c.log, ctx, globalerr.ErrBadRequest, "GetAll/parse_size")
		return
	}

	resp, err := c.getAllUC.Execute(int64(page), int64(size))
	if err != nil {
		outputError(c.log, ctx, err, "GetAll/usecase.Execute")
		return
	}
	ctx.JSON(http.StatusOK, resp)
	c.log.Info("### Finished OK", "status_code", http.StatusOK)
}

// Handler específico de Atualização
func (c *ClienteController) Update(ctx *gin.Context) {
	c.log.Debug("Entrou controller.Update")
	id, err := getIdParam(c.log, ctx)
	if err != nil {
		outputError(c.log, ctx, globalerr.ErrBadRequest, "Update/getIdParam")
		return
	}
	c.log.Debug("ID: " + id)
	var input *dto.Request
	err = json.NewDecoder(ctx.Request.Body).Decode(&input)
	if err != nil {
		outputError(c.log, ctx, err, "Update/json.NewDecoder")
		return
	}
	resp, err := c.updateUC.Execute(id, input)
	if err != nil {
		outputError(c.log, ctx, err, "Update/usecase.Execute")
		return
	}
	ctx.JSON(http.StatusOK, resp)
	c.log.Info("### Finished OK", "status_code", http.StatusOK)
}

func getIdParam(log logger.ILogger, ctx *gin.Context) (string, error) {
	idParam := ctx.Param("id")
	if idParam == "" {
		log.Error(domainerr.ErrClienteIDInvalid.Error(), "mtd", "getIdParam")
		return "", domainerr.ErrClienteIDInvalid
	}
	return idParam, nil
}
func outputError(log logger.ILogger, ctx *gin.Context, err error, method string) {
	log.Error(err.Error(), "mtd", method)
	dataJErro := dto.OutputDefault{}
	var errHttp int
	switch err {
	case globalerr.ErrDuplicatekey:
		errHttp = http.StatusConflict
		dataJErro.Title = globalerr.ErrHttp409.Error()
		dataJErro.Detail = domainerr.ErrDuplicatekey.Error()
	case globalerr.ErrNotFound, domainerr.ErrClienteNotFound:
		errHttp = http.StatusNotFound
		dataJErro.Title = globalerr.ErrHttp404.Error()
		dataJErro.Detail = domainerr.ErrClienteNotFound.Error()
	case globalerr.ErrSaveInDatabase:
		errHttp = http.StatusInternalServerError
		dataJErro.Title = globalerr.ErrHttp500.Error()
		dataJErro.Detail = globalerr.ErrSaveInDatabase.Error()
	case globalerr.ErrBadRequest:
		errHttp = http.StatusBadRequest
		dataJErro.Title = globalerr.ErrHttp404.Error()
		dataJErro.Detail = globalerr.ErrBadRequest.Error()
	case domainerr.ErrClienteIDInvalid:
		errHttp = http.StatusBadRequest
		dataJErro.Title = globalerr.ErrHttp400.Error()
		dataJErro.Detail = globalerr.ErrBadRequest.Error()
	default:
		errHttp = http.StatusInternalServerError
		dataJErro.Title = globalerr.ErrHttp500.Error()
		dataJErro.Detail = globalerr.ErrSaveInDatabase.Error()
	}
	ctx.JSON(errHttp, dataJErro)
	log.Info("### Finished ERROR", "status_code", errHttp)
}
