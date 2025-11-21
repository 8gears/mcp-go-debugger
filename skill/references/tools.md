# Complete Tool Reference

Comprehensive documentation for all delve-mcp debugging tools.

## Table of Contents

1. [Session Management](#session-management)
2. [Breakpoint Management](#breakpoint-management)
3. [Execution Control](#execution-control)
4. [Variable Inspection](#variable-inspection)
5. [Output Capture](#output-capture)
6. [Tool Response Format](#tool-response-format)

---

## Session Management

### debug

**Purpose:** Debug a Go source file from the beginning.

**Signature:**
```
mcp__delve-mcp__debug(
  file: string,      # Absolute path to Go source file (required)
  args: []string     # Command-line arguments for the program (optional)
)
```

**Parameters:**
- `file` (required): Absolute path to the Go source file to debug
- `args` (optional): Array of command-line arguments to pass to the program

**Behavior:**
- Compiles the Go source file with debug symbols
- Starts the program in debug mode
- Program is paused at entry point (beginning of main)
- Creates a debug session (only one session allowed at a time)

**Response:**
```json
{
  "status": "success",
  "context": {
    "timestamp": "2025-11-21T15:00:00Z",
    "operation": "debug",
    "currentLocation": "At /path/to/main.go:1 in main.main",
    "localVariables": [],
    "stopReason": "process is stopped"
  },
  "target": "/tmp/debug_binary123456"
}
```

**Example:**
```
mcp__delve-mcp__debug(
  file: "/Users/vadim/app/server.go",
  args: ["--port", "8080", "--verbose"]
)
```

**Use When:**
- Debugging a new program from scratch
- Need to debug main() function
- Want to trace execution from the start

**Notes:**
- File path MUST be absolute
- Only one debug session per MCP connection
- Call `close()` when done to cleanup

---

### debug_test

**Purpose:** Debug a specific Go test function.

**Signature:**
```
mcp__delve-mcp__debug_test(
  testfile: string,     # Absolute path to test file (required)
  testname: string,     # Exact test function name (required)
  testflags: []string   # Additional go test flags (optional)
)
```

**Parameters:**
- `testfile` (required): Absolute path to the Go test file
- `testname` (required): Exact name of the test function (e.g., "TestHandleRequest")
- `testflags` (optional): Array of flags to pass to `go test` (e.g., ["-v", "-count=1"])

**Behavior:**
- Compiles test with debug symbols
- Runs only the specified test function
- Program paused at test start
- Automatically adds `-test.run=^TestName$` flag for exact match

**Response:**
```json
{
  "status": "success",
  "context": {
    "timestamp": "2025-11-21T15:00:00Z",
    "operation": "debug_test",
    "currentLocation": "At /path/to/handler_test.go:23 in TestHandleRequest",
    "stopReason": "process is stopped"
  },
  "target": "/tmp/debug_binary123456"
}
```

**Example:**
```
mcp__delve-mcp__debug_test(
  testfile: "/Users/vadim/app/handler_test.go",
  testname: "TestCreateUser",
  testflags: ["-v", "-count=1"]
)
```

**Use When:**
- Debugging failing tests
- Understanding test behavior
- Verifying test setup/teardown

**Notes:**
- Test name must match exactly (case-sensitive)
- Test name includes "Test" prefix (e.g., "TestFoo", not "Foo")
- Can debug table-driven tests by test function name

---

### attach

**Purpose:** Attach debugger to a running Go process.

**Signature:**
```
mcp__delve-mcp__attach(
  pid: number     # Process ID to attach to (required)
)
```

**Parameters:**
- `pid` (required): Process ID of the running Go program

**Behavior:**
- Attaches Delve to the running process
- Process is **immediately paused** (important!)
- All goroutines frozen until `continue()` is called
- Existing process state preserved

**Response:**
```json
{
  "status": "success",
  "context": {
    "timestamp": "2025-11-21T15:00:00Z",
    "operation": "attach",
    "currentLocation": "At /path/to/server.go:67 in runtime.select",
    "stopReason": "process is stopped"
  },
  "pid": 28026
}
```

**Example:**
```
# Find process first
bash: ps aux | grep myserver
→ user 28026 ... myserver

# Then attach
mcp__delve-mcp__attach(pid: 28026)
```

**Use When:**
- Debugging production servers
- Investigating live issues
- Attaching to long-running processes

**Notes:**
- Requires appropriate permissions (same user or root)
- Process is paused immediately - call `continue()` ASAP!
- Can't attach to processes without debug symbols
- Use conditional breakpoints to avoid pausing production traffic

---

### close

**Purpose:** End the debug session and cleanup resources.

**Signature:**
```
mcp__delve-mcp__close()
```

**Parameters:** None

**Behavior:**
- Detaches from debugged process (if attached)
- Stops debug server
- Removes temporary debug binaries
- Cleans up all resources
- Resets session state

**Response:**
```json
{
  "status": "success",
  "context": {
    "timestamp": "2025-11-21T15:00:00Z",
    "operation": "close"
  }
}
```

**Example:**
```
mcp__delve-mcp__close()
```

**Use When:**
- Finished debugging
- Before starting a new debug session
- On error/cleanup

**Notes:**
- **ALWAYS** call this when done debugging
- Failure to call can leave zombie processes
- Required before starting new session
- Safe to call multiple times

---

## Breakpoint Management

### set_breakpoint

**Purpose:** Set a breakpoint at a specific line, optionally with a condition.

**Signature:**
```
mcp__delve-mcp__set_breakpoint(
  file: string,         # Absolute path to source file (required)
  line: number,         # Line number (required)
  condition: string     # Go expression condition (optional)
)
```

**Parameters:**
- `file` (required): Absolute path to the source file
- `line` (required): Line number to break at
- `condition` (optional): Go expression that must be true to trigger breakpoint

**Behavior:**
- Sets breakpoint at specified location
- If condition provided, only breaks when condition is true
- Returns breakpoint ID for later reference
- Can only set one breakpoint per line

**Response:**
```json
{
  "status": "success",
  "context": {
    "timestamp": "2025-11-21T15:00:00Z",
    "operation": "set_breakpoint",
    "stopReason": "process is stopped"
  },
  "breakpoint": {
    "id": 1,
    "status": "enabled",
    "location": "At /path/to/handler.go:45 in handleRequest",
    "condition": "userID > 1000",
    "hitCount": 0
  }
}
```

**Examples:**

Simple breakpoint:
```
mcp__delve-mcp__set_breakpoint(
  file: "/app/handler.go",
  line: 45
)
```

Conditional breakpoint:
```
mcp__delve-mcp__set_breakpoint(
  file: "/app/handler.go",
  line: 45,
  condition: "username == \"admin\""
)
```

Multiple conditions:
```
mcp__delve-mcp__set_breakpoint(
  file: "/app/handler.go",
  line: 45,
  condition: "count > 100 && status == StatusFailed"
)
```

Error condition:
```
mcp__delve-mcp__set_breakpoint(
  file: "/app/handler.go",
  line: 67,
  condition: "err != nil"
)
```

**Use When:**
- Stopping at specific code location
- Filtering breakpoints by condition
- Debugging specific scenarios only

**Notes:**
- File path must be absolute
- Line must be executable (not comment/blank line)
- Condition uses Go expression syntax
- String comparisons need escaped quotes: `"name == \"Alice\""`
- Can reference local variables only
- One breakpoint per line (Delve limitation)

---

### list_breakpoints

**Purpose:** List all currently set breakpoints.

**Signature:**
```
mcp__delve-mcp__list_breakpoints()
```

**Parameters:** None

**Behavior:**
- Returns all breakpoints currently set
- Includes system breakpoints (panics, fatal errors)
- Shows hit count for each breakpoint

**Response:**
```json
{
  "status": "success",
  "context": {
    "timestamp": "2025-11-21T15:00:00Z",
    "operation": "list_breakpoints",
    "stopReason": "process is stopped"
  },
  "breakpoints": [
    {
      "id": 1,
      "status": "enabled",
      "location": "At /app/handler.go:45 in handleRequest",
      "condition": "userID > 1000",
      "hitCount": 3
    },
    {
      "id": 2,
      "status": "enabled",
      "location": "At /app/service.go:67 in processUser",
      "condition": "",
      "hitCount": 0
    }
  ]
}
```

**Example:**
```
mcp__delve-mcp__list_breakpoints()
```

**Use When:**
- Verifying breakpoints are set
- Finding breakpoint IDs for removal
- Checking breakpoint hit counts
- Debugging breakpoint issues

**Notes:**
- Negative IDs (-1, -2) are system breakpoints
- hitCount shows how many times breakpoint was hit
- System breakpoints catch panics/fatal errors

---

### remove_breakpoint

**Purpose:** Remove a breakpoint by its ID.

**Signature:**
```
mcp__delve-mcp__remove_breakpoint(
  id: number     # Breakpoint ID to remove (required)
)
```

**Parameters:**
- `id` (required): The breakpoint ID (from `list_breakpoints()`)

**Behavior:**
- Removes the specified breakpoint
- Breakpoint will no longer trigger
- ID is no longer valid after removal

**Response:**
```json
{
  "status": "success",
  "context": {
    "timestamp": "2025-11-21T15:00:00Z",
    "operation": "remove_breakpoint"
  }
}
```

**Example:**
```
# First list breakpoints
mcp__delve-mcp__list_breakpoints()
→ Breakpoint 1 at handler.go:45

# Remove it
mcp__delve-mcp__remove_breakpoint(id: 1)
```

**Use When:**
- Cleaning up temporary breakpoints
- Removing unwanted breakpoints
- Before closing debug session

**Notes:**
- Cannot remove system breakpoints (negative IDs)
- Removing non-existent ID returns error

---

## Execution Control

### continue

**Purpose:** Resume program execution until next breakpoint or completion.

**Signature:**
```
mcp__delve-mcp__continue()
```

**Parameters:** None

**Behavior:**
- Resumes program execution
- Runs until:
  - Breakpoint is hit
  - Program panics
  - Program completes
  - Manual interrupt
- Returns immediately if already running

**Response when breakpoint hit:**
```json
{
  "status": "success",
  "context": {
    "timestamp": "2025-11-21T15:00:00Z",
    "operation": "continue",
    "currentLocation": "At /app/handler.go:45 in handleRequest",
    "localVariables": [
      {"name": "userID", "value": "123", "type": "string"},
      {"name": "err", "value": "nil", "type": "error"}
    ],
    "stopReason": "hit breakpoint"
  }
}
```

**Response when program completes:**
```json
{
  "status": "success",
  "context": {
    "timestamp": "2025-11-21T15:00:00Z",
    "operation": "continue",
    "stopReason": "program exited"
  }
}
```

**Example:**
```
mcp__delve-mcp__continue()
```

**Use When:**
- Starting program execution (after debug/attach)
- Resuming after breakpoint
- Running to next breakpoint
- Continuing production process (after attach)

**Notes:**
- MUST call after `debug()`, `debug_test()`, or `attach()`
- For servers, may run indefinitely until breakpoint
- For tests, runs until test completes

---

### step

**Purpose:** Execute next line of code, stepping into function calls.

**Signature:**
```
mcp__delve-mcp__step()
```

**Parameters:** None

**Behavior:**
- Executes one source line
- Enters function calls (shows function internals)
- Stops at next executable line

**Response:**
```json
{
  "status": "success",
  "context": {
    "timestamp": "2025-11-21T15:00:00Z",
    "operation": "step",
    "currentLocation": "At /app/handler.go:46 in validateUser",
    "localVariables": [...],
    "stopReason": "step complete"
  }
}
```

**Example:**
```
# At line 45: result := validateUser(user)
mcp__delve-mcp__step()
→ Now inside validateUser() at line 1
```

**Use When:**
- Want to see function implementation
- Debugging function internals
- Understanding detailed execution flow

**Notes:**
- Slow for deep call stacks
- Use `step_over()` to skip function details
- Use `step_out()` to exit current function

---

### step_over

**Purpose:** Execute next line of code, stepping over function calls.

**Signature:**
```
mcp__delve-mcp__step_over()
```

**Parameters:** None

**Behavior:**
- Executes one source line
- Skips over function calls (doesn't enter them)
- Stops at next line in current function

**Response:**
```json
{
  "status": "success",
  "context": {
    "timestamp": "2025-11-21T15:00:00Z",
    "operation": "step_over",
    "currentLocation": "At /app/handler.go:46 in handleRequest",
    "localVariables": [...],
    "stopReason": "step complete"
  }
}
```

**Example:**
```
# At line 45: result := validateUser(user)
mcp__delve-mcp__step_over()
→ Now at line 46 (validateUser executed, stayed in handleRequest)
```

**Use When:**
- Debugging at current function level
- Skipping function details
- Fast execution through code

**Notes:**
- Faster than `step()`
- Use when you trust the called function
- Still stops at breakpoints inside called functions

---

### step_out

**Purpose:** Execute until current function returns.

**Signature:**
```
mcp__delve-mcp__step_out()
```

**Parameters:** None

**Behavior:**
- Runs until current function returns
- Stops at the line after the function call
- Returns to calling function

**Response:**
```json
{
  "status": "success",
  "context": {
    "timestamp": "2025-11-21T15:00:00Z",
    "operation": "step_out",
    "currentLocation": "At /app/handler.go:46 in handleRequest",
    "localVariables": [...],
    "stopReason": "step complete"
  }
}
```

**Example:**
```
# Inside validateUser() at line 23
mcp__delve-mcp__step_out()
→ Back in handleRequest() at line 46 (after validateUser call)
```

**Use When:**
- Deep in call stack, want to go back up
- Accidentally stepped into unwanted function
- Finished inspecting current function

**Notes:**
- Useful after accidentally using `step()` instead of `step_over()`
- May take time if function has loops

---

## Variable Inspection

### eval_variable

**Purpose:** Evaluate and inspect a variable's value.

**Signature:**
```
mcp__delve-mcp__eval_variable(
  name: string,      # Variable name or expression (required)
  depth: number      # Recursion depth for nested structures (optional, default: 1)
)
```

**Parameters:**
- `name` (required): Variable name or Go expression
- `depth` (optional): How deep to traverse nested structures (default: 1)

**Behavior:**
- Evaluates the expression in current scope
- Returns value, type, and kind
- Recursively expands nested structures up to `depth`

**Response:**
```json
{
  "status": "success",
  "context": {
    "timestamp": "2025-11-21T15:00:00Z",
    "operation": "eval_variable",
    "stopReason": "process is stopped"
  },
  "variable": {
    "name": "user",
    "value": "{ID: 123, Name: \"Alice\", Email: \"alice@example.com\"}",
    "type": "*User",
    "kind": "pointer"
  }
}
```

**Examples:**

Simple variable:
```
mcp__delve-mcp__eval_variable(name: "count")
→ {value: "42", type: "int", kind: "int"}
```

Struct (shallow):
```
mcp__delve-mcp__eval_variable(name: "user", depth: 1)
→ {value: "{ID: 123, Name: \"Alice\"}", type: "*User", kind: "pointer"}
```

Struct (deep):
```
mcp__delve-mcp__eval_variable(name: "user", depth: 3)
→ {
    value: "{
      ID: 123,
      Name: \"Alice\",
      Email: \"alice@example.com\",
      Profile: {
        Avatar: \"url...\",
        Bio: \"...\"
      }
    }",
    type: "*User",
    kind: "pointer"
  }
```

Expression:
```
mcp__delve-mcp__eval_variable(name: "len(items)")
→ {value: "10", type: "int", kind: "int"}
```

Map access:
```
mcp__delve-mcp__eval_variable(name: "cache[\"user:123\"]", depth: 2)
→ {value: "{Name: \"Alice\", ...}", type: "interface{}", kind: "interface"}
```

Slice element:
```
mcp__delve-mcp__eval_variable(name: "users[0]", depth: 2)
→ {value: "{ID: 1, Name: \"Alice\"}", type: "User", kind: "struct"}
```

**Use When:**
- Inspecting variable values
- Checking struct fields
- Evaluating expressions
- Comparing expected vs actual

**Notes:**
- Only local variables in scope can be accessed
- Depth 1 = shallow (fast), 5+ = deep (slow)
- Can call simple methods: `user.IsAdmin()`
- Cannot modify variables (read-only)

---

## Output Capture

### get_debugger_output

**Purpose:** Retrieve stdout/stderr output from the debugged program.

**Signature:**
```
mcp__delve-mcp__get_debugger_output()
```

**Parameters:** None

**Behavior:**
- Returns all captured stdout and stderr
- Output is buffered during execution
- Includes print statements, log output, etc.

**Response:**
```json
{
  "status": "success",
  "context": {
    "timestamp": "2025-11-21T15:00:00Z",
    "operation": "get_debugger_output"
  },
  "output": "[15:00:01] Server starting...\n[15:00:02] Listening on :8080\n[15:00:05] Request received from 192.168.1.1\n"
}
```

**Example:**
```
mcp__delve-mcp__get_debugger_output()
```

**Use When:**
- Checking program log output
- Verifying print statements
- Correlating debug state with logs
- Debugging output-based issues

**Notes:**
- Output is cumulative (all output since start)
- May be large for verbose programs
- Captured regardless of breakpoints

---

## Tool Response Format

All tools return a consistent response structure:

### Success Response
```json
{
  "status": "success",
  "context": {
    "timestamp": "2025-11-21T15:00:00Z",
    "operation": "tool_name",
    "currentLocation": "At /path/to/file.go:45 in functionName",
    "localVariables": [
      {
        "name": "varName",
        "value": "value",
        "type": "type",
        "scope": "local|argument",
        "kind": "string|int|struct|..."
      }
    ],
    "stopReason": "process is stopped|hit breakpoint|step complete|..."
  },
  // Tool-specific fields...
}
```

### Error Response
```json
{
  "status": "error",
  "context": {
    "timestamp": "2025-11-21T15:00:00Z",
    "operation": "tool_name",
    "error": "Error message here"
  }
}
```

### Context Fields

- `timestamp`: ISO 8601 timestamp of operation
- `operation`: Name of the operation executed
- `currentLocation`: Where execution is currently stopped (if applicable)
- `localVariables`: Array of local variables at current location (automatic)
- `stopReason`: Why execution stopped
- `error`: Error message (only in error responses)

### Stop Reasons

- `"process is stopped"` - Process paused (after debug/attach)
- `"hit breakpoint"` - Breakpoint was triggered
- `"step complete"` - Step operation finished
- `"program exited"` - Program completed execution
- `"conditional breakpoint: <condition>"` - Conditional breakpoint hit

---

## Quick Reference

| Tool | Purpose | Key Parameters |
|------|---------|----------------|
| `debug` | Debug source file | `file`, `args` |
| `debug_test` | Debug test function | `testfile`, `testname`, `testflags` |
| `attach` | Attach to process | `pid` |
| `close` | End session | - |
| `set_breakpoint` | Set breakpoint | `file`, `line`, `condition` |
| `list_breakpoints` | List breakpoints | - |
| `remove_breakpoint` | Remove breakpoint | `id` |
| `continue` | Resume execution | - |
| `step` | Step into | - |
| `step_over` | Step over | - |
| `step_out` | Step out | - |
| `eval_variable` | Inspect variable | `name`, `depth` |
| `get_debugger_output` | Get output | - |
