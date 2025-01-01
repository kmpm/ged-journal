package collector

import (
	"bufio"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/kmpm/ged-journal/internal/watcher"
	"github.com/kmpm/ged-journal/public/journal"
)

var journalFilePattern = regexp.MustCompile(`^Journal\.\d{4}\-\d{2}\-\d{2}T\d{6}\.\d{2}\.log$`)

type Publisher func(subject string, data []byte, compress bool) error

type Collector struct {
	logPath        string
	w              *watcher.Watcher
	currentJournal string
	pub            Publisher
	prefix         string
}

func New(logPath string, pub Publisher) (*Collector, error) {

	a := &Collector{
		logPath: logPath,
		pub:     pub,
		prefix:  "ged.",
	}

	w, err := watcher.New(logPath)
	if err != nil {
		return nil, err
	}
	a.w = w
	go a.watchWorker()
	return a, nil
}

func (cr *Collector) Close() error {
	if cr.w != nil {
		cr.w.Close()
	}
	return nil
}

func (cr *Collector) watchWorker() {
	go cr.w.Watch()
	var ctx context.Context
	var cancel context.CancelFunc
	for event := range cr.w.Events {
		slog.Debug("something happened", "event", event)
		switch event.Name {
		case "Backpack.json":
			cr.publishJSON(event.Path, "global.backpack")
		case "Cargo.json":
			cr.publishJSON(event.Path, "global.cargo")
		case "Market.json":
			cr.publishJSON(event.Path, "global.market")
		case "ModulesInfo.json":
			cr.publishJSON(event.Path, "global.modulesinfo")
		case "NavRoute.json":
			cr.publishJSON(event.Path, "global.navroute")
		case "Outfitting.json":
			cr.publishJSON(event.Path, "global.outfitting")
		case "ShipLocker.json":
			cr.publishJSON(event.Path, "global.shiplocker")
		case "Shipyard.json":
			cr.publishJSON(event.Path, "global.shipyard")
		case "Status.json":
			cr.publishJSON(event.Path, "global.status")
		default:
			if journalFilePattern.MatchString(event.Name) {
				if cr.currentJournal != event.Name {
					if cancel != nil {
						cancel()
					}
					ctx, cancel = context.WithCancel(context.Background())
					cr.currentJournal = event.Name
					go cr.scanJournal(ctx, event.Path)
				}
			} else {
				slog.Warn("unknown file", "file", event.Name)
			}
		}
	}
	if cancel != nil {
		cancel()
	}
	slog.Debug("watcher event channel closed")
}
func (cr *Collector) publishJSON(filename string, subject string) error {
	var err error
	defer func() {
		if err != nil {
			slog.Warn("failed to publish json file", "error", err, "subject", subject, "filename", filename)
		}
	}()
	var obj map[string]any
	obj, err = cr.readJSON(filename)
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
	err = cr.pub(subject, data, false)

	return err
}

func (cr *Collector) readJSON(filename string) (map[string]any, error) {
	data, err := os.ReadFile(filename)
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

func (cr *Collector) scanJournal(ctx context.Context, filename string) error {
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

				err = cr.pub(strings.ToLower("journal.event."+obj.Event), []byte(line), false)
				if err != nil {
					slog.Warn("failed to publish journal entry", "line", line, "error", err)
				}
			}
		}
	}()
	return nil
}
