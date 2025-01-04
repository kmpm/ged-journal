package state

import "github.com/perimeterx/marshmallow"

func UnmarsmallowJSON(data []byte, v interface{}) (map[string]interface{}, error) {
	return marshmallow.Unmarshal(data, v, marshmallow.WithExcludeKnownFieldsFromMap(true))
}
