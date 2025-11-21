# Quick Start: Debugging the Example Web Server

This guide demonstrates how to debug the example Hello World web server using the MCP Go Debugger with conditional breakpoints.

## Files

- `webserver.go` - A simple HTTP server (port 8080) that:
  - Tracks request count
  - Accepts a `name` query parameter
  - Prints formatted messages to stdout
  - Responds with personalized greetings

## Prerequisites

```bash
# Build the debugger
cd /mcp-go-debugger
make build
```

## Method 1: Using Go API Directly

Create a simple test program:

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

## Method 2: Using MCP Server (Recommended)

### Step 1: Start the MCP Server

```bash
cd /mcp-go-debugger
./bin/mcp-go-debugger
```

The server communicates via JSON-RPC on stdio.

### Step 2: Send Debug Commands

Send JSON-RPC requests to stdin:

#### 2.1 Start Debug Session

```json
{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"debug","arguments":{"file":"/mcp-go-debugger/testdata/webserver/webserver.go"}}}
```

#### 2.2 Set Regular Breakpoint at Line 14

```json
{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"set_breakpoint","arguments":{"file":"/mcp-go-debugger/testdata/webserver/webserver.go","line":14}}}
```

#### 2.3 Set Conditional Breakpoint

Only break when `requestCount > 2`:

```json
{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"set_breakpoint","arguments":{"file":"/mcp-go-debugger/testdata/webserver/webserver.go","line":26,"condition":"requestCount > 2"}}}
```

#### 2.4 List All Breakpoints

```json
{"jsonrpc":"2.0","id":4,"method":"tools/call","params":{"name":"list_breakpoints","arguments":{}}}
```

Response will show:
```json
{
  "breakpoints": [
    {"id": 1, "location": "webserver.go:14", "condition": ""},
    {"id": 2, "location": "webserver.go:26", "condition": "requestCount > 2"}
  ]
}
```

#### 2.5 Continue Execution

```json
{"jsonrpc":"2.0","id":5,"method":"tools/call","params":{"name":"continue","arguments":{}}}
```

The server will start listening on http://localhost:8080

### Step 3: Trigger Breakpoints

In another terminal, make HTTP requests:

```bash
# Request 1 - will hit BP at line 14
curl http://localhost:8080/?name=Alice

# Request 2 - will hit BP at line 14
curl http://localhost:8080/?name=Bob

# Request 3 - will hit BOTH breakpoints (requestCount > 2)
curl http://localhost:8080/?name=Charlie

# Request 4 - will hit conditional BP (requestCount > 2)
curl http://localhost:8080/?name=Dave
```

### Step 4: Inspect Variables

When stopped at a breakpoint:

```json
{"jsonrpc":"2.0","id":6,"method":"tools/call","params":{"name":"eval_variable","arguments":{"name":"requestCount"}}}
```

```json
{"jsonrpc":"2.0","id":7,"method":"tools/call","params":{"name":"eval_variable","arguments":{"name":"name"}}}
```

```json
{"jsonrpc":"2.0","id":8,"method":"tools/call","params":{"name":"eval_variable","arguments":{"name":"message"}}}
```

### Step 5: Step Through Code

```json
{"jsonrpc":"2.0","id":9,"method":"tools/call","params":{"name":"step_over","arguments":{}}}
```

```json
{"jsonrpc":"2.0","id":10,"method":"tools/call","params":{"name":"step","arguments":{}}}
```

### Step 6: Get Program Output

```json
{"jsonrpc":"2.0","id":11,"method":"tools/call","params":{"name":"get_debugger_output","arguments":{}}}
```

This returns captured stdout/stderr from the server.

### Step 7: Close Session

```json
{"jsonrpc":"2.0","id":12,"method":"tools/call","params":{"name":"close","arguments":{}}}
```

## Conditional Breakpoint Examples

The debugger supports any valid Go expression as a condition:

```json
// Only break when count exceeds threshold
{"condition": "requestCount > 5"}

// Only break for specific names
{"condition": "name == \"Alice\""}

// Only break when name is empty
{"condition": "name == \"\""}

// Combine conditions
{"condition": "requestCount > 2 && name != \"World\""}

// Check request properties
{"condition": "len(r.URL.Query()) > 1"}
```

## Debugging Tips

1. **Line 12** (`requestCount++`) - Watch counter increment
2. **Line 14-17** - Inspect HTTP request details
3. **Line 20-23** - Check query parameter extraction
4. **Line 26** - Verify message formatting
5. **Line 36** - Confirm response sent

## What Gets Demonstrated

✓ Starting a debug session from source file
✓ Setting regular breakpoints
✓ Setting conditional breakpoints
✓ Listing all breakpoints with conditions
✓ Continuing execution
✓ Inspecting variables at breakpoints
✓ Stepping through code
✓ Capturing program output (stdout/stderr)
✓ Closing debug sessions cleanly

## Expected Output

When requests hit breakpoints, you'll see:

1. **Breakpoint hit** - Current location shown
2. **Local variables** - Automatically included in context
3. **Stop reason** - "breakpoint hit" or "conditional breakpoint: requestCount > 2"
4. **Execution position** - File, line, function name

The conditional breakpoint will only trigger on requests #3, #4, etc. (when `requestCount > 2`), while the regular breakpoint triggers on every request.
