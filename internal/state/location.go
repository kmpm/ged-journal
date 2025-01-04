package state

type StarPos [3]float32

type System struct {
	StarSystem    string  `json:"StarSystem,omitzero"`
	SystemAddress int     `json:"SystemAddress,omitzero"`
	StarPos       StarPos `json:"StarPos,omitempty"`
	StarClass     string  `json:"StarClass,omitzero"`
}
type Body struct {
	Body     string `json:"Body,omitzero"`
	BodyID   int    `json:"BodyID,omitzero"`
	BodyType string `json:"BodyType,omitzero"`
}

type Station struct {
	StationName string `json:"StationName,omitzero"`
	StationType string `json:"StationType,omitzero"`
	MarketID    int    `json:"MarketID,omitzero"`
}

type Location struct {
	changed bool
	System
	Body
	Station
	Longitude float32 `json:"Longitude,omitzero"`
	Latitude  float32 `json:"Latitude,omitzero"`
	Heading   int     `json:"Heading,omitzero"`
	Altitude  int     `json:"Altitude,omitzero"`
}

// Update updates the location with the values from v if they are different
// from the current values.
func (l *Location) Update(v Location) error {

	if l == nil {
		panic("nil location")
	}
	err := l.UpdateSystem(v.System)
	if err != nil {
		return err
	}
	err = l.UpdateBody(v.Body)
	if err != nil {
		return err
	}
	err = l.UpdateStation(v.Station)
	if err != nil {
		return err
	}
	if l.Longitude != v.Longitude {
		l.Longitude = v.Longitude
		l.changed = true
	}
	if l.Latitude != v.Latitude {
		l.Latitude = v.Latitude
		l.changed = true
	}
	if l.Heading != v.Heading {
		l.Heading = v.Heading
		l.changed = true
	}
	if l.Altitude != v.Altitude {
		l.Altitude = v.Altitude
		l.changed = true
	}
	return nil
}

func (l *Location) UpdateSystem(v System) error {
	if l.StarSystem != v.StarSystem {
		l.StarSystem = v.StarSystem
		l.changed = true
		// if system changes then clear the rest of the location
		l.UpdateBody(Body{})
	}
	if l.SystemAddress != v.SystemAddress {
		l.SystemAddress = v.SystemAddress
		l.changed = true
	}
	if l.StarPos != v.StarPos {
		l.StarPos = v.StarPos
		l.changed = true
	}
	if l.StarClass != v.StarClass {
		l.StarClass = v.StarClass
		l.changed = true
	}
	return nil
}

func (l *Location) UpdateBody(v Body) error {

	if l == nil {
		panic("nil location")
	}

	if l.Body.Body != v.Body {
		l.Body = v
		l.changed = true
		// if body changes then clear the rest of the location
		l.UpdateStation(Station{})
		l.Longitude = 0
		l.Latitude = 0
		l.Heading = 0
		l.Altitude = 0
	}

	return nil
}

func (l *Location) UpdateStation(v Station) error {

	if l == nil {
		panic("nil location")
	}

	if l.Station.StationName != v.StationName {
		l.StationName = v.StationName
		l.changed = true
	}
	if l.Station.MarketID != v.MarketID {
		l.Station.MarketID = v.MarketID
		l.changed = true
	}
	if l.Station.StationType != v.StationType {
		l.Station.StationType = v.StationType
		l.changed = true
	}

	return nil
}
