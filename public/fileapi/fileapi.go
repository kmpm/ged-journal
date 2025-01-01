package fileapi

import (
	"bufio"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/kmpm/ged-journal/internal/watcher"
	"github.com/kmpm/ged-journal/public/journal"
	"github.com/nats-io/nats.go"
)

var journalFilePattern = regexp.MustCompile(`^Journal\.\d{4}\-\d{2}\-\d{2}T\d{6}\.\d{2}\.log$`)

type API struct {
	logPath        string
	w              *watcher.Watcher
	mu             sync.Mutex
	currentJournal string
	nc             *nats.Conn
	Status         *Status
}

func New(logPath string, nc *nats.Conn) (*API, error) {

	a := &API{
		logPath: logPath,
		nc:      nc,
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

func (a *API) Close() error {
	if a.w != nil {
		a.w.Close()
	}
	return nil
}

func (a *API) watchWorker() {
	go a.w.Watch()
	var ctx context.Context
	var cancel context.CancelFunc
	var err error
	var data []byte
	for event := range a.w.Events {
		slog.Debug("something happened", "event", event)
		switch event.Name {
		case "Backpack.json":
			a.publishJSON(event.Path, "global.backpack")
		case "Cargo.json":
			a.publishJSON(event.Path, "global.cargo")
		case "Market.json":
			a.publishJSON(event.Path, "global.market")
		case "ModulesInfo.json":
			a.publishJSON(event.Path, "global.modulesinfo")
		case "NavRoute.json":
			a.publishJSON(event.Path, "global.navroute")
		case "Outfitting.json":
			a.publishJSON(event.Path, "global.outfitting")
		case "ShipLocker.json":
			a.publishJSON(event.Path, "global.shiplocker")
		case "Shipyard.json":
			a.publishJSON(event.Path, "global.shipyard")
		case "Status.json":
			data, err = a.readFile(event.Path)
			if err != nil {
				slog.Warn("failed to read status file", "error", err)
				continue
			}
			if len(data) == 0 {
				continue
			}
			s, err := GetStatusFromBytes(data)
			if err != nil {
				slog.Warn("failed to parse status file", "error", err)
				continue
			}
			s.ExpandFlags()
			data, err = json.Marshal(s)
			if err != nil {
				slog.Warn("failed to marshal status file", "error", err)
				continue
			}
			a.nc.Publish("global.status", data)

		default:
			if journalFilePattern.MatchString(event.Name) {
				if a.currentJournal != event.Name {
					if cancel != nil {
						cancel()
					}
					ctx, cancel = context.WithCancel(context.Background())
					a.currentJournal = event.Name
					go a.scanJournal(ctx, event.Path)
				}
			} else {
				slog.Warn("unknown file", "file", event.Name)
			}
		}
		// a.Refresh()
	}
	if cancel != nil {
		cancel()
	}
	slog.Debug("watcher event channel closed")
}
func (a *API) publishJSON(filename string, subject string) error {
	var err error
	defer func() {
		if err != nil {
			slog.Warn("failed to publish json file", "error", err, "subject", subject, "filename", filename)
		}
	}()
	var obj map[string]any
	obj, err = a.readJSON(filename)
	if err != nil {
		return err
	}
	if obj == nil {
		return nil
	}
	data, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	err = a.nc.Publish(subject, data)

	return err
}

func (a *API) readJSON(filename string) (map[string]any, error) {
	data, err := a.readFile(filename)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, nil
	}
	obj := map[string]any{}
	err = json.Unmarshal(data, &obj)
	if err != nil {
		return nil, err
	}
	return obj, nil
}

func (a *API) readFile(filename string) ([]byte, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return io.ReadAll(file)
}

func (a *API) scanJournal(ctx context.Context, filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}

	go func() {
		slog.Info("scanJournal start", "filename", filename)
		defer file.Close()
		defer slog.Info("scanJournal done", "filename", filename)
		reader := bufio.NewReader(file)
		for {
			select {
			case <-ctx.Done():
				slog.Debug("scanJournal cancelled")
				return
			default:
				line, err := reader.ReadString('\n')
				if err != nil {
					if err == io.EOF {
						// without this sleep you would hogg the CPU
						time.Sleep(500 * time.Millisecond)
						continue
					}
					slog.Warn("failed to read journal entry", "error", err)
					break
				}
				obj, err := journal.Parse([]byte(line))
				if err != nil {
					slog.Warn("failed to parse journal entry", "line", line, "error", err)
					continue
				}

				err = a.nc.Publish(strings.ToLower("journal.event."+obj.Event), []byte(line))
				if err != nil {
					slog.Warn("failed to publish journal entry", "line", line, "error", err)
				}
			}
		}
	}()
	return nil
}

func (a *API) Refresh() error {
	a.mu.Lock()
	defer a.mu.Unlock()
	s, err := a.GetStatus()
	if err != nil {
		return err
	}
	a.Status = s
	return nil
}
