# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

MCP Go Debugger is a Model Context Protocol (MCP) server that provides Go debugging capabilities through the Delve debugger. It exposes debugging tools (launch, attach, breakpoints, stepping, variable inspection) via the MCP protocol for integration with AI assistants like Claude.

## Build and Development Commands

### Building
```bash
make build          # Compiles to bin/mcp-go-debugger
make install        # Installs to $GOPATH/bin
```

### Testing
```bash
make test           # Run all tests with go test -v ./...
go test ./pkg/mcp   # Run tests for specific package
```

### Running
```bash
make run            # Build and run the binary
./bin/mcp-go-debugger  # Run directly (starts stdio MCP server)
```

### Cleanup
```bash
make clean          # Remove binaries and clean build artifacts
```

## Architecture

### Core Components

1. **MCP Server Layer** (`pkg/mcp/server.go`)
   - Entry point for MCP protocol communication
   - Registers all debugging tools (launch, attach, set_breakpoint, etc.)
   - Handles tool call routing and parameter marshaling
   - Converts Delve responses to JSON for MCP clients

2. **Debugger Client** (`pkg/debugger/client.go`)
   - Wraps Delve's RPC client
   - Manages debug server lifecycle (starting/stopping)
   - Captures program stdout/stderr via pipes
   - Maintains debug session state (target, pid, server instance)

3. **Debug Operations** (split across `pkg/debugger/` files)
   - `program.go`: Launch, attach, debug source files/tests, close sessions
   - `breakpoints.go`: Set, list, remove breakpoints
   - `execution.go`: Continue, step, step over, step out
   - `variables.go`: Evaluate variables, list local/function args
   - `output.go`: Capture and retrieve program output
   - `helpers.go`: Utility functions for location formatting, state inspection

4. **Type System** (`pkg/types/debug_types.go`)
   - Defines all response types (LaunchResponse, BreakpointResponse, etc.)
   - Each response includes a `DebugContext` with timestamp, operation, current location, local variables, stop reason
   - Types are LLM-friendly: internal Delve types excluded from JSON, human-readable fields included
   - All Delve API types (`*api.DebuggerState`, `*api.Variable`, `*api.Breakpoint`) stored but not serialized

### Key Design Patterns

**Response Structure**: All operations return structured responses with:
- `status`: "success" or "error"
- `context`: DebugContext with current execution state
- Operation-specific fields (breakpoint details, variable values, etc.)

**Delve Integration**:
- Starts Delve server with RPC API on dynamic port
- Uses `rpc2.RPCClient` for all debugger communication
- Configures redirectors for stdout/stderr capture
- Builds test/debug binaries with `-gcflags all=-N` (disables optimizations)

**Output Capture**:
- Creates pipes via `proc.Redirector()` for stdout/stderr
- Goroutines scan output and write to buffers + channels
- Output available via `get_debugger_output` tool

## Important Implementation Details

### Breakpoints
- `set_breakpoint` supports optional conditional breakpoints via `condition` parameter
- Conditions use Go expression syntax (e.g., `count > 5`, `username == "admin"`)
- Condition is passed directly to Delve's `api.Breakpoint.Cond` field
- All breakpoint responses include the condition (if set) in the `Condition` field

### Test Debugging
- `debug_test` tool compiles tests with `gobuild.GoTestBuildCombinedOutput`
- Changes to test directory before building (test packages expect to be built from their location)
- Escapes test name with `regexp.QuoteMeta` and uses `-test.run=^TestName$` for exact matching
- Always includes `-test.v` for verbose output

### Variable Evaluation
- Uses `api.EvalScope` with goroutine ID and frame number
- `LoadConfig.MaxStructFields = -1` loads all struct fields
- Special formatting for structs (field:value pairs) and arrays/slices
- `depth` parameter controls `MaxVariableRecurse` for nested structures

### Binary Cleanup
- Compiled binaries stored in temp paths via `gobuild.DefaultDebugBinaryPath`
- Cleaned up in `Close()` with `gobuild.Remove(c.target)`
- Critical: Always clean up debug binaries after session ends

### Session State Management
- Only one debug session per client instance
- Check `c.client != nil` to determine if session active
- `Close()` must: detach client, stop server, clean up binary, reset state
- Timeouts on detach/stop operations to prevent hangs

## Common Development Workflows

### Adding a New Tool
1. Define tool schema in `registerTools()` using `mcp.NewTool()`
2. Add handler method to `MCPDebugServer` (receives `context.Context` and `mcp.CallToolRequest`)
3. Extract parameters from `request.Params.Arguments` with type assertions
4. Call corresponding `debugger.Client` method
5. Return `*mcp.CallToolResult` via `newToolResultJSON()`

### Modifying Response Types
1. Update type definition in `pkg/types/debug_types.go`
2. Ensure LLM-friendly fields are exported and have JSON tags
3. Exclude internal Delve types with `json:"-"`
4. Update response creation helper in relevant `pkg/debugger/` file

### Debugging Test Failures
The repository includes test programs in `testdata/calculator/`:
- `calculator.go`: Simple calculator functions
- `calculator_test.go`: Test suite for debugging

Use these for manual testing:
```bash
# Debug a test function
./bin/mcp-go-debugger
# Then via MCP: debug_test with testdata/calculator/calculator_test.go, TestAdd
```

## Dependencies

- `github.com/go-delve/delve`: Delve debugger (v1.24.1)
- `github.com/mark3labs/mcp-go`: MCP Go SDK (v0.15.0)

## Logging

Uses custom logger package (`pkg/logger/logger.go`). Logs written to:
- `mcp-go-debugger.log` in working directory
- stderr (for MCP protocol errors)

Debug level logging throughout for troubleshooting.