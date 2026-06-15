//go:build !windows

package main

import (
	"os"
	"syscall"
)

func execMpv(mpv string, cmdargs []string) {
	syscall.Exec(mpv, cmdargs, os.Environ())
}
