# Fundamental Guide to Software Debugging

## Core Philosophy

**Debugging is a scientific process**: Form hypothesis → Test → Observe → Refine

Never randomly change code hoping it works. Use the debugger to **understand** before fixing.

## The Debugging Mindset

### 1. Reproduce First
- **No reproduction = No real debugging**
- Minimize the reproduction case
- Document exact steps
- Save input data that triggers the bug

### 2. Binary Search Strategy
Cut the problem space in half repeatedly:
```
Working state ←────[???]────→ Broken state
                     ↓
              Find midpoint
```

### 3. Question Assumptions
The bug is usually in the code you're **certain** works correctly.

## Essential Debugger Concepts

### Breakpoints - Your Primary Tool

**Types:**
- **Line breakpoints**: Stop at specific code line
- **Conditional breakpoints**: Stop only when condition true
- **Watchpoints**: Stop when data changes (via conditional breakpoints)
- **Function breakpoints**: Stop at function entry
- **Panic breakpoints**: Stop on panic (automatic in delve-mcp)

**Strategy:**
```
Start wide → Narrow down → Pinpoint exact location
```

### Execution Control

```
┌─────────┐
│Continue │ ──→ Execute until breakpoint
└─────────┘
     ↓
┌─────────┐
│StepOver│ ──→ One line forward
│   (n)   │     (don't enter functions)
└─────────┘
     ↓
┌─────────┐
│  Step   │ ──→ One line forward
│   (s)   │     (enter functions)
└─────────┘
     ↓
┌─────────┐
│ StepOut │ ──→ Run until current function returns
│  (so)   │
└─────────┘
```

## The Systematic Approach

### Phase 1: Observe the Failure

```go
// Code with bug
func ProcessData(data []byte) error {
    result := transform(data)  // Something goes wrong here
    return save(result)
}
```

**MCP Debug Session:**
```
# Start debugging
mcp__delve-mcp__debug(file: "/path/to/main.go")

# Set breakpoint at ProcessData
mcp__delve-mcp__set_breakpoint(file: "/path/to/main.go", line: 23)

# Continue to breakpoint
mcp__delve-mcp__continue()

# When hit, examine context (automatic in response):
→ context.localVariables shows current state
→ context.currentLocation shows stack position

# Inspect suspicious data
mcp__delve-mcp__eval_variable(name: "data")
mcp__delve-mcp__eval_variable(name: "len(data)")
```

**Document what you see** - Don't trust memory.

### Phase 2: Trace Backwards

Work backwards from failure to cause:

```go
func FailingFunction(input *Config) error {
    // This returns error
    validated := validate(input)
    if validated == nil {
        return errors.New("validation failed")
    }
    // ...
}
```

**MCP Debug Session:**
```
# Start from known bad state
mcp__delve-mcp__set_breakpoint(file: "/path/to/service.go", line: 45)
mcp__delve-mcp__continue()

# Step backwards mentally:
# "For this to fail, validated must be nil"
# "For validated to be nil, validate must return nil"

# Set breakpoint at validate function
mcp__delve-mcp__set_breakpoint(file: "/path/to/validator.go", line: 12)

# Restart and trace
mcp__delve-mcp__close()
mcp__delve-mcp__debug(file: "/path/to/main.go")
mcp__delve-mcp__continue()
```

### Phase 3: Binary Search for Root Cause

```go
func ProcessList(items []Item) []Result {
    // Works at start, corrupted at end
    results := make([]Result, 0)

    for i, item := range items {  // Where does it break?
        processed := processItem(item)
        results = append(results, processed)
    }

    return results
}
```

**MCP Debug Session:**
```
# Set breakpoint at loop entry
mcp__delve-mcp__set_breakpoint(file: "/path/to/processor.go", line: 23)
mcp__delve-mcp__continue()

# Check total items
mcp__delve-mcp__eval_variable(name: "len(items)")
→ 100 items

# Binary search: Set conditional breakpoint at midpoint
mcp__delve-mcp__set_breakpoint(
  file: "/path/to/processor.go",
  line: 25,
  condition: "i == 50"
)

mcp__delve-mcp__continue()

# Check if results are corrupted
mcp__delve-mcp__eval_variable(name: "results", depth: 1)

# If corrupted: problem is in first half (0-50)
# If good: problem is in second half (50-100)
# Remove old breakpoint, set new midpoint
mcp__delve-mcp__remove_breakpoint(id: 2)
mcp__delve-mcp__set_breakpoint(
  file: "/path/to/processor.go",
  line: 25,
  condition: "i == 25"  # or i == 75
)
# Repeat until found
```

## Data Inspection Techniques

### 1. Variable Evolution

Track how data changes:

```go
type State struct {
    Counter int
    Buffer  []byte
    Config  *Config
}

func (s *State) Process() {
    s.Counter++  // Track this
    s.Buffer = append(s.Buffer, 0xFF)
    // ...
}
```

