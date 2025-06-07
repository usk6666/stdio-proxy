package main

import (
	"os/exec"
	"testing"
)

func TestExecCommand(t *testing.T) {
	cmd := exec.Command("./stdio-proxy", "exec", "echo", "hello")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("exec command failed: %v", err)
	}
	expected := "hello\n"
	if string(output) != expected {
		t.Errorf("unexpected output: got %q, want %q", string(output), expected)
	}
}

func TestShellCommand(t *testing.T) {
	cmd := exec.Command("./stdio-proxy", "shell", "echo hello")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("shell command failed: %v", err)
	}
	expected := "hello\n"
	if string(output) != expected {
		t.Errorf("unexpected output: got %q, want %q", string(output), expected)
	}
}