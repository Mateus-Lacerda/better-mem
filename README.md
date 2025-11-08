# Better-Mem

Memory management system for LLMs with automatic classification into short-term and long-term memories.

## Project Structure

- `cmd/api` - Main REST API
- `cmd/worker` - Worker for asynchronous message processing
- `internal/` - Internal application code
- `inference/` - ML inference service (Python)
- `demo/` - Demo application package

## Quick Start

### 1. Start services

```bash
docker-compose up -d
```

### 2. Build and run the API

```bash
./scripts/build.sh
./bin/api
```

### 3. Build and run the Worker

```bash
./scripts/build.sh
./bin/worker
```

### 4. Test with CLI Demo

```bash
./demo/build_demo.sh
./demo/demo.o
```

## API Documentation

Access Swagger documentation at: http://localhost:8080/swagger/index.html

## Main Endpoints

- `POST /api/v1/chat` - Create a new chat
- `POST /api/v1/message` - Send message for processing
- `POST /api/v1/memory/chat/{chat_id}/fetch` - Fetch relevant memories
- `GET /api/v1/memory/short-term/chat/{chat_id}` - List short-term memories
- `GET /api/v1/memory/long-term/chat/{chat_id}` - List long-term memories

## Technologies

- **Go** - API and Worker
- **Python** - ML inference service
- **MongoDB** - Main database
- **Redis** - Task queue (Asynq)
- **Qdrant** - Vector database
- **gRPC** - Inter-service communication
