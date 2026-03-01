# MCP (Model Context Protocol) - Technical Guide

> Connects agents to external tools and data sources. Real code execution, not LLM context.

---

## Overview

MCP (Model Context Protocol) is an **open protocol** for connecting AI agents to external tools, databases, and services. It's the "USB-C for AI" - a standardized way to integrate any data source without custom glue code.

**Key point:** MCP enables real tool execution. Agents can call actual functions/APIs, not just context in prompts.

---

## Architecture

```
┌─────────────────────────────────────────────────────────┐
│                     Host (AI Agent)                      │
│  ┌─────────────────────────────────────────────────┐    │
│  │               MCP Client(s)                     │    │
│  └─────────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────┘
                           │
              JSON-RPC 2.0 (over stdio / SSE / HTTP)
                           │
┌─────────────────────────────────────────────────────────┐
│                      MCP Server                          │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐             │
│  │  Tools   │  │Resources │  │ Prompts  │             │
│  └──────────┘  └──────────┘  └──────────┘             │
└─────────────────────────────────────────────────────────┘
```

---

## Three Components

### 1. Hosts

The AI application that initiates connections.

- Examples: Claude Desktop, Claude Code, Cursor, Codex, Copilot
- Manages MCP client lifecycle
- Enforces security policies

### 2. Clients

Connectors within the host application.

- One client per MCP server connection
- Handles protocol communication
- Manages stateful sessions

### 3. Servers

Standalone programs that expose capabilities.

- One server per integration (GitHub, Slack, database, etc.)
- Provides tools, resources, prompts
- Can be local (stdio) or remote (HTTP/SSE)

---

## Three Primitives (What Servers Expose)

### Tools

Executable functions the AI can call.

```json
{
  "name": "get_weather",
  "description": "Get weather for a location",
  "inputSchema": {
    "type": "object",
    "properties": {
      "location": { "type": "string" }
    }
  }
}
```

### Resources

Data the AI can read (files, DB records, API responses).

```json
{
  "uri": "file:///src/main.go",
  "name": "Main Go file",
  "mimeType": "text/plain"
}
```

### Prompts

Reusable prompt templates.

```json
{
  "name": "review_code",
  "description": "Review code for issues",
  "arguments": [{ "name": "file", "type": "string" }]
}
```

---

## Protocol Details

### Transport Mechanisms

| Type | Use Case |
|------|----------|
| **stdio** | Local servers (child process) |
| **HTTP + SSE** | Remote servers |
| **WebSocket** | Real-time bidirectional |

### Message Format

All messages use **JSON-RPC 2.0**:

```json
// Request
{
  "jsonrpc": "2.0",
  "id": "req-1",
  "method": "tools/call",
  "params": {
    "name": "get_weather",
    "arguments": { "location": "Tokyo" }
  }
}

// Response
{
  "jsonrpc": "2.0",
  "id": "req-1",
  "result": {
    "content": [{ "type": "text", "text": "22°C, sunny" }]
  }
}
```

### Lifecycle

1. **Initialize** - Client and server negotiate protocol version
2. **Connected** - Server announces capabilities (tools, resources, prompts)
3. **Interact** - Client calls tools, reads resources
4. **Close** - Clean shutdown

---

## Configuration

### Claude Code

```json
// ~/.claude/settings.json
{
  "mcpServers": {
    "github": {
      "command": "npx",
      "args": ["-y", "@modelcontextprotocol/server-github"]
    },
    "filesystem": {
      "command": "npx",
      "args": ["-y", "@modelcontextprotocol/server-filesystem", "/path/to/allowed"]
    }
  }
}
```

### Environment Variables

```json
{
  "mcpServers": {
    "my-server": {
      "command": "node",
      "args": ["/path/to/server.js"],
      "env": {
        "API_KEY": "secret"
      }
    }
  }
}
```

---

## Function Calling vs MCP

| Aspect | Function Calling | MCP |
|--------|-----------------|-----|
| Setup | Manual per function | Standard protocol |
| Discovery | Static | Dynamic |
| Flexibility | Limited to predefined | Discover new tools |
| Structure | Fixed inputs/outputs | Flexible schema |
| Integration | Custom glue code | Protocol-based |

**MCP standardizes how tools are served and discovered, not just what they do.**

---

## Popular MCP Servers

