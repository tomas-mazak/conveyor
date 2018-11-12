package conveyor

import (
	"github.com/convox/inotify"
	"golang.org/x/sys/unix"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

func NewWatcher(filePath string, logger Logger) *inotify.Watcher {
	watcher, err := inotify.NewWatcher()
	if err != nil {
		logger.LogFatal(err)
	}
	// we are interested in these events:
	// IN_MODIFY: file was written into (or truncated, what should not happen) -- continue reading
	// IN_MOVE_SELF: file was moved, this means log rotation occured, read the file to the end, close and delete it
	// IN_DELETE_SELF: file was deleted, we can stop watching it
	err = watcher.AddWatch(filePath, unix.IN_MODIFY | unix.IN_DELETE_SELF | unix.IN_MOVE_SELF)
	if err != nil {
		logger.LogFatal(err)
	}
	return watcher
}

func WatchDirectory(dir string, suffix string, logger Logger) {
	// First, read the directory and start watching all files matching the pattern (delete files not matching)
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		logger.LogFatal(err)
	}
	for _, f := range files {
		if strings.HasSuffix(f.Name(), suffix) {
			go Tail(path.Join(dir, f.Name()), logger)
		} else {
			os.Remove(path.Join(dir, f.Name()))
		}
	}

	// Then, start watching the directory using Inotify, if a new matching file was created,
	// start watching it immediately
	// TODO: This is not consistent! A file can be created between listing and start watching, fix!!
	watcher, err := inotify.NewWatcher()
	if err != nil {
		logger.LogFatal(err)
	}
	defer watcher.Close()
	err = watcher.AddWatch(dir, unix.IN_CREATE | unix.IN_MOVED_TO)
	if err != nil {
		logger.LogFatal(err)
	}
	for {
		select {
		case ev := <-watcher.Event:
			if strings.HasSuffix(ev.Name, suffix) {
				// New file matching the pattern was created, let's tail it
				go Tail(ev.Name, logger)
			} else {
				// New file NOT matching the pattern was created, delete it
				os.Remove(ev.Name)
			}
		case err := <-watcher.Error:
			logger.LogError(err)
		}
	}
}
