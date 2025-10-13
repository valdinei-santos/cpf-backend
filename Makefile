# Variables
APP_NAME=cpf-backend
#VERSION=0.1.0


.PHONY: help run build mock test cover docs clean docker-build docker-start docker-stop

# Tasks
default: help

help: ## Exibe esta mensagem de ajuda
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
run: ## Roda o projeto
	@echo "Rodando o projeto..."
	@go run cmd/api/main.go
build: ## Cria o executável da aplicação
	@echo "Construindo o executável..."
	@go build -o $(APP_NAME) cmd/api/main.go
mock: ## Recria os mocks usados pelos tests da aplicação
	@echo "Recriando os mocks..."
	mockgen -source=internal/infra/logger/interfaces.go -destination=internal/infra/logger/mocks/mocks.go -package=mocks

	mockgen -source=internal/modules/cliente/infra/repository/interfaces.go -destination=internal/modules/cliente/infra/repository/mocks/mocks.go -package=mocks
	mockgen -source=internal/modules/cliente/usecases/create/interfaces.go -destination=internal/modules/cliente/usecases/create/mocks/mocks.go -package=mocks
	mockgen -source=internal/modules/cliente/usecases/delete/interfaces.go -destination=internal/modules/cliente/usecases/delete/mocks/mocks.go -package=mocks
	mockgen -source=internal/modules/cliente/usecases/get/interfaces.go -destination=internal/modules/cliente/usecases/get/mocks/mocks.go -package=mocks
	mockgen -source=internal/modules/cliente/usecases/getall/interfaces.go -destination=internal/modules/cliente/usecases/getall/mocks/mocks.go -package=mocks
	mockgen -source=internal/modules/cliente/usecases/update/interfaces.go -destination=internal/modules/cliente/usecases/update/mocks/mocks.go -package=mocks

	go mod tidy
test: ## Executa os test automatizados da aplicação
	@echo "Executando os testes automatizados..."
	@go test ./...
cover: ## Mostra a cobertura de testes da aplicação
	go test -v -cover ./...
lint: ## Roda a análise estática do código. Precisa ter a ferramenta golangci-lint instalada.
	golangci-lint run ./...
docs: ## Gera a documentação OpenAPI (Swagger) dos endpoints da aplicação
	@swag init -g ./cmd/api/main.go -o ./docs
	@echo "Documentação gerada na pasta docs"
clean: ## Remove o binário e arquivos de cache Go
	rm -f $(APP_NAME)
	go clean
	@echo "Limpeza completa: Executável $(APP_NAME) removido"

# ==========
# DOCKER
# ==========

# Variaveis Docker
CONTAINER_NAME=cont-${APP_NAME}
HOST_PORT=8888
NETWORK_NAME=net_backend
ENV_FILE=.env
REGISTRY_DNS=localhost
REGISTRY_PORT=5000
REGISTRY_USER=

get_version_command := read -p "Digite a versão da imagem (ex: 1.0.0): " VERSION; echo "$$VERSION"

docker-build: ## Cria a imagem docker com o executável da api. Solicita a TAG da imagem a ser gerada
	@VERSION=$$( $(get_version_command) ); \
	echo "Criando a imagem ${APP_NAME}:$$VERSION" ; \
	docker build -t ${APP_NAME}:$$VERSION .
docker-start: ## Inicia o container da aplicação com a imagem da TAG informada
	@VERSION=$$( $(get_version_command) ); \
	@echo "Subindo o container ${CONTAINER_NAME}..."; \
	echo "Iniciando o container ${CONTAINER_NAME} com imagem ${APP_NAME}:$$VERSION" ; \
	echo "docker run --rm --name ${CONTAINER_NAME} --network ${NETWORK_NAME} -p 0.0.0.0:${HOST_PORT}:8800 --env-file ${ENV_FILE} -d ${APP_NAME}:$$VERSION" ; \
	docker run --rm --name ${CONTAINER_NAME} --network ${NETWORK_NAME} -p 0.0.0.0:${HOST_PORT}:8800 --env-file ${ENV_FILE} -d ${APP_NAME}:$$VERSION
docker-stop: ## Para o container da aplicação
	@echo "Parando o container ${CONTAINER_NAME}..."
	docker stop ${CONTAINER_NAME}
docker-restart: docker-stop docker-start ## Executa o docker-stop e docker-start na sequencia
	@echo "Reiniciando container ${CONTAINER_NAME}..."
docker-rebuild: ## Cria a imagem docker, para o container atual e sobe com a nova imagem. Solicita a TAG da imagem
	@VERSION=$$( $(get_version_command) ); \
	echo "Criando a imagem ${APP_NAME}:$$VERSION" ; \
	docker build -t ${APP_NAME}:$$VERSION . ; \
	echo "Parando o container ${CONTAINER_NAME}..." ; \
	docker stop ${CONTAINER_NAME} ; \
	echo "Subindo o container ${CONTAINER_NAME}..."; \
	echo "Iniciando o container ${CONTAINER_NAME} com imagem ${APP_NAME}:$$VERSION" ; \
	echo "docker run --rm --name ${CONTAINER_NAME} --network ${NETWORK_NAME} -p 0.0.0.0:${HOST_PORT}:8800 --env-file ${ENV_FILE} -d ${APP_NAME}:$$VERSION" ; \
	docker run --rm --name ${CONTAINER_NAME} --network ${NETWORK_NAME} -p 0.0.0.0:${HOST_PORT}:8800 --env-file ${ENV_FILE} -d ${APP_NAME}:$$VERSION
docker-ps: ## Listar o container rodando
	-docker ps | grep ${CONTAINER_NAME}
docker-image: ## Para ver as imagens existentes
	docker images | grep $(APP_NAME)
docker-logs: ## Mostra o log do container
	docker logs ${CONTAINER_NAME}
docker-logs-f: ## Mostra o log do container com opção -f
	docker logs -f ${CONTAINER_NAME}
docker-exec: ## Entra dentro do container com ash
	docker exec -it ${CONTAINER_NAME} ash
#docker-exec-bash: ## Entra dentro do container com bash
#       docker exec -it ${CONTAINER_NAME} bash
#docker-push: ## Envia a imagem para o registry local
#        @echo "Não esqueça de fazer o \"make test\" antes de enviar a imagem para o Repository."
#        @VERSION=$$( $(get_version_command) ); \
#        echo "Forneça a senha do usuário do Registry (${REGISTRY_USER})" ; \
#        docker login ${REGISTRY_DNS}:${REGISTRY_PORT} -u ${REGISTRY_USER} ; \
#        docker tag ${APP_NAME}:$$VERSION ${REGISTRY_DNS}:${REGISTRY_PORT}/${APP_NAME}:$$VERSION ; \
#        echo "tag $$VERSION gerada..." ; \
#        docker push ${REGISTRY_DNS}:${REGISTRY_PORT}/${APP_NAME}:$$VERSION ; \
#        echo "push concluído..."