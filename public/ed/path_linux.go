package ed

import (
	"fmt"
	"os"
	"path/filepath"
)

// https://wiki.archlinux.org/title/XDG_Base_Directory

func SteamAppPath(p PathID) (string, error) {
	var datahome string
	// use XDG_DATA_HOME or fallback to $HOME/.local/share
	if datahome = os.Getenv("XDG_DATA_HOME"); datahome == "" {
		datahome = filepath.Join(os.Getenv("HOME"), ".local", "share")
	}
	basePath := filepath.Join(datahome, "Steam", "steamapps", "compatdata", appID)
	switch p {
	case LocalAppDataPath:
		return filepath.Join(basePath, "pfx", "drive_c", "users", "steamuser", "AppData", "Local", "Frontier Developments", "Elite Dangerous"), nil
	case RoamingAppDataPath:
		return filepath.Join(basePath, "pfx", "drive_c", "users", "steamuser", "AppData", "Roaming", "Frontier Developments", "Elite Dangerous"), nil
	case PicturesPath:
		return filepath.Join(basePath, "pfx", "drive_c", "users", "steamuser", "Pictures", "Frontier Developments", "Elite Dangerous"), nil
	case SavedGamesPath:
		return filepath.Join(basePath, "pfx", "drive_c", "users", "steamuser", "Saved Games", "Frontier Developments", "Elite Dangerous"), nil
	default:
		return "", fmt.Errorf("unsupported path ID: %d", p)
	}
}
