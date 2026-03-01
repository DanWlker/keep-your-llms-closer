# Subagents - Technical Guide

> Specialized AI assistants with their own context and tool permissions that the main agent can delegate tasks to.

---

## Overview

Subagents are **specialized AI "personas"** that the main agent can delegate tasks to. Each has its own:
- Context window (separate from main conversation)
- System prompt (specific instructions)
- Tool permissions (can be restricted)

**Key point:** Subagents are like hiring specialists for your team. The main agent orchestrates, subagents execute in their own isolated context.

---

## File Structure

```
.claude/
├── agents/                    # Project-level subagents
│   ├── code-review.md
│   ├── security-audit.md
│   └── docs-writer.md
~/.claude/agents/             # User-level (global) subagents
    ├── market-researcher.md
    └── data-analyst.md
```

---

## Subagent Format

```markdown
---
name: code-reviewer
description: Reviews code for bugs, security issues, and best practices.
model: claude-sonnet-4-20250514
tools:
  - Read
  - Glob
  - Grep
  - Bash:git
---

# Code Reviewer

You are an expert code reviewer focused on:
- Finding bugs before they ship
- Security vulnerabilities
- Performance issues
- Code clarity and maintainability

When reviewing:
1. Read the changed files thoroughly
2. Check for edge cases
3. Look for common vulnerability patterns
4. Suggest specific improvements

Always provide concrete examples in your feedback.
```

---

## Frontmatter Fields

| Field | Required | Description |
|-------|----------|-------------|
| `name` | Yes | Unique identifier (lowercase, hyphenated) |
| `description` | Yes | Used for automatic delegation routing |
| `model` | No | Specific model to use |
| `tools` | No | Array of allowed tools |
| `maxTurns` | No | Max conversation turns |

### Tool Restriction Example

```yaml
---
name: read-only-explorer
description: Explores codebase without making changes
tools:
  - Read
  - Glob
  - Grep
---
```

This subagent can ONLY read files - it cannot write or execute commands.

---

## How Subagents Work

### Delegation Flow

```
User: "Review the authentication code"

Main Agent:
  1. Analyzes request
  2. Matches against subagent descriptions
  3. Finds "code-reviewer" matches
  4. Spawns subagent with task

Subagent (isolated context):
  1. Receives task + its system prompt
  2. Works independently
  3. Returns only final results

Main Agent:
  1. Receives subagent results
  2. Synthesizes response
  3. Presents to user
```

### Automatic vs Explicit Invocation

**Automatic:**
```
User: "Check for SQL injection vulnerabilities"
Main Agent: "I'll use the security-auditor subagent for this..."
[Delegates to security-auditor]
```

**Explicit:**
```
User: "/subagent security-auditor Check for SQL injection"
[Directly invokes security-auditor]
```

---

## Use Cases

### 1. Context Isolation

Main chat won't get polluted with exploration context:
- Subagent does repo scan
- Only final results return to main thread
- Main context stays clean

### 2. Tool Restriction

Read-only agents for safe exploration:
```yaml
tools: [Read, Glob, Grep]
```

### 3. Parallel Execution

Multiple subagents can work simultaneously:
```
Main Agent:
  - Spawns market-researcher subagent
  - Spawns technical-analyst subagent
  - Both run in parallel
  - Results synthesized together
```

### 4. Specialized Expertise

Different subagents for different tasks:
- `security-auditor` - security focus
- `code-reviewer` - general review
- `docs-writer` - documentation
- `test-generator` - test creation

---

## Subagents vs Other Mechanisms

| Aspect | Subagents | Skills | Commands |
|--------|-----------|--------|----------|
| **Context** | Isolated (own window) | Shared (in main context) | Shared |
| **Invocation** | Auto or explicit | Auto (description match) | Explicit `/` |
| **Tool restriction** | Yes | No | No |
| **Parallel execution** | Yes | No | No |
| **Persistence** | No (per task) | No (per task) | No |

### Subagents vs MCP

| Aspect | Subagents | MCP |
|--------|-----------|-----|
| **What it is** | Another LLM | External API/tool |
| **Can it think** | Yes | No |
| **Tool access** | Defined in config | Defined by server |
| **Use case** | Specialized expertise | External integration |

---

## Practical Example

### Project Structure

```
.claude/
├── agents/
│   ├── explorer.md      # Repo exploration
│   ├── reviewer.md      # Code review
│   └── debugger.md      # Bug investigation
```

### explorer.md

```markdown
---
name: explorer
description: Explores codebase structure, finds files, understands architecture
tools: [Read, Glob, Grep]
---

You are a codebase explorer. Your job is to:
- Map directory structure
- Find relevant files
- Understand how components connect

Be thorough but efficient. Report file locations and relationships.
```

### Usage

```
User: "Find all authentication related files"
Main Agent: "I'll ask the explorer subagent to find this..."
[Delegates to explorer subagent]

Explorer returns:
- Found 3 auth files:
  - src/auth/login.go
  - src/auth/jwt.go
  - src/auth/middleware.go

Main Agent: "Found 3 authentication files..."
```

---

## Key Takeaways

1. **Isolated context** - Subagent has own context window, doesn't pollute main
2. **Tool restriction** - Can limit to read-only for safe exploration
3. **Automatic delegation** - Main agent chooses based on description
4. **Parallel execution** - Multiple subagents can work simultaneously
5. **Returns results only** - Main agent gets summary, not full conversation

---

## References

- Claude Code subagents docs
- Define in `.claude/agents/` (project) or `~/.claude/agents/` (global)
