# Real-World Debugging Examples

Complete debugging scenarios from the mcp-go-debugger examples directory.

## Table of Contents

1. [Example Web Server](#example-web-server)
2. [Test Debugging](#test-debugging)
3. [Production Debugging](#production-debugging)
4. [Common Scenarios](#common-scenarios)

---

## Example Web Server

### Scenario: Debugging the Hello World HTTP Server

**File:** `testdata/webserver/webserver.go`

**Description:** Simple HTTP server that:
- Listens on port 8080
- Tracks request count
- Accepts `name` query parameter
- Returns personalized greeting

### Complete Debugging Session

**Step 1: Launch the debugger**
```
mcp__delve-mcp__debug(file: "/mcp-go-debugger/testdata/webserver/webserver.go")

Response:
{
  "status": "success",
  "context": {
    "operation": "debug",
    "currentLocation": "At webserver.go:43 in main.main",
    "stopReason": "process is stopped"
  }
}
```

**Step 2: Set breakpoints at key locations**
```
# Track request count increment
mcp__delve-mcp__set_breakpoint(
  file: "/mcp-go-debugger/testdata/webserver/webserver.go",
  line: 13
)

# Check query parameter handling
mcp__delve-mcp__set_breakpoint(
  file: "/mcp-go-debugger/testdata/webserver/webserver.go",
  line: 22
)

# Conditional: only break for Alice
mcp__delve-mcp__set_breakpoint(
  file: "/mcp-go-debugger/testdata/webserver/webserver.go",
  line: 24,
  condition: "name == \"Alice\""
)

# Conditional: only after 2 requests
mcp__delve-mcp__set_breakpoint(
  file: "/mcp-go-debugger/testdata/webserver/webserver.go",
  line: 28,
  condition: "requestCount > 2"
)
```

**Step 3: Start the server**
```
mcp__delve-mcp__continue()

Response:
{
  "status": "success",
  "context": {
    "operation": "continue",
    "stopReason": "server running"
  }
}
```

**Step 4: Make HTTP requests**
```bash
# Terminal 2
curl http://localhost:8080/
curl http://localhost:8080/?name=Alice
curl http://localhost:8080/?name=Bob
curl http://localhost:8080/?name=Charlie
```

**Step 5: Hit first breakpoint (requestCount++)**
```
# After first curl
Response:
{
  "status": "success",
  "context": {
    "currentLocation": "At webserver.go:13 in main.helloHandler",
    "localVariables": [
      {"name": "w", "value": "...", "type": "net/http.ResponseWriter"},
      {"name": "r", "value": "...", "type": "*net/http.Request"}
    ],
    "stopReason": "hit breakpoint"
  }
}

# Inspect
mcp__delve-mcp__eval_variable(name: "requestCount")
→ {value: "0", type: "int"}  # Before increment

mcp__delve-mcp__step_over()  # Execute increment

mcp__delve-mcp__eval_variable(name: "requestCount")
→ {value: "1", type: "int"}  # After increment

mcp__delve-mcp__continue()
```

**Step 6: Hit conditional breakpoint (name == "Alice")**
```
# After curl with name=Alice
Response:
{
  "status": "success",
  "context": {
    "currentLocation": "At webserver.go:24 in main.helloHandler",
    "localVariables": [
      {"name": "name", "value": "Alice", "type": "string"}
    ],
    "stopReason": "hit breakpoint (conditional: name == \"Alice\")"
  }
}

# Inspect the request
mcp__delve-mcp__eval_variable(name: "r.URL.Query()", depth: 2)
→ {
    value: "map[string][]string{
      \"name\": [\"Alice\"]
    }"
  }

mcp__delve-mcp__eval_variable(name: "name")
→ {value: "Alice", type: "string"}

mcp__delve-mcp__continue()
```

**Step 7: Hit conditional breakpoint (requestCount > 2)**
```
# After 3rd curl
Response:
{
  "status": "success",
  "context": {
    "currentLocation": "At webserver.go:28 in main.helloHandler",
    "localVariables": [
      {"name": "name", "value": "Charlie", "type": "string"},
      {"name": "requestCount", "value": "3", "type": "int"}
    ],
    "stopReason": "hit breakpoint (conditional: requestCount > 2)"
  }
}

mcp__delve-mcp__eval_variable(name: "message")
→ {value: "Hello, Charlie! This is request #3", type: "string"}

mcp__delve-mcp__continue()
```

**Step 8: Get program output**
```
mcp__delve-mcp__get_debugger_output()

Response:
{
  "output": "Starting Hello World web server on http://localhost:8080\nTry visiting:\n  http://localhost:8080/\n  http://localhost:8080/?name=Alice\n  http://localhost:8080/?name=Bob\n\nServer is listening on port :8080...\n[15:04:05] Request #1 received from [::1]:50123\n[15:04:05] Sending response: Hello, World! This is request #1\n[15:04:05] Response sent successfully\n[15:04:06] Request #2 received from [::1]:50124\n[15:04:06] Sending response: Hello, Alice! This is request #2\n[15:04:06] Response sent successfully\n[15:04:07] Request #3 received from [::1]:50125\n[15:04:07] Sending response: Hello, Charlie! This is request #3\n[15:04:07] Response sent successfully\n"
}
```

**Step 9: Close session**
```
mcp__delve-mcp__close()
```

### Key Learnings

1. **Conditional breakpoints** filter noise in busy handlers
2. **Local variables** appear automatically in context
3. **Query parameters** accessible via `r.URL.Query()`
4. **Request count** tracked across multiple requests
5. **Output capture** shows all print statements

---

## Test Debugging

### Scenario: Debug Calculator Test

**File:** `testdata/calculator/calculator_test.go`

**Description:** Unit tests for calculator functions (Add, Subtract, Multiply, Divide).

### Debugging a Failing Test

**Step 1: Start test debug session**
```
mcp__delve-mcp__debug_test(
  testfile: "/mcp-go-debugger/testdata/calculator/calculator_test.go",
  testname: "TestAdd",
  testflags: ["-v"]
)

Response:
{
  "status": "success",
  "context": {
    "operation": "debug_test",
    "currentLocation": "At calculator_test.go:7 in TestAdd",
    "stopReason": "process is stopped"
  }
}
```

**Step 2: Set breakpoints**
```
# In test
mcp__delve-mcp__set_breakpoint(
  file: "/mcp-go-debugger/testdata/calculator/calculator_test.go",
  line: 8  # At first test case
)

# In implementation
mcp__delve-mcp__set_breakpoint(
  file: "/mcp-go-debugger/testdata/calculator/calculator.go",
  line: 4  # Inside Add function
)
```

**Step 3: Run test**
```
mcp__delve-mcp__continue()

Response:
{
  "status": "success",
  "context": {
    "currentLocation": "At calculator_test.go:8 in TestAdd",
    "localVariables": [
      {"name": "t", "value": "...", "type": "*testing.T"}
    ],
    "stopReason": "hit breakpoint"
  }
}
```

**Step 4: Inspect test data**
```
# Check test table
mcp__delve-mcp__eval_variable(name: "tests", depth: 2)
→ {
    value: "[]struct{...}{
      {a: 2, b: 3, expected: 5},
      {a: -1, b: 1, expected: 0},
      {a: 0, b: 0, expected: 0}
    }"
  }

# Step to loop iteration
mcp__delve-mcp__step_over()
mcp__delve-mcp__step_over()

# Check current test case
mcp__delve-mcp__eval_variable(name: "tt", depth: 1)
→ {value: "{a: 2, b: 3, expected: 5}"}
```

**Step 5: Step into Add function**
```
mcp__delve-mcp__continue()

Response:
{
  "status": "success",
  "context": {
    "currentLocation": "At calculator.go:4 in Add",
    "localVariables": [
      {"name": "a", "value": "2", "type": "int"},
      {"name": "b", "value": "3", "type": "int"}
    ],
    "stopReason": "hit breakpoint"
  }
}

# Verify inputs
mcp__delve-mcp__eval_variable(name: "a")
→ {value: "2", type: "int"}

mcp__delve-mcp__eval_variable(name: "b")
→ {value: "3", type: "int"}

# Execute addition
mcp__delve-mcp__step_over()

# Check result
mcp__delve-mcp__eval_variable(name: "a + b")
→ {value: "5", type: "int"}

mcp__delve-mcp__step_out()  # Return to test
```

**Step 6: Verify assertion**
```
# Back in test
mcp__delve-mcp__eval_variable(name: "result")
→ {value: "5", type: "int"}

mcp__delve-mcp__eval_variable(name: "tt.expected")
→ {value: "5", type: "int"}

# They match! Test should pass
mcp__delve-mcp__continue()
```

**Step 7: Complete test**
```
Response:
{
  "status": "success",
  "context": {
    "stopReason": "program exited"
  }
}

mcp__delve-mcp__close()
```

### Table-Driven Test Pattern

For table-driven tests, set conditional breakpoint on specific cases:

```
mcp__delve-mcp__set_breakpoint(
  file: "calculator_test.go",
  line: 10,
  condition: "tt.a == -1"  # Only break on negative number test
)
```

---

## Production Debugging

### Scenario: Attach to Running Server

**Context:** The webserver.go server is already running on PID 28026.

**Step 1: Find process**
```bash
ps aux | grep webserver
→ user 28026  0.5  1.2  webserver

# Or
pgrep -fl webserver
→ 28026 ./webserver
```

**Step 2: Attach to process**
```
mcp__delve-mcp__attach(pid: 28026)

Response:
{
  "status": "success",
  "context": {
    "operation": "attach",
    "currentLocation": "At runtime/select.go:590 in runtime.selectgo",
    "stopReason": "process is stopped"
  },
  "pid": 28026
}
```

**IMPORTANT:** Server is now paused! Users cannot access it.

**Step 3: Set conditional breakpoints (production-safe)**
```
# Only break for specific user
mcp__delve-mcp__set_breakpoint(
  file: "/mcp-go-debugger/testdata/webserver/webserver.go",
  line: 24,
  condition: "name == \"Alice\""
)

# Only break on errors (if there were any)
# set_breakpoint(file: "...", line: X, condition: "err != nil")
```

**Step 4: Resume server immediately**
```
mcp__delve-mcp__continue()

Response:
{
  "status": "success",
  "context": {
    "operation": "continue",
    "stopReason": "running"
  }
}
```

Server is now handling requests again!

**Step 5: Trigger breakpoint**
```bash
# Terminal 2
curl http://localhost:8080/?name=Alice
```

**Step 6: Inspect production state**
```
Response:
{
  "status": "success",
  "context": {
    "currentLocation": "At webserver.go:24 in main.helloHandler",
    "localVariables": [
      {"name": "name", "value": "Alice", "type": "string"},
      {"name": "requestCount", "value": "127", "type": "int"}
    ],
    "stopReason": "hit breakpoint (conditional: name == \"Alice\")"
  }
}

# Check production values
mcp__delve-mcp__eval_variable(name: "requestCount")
→ {value: "127"}  # Server has handled 127 requests

mcp__delve-mcp__eval_variable(name: "r.RemoteAddr")
→ {value: "192.168.1.50:55123"}  # Real client IP

mcp__delve-mcp__eval_variable(name: "r.Header", depth: 2)
→ {value: "map[string][]string{...}"}  # Real headers

# Capture findings quickly (< 30 seconds)
```

**Step 7: Resume immediately**
```
mcp__delve-mcp__continue()
```

**Step 8: Clean up**
```
# Remove breakpoint
mcp__delve-mcp__list_breakpoints()
→ [{id: 1, condition: "name == \"Alice\""}]

mcp__delve-mcp__remove_breakpoint(id: 1)

# Detach
mcp__delve-mcp__close()
```

Server continues running normally.

### Production Safety Rules

1. ✓ Use conditional breakpoints only
2. ✓ Resume within 30 seconds
3. ✓ Remove breakpoints when done
4. ✓ Always call close()
5. ✗ Never set unconditional breakpoints on hot paths
6. ✗ Never leave process paused

---

## Common Scenarios

### Scenario 1: Finding Nil Pointer Bug

**Problem:** Server panics with nil pointer dereference.

**Debugging Session:**
```
mcp__delve-mcp__debug(file: "/app/server.go")

# Set breakpoint before suspected line
mcp__delve-mcp__set_breakpoint(file: "/app/handler.go", line: 44)

mcp__delve-mcp__continue()
→ Trigger the panic scenario

# Check all pointers
mcp__delve-mcp__eval_variable(name: "user")
→ nil  # Found it!

mcp__delve-mcp__eval_variable(name: "session")
→ &{...}  # Not nil

mcp__delve-mcp__eval_variable(name: "config")
→ &{...}  # Not nil

# Found: user is nil, causes crash at line 45
```

**Fix:** Add nil check before line 45.

---

### Scenario 2: Wrong Calculation

**Problem:** Total is always 0 instead of sum of items.

**Debugging Session:**
```
mcp__delve-mcp__debug(file: "/app/calculator.go")

mcp__delve-mcp__set_breakpoint(file: "/app/calculator.go", line: 23)  # total = 0
mcp__delve-mcp__set_breakpoint(file: "/app/calculator.go", line: 25)  # loop

mcp__delve-mcp__continue()

# Check initial value
mcp__delve-mcp__eval_variable(name: "total")
→ 0

# Enter loop
mcp__delve-mcp__continue()

# Check items
mcp__delve-mcp__eval_variable(name: "items")
→ []  # Empty!

# Found: items slice is empty, that's why total is 0
```

**Fix:** Check why items is empty (not passed correctly, filtering issue, etc.).

---

### Scenario 3: Unexpected Branch

**Problem:** Code takes else branch when it should take if branch.

**Debugging Session:**
```
mcp__delve-mcp__debug(file: "/app/auth.go")

mcp__delve-mcp__set_breakpoint(file: "/app/auth.go", line: 45)  # if user.IsAdmin()

mcp__delve-mcp__continue()

# Check condition
mcp__delve-mcp__eval_variable(name: "user.IsAdmin()")
→ false  # Condition is false

mcp__delve-mcp__eval_variable(name: "user.Role")
→ "user"  # Not admin!

mcp__delve-mcp__eval_variable(name: "user", depth: 2)
→ {ID: 123, Name: "Alice", Role: "user", ...}

# Found: User role is "user", not "admin"
```

**Fix:** Check why user doesn't have admin role (database, test data, etc.).

---

### Scenario 4: Slow Performance

**Problem:** Handler is slow, need to find bottleneck.

**Debugging Session:**
```
mcp__delve-mcp__attach(pid: 28026)

mcp__delve-mcp__set_breakpoint(file: "/app/handler.go", line: 23)  # Entry
mcp__delve-mcp__set_breakpoint(file: "/app/handler.go", line: 45)  # DB query
mcp__delve-mcp__set_breakpoint(file: "/app/handler.go", line: 67)  # Processing
mcp__delve-mcp__set_breakpoint(file: "/app/handler.go", line: 89)  # Exit

mcp__delve-mcp__continue()

# Trigger request
# Note timestamps in responses:
→ 15:00:00.000 - Entry
→ 15:00:05.234 - After DB query (5 seconds!)
→ 15:00:05.345 - After processing (fast)
→ 15:00:05.456 - Exit

# Found: DB query takes 5 seconds
```

**Fix:** Optimize database query (add index, reduce data, etc.).

---

## Next Steps

- See `SKILL.md` for workflow guidance
- See `workflows.md` for detailed procedures
- See `patterns.md` for debugging patterns
- See `tools.md` for complete tool reference
