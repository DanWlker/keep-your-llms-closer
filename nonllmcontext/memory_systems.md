# Memory Systems

> How agents store and retrieve information across interactions.

---

## The Problem

LLMs have no memory by default:
- Each conversation is fresh
- Can't remember previous sessions
- Context window limits what can be remembered

**Memory systems** solve this.

---

## Types of Memory

### Short-Term Memory (Working Memory)

What's currently in the context window:

```
Current conversation
User: What was my name?
LLM: You told me earlier it's Daniel.
         ↑ This is in context
```

### Long-Term Memory

Persists across sessions:

```
Session 1:
User: My name is Daniel.
LLM: Nice to meet you, Daniel!

Session 2 (tomorrow):
User: What's my name?
LLM: Your name is Daniel.
         ↑ Retrieved from long-term memory
```

---

## Memory Architectures

### 1. Simple History

Store all messages, replay:

```python
memory = []

def chat(user_input):
    memory.append({"role": "user", "content": user_input})
    
    response = llm.chat(memory)
    
    memory.append({"role": "assistant", "content": response})
    
    return response
```

**Problem:** Eventually hits context limits.

### 2. Summarization

Summarize old messages:

```python
def chat(user_input):
    if len(messages) > 10:
        summary = llm.summarize(messages[:-5])
        messages = [{"role": "system", "content": summary}] + messages[-5:]
    
    return llm.chat(messages + user_input)
```

### 3. Vector Memory (RAG-based)

Store embeddings, retrieve relevant:

```python
# Store
memory_embeddings.add({
    "content": "User's name is Daniel",
    "timestamp": "2024-01-01"
})

# Retrieve
def get_context(user_input):
    query_emb = embed(user_input)
    relevant = memory_embeddings.search(query_emb, top_k=3)
    return relevant
```

---

## Memory Types in Agents

### 1. Conversation Memory

What happened in current session:

```python
conversation_memory = [
    {"role": "user", "content": "Create a file"},
    {"role": "assistant", "content": "I'll create it"},
    {"role": "tool", "content": "File created: app.py"},
]
```

### 2. Entity Memory

Facts about entities:

```python
entity_memory = {
    "user": {"name": "Daniel", "preferences": {...}},
    "project": {"name": "myapp", "framework": "react"},
}
```

### 3. Procedural Memory

How to do things:

```python
procedural_memory = {
    "run_tests": "Use npm test in project root",
    "deploy": "Run npm run build then npm run deploy",
}
```

---

## Implementation Patterns

### sliding Window

Only remember last N messages:

```python
def chat(messages):
    return messages[-10:]  # Last 10 messages
```

### Importance-Based

Remember important, forget trivial:

```python
def should_remember(message):
    # Mark important messages
    if "remember" in message or "don't forget" in message:
        return True
    return importance_score(message) > 0.7
```

### Time-Based

Remember recent, fade old:

```python
def get_recent_memory(hours=24):
    return memory.filter(
        timestamp > now() - hours
    )
```

---

## Tools for Memory

| Tool | Description |
|------|-------------|
| Mem0 | Memory layer for AI apps |
| LangChain memory | Conversation memory |
| Zep | Long-term memory |
| Graph memory | Knowledge graphs |

---

## Memory in Agents

Full agent memory hierarchy:

```
┌─────────────────────────────────────────┐
│              Agent                       │
├─────────────────────────────────────────┤
│  Working memory (context)                │
│  - Current task                          │
│  - Recent conversation                  │
│  - Current plan                         │
├─────────────────────────────────────────┤
│  Short-term memory (session)             │
│  - This conversation                    │
│  - Files created                        │
│  - Commands run                          │
├─────────────────────────────────────────┤
│  Long-term memory (persistent)           │
│  - User preferences                     │
│  - Past projects                        │
│  - Learned facts                        │
└─────────────────────────────────────────┘
```

---

## Retrieval

### How to find relevant memory

```python
def retrieve(query):
    # 1. Vector similarity
    similar = vector_db.search(embed(query))
    
    # 2. Time decay
    recent = weighted_by_time(similar)
    
    # 3. Importance
    important = filter_by_importance(recent)
    
    return important[:5]
```

---

## Key Takeaways

1. **Short-term** = Context window (current session)
2. **Long-term** = Persists across sessions
3. **Entity** = Facts about things
4. **Procedural** = How to do things
5. **Retrieval** = Vector search + time + importance

---

## References

- Mem0 docs
- LangChain memory
- "MemGPT" paper
