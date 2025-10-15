# CPF-MANAGEMENT - API REST cpf-backend

API REST feita em Golang para gerenciar clientes(CPF/CNPJ). Veja mais detalhes no [`README.md`](./README.md)

## Pré-requisitos

- Sistema operacional Linux.
- Ferramenta `make` do Linux instalada na sua distribuição.
- Go versão 1.24 ou superior instalado.

## 1. Instalação

Clone o repositório e navegue até a pasta do projeto:
```bash
git clone https://github.com/valdinei-santos/cpf-backend.git
cd cpf-backend
```

Para baixar todas as dependências do projeto, use o comando:
```bash
go mod tidy
```

## 2. Configuração
Para rodar a API, você precisa criar um arquivo de variáveis de ambiente `.env`. Copie o arquivo de exemplo `env.example`, pois nesse momento você só precisa da definição da PORTA que a API vai rodar:
```bash
cp .env.exemplo .env
```

## 3. Execução

Para compilar a API, use o seguinte comando.
```bash
make build
```

Para rodar a API, você executa o arquivo executável `cpf-backend` que foi gerado:
```bash
./cpf-backend
```

## 4. Execução com docker

Para criar a imagem com o código fonto compilado - Usa o docker-compose.yaml do projeto.
```bash
make docker-build
```

Para subir o container com a imagem criada - Usa o docker-compose.yaml do projeto.
```bash
make docker-start
```

## 5. Testes nos endpoints
Com a API em execução você pode fazer testes básicos usuando sua ferramenta preferida.
Seguem alguns endpoints de exemplo.
- GET http://localhost:8889/ping
- GET http://localhost:8889/status
- POST http://localhost:8889/api/v1/cliente/ --> Body com JSON conforme documentação OpenAPI
- DELETE http://localhost:8889/api/v1/cliente/0d605862-91e8-11f0-9140-00155d6d572f
- GET http://localhost:8889/api/v1/cliente?page=1&limit=2
- GET http://localhost:8889/api/v1/cliente/0d605862-91e8-11f0-9140-00155d6d572f
- PUT http://localhost:8889/api/v1/cliente/0d605862-91e8-11f0-9140-00155d6d572f


## 6. Testes automatizados
Para rodar os testes unitários e de integração do projeto, siga os passos abaixo:

1. Navegue até o diretório do projeto. Caso a API esteja rodando você precisa parar ela com CTRL+C:
```bash
cd cpf-backend
```

2. Execute todos os testes de unidade:
```bash
make test-u
```

3. Para rodar um arquivo de teste específico, use o comando **go test "nome-do-arquivo"**, conforme abaixo:
```bash
go test modules/cliente/usecases/delete/usecase_test.go
```

4. Para rodar um caso de teste específico, use o comando **go test -run "nome-do-caso-de-teste"**, conforme abaixo:
```bash
go test -run "Deve retornar sucesso ao excluir um cliente" modules/cliente/usecases/delete/usecase_test.go
```

5. Para rodar o teste de integração de todos os endpoints:
```bash
go test -v -tags=integration ./cmd/api/routes
```

6. Para ver a cobertura dos testes na aplicação:
```bash
make cover
``` 

### Estrutura de Testes
O projeto inclui testes automatizados para os seguintes pacotes:

- **cmd/api/routes**: Faz os testes de integração de todos os endpoints.
- **modules/cliente/usecases/create**: Faz testes de unidade do usecase **create**
- **modules/cliente/usecases/delete**: Faz testes de unidade do usecase **delete**
- **modules/cliente/usecases/get**: Faz testes de unidade do usecase **get**
- **modules/cliente/usecases/getall**: Faz testes de unidade do usecase **getall**
- **modules/cliente/usecases/update**: Faz testes de unidade do usecase **update**


## 7. Documentação da API
O link para acessar a documentação está disponível no `README.md`, mas caso algum alteração seja feita no código da API e você precise recriar a documentação, o comando abaixo deverá ser executado:
```bash
make docs
```