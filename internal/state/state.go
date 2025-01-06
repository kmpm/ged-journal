package state

import (
	"reflect"
	"sync"

	"github.com/kmpm/ged-journal/public/messages"
)

type Player struct {
	changed bool
	Cmdr    string `json:"Cmdr"`
	FID     string `json:"FID"`
}

type Game struct {
	Language string `json:"Language"`
	Version  string `json:"Version"`
	Build    string `json:"Build"`
	Horizons bool   `json:"Horizons"`
	Odyssey  bool   `json:"Odyssey"`
	Mode     string `json:"GameMode"`
}
type Ship struct {
	Name string `json:"Name"`
	ID   int    `json:"ID"`
}

type State struct {
	mu        sync.Mutex
	changed   bool
	Timestamp string `json:"Timestamp"`
	Docked    bool   `json:"Docked"`
	Credits   int    `json:"Credits"`

	Location Location `json:"Location"`
	Player   Player   `json:"Player"`
	Game     Game     `json:"Game"`
	Ship     Ship     `json:"Ship"`
	RawFlags uint32   `json:"Flags"`
	Pips     [3]int32 `json:"Pips"`
}

func (s *State) Lock() {
	s.mu.Lock()
}
func (s *State) Unlock() {
	s.mu.Unlock()
}

func (s *State) Flags() (f messages.StatusFlags) {
	messages.ExpandFlags(s.RawFlags, &f)
	return
}

func New() *State {
	return &State{}
}

// func (s *State) MarshalJSON() ([]byte, error) {
// 	s.Mu.Lock()
// 	defer s.Mu.Unlock()

// 	return json.Marshal(&struct {
// 		Timestamp string   `json:"Timestamp"`
// 		Game      Game     `json:"Game"`
// 		Docked    bool     `json:"Docked"`
// 		Location  Location `json:"Location"`
// 		Player    Player   `json:"Player"`
// 	}{
// 		Timestamp: s.timestamp,
// 		Game:      s.game,
// 		Docked:    s.docked,
// 		Location:  s.location,
// 		Player:    s.player,
// 	})
// }

func (s *State) SetLocation(v Location) error {
	s.Lock()
	defer s.Unlock()
	s.Location.Update(v)

	return nil
}

func (s *State) UpdatePlayer(p Player) error {
	s.Lock()
	defer s.Unlock()
	if !reflect.DeepEqual(s.Player, p) {
		p.changed = true
		s.Player = p
	}
	s.changed = s.changed || p.changed
	return nil
}
