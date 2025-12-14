package ed

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
)

const appID = "359320"

type PathID int

const (
	LocalAppDataPath PathID = iota
	RoamingAppDataPath
	PicturesPath
	SavedGamesPath
)

func (p PathID) String() string {
	switch p {
	case LocalAppDataPath:
		return "LocalAppDataPath"
	case RoamingAppDataPath:
		return "RoamingAppDataPath"
	case PicturesPath:
		return "PicturesPath"
	case SavedGamesPath:
		return "SavedGamesPath"
	default:
		return fmt.Sprintf("UnknownPathID: %d", p)
	}
}

func pathExists(path string) bool {
	stat, err := os.Stat(path)
	return err == nil && stat.IsDir()
}

func AppPath(p PathID) (string, error) {
	hp, homeErr := HomePath(p)
	sp, steamErr := SteamAppPath(p)
	if homeErr != nil && steamErr != nil {
		return "", fmt.Errorf("failed to get path: %v, %v", homeErr, steamErr)
	}
	if pathExists(hp) {
		return hp, nil
	}
	if pathExists(sp) {
		return sp, nil
	}
	return "", fmt.Errorf("failed to find path")
}

func HomePath(p PathID) (string, error) {

	currUser, _ := user.Current()
	homeDir := currUser.HomeDir

	if p == SavedGamesPath {
		return filepath.FromSlash(homeDir + "/Saved Games/Frontier Developments/Elite Dangerous"), nil
	}
	return "", fmt.Errorf("invalid path ID")
}
