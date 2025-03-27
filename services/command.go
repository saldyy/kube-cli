package services

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os/exec"
)

type CommandOptions struct {
	Args       []string
	WithOutput bool
}

func RunCommand(command string, opt CommandOptions) {
	args := opt.Args
	shouldPipeOutput := opt.WithOutput

	cmd := exec.Command(command, args...)

	if shouldPipeOutput {
		runWithOutput(cmd)
		return
	}

	err := cmd.Run()

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err != nil {
		fmt.Printf("Cannot run command: %s, args: %v\n\n", command, args)
		fmt.Printf("Error: %v\n", err)
		fmt.Printf("Stderr: %s\n", stderr.String())
	}
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
