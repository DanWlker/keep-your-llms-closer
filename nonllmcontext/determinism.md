# Seed / Determinism

> How to get reproducible outputs from LLMs.

---

## The Problem

LLMs are probabilistic - same input can give different outputs:

```
User: What is 2+2?
LLM: 4

User: What is 2+2?
LLM: It's 4

User: What is 2+2?
LLM: Four
```

Even with same prompt, randomness kicks in.

---

## Seed

### What It Does

Seed forces the LLM to use the same randomness:

```
User: What is 2+2?
LLM: 4 (seed=42)

User: What is 2+2?
LLM: 4 (seed=42)  ← Same!

User: What is 2+2?
LLM: It's 4 (seed=99)  ← Different seed, different output
```

### How It Works

```python
# Same seed = same output
response1 = llm.chat(prompt, seed=42)
response2 = llm.chat(prompt, seed=42)

assert response1 == response2  # Usually passes
```

### What "Same" Means

With same seed, you get:
- Same tokens (deterministic sampling)
- Same "thinking" tokens if applicable

---

## When Seeds Matter

| Use Case | Why Seed |
|----------|----------|
| Testing | Reproduce bugs |
| Debugging | Get same output to investigate |
| Caching | Cache responses by prompt+seed |
| A/B testing | Compare prompt changes fairly |

---

## Temperature = 0

### What It Does

Temperature controls randomness. Setting it to 0 makes output more deterministic:

```
Temperature 0:     Most probable token always chosen
Temperature 0.7:  Some randomness
Temperature 1.0:  High randomness
```

### Seed + Temperature 0

Most deterministic:

```python
response = llm.chat(
    prompt,
    temperature=0,
    seed=42
)
```

---

## Limitations

### Not Fully Deterministic

Even with seed + temperature 0, may differ across:

| Factor | Why |
|--------|-----|
| Different model versions | Weights may differ |
| Different hardware | Floating point precision |
| Different API versions | Internal changes |
| Batching | Parallel processing |

### What Usually Works

- Same model + same version + same seed = same output
- Good enough for testing/debugging

### What Usually Fails

- Cross-model (GPT-4 vs Claude)
- Cross-version (Claude 3.5 vs 4)
- Production use (API may change)

---

## Practical Use

### Testing

```python
def test_llm_response():
    response = llm.chat(
        "What is 2+2?",
        temperature=0,
        seed=42
    )
    assert response == "4"  # Will pass consistently
```

### Caching

```python
cache = {}

def get_response(prompt):
    seed = hash(prompt)  # Consistent seed per prompt
    
    if prompt in cache:
        return cache[prompt]
    
    response = llm.chat(prompt, seed=seed)
    cache[prompt] = response
    return response
```

---

## Key Takeaways

1. **Seed** = Forces same randomness source = more reproducible
2. **Temperature 0** = Always pick most probable token
3. **Together** = Most deterministic output
4. **Not guaranteed** = May differ across versions/hardware
5. **Use for** = Testing, debugging, caching

---

## References

- OpenAI reproducibility docs
- Anthropic seed support
