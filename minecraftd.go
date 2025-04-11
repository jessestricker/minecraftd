package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

func main() {
	if err := run(); err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := loadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	srv, err := newServer()
	if err != nil {
		return err
	}
	defer srv.close()

	handleSigterm(func() {
		srv.send("stop")
	})

	return srv.run(cfg)
}

type server struct {
	stdinRead  *os.File
	stdinWrite *os.File
}

func newServer() (*server, error) {
	stdinRead, stdinWrite, err := os.Pipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdin pipe: %w", err)
	}
	return &server{stdinRead, stdinWrite}, nil
}

func (s server) close() {
	s.stdinWrite.Close()
	s.stdinRead.Close()
}

func (s server) run(cfg *config) error {
	args := cfg.serverCommand()
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Dir = cfg.Server.Dir
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
