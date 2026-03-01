# RAG and Vector Databases

> How LLMs access external knowledge through retrieval.

---

## The Problem

LLMs have fixed knowledge - they can't read your docs, codebase, or internal data.

**Solution:** Give them a way to look stuff up.

---

## How It Works

```
User: "What's our authentication flow?"

1. Convert question to embedding (vector)
2. Search vector DB for similar text
3. Return top matches as context
4. LLM reads: question + retrieved docs
5. LLM generates answer
```

---

## Embeddings

**Embeddings** convert text to numerical vectors (arrays of numbers):

```
"hello world" → [0.1, -0.3, 0.5, 0.8, ...]
"hi there"    → [0.12, -0.29, 0.48, 0.79, ...]
```

Similar text = similar vectors = close in vector space.

---

## Vector Databases

Specialized databases that:
1. Store embeddings
2. Search by similarity (not exact match)
3. Return top-K results

| Database | Description |
|----------|-------------|
| Pinecone | Managed vector DB |
| Weaviate | Open source |
| Chroma | Lightweight, Python |
| pgvector | PostgreSQL extension |
| Qdrant | Open source |

---

## RAG Pipeline

```
┌──────────────┐     ┌──────────────┐     ┌──────────────┐
│  User Query  │────▶│  Embedding   │────▶│   Vector DB  │
└──────────────┘     │   (encode)   │     │   (search)   │
                     └──────────────┘     └──────────────┘
                                                  │
                                                  ▼
┌──────────────┐     ┌──────────────┐     ┌──────────────┐
│    Answer    │◀────│     LLM      │◀────│   Context   │
└──────────────┘     │  (generate)  │     │  (retrieved)│
                     └──────────────┘     └──────────────┘
```

---

## Implementation

### 1. Indexing (Build the DB)

```python
# 1. Load documents
docs = load_docs("./docs")

# 2. Split into chunks
chunks = split(docs, chunk_size=500)

# 3. Embed each chunk
embeddings = embed_model.encode(chunks)

# 4. Store in vector DB
vector_db.add(embeddings, chunks)
```

### 2. Querying (At runtime)

```python
# 1. Embed user question
query_embedding = embed_model.encode(user_question)

# 2. Search vector DB
results = vector_db.search(query_embedding, top_k=5)

# 3. Build context
context = "\n\n".join([r.text for r in results])

# 4. Prompt LLM
response = llm.chat(f"""
Context:
{context}

Question: {user_question}
""")
```

---

## Chunking Strategies

How you split documents matters:

| Strategy | Description | Use case |
|----------|-------------|----------|
| Fixed size | Every N characters | Simple |
| By paragraph | Split on newlines | Prose |
| By heading | Split on headers | Markdown |
| Recursive | Multi-level splitting | Code |
| Semantic | Split by meaning | Complex docs |

---

## Embedding Models

| Model | Dimensions | Notes |
|-------|-----------|-------|
| OpenAI text-embedding-3-small | 1536 | Good quality |
| text-embedding-3-large | 3072 | Better quality |
| Cohere embed-multilingual | 1024 | Multilingual |
| BGE | 1024 | Open source |

---

## RAG + MCP

MCP can serve as the "R" in RAG:

```
MCP server with resources
        ↓
User asks question
        ↓
Framework retrieves via MCP
        ↓
LLM generates answer
```

---

## Key Takeaways

1. **Embeddings** - Convert text to vectors
2. **Vector DB** - Search by similarity
3. **RAG** - Retrieve, then generate
4. **Chunking** - How you split docs matters
5. **Not perfect** - Can miss relevant info

---

## References

- Pinecone docs
- LangChain RAG guides
- Chroma (Python vector DB)
