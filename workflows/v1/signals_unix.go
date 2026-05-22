//go:build unix

package workflows

import (
	"os"
	"syscall"
)

func runnerShutdownSignals() []os.Signal {
	return []os.Signal{syscall.SIGTERM, syscall.SIGINT, syscall.SIGTSTP, syscall.SIGQUIT}
}
