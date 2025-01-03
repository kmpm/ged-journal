package jsonmap

// HasKeys checks if a map has all the keys in a list.
// Returns a list of missing keys and a boolean indicating if all keys are present.
func HasKeys(m map[string]interface{}, k []string) ([]string, bool) {
	hasAll := true
	out := []string{}
	for _, key := range k {
		if _, ok := m[key]; !ok {
			hasAll = false
			out = append(out, key)
		}
	}
	return out, hasAll
}
