# Critical Console Bugs

## Bug 1: ESC Doesn't Work During Operations

**Problem:** When an operation is running (Build, Generate, etc.), pressing ESC does nothing.

**Current Code:**
```go
case "esc":
    if !a.asyncState.IsActive() {
        a.mode = ModeMenu
        // ...
    }
    return a, nil
```

ESC only works when `!asyncState.IsActive()`. If operation is running, it just returns without doing anything.

**Expected Behavior:** 
- ESC should ABORT the running operation
- Or at minimum, show "Operation aborted" and allow returning to menu

**Root Cause:** Operations run synchronously in a tea.Cmd and there's no process kill mechanism.

**Fix Required:**
1. Store the running command process
2. When ESC pressed during operation, kill the process
3. Set asyncState to inactive
4. Return to menu

---

## Bug 2: Console Flooded with Lines

**Problem:** Console shows 9k+ lines, making it unusable.

**Possible Causes:**
1. Legitimate build output from cmake (verbose mode)
2. Buffer not being cleared between operations
3. Streaming duplicating lines

**Current Buffer Code:**
```go
func (a *Application) startBuildOperation() (tea.Model, tea.Cmd) {
    a.mode = ModeConsole
    a.asyncState.Start()
    a.outputBuffer.Clear()  // <-- This clears the buffer
    // ...
}
```

Buffer IS being cleared. So 9k lines is likely legitimate cmake output.

**Fix Required:**
- Add line limit to OutputBuffer (circular buffer already has 1000 line limit)
- Or add filtering to reduce verbosity
- Or add collapsible sections

---

## Bug 3: Operation Never Completes

**Problem:** Operation runs forever, ESC doesn't work, console fills with lines.

**Root Cause:** The streaming goroutines might be hanging, or cmd.Wait() never returns.

**In ops/build.go:**
```go
// Stream stdout in goroutine
go func() {
    scanner := bufio.NewScanner(stdout)
    for scanner.Scan() {
        line := scanner.Text()
        if line != "" {
            outputCallback(line, ui.TypeStdout)
        }
    }
}()

// Wait for command to complete
err = cmd.Wait()
```

**Issues:**
1. Goroutines leak - they keep running even after cmd.Wait() returns
2. No timeout mechanism
3. No way to cancel/kill the process

**Fix Required:**
1. Add context with cancel
2. Kill process on abort
3. Properly close goroutines

---

## Summary of Required Fixes

### Immediate (Critical)
1. **Fix GetArrayIndex/GetVisibleIndex** - Add `IsSelectable` check (causes wrong operation to trigger)
2. **Add process kill on ESC** - Allow aborting operations
3. **Add timeout** - Prevent operations from running forever

### Short Term
4. **Fix compilation errors** - RegenerateCompleteMsg undefined, signature mismatches
5. **Add line limiting** - Prevent console flooding
6. **Fix goroutine leaks** - Proper cleanup

### Files to Fix
- `internal/app/app.go` - GetArrayIndex, GetVisibleIndex, ESC handler
- `internal/ops/build.go` - Add process kill, timeout, goroutine cleanup
- `internal/ops/setup.go` - Add process kill, timeout, goroutine cleanup
- `internal/app/op_*.go` - Fix compilation errors
- `internal/app/messages.go` - Add RegenerateCompleteMsg
