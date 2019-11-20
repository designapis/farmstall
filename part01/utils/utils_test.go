package utils

import (
	"testing"

	"encoding/json"
	"gotest.tools/assert"
	is "gotest.tools/assert/cmp"
)

func TestMarshalNullStringToString(t *testing.T) {
	var ns NullString
	ns = "hello"

	jsonBytes, _ := json.Marshal(ns)
	assert.Assert(t, is.Equal(string(jsonBytes), "\"hello\""), "should serialize to a string")
}

func TestUnmarshalNullStringFromString(t *testing.T) {
	jsonStr := `"hello"`
	var ns NullString

	json.Unmarshal([]byte(jsonStr), &ns)
	assert.Assert(t, is.Equal(string(ns), "hello"), "should deserialize from string")
}

func TestMarshalNullStringToNull(t *testing.T) {
	var ns NullString
	ns = ""

	jsonBytes, _ := json.Marshal(ns)
	assert.Assert(t, is.Equal(string(jsonBytes), "null"), "should serialize to a string")
}

func TestUnmarshalNullStringFromNull(t *testing.T) {
	jsonStr := `null`
	var ns NullString

	json.Unmarshal([]byte(jsonStr), &ns)
	assert.Assert(t, is.Equal(string(ns), ""), "should deserialize from null")
}
