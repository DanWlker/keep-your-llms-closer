# Commands - Technical Guide

> Quick-triggered custom instructions via slash commands.

---

## Overview

Commands are **saved prompts** stored as markdown files. They become slash commands you can invoke with `/command-name`.

**Key point:** Commands are essentially prompts stored in files. When invoked, their content is injected into the conversation.

---

## File Structure

```
.claude/
├── commands/
│   ├── review.md        → /review
│   ├── test.md         → /test
│   └── git/
│       └── commit.md   → /git:commit
```

**Location determines scope:**
- `~/.claude/commands/` - Global (all projects)
- `.claude/commands/` - Project-specific

---

## Format

```markdown
# Command name (optional)

Instructions for Claude...

Run these steps:
1. Do this
2. Then do that
3. Report results
```

### With Arguments

```markdown
# Code Review

Review the following code for security issues:
$ARGUMENTS

Be specific about:
- SQL injection
- XSS vulnerabilities
- Hardcoded secrets
```

**Usage:** `/review src/auth/login.go`

---

## How Tools Use Commands

### 1. Invocation (All Tools)

```
User types: /review

Claude reads: .claude/commands/review.md

Content is injected into conversation as instructions
```

**This is it.** Just reads the file and uses content as prompt.

---

### 2. Built-in Commands

| Command | Description |
|---------|-------------|
| `/clear` | Fresh start |
| `/compact` | Compress history |
| `/help` | Show commands |
| `/config` | Show configuration |
| `/mcp` | MCP servers |
| `/review` | Request code review |
| `/exit` | End session |

---

### 3. Parameterized Commands

```markdown
# Create Component

Create a new React component at $ARGUMENTS

Requirements:
- Use TypeScript
- Include props interface
- Add basic styles
```

**Usage:** `/create-component components/Button.tsx`

---

## Comparison with Skills

| Aspect | Commands | Skills |
|--------|----------|--------|
| Trigger | Explicit `/command` | Auto-detect via description |
| Scope | Single prompt | Full workflow |
| Parameters | Via `$ARGUMENTS` | Via skill instructions |
| Status | Merged into Skills (recent) | Active development |

---

## OpenCode Commands vs Claude Code Commands

These are **different** things:

| Aspect | OpenCode | Claude Code |
|--------|----------|-------------|
| **Type** | TUI commands (application-level) | Prompt templates (LLM context) |
| **Built-in** | `/init`, `/help`, `/new`, `/undo`, `/redo`, `/exit` | `/clear`, `/compact`, `/review`, `/mcp` |
| **Custom** | `.opencode/commands/` | `.claude/commands/` |
| **What they do** | Application actions | Content injected into LLM prompt |

### OpenCode Slash Commands

OpenCode's `/` commands are **TUI shortcuts** for application actions:

- `/init` - Create/update AGENTS.md
- `/help` - Show help
- `/new` / `/clear` - New session
- `/undo` / `/redo` - Undo/redo messages
- `/exit` - Exit

These are **not** the same as Claude Code's Commands which inject content into the LLM prompt.

---

## Key Takeaways

1. **Commands are prompts in files** - Content injected when invoked
2. **Explicit trigger** - User must type `/command`
3. **Simple mechanism** - Just reads .md file
4. **Recently merged** - In new Claude Code versions, Commands merged into Skills

---

## References

- Claude Code docs
- Custom commands stored in `.claude/commands/`
- OpenCode commands: https://opencode.ai/docs/commands/
