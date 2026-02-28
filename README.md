# Labs Auction - GoExpert

Sistema de leilões em Go com fechamento automático. Permite criar leilões, realizar lances e encerra automaticamente os leilões após o tempo configurado.

## Pré-requisitos

- Go 1.20+
- Docker e Docker Compose (para execução via containers)
- MongoDB (quando rodar localmente sem Docker)

## Como Rodar

### Via Docker Compose (recomendado)

```bash
docker-compose up --build
```

A API ficará disponível em `http://localhost:8080`.

### Localmente

1. Configure o arquivo `cmd/auction/.env` (veja seção de variáveis de ambiente).
2. Inicie o MongoDB (ou use um banco remoto).
3. Execute:

```bash
go run cmd/auction/main.go
```

## Variáveis de Ambiente

Configure no arquivo `cmd/auction/.env`:

| Variável | Descrição | Exemplo |
|----------|-----------|---------|
| `AUCTION_DURATION` | Duração do leilão até fechamento automático | `20s`, `5m` |
| `AUCTION_INTERVAL` | Usado na validação de bids e como fallback de `AUCTION_DURATION` | `20s` |
| `BATCH_INSERT_INTERVAL` | Intervalo para batch de inserção de bids | `20s` |
| `MAX_BATCH_SIZE` | Tamanho máximo do batch de bids | `4` |
| `MONGODB_URL` | URL de conexão do MongoDB | `mongodb://admin:admin@mongodb:27017/auctions?authSource=admin` |
| `MONGODB_DB` | Nome do banco de dados | `auctions` |
| `MONGO_INITDB_ROOT_USERNAME` | Usuário root do MongoDB (para Docker) | `admin` |
| `MONGO_INITDB_ROOT_PASSWORD` | Senha root do MongoDB (para Docker) | `admin` |

### Configuração de Tempo

- `AUCTION_DURATION` define por quanto tempo o leilão permanece ativo até ser fechado automaticamente.
- Recomenda-se manter `AUCTION_DURATION` e `AUCTION_INTERVAL` com o mesmo valor para consistência na validação de lances.
- Valores aceitos: `s` (segundos), `m` (minutos), `h` (horas). Ex.: `30s`, `5m`, `1h`.

### Exemplo de `.env`

```env
BATCH_INSERT_INTERVAL=20s
MAX_BATCH_SIZE=4
AUCTION_INTERVAL=20s
AUCTION_DURATION=20s

MONGO_INITDB_ROOT_USERNAME=admin
MONGO_INITDB_ROOT_PASSWORD=admin
MONGODB_URL=mongodb://admin:admin@mongodb:27017/auctions?authSource=admin
MONGODB_DB=auctions
```

Para rodar localmente (MongoDB na máquina):

```env
MONGODB_URL=mongodb://localhost:27017
MONGODB_DB=auctions
```

## Testes

Execute todos os testes:

```bash
go test ./...
```

Teste de fechamento automático (requer MongoDB rodando):

```bash
go test ./internal/infra/database/auction/... -v -run TestCreateAuction_AutoCloseAfterDuration
```

Com Docker Compose, suba o MongoDB antes:

```bash
docker-compose up -d mongodb
go test ./internal/infra/database/auction/... -v -run TestCreateAuction_AutoCloseAfterDuration
```

## Endpoints

- `POST /auction` - Criar leilão
- `GET /auction` - Listar leilões
- `GET /auction/:auctionId` - Buscar leilão por ID
- `GET /auction/winner/:auctionId` - Buscar vencedor do leilão
- `POST /bid` - Registrar lance
- `GET /bid/:auctionId` - Listar lances de um leilão
- `GET /user/:userId` - Buscar usuário
