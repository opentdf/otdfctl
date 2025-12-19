# Agent Documentation

This directory contains documentation for AI coding agents (Claude Code, Cursor, Zed, etc.) working on the otdfctl codebase.

## Purpose

These documents use **progressive disclosure** to provide context to AI agents only when needed, keeping the main `CLAUDE.md` file concise and universally applicable.

## Structure

- **`agent-setup-instructions.md`** - Guide for setting up agent documentation (this structure)
- **`building_and_testing.md`** - Build process, testing, BATS tests, TestRail integration
- **`project_structure.md`** - Detailed codebase layout and component descriptions
- **`development_workflow.md`** - Adding commands, handlers, documentation, error handling
- **`tui_development.md`** - TUI framework (⚠️ work in progress, avoid unless instructed)

## How to Use

**For AI Agents:**

When starting a task, the agent should:
1. Read `/CLAUDE.md` for overview and essential context
2. Identify which agent docs are relevant to the current task
3. Read only the relevant documents from this directory
4. Proceed with the task using the focused context

**For Developers:**

When updating the project:
- Keep `CLAUDE.md` concise (< 60 lines ideally)
- Add detailed context to specific files in `agents/`
- Use file:line references instead of code snippets (avoid staleness)
- Update these docs when making significant architectural changes

## Philosophy

Following the principles from `agent-setup-instructions.md`:

1. **Less is More** - Include only necessary instructions
2. **Universal Applicability** - Keep `CLAUDE.md` relevant to all tasks
3. **Progressive Disclosure** - Let agents discover details when needed
4. **Pointers Over Copies** - Reference code locations, don't duplicate code
5. **High Leverage** - These docs affect every session, so craft them carefully

## What Goes Where

### CLAUDE.md (Root)
- WHAT: Tech stack, project structure overview
- WHY: Project purpose
- HOW: Essential commands, key patterns
- Pointers to these detailed docs

### agents/ (This Directory)
- Detailed build/test instructions
- Complete project structure breakdown
- Step-by-step development workflows
- Framework-specific guidance (TUI, testing, etc.)
- Task-specific context that's not universally needed

## Maintenance

When you modify the codebase:

- **New commands**: Update `development_workflow.md` if patterns change
- **New packages**: Update `project_structure.md`
- **Build changes**: Update `building_and_testing.md`
- **Architecture changes**: Update relevant docs and `CLAUDE.md` pointers

Keep the docs accurate. Outdated documentation is worse than no documentation.
