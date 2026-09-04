//go:build !darwin && !linux

package client

func osVersion() string {
	return ""
}
