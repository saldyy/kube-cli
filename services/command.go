package services

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
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

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if shouldPipeOutput {
		runWithOutput(cmd)
		return
	}

	err := cmd.Run()

	if err != nil {
		log.Printf("Cannot run command: %s, args: %v\n\n", command, args)
		log.Printf("Error: %v\n", err)
		log.Printf("Stderr: %s\n", stderr.String())
	}
}

func runWithOutput(cmd *exec.Cmd) {
	// Get output pipes
	stdoutPipe, err := cmd.StdoutPipe()

	if err != nil {
		log.Fatal("Error getting StdoutPipe:", err)
	}

	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		log.Fatal("Error getting StderrPipe:", err)
		return
	}

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
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
