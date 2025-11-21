---
name: go-debug
description: Comprehensive guide for debugging Go applications using the delve-mcp MCP server. Use when debugging Go programs, investigating bugs, analyzing test failures, stepping through code execution, inspecting variables at runtime, or attaching to running processes. Provides workflow-oriented guidance for effective debugging with conditional breakpoints, variable inspection, and execution control.
---

# Go Debugging with Delve MCP

Complete debugging solution for Go applications using the Delve debugger via the Model Context Protocol (MCP). This skill teaches you how to effectively debug Go programs, set breakpoints, inspect variables, and diagnose issues.

**MCP Server:** delve-mcp
**Debugger:** Delve (github.com/go-delve/delve)
**Protocol:** MCP 2024-11-05

## Core Philosophy

**Debugging is a scientific process**: Form hypothesis → Test → Observe → Refine

Key principles:
1. **Reproduce first** - No reproduction = No real debugging
2. **Binary search** - Cut problem space in half repeatedly
3. **Question assumptions** - The bug is usually in code you're certain works
4. **Trust the debugger** - It shows what IS happening, not what you THINK is happening

See `references/fundamentals.md` for complete debugging methodology and mental models.

## Quick Start

### Minimal Debugging Session

Debug a Go program in 3 commands:

```
1. Launch debugger:
   mcp__delve-mcp__debug(file: "/path/to/main.go")

2. Set breakpoint:
   mcp__delve-mcp__set_breakpoint(file: "/path/to/main.go", line: 42)

3. Continue execution:
   mcp__delve-mcp__continue()
```

When breakpoint hits, inspect variables:

```
mcp__delve-mcp__eval_variable(name: "myVar")
```

Close session when done:

```
mcp__delve-mcp__close()
```

### Common Use Cases

**Debug a crashing program:**
```
1. debug(file: "main.go")
2. set_breakpoint at suspected crash location
3. continue → inspect variables when stopped
4. step_over to trace execution flow
```

**Debug a failing test:**
```
1. debug_test(testfile: "handler_test.go", testname: "TestHandleRequest")
2. set_breakpoint in test function
3. set_breakpoint in code being tested
4. continue → compare expected vs actual values
```

**Attach to running server:**
```
1. Find process: ps aux | grep myserver
2. attach(pid: 12345)
3. set_breakpoint at handler function
4. continue → wait for request to trigger breakpoint
```

## Detection

Use this skill when:
- User mentions debugging, breakpoints, or stepping through code
- Investigating crashes, panics, or unexpected behavior
- Analyzing test failures or race conditions
- Inspecting runtime state or variable values
- User mentions "delve", "dlv", or debugging tools
- User asks to "pause", "stop at", or "inspect" during execution

Check if delve-mcp server is available before using debugging tools.

## Core Workflows

### 1. Basic Debugging Workflow

**Goal:** Launch a program, set breakpoints, and inspect execution

**Steps:**
1. **MUST** start a debug session first:
   - For source files: `debug(file: "/abs/path/to/file.go")`
   - For running processes: `attach(pid: <process_id>)`
   - For tests: `debug_test(testfile: "/path/to/test.go", testname: "TestFoo")`

2. **SHOULD** set breakpoints before continuing:
   - `set_breakpoint(file: "/path/to/file.go", line: 42)`
   - For conditional breaks: `set_breakpoint(..., condition: "count > 5")`

3. **MUST** continue execution to start the program:
   - `continue()` → program runs until breakpoint or completion

4. **MAY** inspect state when stopped:
   - `eval_variable(name: "varName")` → see current value
   - Check `context.localVariables` in response → automatic locals
   - Check `context.currentLocation` → where execution stopped

5. **MAY** control execution:
   - `step()` → step into function calls
   - `step_over()` → step over function calls
   - `step_out()` → step out of current function
   - `continue()` → resume until next breakpoint

6. **MUST** close session when done:
   - `close()` → cleanup resources and terminate debugger

**Example sequence:**
```
debug(file: "/app/server.go")
→ Session started

set_breakpoint(file: "/app/server.go", line: 45)
→ Breakpoint 1 set at server.go:45

continue()
→ Running... [stopped at breakpoint]
→ context.localVariables shows: request, response, err

eval_variable(name: "request", depth: 2)
→ {Method: "POST", URL: "/api/users", Body: {...}}

step_over()
→ Moved to line 46

close()
→ Session closed
```

