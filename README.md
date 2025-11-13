# Go Payments API

![Go Version](https://img.shields.io/badge/Go-1.25.4-00ADD8?style=flat&logo=go)
![License](https://img.shields.io/badge/License-MIT-green.svg)
![Architecture](https://img.shields.io/badge/Architecture-Clean%20Architecture-blue)

API de pagamentos desenvolvida em Go seguindo princípios de Clean Architecture, com suporte a PostgreSQL, Kafka e observabilidade completa com OpenTelemetry e Jaeger.

## 📋 Índice

- [Características](#-características)
- [Tecnologias](#-tecnologias)
- [Arquitetura](#-arquitetura)
- [Pré-requisitos](#-pré-requisitos)
- [Instalação](#-instalação)
- [Configuração](#-configuração)
- [Uso](#-uso)
- [Endpoints](#-endpoints)
- [Desenvolvimento](#-desenvolvimento)
- [Testes](#-testes)
- [Monitoramento](#-monitoramento)

## ✨ Características

- ✅ **Clean Architecture** - Separação clara entre camadas de domínio, aplicação e infraestrutura
- ✅ **Dependency Injection** - Usando Google Wire para injeção de dependências
- ✅ **Database Migrations** - Sistema automático de migrations ao iniciar a aplicação
- ✅ **Event-Driven** - Publicação de eventos no Kafka quando pagamentos são criados
- ✅ **Observabilidade** - Tracing distribuído com OpenTelemetry e Jaeger
- ✅ **API Documentation** - Swagger/OpenAPI automático
- ✅ **Docker Compose** - Infraestrutura completa containerizada
- ✅ **Validação** - Validação de entrada com go-playground/validator

## 🛠 Tecnologias

### Backend
- **[Go 1.25.4](https://go.dev/)** - Linguagem de programação
- **[Gin](https://gin-gonic.com/)** - Framework HTTP
- **[Google Wire](https://github.com/google/wire)** - Dependency injection
- **[Logrus](https://github.com/sirupsen/logrus)** - Logging estruturado

### Banco de Dados
- **[PostgreSQL 15](https://www.postgresql.org/)** - Banco de dados relacional
- **[database/sql](https://pkg.go.dev/database/sql)** - Driver nativo do Go

### Mensageria
- **[Apache Kafka](https://kafka.apache.org/)** - Sistema de mensageria
- **[Zookeeper](https://zookeeper.apache.org/)** - Coordenação do Kafka
- **[Kafka UI](https://github.com/provectus/kafka-ui)** - Interface web para Kafka
- **[segmentio/kafka-go](https://github.com/segmentio/kafka-go)** - Cliente Kafka para Go

### Observabilidade
- **[OpenTelemetry](https://opentelemetry.io/)** - Instrumentação de observabilidade
- **[Jaeger](https://www.jaegertracing.io/)** - Distributed tracing
- **[OTEL Collector](https://opentelemetry.io/docs/collector/)** - Coletor de telemetria

### Documentação
- **[Swagger/OpenAPI](https://swagger.io/)** - Documentação interativa da API
- **[swaggo](https://github.com/swaggo/swag)** - Gerador de docs Swagger para Go

### DevOps
- **[Docker](https://www.docker.com/)** - Containerização
- **[Docker Compose](https://docs.docker.com/compose/)** - Orquestração de containers

## 🏗 Arquitetura

```
go-payments-api/
├── cmd/
│   └── server/           # Entrypoint da aplicação
├── internal/
│   ├── domain/           # Camada de Domínio (Entidades, DTOs)
│   │   ├── entity/       # Entidades de negócio
│   │   └── dto/          # Data Transfer Objects
│   ├── application/      # Camada de Aplicação (Use Cases)
│   │   ├── usecase/      # Casos de uso (regras de negócio)
│   │   └── gateway/      # Interfaces de repositório
│   ├── infrastructure/   # Camada de Infraestrutura
│   │   ├── api/          # Handlers HTTP
│   │   │   ├── handler/  # Handlers de rotas
│   │   │   └── middleware/ # Middlewares
│   │   ├── database/     # Implementações de banco de dados
│   │   │   └── postgres/ # Repository do PostgreSQL
│   │   └── messaging/    # Mensageria
│   │       └── kafka/    # Publisher Kafka
│   └── settings/         # Configurações da aplicação
├── pkg/                  # Pacotes reutilizáveis
│   ├── api/              # Abstrações de API
│   ├── base/             # Interfaces base (UseCase, Repository)
│   ├── errors/           # Tratamento de erros
│   ├── log/              # Logging
│   ├── metrics/          # Métricas e tracing
│   └── validator/        # Validação
├── di/                   # Dependency Injection (Wire)
├── scripts/              # Scripts utilitários
│   └── migrations/       # Migrations SQL
├── docs/                 # Documentação Swagger
└── test/                 # Utilitários de teste
```

### Fluxo de uma Requisição

```
HTTP Request
    ↓
[Handler] → Valida entrada
    ↓
[Use Case] → Lógica de negócio
    ↓
[Repository] → Persiste no PostgreSQL
    ↓
[Publisher] → Publica evento no Kafka
    ↓
HTTP Response
```

## 📦 Pré-requisitos

- [Go 1.25.4+](https://go.dev/dl/)
- [Docker](https://www.docker.com/get-started)
- [Docker Compose](https://docs.docker.com/compose/install/)
- [Make](https://www.gnu.org/software/make/) (opcional, mas recomendado)

## 🚀 Instalação

### 1. Clone o repositório

```bash
git clone https://github.com/seu-usuario/go-payments-api.git
cd go-payments-api
```

### 2. Instale as dependências

```bash
go mod download
```

### 3. Configure as variáveis de ambiente

```bash
cp .env.example .env
```

### 4. Suba a infraestrutura (Docker)

```bash
docker-compose up -d
```

Aguarde ~30 segundos para os serviços iniciarem completamente.

### 5. Gere o código Wire e Swagger

```bash
make wire
make docs
```

### 6. Execute a aplicação

```bash
go run cmd/server/main.go
```

A API estará disponível em: **http://localhost:8080**

## ⚙️ Configuração

As configurações são feitas através de variáveis de ambiente no arquivo `.env`:

```env
# Ambiente
ENVIRONMENT=local

# HTTP Server
HTTP_SERVER_PORT=:8080
HTTP_SERVER_READ_TIMEOUT=15s
HTTP_SERVER_WRITE_TIMEOUT=15s

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=payments_user
DB_PASSWORD=payments_pass
DB_NAME=payments

# Kafka (use porta 29092 quando rodar FORA do Docker)
KAFKA_BROKERS=localhost:29092

# Observabilidade
OTEL_SERVICE_NAME=go-payments-api
OTEL_EXPORTER_OTLP_ENDPOINT=http://localhost:4317
```

## 📖 Uso

### Criar um Pagamento

```bash
curl -X POST http://localhost:8080/v1/payments/payments \
  -H "Content-Type: application/json" \
  -d '{
    "amount": 150.75,
    "method": "PIX"
  }'
```

**Resposta:**

```json
{
  "id": 1,
  "amount": 150.75,
  "method": "PIX",
  "status": "CREATED",
  "created_at": "2024-11-13T10:30:00Z"
}
```

### Verificar Health Check

```bash
curl http://localhost:8080/v1/payments/health
```

**Resposta:**

```json
{
  "ok": true
}
```

## 📚 Endpoints

| Método | Endpoint | Descrição |
|--------|----------|-----------|
| `GET` | `/v1/payments/health` | Health check da aplicação |
| `POST` | `/v1/payments/payments` | Criar novo pagamento |
| `GET` | `/docs/payments` | Documentação Swagger |

### Documentação Interativa

Acesse a documentação Swagger em: **http://localhost:8080/docs/payments**

## 🔧 Desenvolvimento

### Comandos Make

```bash
# Gerar código Wire (dependency injection)
make wire

# Gerar documentação Swagger
make docs

# Criar tópico Kafka manualmente
make kafka-create-topic

# Consumir eventos Kafka
make kafka-consume

# Listar tópicos Kafka
make kafka-topics
```

### Adicionar Nova Migration

1. Crie um arquivo SQL em `scripts/migrations/` com prefixo numérico:

```bash
# Exemplo: scripts/migrations/002_add_user_id.sql
ALTER TABLE payments ADD COLUMN user_id BIGINT;
CREATE INDEX idx_payments_user_id ON payments(user_id);
```

2. Reinicie a aplicação - a migration será aplicada automaticamente

### Estrutura de um Use Case

```go
type CreatePayment = base.UseCase[dto.CreatePaymentInput, *dto.CreatePaymentOutput]

type CreatePaymentImplementation struct {
    repository repository.PaymentRepository
    publisher  kafka.Publisher
}

func (uc *CreatePaymentImplementation) Execute(ctx context.Context, input dto.CreatePaymentInput) (*dto.CreatePaymentOutput, error) {
    // 1. Validação
    // 2. Lógica de negócio
    // 3. Persistência
    // 4. Publicação de evento
    // 5. Retorno
}
```

### Adicionar Novo Handler

1. Crie o handler em `internal/infrastructure/api/handler/`:

```go
type MyHandler struct {
    UseCase   usecase.MyUseCase
    Presenter api.Presenter
}

func (h *MyHandler) Handlefunc(ctx *gin.Context) {
    return func(ctx *gin.Context) {
        // implementação
    }
}
```

2. Registre em `di/inject_api_handlers.go`:

```go
var apiHandlersSet = wire.NewSet(
    // ... handlers existentes
    wire.Struct(new(handler.MyHandler), "*"),
)
```

3. Adicione a rota em `internal/infrastructure/api/routes.go`:

```go
base.POST("/my-route", a.MyHandler.Handle())
```

4. Regenere o Wire:

```bash
make wire
```

## 🧪 Testes

```bash
# Executar todos os testes
go test ./...

# Executar testes com cobertura
go test -cover ./...

# Executar testes de um pacote específico
go test ./internal/application/usecase/...

# Executar testes com modo verbose
go test -v ./...
```

### Estrutura de Teste

```go
func TestCreatePaymentUseCase(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockRepo := repository.NewMockPaymentRepository(ctrl)
    mockPublisher := kafka.NewMockPublisher(ctrl)

    uc := usecase.NewCreatePaymentUseCase(mockRepo, mockPublisher)

    // Configurar mocks e assertions
}
```

## 📊 Monitoramento

### Interfaces Web Disponíveis

| Serviço | URL | Descrição |
|---------|-----|-----------|
| **API** | http://localhost:8080 | Aplicação principal |
| **Swagger** | http://localhost:8080/docs/payments | Documentação interativa |
| **Jaeger UI** | http://localhost:16686 | Distributed tracing |
| **Kafka UI** | http://localhost:8081 | Interface do Kafka |

### Verificar Eventos Kafka

#### Via Terminal:

```bash
docker exec -it go-payments-kafka kafka-console-consumer \
  --bootstrap-server kafka:9092 \
  --topic payment.events \
  --from-beginning \
  --property print.key=true \
  --property key.separator=": "
```

#### Via Kafka UI:

1. Acesse http://localhost:8081
2. Navegue até **Topics → payment.events → Messages**

### Visualizar Traces

1. Acesse http://localhost:16686 (Jaeger UI)
2. Selecione o serviço **go-payments-api**
3. Clique em "Find Traces"
4. Visualize o trace completo da requisição

### Logs da Aplicação

A aplicação usa logging estruturado com níveis:

```
[2024-11-13T10:30:00] [INFO] [usecase.Execute()] [create_payment.go:45] 🔵 Starting payment creation - Amount: 150.75, Method: PIX
[2024-11-13T10:30:00] [INFO] [usecase.Execute()] [create_payment.go:60] 💾 Saving payment to database...
[2024-11-13T10:30:00] [INFO] [usecase.Execute()] [create_payment.go:68] ✅ Payment saved to database with ID: 1
[2024-11-13T10:30:00] [INFO] [usecase.Execute()] [create_payment.go:82] 📤 Publishing event to Kafka - Topic: payment.events, Key: 1
[2024-11-13T10:30:00] [INFO] [usecase.Execute()] [create_payment.go:88] ✅ Event published successfully to Kafka
```

## 🗂 Banco de Dados

### Schema

```sql
CREATE TABLE payments (
    id BIGSERIAL PRIMARY KEY,
    amount DECIMAL(10, 2) NOT NULL,
    method VARCHAR(20) NOT NULL,
    status VARCHAR(20) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_payments_status ON payments(status);
CREATE INDEX idx_payments_created_at ON payments(created_at);
```

### Conectar ao PostgreSQL

```bash
docker exec -it go-payments-postgres psql -U payments_user -d payments
```

### Verificar Migrations Aplicadas

```sql
SELECT * FROM schema_migrations ORDER BY version;
```

## 🐛 Troubleshooting

### Erro: "password authentication failed"

Verifique se as credenciais no `.env` correspondem ao `docker-compose.yml`:

```bash
docker-compose down -v
docker-compose up -d
```

### Erro: "Unknown Topic Or Partition"

Crie o tópico Kafka manualmente:

```bash
docker exec -it go-payments-kafka kafka-topics \
  --bootstrap-server kafka:9092 \
  --create \
  --topic payment.events \
  --partitions 3 \
  --replication-factor 1
```

### Kafka não conecta

Aguarde ~30 segundos após `docker-compose up` para os serviços iniciarem.

Verifique se os containers estão rodando:

```bash
docker-compose ps
```

### Wire não gera código

Instale o Wire globalmente:

```bash
go install github.com/google/wire/cmd/wire@latest
```

## 📝 Licença

Este projeto está sob a licença MIT. Veja o arquivo [LICENSE](LICENSE) para mais detalhes.

## 👥 Contribuindo

1. Fork o projeto
2. Crie uma branch para sua feature (`git checkout -b feature/AmazingFeature`)
3. Commit suas mudanças (`git commit -m 'Add some AmazingFeature'`)
4. Push para a branch (`git push origin feature/AmazingFeature`)
5. Abra um Pull Request

## 📧 Contato

Diego - [@luckydied](https://x.com/luckydied)

Link do Projeto: [https://github.com/didiegovieira/go-payments-api](https://github.com/didiegovieira/go-payments-api)

---

⭐️ Se este projeto te ajudou, considere dar uma estrela!