package main

import (
	"fmt"
	"os/exec"
	"testing"
)

func TestCheckpoint(t *testing.T) {
	cmd := exec.Command("podman", "container", "checkpoint", "76067ded0c7c", "-e", "/tmp/checkpoint.tar.gz", "-R")
	fmt.Println(cmd.String())
	out, stdErr := cmd.CombinedOutput()
	fmt.Println(string(out), (stdErr))
}
