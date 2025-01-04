package state

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_SystemUnmarshall(t *testing.T) {
	type args struct {
		data []byte
	}
	empty := map[string]interface{}{}
	tests := []struct {
		name       string
		args       args
		want       System
		wantResult map[string]interface{}
		wantErr    bool
	}{
		{"empty", args{[]byte(`{}`)}, System{}, empty, false},
		{"starsystem", args{[]byte(`{"StarSystem":"Sol"}`)}, System{StarSystem: "Sol"}, empty, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotLoc System
			gotRes, err := UnmarsmallowJSON(tt.args.data, &gotLoc)
			if (err != nil) != tt.wantErr {
				t.Errorf("Location.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.want, gotLoc, "Location is not equal")
			assert.Equal(t, tt.wantResult, gotRes, "Result is not equal")

		})
	}
}

func Test_LocationUnmarshall(t *testing.T) {
	type args struct {
		data []byte
	}
	empty := map[string]interface{}{}
	tests := []struct {
		name       string
		args       args
		want       Location
		wantResult map[string]interface{}
		wantErr    bool
	}{
		{"empty", args{[]byte(`{}`)}, Location{}, empty, false},
		{"starsystem", args{[]byte(`{"StarSystem":"Sol"}`)}, Location{System: System{StarSystem: "Sol"}}, empty, false},
		{"body", args{[]byte(`{"Body":"Earth"}`)}, Location{Body: Body{Body: "Earth"}}, empty, false},
		{"station", args{[]byte(`{"StationName":"Jameson Memorial"}`)}, Location{Station: Station{StationName: "Jameson Memorial"}}, empty, false},
		{"system with pos", args{[]byte(`{"StarSystem":"Sol","StarPos":[1.2,3.4,-5.6]}`)}, Location{System: System{StarSystem: "Sol", StarPos: StarPos{1.2, 3.4, -5.6}}}, empty, false},
		{"lon/latt", args{[]byte(`{"Longitude":1.2,"Latitude":3.4}`)}, Location{Longitude: 1.2, Latitude: 3.4}, empty, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotLoc Location
			gotRes, err := UnmarsmallowJSON(tt.args.data, &gotLoc)
			if (err != nil) != tt.wantErr {
				t.Errorf("Location.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.want, gotLoc, "Location is not equal")
			assert.Equal(t, tt.wantResult, gotRes, "Result is not equal")

		})
	}
}

func Test_SystemMarshal(t *testing.T) {
	t.Skip("will not work until go 1.24 (feb 2025)")
	type args struct {
		sys System
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"empty", args{System{}}, `{}`, false},
		{"starsystem", args{System{StarSystem: "Sol"}}, `{"StarSystem":"Sol"}`, false},
		{"system with pos", args{System{StarSystem: "Sol", StarPos: StarPos{1.2, 3.4, -5.6}}}, `{"StarSystem":"Sol","StarPos":[1.2,3.4,-5.6]}`, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := json.Marshal(tt.args.sys)
			if (err != nil) != tt.wantErr {
				t.Errorf("System.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, string(got), "System is not equal")
		})
	}
}

func Test_LocationMarshal(t *testing.T) {
	t.Skip("will not work until go 1.24 (feb 2025)")
	type args struct {
		loc Location
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"empty", args{Location{}}, `{}`, false},
		{"starsystem", args{Location{System: System{StarSystem: "Sol"}}}, `{"StarSystem":"Sol"}`, false},
		{"body", args{Location{Body: Body{Body: "Earth"}}}, `{"Body":"Earth"}`, false},
		{"station", args{Location{Station: Station{StationName: "Jameson Memorial"}}}, `{"StationName":"Jameson Memorial"}`, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := json.Marshal(tt.args.loc)
			if (err != nil) != tt.wantErr {
				t.Errorf("Location.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, string(got), "Location is not equal")
		})
	}
}

func TestLocation_Update(t *testing.T) {
	type args struct {
		v Location
	}
	tests := []struct {
		name    string
		before  *Location
		args    args
		want    *Location
		wantErr bool
	}{
		// {"nil", nil, args{Location{}}, nil, false},
		{"System",
			&Location{},
			args{Location{System: System{StarSystem: "foo"}}},
			&Location{changed: true, System: System{StarSystem: "foo"}},
			false},
		{"Change System",
			&Location{System: System{StarSystem: "foo"}},
			args{Location{System: System{StarSystem: "bar"}}},
			&Location{changed: true, System: System{StarSystem: "bar"}},
			false},
		{"Change System clears rest",
			&Location{System: System{StarSystem: "foo"}, Body: Body{Body: "earth"}, Station: Station{StationName: "Jameson Memorial"}, Longitude: 1.2, Latitude: 3.4, Heading: 5, Altitude: 6},
			args{Location{System: System{StarSystem: "bar"}}},
			&Location{changed: true, System: System{StarSystem: "bar"}},
			false},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.before.Update(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("Location.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
			assert.Equal(t, tt.want, tt.before)
		})
	}
}
