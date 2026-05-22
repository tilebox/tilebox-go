//go:build !unix && !windows

package workflows

import "os"

func runnerShutdownSignals() []os.Signal {
	return []os.Signal{os.Interrupt}
}
