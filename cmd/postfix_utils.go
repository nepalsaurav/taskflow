package cmd

import (
	"os/exec"
	"strings"
)

// checkPostfixRunning ensures Postfix master process is running
func CheckPostfixRunning() bool {
	cmd := exec.Command("ps", "-C", "master", "-o", "pid=")
	out, _ := cmd.Output()
	if strings.TrimSpace(string(out)) == "" {
		return false
	}
	return true
}
