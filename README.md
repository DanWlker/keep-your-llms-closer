# AI Agent Configuration & Capabilities Guide

> Overview of markdown configuration files and capability extensions for AI coding agents.

---

## 1. Markdown Configuration Files

Prioritized list of configuration files that tell AI agents how to work in your project.

### Priority Order

| Priority | File | Tool(s) | Location |
|----------|------|---------|----------|
| 1 | `AGENTS.md` or `AGENT.md` | Universal (Codex, Amp, etc.) | Project root |
| 2 | `CLAUDE.md` | Claude Code | Project/user/enterprise level |
| 3 | `.cursor/rules/*.mdc` | Cursor IDE | `.cursor/rules/` |
| 4 | `.cursorrules` | Cursor IDE (legacy) | Project root |
| 5 | `.github/copilot-instructions.md` | GitHub Copilot | `.github/` |
| 6 | `.windsurfrules` | Windsurf | Project root |

### Key Formats

**AGENTS.md** - Universal standard (July 2025)
- Vendor-neutral format
- Supported by: OpenAI Codex, Google Jules, Cursor, Amp, Factory
- One file that any agent can read

**CLAUDE.md** - Claude Code specific
- Hierarchy: Enterprise → User → Project → Directory
- First ~200 lines auto-loaded
- Supports `@file` imports

