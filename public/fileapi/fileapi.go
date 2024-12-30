package fileapi

import (
	"bufio"
	"fmt"
	"io"
	"log/slog"
	"os"
	"regexp"
	"sync"
	"time"

	"github.com/kmpm/ged-journal/internal/watcher"
)

var journalFilePattern = regexp.MustCompile(`^Journal\.\d{4}\-\d{2}\-\d{2}T\d{6}\.\d{2}\.log$`)

type Api struct {
	logPath        string
	w              *watcher.Watcher
	mu             sync.Mutex
	currentJournal string
	Status         *Status
}

func New(logPath string) (*Api, error) {

	a := &Api{
		logPath: logPath,
	}
	err := a.Refresh()
	if err != nil {
		return nil, err
	}
	w, err := watcher.New(logPath)
	if err != nil {
		return nil, err
	}
	a.w = w
	go a.watchWorker()
	return a, nil
}

func (a *Api) Close() error {
	if a.w != nil {
		a.w.Close()
	}
	return nil
}

func (a *Api) watchWorker() {
	go a.w.Watch()
	for event := range a.w.Events {
		slog.Info("something happened", "event", event)
		if journalFilePattern.MatchString(event.Name) {
			if a.currentJournal != event.Name {
				a.currentJournal = event.Name
			}
		}
		// a.Refresh()
	}
	slog.Debug("watcher event channel closed")
}

func (a *Api) scanJournal(filename string) error {
	file, err := os.Open("./test.log")
	if err != nil {
		return err
	}
	go func() {
		defer file.Close()
		reader := bufio.NewReader(file)
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					// without this sleep you would hogg the CPU
					time.Sleep(500 * time.Millisecond)
					continue
				}

				break
			}

			fmt.Printf("%s\n", string(line))
		}
	}()

}

func (a *Api) Refresh() error {
	a.mu.Lock()
	defer a.mu.Unlock()
	s, err := a.GetStatus()
	if err != nil {
		return err
	}
	a.Status = s
	return nil
}
