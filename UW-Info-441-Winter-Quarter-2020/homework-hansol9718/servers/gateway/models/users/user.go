package users

import (
	"net/mail"
	"fmt"
	"strings"
	"crypto/md5"
	"encoding/hex"
	"golang.org/x/crypto/bcrypt"
	"errors"
)
//gravatarBasePhotoURL is the base URL for Gravatar image requests.
//See https://id.gravatar.com/site/implement/images/ for details
const gravatarBasePhotoURL = "https://www.gravatar.com/avatar/"

//bcryptCost is the default bcrypt cost to use when hashing passwords
var bcryptCost = 13

//User represents a user account in the database
type User struct {
	ID        int64  `json:"id"`
	Email     string `json:"-"` //never JSON encoded/decoded
	PassHash  []byte `json:"-"` //never JSON encoded/decoded
	UserName  string `json:"userName"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	PhotoURL  string `json:"photoURL"`
}

//Credentials represents user sign-in credentials
type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

//NewUser represents a new user signing up for an account
type NewUser struct {
	Email        string `json:"email"`
	Password     string `json:"password"`
	PasswordConf string `json:"passwordConf"`
	UserName     string `json:"userName"`
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
}

//Updates represents allowed updates to a user profile
type Updates struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

//Validate validates the new user and returns an error if
//any of the validation rules fail, or nil if its valid
func (nu *NewUser) Validate() error {
	
	_, err := mail.ParseAddress(nu.Email)
	if err != nil {
		return fmt.Errorf("Error with email: %v", err)
	}
	if len(nu.Password) < 6 {
		return fmt.Errorf("Password must be at least 6 characters")
	}
	if nu.Password != nu.PasswordConf {
		return fmt.Errorf("Passwords do not match")
	}
	if len(nu.UserName) == 0 {
		return fmt.Errorf("UserName cannot be empty")
	}
	if strings.Contains(nu.UserName, " ") {
		return fmt.Errorf("UserName must not contain spaces")
	}
	return nil
}

//ToUser converts the NewUser to a User, setting the
//PhotoURL and PassHash fields appropriately
func (nu *NewUser) ToUser() (*User, error) {

	err := nu.Validate()
	if err != nil {
		return nil, err
	}

	email := strings.TrimSpace(strings.ToLower(nu.Email))
	h := md5.New()
	h.Write([]byte(email))
	hash := h.Sum(nil)
	

	user := &User {
		ID: 0,
		Email: nu.Email,
		UserName: nu.UserName,
		FirstName: nu.FirstName,
		LastName: nu.LastName,
		PhotoURL: gravatarBasePhotoURL + hex.EncodeToString(hash),
	}


	err = user.SetPassword(nu.Password)
	if err != nil  {
		return nil, err
	}

	return user, nil
}

//FullName returns the user's full name, in the form:
// "<FirstName> <LastName>"
//If either first or last name is an empty string, no
//space is put between the names. If both are missing,
//this returns an empty string
func (u *User) FullName() string {
	//TODO: implement according to comment above
	if u.FirstName != "" && u.LastName != "" {
		return u.FirstName + " " + u.LastName
	}
	return strings.TrimSpace(u.FirstName + " " + u.LastName)
}

//SetPassword hashes the password and stores it in the PassHash field
func (u *User) SetPassword(password string) error {
	//TODO: use the bcrypt package to generate a new hash of the password
	//https://godoc.org/golang.org/x/crypto/bcrypt
	h, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)

	if err != nil {
		return fmt.Errorf("Error with hashing password: %v", err)
	}
	u.PassHash = h
	return nil
}

//Authenticate compares the plaintext password against the stored hash
//and returns an error if they don't match, or nil if they do
func (u *User) Authenticate(password string) error {
	//TODO: use the bcrypt package to compare the supplied
	//password with the stored PassHash
	//https://godoc.org/golang.org/x/crypto/bcrypt
	err := bcrypt.CompareHashAndPassword(u.PassHash, []byte(password))
	if err != nil {
		return fmt.Errorf("Hash and password don't match: %v", err)
	}
	return nil
}

//ApplyUpdates applies the updates to the user. An error
//is returned if the updates are invalid
func (u *User) ApplyUpdates(updates *Updates) error {
	if updates.FirstName == "" || updates.LastName == "" {
		return errors.New("FirstName or LastName not provided")
	}
	u.FirstName = updates.FirstName
	u.LastName = updates.LastName

	return nil
}
