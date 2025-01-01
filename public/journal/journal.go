package journal

import (
	"errors"

	"github.com/perimeterx/marshmallow"
)

type JournalEntry struct {
	Timestamp string `json:"timestamp"`
	Event     string `json:"event"`
	// Data      map[string]json.RawMessage
}

type GenericEntry struct {
	Timestamp string `json:"timestamp"`
	Event     string `json:"event"`
	Data      map[string]interface{}
}

func Parse(data []byte) (obj *GenericEntry, err error) {
	var entry JournalEntry
	var result map[string]interface{}
	result, err = marshmallow.Unmarshal(data, &entry)
	if err != nil {
		return nil, err
	}
	if result["event"] == nil || result["event"] == "" {
		return nil, errors.New("missing event field")
	}

	switch result["event"] {
	default:
		return &GenericEntry{
			Timestamp: entry.Timestamp,
			Event:     entry.Event,
			Data:      result,
		}, nil

	}
	// if entry.Event == "" {
	// 	slog.Warn("missing event field", "entry", entry, "result", result)
	// 	for key, v := range result {
	// 		slog.Warn("result key", "key", key, "value", v)

	// 	}
	// 	return nil, errors.New("missing event field")
	// }
	// // // second step: unmarshal the data field
	// // if err = json.Unmarshal(data, &entry.Data); err != nil {
	// // 	return nil, err
	// // }
	// slog.Info("parsed journal entry", "entry", entry, "result", result)

	// return &entry, nil
}
