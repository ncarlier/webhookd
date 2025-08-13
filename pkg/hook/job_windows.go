//go:build windows

package hook

import (
	"fmt"
	"os/exec"
	"syscall"
)

// setupProcessGroup sets up the command to run in a new process group
func setupProcessGroup(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{
		CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP,
	}
}

// killProcess sends a SIGKILL to the process group of the given PID
func killProcess(pid int) error {
	return exec.Command("taskkill", "/F", "/T", "/PID", fmt.Sprint(pid)).Run()
}