**MCP Debug Session:**
```
# Track variable changes with conditional breakpoint
mcp__delve-mcp__set_breakpoint(
  file: "/path/to/state.go",
  line: 12,  # After s.Counter++
  condition: "s.Counter > 100"  # Only when suspicious
)

mcp__delve-mcp__continue()

# Check value each time
mcp__delve-mcp__eval_variable(name: "s.Counter")
mcp__delve-mcp__eval_variable(name: "len(s.Buffer)")
```

### 2. Memory & Pointer Inspection

```go
type Node struct {
    Value int
    Next  *Node
}

func Traverse(head *Node) {
    current := head
    for current != nil {
        process(current.Value)
        current = current.Next
    }
}
```

**MCP Debug Session:**
```
# Pointer examination
mcp__delve-mcp__set_breakpoint(file: "/path/to/list.go", line: 23)
mcp__delve-mcp__continue()

# Inspect pointer chain
mcp__delve-mcp__eval_variable(name: "head", depth: 1)
mcp__delve-mcp__eval_variable(name: "head.Next", depth: 1)
mcp__delve-mcp__eval_variable(name: "head.Next.Value")

# Follow pointer chain in loop
mcp__delve-mcp__eval_variable(name: "current")
mcp__delve-mcp__eval_variable(name: "current.Value")
mcp__delve-mcp__eval_variable(name: "current.Next")

# Check for cycles or nil
mcp__delve-mcp__eval_variable(name: "current.Next == nil")
```

### 3. State Validation

```go
func Process(data []byte) ([]byte, error) {
    // Debug assertions
    if len(data) == 0 {
        panic("Precondition failed: empty data")
    }

    result := transform(data)

    if len(result) > len(data)*2 {
        panic("Postcondition failed: result too large")
    }

    return result, nil
}
```

**MCP Debug Session:**
```
# Set breakpoints at assertions
mcp__delve-mcp__set_breakpoint(
  file: "/path/to/processor.go",
  line: 3,
  condition: "len(data) == 0"  # Catch precondition violation
)

mcp__delve-mcp__set_breakpoint(
  file: "/path/to/processor.go",
  line: 8,
  condition: "len(result) > len(data)*2"  # Catch postcondition violation
)

mcp__delve-mcp__continue()
```

## Advanced Debugging Patterns

### Pattern 1: Differential Debugging

Compare working vs. broken:

```go
func HandleRequest(req Request) Response {
    // Sometimes works, sometimes fails
    validated := validate(req)
    processed := process(validated)
    return format(processed)
}
```

**MCP Debug Session:**
```
# Run with working input
mcp__delve-mcp__debug(file: "/path/to/main.go", args: ["--test-mode", "good"])
mcp__delve-mcp__set_breakpoint(file: "/path/to/handler.go", line: 23)
mcp__delve-mcp__continue()

# Save the state
mcp__delve-mcp__eval_variable(name: "req", depth: 2)
# → Document output

mcp__delve-mcp__close()

# Run with broken input
mcp__delve-mcp__debug(file: "/path/to/main.go", args: ["--test-mode", "bad"])
mcp__delve-mcp__set_breakpoint(file: "/path/to/handler.go", line: 23)
mcp__delve-mcp__continue()

# Compare state
mcp__delve-mcp__eval_variable(name: "req", depth: 2)
# → Compare with previous output to find difference
```

### Pattern 2: Goroutine Debugging

```go
func ConcurrentBug() {
    var wg sync.WaitGroup
    shared := make([]int, 0)

    for i := 0; i < 10; i++ {
        wg.Add(1)
        go func(n int) {
            defer wg.Done()
            shared = append(shared, n)  // Race condition
        }(i)
    }

    wg.Wait()
}
```

**MCP Debug Session:**
```
# Track goroutines - check context.localVariables in responses
mcp__delve-mcp__set_breakpoint(file: "/path/to/concurrent.go", line: 8)
mcp__delve-mcp__continue()

# Each time breakpoint hits, it's a different goroutine
# Check shared state
mcp__delve-mcp__eval_variable(name: "shared", depth: 1)
mcp__delve-mcp__eval_variable(name: "n")

# Set condition to catch race
mcp__delve-mcp__set_breakpoint(
  file: "/path/to/concurrent.go",
  line: 8,
  condition: "len(shared) > 5"  # Catch when race likely
)
```

**Note:** For race conditions, prefer `go run -race` for detection, then use debugger to understand.

### Pattern 3: Hypothesis Testing

```go
func BuggyFunction(index int, data []byte) {
    // Hypothesis: Crashes when index >= len(data)

    // Process...
    result := data[index]
}
```

