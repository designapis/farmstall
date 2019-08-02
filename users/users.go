package users

import (
	"farmstall/passwords"
	"fmt"
	"github.com/google/uuid"
	"math/rand"
	"time"
)

const EDENIEDBYCODE = "e4033"
const ENOTFOUND = "e4040"

type UUID = string

type UserMap map[UUID]User

type User struct {
	Uuid     UUID   `json:"uuid"`
	Username string `json:"username"`
	FullName string `json:"userId"`
}

type Users struct {
	Users     map[UUID]User `json:"users"`
	Passwords *passwords.PasswordStore
	Tokens    map[UUID]string
}

type NewUser struct {
	Uuid     UUID   `json:"uuid"`
	Username string `json:"username"`
	FullName string `json:"userId"`
	Password string `json:"password"`
}

type UserLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserError struct {
	Msg  string `json:"error"`
	Uuid UUID   `json:"uuid"`
	Code string `json:"code"`
}

func (us *Users) CreateToken(ul UserLogin) (string, error) {
	user, foundUser := us.GetUserByUsername(ul.Username)
	if !foundUser {
		return "", UserError{Msg: "User not found"}
	}

	_, verifyErr := us.Passwords.Verify(user.Uuid, ul.Password)

	if verifyErr != nil {
		return "", verifyErr
	}

	token := RandomString(10)
	us.Tokens[user.Uuid] = token

	return token, nil
}

func (us *Users) GetUserByUsername(username string) (*User, bool) {
	var user *User
	for _, testUser := range us.Users {
		if testUser.Username == username {
			user = &testUser
			break
		}
	}

	if user == nil {
		return nil, false
	}

	return user, true
}

func (e UserError) Error() string {
	return fmt.Sprintf("%s: %s %s", e.Uuid, e.Msg, e.Code)
}

func (us *Users) AddUser(nu NewUser) (*User, error) {
	uuidVal := uuid.New().String()
	u := User{
		FullName: nu.FullName,
		Username: nu.Username,
		Uuid:     uuidVal,
	}

	us.Users[uuidVal] = u
	us.Passwords.Add(uuidVal, nu.Password)

	return &u, nil
}

func (us *Users) GetUser(id UUID) (*User, error) {
	var user User
	var ok bool
	user, ok = us.Users[id]
	if !ok {
		reErr := &UserError{Uuid: id, Msg: "User not found", Code: ENOTFOUND}
		return nil, reErr
	}
	return &user, nil
}

func (us *Users) DeleteUser(id UUID) error {
	if _, ok := us.Users[id]; !ok {
		reErr := UserError{Uuid: id, Msg: "User not found", Code: ENOTFOUND}
		return reErr
	}
	delete(us.Users, id)
	return nil
}

func (us *Users) GetUsers() (*[]User, error) {
	v := make([]User, 0, len(us.Users))
	for _, value := range us.Users {
		v = append(v, value)
	}
	return &v, nil
}

func NewUsers() *Users {
	us := Users{
		Users:     UserMap{},
		Passwords: passwords.NewPasswordStore(),
		Tokens:    make(map[UUID]string),
	}
	return &us
}

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func RandomStringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func RandomString(length int) string {
	return RandomStringWithCharset(length, charset)
}
