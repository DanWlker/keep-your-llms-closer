package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/pgvector/pgvector-go"
	"github.com/sashabaranov/go-openai"
)

func main() {
	ctx := context.Background()

	// Connect to PostgreSQL
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		connStr = "postgres://postgres:postgres@localhost:5432/rag?sslmode=disable"
	}

	conn, err := pgx.Connect(ctx, connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer conn.Close(ctx)

	// Initialize pgvector
	_, err = conn.Exec(ctx, "CREATE EXTENSION IF NOT EXISTS vector")
	if err != nil {
		log.Fatalf("Failed to create vector extension: %v", err)
	}

	// Create table if not exists
	_, err = conn.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS document_chunks (
			id SERIAL PRIMARY KEY,
			content TEXT NOT NULL,
			embedding vector(1536),
			metadata JSONB,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	// Create index if not exists
	_, err = conn.Exec(ctx, `
		CREATE INDEX IF NOT EXISTS document_chunks_idx 
		ON document_chunks USING ivfflat (embedding vector_cosine_ops)
	`)
	if err != nil {
		log.Fatalf("Failed to create index: %v", err)
	}

	// Initialize OpenAI client
	openaiKey := os.Getenv("OPENAI_API_KEY")
	if openaiKey == "" {
		log.Fatal("OPENAI_API_KEY not set")
	}
	client := openai.NewClient(openaiKey)

	// Example: Add a document
	content := `Go is a programming language. 
It was created at Google in 2009 by Robert Griesemer, Rob Pike, and Ken Thompson.
Go is known for its simplicity, concurrency support, and fast compilation.`

	err = AddDocument(ctx, client, conn, "about-go", content, nil)
	if err != nil {
		log.Fatalf("Failed to add document: %v", err)
	}

	// Example: Query
	question := "Who created Go?"
	answer, err := QueryRAG(ctx, client, conn, question)
	if err != nil {
		log.Fatalf("Query failed: %v", err)
	}

	fmt.Printf("Question: %s\n", question)
	fmt.Printf("Answer: %s\n", answer)
}

// AddDocument chunks text, embeds it, and stores in pgvector
func AddDocument(ctx context.Context, client *openai.Client, conn *pgx.Conn, docID, content string, metadata map[string]interface{}) error {
	// Simple chunking - split by sentences (improve for production)
	chunks := chunkText(content, 500)

	for _, chunk := range chunks {
		// Get embedding from OpenAI
		resp, err := client.CreateEmbeddings(ctx, openai.EmbeddingRequest{
			Model: openai.Embedding3Small,
			Input: []string{chunk},
		})
		if err != nil {
			return fmt.Errorf("failed to create embedding: %w", err)
		}

		embedding := pgvector.NewVector(resp.Data[0].Embedding)

		// Store in pgvector
		_, err = conn.Exec(ctx,
			"INSERT INTO document_chunks (content, embedding, metadata) VALUES ($1, $2, $3)",
			chunk, embedding, metadata,
		)
		if err != nil {
			return fmt.Errorf("failed to insert chunk: %w", err)
		}
	}

	return nil
}

// QueryRAG takes a question, searches vectors, and asks LLM
func QueryRAG(ctx context.Context, client *openai.Client, conn *pgx.Conn, question string) (string, error) {
	// Embed the question
	resp, err := client.CreateEmbeddings(ctx, openai.EmbeddingRequest{
		Model: openai.Embedding3Small,
		Input: []string{question},
	})
	if err != nil {
		return "", fmt.Errorf("failed to create question embedding: %w", err)
	}

	questionEmbedding := pgvector.NewVector(resp.Data[0].Embedding)

	// Search pgvector for similar chunks
	rows, err := conn.Query(ctx,
		`SELECT content FROM document_chunks 
		 ORDER BY embedding <=> $1 
		 LIMIT 5`,
		questionEmbedding,
	)
	if err != nil {
		return "", fmt.Errorf("failed to search: %w", err)
	}
	defer rows.Close()

	var chunks []string
	for rows.Next() {
		var content string
		if err := rows.Scan(&content); err != nil {
			return "", err
		}
		chunks = append(chunks, content)
	}

	if len(chunks) == 0 {
		return "No relevant context found.", nil
	}

	// Build context for LLM
	context := strings.Join(chunks, "\n\n")

	// Call LLM with context
	prompt := fmt.Sprintf(`You are a helpful assistant. Use the following context to answer the question.

Context:
%s

Question: %s

Answer:`, context, question)

	chatResp, err := client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: openai.GPT4oMini,
		Messages: []openai.ChatCompletionMessage{
			{Role: "user", Content: prompt},
		},
	})
	if err != nil {
		return "", fmt.Errorf("failed to call LLM: %w", err)
	}

	return chatResp.Choices[0].Message.Content, nil
}

// chunkText splits text into chunks of approximately maxSize characters
// This is a simple implementation - consider using more sophisticated chunking for production
func chunkText(text string, maxSize int) []string {
	if len(text) <= maxSize {
		return []string{text}
	}

	var chunks []string
	paragraphs := strings.Split(text, "\n")

	var current strings.Builder
	for _, p := range paragraphs {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}

		if current.Len()+len(p) > maxSize {
			if current.Len() > 0 {
				chunks = append(chunks, current.String())
				current.Reset()
			}
		}
		current.WriteString(p)
		current.WriteString("\n")
	}

	if current.Len() > 0 {
		chunks = append(chunks, current.String())
	}

	return chunks
}
