package agent

import (
	"log/slog"

	"github.com/kmpm/ged-journal/public/messages"
)

func init() {
	RegisterEventHandlers(onChainedHandlers(onSystemHandler, onBodyHandler),
		"ApproachBody",
		"FSDJump",
		"SupercruiseExit")
	RegisterEventHandlers(onChainedHandlers(onSystemHandler, onBodyHandler, onStationHandler), "Location")
	RegisterEventHandlers(onGamehandler, "Fileheader", "LoadGame")
	RegisterEventHandlers(onCommanderHandler, "Commander", "LoadGame")
	RegisterEventHandlers(onStatusHandler, "Status")
	RegisterEventHandlers(onDockHandler, "Docked", "Undocked")
}

func onChainedHandlers(handlers ...EventHandler) EventHandler {
	return func(e messages.Event, fields map[string]interface{}, state *State) (bool, error) {
		return chainHandlers(e, fields, state, handlers...)
	}
}

func chainHandlers(e messages.Event, fields map[string]interface{}, state *State, handlers ...EventHandler) (bool, error) {
	c := false
	for _, h := range handlers {
		hc, err := h(e, fields, state)
		if err != nil {
			return c, err
		}
		c = c || hc
	}
	return c, nil
}

func onDockHandler(e messages.Event, fields map[string]interface{}, state *State) (bool, error) {
	slog.Info("onDockHandler", "event", e.Event, "fields", fields)
	c := false
	var err error
	var hc bool
	if e.Event == "Docked" {
		state.Docked = true
		hc, err = chainHandlers(e, fields, state, onSystemHandler, onBodyHandler, onStationHandler)
		c = c || hc
	} else if e.Event == "Undocked" {
		state.Docked = false
		state.Station = messages.Station{}
		c = true
	} else {
		slog.Warn("onDockHandler with unknown event", "event", e.Event, "fields", fields)
	}
	return c, err
}

func onStatusHandler(_ messages.Event, fields map[string]interface{}, state *State) (bool, error) {
	prevFlags := state.Status.RawFlags
	_, err := messages.PopFromJSONMap(fields, state.Status)
	if err != nil {
		return false, err
	}
	if prevFlags != state.Status.RawFlags {
		// slog.Info("Status flags changed", "old", prevFlags, "new", state.Status.RawFlags)
	}
	return true, nil
}

func onGamehandler(e messages.Event, fields map[string]interface{}, state *State) (bool, error) {
	// Example event handler logic
	state.Game.Version = fields["gameversion"].(string)
	state.Game.Build = fields["build"].(string)
	state.Game.Language = fields["language"].(string)
	state.Game.Horizons = getBool(fields, "Horizons")
	state.Game.Odyssey = getBool(fields, "Odyssey")

	return true, nil
}

func onCommanderHandler(e messages.Event, fields map[string]interface{}, state *State) (bool, error) {
	// slog.Info("onCommanderHandler", "event", e.Event)
	if v, ok := fields["Name"]; ok && v.(string) != state.Cmdr {
		state.Cmdr = v.(string)
	}
	if v, ok := fields["Commander"]; ok && v.(string) != state.Cmdr {
		state.Cmdr = v.(string)
	}
	state.FID = fields["FID"].(string)
	return true, nil
}

func onSystemHandler(e messages.Event, fields map[string]interface{}, state *State) (bool, error) {
	// slog.Info("onSystemHandler", "event", e.Event)
	c := false
	if _, ok := fields["StarSystem"]; ok {
		s := messages.System{}
		err := messages.GetFromJSONMap(fields, &s)
		if err != nil {
			return false, err
		}
		if s.StarSystem != state.System.StarSystem {
			//New system
			slog.Info("StarSystem changed", "old", state.System, "new", s)
			state.System = s
			c = true
		} else {
			//Same system
			if s.SystemAddress != state.System.SystemAddress && s.SystemAddress != 0 {
				slog.Info("System address changed", "old", state.System.SystemAddress, "new", s.SystemAddress)
				state.System.SystemAddress = s.SystemAddress
				c = true
			}
			if len(s.StarPos) > 0 {
				slog.Info("StarPos changed", "old", state.System.StarPos, "new", s.StarPos)
				state.System.StarPos = s.StarPos
				c = true
			}
		}
		// if system changes then body and station must change as well
		state.Station = messages.Station{}
		state.Body = messages.Body{}
		c = true
	} else if !ok {
		slog.Warn("onSystemHandler without StarSystem", "event", e.Event, "fields", fields)
	}

	return c, nil
}

func onStationHandler(e messages.Event, fields map[string]interface{}, state *State) (bool, error) {
	// slog.Info("onStationHandler", "event", e.Event)
	if _, ok := fields["StationName"]; ok {
		s := messages.Station{}
		err := messages.GetFromJSONMap(fields, &s)
		if err != nil {
			return false, err
		}
		if s.StationName != state.Station.StationName {
			slog.Info("Station changed", "old", state.Station, "new", s)
			state.Station = s
			return true, nil
		}
		if s.MarketID != state.Station.MarketID {
			slog.Info("MarketID changed", "old", state.Station.MarketID, "new", s.MarketID)
			state.Station.MarketID = s.MarketID
		}
		if s.StationType != state.Station.StationType {
			slog.Info("StationType changed", "old", state.Station.StationType, "new", s.StationType)
			state.Station.StationType = s.StationType
		}
	} else if !ok {
		slog.Warn("onStationHandler without StationName", "event", e.Event, "fields", fields)
		state.Station.StationName = ""
		state.Station.MarketID = 0
		state.Station.StationType = ""
	}

	return true, nil
}

func onBodyHandler(e messages.Event, fields map[string]interface{}, state *State) (bool, error) {
	// slog.Info("onBodyHandler", "event", e.Event)
	c := false
	if v, ok := fields["Body"]; ok && v.(string) != state.Body.Body {
		state.Body.Body = v.(string)
		c = true
	} else if !ok {
		slog.Warn("onBodyHandler without Body", "event", e.Event, "fields", fields)
		state.Body.Body = ""
	}

	if v, ok := fields["BodyID"]; ok {
		state.Body.BodyID = getIntNumber(v)
		c = true
	}
	if v, ok := fields["BodyType"]; ok && v.(string) != state.Body.BodyType {
		state.Body.BodyType = v.(string)
		c = true
	}
	return c, nil
}

func getBool(fields map[string]interface{}, key string) bool {
	if fields[key] != nil {
		return fields[key].(bool)
	}
	return false
}

func getIntNumber(in any) int {
	switch in.(type) {
	case int:
		return in.(int)
	case float64:
		return int(in.(float64))
	default:
		return 0
	}
}

func getInt64Number(in any) int64 {
	switch in.(type) {
	case int:
		return int64(in.(int))
	case float64:
		return int64(in.(float64))
	default:
		return 0
	}
}
