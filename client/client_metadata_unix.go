//go:build darwin || linux

package client

import (
	"bytes"

	"golang.org/x/sys/unix"
)

func osVersion() string {
	var system unix.Utsname
	if unix.Uname(&system) != nil {
		return ""
	}
	return string(bytes.TrimRight(system.Release[:], "\x00"))
}
