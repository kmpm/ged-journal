package messages

import (
	"github.com/perimeterx/marshmallow"
)

type Event struct {
	Timestamp string `json:"timestamp"`
	Event     string `json:"event"`
}

func init() {
	marshmallow.EnableCache()
}

// GetEventComponent unmarshals a byte slice into an Event struct and returns all the data in a JSONMap.
func GetEventComponent(data []byte, event *Event) (map[string]any, error) {
	result, err := marshmallow.Unmarshal(data, event)
	return result, err
}

func PopFromData(data []byte, v any) (map[string]interface{}, error) {
	return marshmallow.Unmarshal(data, v, marshmallow.WithExcludeKnownFieldsFromMap(true))
}

// PopFromJSONMap unmarshals a map[string]interface{} into a struct and returns whats left.
func PopFromJSONMap(data map[string]interface{}, v any) (map[string]interface{}, error) {
	return marshmallow.UnmarshalFromJSONMap(data, v, marshmallow.WithExcludeKnownFieldsFromMap(true))
}

func GetFromJSONMap(data map[string]interface{}, v any) error {
	_, err := marshmallow.UnmarshalFromJSONMap(data, v)
	return err
}

type FileHeaderEvent struct {
	Event
	Part        int    `json:"part"`
	Language    string `json:"language"`
	Gameversion string `json:"gameversion"`
	Build       string `json:"build"`
	Odyssey     bool   `json:"Odyssey,omitempty"`
	Horizons    bool   `json:"Horizons,omitempty"`
}

type CommanderEvent struct {
	Event
	FID  string `json:"FID"`
	Name string `json:"Name"`
}

type StarPos []float32

type System struct {
	StarSystem    string  `json:"StarSystem"`
	SystemAddress int     `json:"SystemAddress"`
	StarPos       StarPos `json:"StarPos,omitempty"`
	StarClass     string  `json:"StarClass,omitempty"`
}
type Body struct {
	Body     string `json:"Body"`
	BodyID   int    `json:"BodyID"`
	BodyType string `json:"BodyType"`
}

type Station struct {
	StationName string `json:"StationName"`
	StationType string `json:"StationType"`
	MarketID    int64  `json:"MarketID,omitempty"`
}

type StateEntry struct {
	State string `json:"State"`
	Trend int    `json:"Trend"`
}

type Faction struct {
	Name                string       `json:"Name"`
	FactionState        string       `json:"FactionState"`
	Government          string       `json:"Government,omitempty"`
	Influence           float32      `json:"Influence,omitempty"`
	Allegiance          string       `json:"Allegiance,omitempty"`
	Happiness           string       `json:"Happiness,omitempty"`
	Happiness_Localised string       `json:"Happiness_Localised,omitempty"`
	MyReputation        float32      `json:"MyReputation,omitempty"`
	PendingStates       []StateEntry `json:"PendingStates,omitempty"`
	ActiveStates        []StateEntry `json:"ActiveStates,omitempty"`
}

type RouteStep struct {
	System
	StarClass string `json:"StarClass"`
}

type NavRouteEvent struct {
	Event
	Route []RouteStep `json:"Route"`
}

type ApproachBodyEvent struct {
	Event
	System
	Body
}

type ApproachSettlementEvent struct {
	Event
	Name                        string         `json:"Name"`
	MarketID                    int            `json:"MarketID"`
	StationFaction              Faction        `json:"StationFaction"`
	StationGovernment           string         `json:"StationGovernment"`
	StationGovernment_Localised string         `json:"StationGovernment_Localised"`
	StationAllegiance           string         `json:"StationAllegiance"`
	StationServices             []string       `json:"StationServices"`
	StationEconomy              string         `json:"StationEconomy"`
	StationEconomy_Localised    string         `json:"StationEconomy_Localised"`
	StationEconomies            []EconomyEntry `json:"StationEconomies"`
	SystemAddress               int64          `json:"SystemAddress"`
	BodyID                      int            `json:"BodyID"`
	BodyName                    string         `json:"BodyName"`
	Longitude                   float32        `json:"Longitude"`
	Latitude                    float32        `json:"Latitude"`
}

type FSDJumpEvent struct {
	Event
	System
	Body
	Factions                      []Faction `json:"Factions"`
	FuelLevel                     float32   `json:"FuelLevel"`
	FuelUsed                      float32   `json:"FuelUsed"`
	JumpDist                      float32   `json:"JumpDist"`
	Multicrew                     bool      `json:"Multicrew,omitempty"`
	Population                    float32   `json:"Population,omitempty"`
	SystemAlliance                string    `json:"SystemAllegiance,omitempty"`
	SystemEconomy                 string    `json:"SystemEconomy,omitempty"`
	SystemEconomy_Localised       string    `json:"SystemEconomy_Localised,omitempty"`
	SystemFaction                 Faction   `json:"SystemFaction,omitempty"`
	SystemGovernment              string    `json:"SystemGovernment,omitempty"`
	SystemGovernment_Localised    string    `json:"SystemGovernment_Localised,omitempty"`
	SystemSecondEconomy           string    `json:"SystemSecondEconomy,omitempty"`
	SystemSecondEconomy_Localised string    `json:"SystemSecondEconomy_Localised,omitempty"`
	SystemSecurity                string    `json:"SystemSecurity,omitempty"`
	SystemSecurity_Localised      string    `json:"SystemSecurity_Localised,omitempty"`
	Taxi                          bool      `json:"Taxi,omitempty"`
}

