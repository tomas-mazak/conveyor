package conveyor

import (
	"bufio"
	"github.com/convox/inotify"
	"golang.org/x/sys/unix"
	"io"
	"os"
	"path"
)

func getEvent(watcher *inotify.Watcher, logger Logger) bool {
	select {
		case ev := <-watcher.Event:
			// return true if event was different from IN_MODIFY (all other events mean the file was rotated)
			return ev.Mask != unix.IN_MODIFY
		case err := <-watcher.Error:
			logger.LogFatal(err)
	}
	panic("unreachable")
}

func Tail(fileName string, logger Logger) {
	prefix := path.Base(fileName) + ": "

	watcher := NewWatcher(fileName, logger)

	f, err := os.Open(fileName)
	if err != nil {
		logger.LogFatal(err)
	}
	reader := bufio.NewReader(f)

	rotated := false
	for {
		offset, err := f.Seek(0, 1) // move by 0 from current offset - return the current position
		if err != nil {
			logger.LogFatal(err)
		}
		line, err := reader.ReadString('\n')
		if err == nil {
			logger.Log(prefix + line)
		} else if err == io.EOF {
			// If EOF was reached and the file was already rotated, close and delete it
			if rotated {
				f.Close()
				watcher.Close()
				return
			}
			// If EOF was reached and non-empty string was read, we have a partial line, return to the beginning
			// and try again
			if line != "" {
				f.Seek(offset, 0)
			}
			// Wait for an event to occur
			rotated = getEvent(watcher, logger)
		} else {
			logger.LogFatal(err)
		}
	}
}
