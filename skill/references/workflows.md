# Detailed Debugging Workflows

This document provides comprehensive step-by-step procedures for common Go debugging scenarios using the delve-mcp server.

## Table of Contents

1. [Interactive Development Debugging](#interactive-development-debugging)
2. [CI/CD Test Failure Investigation](#cicd-test-failure-investigation)
3. [Production Issue Investigation](#production-issue-investigation)
4. [Race Condition Debugging](#race-condition-debugging)
5. [Memory Leak Investigation](#memory-leak-investigation)
6. [Goroutine Debugging](#goroutine-debugging)
7. [Integration Test Debugging](#integration-test-debugging)

---

## Interactive Development Debugging

**Scenario:** You're developing a new feature and want to understand how your code executes.

### Workflow Steps

**1. Prepare your code**
- Identify the entry point (usually `main.go` or a specific function)
- Know which variables you want to inspect
- Identify critical execution points

**2. Start debug session**
```
debug(file: "/absolute/path/to/main.go", args: ["--config", "dev.yaml"])
```

**Response:**
```json
{
  "status": "success",
  "context": {
    "operation": "debug",
    "stopReason": "process is stopped",
    "currentLocation": "At /path/to/main.go:1 in main.main"
  }
}
```

**3. Set strategic breakpoints**
```
# Entry point of your new feature
set_breakpoint(file: "/path/to/handler.go", line: 45)

# Business logic
set_breakpoint(file: "/path/to/service.go", line: 123)

# Error handling
set_breakpoint(file: "/path/to/handler.go", line: 78)
```

**4. Start execution**
```
continue()
```

**5. When breakpoint hits, inspect state**
```
# Check automatic local variables
→ Response includes context.localVariables

# Inspect specific variables
eval_variable(name: "user", depth: 2)
eval_variable(name: "request", depth: 1)
eval_variable(name: "err")
```

**6. Navigate through code**
```
# Step into function to see implementation
step()

# Step over helper functions you trust
step_over()

# Step out if you're deep in a call stack
step_out()

# Continue to next breakpoint
continue()
```

**7. Verify output**
```
get_debugger_output()
→ See print statements and log output
```

**8. Close session**
```
close()
```

### Example Session

```
debug(file: "/app/server.go", args: ["--port", "8080"])
→ Debug session started

set_breakpoint(file: "/app/handlers/user.go", line: 23)
→ Breakpoint 1 set at CreateUser function

set_breakpoint(file: "/app/handlers/user.go", line: 45)
→ Breakpoint 2 set at validation logic

continue()
→ Server started, listening...

# Trigger endpoint: curl -X POST http://localhost:8080/users -d '{"name":"Alice"}'

→ Stopped at handlers/user.go:23
→ localVariables: req={Name: "Alice"}, w=ResponseWriter

eval_variable(name: "req", depth: 2)
→ {
    Name: "Alice",
    Email: "",
    Role: ""
  }
  # Found issue: Email is empty!

step_over()  # Move to next line
step_over()  # Move to validation

→ Stopped at handlers/user.go:45
→ localVariables: valid=false, errors=["email is required"]

# Confirmed: Validation correctly catches missing email

continue()
→ Program continued

close()
→ Session closed
```

---

## CI/CD Test Failure Investigation

**Scenario:** A test fails in CI but passes locally. You need to understand why.

### Workflow Steps

**1. Identify the failing test**
```bash
# From CI logs or local run
go test -v ./...
→ FAIL: TestUserCreation (0.01s)
    handler_test.go:45: expected status 201, got 500
```

**2. Start test debug session**
```
debug_test(
  testfile: "/absolute/path/to/handlers/user_test.go",
  testname: "TestUserCreation",
  testflags: ["-v"]
)
```

**3. Set breakpoints in test and implementation**
```
# In the test where assertion fails
set_breakpoint(file: "/path/to/handlers/user_test.go", line: 45)

# In the implementation being tested
set_breakpoint(file: "/path/to/handlers/user.go", line: 30)

# At error handling paths
set_breakpoint(file: "/path/to/handlers/user.go", line: 67, condition: "err != nil")
```

**4. Start test execution**
```
continue()
→ Test starts running
```

**5. Inspect test setup**
```
→ Stopped at user_test.go:45

eval_variable(name: "expected")
→ 201

eval_variable(name: "actual")
→ 500

eval_variable(name: "response", depth: 2)
→ {Status: 500, Body: "Internal Server Error", Headers: {...}}
```

**6. Step into implementation**
```
continue()
→ Stopped at user.go:30

eval_variable(name: "req", depth: 3)
→ {Name: "TestUser", Email: "", CreatedAt: null}
  # Email is empty in test!

step_over()
step_over()

→ Stopped at user.go:67 (err != nil)

eval_variable(name: "err")
→ "email validation failed"
  # Found root cause!
```

**7. Verify fix hypothesis**
```
# Note the issue, close session
close()

# Fix test to include email
# Re-run without debugging to verify
```

### Common Test Debugging Patterns

**Pattern 1: Compare Expected vs Actual**
```
# At assertion point
eval_variable(name: "expected")
eval_variable(name: "actual")
eval_variable(name: "diff")  # if test creates a diff
```

**Pattern 2: Trace Mock Behavior**
```
# Check if mocks are configured correctly
eval_variable(name: "mockDB")
eval_variable(name: "mockDB.calls")  # if tracking calls
```

**Pattern 3: Inspect Test Fixtures**
```
# Verify test data is correct
eval_variable(name: "testUser", depth: 3)
eval_variable(name: "testData", depth: 2)
```

---

## Production Issue Investigation

**Scenario:** Users report errors in production. You need to debug a live server without downtime.

### Workflow Steps

**1. Find the running process**
```bash
# On production server
ps aux | grep myserver
→ user  28026  0.5  2.1  myserver

# Or use pgrep
pgrep -fl myserver
→ 28026 /usr/bin/myserver
```

**2. Attach to the process**
```
attach(pid: 28026)
```

**Response:**
```json
{
  "status": "success",
  "context": {
    "operation": "attach",
    "stopReason": "process is stopped"
  },
  "pid": 28026
}
```

**IMPORTANT:** Process is now paused! Users cannot access the server until you continue.

**3. Set conditional breakpoints (to filter noise)**
```
# Only break for specific user
set_breakpoint(
  file: "/app/handlers/api.go",
  line: 45,
  condition: "userID == \"problem_user_123\""
)

# Only break on errors
set_breakpoint(
  file: "/app/handlers/api.go",
  line: 89,
  condition: "err != nil"
)

# Only break on specific endpoints
set_breakpoint(
  file: "/app/router.go",
  line: 123,
  condition: "path == \"/api/orders\""
)
```

**4. Resume server (critical!)**
```
continue()
→ Server resumed, handling requests normally
```

**5. Wait for breakpoint to trigger**
```
# Monitor logs or metrics
# Reproduce issue manually if possible
# Wait for condition to occur naturally
```

**6. When stopped, inspect state quickly**
```
→ Stopped at api.go:45 (condition: userID == "problem_user_123")

# Check local variables (automatic)
→ context.localVariables shows all locals

# Inspect specific values
eval_variable(name: "request", depth: 2)
eval_variable(name: "session")
eval_variable(name: "err")
```

**7. Capture findings and resume**
```
# Take notes on variable values
# Copy important data

# Resume server immediately
continue()
```

**8. Remove breakpoints when done**
```
list_breakpoints()
→ [{id: 1, ...}, {id: 2, ...}]

remove_breakpoint(id: 1)
remove_breakpoint(id: 2)

continue()
```

**9. Detach cleanly**
```
close()
→ Server continues running normally
```

### Production Safety Checklist

- ✓ Use conditional breakpoints to avoid pausing for all traffic
- ✓ Resume execution immediately after attaching
- ✓ Inspect variables quickly (< 30 seconds at breakpoint)
- ✓ Remove breakpoints when investigation is complete
- ✓ Always call `close()` to cleanly detach
- ✗ Never leave process paused (users can't access service)
- ✗ Never set unconditional breakpoints on hot paths
- ✗ Never use deep variable inspection (slow)

---

## Race Condition Debugging

**Scenario:** Intermittent failures suggest a race condition between goroutines.

### Workflow Steps

**1. Enable race detector first**
```bash
# Run with race detector to confirm
go test -race ./...
# or
go run -race main.go
```

**2. Debug with breakpoints at shared data access**
```
debug(file: "/path/to/main.go")

# Set breakpoints at writes to shared data
set_breakpoint(file: "/app/cache.go", line: 45)  # cache.Set()
set_breakpoint(file: "/app/cache.go", line: 67)  # cache.Delete()

# Set breakpoints at reads
set_breakpoint(file: "/app/cache.go", line: 89)  # cache.Get()
```

**3. Inspect goroutine context**
```
continue()
→ Stopped at cache.go:45

# Check which goroutine we're in
eval_variable(name: "key")
eval_variable(name: "value")

# Check concurrent state
eval_variable(name: "cache.mu")  # Is mutex locked?
eval_variable(name: "cache.data", depth: 1)  # Current cache state
```

**4. Step carefully through critical sections**
```
# Step through mutex acquisition
step()
→ cache.mu.Lock()

step()
→ cache.data[key] = value

step()
→ cache.mu.Unlock()

# Verify state after operation
eval_variable(name: "cache.data", depth: 1)
```

**5. Look for missing synchronization**
```
# Common issues:
# - Reading without lock
# - Writing without lock
# - Unlocking before write completes
# - Wrong mutex for data structure
```

### Race Detection Patterns

**Pattern 1: Verify Mutex Usage**
```
# Before critical section
eval_variable(name: "mu.state")  # Should be 0 (unlocked)

# After Lock()
eval_variable(name: "mu.state")  # Should be 1 (locked)

# After Unlock()
eval_variable(name: "mu.state")  # Should be 0 (unlocked)
```

**Pattern 2: Check Goroutine Coordination**
```
# At channel send
eval_variable(name: "ch")
→ {len: 0, cap: 10}  # Channel state

# At channel receive
eval_variable(name: "ch")
→ {len: 1, cap: 10}  # Message waiting
```

---

## Memory Leak Investigation

**Scenario:** Application memory grows over time. Need to find what's not being freed.

### Workflow Steps

**1. Attach to running process with memory issue**
```
attach(pid: 28026)
```

**2. Set breakpoints at allocation points**
```
# Where objects are created
set_breakpoint(file: "/app/cache.go", line: 23)  # cache creation

# Where objects should be freed
set_breakpoint(file: "/app/cache.go", line: 56)  # cache cleanup
```

**3. Inspect data structure growth**
```
continue()
→ Stopped at cache.go:23

eval_variable(name: "cache.data", depth: 1)
→ map[string]interface{}{
    "key1": {...},
    "key2": {...},
    ...
    # Count entries
  }

eval_variable(name: "len(cache.data)")
→ 10000  # Growing over time?
```

**4. Check if cleanup is running**
```
# Continue to cleanup breakpoint
continue()

→ Does breakpoint hit?
  YES: Cleanup runs, but may be insufficient
  NO: Cleanup not running (found the bug!)
```

**5. Inspect cleanup logic**
```
→ Stopped at cache.go:56

step()
step()
→ In cleanup function

eval_variable(name: "toDelete")
→ []string{"key1", "key2"}  # Only 2 keys?

eval_variable(name: "cache.data", depth: 1)
→ map still has 10000 entries  # Cleanup is ineffective!
```

### Memory Leak Indicators

- Maps/slices that only grow, never shrink
- Goroutines that start but never exit
- Channels that fill up but never drain
- Objects with no cleanup path
- Timers/tickers that are never stopped

---

## Goroutine Debugging

**Scenario:** Application hangs or has deadlock. Need to inspect goroutine state.

### Workflow Steps

**1. Attach to hanging process**
```
attach(pid: 28026)
```

**2. Check where execution is stuck**
```
# Process will be paused
→ context.currentLocation shows where it stopped

# If it's in a select/channel operation
eval_variable(name: "ch")
→ {len: 0, cap: 0}  # Unbuffered channel, waiting for sender/receiver
```

**3. Set breakpoints at goroutine creation**
```
set_breakpoint(file: "/app/worker.go", line: 34)  # go startWorker()

continue()

→ Stopped at worker.go:34
eval_variable(name: "workerCount")
→ 1000  # Too many workers?
```

**4. Inspect channel operations**
```
# At channel send
set_breakpoint(file: "/app/worker.go", line: 56)

continue()
→ Stopped at channel send

eval_variable(name: "workCh")
→ {len: 100, cap: 100}  # Channel is full! Sender will block
```

**5. Trace deadlock**
```
# Goroutine 1: waiting to send on workCh (full)
# Goroutine 2: waiting to receive on resultCh (empty)
# Deadlock if no consumer for resultCh!

eval_variable(name: "resultCh")
→ {len: 0, cap: 0}  # No one receiving results
```

### Goroutine Debugging Patterns

**Pattern 1: Check Channel State**
```
eval_variable(name: "ch")
→ {
    len: 5,    # Current items in channel
    cap: 10    # Channel capacity
  }

# Interpretations:
# len == cap → Channel full, senders will block
# len == 0 → Channel empty, receivers will block
```

**Pattern 2: Verify Goroutine Completion**
```
# At goroutine start
eval_variable(name: "activeCount")
→ 10

# At goroutine end
eval_variable(name: "activeCount")
→ 9  # Decremented? Goroutine exited properly
```

**Pattern 3: Inspect Context Cancellation**
```
eval_variable(name: "ctx")
→ {Done: <chan>, Err: nil}  # Context still active

eval_variable(name: "ctx.Err()")
→ "context canceled"  # Context was cancelled
```

---

## Integration Test Debugging

**Scenario:** Integration test fails due to complex interactions between components.

### Workflow Steps

**1. Start integration test in debug mode**
```
debug_test(
  testfile: "/path/to/integration_test.go",
  testname: "TestUserRegistrationFlow",
  testflags: ["-v", "-count=1"]
)
```

**2. Set breakpoints across layers**
```
# Test level
set_breakpoint(file: "/path/to/integration_test.go", line: 67)

# API level
set_breakpoint(file: "/app/handlers/register.go", line: 23)

# Service level
set_breakpoint(file: "/app/services/user.go", line: 45)

# Database level
set_breakpoint(file: "/app/db/users.go", line: 78)
```

**3. Trace request flow**
```
continue()
→ Stopped at integration_test.go:67

eval_variable(name: "request", depth: 2)
→ {Username: "testuser", Email: "test@example.com"}

continue()
→ Stopped at handlers/register.go:23

eval_variable(name: "req", depth: 2)
→ Request arrived at handler

step_over()
continue()
→ Stopped at services/user.go:45

eval_variable(name: "user", depth: 2)
→ {Username: "testuser", Email: "test@example.com", HashedPassword: "..."}
  # Verify transformation

continue()
→ Stopped at db/users.go:78

eval_variable(name: "query")
→ "INSERT INTO users ..."
  # Verify SQL generated correctly
```

**4. Inspect state at each layer**
```
# At each breakpoint, verify:
# - Data transformation is correct
# - No data loss between layers
# - Errors are propagated properly
# - Side effects occur as expected
```

**5. Check integration points**
```
# HTTP client
eval_variable(name: "httpClient.Timeout")

# Database connection
eval_variable(name: "db.Stats()")

# External services
eval_variable(name: "apiClient.BaseURL")
```

---

## Advanced Techniques

### Conditional Breakpoint Strategies

**Break on Error Paths Only**
```
set_breakpoint(
  file: "/app/handler.go",
  line: 89,
  condition: "err != nil"
)
```

**Break on Specific Data**
```
set_breakpoint(
  file: "/app/processor.go",
  line: 45,
  condition: "len(items) == 0"
)
```

**Break on State Changes**
```
set_breakpoint(
  file: "/app/state.go",
  line: 67,
  condition: "oldState != newState"
)
```

### Variable Inspection Strategies

**Shallow Inspection (Fast)**
```
eval_variable(name: "user", depth: 1)
→ Quick overview of top-level fields
```

**Medium Inspection (Balanced)**
```
eval_variable(name: "request", depth: 2)
→ See nested objects, but not too deep
```

**Deep Inspection (Slow, Thorough)**
```
eval_variable(name: "config", depth: 5)
→ Complete structure, use sparingly
```

### Output Debugging

**Correlate Debug State with Logs**
```
# Set breakpoint
set_breakpoint(file: "/app/handler.go", line: 45)

continue()
→ Stopped

# Inspect state
eval_variable(name: "requestID")
→ "req-12345"

# Check output for related logs
get_debugger_output()
→ [INFO] Processing request req-12345
  [DEBUG] User authenticated: user-789
  [ERROR] Database timeout
  # Correlate debug state with log messages
```

---

## Workflow Selection Guide

| Scenario | Workflow |
|----------|----------|
| Developing new feature | Interactive Development Debugging |
| Test fails in CI | CI/CD Test Failure Investigation |
| Production error reports | Production Issue Investigation |
| Intermittent failures | Race Condition Debugging |
| Memory grows over time | Memory Leak Investigation |
| Application hangs | Goroutine Debugging |
| Multi-component test fails | Integration Test Debugging |

---

## Next Steps

- See `patterns.md` for common debugging patterns
- See `tools.md` for complete tool reference
- See `examples.md` for real-world examples
