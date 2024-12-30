package fileapi

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GetStatusFromPath(t *testing.T) {
	type args struct {
		logPath string
	}
	tests := []struct {
		name    string
		args    args
		want    *Status
		wantErr bool
	}{
		{"Test 1", args{"../../testdata/set1/"}, &Status{Event: "Status", Timestamp: "2024-12-29T23:22:17Z"}, false},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetStatusFromPath(tt.args.logPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("getStatusFromPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.want, got)
			// if !reflect.DeepEqual(got, tt.want) {
			// 	t.Errorf("getStatusFromPath() = %+v, want %v", got, tt.want)
			// }
		})
	}
}
