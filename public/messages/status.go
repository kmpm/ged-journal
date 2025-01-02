package messages

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/kmpm/ged-journal/internal/flags"
)

// Fuel contains fuel readouts for the ship.
type Fuel struct {
	Main      float64 `json:"FuelMain"`
	Reservoir float64 `json:"FuelReservoir"`
}

// StatusFlags contains boolean flags describing the player and ship.
type StatusFlags struct {
	Docked                    bool
	Landed                    bool
	LandingGearDown           bool
	ShieldsUp                 bool
	Supercruise               bool
	FlightAssistOff           bool
	HardpointsDeployed        bool
	InWing                    bool
	LightsOn                  bool
	CargoScoopDeployed        bool
	SilentRunning             bool
	ScoopingFuel              bool
	SRVHandbrake              bool
	SRVTurret                 bool
	SRVUnderShip              bool
	SRVDriveAssist            bool
	FSDMassLocked             bool
	FSDCharging               bool
	FSDCooldown               bool
	LowFuel                   bool
	Overheating               bool
	HasLatLong                bool
	IsInDanger                bool
	BeingInterdicted          bool
	InMainShip                bool
	InFighter                 bool
	InSRV                     bool
	InAnalysisMode            bool
	NightVision               bool
	AltitudeFromAverageRadius bool
	FSDJump                   bool
	SRVHighBeam               bool
}

// ExpandFlags parses the RawFlags and sets the Flags values accordingly.
func (status *Status) ExpandFlags() {
	status.Flags.Docked = status.RawFlags&flags.Docked != 0
	status.Flags.Landed = status.RawFlags&flags.Landed != 0
	status.Flags.LandingGearDown = status.RawFlags&flags.LandingGearDown != 0
	status.Flags.ShieldsUp = status.RawFlags&flags.ShieldsUp != 0
	status.Flags.Supercruise = status.RawFlags&flags.Supercruise != 0
	status.Flags.FlightAssistOff = status.RawFlags&flags.FlightAssistOff != 0
	status.Flags.HardpointsDeployed = status.RawFlags&flags.HardpointsDeployed != 0
	status.Flags.InWing = status.RawFlags&flags.InWing != 0
	status.Flags.LightsOn = status.RawFlags&flags.LightsOn != 0
	status.Flags.CargoScoopDeployed = status.RawFlags&flags.CargoScoopDeployed != 0
	status.Flags.SilentRunning = status.RawFlags&flags.SilentRunning != 0
	status.Flags.ScoopingFuel = status.RawFlags&flags.ScoopingFuel != 0
	status.Flags.SRVHandbrake = status.RawFlags&flags.SRVHandbrake != 0
	status.Flags.SRVTurret = status.RawFlags&flags.SRVTurret != 0
	status.Flags.SRVUnderShip = status.RawFlags&flags.SRVUnderShip != 0
	status.Flags.SRVDriveAssist = status.RawFlags&flags.SRVDriveAssist != 0
	status.Flags.FSDMassLocked = status.RawFlags&flags.FSDMassLocked != 0
	status.Flags.FSDCharging = status.RawFlags&flags.FSDCharging != 0
	status.Flags.FSDCooldown = status.RawFlags&flags.FSDCooldown != 0
	status.Flags.LowFuel = status.RawFlags&flags.LowFuel != 0
	status.Flags.Overheating = status.RawFlags&flags.Overheating != 0
	status.Flags.HasLatLong = status.RawFlags&flags.HasLatLong != 0
	status.Flags.IsInDanger = status.RawFlags&flags.IsInDanger != 0
	status.Flags.BeingInterdicted = status.RawFlags&flags.BeingInterdicted != 0
	status.Flags.InMainShip = status.RawFlags&flags.InMainShip != 0
	status.Flags.InFighter = status.RawFlags&flags.InFighter != 0
	status.Flags.InSRV = status.RawFlags&flags.InSRV != 0
	status.Flags.InAnalysisMode = status.RawFlags&flags.InAnalysisMode != 0
	status.Flags.NightVision = status.RawFlags&flags.NightVision != 0
	status.Flags.AltitudeFromAverageRadius = status.RawFlags&flags.AltitudeFromAverageRadius != 0
	status.Flags.FSDJump = status.RawFlags&flags.FSDJump != 0
	status.Flags.SRVHighBeam = status.RawFlags&flags.SRVHighBeam != 0
}

// Status represents the current state of the player and ship.
type Status struct {
	Event
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

func getStatusFromPath(logPath string) (Status, error) {
	status := Status{}
	statusFilePath := filepath.FromSlash(logPath + "/Status.json")
	f, err := os.Open(statusFilePath)
	if err != nil {
		return status, errors.New("couldn't open Status.json file: " + err.Error())
	}
	defer f.Close()
	data, err := io.ReadAll(f)
	if err != nil {
		return status, errors.New("couldn't read Status.json file: " + err.Error())
	}

	err = GetStatus(data, &status)
	return status, err
}

// GetStatus reads the current player and ship status from the string contained in the byte array.
func GetStatus(content []byte, status *Status) error {
	if err := json.Unmarshal(content, status); err != nil {
		return errors.New("couldn't unmarshal Status.json file: " + err.Error())
	}
	status.ExpandFlags()
	return nil
}
