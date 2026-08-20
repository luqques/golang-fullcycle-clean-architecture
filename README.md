# Clean Architecture: Listagem de Orders

Projeto exclusivo para o desafio de listagem de Orders usando Go, Clean Architecture, REST, gRPC, GraphQL, MySQL, Docker e Docker Compose.

## Execução

```bash
docker compose up --build
```

Esse comando único sobe o MySQL, aguarda o banco ficar saudável, aplica as migrações automaticamente e inicia a aplicação.

## Portas

| Serviço | Porta | Endpoint / Service |
|---|---:|---|
| REST | `8000` | `GET /order` e `POST /order` |
| GraphQL | `8080` | `POST /query` |
| gRPC | `50051` | `orders.OrderService/ListOrders` |
| MySQL | `3306` | Banco `orders` |

## Testes rápidos

### REST

Criar uma order:

```bash
curl -X POST http://localhost:8000/order \
  -H "Content-Type: application/json" \
  -d '{"price":100.50,"tax":10.25}'
```

Listar orders:

```bash
curl http://localhost:8000/order
```

### GraphQL

Criar uma order:

```bash
curl -X POST http://localhost:8080/query \
  -H "Content-Type: application/json" \
  -d '{"query":"mutation { createOrder(input: { price: 200.00, tax: 20.00 }) { id price tax finalPrice } }"}'
```

Listar orders:

```bash
curl -X POST http://localhost:8080/query \
  -H "Content-Type: application/json" \
  -d '{"query":"query { listOrders { id price tax finalPrice } }"}'
```

### gRPC

Com `grpcurl` instalado:

```bash
grpcurl -plaintext localhost:50051 orders.OrderService/ListOrders
```

## Estrutura

```text
cmd/orders/main.go               # bootstrap da aplicação
internal/entity                  # entidades de domínio
internal/usecase                 # casos de uso
internal/infra/database          # repositório MySQL
internal/infra/web               # handlers REST
internal/infra/grpc              # service gRPC
internal/infra/graphql           # schema e resolvers GraphQL
migrations                       # migrações SQL
proto                            # contrato protobuf
api.http                         # requisições prontas
Dockerfile                       # build da aplicação Go
docker-compose.yaml              # app + MySQL
```

## Observações de arquitetura

O caso de uso `ListOrdersUseCase` é único e fica em `internal/usecase/list_orders.go`.

As três interfaces de entrada apenas adaptam transporte para aplicação:

- REST chama `ListOrdersUseCase.Execute()` no endpoint `GET /order`.
- gRPC chama `ListOrdersUseCase.Execute()` no service `OrderService/ListOrders`.
- GraphQL chama `ListOrdersUseCase.Execute()` na query `listOrders`.

A regra de listagem não fica duplicada nos handlers.
