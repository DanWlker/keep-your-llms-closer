-- Enable pgvector extension
CREATE EXTENSION IF NOT EXISTS vector;

-- Table for document chunks
CREATE TABLE IF NOT EXISTS document_chunks (
    id SERIAL PRIMARY KEY,
    content TEXT NOT NULL,
    embedding vector(1536),
    metadata JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Index for vector similarity search
-- Using IVFFlat for faster search on larger datasets
CREATE INDEX ON document_chunks USING ivfflat (embedding vector_cosine_ops);

-- Optional: Table for tracking documents
CREATE TABLE IF NOT EXISTS documents (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    source TEXT,
    metadata JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
