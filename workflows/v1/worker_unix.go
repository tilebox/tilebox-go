//go:build unix

package workflows

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

func listenWorker(address string) (net.Listener, func() error, error) {
	socketPath, err := workerSocketPath(address)
	if err != nil {
		return nil, nil, err
	}

	if existing, err := os.Lstat(socketPath); err == nil { //nolint:gosec // The validated absolute path is the local worker socket contract.
		if existing.Mode()&os.ModeSocket == 0 {
			return nil, nil, errors.New("worker socket path already exists and is not a socket")
		}

		dialContext, cancelDial := context.WithTimeout(context.Background(), 100*time.Millisecond)
		connection, dialErr := (&net.Dialer{}).DialContext(dialContext, "unix", socketPath)
		cancelDial()
		if dialErr == nil {
			_ = connection.Close()
			return nil, nil, errors.New("worker socket is already in use")
		}
		if !errors.Is(dialErr, syscall.ECONNREFUSED) && !errors.Is(dialErr, os.ErrNotExist) {
			return nil, nil, errors.New("cannot verify existing worker socket")
		}
		if err := os.Remove(socketPath); err != nil && !errors.Is(err, os.ErrNotExist) { //nolint:gosec // Only a verified stale Unix socket is removed.
			return nil, nil, errors.New("failed to remove stale worker socket")
		}
	} else if !errors.Is(err, os.ErrNotExist) {
		return nil, nil, errors.New("failed to inspect worker socket path")
	}

	listener, err := net.ListenUnix("unix", &net.UnixAddr{Name: socketPath, Net: "unix"})
	if err != nil {
		return nil, nil, errors.New("failed to bind worker Unix socket")
	}
	if err := os.Chmod(socketPath, 0o600); err != nil { //nolint:gosec // The worker must restrict its validated local socket.
		_ = listener.Close()
		_ = os.Remove(socketPath) //nolint:gosec // The validated socket created immediately above is removed.
		return nil, nil, errors.New("failed to secure worker Unix socket")
	}
	createdSocket, err := os.Lstat(socketPath) //nolint:gosec // The validated socket created immediately above is inspected.
	if err != nil {
		_ = listener.Close()
		_ = os.Remove(socketPath) //nolint:gosec // The validated socket created immediately above is removed.
		return nil, nil, errors.New("failed to inspect bound worker Unix socket")
	}

	cleanup := func() error {
		closeErr := listener.Close()
		if errors.Is(closeErr, net.ErrClosed) {
			closeErr = nil
		}

		currentSocket, err := os.Lstat(socketPath) //nolint:gosec // The validated socket owned by this listener is inspected.
		switch {
		case errors.Is(err, os.ErrNotExist):
			return closeErr
		case err != nil:
			return errors.Join(closeErr, errors.New("failed to inspect worker socket during cleanup"))
		case !os.SameFile(createdSocket, currentSocket):
			return errors.Join(closeErr, errors.New("worker socket path was replaced before cleanup"))
		case currentSocket.Mode()&os.ModeSocket == 0:
			return errors.Join(closeErr, errors.New("worker socket path is no longer a socket"))
		default:
			return errors.Join(closeErr, os.Remove(socketPath)) //nolint:gosec // SameFile verified that this is the socket created above.
		}
	}
	return listener, cleanup, nil
}

func workerSocketPath(address string) (string, error) {
	address = strings.TrimSpace(address)
	if address == "" {
		return "", fmt.Errorf("%s is not set", workerAddressEnvironmentVariable)
	}

	var socketPath string
	switch {
	case strings.HasPrefix(address, "unix://"):
		socketPath = strings.TrimPrefix(address, "unix://")
	case strings.HasPrefix(address, "unix:"):
		socketPath = strings.TrimPrefix(address, "unix:")
	default:
		return "", errors.New("worker address must use the unix:// transport")
	}

	socketPath = filepath.Clean(socketPath)
	if !filepath.IsAbs(socketPath) || socketPath == string(filepath.Separator) {
		return "", errors.New("worker Unix socket path must be absolute")
	}
	return socketPath, nil
}
