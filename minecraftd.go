package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	serverCmd := os.Args[1:]
	if len(serverCmd) == 0 {
		fmt.Println("error: missing server command")
		fmt.Println("usage: minecraftd <server-cmd ...>")
		os.Exit(2)
	}

	if err := run(serverCmd); err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
}

func run(serverCmd []string) error {
	srv, err := newServer(serverCmd)
	if err != nil {
		return err
	}
	defer srv.close()

	handleSigterm(func() {
		srv.send("say Stopping server in 3s ...")
		time.Sleep(3 * time.Second)
		srv.send("stop")
	})

	return srv.run()
}

type server struct {
	command    []string
	stdinRead  *os.File
	stdinWrite *os.File
}

func newServer(command []string) (*server, error) {
	stdinRead, stdinWrite, err := os.Pipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdin pipe: %w", err)
	}
	return &server{command, stdinRead, stdinWrite}, nil
}

func (s server) close() {
	s.stdinWrite.Close()
	s.stdinRead.Close()
}

func (s server) run() error {
	cmd := exec.Command(s.command[0], s.command[1:]...)
	cmd.Stdin = s.stdinRead
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (s server) send(cmd string) error {
	_, err := io.WriteString(s.stdinWrite, cmd+"\n")
	return err
}

func handleSigterm(handler func()) {
	channel := make(chan os.Signal, 1)
	go func() {
		<-channel
		handler()
		signal.Stop(channel)
	}()
	signal.Notify(channel, syscall.SIGTERM)
}