type DockingDeniedEvent struct {
	Event
	Station
	Reason string `json:"Reason"`
}

type UndockedEvent struct {
	Event
	System
	Station
	Taxi      bool `json:"Taxi,omitempty"`
	Multicrew bool `json:"Multicrew,omitempty"`
	MarketID  int  `json:"MarketID,omitempty"`
}

type SupercruiseExitEvent struct {
	Event
	System
	Body
	Taxi      bool `json:"Taxi,omitempty"`
	Multicrew bool `json:"Multicrew,omitempty"`
}

type SupercruiseEntryEvent struct {
	Event
	System
	Taxi      bool `json:"Taxi,omitempty"`
	Multicrew bool `json:"Multicrew,omitempty"`
}

type SupercruiseDestinationDropEvent struct {
	Event
	Type     string `json:"Type"`
	Threat   int64  `json:"Threat"`
	MarketID int64  `json:"MarketID"`
}

type StartJumpEvent struct {
	Event
	System
	Taxi     bool   `json:"Taxi,omitempty"`
	JumpType string `json:"JumpType"`
}

type ShipTargetedEvent struct {
	Event
	TargetLocked       bool    `json:"TargetLocked"`
	Ship               string  `json:"Ship"`
	ShipLocalised      string  `json:"Ship_Localised"`
	ScanStage          int     `json:"ScanStage"`
	PilotName          string  `json:"PilotName"`
	PilotNameLocalised string  `json:"PilotName_Localised"`
	PilotRank          string  `json:"PilotRank"`
	ShieldHealth       float32 `json:"ShieldHealth,omitempty"`
	HullHealth         float32 `json:"HullHealth,omitempty"`
	Faction            string  `json:"Faction,omitempty"`
	LegalStatus        string  `json:"LegalStatus,omitempty"`
}

type ShipLockerEvent struct {
	Event
	Items       []any `json:"Items"`
	Components  []any `json:"Components"`
	Consumables []any `json:"Consumables"`
	Data        []any `json:"Data"`
}

type ReputationEvent struct {
	Event
	Empire      float32 `json:"Empire"`
	Federation  float32 `json:"Federation"`
	Independent float32 `json:"Independent"`
	Alliance    float32 `json:"Alliance"`
}

type RepairEvent struct {
	Event
	Items []string `json:"Items"`
	Cost  int      `json:"Cost"`
}

type RefuelAllEvent struct {
	Event
	Cost   int     `json:"Cost"`
	Amount float32 `json:"Amount"`
}

type ReceiveTextEvent struct {
	Event
	Channel          string `json:"Channel"`
	From             string `json:"From"`
	FromLocalised    string `json:"From_Localised"`
	Message          string `json:"Message"`
	MessageLocalised string `json:"Message_Localised"`
}

type Engineer struct {
	Engineer     string `json:"Engineer"`
	EngineerID   int    `json:"EngineerID"`
	RankProgress int    `json:"RankProgress,omitempty"`
	Progress     string `json:"Progress"`
	Rank         int    `json:"Rank"`
}

type EngineerProgressEvent struct {
	Event
	Engineers []Engineer `json:"Engineers"`
}

type RankEvent struct {
	Event
	Combat       int `json:"Combat"`
	Trade        int `json:"Trade"`
	Explore      int `json:"Explore"`
	Soldier      int `json:"Soldier"`
	Exobiologist int `json:"Exobiologist"`
	Empire       int `json:"Empire"`
	Federation   int `json:"Federation"`
	CQC          int `json:"CQC"`
}

type PromotionEvent struct {
	Event
	Trade   int `json:"Trade,omitempty"`
	Explore int `json:"Explore,omitempty"`
	CQC     int `json:"CQC,omitempty"`
}

type MissionRedirectedEvent struct {
	Event
	MissionID             int64  `json:"MissionID"`
	Name                  string `json:"Name"`
	LocalisedName         string `json:"LocalisedName"`
	NewDestinationStation string `json:"NewDestinationStation"`
	NewDestinationSystem  string `json:"NewDestinationSystem"`
	OldDestinationStation string `json:"OldDestinationStation"`
	OldDestinationSystem  string `json:"OldDestinationSystem"`
}

