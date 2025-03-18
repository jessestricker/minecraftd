package main

import (
	"errors"
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

const (
	FIFO_FILE = "minecraft.fifo"
)

func main() {
	log.SetPrefix("(minecraftd) ")
	if err := run(); err != nil {
		log.Fatalln("error:", err)
	}
}

func run() error {
	fifo, err := newFifo(FIFO_FILE)
	if err != nil {
		return err
	}
	defer fifo.closeAndRemove()

	handleSigterm(func() {
		stopServer(fifo.w)
	})

	return runServer(fifo.r)
}

func runServer(stdin io.Reader) error {
	cmd := exec.Command("/usr/bin/java", "-jar", "server.jar", "--nogui")
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

func handleSigterm(f func()) {
	sigterms := make(chan os.Signal, 1)
	signal.Notify(sigterms, syscall.SIGTERM)
	go func() {
		<-sigterms
		f()
	}()
}

type fifo struct {
	path string
	r    io.ReadCloser
	w    io.WriteCloser
}

func newFifo(path string) (fifo, error) {
	if err := syscall.Mkfifo(path, 0o600); err != nil {
		return fifo{}, err
	}

	var wg sync.WaitGroup

	// open read end
	wg.Add(1)
	var (
		r  *os.File
		re error
	)
	go func() {
		defer wg.Done()
		r, re = os.OpenFile(path, os.O_RDONLY, 0)
	}()

	// open write end
	var we error
	w, we := os.OpenFile(path, os.O_WRONLY, 0)

	wg.Wait()

	// handle errors
	if err := errors.Join(re, we); err != nil {
		if r != nil {
			r.Close()
		}
		if w != nil {
			w.Close()
		}
		return fifo{}, err
	}

	return fifo{path, r, w}, nil
}

func (f fifo) closeAndRemove() {
	f.w.Close()
	f.r.Close()
	os.Remove(f.path)
}
