package data

import (
	"encoding/json"
	"io"
)

// ToJson serializes the given interface into a string based JSON format
func ToJson(i interface{}, w io.Writer) error {
	e := json.NewEncoder(w)

	return e.Encode(i)
}
