package main

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/kmpm/ged-journal/internal/misccli"
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
	nc, err := misccli.ConnectNATS(cli.Nats, cli.NatsContext, cli.NatsCreds)
	if err != nil {
		panic(err)
	}
	defer func() {
		nc.Flush()
		nc.Close()
	}()
	fmt.Println("Getting ready to send messages")
	time.Sleep(5 * time.Second)
	var prevEpoc, thisEpoc int
	var subject string
	length := len(entries)
	if length == 0 {
		slog.Info("No files found")
		fmt.Println("No files found")
		return nil
	}
	count := 0
	tick := time.NewTicker(2 * time.Second)
	go func() {
		prevCount := count
		start := time.Now()
		for range tick.C {
			elapsed := time.Since(start)
			msgsPerSec := float64(count-prevCount) / elapsed.Seconds()
			fmt.Printf("Progress: %d/%d, Rate: %d msg/s\n", count, length, int(msgsPerSec))
			slog.Info("Progress", "count", count, "length", length, "rate", int(msgsPerSec))
			start = time.Now()
			prevCount = count

		}
	}()
	defer tick.Stop()

	fmt.Println("Sending messages")
	for _, e := range entries {
		count++
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
	fmt.Println("All messages sent")
	return nil
}
