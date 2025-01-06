package messages

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/kmpm/ged-journal/internal/jsonmap"
	"github.com/stretchr/testify/assert"
)

const (
	fileheader1   = `{ "timestamp":"2025-01-01T18:36:31Z", "event":"Fileheader", "part":1, "language":"English/UK", "Odyssey":true, "gameversion":"4.0.0.1904", "build":"r308767/r0 " }`
	commander1    = `{ "timestamp":"2025-01-01T18:37:15Z", "event":"Commander", "FID":"F9999", "Name":"Jameson" }`
	navroute1     = `{"Route":[{"StarClass":"K","StarPos":[27.4375,80.375,-62.40625],"StarSystem":"Ross 94","SystemAddress":358797611722},{"StarClass":"M","StarPos":[24.1875,80.5,-60.09375],"StarSystem":"Djakam","SystemAddress":2870246450577}],"event":"NavRoute","timestamp":"2025-01-01T19:15:31Z"}`
	approachbody1 = `{ "timestamp":"2025-01-01T19:29:31Z", "event":"ApproachBody", "StarSystem":"Ross 94", "SystemAddress":358797611722, "Body":"Ross 94 1 b", "BodyID":15 }`
)

func TestFromBytes(t *testing.T) {
	type args struct {
		data  []byte
		event any
	}
	tests := []struct {
		name    string
		args    args
		want    any
		wantErr bool
	}{
		{"approachbody", args{[]byte(approachbody1), &ApproachBodyEvent{}}, &ApproachBodyEvent{Event: Event{Timestamp: "2025-01-01T19:29:31Z", Event: "ApproachBody"}, System: System{StarSystem: "Ross 94", SystemAddress: 358797611722, StarPos: StarPos(nil)}, Body: Body{Body: "Ross 94 1 b", BodyID: 15, BodyType: ""}}, false},
		{"fileheader", args{[]byte(fileheader1), &FileHeaderEvent{}}, &FileHeaderEvent{Event: Event{Timestamp: "2025-01-01T18:36:31Z", Event: "Fileheader"}, Part: 1, Language: "English/UK", Gameversion: "4.0.0.1904", Build: "r308767/r0 ", Odyssey: true}, false},
		{"navroute", args{[]byte(navroute1), &NavRouteEvent{}}, &NavRouteEvent{Event: Event{Timestamp: "2025-01-01T19:15:31Z", Event: "NavRoute"}, Route: []RouteStep{{StarClass: "K", System: System{StarSystem: "Ross 94", SystemAddress: 358797611722, StarPos: StarPos{27.4375, 80.375, -62.40625}}}, {StarClass: "M", System: System{StarSystem: "Djakam", SystemAddress: 2870246450577, StarPos: StarPos{24.1875, 80.5, -60.09375}}}}}, false},
		// {"empty", args{[]byte{}, &Event{}}, true},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PopFromData(tt.args.data, tt.args.event)
			if (err != nil) != tt.wantErr {
				t.Errorf("PopFromData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, tt.args.event)
			assert.Equal(t, got, map[string]interface{}{}, "should be empty")
		})
	}
}

func assertHas(t *testing.T, result map[string]interface{}, key string) {
	t.Helper()
	if _, ok := result[key]; !ok {
		t.Errorf("missing key %q", key)
	}
}

func fileTest(t *testing.T, path string, wantErr bool) {
	// t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}

	bytesTest(t, data, wantErr)
}

