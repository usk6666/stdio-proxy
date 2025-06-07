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

	"github.com/urfave/cli/v2"
)

func execCommand(name string, args []string) {
	log.Println("stdio-proxy-msg: Starting proxy...")
	log.Printf("stdio-proxy-msg: Name: %s", name)
	log.Printf("stdio-proxy-msg: Args: %v", args)
	cmd := exec.Command(name, args...)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Fatalf("stdio-proxy-msg: Error creating stdin pipe: %v", err)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatalf("stdio-proxy-msg: Error creating stdout pipe: %v", err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Fatalf("stdio-proxy-msg: Error creating stderr pipe: %v", err)
	}

	if err := cmd.Start(); err != nil {
		log.Fatalf("stdio-proxy-msg: Error starting process: %v", err)
	}
	defer cmd.Process.Kill()

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
					log.Println("stdio-proxy-msg: stdin EOF reached, closing stdin pipe")
					return
				}
				log.Printf("stdio-proxy-msg: Error reading from stdin: %v", err)
				return
			}
			log.Printf("stdin: %s", input)

			_, err = io.WriteString(stdin, input)
			if err != nil {
				log.Printf("stdio-proxy-msg: Error writing to stdin: %v", err)
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
					log.Println("stdio-proxy-msg: stdout EOF reached, closing stdout pipe")
					return
				}
				log.Printf("stdio-proxy-msg: Error reading from stdout: %v", err)
				return
			}

			log.Printf("stdout: %s", output)
			fmt.Print(output)
		}
	}()

	go func() {
		defer wg.Done()
		reader := bufio.NewReader(stderr)
		for {
			errOutput, err := reader.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					log.Println("stdio-proxy-msg: stderr EOF reached, closing stderr pipe")
					return
				}
				log.Printf("stdio-proxy-msg: Error reading from stderr: %v", err)
				return
			}
			log.Printf("stderr: %s", errOutput)
			fmt.Print(errOutput)
		}
	}()

	wg.Wait()

	if err := cmd.Wait(); err != nil {
		log.Printf("stdio-proxy-msg: Process finished with error: %v", err)
	} else {
		log.Println("stdio-proxy-msg: Process finished successfully")
	}
}

func main() {
	var outputFlag = &cli.StringFlag{
		Name:    "output",
		Value:   "stdio-proxy.log",
		Usage:   "Path to the output log file",
		EnvVars: []string{"STDIO_PROXY_OUTPUT"},
	}

	app := &cli.App{
		Name:  "stdio-proxy",
		Usage: "Proxy for logging stdio",
		UsageText: "stdio-proxy exec|shell [options] command",
		Flags: []cli.Flag{
			outputFlag,
		},
		Commands: []*cli.Command{
			{
				Name:      "exec",
				Usage:     "Execute a command directly",
				UsageText: "stdio-proxy exec [options] command",
				Flags: []cli.Flag{
					outputFlag,
				},
				Action: func(c *cli.Context) error {
					if c.NArg() < 1 {
						fmt.Println("You must provide a command to run.")
						return nil
					}

					logFilePath := c.String("output")
					logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
					if err != nil {
						log.Fatalf("stdio-proxy-msg: Error opening log file: %v", err)
					}
					defer logFile.Close()
					log.SetOutput(logFile)
					log.SetFlags(log.Ldate | log.Ltime)

					targetCommand := c.Args().Slice()
					execCommand(targetCommand[0], targetCommand[1:])

					return nil
				},
			},
			{
				Name:      "shell",
				Usage:     "Execute a command in a shell",
				UsageText: "stdio-proxy shell [options] command",
				Flags: []cli.Flag{
					outputFlag,
				},
				Action: func(c *cli.Context) error {
					if c.NArg() < 1 {
						fmt.Println("You must provide a command to run.")
						return nil
					}

					logFilePath := c.String("output")
					logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
					if err != nil {
						log.Fatalf("stdio-proxy-msg: Error opening log file: %v", err)
					}
					defer logFile.Close()
					log.SetOutput(logFile)
					log.SetFlags(log.Ldate | log.Ltime)

					targetCommand := c.Args().Slice()
					shell := os.Getenv("SHELL")
					if shell == "" {
						shell = "/bin/sh"
					}
					execCommand(shell, append([]string{"-c"}, strings.Join(targetCommand, " ")))

					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
	}
}