### 2. Test Debugging Workflow

**Goal:** Debug failing tests to understand why they fail

**Steps:**
1. **MUST** identify the failing test name (exact match required)

2. **MUST** use `debug_test` with absolute path:
   ```
   debug_test(
     testfile: "/abs/path/to/handler_test.go",
     testname: "TestHandleRequest"
   )
   ```

3. **SHOULD** set strategic breakpoints:
   - Breakpoint in test function (verify test logic)
   - Breakpoint in code under test (verify implementation)
   - Breakpoint at assertion/comparison (compare expected vs actual)

4. **MUST** continue to start test execution:
   ```
   continue()
   ```

5. **SHOULD** inspect test-specific variables:
   - `expected` vs `actual` values
   - Test fixtures or mock data
   - Error values or return codes

6. **MAY** step through test execution:
   - `step()` into functions being tested
   - `step_over()` past setup code
   - `step_out()` from helper functions

**Example sequence:**
```
debug_test(testfile: "/app/handler_test.go", testname: "TestCreateUser")
→ Test debug session started

set_breakpoint(file: "/app/handler_test.go", line: 23)  # in test
→ Breakpoint 1 set

set_breakpoint(file: "/app/handler.go", line: 67)  # code under test
→ Breakpoint 2 set

continue()
→ Stopped at handler_test.go:23
→ localVariables: expected={Name: "Alice"}, actual=nil

step()
→ Entering CreateUser function

continue()
→ Stopped at handler.go:67
→ localVariables: name="Alice", valid=false  # Found the bug!

close()
```

### 3. Production Debugging Workflow

**Goal:** Debug a running server without restarting it

**Steps:**
1. **MUST** find the process ID:
   ```bash
   ps aux | grep myserver
   # or
   pgrep -f myserver
   ```

2. **MUST** attach to the process:
   ```
   attach(pid: 28026)
   ```

3. **SHOULD** set conditional breakpoints to filter noise:
   ```
   set_breakpoint(
     file: "/app/handler.go",
     line: 45,
     condition: "userID == \"problem_user\""
   )
   ```

4. **MUST** continue to resume the server:
   ```
   continue()
   ```

5. **WAIT** for breakpoint to trigger (may need external requests)

6. **MAY** inspect runtime state:
   - Check variables related to the issue
   - Examine request/response objects
   - Verify configuration values

7. **SHOULD** remove breakpoints when investigation is done:
   ```
   list_breakpoints()  # find breakpoint IDs
   remove_breakpoint(id: 1)
   ```

8. **MUST** close cleanly to avoid leaving process paused:
   ```
   close()
   ```

**Example sequence:**
```
# Terminal 1: Find running server
bash: ps aux | grep example
→ PID 28026

attach(pid: 28026)
→ Attached to process 28026

set_breakpoint(
  file: "/app/example.go",
  line: 24,
  condition: "name == \"Alice\""
)
→ Breakpoint 1 set (conditional)

continue()
→ Server resumed, waiting for request...

# Terminal 2: Trigger the condition
bash: curl http://localhost:8080/?name=Alice

# Back to debugger
→ Stopped at example.go:24
→ localVariables: name="Alice", requestCount=5

eval_variable(name: "message")
→ "Hello, Alice! This is request #5"

continue()
→ Resumed

close()
→ Detached from process
```

### 4. Advanced Debugging Workflow

**Goal:** Use conditional breakpoints and deep inspection for complex issues

**Conditional Breakpoint Strategies:**

1. **Filter by value:**
   ```
   condition: "count > 100"
   condition: "username == \"admin\""
   condition: "err != nil"
   ```

2. **Filter by state:**
   ```
   condition: "len(items) == 0"
   condition: "status == StatusFailed"
   condition: "!initialized"
   ```

3. **Combine conditions:**
   ```
   condition: "userID > 1000 && role == \"guest\""
   condition: "len(errors) > 0 || status == 500"
   ```

**Deep Variable Inspection:**

Use `depth` parameter for nested structures:

```
eval_variable(name: "request", depth: 1)  # shallow (default)
→ {Method: "...", URL: "...", Header: {...}}

eval_variable(name: "request", depth: 3)  # deep
→ {
    Method: "POST",
    URL: "/api/users",
    Header: {
      "Content-Type": ["application/json"],
      "Authorization": ["Bearer ..."]
    },
    Body: {...}
  }
```

