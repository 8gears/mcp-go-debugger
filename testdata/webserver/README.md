# MCP Go Debugger - Web Server Example

This example demonstrates how to debug a simple HTTP web server using the MCP Go Debugger with conditional breakpoints.

## Overview

The example includes:
- `webserver.go` - A simple HTTP server (port 8080) that:
  - Tracks request count with a global counter
  - Accepts a `name` query parameter
  - Prints formatted messages to stdout
  - Responds with personalized greetings

## Prerequisites

```bash
# Build the debugger
cd /mcp-go-debugger
make build
```

## Usage

### Method 1: Using MCP Server (Recommended)

#### Step 1: Start the MCP Server

```bash
cd /mcp-go-debugger
./bin/mcp-go-debugger
```

The server communicates via JSON-RPC on stdio.

#### Step 2: Start Debug Session

```json
{
  "tool": "debug",
  "arguments": {
    "file": "/mcp-go-debugger/testdata/webserver/webserver.go"
  }
}
```

#### Step 3: Set Breakpoints

**Regular breakpoint:**
```json
{
  "tool": "set_breakpoint",
  "arguments": {
    "file": "/mcp-go-debugger/testdata/webserver/webserver.go",
    "line": 14
  }
}
```

**Conditional breakpoint** (only breaks when `requestCount > 2`):
```json
{
  "tool": "set_breakpoint",
  "arguments": {
    "file": "/mcp-go-debugger/testdata/webserver/webserver.go",
    "line": 26,
    "condition": "requestCount > 2"
  }
}
```

#### Step 4: Continue Execution

```json
{
  "tool": "continue"
}
```

The server will start listening on http://localhost:8080

#### Step 5: Trigger Breakpoints

In another terminal, make HTTP requests:

```bash
# Request 1 - hits breakpoint at line 14
curl http://localhost:8080/?name=Alice

# Request 2 - hits breakpoint at line 14
curl http://localhost:8080/?name=Bob

# Request 3 - hits BOTH breakpoints (requestCount > 2)
curl http://localhost:8080/?name=Charlie

# Request 4 - hits conditional breakpoint only
curl http://localhost:8080/?name=Dave
```

#### Step 6: Inspect Variables

When stopped at a breakpoint:

```json
{
  "tool": "eval_variable",
  "arguments": {
    "name": "requestCount"
  }
}
```

```json
{
  "tool": "eval_variable",
  "arguments": {
    "name": "name"
  }
}
```

```json
{
  "tool": "eval_variable",
  "arguments": {
    "name": "message"
  }
}
```

#### Step 7: Step Through Code

```json
{"tool": "step_over"}
```

```json
{"tool": "step"}
```

```json
{"tool": "step_out"}
```

#### Step 8: Get Program Output

```json
{
  "tool": "get_debugger_output"
}
```

This returns captured stdout/stderr from the server.

#### Step 9: Close Session

```json
{
  "tool": "close"
}
```

This will:
- Detach from the debugged process
- Clean up temporary debug binaries
- Close the debug server

### Method 2: Using Go API Directly

Create a test program:

```go
package main

import (
	"fmt"
	"path/filepath"
	"github.com/sunfmin/mcp-go-debugger/pkg/debugger"
)

func main() {
	// Get absolute path
	webserverPath, _ := filepath.Abs("webserver.go")

	// Create debugger client
	client := debugger.NewClient()

	// Start debug session
	debugResp := client.DebugSourceFile(webserverPath, nil)
	fmt.Printf("Debug session started: %s\n", debugResp.Status)

	// Set regular breakpoint
	bp1 := client.SetBreakpoint(webserverPath, 14, "")
	fmt.Printf("Breakpoint %d set\n", bp1.Breakpoint.ID)

	// Set conditional breakpoint
	bp2 := client.SetBreakpoint(webserverPath, 26, "requestCount > 2")
	fmt.Printf("Conditional BP %d set: %s\n", bp2.Breakpoint.ID, bp2.Breakpoint.Condition)

	// List breakpoints
	list := client.ListBreakpoints()
	fmt.Printf("Total breakpoints: %d\n", len(list.Breakpoints))

	// Continue execution
	client.Continue()

	// ... server is now running in debug mode on :8080
}
```

Run it:
```bash
cd /mcp-go-debugger/testdata/webserver
go run debug_webserver.go
```

While it's running, trigger requests in another terminal:
```bash
curl http://localhost:8080/?name=TestUser
```

## Conditional Breakpoint Examples

The debugger supports any valid Go expression as a condition:

```go
// Only break when request count exceeds threshold
condition: "requestCount > 5"

// Only break for specific names
condition: "name == \"Alice\""

// Only break when name is empty
condition: "name == \"\""

// Combine conditions
condition: "requestCount > 2 && name != \"World\""

// Check request properties
condition: "len(r.URL.Query()) > 1"
```

## Debugging Tips

### Good Breakpoint Locations

- **Line 12** (`requestCount++`) - Watch counter increment
- **Line 14-17** - Inspect HTTP request details
- **Line 20-23** - Check query parameter extraction
- **Line 26** - Verify message formatting
- **Line 36** - Confirm response sent

### Variables to Inspect

When stopped at a breakpoint, you can evaluate:

- `requestCount` - Total number of requests processed
- `name` - The extracted name parameter
- `message` - The response message being prepared
- `r.RemoteAddr` - Client's remote address
- `r.URL.Query()` - All query parameters

### Example Debugging Session

```bash
# Set breakpoint at request handler start
set_breakpoint(file="webserver.go", line=12)

# Continue until breakpoint
continue

# Step through incrementing counter
step

# Evaluate current count
eval_variable("requestCount")

# Step to query parameter extraction
step_over

# Inspect the name variable
eval_variable("name")

# Continue to response
continue
```

## What Gets Demonstrated

✓ Starting a debug session from source file
✓ Setting regular breakpoints
✓ Setting conditional breakpoints
✓ Listing all breakpoints with conditions
✓ Continuing execution
✓ Inspecting variables at breakpoints
✓ Stepping through code (step, step_over, step_out)
✓ Capturing program output (stdout/stderr)
✓ Closing debug sessions cleanly

## Expected Output

When requests hit breakpoints, you'll see:

1. **Breakpoint hit** - Current location shown
2. **Local variables** - Automatically included in context
3. **Stop reason** - "breakpoint hit" or "conditional breakpoint: requestCount > 2"
4. **Execution position** - File, line, function name

The conditional breakpoint will only trigger on requests #3, #4, etc. (when `requestCount > 2`), while the regular breakpoint triggers on every request.

## List All Breakpoints

```json
{
  "tool": "list_breakpoints"
}
```

Response shows all breakpoints with their conditions:
```json
{
  "breakpoints": [
    {"id": 1, "location": "webserver.go:14", "condition": ""},
    {"id": 2, "location": "webserver.go:26", "condition": "requestCount > 2"}
  ]
}
```