**MCP Debug Session:**
```
# Test hypothesis with conditional breakpoint
mcp__delve-mcp__set_breakpoint(
  file: "/path/to/buggy.go",
  line: 5,
  condition: "index >= len(data)"  # Test exact hypothesis
)

mcp__delve-mcp__continue()

# Does it trigger breakpoint?
# If yes: Hypothesis confirmed - index out of bounds
# If no: Revise hypothesis, test different condition
```

## Dealing with Specific Bug Types

### Nil Pointer Dereference

```go
type Service struct {
    client *Client
}

func (s *Service) DoWork() {
    s.client.Send()  // Potential nil panic
}
```

**MCP Debug Session:**
```
# Catch before panic
mcp__delve-mcp__set_breakpoint(
  file: "/path/to/service.go",
  line: 7,
  condition: "s.client == nil"  # Only break if nil
)

mcp__delve-mcp__continue()

# If breakpoint hits, found the problem
mcp__delve-mcp__eval_variable(name: "s", depth: 2)
mcp__delve-mcp__eval_variable(name: "s.client")
→ nil  # Confirmed
```

### Slice Bounds Errors

```go
func GetElement(items []string, index int) string {
    return items[index]  // Potential panic
}
```

**MCP Debug Session:**
```
# Debug bounds issue
mcp__delve-mcp__set_breakpoint(
  file: "/path/to/getter.go",
  line: 2,
  condition: "index >= len(items)"  # Catch out of bounds
)

mcp__delve-mcp__continue()

# When hit, inspect
mcp__delve-mcp__eval_variable(name: "len(items)")
mcp__delve-mcp__eval_variable(name: "cap(items)")
mcp__delve-mcp__eval_variable(name: "index")
→ Found: index=10, len=5
```

### Race Conditions

```go
type Counter struct {
    mu    sync.Mutex
    value int
}

func (c *Counter) Inc() {
    // c.mu.Lock()  // Forgotten!
    c.value++
    // c.mu.Unlock()
}
```

**MCP Debug Session:**
```
# First, build with race detector
# go build -race main.go

# In debugger, track mutex state
mcp__delve-mcp__set_breakpoint(file: "/path/to/counter.go", line: 9)
mcp__delve-mcp__continue()

# Check mutex state (not locked!)
mcp__delve-mcp__eval_variable(name: "c.mu")
mcp__delve-mcp__eval_variable(name: "c.value")

# Each continue() might be different goroutine
mcp__delve-mcp__continue()
mcp__delve-mcp__eval_variable(name: "c.value")
→ Value changed unexpectedly (race detected)
```

### Channel Deadlocks

```go
func DeadlockExample() {
    ch := make(chan int)  // Unbuffered

    ch <- 42  // Blocks forever
    val := <-ch
    fmt.Println(val)
}
```

**MCP Debug Session:**
```
# Attach to hanging process
mcp__delve-mcp__attach(pid: 28026)

# Check where it's stuck
→ context.currentLocation: "At deadlock.go:4 in DeadlockExample"
# Stuck at channel send

# Inspect channel state
mcp__delve-mcp__eval_variable(name: "ch")
→ Channel state (cap=0, no receiver)

# Identified: Unbuffered channel with no receiver
```

## The Debugging Workflow

### 1. Preparation

```bash
# Build with debug symbols, no optimization
go build -gcflags="all=-N -l" -o myapp

# Enable race detection for concurrency bugs
go build -race -o myapp

# Set up verbose logging
export LOG_LEVEL=debug
```

### 2. Initial Investigation

**MCP Debug Session:**
```
# Start debugging
mcp__delve-mcp__debug(file: "/path/to/main.go", args: ["--config", "config.yaml"])

# Set up wide net - panic breakpoints are automatic in delve-mcp

# Run
mcp__delve-mcp__continue()
```

### 3. Narrowing Down

```go
func ComplexFunction(input Data) Result {
    step1 := validate(input)     // OK?
    step2 := transform(step1)    // OK?
    step3 := optimize(step2)     // OK?
    return finalize(step3)       // Fails here
}
```

**MCP Debug Session:**
```
# Found general area, now be surgical
mcp__delve-mcp__list_breakpoints()
mcp__delve-mcp__remove_breakpoint(id: 1)  # Clear old breakpoints

# Surgical breakpoints with condition
mcp__delve-mcp__set_breakpoint(
  file: "/path/to/complex.go",
  line: 4,
  condition: "step2.IsValid() == false"  # Only when problem occurs
)

mcp__delve-mcp__continue()
```

### 4. Root Cause Analysis

**MCP Debug Session:**
```
# At the bug location
→ context.currentLocation shows exactly where we are
→ context.localVariables shows all locals automatically

# Inspect specific variables deeply
mcp__delve-mcp__eval_variable(name: "step3", depth: 3)
mcp__delve-mcp__eval_variable(name: "input", depth: 2)

# Trace back through execution
mcp__delve-mcp__step_out()  # Go to caller
mcp__delve-mcp__eval_variable(name: "result")  # What was returned?
```

