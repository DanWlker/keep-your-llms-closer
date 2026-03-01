# Hooks - Technical Guide

> Execute code at specific points in the agent's workflow.

---

## Overview

Hooks are **shell commands** (or LLM prompts) that run automatically at specific points in Claude Code's lifecycle. Unlike Skills which are advisory, hooks provide **deterministic control** - they always execute.

**Key point:** Hooks execute real code at defined events. They can block actions, format files, run tests, send notifications, etc.

---

## How Hooks Work

```
Claude Code Session
       ↓
Event fires (e.g., PostToolUse)
       ↓
Matcher checks if hook applies
       ↓
Hook executes (shell command / HTTP / LLM prompt)
       ↓
Return JSON (approve/block/modify)
       ↓
Claude continues or stops
```

---

## Hook Events (14 Lifecycle Points)

| Event | When It Fires |
|-------|---------------|
| `Start` | Session begins |
| `Stop` | Session ends |
| `ToolUse` | Before a tool runs |
| `PostToolUse` | After a tool runs |
| `Notification` | When Claude needs to notify user |
| `Stop` | Session ends |
| `StopSubagent` | Subagent completes |
| `Pause` | Session paused |
| `Resume` | Session resumed |
| `PreRead` | Before reading files |
| `PostRead` | After reading files |
| `PreEdit` | Before editing |
| `PostEdit` | After editing |
| `PermissionRequest` | When Claude requests permission |

---

## Hook Types

### 1. Command Hooks (Shell)

```json
{
  "hooks": {
    "PostToolUse": [
      {
        "matcher": "Edit|Write",
        "command": "prettier --write ${event.path}"
      }
    ]
  }
}
```

### 2. HTTP Hooks

```json
{
  "hooks": {
    "Stop": [
      {
        "matcher": ".*",
        "url": "https://your-server.com/webhook",
        "method": "POST"
      }
    ]
  }
}
```

### 3. Prompt Hooks (LLM)

```json
{
  "hooks": {
    "Start": [
      {
        "matcher": ".*",
        "prompt": "You are a CI bot. Check if the following git diff has tests: ${gitDiff}"
      }
    ]
  }
}
```

---

## JSON Input/Output

### Input (via stdin)

```json
{
  "event": "PostToolUse",
  "path": "/path/to/file.go",
  "tool_name": "Write",
  "tool_input": { ... }
}
```

### Output (control behavior)

```json
// Allow (continue)
{ "ok": true }

// Block (stop execution)
{ "ok": false, "reason": "Security violation" }

// Modify (change behavior)
{ "ok": true, "override": { ... } }
```

### Exit Codes

| Code | Effect |
|------|--------|
| 0 | Allow / Success |
| 1 | Error (log only) |
| 2 | Block the action |

---

## Async Hooks (Jan 2026)

```json
{
  "matcher": ".*",
  "command": "notify-send 'Done'",
  "async": true
}
```

Run hooks in background without blocking Claude.

---

## Configuration Locations

```
~/.claude/settings.json     # Global (all projects)
.claude/settings.json       # Project-specific
```

---

## Common Patterns

### Auto-Format After Write

```json
{
  "hooks": {
    "PostToolUse": [
      {
        "matcher": "Edit|Write",
        "command": "prettier --write ${event.path}"
      }
    ]
  }
}
```

### Typecheck on Every Change

```json
{
  "hooks": {
    "PostToolUse": [
      {
        "matcher": "Write|Edit",
        "command": "bun run typecheck"
      }
    ]
  }
}
```

### Block Dangerous Commands

```json
{
  "hooks": {
    "ToolUse": [
      {
        "matcher": "Bash.*rm -rf",
        "command": "echo 'Blocked: rm -rf is dangerous'"
      }
    ]
  }
}
```

### Session Start Context

```json
{
  "hooks": {
    "Start": [
      {
        "matcher": ".*",
        "command": "cat .env.example"
      }
    ]
  }
}
```

---

## Hooks vs Skills vs Commands

| Aspect | Hooks | Skills | Commands |
|--------|-------|--------|----------|
| **Trigger** | Automatic event | Auto-detect | Explicit `/` |
| **Execution** | Real shell commands | LLM context | LLM context |
| **Deterministic** | Yes - always runs | No - model may ignore | No - model may ignore |
| **Blocking** | Can stop actions | Cannot | Cannot |
| **Purpose** | Automation/guardrails | Guidance | Reusable prompts |

---

## How Hooks Actually Work (Technical)

### The Key Insight

**Hooks operate on FRAMEWORK EVENTS, not LLM text output.**

The LLM never knows about hooks. They're purely a framework-level mechanism.

### The Real Flow

```
User types something in Claude Code
    ↓
Framework receives input → sends to LLM
    ↓
LLM generates text with potential tool call
    ↓
Framework decides to use "Bash" tool with command "rm -rf /"
    ↓
BEFORE executing (PreToolUse event fires)
    ↓
Framework checks: "Any hooks match this event?"
    ↓
Yes! Block hook matches "Bash.*rm -rf"
    ↓
Hook executes → returns exit code 2 (BLOCK)
    ↓
Framework BLOCKS the tool execution
    ↓
User sees: "Command blocked: rm -rf is dangerous"
    ↓
LLM is told the tool failed → generates apology/alternative
```

### How Hooks Stop Bad Actions

Hooks use **PreToolUse** event to intercept BEFORE execution:

```json
{
  "hooks": {
    "ToolUse": [
      {
        "matcher": "Bash.*rm -rf",
        "command": "echo 'Blocked'",
        "env": { "DECISION": "block" }
      }
    ]
  }
}
```

Exit codes control behavior:

| Exit Code | Effect |
|-----------|--------|
| 0 | Allow the tool to execute |
| 1 | Error (logs only, doesn't block) |
| 2 | BLOCK the tool execution |

### Prompt Hooks (LLM-Powered)

There's one exception: **Prompt hooks** can involve the LLM:

```json
{
  "hooks": {
    "Start": [
      {
        "matcher": ".*",
        "prompt": "Analyze this git diff: ${gitDiff}"
      }
    ]
  }
}
```

This asks the LLM to analyze something and return a decision. But the LLM in the hook is **separate** from the main conversation.

---

## Hooks vs MCP: Implementation Difference

| Aspect | Hooks | MCP |
|--------|-------|-----|
| **What triggers it** | Framework event | LLM generates tool call |
| **LLM involvement** | None (except prompt hooks) | Yes - decides to call |
| **What executes** | Shell commands | HTTP/API calls |
| **Blocking** | Exit code 2 | Cannot block LLM |
| **When it runs** | At defined events | When LLM decides |

### Both NOT on LLM Text

Neither hooks nor MCP work on "LLM text" in the way you might think:

- **Hooks**: Pure framework events. LLM never sees them.
- **MCP**: LLM generates structured JSON, framework parses and executes.

The framework is the key - it's what makes agents "agentic" beyond just text generation.

---

## Key Takeaways

1. **Real execution** - Hooks run actual shell commands, not just prompts
2. **Deterministic** - Always execute at defined events
3. **Can block** - Exit code 2 stops the action
4. **14 events** - Hook into full lifecycle
5. **Async support** - Non-blocking execution since Jan 2026

---

## References

- Claude Code Hooks docs: https://code.claude.com/docs/en/hooks
- Hooks guide: https://code.claude.com/docs/en/hooks-guide
