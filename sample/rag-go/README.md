# Go RAG Example

A simple Go-native RAG (Retrieval Augmented Generation) application using:
- PostgreSQL with pgvector
- OpenAI API for embeddings and LLM

## Prerequisites

1. PostgreSQL with pgvector extension
2. Go 1.21+
3. OpenAI API key

## Setup

### 1. Create Database

```bash
# Connect to PostgreSQL
psql -U postgres

# Create database and enable pgvector
CREATE DATABASE rag;
\c rag
CREATE EXTENSION IF NOT EXISTS vector;
\i schema.sql
```

### 2. Configure Environment

```bash
cp .env.example .env
# Edit .env with your DATABASE_URL and OPENAI_API_KEY
```

### 3. Install Dependencies

```bash
go mod tidy
```

### 4. Run

```bash
go run main.go
```

## Project Structure

```
rag-go/
├── main.go         # RAG implementation
├── schema.sql      # Database schema
├── .env.example    # Environment template
└── go.mod          # Go module
```

## Key Functions

- `AddDocument()` - Chunks text, embeds, stores in pgvector
- `QueryRAG()` - Embeds question, searches vectors, asks LLM
- `chunkText()` - Simple text chunking

## Notes

- This is a basic example for learning
- For production: improve chunking, add error handling, use proper config
- Consider: LlamaIndex Go bindings, more sophisticated chunking, caching
