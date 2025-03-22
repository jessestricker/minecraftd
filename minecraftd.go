package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"time"
)

const (
	FIFO_FILE       = "minecraft.stdin"
	FIFO_FILE_PERMS = 0o600 // u=rw,g=,o=
)

func main() {
	log.SetPrefix("[minecraftd] ")

	serverCmd := os.Args[1:]
	if len(serverCmd) == 0 {
		log.Println("error: missing server command")
		log.Println("usage: minecraftd <server-cmd ...>")
		os.Exit(2)
	}

	if err := run(serverCmd); err != nil {
		log.Println("error:", err)
		os.Exit(1)
	}
}

func run(serverCmd []string) error {
	fifo, err := newFifo(FIFO_FILE, FIFO_FILE_PERMS)
	if err != nil {
		return fmt.Errorf("failed to create FIFO: %w", err)
	}
	defer fifo.closeAndRemove()

	handleSigterm(func() {
		stopServer(fifo.w)
	})

	return runServer(fifo.r, serverCmd)
}

func runServer(stdin io.Reader, serverCmd []string) error {
	cmd := exec.Command(serverCmd[0], serverCmd[1:]...)
	cmd.Stdin = stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func stopServer(stdin io.Writer) {
	log.Println("stopping server in 3s ...")

	sendCommand(stdin, "say Stopping server in 3s ...")
	time.Sleep(3 * time.Second)
	sendCommand(stdin, "stop")
}

func sendCommand(stdin io.Writer, cmd string) {
	if _, err := io.WriteString(stdin, cmd+"\n"); err != nil {
		log.Println("error: failed to send server command:", err)
	}
}
