package services

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os/exec"
)

const PROFILE_NAME = "saldyy"

func InitCluster() {
	args := []string{
		"start", "-p", PROFILE_NAME, "--addons", "metallb,metrics-server", "--cpus=4", "--memory=16GB",
	}

	cmd := exec.Command("minikube", args...)
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

func DestroyCluster() {
	cmd := exec.Command("minikube", "delete", "-p", PROFILE_NAME)
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}