**.cursor/rules/*.mdc** - Modern Cursor
- YAML frontmatter with `description`, `globs`, `alwaysApply`
- More granular than legacy `.cursorrules`

**.github/copilot-instructions.md** - GitHub Copilot
- Supports path-specific rules via `.github/instructions/*.instructions.md`

---

## 2. Capability Extensions

Systems that extend what AI agents can do.

### Skills (SKILL.md)

Modular packages of instructions + tools.

- **Created:** Oct 2025 by Anthropic
- **Format:** `SKILL.md` folder with resources
- **Adopted by:** Claude Code, Codex CLI, Copilot, Cursor, Windsurf, Gemini CLI, 20+ agents
- **Storage:** `~/.claude/skills/` (Claude), `~/.codex/skills/` (Codex)
- **Marketplace:** skills.sh (280K+ skills)

**Install:**
```bash
npx skills-installer install @owner/skill-name
```

**Top skills:**
| Skill | Installs | Purpose |
|-------|----------|---------|
| find-skills | 322K | Discovers other skills |
| vercel-react-best-practices | 167K | React/Next.js patterns |

---

### MCP (nonllmcontext/mcp.md)

Connects agents to external tools/data sources - the "USB-C for AI".

- **Created:** Nov 2024 by Anthropic
- **Purpose:** Standardized connection to databases, APIs, filesystems
- **Format:** JSON-RPC 2.0 over stdio/SSE/HTTP
- **Adopted by:** Claude, Codex, Copilot, Cursor, most agents
- **Note:** Real tool execution via client-server architecture, not LLM context

**Architecture:**
```
AI Agent (MCP Client)
       ↕
MCP Server (external tool/service)
```

---

### Commands

Quick-triggered custom instructions.

- **Location:** `.claude/commands/`
- **Format:** Markdown files starting with `/`
- **Note:** Now merged into Skills in recent Claude Code versions

---

### Hooks (nonllmcontext/hooks.md)

Execute real code at specific points in agent workflow (not LLM context).

- **Location:** `.claude/hooks/` (in `settings.json`)
- **Use:** Pre/post task automation, notifications, blocking actions
- **Note:** These are shell commands that execute, not prompts for the LLM

---

### Subagents (subagents.md)

Delegate complex tasks to specialized agent instances with isolated context.

- **Location:** `.claude/agents/` (project) or `~/.claude/agents/` (global)
- **Use:** Spawn focused agents, tool restriction, parallel execution
- **Note:** Has own context window, tool restrictions, returns only results

---

### Plugins (nonllmcontext/plugins.md)

Bundles of commands + skills + MCP + hooks for distribution.

- **Location:** Installed via `/plugin install`
- **Use:** Team standardization, redistributable packages
- **Note:** Distribution layer - bundles other extensions for sharing

---

## 3. How They Work Together

### Capability Hierarchy

```
CLAUDE.md (project context)
        ↓
      Skills (reusable patterns)
        ↓
    Commands (custom triggers)
        ↓
      MCP (external tools)
        ↓
    Hooks (workflow automation)
        ↓
    Subagents (task delegation)
        ↓
    Plugins (bundled distribution)
```

### When to Use What

| Scenario | Use | Detailed Notes |
|----------|-----|----------------|
| Project conventions | CLAUDE.md | |
| Reusable patterns | Skills | [skills.md](skills.md) |
| Custom triggers | Commands | [commands.md](commands.md) |
| External integrations | MCP | [nonllmcontext/mcp.md](nonllmcontext/mcp.md) |
| Workflow automation | Hooks | [nonllmcontext/hooks.md](nonllmcontext/hooks.md) |
| Complex delegation | Subagents | [subagents.md](subagents.md) |
| Team distribution | Plugins | [nonllmcontext/plugins.md](nonllmcontext/plugins.md) |

---

## 4. LLM Concepts

Additional knowledge about how LLMs and agents work.

### [nonllmcontext/framework.md](nonllmcontext/framework.md)

The orchestration layer that coordinates LLM, tools, subagents, plugins, and everything else.

### [nonllmcontext/rag.md](nonllmcontext/rag.md)

Retrieval Augmented Generation - how LLMs access external knowledge through vector databases.

### [nonllmcontext/llm_parameters.md](nonllmcontext/llm_parameters.md)

Controls that affect how LLMs generate output (temperature, top-p, etc.).

### [nonllmcontext/system_tokens_hallucinations.md](nonllmcontext/system_tokens_hallucinations.md)

System prompts, token limits, and hallucinations - core LLM concepts.

### [nonllmcontext/finetuning.md](nonllmcontext/finetuning.md)

Fine-tuning - customizing model behavior with your own data.

### [nonllmcontext/lora.md](nonllmcontext/lora.md)

LoRA / QLoRA - efficient fine-tuning with low-rank adapters and quantization.

### [nonllmcontext/function_calling.md](nonllmcontext/function_calling.md)

Function calling - structured JSON output for tools, and how it relates to MCP.

### [nonllmcontext/determinism.md](nonllmcontext/determinism.md)

Seed / determinism - reproducible outputs from LLMs.

### [nonllmcontext/prompt_injection.md](nonllmcontext/prompt_injection.md)

Prompt injection - security vulnerability via malicious prompts.

### [nonllmcontext/guardrails.md](nonllmcontext/guardrails.md)

Guardrails - filtering unsafe outputs and inputs.

---

## 5. References

- **AGENTS.md spec:** https://github.com/agentmd/agent.md
- **Skills spec:** https://agentskills.io
- **Skills marketplace:** https://skills.sh
- **MCP spec:** https://modelcontextprotocol.io
- **MCP servers:** https://github.com/modelcontextprotocol/servers
- **Claude Plugins registry:** https://claude-plugins.dev
- **Cursor Rules docs:** https://cursor.com/docs/context/rules
- **Copilot instructions:** https://docs.github.com/copilot/customizing-copilot/adding-custom-instructions

---

## 6. Timeline

| Date | Milestone |
|------|-----------|
| Nov 2024 | Anthropic introduces MCP |
| Jul 2025 | AGENTS.md standard published (Sourcegraph Amp team) |
| Oct 2025 | Anthropic introduces Skills for Claude |
| Dec 2025 | Open standard for Skills published at agentskills.io |
| Dec 2025 | OpenAI and Microsoft adopt SKILL.md |
| Jan 2026 | Google Gemini CLI support; skills.sh launches |
| Feb 2026 | Ecosystem surpasses 280K skills |

---

## TODO: Advanced Topics

Topics to explore later:

| Topic | Description | Status |
|-------|-------------|--------|
| Fine-tuning | Customizing model behavior with your own data | ✅ Done |
| LoRA / QLoRA | Efficient fine-tuning (low-rank adaptation) | ✅ Done |
| Function calling | Structured JSON output for tools | ✅ Done |
| JSON mode | Force valid JSON responses | ✅ Done |
| Seed/determinism | Reproducible outputs | ✅ Done |
| Prompt injection | Security vulnerability via prompts | ✅ Done |
| Guardrails | Filtering unsafe outputs | ✅ Done |
| Evaluation | How to measure LLM quality | |
| Benchmark datasets | MMLU, HumanEval, etc. | |
| Quantization | Smaller models (4-bit, 8-bit) | |
| Distillation | Smaller model trained on larger one | |
| Multi-modal | Vision, audio in LLMs | |
| Agents | Full agent architectures | |
| Memory systems | Short-term vs long-term |
