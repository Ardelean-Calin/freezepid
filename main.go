package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	pid := flag.Int("pid", 1, "PID to try to freeze")
	timeout := flag.Duration("timeout", time.Minute, "Duration to wait with process hang-up [default: 1 minute]")

	flag.Parse()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	done := make(chan bool, 1)
	errChan := make(chan error, 1)

	fmt.Printf("Freezing PID %d\n", *pid)

	go func() {
		err := syscall.PtraceAttach(*pid)
		if err != nil {
			errChan <- fmt.Errorf("PtraceAttach failed: %v", err)
			return
		}
		syscall.PtraceCont(*pid, int(syscall.SIGINT))
		fmt.Println("Process frozen successfully")

		// Exit if SIGINT or if a minute passed
		select {
		case <-sigs:
			done <- true
		case <-time.After(*timeout):
			done <- true
		}
	}()
	select {
	case <-done:
	case err := <-errChan:
		log.Fatalf("Error: %v", err)
	}

	err := syscall.PtraceDetach(*pid)
	if err != nil {
		panic(err)
	}
	fmt.Println("Process unfrozen successfully")
}
