# Agents

> Full agent architectures - how LLMs become autonomous systems that can plan, execute, and learn.

---

## What Is an Agent

An agent = LLM + tools + memory + loop

```
Agent = LLM (brain) + Tools (hands) + Memory (context) + Loop (autonomy)
```

Unlike a simple LLM that just responds once, an agent can:
- Plan multi-step tasks
- Use tools autonomously
- Make decisions based on results
- Learn from feedback

---

## Agent Architecture

```
┌─────────────────────────────────────────────────────────┐
│                      Agent                               │
├─────────────────────────────────────────────────────────┤
│  ┌─────────────┐                                        │
│  │  LLM        │ ← The "brain"                        │
│  └─────────────┘                                        │
│         ↓                                               │
│  ┌─────────────┐                                        │
│  │  Planner    │ ← Breaks tasks into steps             │
│  └─────────────┘                                        │
│         ↓                                               │
│  ┌─────────────┐                                        │
│  │  Executor   │ ← Calls tools, takes actions           │
│  └─────────────┘                                        │
│         ↓                                               │
│  ┌─────────────┐                                        │
│  │  Observer   │ ← Checks results, decides next step    │
│  └─────────────┘                                        │
│         ↓                                               │
│  ┌─────────────┐                                        │
│  │  Memory    │ ← Stores context, history             │
│  └─────────────┘                                        │
└─────────────────────────────────────────────────────────┘
```

---

## The Agent Loop

```
1. Receive task
        ↓
2. Plan (LLM decides steps)
        ↓
3. Execute (call tools)
        ↓
4. Observe (get results)
        ↓
5. Decide (continue or done?)
        ↓
6. Repeat until complete
```

### Example

```
User: "Create a todo app"

Agent loop:
1. Plan: Need to create files, set up project
2. Execute: Create SPEC.md
3. Observe: Created successfully
4. Plan: Next, create main.go
5. Execute: Write main.go
6. Observe: Created
7. Plan: Need to add dependencies
8. Execute: Run go mod init
... (continues until done)
```

---

## Types of Agents

### 1. ReAct Agents

Think (reason) → Act (use tools) → Observe (get results)

```python
while not done:
    thought = llm.think(context)
    action = llm.decide_action(thought)
    result = execute(action)
    context += result
```

### 2. Plan-and-Execute

Plan first, then execute step by step

```
1. LLM creates full plan
2. Execute each step
3. Can replan if needed
```

### 3. Reflexive Agents

Can stop/adapt based on results

```
If step fails → replan
If tool error → try alternative
If context too long → summarize
```

### 4. Multi-Agent

Multiple agents work together

```
Agent A (planner) → Agent B (coder) → Agent C (reviewer)
```

---

## Tools in Agents

Agents need tools to be useful:

| Tool | Purpose |
|------|---------|
| Read/Write files | Interact with filesystem |
| Bash | Run commands |
| Search | Find information |
| Browser | Navigate web |
| MCP | Connect to external services |

---

## Key Components

### Planning

```python
def plan(task, context):
    return llm.chat(f"""
    Task: {task}
    History: {context}
    
    What's the next step? List concrete actions.
    """)
```

### Execution

```python
def execute(action):
    tool_name, args = parse(action)
    return tools[tool_name](args)
```

### Feedback Loop

```python
def should_continue(result, max_steps):
    if result.is_error:
        return "replan"
    if result.success:
        return "continue"
    if step_count > max_steps:
        return "stop"
    return "continue"
```

---

## Agent Frameworks

| Framework | Description |
|-----------|-------------|
| LangChain | Python, popular |
| LangGraph | Graph-based workflows |
| AutoGen | Multi-agent Microsoft |
| Claude Agent SDK | Anthropic's |
| OpenAI Agents | OpenAI's |
| CrewAI | Multi-agent orchestrator |

---

## Challenges

| Challenge | Description |
|-----------|-------------|
| Planning failure | LLM can't plan complex tasks |
| Tool errors | Wrong tool use |
| Loops | Agent gets stuck repeating |
| Context limits | Too much history |
| Reliability | Not deterministic |

---

## Key Takeaways

1. **Agent** = LLM + tools + loop + memory
2. **Loop** = Plan → Execute → Observe → Decide
3. **Tools** = What agent can do
4. **Memory** = What agent remembers
5. **Multi-agent** = Multiple agents collaborating

---

## References

- LangChain docs
- ReAct paper: "ReAct: Synergizing Reasoning and Acting in Language Models"
- AutoGen docs
- Anthropic Agent SDK
