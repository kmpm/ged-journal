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

// StatusEvent represents the current state of the player and ship.
type StatusEvent struct {
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

// ExpandFlags parses the RawFlags and sets the Flags values accordingly.
func (status *StatusEvent) ExpandFlags() {
	ExpandFlags(status.RawFlags, &status.Flags)
}

func ExpandFlags(raw uint32, f *StatusFlags) {
	f.Docked = raw&flags.Docked != 0
	f.Landed = raw&flags.Landed != 0
	f.LandingGearDown = raw&flags.LandingGearDown != 0
	f.ShieldsUp = raw&flags.ShieldsUp != 0
	f.Supercruise = raw&flags.Supercruise != 0
	f.FlightAssistOff = raw&flags.FlightAssistOff != 0
	f.HardpointsDeployed = raw&flags.HardpointsDeployed != 0
	f.InWing = raw&flags.InWing != 0
	f.LightsOn = raw&flags.LightsOn != 0
	f.CargoScoopDeployed = raw&flags.CargoScoopDeployed != 0
	f.SilentRunning = raw&flags.SilentRunning != 0
	f.ScoopingFuel = raw&flags.ScoopingFuel != 0
	f.SRVHandbrake = raw&flags.SRVHandbrake != 0
	f.SRVTurret = raw&flags.SRVTurret != 0
	f.SRVUnderShip = raw&flags.SRVUnderShip != 0
	f.SRVDriveAssist = raw&flags.SRVDriveAssist != 0
	f.FSDMassLocked = raw&flags.FSDMassLocked != 0
	f.FSDCharging = raw&flags.FSDCharging != 0
	f.FSDCooldown = raw&flags.FSDCooldown != 0
	f.LowFuel = raw&flags.LowFuel != 0
	f.Overheating = raw&flags.Overheating != 0
	f.HasLatLong = raw&flags.HasLatLong != 0
	f.IsInDanger = raw&flags.IsInDanger != 0
	f.BeingInterdicted = raw&flags.BeingInterdicted != 0
	f.InMainShip = raw&flags.InMainShip != 0
	f.InFighter = raw&flags.InFighter != 0
	f.InSRV = raw&flags.InSRV != 0
	f.InAnalysisMode = raw&flags.InAnalysisMode != 0
	f.NightVision = raw&flags.NightVision != 0
	f.AltitudeFromAverageRadius = raw&flags.AltitudeFromAverageRadius != 0
	f.FSDJump = raw&flags.FSDJump != 0
	f.SRVHighBeam = raw&flags.SRVHighBeam != 0
}

func getStatusFromPath(logPath string) (StatusEvent, error) {
	status := StatusEvent{}
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
func GetStatus(content []byte, status *StatusEvent) error {
	if err := json.Unmarshal(content, status); err != nil {
		return errors.New("couldn't unmarshal Status.json file: " + err.Error())
	}
	status.ExpandFlags()
	return nil
}
