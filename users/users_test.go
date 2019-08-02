package users

import (
	"testing"

	// "github.com/google/uuid"
	"gotest.tools/assert"
	is "gotest.tools/assert/cmp"
)

func TestGetUsersEmpty(t *testing.T) {
	users := NewUsers()
	res, err := users.GetUsers()
	assert.NilError(t, err, "should have no errors")
	assert.Assert(t, is.Len(*res, 0), "should return empty list of users")
}

func TestAddOneUser(t *testing.T) {
	us := NewUsers()
	_, err := us.AddUser(NewUser{
		Username: "ponelat",
		FullName: "Josh Ponelat",
		Password: "password",
	})
	assert.NilError(t, err, "should have no errors")
	assert.Assert(t, is.Len(us.Users, 1), "should update the list of users")
}

func TestAddThenGetOneUser(t *testing.T) {
	us := NewUsers()
	addedUser, err := us.AddUser(NewUser{
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

func TestGetTokenFromUser(t *testing.T) {
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

// func TestDeleteOnUser(t *testing.T) {
// 	users := NewUsers()
// 	addedUser, _ := users.AddUser(User{
// 		Message: "good",
// 		Rating:  5,
// 	})

// 	err := users.DeleteUser(addedUser.Uuid)
// 	assert.NilError(t, err, "should have no errors")
// 	assert.Assert(t, is.Len(users.Users, 0), "should have no users")
// }

// func TestUpdateUser(t *testing.T) {
// 	users := NewUsers()
// 	oriUser, _ := users.AddUser(User{
// 		Message: "poor",
// 		Rating:  1,
// 	})

// 	newUser := User{
// 		Message: "good",
// 		Rating:  5,
// 	}

// 	_, err := users.UpdateUser(oriUser.Uuid, newUser)
// 	assert.NilError(t, err, "should have no errors")

// 	updatedUser, err := users.GetUser(oriUser.Uuid)
// 	assert.NilError(t, err, "should have no errors")

// 	assert.Assert(t, is.Equal(updatedUser.Message, "good"), "should match the updated user's Message")
// 	assert.Assert(t, is.Equal(updatedUser.Rating, 5), "should match the updated user's Rating")
// 	assert.Assert(t, is.Len(users.Users, 1), "should not add any more users")
// }

// func TestUpdateUserInvalidUuid(t *testing.T) {
// 	users := NewUsers()
// 	newUser := User{
// 		Message: "good",
// 		Rating:  5,
// 		Uuid:    uuid.New().String(),
// 	}
// 	_, err := users.UpdateUser(newUser.Uuid, newUser)
// 	assert.ErrorContains(t, err, "User does not exist", "should return an error")
// }

// func TestGetAllUsers(t *testing.T) {
// 	users := NewUsers()
// 	users.AddUser(User{
// 		Message: "good",
// 		Rating:  5,
// 	})
// 	users.AddUser(User{
// 		Message: "average",
// 		Rating:  3,
// 	})
// 	users.AddUser(User{
// 		Message: "poor",
// 		Rating:  1,
// 	})

// 	allUsers, err := users.GetUsers()
// 	assert.NilError(t, err, "should have no errors")

// 	assert.Assert(t, is.Len(*allUsers, 3), "should equal the number of users added")
// }

// func TestGetUsersByMaxRatingBelowOrEquals(t *testing.T) {
// 	users := NewUsers()
// 	users.AddUser(User{
// 		Message: "good",
// 		Rating:  5,
// 	})
// 	users.AddUser(User{
// 		Message: "average",
// 		Rating:  3,
// 	})
// 	users.AddUser(User{
// 		Message: "poor",
// 		Rating:  1,
// 	})

// 	filters := UserFilters{
// 		MaxRating: 3,
// 	}

// 	allUsers, err := users.GetUsersFiltered(filters)

// 	assert.NilError(t, err, "should have no errors")

// 	assert.Assert(t, is.Len(*allUsers, 2), "should equal two, for the users with ratings 1 and 3")
// }

// func TestGetUsersByMaxRatingEquals(t *testing.T) {
// 	users := NewUsers()
// 	users.AddUser(User{
// 		Message: "good",
// 		Rating:  5,
// 	})
// 	users.AddUser(User{
// 		Message: "average",
// 		Rating:  3,
// 	})
// 	users.AddUser(User{
// 		Message: "poor",
// 		Rating:  1,
// 	})

// 	filters := UserFilters{
// 		MaxRating: 1,
// 	}

// 	allUsers, err := users.GetUsersFiltered(filters)

// 	assert.NilError(t, err, "should have no errors")

// 	assert.Assert(t, is.Len(*allUsers, 1), "should equal one, for the user with rating 1")
// }
