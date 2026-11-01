//go:build !linux

package hostedrunner

import (
	"context"
	"errors"
	"io"
)

// DialVSock is available only inside Linux guests.
func DialVSock(context.Context, uint32) (io.ReadWriteCloser, error) {
	return nil, errors.New("hosted lifecycle vsock requires Linux")
}
