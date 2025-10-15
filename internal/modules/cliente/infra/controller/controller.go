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

// Create - Controlador para criar um cliente
func Create(log logger.ILogger, ctx *gin.Context, useCase create.IUsecase) {
	log.Debug("Entrou controller.Create")
	var input *dto.Request
	err := json.NewDecoder(ctx.Request.Body).Decode(&input)
	if err != nil {
		outputError(log, ctx, globalerr.ErrBadRequest, "Create/json.Decode")
		return
	}
	resp, err := useCase.Execute(input)
	if err != nil {
		outputError(log, ctx, err, "Create/usecase.Execute")
		return
	}

	ctx.JSON(http.StatusCreated, resp)
	log.Info("### Finished OK", "status_code", http.StatusCreated)
}

// Delete - Controlador para deletar um cliente
func Delete(log logger.ILogger, ctx *gin.Context, useCase delete.IUsecase) {
	log.Debug("Entrou controller.Delete")
	id, err := getIdParam(log, ctx)
	if err != nil {
		return
	}
	log.Debug("ID: " + id)
	err = useCase.Execute(id)
	if err != nil {
		outputError(log, ctx, err, "Delete/usecase.Execute")
		return
	}

	// Retorna a resposta padrão
	result := &dto.OutputDefault{}

	ctx.JSON(http.StatusNoContent, result)
	log.Info("### Finished OK", "status_code", http.StatusNoContent)
}

// Get - Controlador para obter um cliente por ID
func Get(log logger.ILogger, ctx *gin.Context, useCase get.IUsecase) {
	log.Debug("Entrou controller.Get")
	id, err := getIdParam(log, ctx)
	if err != nil {
		outputError(log, ctx, globalerr.ErrBadRequest, "Get/getIdParam")
		return
	}
	log.Debug("ID: " + id)
	resp, err := useCase.Execute(id)
	if err != nil {
		outputError(log, ctx, err, "Get/usecase.Execute")
		return
	}
	ctx.JSON(http.StatusOK, resp)
	log.Info("### Finished OK", "status_code", http.StatusOK)
}

// GetAll - Controlador para obter todos os clientes
func GetAll(log logger.ILogger, ctx *gin.Context, useCase getall.IUsecase) {
	log.Debug("Entrou controller.GetAll")

	// Pega os parâmetros de paginação (page) da query string
	page, err := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		outputError(log, ctx, globalerr.ErrBadRequest, "GetAll/parse_size")
		return
	}

	// Pega os parâmetros de paginação (size) da query string
	size, err := strconv.Atoi(ctx.DefaultQuery("size", "10"))
	if err != nil || size < 1 {
		outputError(log, ctx, globalerr.ErrBadRequest, "GetAll/parse_size")
		return
	}

	resp, err := useCase.Execute(int64(page), int64(size))
	if err != nil {
		outputError(log, ctx, err, "GetAll/usecase.Execute")
		return
	}
	ctx.JSON(http.StatusOK, resp)
	log.Info("### Finished OK", "status_code", http.StatusOK)
}

// Update - Controlador para alterar um cliente pelo ID
func Update(log logger.ILogger, ctx *gin.Context, useCase update.IUsecase) {
	log.Debug("Entrou controller.Update")
	id, err := getIdParam(log, ctx)
	if err != nil {
		outputError(log, ctx, globalerr.ErrBadRequest, "Update/getIdParam")
		return
	}
	log.Debug("ID: " + id)
	var input *dto.Request
	err = json.NewDecoder(ctx.Request.Body).Decode(&input)
	if err != nil {
		outputError(log, ctx, err, "Update/json.NewDecoder")
		return
	}
	resp, err := useCase.Execute(id, input)
	if err != nil {
		outputError(log, ctx, err, "Update/usecase.Execute")
		return
	}
	ctx.JSON(http.StatusOK, resp)
	log.Info("### Finished OK", "status_code", http.StatusOK)
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
