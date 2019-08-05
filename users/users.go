package users

import (
	"farmstall/passwords"
	"farmstall/problems"
	"fmt"
	"github.com/google/uuid"
	"math/rand"
	"time"
)

type UUID = string
type UserMap map[UUID]User

const BASE_PATH = "/users"

type User struct {
	Uuid     UUID   `json:"uuid"`
	Username string `json:"username"`
	FullName string `json:"fullName"`
}

type Users struct {
	Users     map[UUID]User `json:"users"`
	Passwords *passwords.PasswordStore
	Tokens    map[UUID]string
}

type NewUser struct {
	Uuid     UUID   `json:"uuid"`
	Username string `json:"username"`
	FullName string `json:"fullName"`
	Password string `json:"password"`
}

type UserLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type TokenResponse struct {
	Token string `json:"token"`
}

func (us *Users) CreateToken(ul UserLogin) (string, error) {
	user, userErr := us.GetFirstUserByUsername(ul.Username)
	if userErr != nil {
		return "", userErr
	}

	_, verifyErr := us.Passwords.Verify(user.Uuid, ul.Password)

	if verifyErr != nil {
		return "", problems.InvalidCreds(problems.ProblemJson{
			Detail: "Username or password is invalid",
		})
	}

	token := RandomString(10)
	us.Tokens[user.Uuid] = token

	return token, nil
}

func (us *Users) GetFirstUserByUsername(username string) (*User, error) {
	var user *User
	for _, testUser := range us.Users {
		if testUser.Username == username {
			user = &testUser
			break
		}
	}

	if user == nil {
		return nil, problems.NotFound(problems.ProblemJson{
			Detail: fmt.Sprintf("No user with username, %s, found", username),
		})
	}

	return user, nil
}

func (us *Users) AddUser(nu NewUser) *User {
	uuidVal := uuid.New().String()
	u := User{
		FullName: nu.FullName,
		Username: nu.Username,
		Uuid:     uuidVal,
	}

	us.Users[uuidVal] = u
	us.Passwords.Add(uuidVal, nu.Password)

	return &u
}

func (us *Users) GetUser(id UUID) (*User, error) {
	var user User
	var ok bool
	user, ok = us.Users[id]
	if !ok {
		return nil, problems.NotFound(problems.ProblemJson{
			Instance: BASE_PATH + "/" + id,
			Detail:   fmt.Sprintf("User with uuid, %s, does not exist.", id),
		})
	}
	return &user, nil
}

func (us *Users) DeleteUser(id UUID) error {
	if _, ok := us.Users[id]; !ok {
		return problems.NotFound(problems.ProblemJson{
			Instance: BASE_PATH + "/" + id,
			Detail:   fmt.Sprintf("User with uuid, %s, does not exist.", id),
		})
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
