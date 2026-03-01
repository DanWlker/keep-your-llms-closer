# Plugins - Technical Guide

> Distribution packages that bundle commands, skills, hooks, agents, and MCP for sharing across projects/teams.

---

## Overview

Plugins are the **distribution layer** for Claude Code extensions. They bundle any combination of:
- Slash commands
- Skills
- Subagents
- Hooks
- MCP servers

**Key point:** Plugins are for sharing. If you want your team to use the same commands/skills/hooks, package them as a plugin.

---

## File Structure

```
my-plugin/
├── .claude-plugin/
│   └── plugin.json          # Required manifest
├── commands/                # Slash commands
│   └── my-command.md
├── agents/                 # Subagent definitions
│   └── my-agent.md
├── skills/                 # Skills
│   └── my-skill/
│       └── SKILL.md
├── hooks/
│   └── hooks.json          # Hook configuration
└── .mcp.json               # MCP server definitions
```

---

## Plugin Manifest (plugin.json)

```json
{
  "name": "my-plugin",
  "version": "1.0.0",
  "description": "My custom plugin for team workflows",
  "author": "Your Name",
  "homepage": "https://github.com/you/my-plugin"
}
```

---

## What Plugins Can Contain

| Component | Description | How it works |
|-----------|-------------|--------------|
| **Commands** | Slash commands | Becomes `/plugin:command` |
| **Agents** | Subagent definitions | Loaded automatically |
| **Skills** | Auto-discovered capabilities | Activates based on description |
| **Hooks** | Event handlers | Merges with existing hooks |
| **MCP** | Server configurations | Registers MCP servers |

---

## Installation

### From Marketplace

```
/plugin install owner/plugin-name
```

### From GitHub

```
/plugin install https://github.com/owner/repo
```

### Local Development

```
claude --plugin ./path/to/plugin
```

---

## How Plugins Work

### Installation Flow

```
User: /plugin install my-team/plugin

Framework:
  1. Downloads/parses plugin
  2. Validates plugin.json manifest
  3. Copies components to appropriate directories
  4. Registers commands, agents, hooks, MCP
  5. Skills auto-discovered on next prompt
```

### Component Merging

| Component | Merge Behavior |
|-----------|---------------|
| Commands | Added (namespaced as `/plugin:command`) |
| Agents | Added to agent pool |
| Skills | Auto-discovered |
| Hooks | Merged with existing |
| MCP | Registered alongside existing |

---

## Example: Team Standardization Plugin

```
team-standards/
├── .claude-plugin/
│   └── plugin.json
├── commands/
│   ├── deploy.md
│   └── ci-check.md
├── agents/
│   └── security-reviewer.md
├── skills/
│   └── team-conventions/
│       └── SKILL.md
└── hooks/
    └── hooks.json
```

### hooks.json

```json
{
  "PostToolUse": [
    {
      "matcher": "Write|Edit",
      "command": "team-lint ${event.path}"
    }
  ]
}
```

---

## Plugins vs Skills

| Aspect | Plugins | Skills |
|--------|---------|--------|
| **Purpose** | Distribution | Capability |
| **Scope** | Claude Code only | All Claude products |
| **Contains** | Commands, agents, hooks, MCP, skills | Just skill instructions |
| **Activation** | Varies by component | Auto (description match) |
| **Use case** | Team standardization | Teaching expertise |

---

## When to Use Plugins

| Scenario | Use |
|----------|-----|
| Share team configurations | Plugins |
| Package domain workflows | Plugins |
| Distribute skills | Either (plugins for teams) |
| Personal use | Direct in `.claude/` |
| Cross-product | Skills |

---

## Subagents in Plugins: Known Limitations

### MCP Access

**Issue:** Custom subagents defined in plugins cannot access MCP tools (bug #13605).

```yaml
# This subagent in a plugin CANNOT use MCP tools
---
name: my-agent
description: Does something with external API
tools:
  - Read
  - mcp__my-server__my-tool  # Won't work in plugins
---
```

**Workaround:** Use non-plugin subagents for MCP access.

### Best Practices for Plugin Subagents

| Do | Don't |
|----|-------|
| Use for commands, skills, hooks | Expect MCP to work in plugin agents |
| Read-only tool access | Give Write/Bash unless necessary |
| Define in `.claude/agents/` for MCP | Rely on plugin agents for external APIs |

---

## Key Takeaways

1. **Distribution layer** - Bundle extensions for sharing
2. **Manifest required** - `.claude-plugin/plugin.json`
3. **Components merged** - Hooks combine, commands namespaced
4. **Team standardization** - Main use case
5. **Install via** - `/plugin install` command

---

## References

- Claude Code plugins docs
- Plugin marketplace: https://plugins.claude.ai
- Structure: https://claude-plugins.dev
