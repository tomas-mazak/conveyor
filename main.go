package main

import (
	"github.com/tomas-mazak/conveyor/conveyor"
	"os"
	"os/signal"
	"syscall"
)

const OutputBufferSize = 4096

func main() {
	/*

	Main routine: read chan and write to stdout
	WatchDirectory routine: WatchDirectory directory, start tail routine for each new file
	Tail routine(s): read file line-by-line, prepend filename, send to chan
	  - when file renamed, finish reading, close and delete the file, exit
	 */

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)

	// TODO: consider if dying in case of log read failure is the best thing to do (kube will kill the whole pod)
	// TODO: do the logger handling nicer
	logger := conveyor.Logger{Ch: make(chan string, OutputBufferSize)}

	go conveyor.WatchDirectory("/tmp/logs", ".log", logger)

	for {
		select {
		case line := <-logger.Ch:
			os.Stdout.WriteString(line)
		case <-exit:
			// terminate
			return
		}
	}
}
