//go:build windows

package main

import (
	"log"
	"os"
	"os/exec"
)

func execMpv(mpv string, cmdargs []string) {
	cmd := exec.Command(mpv, cmdargs[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Println("mpv exited with error:", err)
	}
}
