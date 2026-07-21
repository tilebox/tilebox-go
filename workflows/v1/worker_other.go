//go:build !unix

package workflows

import (
	"errors"
	"net"
)

func listenWorker(string) (net.Listener, func() error, error) {
	return nil, nil, errors.New("execution-only workflow workers require Unix domain socket support")
}
