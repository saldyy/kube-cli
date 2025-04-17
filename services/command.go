package services

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"strings"
)

type CommandOptions struct {
	Args       []string
	WithOutput bool
}

func RunWithLiveOutput(command string, args []string) error {
	cmd := exec.Command(command, args...)

	// Get output pipes
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("error getting StdoutPipe: %v", err)
	}

	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("error getting StderrPipe: %v", err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("error starting command: %v", err)
	}

	// Stream output in separate goroutines
	go streamOutput(&stdoutPipe, "STDOUT")
	go streamOutput(&stderrPipe, "STDERR")

	// Wait for command to complete
	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("error waiting for command: %v", err)
	}

	return nil
}

func RunSilent(command string, args []string) (string, error) {
	cmd := exec.Command(command, args...)

	var stderr bytes.Buffer
	var stdout bytes.Buffer
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("cannot run command %s with args %v: %v\nStderr: %s",
			command, args, err, stderr.String())
	}

	return strings.TrimSpace(stdout.String()), nil
}

func runWithOutput(cmd *exec.Cmd) {
	fmt.Printf("INheree\n")
	// Get output pipes
	stdoutPipe, err := cmd.StdoutPipe()

	if err != nil {
		fmt.Printf("Error getting StdoutPipe: %v\n", err)
	}

	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		fmt.Printf("Error getting StderrPipe: %v\n", err)
	}

	if err := cmd.Start(); err != nil {
		fmt.Println(err)
	}
	// Stream output in separate goroutines
	go streamOutput(&stdoutPipe, "STDOUT")
	go streamOutput(&stderrPipe, "STDERR")

	// Wait for command to complete
	if err := cmd.Wait(); err != nil {
		fmt.Println("Error waiting for command:", err)
	}
}

// streamOutput reads from the provided pipe and prints to the console
func streamOutput(pipe *io.ReadCloser, prefix string) {
	scanner := bufio.NewScanner(*pipe)
	for scanner.Scan() {
		fmt.Printf("[%s] %s\n", prefix, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading from %s: %v\n", prefix, err)
	}
}
