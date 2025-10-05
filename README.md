# Better-Mem

Sistema de gerenciamento de memórias para LLMs com classificação automática em memórias de curto e longo prazo.

## Estrutura do Projeto

- `cmd/api` - API REST principal
- `cmd/worker` - Worker para processamento assíncrono de mensagens
- `cmd/demo` - Aplicação CLI de demonstração
- `internal/` - Código interno da aplicação
- `inference/` - Serviço de inferência ML (Python)
- `demo/` - Pacote da aplicação demo

## Quick Start

### 1. Iniciar os serviços

```bash
docker-compose up -d
```

### 2. Build e executar a API

```bash
./scripts/build.sh
./target/api.o
```

### 3. Build e executar o Worker

```bash
./scripts/build.sh
./target/worker.o
```

### 4. Testar com a Demo CLI

```bash
./demo/build_demo.sh
./demo/demo.o
```

Veja mais detalhes em [demo/README.md](demo/README.md)

## Documentação da API

Acesse a documentação Swagger em: http://localhost:8080/swagger/index.html

## Endpoints Principais

- `POST /api/v1/chat` - Criar um novo chat
- `POST /api/v1/message` - Enviar mensagem para processamento
- `POST /api/v1/memory/chat/{chat_id}/fetch` - Buscar memórias relevantes
- `GET /api/v1/memory/short-term/chat/{chat_id}` - Listar memórias de curto prazo
- `GET /api/v1/memory/long-term/chat/{chat_id}` - Listar memórias de longo prazo

## Tecnologias

- **Go** - API e Worker
- **Python** - Serviço de inferência ML
- **MongoDB** - Banco de dados principal
- **Redis** - Fila de tarefas (Asynq)
- **Qdrant** - Banco de dados vetorial
- **gRPC** - Comunicação entre serviços
