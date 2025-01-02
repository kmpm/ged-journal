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
	// pattern is <subject>_<epoc>.json
	// split the name by _
	parts := strings.Split(name, "_")
	// get the last part and split the last part by .
	parts = strings.Split(parts[len(parts)-1], ".")
	// get the first part
	unix, err := strconv.Atoi(parts[0])
	if err != nil {
		panic(err)
	}

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
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		// get subject from name
		parts := strings.Split(e.Name(), "_")
		parts = strings.Split(parts[0], "-")
		subject := strings.Join(parts, ".")

		if !strings.HasPrefix(subject, prefix) {
			slog.Info("Subject does not have prefix", "subject", subject, "prefix", prefix)
			subject = prefix + subject
		}
		data, err := os.ReadFile(filepath.FromSlash(folder + "/" + e.Name()))
		if err != nil {
			panic(err)
		}
		slog.Info("Sending message", "subject", subject, "name", e.Name())
		err = nc.Publish(subject, data)
		if err != nil {
			panic(err)
		}
		time.Sleep(cli.Delay)
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
