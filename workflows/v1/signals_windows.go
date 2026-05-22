//go:build windows

package workflows

import (
	"os"
	"syscall"
)

func runnerShutdownSignals() []os.Signal {
	return []os.Signal{syscall.SIGTERM, syscall.SIGINT}
}
