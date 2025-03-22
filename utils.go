package main

import (
	"errors"
	"io"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

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

func newFifo(path string, mode uint32) (fifo, error) {
	if err := syscall.Mkfifo(path, mode); err != nil {
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
