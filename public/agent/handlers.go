package agent

import (
	"log/slog"

	"github.com/kmpm/ged-journal/internal/state"
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
	return func(e messages.Event, fields map[string]interface{}, s *state.State) (bool, error) {
		return chainHandlers(e, fields, s, handlers...)
	}
}

func chainHandlers(e messages.Event, fields map[string]interface{}, s *state.State, handlers ...EventHandler) (bool, error) {
	c := false
	for _, h := range handlers {
		hc, err := h(e, fields, s)
		if err != nil {
			return c, err
		}
		c = c || hc
	}
	return c, nil
}

func onDockHandler(e messages.Event, fields map[string]interface{}, s *state.State) (bool, error) {
	// slog.Info("onDockHandler", "event", e.Event, "fields", fields)
	c := false
	var err error
	var hc bool
	if e.Event == "Docked" {
		s.Docked = true
		hc, err = chainHandlers(e, fields, s, onSystemHandler, onBodyHandler, onStationHandler)
		c = c || hc
	} else if e.Event == "Undocked" {
		s.Docked = false
		s.Location.UpdateStation(state.Station{})
		c = true
	} else {
		slog.Warn("onDockHandler with unknown event", "event", e.Event, "fields", fields)
	}
	return c, err
}

func onStatusHandler(_ messages.Event, fields map[string]interface{}, s *state.State) (bool, error) {
	prevFlags := s.RawFlags
	se := messages.StatusEvent{}
	err := messages.GetFromJSONMap(fields, &se)
	if err != nil {
		return false, err
	}
	if prevFlags != se.RawFlags {
		// slog.Info("Status flags changed", "old", prevFlags, "new", state.Status.RawFlags)
		s.RawFlags = se.RawFlags
	}
	return false, nil
}

func onGamehandler(e messages.Event, fields map[string]interface{}, s *state.State) (bool, error) {
	// Example event handler logic
	s.Game.Version = fields["gameversion"].(string)
	s.Game.Build = fields["build"].(string)
	s.Game.Language = fields["language"].(string)
	s.Game.Horizons = getBool(fields, "Horizons")
	s.Game.Odyssey = getBool(fields, "Odyssey")

	return true, nil
}

func onCommanderHandler(e messages.Event, fields map[string]interface{}, s *state.State) (bool, error) {
	// slog.Info("onCommanderHandler", "event", e.Event)

	p := state.Player{}
	if v, ok := fields["Name"]; ok {
		p.Cmdr = v.(string)
	}
	if v, ok := fields["Commander"]; ok {
		p.Cmdr = v.(string)
	}
	if v, ok := fields["FID"]; ok {
		p.FID = v.(string)
	}
	s.UpdatePlayer(p)
	return true, nil
}

func onSystemHandler(e messages.Event, fields map[string]interface{}, s *state.State) (bool, error) {
	// slog.Info("onSystemHandler", "event", e.Event)
	c := false
	if _, ok := fields["StarSystem"]; ok {
		l := state.Location{}
		if err := messages.GetFromJSONMap(fields, &l); err != nil {
			return false, err
		}
		s.Location.Update(l)

	} else if !ok {
		slog.Warn("onSystemHandler without StarSystem", "event", e.Event, "fields", fields)
	}

	return c, nil
}

func onStationHandler(e messages.Event, fields map[string]interface{}, s *state.State) (bool, error) {
	// slog.Info("onStationHandler", "event", e.Event)
	if _, ok := fields["StationName"]; ok {
		st := state.Station{}
		if err := messages.GetFromJSONMap(fields, &st); err != nil {
			return false, err
		}
		s.Location.UpdateStation(st)
		return true, nil
	}
	return false, nil
}

func onBodyHandler(e messages.Event, fields map[string]interface{}, s *state.State) (bool, error) {
	// slog.Info("onBodyHandler", "event", e.Event)
	c := false
	if _, ok := fields["Body"]; ok {
		bo := state.Body{}
		if err := messages.GetFromJSONMap(fields, &bo); err != nil {
			return false, err
		}
		s.Location.UpdateBody(bo)
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

// func getIntNumber(in any) int {
// 	switch in.(type) {
// 	case int:
// 		return in.(int)
// 	case float64:
// 		return int(in.(float64))
// 	default:
// 		return 0
// 	}
// }

// func getInt64Number(in any) int64 {
// 	switch in.(type) {
// 	case int:
// 		return int64(in.(int))
// 	case float64:
// 		return int64(in.(float64))
// 	default:
// 		return 0
// 	}
// }
