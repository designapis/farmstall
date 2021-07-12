package users

import (
	"farmstall/passwords"
	"farmstall/problems"
	"fmt"
	"github.com/google/uuid"
	_ "log"
	"math/rand"
	"time"
)

type UserMap map[string]User

const BASE_PATH = "/users"

type User struct {
	Uuid     string `json:"uuid"`
	Username string `json:"username"`
	FullName string `json:"fullName"`
}

type Users struct {
	Users     map[string]User `json:"users"`
	Passwords *passwords.PasswordStore
	Tokens    map[string]string
}

type NewUser struct {
	Uuid     string `json:"uuid"`
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

func (us *Users) CreateToken(ul UserLogin, tokenOverride string) (string, error) {
	user, userErr := us.GetUserByUsername(ul.Username)
	if userErr != nil {
		return "", userErr
	}

	_, verifyErr := us.Passwords.Verify(user.Uuid, ul.Password)

	if verifyErr != nil {
		return "", problems.InvalidCreds(problems.ProblemJson{
			Detail: "Username or password is invalid",
		})
	}

	var token string
	if tokenOverride != "" {
		token = tokenOverride
	} else {
		token = RandomString(10)
	}
	us.Tokens[user.Uuid] = token

	return token, nil
}

func (us *Users) GetUserByUsername(username string) (*User, error) {
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

func (us *Users) AddUser(nu NewUser) (*User, error) {
	existingUser, _ := us.GetUserByUsername(nu.Username)

	if existingUser != nil {
		return nil, problems.CreateAlreadyExists(problems.ProblemJson{
			Instance: BASE_PATH + "/" + nu.Username,
			Detail:   fmt.Sprintf("User with username, %s, already exists.", nu.Username),
		})
	}

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

func (us *Users) GetUser(id string) (*User, error) {
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

func (us *Users) DeleteUser(id string) error {
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

func (us *Users) UserFromToken(token string) (*User, error) {
	for id, tok := range us.Tokens {
		if tok == token {
			user := us.Users[id]
			return &user, nil
		}
	}

	return nil, problems.InvalidCreds(problems.ProblemJson{
		Detail: "Invalid token",
	})
}

func NewUsers() *Users {
	us := Users{
		Users:     UserMap{},
		Passwords: passwords.NewPasswordStore(),
		Tokens:    make(map[string]string),
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
