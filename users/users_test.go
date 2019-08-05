package users

import (
	"gotest.tools/assert"
	is "gotest.tools/assert/cmp"
	"testing"
)

func TestGetUsersEmpty(t *testing.T) {
	users := NewUsers()
	res, err := users.GetUsers()
	assert.NilError(t, err, "should have no errors")
	assert.Assert(t, is.Len(*res, 0), "should return empty list of users")
}

func TestAddOneUser(t *testing.T) {
	us := NewUsers()
	us.AddUser(NewUser{
		Username: "ponelat",
		FullName: "Josh Ponelat",
		Password: "password",
	})
	assert.Assert(t, is.Len(us.Users, 1), "should update the list of users")
}

func TestAddOneUserResponse(t *testing.T) {
	us := NewUsers()
	user := us.AddUser(NewUser{
		Username: "ponelat",
		FullName: "Josh Ponelat",
		Password: "password",
	})
	uuid := user.Uuid
	assert.Assert(t, is.Equal(*user, User{
		FullName: "Josh Ponelat",
		Username: "ponelat",
		Uuid:     uuid,
	}), "should return a User object")
}

func TestAddThenGetOneUser(t *testing.T) {
	us := NewUsers()
	addedUser := us.AddUser(NewUser{
		Username: "ponelat",
		FullName: "Josh Ponelat",
		Password: "password",
	})

	gottenUser, err := us.GetUser(addedUser.Uuid)

	assert.NilError(t, err, "should have no errors")
	assert.Assert(t, is.DeepEqual(addedUser, gottenUser), "should match the user in memory")
}

func TestAddTwoThenGetAllUsers(t *testing.T) {
	users := NewUsers()
	users.AddUser(NewUser{
		Username: "ponelat",
		FullName: "Josh Ponelat",
		Password: "password",
	})
	users.AddUser(NewUser{
		Username: "bgerh",
		FullName: "Bob Gerhard",
		Password: "password",
	})

	allUsers, err := users.GetUsers()

	assert.NilError(t, err, "should have no errors")
	assert.Assert(t, is.Len(*allUsers, 2), "should match the number of users added")
}

func TestCreateTokenFromUsernameAndPassword(t *testing.T) {
	users := NewUsers()
	users.AddUser(NewUser{
		Username: "ponelat",
		FullName: "Josh Ponelat",
		Password: "password",
	})

	token, err := users.CreateToken(UserLogin{
		Username: "ponelat",
		Password: "password",
	})

	assert.NilError(t, err, "should have no errors")
	assert.Assert(t, is.Len(token, 10), "should be a string ten characters long")
}