func bytesTest(t *testing.T, data []byte, wantErr bool) {
	// t.Helper()
	var rest map[string]interface{}
	empty := map[string]interface{}{}
	event := Event{}
	result, err := GetEventComponent(data, &event)
	if (err != nil) != wantErr {
		t.Errorf("GetEvent() error = %v, wantErr %v", err, wantErr)
	}
	if err != nil {
		return
	}

	if fields, ok := EventFields[event.Event]; ok {
		e := Event{}
		rest, err = PopFromJSONMap(result, &e)
		assert.NoError(t, err)
		if missing, ok := jsonmap.HasKeys(rest, fields); !ok {
			t.Errorf("missing keys %v", missing)
		}
		return
	}

	switch event.Event {
	case "Fileheader":
		assertHas(t, result, "language")
	case "FSDJump":
		fj := FSDJumpEvent{}
		rest, err = PopFromJSONMap(result, &fj)
		assert.NoError(t, err)
		assert.NotEmpty(t, fj.StarSystem)
		assert.Len(t, fj.StarPos, 3, "StarPos should have 3 elements")
		assert.GreaterOrEqual(t, len(fj.Factions), 1, "Factions should have at least 1 element")
		for _, f := range fj.Factions {
			assert.NotEmpty(t, f.Name)
			assert.NotEmpty(t, f.Government)
			assert.NotEmpty(t, f.Allegiance)
		}
		assert.Equal(t, empty, rest, "rest should be empty")
	case "DockingDenied":
		dd := DockingDeniedEvent{}
		rest, err := PopFromJSONMap(result, &dd)
		assert.NoError(t, err)
		assert.NotEmpty(t, dd.StationName)
		assert.NotEmpty(t, dd.StationType)
		assert.Equal(t, empty, rest, "rest should be empty")
	case "Undocked":
		ud := UndockedEvent{}
		rest, err := PopFromJSONMap(result, &ud)
		assert.NoError(t, err)
		assert.NotEmpty(t, ud.StationName)
		assert.Equal(t, empty, rest, "rest should be empty")
	case "SupercruiseExit":
		se := SupercruiseExitEvent{}
		rest, err := PopFromJSONMap(result, &se)
		assert.NoError(t, err)
		assert.Equal(t, empty, rest, "rest should be empty")
	case "SupercruiseEntry":
		se := SupercruiseEntryEvent{}
		rest, err := PopFromJSONMap(result, &se)
		assert.NoError(t, err)
		assert.Equal(t, empty, rest, "rest should be empty")
	case "SupercruiseDestinationDrop":
		sdd := SupercruiseDestinationDropEvent{}
		rest, err := PopFromJSONMap(result, &sdd)
		assert.NoError(t, err)
		assert.Equal(t, empty, rest, "rest should be empty")
	case "Statistics":
		s := StatisticsEvent{}
		rest, err := PopFromJSONMap(result, &s)
		assert.NoError(t, err)
		assert.Equal(t, empty, rest, "rest should be empty")
	case "StartJump":
		sj := StartJumpEvent{}
		rest, err := PopFromJSONMap(result, &sj)
		assert.NoError(t, err)
		assert.Equal(t, empty, rest, "rest should be empty")

	case "ShipTargeted":
		st := ShipTargetedEvent{}
		rest, err := PopFromJSONMap(result, &st)
		assert.NoError(t, err)
		assert.Empty(t, rest, "rest should be empty")
		// assert.Equal(t, empty, rest, "rest should be empty")
	case "ShipLocker":
		sl := ShipLockerEvent{}
		rest, err := PopFromJSONMap(result, &sl)
		assert.NoError(t, err)
		assert.Empty(t, rest, "rest should be empty")
	case "Reputation":
		r := ReputationEvent{}
		rest, err := PopFromJSONMap(result, &r)
		assert.NoError(t, err)
		assert.Empty(t, rest, "rest should be empty")
	case "Repair":
		r := RepairEvent{}
		rest, err := PopFromJSONMap(result, &r)
		assert.NoError(t, err)
		assert.Empty(t, rest, "rest should be empty")
	case "RefuelAll":
		r := RefuelAllEvent{}
		rest, err := PopFromJSONMap(result, &r)
		assert.NoError(t, err)
		assert.Empty(t, rest, "rest should be empty")
	case "ReceiveText":
		r := ReceiveTextEvent{}
		rest, err := PopFromJSONMap(result, &r)
		assert.NoError(t, err)
		assert.Empty(t, rest, "rest should be empty")
	case "Rank", "Progress":
		r := RankEvent{}
		rest, err := PopFromJSONMap(result, &r)
		assert.NoError(t, err)
		assert.Empty(t, rest, "rest should be empty")
	case "EngineerProgress":
		r := EngineerProgressEvent{}
		rest, err := PopFromJSONMap(result, &r)
		assert.NoError(t, err)
		assert.Empty(t, rest, "rest should be empty")
	case "Promotion":
		r := PromotionEvent{}
		rest, err := PopFromJSONMap(result, &r)
		assert.NoError(t, err)
		assert.Empty(t, rest, "rest should be empty")
	case "NavRoute", "NavRouteClear":
		nr := NavRouteEvent{}
		rest, err := PopFromJSONMap(result, &nr)
		assert.NoError(t, err)
		assert.Empty(t, rest, "rest should be empty")

	case "MissionRedirected":
		e := MissionRedirectedEvent{}
		rest, err := PopFromJSONMap(result, &e)
		assert.NoError(t, err)
		assert.Empty(t, rest, "rest should be empty")
	case "Materials":
		m := MaterialsEvent{}
		rest, err := PopFromJSONMap(result, &m)
		assert.NoError(t, err)
		assert.Empty(t, rest, "rest should be empty")
	case "Location":
		l := LocationEvent{}
		rest, err := PopFromJSONMap(result, &l)
		assert.NoError(t, err)
		assert.Empty(t, rest, "rest should be empty")
	case "Loadout":
		l := LoadoutEvent{}
		rest, err := PopFromJSONMap(result, &l)
		assert.NoError(t, err)
		assert.Empty(t, rest, "rest should be empty")
	case "Cargo":
		c := CargoEvent{}
		rest, err := PopFromJSONMap(result, &c)
		assert.NoError(t, err)
		assert.Empty(t, rest, "rest should be empty")
	case "Status":
		s := StatusEvent{}
		rest, err := PopFromJSONMap(result, &s)
		assert.NoError(t, err)
		assert.Empty(t, rest, "rest should be empty")
	case "Docked":
		d := DockedEvent{}
		rest, err := PopFromJSONMap(result, &d)
		assert.NoError(t, err)
		assert.Empty(t, rest, "rest should be empty")
	case "ApproachSettlement":
		a := ApproachSettlementEvent{}
		rest, err := PopFromJSONMap(result, &a)
		assert.NoError(t, err)
		assert.Empty(t, rest, "rest should be empty")
	case "ApproachBody":
		a := ApproachBodyEvent{}
		rest, err := PopFromJSONMap(result, &a)
		assert.NoError(t, err)
		assert.Empty(t, rest, "rest should be empty")

	default:
		t.Errorf("unknown event %q", event.Event)
	}
}

func Test_GetEvent_Multi(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"ged-collector-sim", args{"../../../output/ged-collector-sim"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list, err := os.ReadDir(tt.args.path)
			if err != nil {
				t.Skipf("skipping test, error reading directory %v", err)
			}

			for _, f := range list {
				if f.IsDir() {
					continue
				}
				t.Run(f.Name(), func(t *testing.T) {
					fileTest(
						t,
						filepath.FromSlash(tt.args.path+"/"+f.Name()),
						tt.wantErr,
					)
				})
			}

		})
	}
}
