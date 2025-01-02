package messages

import "encoding/json"

type Event struct {
	Timestamp string `json:"timestamp"`
	Event     string `json:"event"`
	// Data      map[string]json.RawMessage
}

type GenericEvent struct {
	Event
	Message json.RawMessage
}

type FileHeader struct {
	Event
	Part        int    `json:"part"`
	Language    string `json:"language"`
	Gameversion string `json:"gameversion"`
	Build       string `json:"build"`
	Odyssey     bool   `json:"Odyssey,omitempty"`
}

func GetFileHeader(data []byte, header *FileHeader) error {
	err := json.Unmarshal(data, header)
	return err
}

type Commander struct {
	Event
	FID  string `json:"FID"`
	Name string `json:"Name"`
}

func GetCommander(data []byte, cmd *Commander) error {
	err := json.Unmarshal(data, cmd)
	return err
}