type MaterialEntry struct {
	Name          string `json:"Name"`
	NameLocalised string `json:"Name_Localised"`
	Count         int    `json:"Count"`
}

type MaterialsEvent struct {
	Event
	Raw          []MaterialEntry `json:"Raw"`
	Manufactured []MaterialEntry `json:"Manufactured"`
	Encoded      []MaterialEntry `json:"Encoded"`
}

type EconomyEntry struct {
	Name          string  `json:"Name"`
	NameLocalised string  `json:"Name_Localised"`
	Proportion    float32 `json:"Proportion"`
}

type LocationEvent struct {
	Event
	Station
	System
	Body
	DistFromStarLS               float32        `json:"DistFromStarLS"`
	Docked                       bool           `json:"Docked"`
	StateFaction                 Faction        `json:"StateFaction"`
	StationGoverment             string         `json:"StationGovernment"`
	StationGovermentLocalised    string         `json:"StationGovernment_Localised"`
	StationAllegiance            string         `json:"StationAllegiance"`
	StationServices              []string       `json:"StationServices"`
	StationEconomy               string         `json:"StationEconomy"`
	StationEconomyLocalised      string         `json:"StationEconomy_Localised"`
	StationEconomies             []EconomyEntry `json:"StationEconomies"`
	Taxi                         bool           `json:"Taxi"`
	Multicrew                    bool           `json:"Multicrew"`
	SystemAllegiance             string         `json:"SystemAllegiance"`
	SystemEconomy                string         `json:"SystemEconomy"`
	SystemEconomyLocalised       string         `json:"SystemEconomy_Localised"`
	SystemSecondEconomy          string         `json:"SystemSecondEconomy"`
	SystemSecondEconomyLocalised string         `json:"SystemSecondEconomy_Localised"`
	SystemGovernment             string         `json:"SystemGovernment"`
	SystemGovernmentLocalised    string         `json:"SystemGovernment_Localised"`
	SystemSecurity               string         `json:"SystemSecurity"`
	SystemSecurityLocalised      string         `json:"SystemSecurity_Localised"`
	Population                   int            `json:"Population"`
	Factions                     []Faction      `json:"Factions"`
	SystemFaction                Faction        `json:"SystemFaction"`
	StationFaction               Faction        `json:"StationFaction"`
}

type ModuleEntry struct {
	Slot         string  `json:"Slot"`
	Item         string  `json:"Item"`
	On           bool    `json:"On"`
	Priority     int     `json:"Priority"`
	Health       float32 `json:"Health"`
	AmmoInClip   int     `json:"AmmoInClip,omitempty"`
	AmmoInHopper int     `json:"AmmoInHopper,omitempty"`
}

type LoadoutEvent struct {
	Event
	Ship          string  `json:"Ship"`
	ShipID        int     `json:"ShipID"`
	ShipName      string  `json:"ShipName"`
	ShipIdent     string  `json:"ShipIdent"`
	HullHealth    float32 `json:"HullHealth"`
	UnladenMass   float32 `json:"UnladenMass"`
	CargoCapacity int     `json:"CargoCapacity"`
	MaxJumpRange  float32 `json:"MaxJumpRange"`
	FuelCapacity  struct {
		Main      float32 `json:"Main"`
		Reservoir float32 `json:"Reservoir"`
	} `json:"FuelCapacity"`
	Rebuy   int           `json:"Rebuy"`
	Modules []ModuleEntry `json:"Modules"`
}

type InventoryEntry struct {
	Name          string `json:"Name"`
	NameLocalised string `json:"Name_Localised"`
	Count         int    `json:"Count"`
	Stolen        int    `json:"Stolen"`
}

type CargoEvent struct {
	Event
	Vessel    string           `json:"Vessel"`
	Count     int              `json:"Count"`
	Inventory []InventoryEntry `json:"Inventory"`
}

type DockedEvent struct {
	Event
	Station
	System
	Multicrew                 bool           `json:"Multicrew,omitempty"`
	Taxi                      bool           `json:"Taxi,omitempty"`
	MarketID                  int            `json:"MarketID"`
	StationFaction            Faction        `json:"StationFaction"`
	StationGoverment          string         `json:"StationGovernment"`
	StationGovermentLocalised string         `json:"StationGovernment_Localised"`
	StationAllegiance         string         `json:"StationAllegiance"`
	StationServices           []string       `json:"StationServices"`
	StationEconomy            string         `json:"StationEconomy"`
	StationEconomyLocalised   string         `json:"StationEconomy_Localised"`
	StationEconomies          []EconomyEntry `json:"StationEconomies"`
	DistFromStarLS            float32        `json:"DistFromStarLS"`
	LandingPads               struct {
		Small  int `json:"Small"`
		Medium int `json:"Medium"`
		Large  int `json:"Large"`
	} `json:"LandingPads"`
}
