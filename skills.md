# Skills (SKILL.md) - Technical Guide

> Deep dive into the Agent Skills format, how it works, and how tools use it.

---

## Overview

Skills are **plain markdown files** that get included in the LLM's context. There is no special execution engine or validation at runtime - the content is simply concatenated into the system prompt.

**Key point:** Skills are just text in context. Tools don't "run" skills - they include them as instructions when passing to the LLM.

---

## File Structure

```
my-skill/
├── SKILL.md              # Required: metadata + instructions
├── scripts/              # Optional: executable scripts
├── references/           # Optional: additional docs
└── assets/              # Optional: static files
```

---

## SKILL.md Format

```markdown
---
name: skill-name
description: What this skill does and when to use it.
version: "1.0.0"
license: MIT
compatibility:
  - Claude Code
  - Cursor
  - Codex
allowed-tools: Read, Write, Edit, Glob, Grep, Bash
user-invocable: true
---

# Skill Title

Instructions for the agent...

## When to Use
- Use when user wants to do X
- Use when task involves Y

## Instructions
1. Step one
2. Step two

## Output Format
Format the output as...
```

---

## Frontmatter Fields

| Field | Required | Description |
|-------|----------|-------------|
| `name` | Yes | 1-64 chars, lowercase letters/numbers/hyphens only |
| `description` | Yes | 1-1024 chars. **Critical for activation** - used for routing |
| `version` | No | Semantic version string |
| `license` | No | SPDX license identifier |
| `compatibility` | No | List of supported agents |
| `allowed-tools` | No | Comma-separated list of permitted tools |
| `user-invocable` | No | Whether users can trigger directly |
| `context` | No | Execution context (e.g., `fork`) |

### Name Validation Rules

- Max 64 characters
- Lowercase letters (a-z), numbers (0-9), hyphens (-) only
- No consecutive hyphens (`--`)
- Must match parent directory name

---

## How Tools Use Skills

### 1. Context Inclusion (All Tools)

The primary mechanism - skills are just text in the prompt:

```
System: You have access to skills...
Skills:
- skill-a: Does X (when to use)
- skill-b: Does Y (when to use)

When skill matches task, read: /path/to/skill/SKILL.md
```

When activated, full SKILL.md content is appended to context.

**This is it.** No validation. No execution engine. Just string concatenation.

---

### 2. Progressive Disclosure

Three-tier loading to manage context size:

| Tier | What's Loaded | Trigger |
|------|---------------|---------|
| 1. Discovery | `name` + `description` (~100 tokens) | Always visible in skill list |
| 2. Activation | Full SKILL.md body (<5,000 tokens) | Task matches description |
| 3. Resources | Files in `scripts/`, `references/` | Referenced in instructions |

---

### 3. Tool-Specific Behaviors

| Tool | Storage Location | Special Behavior |
|------|------------------|------------------|
| Claude Code | `~/.claude/skills/` | Auto-discovers |
| Codex CLI | `~/.codex/skills/` | `--enable-skills` flag |
| Cursor | Project-level `.cursor/rules/` | Integrated with rules |
| Copilot | `.claude/skills/` | Auto-loads |
| Gemini CLI | `~/.gemini/skills/` | Recent support |

---

### 4. Optional: Script Execution

Some tools can execute scripts in `scripts/` directory:

- Scripts are **not** read into LLM context
- Agent can invoke them via tool calls
- Example: `scripts/validate.sh` runs tests without bloating context

**Most tools don't do this.** It's optional behavior.

---

## Validation

### By Tools

Most tools **don't validate** skills. They just include them.

Some libraries (e.g., LM-Kit.NET) offer validation:
```python
SkillRegistry.ValidateAll()  # Checks spec compliance
```

### Manual Validation

You can validate yourself:
- Frontmatter is valid YAML
- `name` matches directory
- `description` is present and descriptive
- No XML angle brackets in frontmatter (`<` or `>`)

---

## Activation Triggers

The model decides when to use a skill based on:

1. **Description** - The `description` field in frontmatter
2. **User request** - Natural language match
3. **Explicit invocation** - User asks for the skill by name

**Key insight:** The `description` is critical. It's the only thing shown during discovery. If the model doesn't match the task to the description, the skill won't activate.

---

## Example: How Claude Code Loads a Skill

```python
# Simplified conceptual flow
def handle_user_message(user_input):
    # 1. Show available skills (just metadata)
    skill_list = load_all_skill_metadata()
    
    # 2. Model decides skill matches
    if model_decides_to_use_skill(user_input, skill_list):
        # 3. Load full SKILL.md into context
        skill_content = read_file(skill.path + "/SKILL.md")
        context += skill_content
        
        # 4. Load any referenced resources
        if "references/style.md" in skill_content:
            context += read_file(skill.path + "/references/style.md")
    
    # 5. Generate response
    return model.generate(context)
```

**That's it.** It's all just string manipulation.

---

## Key Takeaways

1. **Skills are context, not code** - Just markdown text passed to LLM
2. **No execution engine** - Nothing runs; instructions are followed by the model
3. **Description is critical** - Used for routing/activation
4. **Progressive disclosure** - Manages context size
5. **Scripts are optional** - Only some tools execute them
6. **Validation is rare** - Most tools don't validate; they just include

---

## References

- **Spec:** https://agentskills.io/specification
- **Marketplace:** https://skills.sh
- **Validation lib:** LM-Kit.NET (optional validation)
