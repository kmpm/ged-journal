package messages

// EventFields is a map of event names to the fields that are expected to be present in the event.
// More complex events have their own struct and are handled separately.
var EventFields map[string][]string = map[string][]string{
	"ApproachBody":         {"Body", "BodyID", "StarSystem", "SystemAddress"},
	"BuyTradeData":         {"System", "Cost"},
	"Commander":            {"FID", "Name"},
	"DockingGranted":       {"LandingPad", "StationName", "StationType", "MarketID"},
	"DockingRequested":     {"StationName", "StationType", "MarketID", "LandingPads"},
	"EscapeInterdiction":   {"Interdictor", "IsPlayer"},
	"Fileheader":           {"part", "language", "gamebersion", "build"},
	"FSDTarget":            {"Name", "SystemAddress", "StarClass", "RemainingJumpsInRoute"},
	"FSSDiscoveryScan":     {"Progress", "BodyCount", "NonBodyCount", "SystemAddress", "SystemName"},
	"FSSSignalDiscovered":  {"SystemAddress", "SignalName", "SignalType"},
	"FuelScoop":            {"Scooped", "Total"},
	"LeaveBody":            {"Body", "BodyID", "StarSystem", "SystemAddress"},
	"LoadGame":             {"FID", "Commander", "Horizons", "Odyssey", "Ship", "Ship_Localised", "ShipID", "ShipName", "ShipIdent", "FuelLevel", "FuelCapacity", "GameMode", "Credits", "Loan", "language", "gameversion", "build"},
	"Market":               {"StationName", "StationType", "MarketID", "StarSystem"},
	"MarketBuy":            {"MarketID", "Type", "Count", "BuyPrice", "TotalCost"},
	"MarketSell":           {"MarketID", "Type", "Count", "SellPrice", "TotalSale", "AvgPricePaid"},
	"MaterialCollected":    {"Category", "Name", "Count", "Name_Localised"},
	"Missions":             {"Active", "Failed", "Complete"},
	"Music":                {"MusicTrack"},
	"ReservoirReplenished": {"FuelMain", "FuelReservoir"},
	"Shutdown":             {},
}
