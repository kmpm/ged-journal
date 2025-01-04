package main

import (
	"errors"
	"log/slog"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/nats-io/jsm.go/natscontext"
	"github.com/nats-io/nats.go"
)

func epocFromName(name string) (unix int) {
	// pattern is <epoc>_<subject>.json
	// split the name by _
	parts := strings.Split(name, "_")
	// first part is epoc
	unix, err := strconv.Atoi(parts[0])
	if err != nil {
		panic(err)
	}
	return
}

// parseName returns the epoc and subject from the name.
// The name is expected to be in the format <epoc>_<subject>.json
func parseName(name string) (unix int, subject string) {
	// pattern is <epoc>_<subject>.json
	var err error
	// split the name by _
	parts := strings.Split(name, "_")
	// first part is epoc
	unix, err = strconv.Atoi(parts[0])
	if err != nil {
		panic(err)
	}
	subject = strings.ReplaceAll(strings.Split(parts[1], ".")[0], "-", ".")
	return
}

func simulator(cli *Cli, folder, prefix string) error {
	slog.Info("Running Simulator")
	//read all files from the folder
	//for each file, read the content and send it to the nats server
	//wait for the delay
	//repeat
	entries, err := os.ReadDir(folder)
	if err != nil {
		panic(err)
	}
	// sort entries by epoc in name
	sort.SliceStable(entries, func(i, j int) bool {
		unixI := epocFromName(entries[i].Name())
		unixJ := epocFromName(entries[j].Name())
		return unixI < unixJ
	})
	nc, err := connect(cli.Nats, cli.NatsContext)
	if err != nil {
		panic(err)
	}
	var prevEpoc, thisEpoc int
	var subject string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		// get subject from name
		thisEpoc, subject = parseName(e.Name())

		if !strings.HasPrefix(subject, prefix) {
			slog.Info("Subject does not have prefix", "subject", subject, "prefix", prefix)
			subject = prefix + subject
		}
		data, err := os.ReadFile(filepath.FromSlash(folder + "/" + e.Name()))
		if err != nil {
			panic(err)
		}
		// relative message spacing same as files
		if prevEpoc != 0 {
			time.Sleep(time.Duration(thisEpoc-prevEpoc) * cli.Delay)
		}
		prevEpoc = thisEpoc
		slog.Info("Sending message", "subject", subject, "name", e.Name())
		err = nc.Publish(subject, data)
		if err != nil {
			panic(err)
		}
		// time.Sleep(cli.Delay)
	}
	slog.Info("All messages sent")
	return nil
}

func connect(uri, context string) (nc *nats.Conn, err error) {
	if context != "" {
		nc, err = natscontext.Connect("nats_development", nil)
	} else if uri != "" {
		nc, err = nats.Connect(uri)
	} else {
		return nil, errors.New("no nats server address provided")
	}
	if err != nil {
		return nil, err
	}
	nc.SetClosedHandler(func(_ *nats.Conn) {
		slog.Error("nats connection closed")
	})
	nc.SetDisconnectHandler(func(_ *nats.Conn) {
		slog.Error("nats connection disconnected")
	})
	nc.SetDisconnectErrHandler(func(_ *nats.Conn, err error) {
		slog.Error("nats connection disconnected", "error", err)
	})
	nc.SetReconnectHandler(func(_ *nats.Conn) {
		slog.Info("nats connection reconnected")
	})

	return nc, nil
}
