package watcher

import (
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
)

type Event struct {
	Name string
	Path string
	At   time.Time
}

type Watcher struct {
	fs      *fsnotify.Watcher
	basedir string
	Events  chan Event
}

func New(basedir string) (*Watcher, error) {
	// Create new watcher.
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	w := &Watcher{
		fs:      watcher,
		basedir: basedir,
		Events:  make(chan Event, 2),
	}
	if info, err := os.Stat(basedir); err != nil {
		return nil, err
	} else if !info.IsDir() {
		return nil, os.ErrNotExist
	}
	w.Add(basedir)
	return w, nil

}

func (w *Watcher) Close() {
	close(w.Events)
	w.fs.Close()
}

func (w *Watcher) Add(path string) error {
	return w.fs.Add(path)
}

func (w *Watcher) Watch() {
	// Start listening for events.
	for {
		select {
		case event, ok := <-w.fs.Events:
			if !ok {
				slog.Debug("fsnotify event channel closed")
				return
			}
			slog.Debug("watcher event", "event_op", event.Op, "event_name", event.Name)
			if event.Has(fsnotify.Write) || event.Has(fsnotify.Create) {
				// slog.Info("file modified", "file", event.Name)

				w.Events <- Event{Name: filepath.Base(event.Name), Path: event.Name, At: time.Now()}
			}
			// if event.Has(fsnotify.Create) {
			// 	log.Println("created file:", event.Name)
			// }
		case err, ok := <-w.fs.Errors:
			if !ok {
				slog.Debug("fsnotify error channel closed")
				return
			}
			log.Println("error:", err)
		}
	}
}
