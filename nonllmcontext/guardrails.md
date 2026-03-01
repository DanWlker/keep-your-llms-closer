# Guardrails

> Systems that filter/control LLM outputs to prevent unwanted content.

---

## What It Is

Guardrails = Filters that check LLM outputs (and sometimes inputs) for unwanted content.

```
User Input → LLM → Output → Guardrail Filter → User
                              ↓
                       Blocked/Modified/Allowed
```

---

## Types of Guardrails

### Input Guardrails

Check user input before it reaches the LLM:

```python
def input_guardrail(user_input):
    if contains_pii(user_input):
        block()
    if is_dangerous_request(user_input):
        require_approval()
    return user_input
```

### Output Guardrails

Check LLM response before returning to user:

```python
def output_guardrail(llm_response):
    if contains_harmful_content(llm_response):
        return "I can't help with that."
    if contains_pii(llm_response):
        redact_pii(llm_response)
    return llm_response
```

### Latency Guardrails

- Timeout checks
- Resource limits

---

## What Guardrails Can Catch

| Category | Examples |
|----------|----------|
| Harmful content | Violence, self-harm, hate speech |
| PII | Names, emails, SSN, addresses |
| Sensitive data | API keys, passwords, secrets |
| Jailbreak attempts | Prompt injection attempts |
|hallucinations | Fact-checking against sources |
| Length limits | Max tokens, response size |

---

## How Guardrails Work

### 1. Rule-Based

```python
def guardrail(output):
    blocked_words = ["bomb", "weapon", "attack"]
    for word in blocked_words:
        if word in output.lower():
            return filtered_output
    return output
```

### 2. ML-Based

```python
# Classifier models
classifier = load_harmful_content_model()

def guardrail(output):
    score = classifier.predict(output)
    if score > 0.8:
        return "Content blocked"
    return output
```

### 3. LLM-Based

```python
def guardrail(output):
    response = llm.chat(f"""
    Does this contain harmful content? Yes/No:
    {output}
    """)
    if "yes" in response.lower():
        return "Content blocked"
    return output
```

---

## Guardrails vs Prompt Injection Defense

| Aspect | Guardrails | Prompt Injection Defense |
|--------|------------|-------------------------|
| **Scope** | Input + Output | Input only |
| **When** | Before/after LLM | Before LLM |
| **What** | Content safety | Instruction manipulation |
| **Example** | Block violent content | Block "ignore instructions" |

They work together but are different.

---

## Popular Guardrail Tools

| Tool | Description |
|------|-------------|
| NVIDIA NeMo Guardrails | Open source, programmable |
| AWS Bedrock Guardrails | AWS managed service |
| Azure AI Content Safety | Microsoft |
| Google Safety Attributes | Google |
| Guardrails AI | Commercial |
| Rebuff | Prompt injection detection |

### NeMo Guardrails Example

```yaml
# config.yml
guards:
  - name: harmful_content
    type: topical
    criteria:
      - not_contains:
          - violence
          - hate speech
          - self-harm
    
  - name: PII
    type: pii
    action: redact
```

---

## Implementation Best Practices

### Layered Approach

```
User Input
    ↓
Input Guardrail (block/filter)
    ↓
LLM
    ↓
Output Guardrail (block/filter)
    ↓
User Response
```

### Fail Safe

```python
def chat_with_guardrails(user_input):
    # Default to safe if guardrail fails
    try:
        filtered_input = input_guardrail(user_input)
        response = llm.chat(filtered_input)
        safe_response = output_guardrail(response)
        return safe_response
    except:
        return "I couldn't process that request safely."
```

---

## Key Takeaways

1. **Guardrails** = Filters on input/output for safety
2. **Input** = Block bad requests before LLM
3. **Output** = Block bad responses after LLM
4. **Types** = Rule-based, ML-based, LLM-based
5. **Layered** = Both input and output protection

---

## References

- NVIDIA NeMo Guardrails docs
- AWS Bedrock Guardrails
- OWASP LLM Top 10
