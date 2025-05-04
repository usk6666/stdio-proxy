package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
)

func main() {
	logFilePath := os.Getenv("LOG_FILE")
	if logFilePath == "" {
		logFilePath = "stdio-logging-proxy.log"
	}
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("Error opening log file: %v", err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.Println("Starting proxy...")

	targetCommandStr := os.Getenv("PROXY_COMMAND")
	if targetCommandStr == "" {
		targetCommandStr = "/bin/bash"
	}
	targetCommand := strings.Split(targetCommandStr, " ")

	cmd := exec.Command(targetCommand[0], targetCommand[1:]...)
	defer cmd.Process.Kill()

	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Fatalf("Error creating stdin pipe: %v", err)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatalf("Error creating stdout pipe: %v", err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Fatalf("Error creating stderr pipe: %v", err)
	}

	if err := cmd.Start(); err != nil {
		log.Fatalf("Error starting process: %v", err)
	}

	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		defer wg.Done()
		defer stdin.Close()
		reader := bufio.NewReader(os.Stdin)
		for {
			input, err := reader.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					log.Println("stdin EOF reached, closing stdin pipe")
					return
				}
				log.Printf("Error reading from stdin: %v", err)
				return
			}
			log.Printf("stdin: %s", input)

			_, err = io.WriteString(stdin, input)
			if err != nil {
				log.Printf("Error writing to stdin: %v", err)
				return
			}
		}
	}()

	go func() {
		defer wg.Done()
		reader := bufio.NewReader(stdout)
		for {
			output, err := reader.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					log.Println("stdout EOF reached, closing stdout pipe")
					return
				}
				log.Printf("Error reading from stdout: %v", err)
				return
			}

			log.Printf("stdout: %s", output)
			fmt.Println(output)
		}
	}()

	go func() {
		defer wg.Done()
		reader := bufio.NewReader(stderr)
		for {
			errOutput, err := reader.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					log.Println("stderr EOF reached, closing stderr pipe")
					return
				}
				log.Printf("Error reading from stderr: %v", err)
				return
			}
			log.Printf("stderr: %s", errOutput)
			fmt.Println(errOutput)
		}
	}()

	wg.Wait()

	if err := cmd.Wait(); err != nil {
		log.Printf("Process finished with error: %v", err)
	} else {
		log.Println("Process finished successfully")
	}
}
