package agent

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
)

// Fuel contains fuel readouts for the ship.
type Fuel struct {
	Main      float64 `json:"FuelMain"`
	Reservoir float64 `json:"FuelReservoir"`
}

// Status represents the current state of the player and ship.
type Status struct {
	Timestamp string      `json:"timestamp"`
	Event     string      `json:"event"`
	Flags     StatusFlags `json:"-"`
	RawFlags  uint32      `json:"Flags"`
	Pips      [3]int32    `json:"Pips"`
	FireGroup int32       `json:"FireGroup"`
	GuiFocus  int32       `json:"GuiFocus"`
	Fuel      Fuel        `json:"Fuel"`
	Cargo     float64     `json:"Cargo"`
	Latitude  float64     `json:"Latitude,omitempty"`
	Longitude float64     `json:"Longitude,omitempty"`
	Heading   int32       `json:"Heading,omitempty"`
	Altitude  int32       `json:"Altitude,omitempty"`
}

func getStatusFromPath(logPath string) (*Status, error) {
	statusFilePath := filepath.FromSlash(logPath + "/Status.json")
	f, err := os.Open(statusFilePath)
	if err != nil {
		return nil, errors.New("couldn't open Status.json file: " + err.Error())
	}
	defer f.Close()
	data, err := io.ReadAll(f)
	if err != nil {
		return nil, errors.New("couldn't read Status.json file: " + err.Error())
	}
	return GetStatusFromBytes(data)
}

// GetStatusFromBytes reads the current player and ship status from the string contained in the byte array.
func GetStatusFromBytes(content []byte) (*Status, error) {
	status := &Status{}
	if err := json.Unmarshal(content, status); err != nil {
		return nil, errors.New("couldn't unmarshal Status.json file: " + err.Error())
	}

	status.ExpandFlags()
	return status, nil
}
