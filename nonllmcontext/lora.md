# LoRA / QLoRA

> Efficient fine-tuning methods that train only small adapters instead of full models.

---

## The Problem with Full Fine-tuning

Full fine-tuning updates ALL model weights:

```
Base model (7B parameters)
        ↓
Fine-tune on your data
        ↓
New model (7B updated weights)
```

Problems:
- Requires huge GPU memory (7B × 2 bytes = 14GB+ for training)
- Time consuming (hours)
- Risk of catastrophic forgetting

---

## LoRA: Low-Rank Adaptation

### The Idea

Instead of changing all weights, add small "adapter" matrices:

```
Base model (frozen)
        ↓
+ Small adapters (trained)
        ↓
Model with new behavior
```

### How It Works

```
Original weight: W (d × d matrix)

LoRA adds: ΔW = B × A

Where:
- B is (d × r)
- A is (r × d)
- r is "rank" (usually 8-32)

Total parameters: d×r + r×d instead of d×d

Example: d=4096, r=8
- Original: 16M parameters
- LoRA: 65K parameters (250x smaller!)
```

### Visual

```
Input → [Frozen weights] → Output
        ↘                ↙
       [LoRA adapters] → added to output
```

### Training

```
1. Forward pass through frozen model
2. Calculate loss
3. Update only LoRA matrices (A and B)
4. At inference: combine frozen + LoRA
```

---

## LoRA Config

### Key Parameters

| Parameter | Typical Value | Notes |
|-----------|--------------|-------|
| rank (r) | 8-32 | Higher = more expressive, more params |
| alpha | r × 2 | Scaling factor |
| dropout | 0.05 | Regularization |
| target modules | q_proj, v_proj | Which layers to apply to |

### Example (Hugging Face)

```python
from peft import LoraConfig

config = LoraConfig(
    r=16,
    lora_alpha=32,
    lora_dropout=0.05,
    target_modules=["q_proj", "v_proj", "k_proj", "o_proj"],
    bias="none",
    task_type="CAUSAL_LM"
)
```

---

## QLoRA: LoRA + Quantization

### The Idea

Combine LoRA with quantization to fine-tune HUGE models on consumer hardware.

```
Base model: 65B parameters (full precision)
        ↓
Quantize to 4-bit
        ↓
Add LoRA adapters
        ↓
Fine-tune!
```

### How It Works

1. **Quantize** - Compress model to 4-bit (uses ~8GB for 65B model)
2. **LoRA** - Add adapters as usual
3. **Train** - Only train adapters, not quantized weights
4. **Inference** - Combine quantized + LoRA

### Bits and Bytes

| Method | Memory for 7B | Memory for 65B |
|--------|--------------|-----------------|
| Full fp16 | 14GB | 130GB |
| LoRA | ~4GB | ~40GB |
| QLoRA | ~2GB | ~8GB |

Consumer GPU: ~24GB

- 7B: Full fine-tune ✅
- 13B: LoRA ✅
- 65B: QLoRA ✅

---

## LoRA vs QLoRA

| Aspect | LoRA | QLoRA |
|--------|------|-------|
| Memory | Medium (needs base model in memory) | Low (quantized) |
| Quality | Excellent | Slightly lower |
| Speed | Fast | Slower (dequantization) |
| Hardware | GPU with 20GB+ | Consumer GPU |
| Can use | Any model | Large models |

---

## LoRA in Practice

### Training Data

Same as fine-tuning - input/output pairs or chat format.

### Multiple LoRAs

You can train multiple LoRAs for different tasks:

```
base_model
    ├── + code_loRA → Code expert
    ├── + chat_loRA → Chat expert  
    └── + creative_loRA → Creative writing
```

Switch by changing which LoRA is loaded.

### Merging

```python
# Merge LoRA into base model
model = model.merge_and_unload()
```

After merging, you get a single model (no adapters).

---

## Tools

| Tool | Description |
|------|-------------|
| PEFT | Hugging Face's LoRA library |
| Axolotl | Popular for LoRA training |
| Unsloth | Very fast LoRA training (2x faster) |
| Llama Factory | Unified fine-tuning |
| QLoRA | Original QLoRA implementation |

---

## Best Practices

### Rank Selection

| Rank | Use Case |
|------|----------|
| 4-8 | Simple tasks, limited data |
| 16-32 | General use (recommended) |
| 64-128 | Complex tasks, lots of data |

### Which Layers to Target

Most impactful:
- `q_proj`, `v_proj` - Attention (most important)
- `k_proj`, `o_proj` - Also helpful

Less impactful but sometimes useful:
- `gate_proj`, `up_proj`, `down_proj` - MLP layers

### Data Quality

Same as fine-tuning: quality > quantity. 100 great examples > 10K mediocre.

---

## Key Takeaways

1. **LoRA** = Train small adapters, keep base model frozen
2. **QLoRA** = Quantize model + add LoRA = fine-tune huge models on CPU
3. **Rank** = Size of adapter (8-32 typical)
4. **Memory** = LoRA needs 20GB+ GPU, QLoRA needs 8GB+
5. **Multiple adapters** = Switch between different behaviors

---

## References

- LoRA paper: "LoRA: Low-Rank Adaptation of Large Language Models"
- QLoRA paper: "QLoRA: Efficient Finetuning of Quantized LLMs"
- PEFT docs: https://huggingface.co/docs/peft
- Unsloth: https://unsloth.ai