| Server | Purpose |
|--------|---------|
| `@modelcontextprotocol/server-github` | GitHub API |
| `@modelcontextprotocol/server-filesystem` | File system access |
| `@modelcontextprotocol/server-brave-search` | Web search |
| `@modelcontextprotocol/server-postgres` | PostgreSQL |
| `@modelcontextprotocol/server-slack` | Slack API |
| `@modelcontextprotocol/server-aws-kb-retrieval` | AWS knowledge base |

**Full list:** https://github.com/modelcontextprotocol/servers

---

## Security

### Current State (as of 2025)

- **No built-in auth** in core spec
- **RFC-9728** adds OAuth2-based authorization
- You must secure MCP servers yourself

### Best Practices

- Restrict filesystem access to specific directories
- Use environment variables for secrets (not hardcoded)
- Consider network isolation for remote servers
- Review MCP server code before use

---

## MCP vs Hooks

| Aspect | MCP | Hooks |
|--------|-----|-------|
| **Type** | External tool integration | Workflow automation |
| **Direction** | Agent → External service | Around agent events |
| **What it does** | Calls APIs, reads data | Runs shell commands |
| **Scope** | Any external system | Local machine |
| **Trigger** | When agent needs data | When events fire |

---

## How MCP Actually Works (Technical)

### The Misconception

LLMs don't "execute" anything. They're text generators. The magic is in the **agent framework** that wraps the LLM.

### The Real Flow

```
1. User prompt → Agent framework
         ↓
2. Framework sends prompt + tool schemas to LLM
         ↓
3. LLM generates text INCLUDING structured tool call
         ↓
4. Framework PARSES the structured output (JSON parsing, not regex)
         ↓
5. Framework EXECUTES the actual MCP tool call
         ↓
6. Result → LLM as context
         ↓
7. LLM generates response with results
```

### What the LLM Actually Sees

The LLM doesn't execute MCP directly. It sees **tool schemas**:

```
You have access to these tools:

function get_weather(location: string) → weather data
function github_create_issue(owner, repo, title) → issue

User: What's the weather in Tokyo?
```

The LLM generates structured output:
```
I'll check the weather for Tokyo.

<tool_call>
{"name": "get_weather", "arguments": {"location": "Tokyo"}}
</tool_call>
```

The **framework** (not the LLM) parses this and executes the actual HTTP request to the MCP server.

---

## MCP vs Hooks: Technical Comparison

### Hooks (Event-Driven)

```
User types something
    ↓
Framework detects "Write tool is about to be used" (PreToolUse event)
    ↓
Framework checks if any hooks match this event
    ↓
If match → Execute shell command (sync or async)
    ↓
Return result → Continue or BLOCK (exit code 2)
    ↓
LLM never knew about it
```

**Key point:** Hooks operate on **framework events**, not LLM output. The LLM is unaware.

### MCP (LLM-Driven)

```
User: "Check the weather"

Framework → LLM: "Here are available tools: get_weather, etc."

LLM → Framework: "I'll call get_weather({location: 'Tokyo'})"

Framework: "Parsing tool call (JSON) → executing HTTP request to weather API..."

API Response → LLM: "22°C, sunny"

LLM → User: "It's 22°C and sunny in Tokyo."
```

**Key point:** MCP operates on **LLM-generated structured output**. The LLM decides WHEN to call based on tool schemas.

### Both are Framework-Level

| Aspect | Hooks | MCP |
|--------|-------|-----|
| **Trigger** | Framework event (PreToolUse, PostToolUse, etc.) | LLM generates structured tool call |
| **LLM involvement** | None - framework decides | LLM decides to call a tool |
| **Parsing** | Event matching | JSON parsing of tool call |
| **Can block** | Yes - exit code 2 | No - LLM decides |
| **Execution** | Shell commands | HTTP/API calls |

---

## Key Takeaways

1. **Real execution** - MCP calls actual functions/APIs, not prompts
2. **Client-server** - Standard protocol over JSON-RPC
3. **Three primitives** - Tools (execute), Resources (read), Prompts (templates)
4. **Transport** - stdio (local) or HTTP/SSE (remote)
5. **Dynamic discovery** - Agent discovers available tools at runtime
6. **Security gap** - No built-in auth (use env vars, network isolation)

---

## References

- **Spec:** https://modelcontextprotocol.io/specification
- **Docs:** https://modelcontextprotocol.io/docs
- **Servers:** https://github.com/modelcontextprotocol/servers
- **Python SDK:** https://github.com/modelcontextprotocol/python-sdk
- **TypeScript SDK:** https://github.com/modelcontextprotocol/typescript-sdk