**Multiple Breakpoints Strategy:**

Set breakpoints at different execution stages:

```
set_breakpoint(file: "handler.go", line: 23)  # Entry point
set_breakpoint(file: "handler.go", line: 67)  # Business logic
set_breakpoint(file: "handler.go", line: 89)  # Error handling
set_breakpoint(file: "handler.go", line: 102) # Return point

list_breakpoints()  # Verify all are set
```

## Decision Trees

### Which debugging tool to use?

```
START: What do you want to debug?

├─ A running process/server?
│  └─> Use: attach(pid: <pid>)
│
├─ A test function?
│  └─> Use: debug_test(testfile: "...", testname: "...")
│
└─ A Go program from source?
   └─> Use: debug(file: "main.go", args: [...])
```

### How to set breakpoints?

```
START: When should execution stop?

├─ Always at this line?
│  └─> Use: set_breakpoint(file: "...", line: 42)
│
├─ Only when a condition is true?
│  └─> Use: set_breakpoint(file: "...", line: 42, condition: "x > 5")
│
└─ Multiple locations?
   └─> Call set_breakpoint multiple times
```

### How to navigate execution?

```
START: Where do you want to go next?

├─ Enter the next function call?
│  └─> Use: step()
│
├─ Skip over function calls (stay at current level)?
│  └─> Use: step_over()
│
├─ Exit current function (return to caller)?
│  └─> Use: step_out()
│
└─ Run until next breakpoint or end?
   └─> Use: continue()
```

## Tool Reference

### Session Management

**debug** - Debug a Go source file
```
debug(file: "/abs/path/to/main.go", args: ["--port", "8080"])
→ Starts debug session, program paused at entry point
```

**debug_test** - Debug a specific test function
```
debug_test(
  testfile: "/abs/path/to/handler_test.go",
  testname: "TestHandleRequest",
  testflags: ["-v", "-count=1"]  # optional
)
→ Starts test in debug mode, paused at test start
```

**attach** - Attach to running process
```
attach(pid: 28026)
→ Attaches to process, pauses execution
```

**close** - End debugging session
```
close()
→ Detaches, cleans up resources, terminates debugger
→ IMPORTANT: Always call this when done!
```

### Breakpoint Management

**set_breakpoint** - Set breakpoint (regular or conditional)
```
set_breakpoint(
  file: "/abs/path/to/handler.go",
  line: 45,
  condition: "userID > 1000"  # optional
)
→ Returns: {id: 1, location: "...", condition: "..."}
```

**list_breakpoints** - Show all breakpoints
```
list_breakpoints()
→ Returns: [{id: 1, location: "...", hitCount: 3}, ...]
```

**remove_breakpoint** - Remove a breakpoint
```
remove_breakpoint(id: 1)
→ Breakpoint removed
```

### Execution Control

**continue** - Resume execution
```
continue()
→ Runs until breakpoint, panic, or program end
```

**step** - Step into next line
```
step()
→ Enters function calls, moves to next source line
```

**step_over** - Step over next line
```
step_over()
→ Skips over function calls, stays at current level
```

**step_out** - Step out of function
```
step_out()
→ Runs until current function returns
```

### Variable Inspection

**eval_variable** - Evaluate a variable
```
eval_variable(name: "request", depth: 2)
→ Returns: {value: {...}, type: "*http.Request"}

# Depth controls nested structure traversal:
# depth: 1 (default) - shallow inspection
# depth: 2-5 - moderate nesting
# depth: 10+ - deep inspection (slow)
```

**get_debugger_output** - Get program stdout/stderr
```
get_debugger_output()
→ Returns captured console output
→ Useful for seeing print statements during debug
```

## Best Practices

### DO's

✓ **Always close sessions** - Call `close()` when done to cleanup resources
✓ **Use absolute paths** - All file paths must be absolute
✓ **Check context** - Inspect `context.localVariables` in responses (automatic)
✓ **Set breakpoints before continue** - Otherwise program may run to completion
✓ **Use conditional breakpoints** - Filter noise in production/high-traffic code
✓ **Verify breakpoint conditions** - Test condition syntax before relying on it
✓ **Check error field** - All responses include `context.error` if something failed

### DON'Ts

