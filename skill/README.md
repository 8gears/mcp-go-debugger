# Go Debugging Skill for Claude Code

This directory contains a comprehensive Claude Code skill for debugging Go applications using the delve-mcp MCP server.

## Installation

Move this entire `skill` directory to your Claude Code skills location:

```bash
# On macOS/Linux
mv skill ~/.claude/skills/go-debug

# On Windows
move skill %USERPROFILE%\.claude\skills\go-debug
```

## Structure

```
go-debug/
├── SKILL.md                    # Main skill file (start here!)
├── references/
│   ├── fundamentals.md        # Core debugging methodology & mental models
│   ├── workflows.md           # Detailed debugging workflows
│   ├── patterns.md            # Common debugging patterns
│   ├── tools.md              # Complete tool reference
│   └── examples.md           # Real-world debugging scenarios
└── assets/
    ├── quickstart.md         # Quick start guide
    └── cheatsheet.md        # One-page reference
```

## Usage

Once installed, Claude Code will automatically use this skill when:
- You mention debugging Go programs
- You ask about setting breakpoints
- You need to inspect variables
- You're investigating bugs or test failures
- You mention "delve" or debugging tools

## What's Included

### Main Skill (SKILL.md)
- Quick start examples
- Detection rules (when to use)
- Core workflows (basic, test, production, advanced)
- Decision trees (which tool to use)
- Tool reference (quick lookup)
- Best practices and troubleshooting

### Detailed References

**fundamentals.md** - Core debugging methodology:
- The debugging mindset (reproduce, binary search, question assumptions)
- Systematic approach (observe, trace backwards, root cause)
- Data inspection techniques (variable evolution, pointers, state validation)
- Advanced patterns (differential debugging, hypothesis testing)
- Bug-specific strategies (nil pointers, races, deadlocks)
- Professional debugging workflow
- Mental models and key principles

**workflows.md** - Complete step-by-step procedures:
- Interactive development debugging
- CI/CD test failure investigation
- Production issue investigation
- Race condition debugging
- Memory leak investigation
- Goroutine debugging
- Integration test debugging

**patterns.md** - Common debugging patterns:
- Finding bugs (crashes, nil pointers, wrong values, logic errors)
- Understanding code flow
- Inspecting data structures
- Error investigation
- Performance issues
- Concurrency issues
- Testing patterns
- HTTP/API debugging

**tools.md** - Complete tool documentation:
- Session management (debug, debug_test, attach, close)
- Breakpoint management (set, list, remove)
- Execution control (continue, step, step_over, step_out)
- Variable inspection (eval_variable, get_debugger_output)
- Tool response format
- Parameter details and examples

**examples.md** - Real-world scenarios:
- Debugging the example web server
- Test debugging with calculator tests
- Production debugging (attach to running process)
- Common scenarios with complete sessions

### Quick Assets

**quickstart.md** - Minimal getting started guide from the examples directory

**cheatsheet.md** - One-page reference card with:
- Command syntax
- Common conditions
- Quick workflows
- Variable depth guide
- Production safety rules
- Troubleshooting

## Prerequisites

This skill requires:
- The delve-mcp MCP server to be configured in your `.mcp.json`
- Go programming environment
- Delve debugger installed

## Features

### Workflow-Oriented Guidance
- Not just tool descriptions - complete debugging procedures
- Step-by-step workflows for common scenarios
- Decision trees for choosing the right approach

### LLM-Friendly Structure
- Clear MUST/SHOULD/MAY rules
- Detection criteria for automatic activation
- Context-aware tool usage
- Error recovery strategies

### Production-Ready Patterns
- Safe production debugging with conditional breakpoints
- Performance debugging strategies
- Concurrency issue investigation
- Real-world examples from the mcp-go-debugger project

### Comprehensive Coverage
- Basic debugging (launch, breakpoint, inspect)
- Test debugging (unit tests, table-driven tests)
- Production debugging (attach to live processes)
- Advanced techniques (conditional breakpoints, deep inspection)

## Skill Metadata

```yaml
name: go-debug
description: Comprehensive guide for debugging Go applications using the delve-mcp MCP server. Use when debugging Go programs, investigating bugs, analyzing test failures, stepping through code execution, inspecting variables at runtime, or attaching to running processes.
```

### Complementary Approach

The skill complements the MCP server:
- **MCP Server** - Low-level debugging operations (set breakpoint, step, eval)
- **Skill** - High-level debugging strategies (find bugs, trace flow, investigate errors)

## Contributing

To update this skill:
1. Edit the relevant markdown files
2. Keep examples synchronized with the `testdata/` directory
3. Test workflows with real debugging scenarios
4. Update the cheatsheet when adding new patterns

## See Also

- [MCP Go Debugger](https://github.com/sunfmin/mcp-go-debugger)
- [Delve Debugger](https://github.com/go-delve/delve)
- [Claude Code Skills](https://docs.claude.ai/code/skills)
