# Installation Instructions

## Quick Install

Move this `skill` directory to your Claude Code skills location:

```bash
# On macOS/Linux
mv skill ~/.claude/skills/go-debug

# On Windows
move skill %USERPROFILE%\.claude\skills\go-debug
```

## Verify Installation

Once installed, Claude Code will automatically detect and use the skill when you:
- Ask about debugging Go programs
- Mention breakpoints or stepping through code
- Need to inspect variables or investigate bugs

## Test the Skill

Try asking Claude Code:
```
Debug testdata/webserver/webserver.go and set a conditional breakpoint when name equals "Alice"
```

Claude should automatically use the go-debug skill to guide the debugging session.

## What You Get

After installation, Claude Code will know how to:
- Debug Go programs from source files
- Debug specific test functions
- Attach to running production processes
- Set conditional breakpoints to filter noise
- Inspect variables at different depth levels
- Navigate through code with step/step_over/step_out
- Investigate crashes, bugs, and performance issues
- Use safe production debugging practices

## File Structure

```
~/.claude/skills/go-debug/
├── SKILL.md                # Main skill (Claude reads this)
├── README.md              # This documentation
├── INSTALL.md            # Installation instructions
├── references/
│   ├── workflows.md      # Detailed debugging procedures
│   ├── patterns.md       # Common debugging patterns
│   ├── tools.md         # Complete tool reference
│   └── examples.md      # Real-world scenarios
└── assets/
    ├── quickstart.md    # Quick start guide
    └── cheatsheet.md   # One-page reference
```

## Size

- **Total:** ~4,500 lines of documentation
- **Main skill:** ~650 lines (core guidance)
- **References:** ~3,200 lines (detailed workflows, patterns, tools, examples)
- **Assets:** ~480 lines (quick reference materials)

## Prerequisites

Ensure you have:
1. Claude Code installed and configured
2. delve-mcp server configured in `.mcp.json`
3. Go and Delve installed

## Updating

To update the skill:
```bash
# Remove old version
rm -rf ~/.claude/skills/go-debug

# Install new version
mv skill ~/.claude/skills/go-debug
```

## Troubleshooting

**Skill not being used?**
- Restart Claude Code after installation
- Check that the directory is named `go-debug` (not `skill`)
- Verify `SKILL.md` exists and has the frontmatter metadata

**Can't find skills directory?**
- Create it: `mkdir -p ~/.claude/skills`
- Check Claude Code documentation for your platform's location

## See Also

- `README.md` - Overview and structure
- `SKILL.md` - The actual skill Claude reads
- `assets/cheatsheet.md` - Quick reference for manual use