✗ **Don't skip close()** - Leaves zombie processes and temp files
✗ **Don't use relative paths** - Will fail to find files
✗ **Don't forget to continue** - Program stays paused after debug/attach
✗ **Don't set too many breakpoints** - Slows execution, use conditionals instead
✗ **Don't use deep depth blindly** - High depth values are slow
✗ **Don't assume breakpoint hit** - Check `context.stopReason` in response

### Breakpoint Gotchas

- **Can't set multiple breakpoints at same line** - Delve limitation
  - Workaround: Set breakpoint on adjacent lines with different conditions

- **Condition syntax is Go expressions** - Use Go syntax, not shell/other languages
  - Correct: `name == "Alice"`
  - Wrong: `name = "Alice"` (assignment, not comparison)

- **String comparisons need quotes** - `condition: "name == \"Alice\""`

- **Can reference local variables only** - Can't use arbitrary expressions

## Troubleshooting

### "Breakpoint exists at this location"

**Problem:** Delve doesn't allow multiple breakpoints at the same line.

**Solution:** Set breakpoint on a different line, or remove existing breakpoint first:
```
list_breakpoints()
remove_breakpoint(id: 1)
set_breakpoint(...)  # Try again
```

### "Process is not running"

**Problem:** Tried to step/continue but no active debug session.

**Solution:** Start a session first with `debug`, `debug_test`, or `attach`.

### "Evaluation error: could not find symbol"

**Problem:** Variable name doesn't exist in current scope.

**Solution:**
- Check `context.localVariables` for available variables
- Ensure you're stopped at a location where the variable exists
- Variable may be out of scope or not yet initialized

### Breakpoint never hits

**Problem:** Set breakpoint but execution never stops there.

**Solution:**
- Verify file path is absolute and correct
- Check line number is executable (not a comment or blank line)
- If using condition, verify condition is correct syntax
- Ensure code path actually executes (not dead code)

### Program exits immediately

**Problem:** Called `continue()` but program exits without hitting breakpoints.

**Solution:**
- Verify breakpoints are set before calling `continue()`
- Check breakpoint locations are reachable (not dead code)
- For servers, may need external trigger (HTTP request, etc.)
- For tests, ensure test name matches exactly

## Integration with Other Tools

### With gopls (Go Language Server)

Use gopls MCP tools to find debugging locations:

```
# Find where a function is defined
go_search(query: "HandleRequest")
→ Found at handler.go:45

# Set breakpoint there
set_breakpoint(file: "/app/handler.go", line: 45)
```

### With Testing Tools

Debug specific test failures identified by `go test`:

```bash
# Run tests first to identify failure
go test -v ./...
→ FAIL: TestHandleRequest

# Then debug the specific test
debug_test(testfile: "...", testname: "TestHandleRequest")
```

### With Process Monitoring

Use system tools to find processes for attach:

```bash
# Find process by name
pgrep -fl myserver

# Or search processes
ps aux | grep myserver

# Then attach
attach(pid: <pid>)
```

## Quick Reference Card

| Task | Tool | Example |
|------|------|---------|
| Debug program | `debug` | `debug(file: "/app/main.go")` |
| Debug test | `debug_test` | `debug_test(testfile: "...", testname: "TestFoo")` |
| Attach to process | `attach` | `attach(pid: 12345)` |
| Set breakpoint | `set_breakpoint` | `set_breakpoint(file: "...", line: 42)` |
| Conditional break | `set_breakpoint` | `set_breakpoint(..., condition: "x > 5")` |
| List breakpoints | `list_breakpoints` | `list_breakpoints()` |
| Remove breakpoint | `remove_breakpoint` | `remove_breakpoint(id: 1)` |
| Start/resume | `continue` | `continue()` |
| Step into | `step` | `step()` |
| Step over | `step_over` | `step_over()` |
| Step out | `step_out` | `step_out()` |
| Inspect variable | `eval_variable` | `eval_variable(name: "user", depth: 2)` |
| Get output | `get_debugger_output` | `get_debugger_output()` |
| End session | `close` | `close()` |

## See Also

- `references/fundamentals.md` - **Core debugging methodology and mental models**
- `references/workflows.md` - Detailed workflow procedures
- `references/patterns.md` - Common debugging patterns
- `references/tools.md` - Complete tool documentation
- `references/examples.md` - Real-world debugging scenarios
- `assets/quickstart.md` - Minimal getting started guide
- `assets/cheatsheet.md` - One-page reference
