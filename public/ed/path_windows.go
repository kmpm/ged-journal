package ed

import (
	"os"
	"path/filepath"
)

func SteamAppPath(p PathID) (string, error) {
	return filepath.Join(os.Getenv("ProgramFiles"), "Steam", "steamapps", "common"), nil
}
