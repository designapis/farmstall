package passwords

import (
	"testing"

	// "github.com/google/uuid"
	"gotest.tools/assert"
	is "gotest.tools/assert/cmp"
)

func TestAddAndCompareGoodPassword(t *testing.T) {
	uuid := "abc"
	store := NewPasswordStore()
	store.Add(uuid, "password")
	res, err := store.Verify(uuid, "password")
	assert.NilError(t, err, "should have no errors")
	assert.Assert(t, res, 0, "should return true, as passwords match")
}

func TestAddAndCompareBadPassword(t *testing.T) {
	uuid := "abc"
	store := NewPasswordStore()
	store.Add(uuid, "password")
	res, err := store.Verify(uuid, "bad")
	assert.Error(t, err, "Hash comparison failed")
	assert.Assert(t, !res, 0, "should return false, as passwords DO NOT match")
}

func TestAddAndCompareMissingPassword(t *testing.T) {
	uuid := "abc"
	store := NewPasswordStore()
	res, err := store.Verify(uuid, "password")
	assert.Error(t, err, "Password not in system")
	assert.Assert(t, is.Equal(res, false), 0, "should return false, as passwords DO NOT match")
}
