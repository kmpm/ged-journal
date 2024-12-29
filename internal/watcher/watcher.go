package watcher

import (
	"log"

	"github.com/fsnotify/fsnotify"
)

type Watcher struct {
	fs *fsnotify.Watcher
}

func New(basedir string) (*Watcher, error) {
	// Create new watcher.
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	w := &Watcher{
		fs: watcher,
	}
	return w, nil

}

func (w *Watcher) Close() {
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
				return
			}
			log.Println("event:", event)
			if event.Has(fsnotify.Write) {
				log.Println("modified file:", event.Name)
			}
		case err, ok := <-w.fs.Errors:
			if !ok {
				return
			}
			log.Println("error:", err)
		}
	}
}
