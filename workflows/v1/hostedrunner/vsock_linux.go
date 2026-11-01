//go:build linux

package hostedrunner

import (
	"context"
	"errors"
	"io"
	"os"
	"time"

	"golang.org/x/sys/unix"
)

// DialVSock connects a Linux guest to host CID 2. Firecracker mediates this on
// its dedicated <uds_path>_<port> listener, never on its privileged API socket.
func DialVSock(ctx context.Context, port uint32) (io.ReadWriteCloser, error) {
	if port == 0 {
		return nil, errors.New("vsock lifecycle port must be nonzero")
	}
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	fd, err := unix.Socket(unix.AF_VSOCK, unix.SOCK_STREAM|unix.SOCK_NONBLOCK|unix.SOCK_CLOEXEC, 0)
	if err != nil {
		return nil, err
	}
	connected := false
	defer func() {
		if !connected {
			_ = unix.Close(fd)
		}
	}()
	err = unix.Connect(fd, &unix.SockaddrVM{CID: unix.VMADDR_CID_HOST, Port: port})
	if err != nil && !errors.Is(err, unix.EINPROGRESS) {
		return nil, err
	}
	for err != nil {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		polls := []unix.PollFd{{Fd: int32(fd), Events: unix.POLLOUT}} //nolint:gosec // Linux file descriptors are signed 32-bit integers.
		n, pollErr := unix.Poll(polls, 100)
		if errors.Is(pollErr, unix.EINTR) {
			continue
		}
		if pollErr != nil {
			return nil, pollErr
		}
		if n == 0 {
			continue
		}
		code, socketErr := unix.GetsockoptInt(fd, unix.SOL_SOCKET, unix.SO_ERROR)
		if socketErr != nil {
			return nil, socketErr
		}
		if code != 0 {
			return nil, unix.Errno(code)
		}
		err = nil
	}
	// NewFile detects O_NONBLOCK and registers with Go's poller. Close unblocks
	// concurrent protocol I/O, including on deadline or reconnect.
	file := os.NewFile(uintptr(fd), "hosted-lifecycle-vsock")
	connected = true
	return file, nil
}
