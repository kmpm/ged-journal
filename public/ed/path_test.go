package ed_test

import (
	"os"
	"strings"
	"testing"

	"github.com/kmpm/ged-journal/public/ed"
)

func assertPathExists(t *testing.T, path string) bool {
	if _, err := os.Stat(path); err != nil {
		t.Errorf("Path %v stat error: %v", path, err)
		return false
	}
	return true
}

func TestSteamAppPath(t *testing.T) {
	tests := []struct {
		name     string
		pathID   ed.PathID
		contains string
		exists   bool
	}{
		{"SteamAppPath", ed.LocalAppDataPath, "AppData/Local", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ed.SteamAppPath(tt.pathID)
			if err != nil {
				t.Errorf("SteamAppPath() error: %v", err)
				return
			}
			if !strings.Contains(got, tt.contains) {
				t.Errorf("SteamAppPath() = %v, want %v", got, tt.contains)
			}
			//check if the path exists and if it is expected
			if assertPathExists(t, got) != tt.exists {
				t.Errorf("SteamAppPath() = %v, want %v", got, tt.exists)
			}

		})
	}
}
