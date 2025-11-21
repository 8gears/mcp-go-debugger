package debugger

import (
	"fmt"
	"github.com/go-delve/delve/service/api"
)

// getFunctionName extracts a human-readable function name from various Delve types
func getFunctionName(thread *api.Thread) string {
	if thread == nil || thread.Function == nil {
		return "unknown"
	}
	return thread.Function.Name()
}

// getFunctionNameFromBreakpoint extracts a human-readable function name from a breakpoint
func getFunctionNameFromBreakpoint(bp *api.Breakpoint) string {
	if bp == nil || bp.FunctionName == "" {
		return "unknown"
	}
	return bp.FunctionName
}

// getBreakpointStatus returns a human-readable breakpoint status
func getBreakpointStatus(bp *api.Breakpoint) string {
	if bp.Disabled {
		return "disabled"
	}
	if bp.TotalHitCount > 0 {
		return "hit"
	}
	return "enabled"
}

// getStateReason returns a human-readable reason for the current state
func getStateReason(state *api.DebuggerState) string {
	if state == nil {
		return "unknown"
	}

	if state.Exited {
		return fmt.Sprintf("process exited with status %d", state.ExitStatus)
	}

	if state.Running {
		return "process is running"
	}

	if state.CurrentThread != nil && state.CurrentThread.Breakpoint != nil {
		return "hit breakpoint"
	}

	return "process is stopped"
}

// getCurrentLocation gets the current location from a DebuggerState
func getCurrentLocation(state *api.DebuggerState) *string {
	if state == nil || state.CurrentThread == nil {
		return nil
	}
	if state.CurrentThread.File == "" || state.CurrentThread.Function == nil {
		return nil
	}

	r := fmt.Sprintf("At %s:%d in %s", state.CurrentThread.File, state.CurrentThread.Line, getFunctionName(state.CurrentThread))
	return &r
}

func getBreakpointLocation(bp *api.Breakpoint) *string {
	r := fmt.Sprintf("At %s:%d in %s", bp.File, bp.Line, getFunctionNameFromBreakpoint(bp))
	return &r
}
