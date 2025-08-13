//go:build !windows

package hook

import (
	"os/exec"
	"syscall"
)

// setupProcessGroup sets up the command to run in a new process group
func setupProcessGroup(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
}

// killProcess sends a SIGKILL to the process group of the given PID
func killProcess(pid int) error {
	return syscall.Kill(-pid, syscall.SIGKILL)
}
