# Fine-tuning

> Customizing model behavior with your own data.

---

## What It Is

Fine-tuning = Taking a pre-trained model and continuing training on your own data to change its behavior.

```
Pre-trained model (knows everything)
        ↓
Fine-tune on your data
        ↓
Model with custom behavior
```

---

## Why Fine-tune

| Use Case | Why Fine-tune |
|----------|---------------|
| Custom style | Write like your company docs |
| Specific tasks | Better at your domain (legal, medical, code) |
| Cost reduction | Smaller model + fine-tuned > bigger model |
| Latency | Faster responses than prompting |

---

## How It Works

### Basic Process

```
1. Prepare training data
   - Input → Output pairs
   - e.g., {"input": "Hello", "output": "Hi there!"}

2. Configure training
   - Learning rate (usually very low)
   - Number of epochs
   - Model to start with

3. Train
   - Model continues learning on your data
   - Weights adjust to minimize loss on your data

4. Evaluate
   - Test on held-out data

5. Deploy
   - Use the fine-tuned model
```

### Training Data Format

```json
[
  {"prompt": "Hello", "completion": "Hi there!"},
  {"prompt": "How are you?", "completion": "I'm doing well, thanks!"},
  {"prompt": "What's up?", "completion": "Not much, how can I help?"}
]
```

Or chat format:

```json
{
  "messages": [
    {"role": "user", "content": "Hello"},
    {"role": "assistant", "content": "Hi there!"}
  ]
}
```

---

## Types of Fine-tuning

### Full Fine-tuning

```
All model weights are updated
```

- Pros: Maximum customization
- Cons: Expensive, can forget original capabilities

### LoRA (Low-Rank Adaptation)

```
Only small "adapters" are trained
Original weights stay frozen
```

- Pros: Fast, cheap, doesn't forget
- Cons: Less customization than full

### QLoRA

```
LoRA + Quantization (4-bit weights)
```

- Pros: Can fine-tune huge models on consumer hardware
- Cons: Slightly less quality

---

## Full Fine-tuning vs LoRA

| Aspect | Full Fine-tuning | LoRA |
|--------|-----------------|------|
| What changes | All weights | Adapter weights only |
| GPU memory | High (full model) | Low (frozen model + small adapter) |
| Training time | Hours | Minutes to hours |
| Performance | Best | Good (sometimes close) |
| Catastrophic forgetting | High risk | Low risk |
| Deployment | Replace model | Add adapter to base |

---

## When to Use What

| Scenario | Recommendation |
|----------|---------------|
| Limited data (< 1000 examples) | Don't fine-tune, use prompting/RAG |
| Medium data (1K - 100K) | LoRA |
| Large data (100K+) | Full fine-tuning or LoRA |
| Production cost critical | LoRA |
| Maximum quality needed | Full fine-tuning |

---

## Key Parameters

### Learning Rate

- Usually 1e-5 to 1e-4 (very small)
- Too high = model forgets everything
- Too low = no learning

### Epochs

- How many times to see all data
- Too many = overfitting
- Too few = underfitting

### Batch Size

- How many examples processed before updating
- Higher = more stable, more memory

---

## Data Quality Matters

### Quantity vs Quality

| Data Size | Quality Needed |
|-----------|----------------|
| 100 examples | Perfect quality |
| 1,000 examples | High quality |
| 10,000+ examples | Moderate OK |

### Common Mistakes

1. **Inconsistent format** - All examples should follow same pattern
2. **No validation set** - Can't measure if it's working
3. **Too little data** - <100 examples usually not worth it
4. **Data leakage** - Test data in training data

---

## Alternatives to Fine-tuning

| Method | When to Use |
|--------|-------------|
| Prompting | Works for most cases |
| Few-shot learning | Add 3-5 examples in prompt |
| RAG | Access external knowledge |
| Fine-tuning | Need consistent style/behavior |

Rule of thumb: Try prompting first. Fine-tune only when prompting doesn't work.

---

## Platforms / Tools

| Tool | Description |
|------|-------------|
| OpenAI fine-tuning | Managed service |
| Hugging Face Transformers | Open source |
| Axolotl | Popular for LoRA |
| Unsloth | Fast LoRA training |
| Llama Factory | Unified fine-tuning |

---

## Evaluation

### Metrics

- Loss (during training)
- Perplexity (how confused model is)
- Human evaluation
- Task-specific metrics (accuracy, F1, etc.)

### Holdout Set

Always split data:
- Training set (80%)
- Validation set (10%)
- Test set (10%)

---

## Key Takeaways

1. **Fine-tuning** = Continue training on your data
2. **LoRA** = Fast, cheap, trains small adapters
3. **Data quality > quantity** - Better 100 great examples than 10K mediocre
4. **Start with prompting** - Only fine-tune when necessary
5. **Catastrophic forgetting** - Model may lose original capabilities

---

## References

- OpenAI fine-tuning docs
- Hugging Face fine-tuning guides
- LoRA paper: "LoRA: Low-Rank Adaptation of Large Language Models"
