package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

func main() {

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s <pid>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nPositional Arguments:\n")
		fmt.Fprintf(os.Stderr, "  pid\tPID to try to freeze\n")
		flag.PrintDefaults()
	}

	flag.Parse()
	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(1)
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	done := make(chan bool, 1)

	pid, err := strconv.Atoi(flag.Arg(0))
	if err != nil {
		panic(err)
	}
	fmt.Printf("Freezing PID %d\n", pid)

	go func() {
		err := syscall.PtraceAttach(pid)
		if err != nil {
			panic(err)
		}

		<-sigs
		done <- true
	}()
	<-done

	err = syscall.PtraceDetach(pid)
	if err != nil {
		panic(err)
	}
	fmt.Println("exiting")
}
