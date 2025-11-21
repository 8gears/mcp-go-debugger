# Common Debugging Patterns

This document catalogs common debugging patterns and scenarios with ready-to-use tool sequences.

## Table of Contents

1. [Finding Bugs](#finding-bugs)
2. [Understanding Code Flow](#understanding-code-flow)
3. [Inspecting Data Structures](#inspecting-data-structures)
4. [Error Investigation](#error-investigation)
5. [Performance Issues](#performance-issues)
6. [Concurrency Issues](#concurrency-issues)
7. [Testing Patterns](#testing-patterns)
8. [HTTP/API Debugging](#httpapi-debugging)

---

## Finding Bugs

### Pattern: Crash Location

**Scenario:** Program crashes with panic/fatal error. Find exact location.

**Tool Sequence:**
```
1. debug(file: "/app/main.go")
2. set_breakpoint at suspected crash location
3. continue()
4. → Crash occurs, inspect context.error and context.currentLocation
```

**Example:**
```
debug(file: "/app/server.go")

# Suspect it crashes when processing requests
set_breakpoint(file: "/app/handlers/api.go", line: 45)

continue()
→ Server starts

# Trigger crash
→ panic: runtime error: invalid memory address or nil pointer dereference
→ context.currentLocation: "At /app/handlers/api.go:47"
→ context.error: "panic: runtime error..."

eval_variable(name: "user")
→ nil  # Found it! User is nil
```

---

### Pattern: Nil Pointer Hunt

**Scenario:** Nil pointer dereference. Find which variable is nil.

**Tool Sequence:**
```
1. set_breakpoint before the crash line
2. continue()
3. eval_variable for each pointer variable
4. Find the nil one
```

**Example:**
```
set_breakpoint(file: "/app/handler.go", line: 44)  # Before crash at 45

continue()
→ Stopped at handler.go:44

# Check all pointer variables
eval_variable(name: "request")
→ &{...}  # Not nil

eval_variable(name: "user")
→ nil  # This is the problem!

eval_variable(name: "session")
→ &{...}  # Not nil

# Found: user is nil at line 44, causes crash at line 45
```

---

### Pattern: Wrong Value

**Scenario:** Variable has unexpected value. Find where it was set wrong.

**Tool Sequence:**
```
1. set_breakpoint where variable is assigned
2. continue()
3. eval_variable before and after assignment
4. step() through assignment logic
```

**Example:**
```
# Expected: count = 10, Actual: count = 0

set_breakpoint(file: "/app/calculator.go", line: 23)  # count = ...

continue()
→ Stopped before assignment

eval_variable(name: "count")
→ 5  # Current value

step()  # Execute assignment

eval_variable(name: "count")
→ 0  # Wrong! Should be 10

# Look at the assignment
eval_variable(name: "items")
→ []  # Empty slice!

# Found: count = len(items), but items is empty
```

---

### Pattern: Logic Error

**Scenario:** Code executes but produces wrong result. Trace logic flow.

**Tool Sequence:**
```
1. set_breakpoint at function entry
2. set_breakpoint at each branch point
3. step_over() through logic
4. eval_variable at each step
```

**Example:**
```
# Function should return true for even numbers, but returns false

set_breakpoint(file: "/app/math.go", line: 12)  # func IsEven(n int)

continue()
→ Stopped at IsEven

eval_variable(name: "n")
→ 4  # Testing with 4 (even)

step_over()  # Execute: result := n % 2 == 1

eval_variable(name: "result")
→ false  # Wrong! Should be true for 4

# Found: Logic is `n % 2 == 1` (checks for odd)
# Should be: `n % 2 == 0` (checks for even)
```

---

## Understanding Code Flow

### Pattern: Trace Execution Path

**Scenario:** Don't understand how code reaches a certain point.

**Tool Sequence:**
```
1. set_breakpoint at the mysterious location
2. continue()
3. Use step_out() to see caller
4. Repeat to trace call stack
```

**Example:**
```
set_breakpoint(file: "/app/service.go", line: 89)

continue()
→ Stopped at service.go:89
→ How did we get here?

step_out()
→ Returned to handler.go:45

step_out()
→ Returned to router.go:23

step_out()
→ Returned to main.go:67

# Call stack: main → router → handler → service
```

---

### Pattern: Follow Function Calls

**Scenario:** Want to see exactly what a function does.

**Tool Sequence:**
```
1. set_breakpoint at function call
2. step() to enter function
3. step_over() through function body
4. eval_variable at interesting points
```

**Example:**
```
set_breakpoint(file: "/app/handler.go", line: 34)  # Before CreateUser()

continue()
→ Stopped at handler.go:34

eval_variable(name: "username")
→ "alice"

step()  # Enter CreateUser
→ Now at service.go:12 (inside CreateUser)

step_over()  # Validate username
eval_variable(name: "valid")
→ true

step_over()  # Hash password
eval_variable(name: "hashedPassword")
→ "$2a$10$..."

step_over()  # Save to database
eval_variable(name: "err")
→ nil  # Success
```

---

### Pattern: Conditional Path Discovery

**Scenario:** Code has multiple branches. Want to see which executes.

**Tool Sequence:**
```
1. set_breakpoint at if/switch statement
2. eval_variable for condition
3. step_over() to see which branch taken
```

**Example:**
```
# Code: if user.IsAdmin() { ... } else { ... }

set_breakpoint(file: "/app/auth.go", line: 45)  # At if statement

continue()
→ Stopped at auth.go:45

eval_variable(name: "user.IsAdmin()")
→ false  # Condition is false

step_over()
→ Jumped to line 52  # Else branch

# Discovered: Code takes else branch (non-admin path)
```

---

## Inspecting Data Structures

### Pattern: Map Contents

**Scenario:** Need to see what's in a map.

**Tool Sequence:**
```
1. set_breakpoint where map is used
2. eval_variable(name: "mapVar", depth: 2)
3. Check specific keys if needed
```

**Example:**
```
set_breakpoint(file: "/app/cache.go", line: 56)

continue()

eval_variable(name: "cache", depth: 2)
→ map[string]interface{}{
    "user:123": {Name: "Alice", Email: "..."},
    "user:456": {Name: "Bob", Email: "..."},
    "session:abc": {...}
  }

# Check specific key
eval_variable(name: "cache[\"user:123\"]", depth: 3)
→ {Name: "Alice", Email: "alice@example.com", Role: "admin"}
```

---

### Pattern: Slice Inspection

**Scenario:** Need to verify slice contents and length.

**Tool Sequence:**
```
1. eval_variable(name: "slice", depth: 1)
2. Check len/cap
3. Inspect specific indices
```

**Example:**
```
eval_variable(name: "users")
→ []User{
    {ID: 1, Name: "Alice"},
    {ID: 2, Name: "Bob"},
    {ID: 3, Name: "Charlie"}
  }

eval_variable(name: "len(users)")
→ 3

eval_variable(name: "cap(users)")
→ 10

eval_variable(name: "users[0]", depth: 2)
→ {ID: 1, Name: "Alice", Email: "alice@example.com"}
```

---

### Pattern: Struct Fields

**Scenario:** Need to see all fields of a complex struct.

**Tool Sequence:**
```
1. eval_variable(name: "structVar", depth: 2-3)
2. Inspect nested structs if needed
```

**Example:**
```
eval_variable(name: "request", depth: 2)
→ http.Request{
    Method: "POST",
    URL: &url.URL{
      Scheme: "http",
      Host: "localhost:8080",
      Path: "/api/users"
    },
    Header: {
      "Content-Type": ["application/json"],
      "Authorization": ["Bearer ..."]
    },
    Body: {...}
  }

# Inspect nested URL
eval_variable(name: "request.URL", depth: 1)
→ &url.URL{Scheme: "http", Host: "localhost:8080", Path: "/api/users"}
```

---

### Pattern: Interface Value

**Scenario:** Variable is interface type, need to see concrete value.

**Tool Sequence:**
```
1. eval_variable(name: "interfaceVar", depth: 2)
2. Check type and value
```

**Example:**
```
# Variable: var response interface{}

eval_variable(name: "response", depth: 2)
→ (*UserResponse){
    User: {ID: 123, Name: "Alice"},
    Token: "abc...",
    ExpiresAt: "2025-01-01T00:00:00Z"
  }

# Can see it's actually *UserResponse with these values
```

---

## Error Investigation

### Pattern: Error Origin

**Scenario:** Error occurs, need to find where it originated.

**Tool Sequence:**
```
1. set_breakpoint(condition: "err != nil")
2. continue()
3. step_out() to find caller
4. Trace back to origin
```

**Example:**
```
set_breakpoint(
  file: "/app/handler.go",
  line: 67,
  condition: "err != nil"
)

continue()
→ Stopped at handler.go:67

eval_variable(name: "err")
→ "database connection failed"

step_out()
→ At service.go:45

eval_variable(name: "err")
→ Same error

step_out()
→ At db.go:23

eval_variable(name: "err")
→ "connection refused"  # Origin!
```

---

### Pattern: Error Context

**Scenario:** Error message lacks context. Need to see what was happening.

**Tool Sequence:**
```
1. set_breakpoint where error is returned
2. eval_variable for error and surrounding context
```

**Example:**
```
set_breakpoint(file: "/app/service.go", line: 89, condition: "err != nil")

continue()
→ Stopped

eval_variable(name: "err")
→ "validation failed"  # Not specific

eval_variable(name: "user")
→ {Name: "", Email: "test@example.com"}  # Name is empty!

eval_variable(name: "validationErrors")
→ ["name is required"]

# Now we know: validation failed because name is empty
```

---

### Pattern: Panic Recovery

**Scenario:** Code panics, need to see state at panic point.

**Tool Sequence:**
```
1. debug(file: "...")
2. continue() until panic
3. Check context.error and context.localVariables
```

**Example:**
```
debug(file: "/app/main.go")

continue()
→ panic: runtime error: index out of range [5] with length 3

→ context.currentLocation: "At /app/processor.go:34"
→ context.localVariables: [
    {name: "items", value: "[1, 2, 3]"},
    {name: "index", value: "5"}
  ]

# Found: Trying to access items[5] but items only has 3 elements
```

---

## Performance Issues

### Pattern: Slow Function

**Scenario:** Function is slow. Find bottleneck.

**Tool Sequence:**
```
1. set_breakpoint at function entry and exit
2. step_over() through function
3. Note which operations take time
```

**Example:**
```
set_breakpoint(file: "/app/service.go", line: 23)  # Entry
set_breakpoint(file: "/app/service.go", line: 78)  # Exit

continue()
→ Stopped at entry (14:30:00.000)

step_over()  # Validate input
→ Fast (14:30:00.001)

step_over()  # Query database
→ Slow! (14:30:05.000)  # 5 seconds!

step_over()  # Process results
→ Fast (14:30:05.100)

# Found: Database query is slow (5 seconds)
```

---

### Pattern: Loop Performance

**Scenario:** Loop is slow. Check iteration efficiency.

**Tool Sequence:**
```
1. set_breakpoint inside loop
2. continue() through several iterations
3. eval_variable to check iteration variables
```

**Example:**
```
set_breakpoint(file: "/app/processor.go", line: 45)  # Inside loop

continue()
→ Iteration 1
eval_variable(name: "i")
→ 0

continue()
→ Iteration 2
eval_variable(name: "i")
→ 1

# If iterations are slow, inspect what happens each time
eval_variable(name: "item")
eval_variable(name: "result")

# Check for inefficiencies:
# - Network calls in loop?
# - Database queries in loop?
# - Expensive computations?
```

---

## Concurrency Issues

### Pattern: Channel Deadlock

**Scenario:** Goroutines deadlock on channels. Find cause.

**Tool Sequence:**
```
1. attach(pid: <hanging_process>)
2. Check context.currentLocation (where it's stuck)
3. eval_variable for channel state
```

**Example:**
```
attach(pid: 28026)
→ Process attached

→ context.currentLocation: "At /app/worker.go:67 in worker.process"
→ Stuck at channel send

eval_variable(name: "workCh")
→ {len: 100, cap: 100}  # Channel is full!

eval_variable(name: "resultCh")
→ {len: 0, cap: 0}  # No one receiving

# Found: workCh is full, send blocks
#        resultCh has no receiver, goroutine stuck
```

---

### Pattern: Race Condition

**Scenario:** Data corruption from concurrent access. Find where.

**Tool Sequence:**
```
1. set_breakpoint at shared data access points
2. Check lock state before access
3. Verify synchronization
```

**Example:**
```
# Shared data: var counter int

set_breakpoint(file: "/app/worker.go", line: 34)  # counter++

continue()
→ Stopped

eval_variable(name: "counter")
→ 42

step()

eval_variable(name: "counter")
→ 43

# But if another goroutine also increments without sync:
# Expected: 43 → 44
# Actual: 43 → 43 (race condition overwrites)

# Look for missing mutex:
eval_variable(name: "mu")
→ Not found  # No mutex! That's the bug
```

---

### Pattern: Goroutine Leak

**Scenario:** Too many goroutines. Find what's not exiting.

**Tool Sequence:**
```
1. set_breakpoint where goroutines are created
2. set_breakpoint where goroutines should exit
3. Check if exit breakpoint is hit
```

**Example:**
```
set_breakpoint(file: "/app/worker.go", line: 23)  # go worker()
set_breakpoint(file: "/app/worker.go", line: 67)  # return (exit)

continue()
→ Stopped at line 23 (goroutine created)

eval_variable(name: "workerCount")
→ 100

continue()
→ Stopped at line 23 again (another goroutine created)

eval_variable(name: "workerCount")
→ 101

# Never hits line 67 (exit)!
# Found: Goroutines are created but never exit
```

---

## Testing Patterns

### Pattern: Test Setup Verification

**Scenario:** Test fails, suspect test setup is wrong.

**Tool Sequence:**
```
1. debug_test(...)
2. set_breakpoint after setup, before test logic
3. eval_variable for test fixtures
```

**Example:**
```
debug_test(testfile: "/app/handler_test.go", testname: "TestCreateUser")

set_breakpoint(file: "/app/handler_test.go", line: 34)  # After setup

continue()
→ Stopped

# Check test fixtures
eval_variable(name: "testDB", depth: 1)
→ nil  # Not initialized!

eval_variable(name: "testServer")
→ &Server{...}  # Initialized

# Found: testDB is nil, setup incomplete
```

---

### Pattern: Mock Verification

**Scenario:** Test fails, verify mock behavior.

**Tool Sequence:**
```
1. set_breakpoint where mock is called
2. eval_variable for mock state
3. Verify expected behavior
```

**Example:**
```
set_breakpoint(file: "/app/service.go", line: 45)  # db.Create()

continue()

# Check if using mock
eval_variable(name: "db")
→ *MockDB{...}  # Yes, using mock

step()  # Execute db.Create()

eval_variable(name: "db.CreateCalled")
→ true  # Mock was called

eval_variable(name: "db.CreateInput")
→ {Name: "Alice", Email: "..."}  # Correct input

eval_variable(name: "db.CreateReturn")
→ nil  # Returns nil error

# Mock is working correctly
```

---

### Pattern: Assertion Failure

**Scenario:** Assertion fails, need to see actual vs expected.

**Tool Sequence:**
```
1. set_breakpoint at assertion line
2. eval_variable for both values
3. Compare to understand difference
```

**Example:**
```
# Test: assert.Equal(t, expected, actual)

set_breakpoint(file: "/app/handler_test.go", line: 67)

continue()

eval_variable(name: "expected")
→ {Status: 201, Body: "{\"id\":123}"}

eval_variable(name: "actual")
→ {Status: 500, Body: "{\"error\":\"validation failed\"}"}

# Found: Expected success (201) but got error (500)
#        Body shows validation failed
```

---

## HTTP/API Debugging

### Pattern: Request Inspection

**Scenario:** API behaves unexpectedly. Inspect incoming request.

**Tool Sequence:**
```
1. set_breakpoint at handler entry
2. eval_variable(name: "request", depth: 2)
3. Check method, headers, body
```

**Example:**
```
set_breakpoint(file: "/app/handlers/api.go", line: 23)

continue()
→ Request received

eval_variable(name: "r", depth: 2)
→ &http.Request{
    Method: "POST",
    URL: &url.URL{Path: "/api/users"},
    Header: {
      "Content-Type": ["application/json"],
      "Authorization": []  # Empty!
    },
    Body: {...}
  }

# Found: Authorization header is missing
```

---

### Pattern: Response Verification

**Scenario:** API returns wrong response. Check what's being sent.

**Tool Sequence:**
```
1. set_breakpoint before response is sent
2. eval_variable for response data
3. Verify status, headers, body
```

**Example:**
```
set_breakpoint(file: "/app/handlers/api.go", line: 89)  # Before w.Write()

continue()

eval_variable(name: "status")
→ 500

eval_variable(name: "responseBody")
→ "{\"error\":\"internal server error\"}"

eval_variable(name: "err")
→ "database connection failed"

# Found: Returning 500 because database connection failed
```

---

### Pattern: Middleware Chain

**Scenario:** Request processed by multiple middleware. Trace flow.

**Tool Sequence:**
```
1. set_breakpoint in each middleware
2. continue() through chain
3. Check modifications at each step
```

**Example:**
```
set_breakpoint(file: "/app/middleware/auth.go", line: 12)
set_breakpoint(file: "/app/middleware/logging.go", line: 23)
set_breakpoint(file: "/app/handlers/api.go", line: 34)

continue()
→ Stopped at auth.go:12

eval_variable(name: "r.Header[\"Authorization\"]")
→ ["Bearer abc123"]

continue()
→ Stopped at logging.go:23

eval_variable(name: "r.Context()")
→ {..., userID: "123"}  # Auth middleware added userID

continue()
→ Stopped at api.go:34

eval_variable(name: "r.Context().Value(\"userID\")")
→ "123"  # userID available in handler

# Traced: auth → logging → handler
```

---

## Pattern Selection Guide

| Problem | Pattern |
|---------|---------|
| Program crashes | Crash Location |
| Nil pointer error | Nil Pointer Hunt |
| Variable has wrong value | Wrong Value |
| Wrong calculation result | Logic Error |
| Don't understand execution path | Trace Execution Path |
| Want to see function internals | Follow Function Calls |
| Need to see which branch executes | Conditional Path Discovery |
| Check map contents | Map Contents |
| Verify slice data | Slice Inspection |
| Inspect complex struct | Struct Fields |
| See interface concrete value | Interface Value |
| Find error origin | Error Origin |
| Understand error context | Error Context |
| Recover from panic | Panic Recovery |
| Slow function | Slow Function |
| Slow loop | Loop Performance |
| Channel deadlock | Channel Deadlock |
| Data race | Race Condition |
| Too many goroutines | Goroutine Leak |
| Test setup issue | Test Setup Verification |
| Mock not working | Mock Verification |
| Assertion fails | Assertion Failure |
| API request issue | Request Inspection |
| API response wrong | Response Verification |
| Middleware issue | Middleware Chain |

---

## Next Steps

- See `workflows.md` for complete debugging workflows
- See `tools.md` for tool reference
- See `examples.md` for real-world scenarios
