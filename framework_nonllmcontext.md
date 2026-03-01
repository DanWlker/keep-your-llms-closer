# Agent Framework - Technical Guide

> The orchestration layer that coordinates LLM, tools, subagents, plugins, and everything else.

---

## Overview

The **Agent Framework** (also called Runtime or Orchestrator) is the glue that makes an "agent" out of an LLM. Without it, you'd just have a text generator.

The framework:
- Parses LLM structured output
- Executes tools (including MCP)
- Spawns subagents
- Loads plugins/skills/commands
- Handles hooks
- Manages context windows

**Key point:** The framework is NOT the LLM. The LLM just generates text. The framework does everything else.

---

## What Each Component Does

| Component | Who loads it | Goes to LLM? | Executed by Framework? |
|-----------|-------------|--------------|----------------------|
| **LLM** | Framework | N/A | No (generates text) |
| **Framework** | You | No | Yes (orchestrates) |
| **MCP servers** | Framework | No (tools) | Yes (HTTP calls) |
| **Plugins** | Framework | Depends | Registers components |
| **Skills** | Framework | Yes | No |
| **Commands** | Framework | Yes | No |
| **Subagents** | Framework | Spawns new | Yes |
| **Hooks** | Framework | No | Yes (shell) |

---

## The Execution Flow

```
User types message
        ↓
Framework prepares context:
  - User message
  - CLAUDE.md (project context)
  - Skills (auto-matched by description)
  - Available tools (Read, Write, Bash, MCP tools)
        ↓
Framework → LLM: "Generate response"
        ↓
LLM generates text (may include structured tool calls)
        ↓
Framework parses output:
  - Tool call? → Execute tool, return result
  - Subagent spawn? → New LLM with agent config
  - Hook event? → Run hook script
        ↓
Result → Back to LLM
        ↓
Framework → User: Final response
```

---

## Framework Architecture

```
┌─────────────────────────────────────────────────────────┐
│                     Your Application                      │
├─────────────────────────────────────────────────────────┤
│                  Agent Framework                         │
│  ┌─────────────────────────────────────────────────┐  │
│  │              Orchestration Layer                  │  │
│  │  - Parse LLM output (structured JSON)            │  │
│  │  - Route to appropriate handlers                 │  │
│  │  - Manage context windows                       │  │
│  │  - Handle state/conversation                    │  │
│  └─────────────────────────────────────────────────┘  │
│           ↓         ↓         ↓         ↓             │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌─────────┐ │
│  │  Tools   │ │ Subagents│ │  Hooks   │ │ Skills  │ │
│  │ Handler  │ │ Spawner  │ │ Executor │ │ Registry│ │
│  └──────────┘ └──────────┘ └──────────┘ └─────────┘ │
│           ↓         ↓         ↓                        │
│  ┌─────────────────────────────────────────────────┐ │
│  │              MCP Client                          │ │
│  │  - Connect to MCP servers                        │ │
│  │  - Call tools, read resources                   │ │
│  └─────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────┘
                          ↓
              ┌────────────────────────┐
              │      External          │
              │  - MCP Servers         │
              │  - Shell commands      │
              │  - APIs               │
              └────────────────────────┘
```

---

## What Happens When You Install a Plugin

```
You: /plugin install owner/plugin

Framework:
1. Download plugin from GitHub/marketplace
2. Read .claude-plugin/plugin.json
3. For each component:
   ├── Skills → Add to skill registry
   ├── Commands → Add to command registry
   ├── Agents → Add to agent pool
   ├── Hooks → Merge with existing hooks
   └── MCP → Register MCP servers
4. Done. Components now available in runtime.
```

---

## What Happens When LLM Generates a Tool Call

```
LLM generates:
<tool_call>{"name": "Read", "arguments": {"file": "main.go"}}</tool_call>

Framework parses:
1. Sees "Read" tool call
2. Looks up tool in registry
3. Executes: os.ReadFile("main.go")
4. Returns: file content
5. Appends to conversation "Tool: result: [content]"
6. Sends back to LLM: "Continue generating"
```

---

## What Happens When Subagent Spawns

```
LLM generates:
<tool_call>{"name": "Task", "arguments": {"subagent": "code-reviewer", "prompt": "Review main.go"}}

Framework:
1. Sees "Task" tool call with subagent
2. Loads code-reviewer agent config
3. Spawns NEW LLM call with:
   - Agent's system prompt
   - Task description
   - Agent's allowed tools
4. Subagent works independently
5. Returns only final result to main
6. Main continues conversation
```

---

## Popular Frameworks

| Framework | Used by |
|-----------|---------|
| Claude Code runtime | Claude Code, CLI |
| Anthropic Agent SDK | Custom agents |
| LangChain | Python agents |
| LangGraph | Complex workflows |
| OpenAI Agents SDK | OpenAI agents |

---

## Framework vs MCP

| Aspect | Framework | MCP |
|--------|-----------|-----|
| **What it is** | Orchestrator | External tool protocol |
| **Scope** | Everything | Just external tools |
| **Who uses MCP** | Framework | Connects to external |
| **Alternative** | N/A | Custom function calling |

---

## Key Takeaways

1. **Framework orchestrates** - It's the thing that does everything
2. **LLM is dumb** - Just generates text, doesn't execute
3. **Framework is smart** - Parses, routes, executes, manages state
4. **MCP is one input** - Framework connects to MCP servers
5. **Plugins are bundles** - Framework loads and registers components

---

## References

- Claude Code architecture
- Anthropic Agent SDK docs
- MCP specification (framework connects to it)
