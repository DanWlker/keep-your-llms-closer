# LLM Parameters

> Controls that affect how LLMs generate output.

---

## Core Parameters

### Temperature

**Controls randomness:**

| Value | Effect |
|-------|--------|
| 0.0 | Deterministic - same input = same output |
| 0.7 | Balanced (recommended for most tasks) |
| 1.0 | Very creative/random |

```python
# Low temperature - factual tasks
llm.chat(prompt, temperature=0.0)

# High temperature - creative writing
llm.chat(prompt, temperature=0.9)
```

**Analogy:** Temperature is like "how much dice rolling" the model does.

---

### Top-P (Nucleus Sampling)

**Controls vocabulary diversity:**

- `top_p=1.0` - Consider all tokens
- `top_p=0.9` - Consider top 90% of probable tokens
- `top_p=0.1` - Only most probable tokens

**Rule of thumb:** Use either temperature OR top_p, not both.

---

### Max Tokens

**Limits output length:**

```python
llm.chat(prompt, max_tokens=1000)  # Stop after 1000 tokens
```

---

### Top-K

**Alternative to top_p:**

- `top_k=1` - Always pick most probable next token
- `top_k=40` - Pick from top 40 tokens

---

## Context Parameters

### Context Window

Maximum input + output tokens:

| Model | Context |
|-------|---------|
| GPT-4 | 128K |
| Claude 3.5 | 200K |
| GPT-4o | 1M |

---

### Presence Penalty

**Discourages repeating same words:**

```python
llm.chat(prompt, presence_penalty=0.5)  # Higher = less repetition
```

---

### Frequency Penalty

**Encourages diversity:**

```python
llm.chat(prompt, frequency_penalty=0.5)  # Higher = more diverse
```

---

## Output Structure

### JSON Mode

Force JSON output:

```python
llm.chat(prompt, response_format={"type": "json_object"})
```

---

### Function Calling

Structured output as tool calls:

```python
response = llm.chat(
    prompt,
    tools=[{
        "type": "function",
        "function": {
            "name": "get_weather",
            "parameters": {"type": "object", "properties": {...}}
        }
    }]
)
```

---

## Summary Table

| Parameter | What it controls | Common values |
|-----------|------------------|---------------|
| Temperature | Randomness | 0.0-1.0 |
| Top-P | Vocabulary selection | 0.7-1.0 |
| Max Tokens | Output length | 100-32000 |
| Top-K | Token candidates | 1-100 |
| Presence Penalty | Word repetition | -2.0 to 2.0 |
| Frequency Penalty | Token diversity | -2.0 to 2.0 |

---

## Recommendations

| Task | Temperature | Top-P |
|------|-------------|-------|
| Code generation | 0.0-0.2 | 0.95 |
| Factual Q&A | 0.0-0.3 | 0.9 |
| Creative writing | 0.7-0.9 | 0.9 |
| Summarization | 0.0-0.3 | 0.9 |
| Agent tool use | 0.0-0.2 | 0.95 |

---

## Key Takeaways

1. **Temperature** - Main control for creativity vs accuracy
2. **Top-P** - Alternative to temperature
3. **Lower = more predictable** - Good for code/facts
4. **Higher = more creative** - Good for writing
5. **Don't use both** - Temperature and top_p together
