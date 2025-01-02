package messages

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	fileheader1 = `{ "timestamp":"2025-01-01T18:36:31Z", "event":"Fileheader", "part":1, "language":"English/UK", "Odyssey":true, "gameversion":"4.0.0.1904", "build":"r308767/r0 " }`
	commander1  = `{ "timestamp":"2025-01-01T18:37:15Z", "event":"Commander", "FID":"F9999", "Name":"Jameson" }`
)

func TestGetFileHeader(t *testing.T) {
	type args struct {
		data   []byte
		header *FileHeader
	}
	tests := []struct {
		name    string
		args    args
		want    *FileHeader
		wantErr bool
	}{
		{"first", args{[]byte(fileheader1), &FileHeader{}}, &FileHeader{Event: Event{Timestamp: "2025-01-01T18:36:31Z", Event: "Fileheader"}, Part: 1, Language: "English/UK", Gameversion: "4.0.0.1904", Build: "r308767/r0 ", Odyssey: true}, false},
		{"empty", args{[]byte(""), &FileHeader{}}, &FileHeader{}, true},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := GetFileHeader(tt.args.data, tt.args.header)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFileHeader() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.want, tt.args.header)

		})
	}
}

func TestGetCommander(t *testing.T) {
	type args struct {
		data []byte
		cmd  *Commander
	}
	tests := []struct {
		name    string
		args    args
		want    *Commander
		wantErr bool
	}{
		{"first", args{[]byte(commander1), &Commander{}}, &Commander{Event: Event{Timestamp: "2025-01-01T18:37:15Z", Event: "Commander"}, FID: "F9999", Name: "Jameson"}, false},
		{"empty", args{[]byte(""), &Commander{}}, &Commander{}, true},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := GetCommander(tt.args.data, tt.args.cmd)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCommander() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, tt.args.cmd)
		})
	}
}
