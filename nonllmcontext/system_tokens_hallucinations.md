# System Prompts, Token Limits, and Hallucinations

> Key LLM concepts: how instructions work, context limits, and why LLMs make things up.

---

## System Prompts

### What It Is

The **system prompt** (or system message) is instructions that shape all responses from the LLM. It's sent with every request.

```
System: "You are a helpful Python expert. Write clean, idiomatic code."

User: "How do I sort a list?"

LLM response: [Uses Python expertise, clean code style]
```

### How It Works

The system prompt is prepended to the conversation:

```
[System prompt]
You are a helpful coding assistant.

[Conversation history]
User: What is a list comprehension?
Assistant: A list comprehension is...

User: How do I sort one?

[Current]
LLM generates: sorted() function, with clean code examples...
```

### What Can Go in System Prompt

- **Role**: "You are a senior software engineer"
- **Tone**: "Be concise and practical"
- **Constraints**: "Never use external libraries"
- **Format**: "Always use TypeScript"
- **Context**: "You have access to the file system"

### Priority

System prompt > User message > Conversation history

The LLM weights system instructions most heavily.

---

## Token Limits

### What Are Tokens

Tokens = pieces of words (not exactly words):

```
"hello world" → 2 tokens
"unconstitutional" → might be 3-4 tokens (subword)
```

~1 token ≈ 0.75 words

### Context Window

The maximum input + output tokens per request:

| Model | Context Window |
|-------|----------------|
| GPT-4 | 128K tokens |
| GPT-4o | 1M tokens |
| Claude 3.5 | 200K tokens |
| Claude 4 | 1M tokens |

### How It Works

```
You: "Here's a 50K token document. Summarize it."
         ↓
LLM: [Can't fit in context]

Framework: "Document too long. Split into chunks."
```

The framework manages truncation when exceeding limits.

### Implications

- Limited context = can't reference entire codebase
- Solutions: RAG, chunking, summarization

---

## Hallucinations

### What It Is

Hallucinations = LLMs generating false or made-up information.

```
User: "Who invented JSON?"
LLM: "JSON was invented by Douglas Crockford in 2005"
                                      ↑
                              Made up (actually standardized 2002)
```

### Why It Happens

**LLMs are text predictors, not truth databases.**

```
Training: "The model learned that after X, Y often follows"
                ↓
Generation: "Given this context, what's most likely next text?"
                ↓
If "JSON" often appears near "Crockford", model generates it
```

The model generates **statistically likely** text, not **verified** text.

### Factors That Increase Hallucinations

| Factor | Why |
|--------|-----|
| Low temperature | Less exploration, more confident (wrong) answers |
| Unfamiliar topics | Less training data to anchor on |
| Long generation | More chances to drift |
| No grounding | No RAG / external info |

### Reducing Hallucinations

| Technique | How it helps |
|-----------|--------------|
| RAG | Ground answers in real data |
| Tool use | Fetch real information |
| Lower temperature | More deterministic |
| Prompt engineering | "If unsure, say you don't know" |
| Citations | Ask for sources |

---

## Key Takeaways

1. **System prompt** - Instructions that shape all responses, weighted most heavily
2. **Token limits** - Context window limits how much can be processed
3. **Hallucinations** - LLM predicts likely text, not truth - use RAG/tools to ground

---

## References

- LLM prompting guides
- Anthropic on Claude's training
- "Attention Is All You Need" paper
