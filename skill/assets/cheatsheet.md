# Go Debugging Cheat Sheet

One-page quick reference for delve-mcp debugging tools.

## Session Management

| Action | Command |
|--------|---------|
| Debug source file | `debug(file: "/abs/path/to/main.go", args: ["--flag"])` |
| Debug test | `debug_test(testfile: "/path/to/test.go", testname: "TestFoo")` |
| Attach to process | `attach(pid: 28026)` |
| Close session | `close()` |

## Breakpoints

| Action | Command |
|--------|---------|
| Set breakpoint | `set_breakpoint(file: "/path/to/file.go", line: 42)` |
| Conditional breakpoint | `set_breakpoint(file: "...", line: 42, condition: "x > 5")` |
| List all breakpoints | `list_breakpoints()` |
| Remove breakpoint | `remove_breakpoint(id: 1)` |

## Execution Control

| Action | Command |
|--------|---------|
| Resume execution | `continue()` |
| Step into function | `step()` |
| Step over function | `step_over()` |
| Step out of function | `step_out()` |

## Variable Inspection

| Action | Command |
|--------|---------|
| Inspect variable | `eval_variable(name: "varName")` |
| Deep inspection | `eval_variable(name: "varName", depth: 3)` |
| Evaluate expression | `eval_variable(name: "len(items)")` |
| Get program output | `get_debugger_output()` |

## Common Conditions

| Condition | Example |
|-----------|---------|
| Value comparison | `condition: "count > 100"` |
| String comparison | `condition: "name == \"Alice\""` |
| Nil check | `condition: "err != nil"` |
| Length check | `condition: "len(items) == 0"` |
| Boolean | `condition: "!initialized"` |
| Combined | `condition: "count > 5 && status == 200"` |

## Quick Workflows

### Debug Program
```
1. debug(file: "/path/to/main.go")
2. set_breakpoint(file: "...", line: 42)
3. continue()
4. eval_variable(name: "myVar")
5. close()
```

### Debug Test
```
1. debug_test(testfile: "...", testname: "TestFoo")
2. set_breakpoint in test and implementation
3. continue()
4. eval_variable(name: "expected")
5. eval_variable(name: "actual")
6. close()
```

### Attach to Production
```
1. attach(pid: 28026)
2. set_breakpoint with condition (production-safe!)
3. continue() immediately
4. Wait for condition to trigger
5. Inspect quickly (< 30 sec)
6. remove_breakpoint(id: 1)
7. close()
```

## Variable Depth Guide

| Depth | Use Case | Speed |
|-------|----------|-------|
| 1 | Quick overview | Fast |
| 2 | Moderate nesting | Balanced |
| 3-5 | Deep structures | Slower |
| 10+ | Complete inspection | Slow |

## Breakpoint Limits

- ✓ One breakpoint per line
- ✓ Multiple conditional breakpoints on different lines
- ✗ Multiple breakpoints on same line (use different lines)

## Production Safety

- ✓ Use conditional breakpoints
- ✓ Resume within 30 seconds
- ✓ Remove breakpoints when done
- ✓ Always call close()
- ✗ Never unconditional breakpoints on hot paths
- ✗ Never leave process paused

## Common Patterns

### Find Nil Pointer
```
set_breakpoint before crash line
eval_variable for each pointer
→ Find which one is nil
```

### Trace Function Calls
```
set_breakpoint at call
step() to enter
step_over() through function
step_out() to return
```

### Compare Expected vs Actual
```
set_breakpoint at assertion
eval_variable(name: "expected")
eval_variable(name: "actual")
→ Compare values
```

### Check Error Origin
```
set_breakpoint(condition: "err != nil")
step_out() repeatedly
→ Trace back to error source
```

### Inspect HTTP Request
```
set_breakpoint at handler entry
eval_variable(name: "r", depth: 2)
→ See method, headers, body
```

## Response Fields

Every tool returns:
```json
{
  "status": "success|error",
  "context": {
    "timestamp": "ISO 8601",
    "operation": "tool_name",
    "currentLocation": "file:line in function",
    "localVariables": [...],  // Automatic!
    "stopReason": "why stopped",
    "error": "error message"   // If error
  }
}
```

## Stop Reasons

- `"process is stopped"` - Paused after debug/attach
- `"hit breakpoint"` - Breakpoint triggered
- `"step complete"` - Step operation finished
- `"program exited"` - Program completed
- `"conditional breakpoint: <condition>"` - Conditional hit

## Troubleshooting

| Problem | Solution |
|---------|----------|
| "Breakpoint exists at this location" | Use different line or remove existing |
| "Process is not running" | Start session with debug/attach first |
| "Could not find symbol" | Variable out of scope or doesn't exist |
| Breakpoint never hits | Check path, line number, condition syntax |
| Program exits immediately | Set breakpoints before continue() |

## Best Practices

### DO
- ✓ Always close sessions
- ✓ Use absolute paths
- ✓ Check context.localVariables (automatic!)
- ✓ Set breakpoints before continue
- ✓ Use conditional breakpoints in production
- ✓ Verify breakpoint conditions

### DON'T
- ✗ Skip close() (leaves zombies)
- ✗ Use relative paths
- ✗ Forget to continue after debug/attach
- ✗ Set too many breakpoints
- ✗ Use deep depth unnecessarily
- ✗ Leave production processes paused

## File Paths

All file paths MUST be absolute:
- ✓ `/Users/vadim/app/server.go`
- ✓ `/home/user/project/main.go`
- ✗ `./server.go` (relative)
- ✗ `server.go` (relative)

## Condition Syntax

Use Go expression syntax:
- ✓ `count > 5`
- ✓ `name == "Alice"`
- ✓ `err != nil`
- ✓ `len(items) == 0`
- ✗ `count = 5` (assignment, not comparison)
- ✗ `name = "Alice"` (assignment)

## Integration

### With gopls
```
go_search(query: "HandleRequest")
→ Found at handler.go:45

set_breakpoint(file: "/app/handler.go", line: 45)
```

### With Testing
```bash
go test -v ./...
→ FAIL: TestHandleRequest

debug_test(testfile: "...", testname: "TestHandleRequest")
```

### With Process Tools
```bash
pgrep -fl myserver
→ 28026

attach(pid: 28026)
```

## See Also

- `SKILL.md` - Complete workflow guide
- `references/workflows.md` - Detailed procedures
- `references/patterns.md` - Common patterns
- `references/tools.md` - Full tool reference
- `references/examples.md` - Real-world examples
- `assets/quickstart.md` - Getting started
