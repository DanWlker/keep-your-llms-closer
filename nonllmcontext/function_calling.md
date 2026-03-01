# Function Calling & JSON Mode

> How LLMs output structured data vs MCP.

---

## Function Calling

### What It Is

Function calling = LLM outputs structured JSON that your code executes.

```
User: What's the weather in Tokyo?

LLM generates:
{
  "name": "get_weather",
  "arguments": {
    "location": "Tokyo"
  }
}

Your code parses this → calls actual function → returns result to LLM
```

### How It Works

```
1. You define functions in the prompt/tool schema
2. LLM decides it needs external data
3. LLM outputs structured JSON
4. Your code parses and executes
5. Result returned to LLM
6. LLM generates final response
```

### Defining Functions

```python
functions = [
    {
        "name": "get_weather",
        "description": "Get weather for a location",
        "parameters": {
            "type": "object",
            "properties": {
                "location": {
                    "type": "string",
                    "description": "City name"
                }
            },
            "required": ["location"]
        }
    }
]

response = llm.chat(
    messages,
    tools=functions
)
```

### Use Cases

| Use Case | Example |
|----------|---------|
| API calls | Get weather, search, etc. |
| Database queries | Search your data |
| Calculations | Run computations |
| File operations | Read/write files |

---

## JSON Mode

### What It Is

JSON mode = Force LLM to output valid JSON.

```
User: Give me a list of colors

Without JSON mode:
"Here are some colors: red, blue, green"

With JSON mode:
{"colors": ["red", "blue", "green"]}
```

### How It Works

```python
response = llm.chat(
    prompt,
    response_format={"type": "json_object"}
)
```

The schema is included in the context - just text telling the LLM what format to output:

```
[Context sent to LLM]
System: You are a helpful assistant.
Output format: JSON

User: Give me a color

LLM generates: {"color": "blue"}
```

### With Schema

For stricter output, provide a schema:

```python
response = llm.chat(
    prompt,
    response_format={
        "type": "json_object",
        "schema": {
            "type": "object",
            "properties": {
                "color": {"type": "string"},
                "hex": {"type": "string"},
                "rgb": {
                    "type": "object",
                    "properties": {
                        "r": {"type": "integer"},
                        "g": {"type": "integer"},
                        "b": {"type": "integer"}
                    }
                }
            },
            "required": ["color"]
        }
    }
)
```

### Limitations

- Still can produce invalid JSON (model may "hallucinate" JSON structure)
- Must provide JSON schema for best results
- Not all models support it
- Schema is just in context - LLM may not follow exactly

### JSON Mode vs Function Calling

| Feature | JSON Mode | Function Calling |
|---------|-----------|------------------|
| Execution | No | Yes |
| Strictness | Can still be invalid | More reliable |
| Use case | Get structured data | Trigger actions |

---

## Function Calling vs JSON Mode

| Aspect | Function Calling | JSON Mode |
|--------|-----------------|-----------|
| **Purpose** | Execute code | Structured output |
| **Trigger** | Model decides | User requests |
| **Schema** | Tool definitions | Output format |
| **Execution** | Yes (your code runs) | No |

---

## Relationship with MCP

```
Function calling = LLM outputs structured data (mechanism)

MCP = Standard protocol for defining functions (standardization)
```

They work together:
- MCP defines tools in a standard way
- Function calling is how LLM requests to call them

But they're NOT the same thing.

---

## Tools Without MCP

You can have function calling without MCP:

```python
# Define your own functions (not MCP)
functions = [
    {
        "name": "my_function",
        "parameters": {...}
    }
]

# LLM outputs JSON → your code executes
```

This is "raw" function calling - works but you define everything manually.

---

## Key Takeaways

1. **Function calling** = LLM outputs JSON → your code executes
2. **JSON mode** = Force valid JSON output
3. **MCP** = Standard way to define functions
4. **MCP uses function calling** but they're different concepts

---

## References

- OpenAI function calling docs
- Anthropic tool use docs
- MCP spec