## Common Go Debugging Mistakes

### 1. **Debugging Optimized Code**
```bash
# Wrong - optimizations break debugging
go build -o app

# Right - disable optimizations
go build -gcflags="all=-N -l" -o app
```

### 2. **Ignoring Goroutine Leaks**
```
# When attaching to process, check context for goroutine info
mcp__delve-mcp__attach(pid: 28026)
→ Look at goroutine counts in responses
# If count keeps growing, you have a leak
```

### 3. **Not Using Race Detector**
```bash
# Always test concurrent code with race detector first
go run -race main.go
# Then use debugger to understand the race
```

### 4. **Forgetting Defer Execution**
```go
func Example() {
    defer cleanup()  // Runs even during debugging step-out
    // ...
}
```
**Note:** `step_out()` will execute defers before returning!

## Debugger Power Features

### Conditional Breakpoints
```
# Complex conditions
mcp__delve-mcp__set_breakpoint(
  file: "/path/to/handler.go",
  line: 45,
  condition: "r.Method == \"POST\" && len(r.Body) > 1024"
)

mcp__delve-mcp__set_breakpoint(
  file: "/path/to/processor.go",
  line: 23,
  condition: "item.Priority > 5 && item.Status == \"pending\""
)
```

### Multiple Hypothesis Testing
```
# Set multiple conditional breakpoints to test different theories
mcp__delve-mcp__set_breakpoint(
  file: "/path/to/service.go",
  line: 45,
  condition: "user == nil"  # Hypothesis 1
)

mcp__delve-mcp__set_breakpoint(
  file: "/path/to/service.go",
  line: 56,
  condition: "len(data) == 0"  # Hypothesis 2
)

mcp__delve-mcp__set_breakpoint(
  file: "/path/to/service.go",
  line: 67,
  condition: "err != nil"  # Hypothesis 3
)

mcp__delve-mcp__continue()
# Whichever breaks first tells you which hypothesis is correct
```

## Logging vs Debugging

**Use logging when:**
- Production issues
- Distributed systems
- Need audit trail
- Performance metrics

**Use debugger when:**
- Complex state inspection
- Need to modify execution (step through)
- Goroutine issues
- Hypothesis testing

## The Professional Approach

### 1. Keep a Debug Log
```markdown
Bug #427: Panic in handler
- Occurs with concurrent requests > 100
- context.localVariables shows nil pointer at handler.go:45
- Traced back with step_out() to shared client not initialized
- Root cause: Race condition in lazy initialization
- Fix: Add sync.Once for initialization
```

### 2. Use Conditional Breakpoints as Debug Helpers
```
# Instead of adding code, use conditional breakpoints
mcp__delve-mcp__set_breakpoint(
  file: "/path/to/service.go",
  line: 45,
  condition: "counter > 1000"  # Debug checkpoint
)
```

### 3. Defensive Programming
```go
// Development - catch early with conditional breakpoints
// Production - return errors gracefully
if ptr == nil {
    return fmt.Errorf("invalid state: nil pointer")
}
```

## Mental Model

Think of debugging as:
1. **Detective work**: Gather clues, form theories
2. **Scientific method**: Hypothesis → Experiment → Conclusion
3. **Binary search**: Systematically eliminate possibilities

## Key Principles

1. **Trust nothing** - Verify every assumption
2. **Change one thing** - Isolate variables (test one hypothesis at a time)
3. **Simplify** - Minimal test case
4. **Document** - Record what you find (future you will thank you)
5. **Learn** - Each bug teaches something about the system

## Quick Reference: MCP Debug Commands

```
# Session Management
Start:          debug(file: "/path/to/main.go")
Test:           debug_test(testfile: "...", testname: "TestFoo")
Attach:         attach(pid: 28026)
Close:          close()

# Breakpoints
Set:            set_breakpoint(file: "...", line: 42)
Conditional:    set_breakpoint(file: "...", line: 42, condition: "x > 10")
List:           list_breakpoints()
Remove:         remove_breakpoint(id: 1)

# Execution Control
Continue:       continue()
Step Over:      step_over()
Step Into:      step()
Step Out:       step_out()

# Inspection
Variable:       eval_variable(name: "varName")
Deep:           eval_variable(name: "varName", depth: 3)
Output:         get_debugger_output()

# Context (Automatic)
Every response includes:
- context.currentLocation (where execution stopped)
- context.localVariables (all local variables)
- context.stopReason (why it stopped)
```

## The Golden Rule

**The debugger shows you what IS happening, not what you THINK is happening.**

Trust the debugger over your assumptions.

Always verify your mental model against actual execution state.
