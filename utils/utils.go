package utils

import (
	"encoding/json"
	_ "log"
)

type NullString string

func (s NullString) MarshalJSON() ([]byte, error) {
	if s == "" {
		return json.Marshal(nil)
	}
	return json.Marshal(string(s))
}

func (s *NullString) UnmarshalJSON(data []byte) error {
	var sPointer *string

	if err := json.Unmarshal(data, &sPointer); err != nil {
		return err
	}
	if sPointer == nil {
		*s = ""
	} else {
		*s = NullString(*sPointer)
	}
	return nil
}
