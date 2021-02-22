package users
import (
	"crypto/md5"
	"encoding/hex"
	"strings"
	"testing"
	"golang.org/x/crypto/bcrypt"
)
func TestValidate(t *testing.T) {
	cases := []struct {
		name          string
		nu            *NewUser
		expectErr     bool
		expectedErr   string
	}{
		{
			"Valid New User",
			&NewUser{
				Email:        "test@uw.edu",
				Password:     "123456",
				PasswordConf: "123456",
				UserName:     "testusername",
				FirstName:    "first",
				LastName:     "last",
			},
			false,
			"",
		},
		{
			"Invalid Email",
			&NewUser{
				Email:        "test!uw.edu",
				Password:     "123456",
				PasswordConf: "123456",
				UserName:     "testusername",
				FirstName:    "first",
				LastName:     "last",
			},
			true,
			"should return error for invalid email address",
		},
		{
			"Password Too Short",
			&NewUser{
				Email:        "test@uw.edu",
				Password:     "123",
				PasswordConf: "123",
				UserName:     "testusername",
				FirstName:    "first",
				LastName:     "last",
			},
			true,
			"should return error for password shorter than 6 characters",
		},
		{
			"Passwords Don't Match",
			&NewUser{
				Email:        "test@uw.edu",
				Password:     "123456",
				PasswordConf: "123457",
				UserName:     "testusername",
				FirstName:    "first",
				LastName:     "last",
			},
			true,
			"should return error if password and password confirmation don't match",
		},
		{
			"Empty User Name",
			&NewUser{
				Email:        "test@uw.edu",
				Password:     "123456",
				PasswordConf: "123456",
				UserName:     "",
				FirstName:    "first",
				LastName:     "last",
			},
			true,
			"should return error for empty username",
		},
		{
			"User Name with Spaces",
			&NewUser{
				Email:        "test@uw.edu",
				Password:     "123456",
				PasswordConf: "123456",
				UserName:     "user name",
				FirstName:    "first",
				LastName:     "last",
			},
			true,
			"should return error for username that contains spaces",
		},
	}
	for _, c := range cases {
		err := c.nu.Validate()
		switch {
		case c.expectErr && err == nil:
			t.Errorf("case %s: expected error: %s, but did not get any error", c.name, c.expectedErr)
		case !c.expectErr && err != nil:
			t.Errorf("case %s: unexpected error: %s", c.name, err)
		}
	}
} 

func TestToUser(t *testing.T) {
	cases := []struct {
		name          string
		nu            *NewUser
		expectErr     bool
		expectedErr   string
	}{
		{
			"Valid Email",
			&NewUser{
				Email:        "test@uw.edu",
				Password:     "123456",
				PasswordConf: "123456",
				UserName:     "testusername",
				FirstName:    "first",
				LastName:     "last",
			},
			false,
			"",
		},
		{
			"Valid Email with Space and Upper Case",
			&NewUser{
				Email:        "tEstT@uw.edu ",
				Password:     "PAssWord@@",
				PasswordConf: "PAssWord@@",
				UserName:     "testusername",
				FirstName:    "first",
				LastName:     "last",
			},
			false,
			"",
		},
	}
	for _, c := range cases {
		u, err := c.nu.ToUser()
		switch {
		case !c.expectErr && err != nil:
			t.Errorf("case %s: unexpected error: %s", c.name, err)
		case !c.expectErr && err == nil:
			email := strings.ToLower(strings.TrimSpace(c.nu.Email))
			h := md5.New()
			h.Write([]byte(email))
			he := h.Sum(nil)
			expectedURL := gravatarBasePhotoURL + hex.EncodeToString(he)
			if expectedURL != u.PhotoURL {
				t.Errorf("case %s: PhotoURLs don't match: expected %s, but got %s", c.name, expectedURL, u.PhotoURL)
			}
			err = bcrypt.CompareHashAndPassword(u.PassHash, []byte(c.nu.Password))
			if err != nil {
				t.Errorf("case %s: error comparing pass hash and password", c.name)
			}
		}
	}
}
func TestFullName(t *testing.T) {
	cases := []struct {
		name         string
		u            *User
		expectedName string
	}{
		{
			"Has Both First and Last Name",
			&User{
				FirstName: "Hansol",
				LastName:  "Trapp",
			},
			"Hansol Trapp",
		},
		{
			"Only First Name",
			&User{
				FirstName: "Hansol",
				LastName: "",
			},
			"Hansol",
		},
		{
			"Only Last Name",
			&User{
				FirstName: "",
				LastName: "Trapp",
			},
			"Trapp",
		},
		{
			"No Name Provided",
			&User{
				FirstName: "",
				LastName: "",
			},
			"",
		},
	}
	for _, c := range cases {
		name := c.u.FullName()
		if c.expectedName != name {
			t.Errorf("case %s: name: %s does not match expected name: %s", c.name, name, c.expectedName)
		}
	}
}


func TestAuthenticate(t *testing.T) {
	cases := []struct {
		name          string
		nu            *NewUser
		test       	  string
		expectErr     bool
		expectedErr   string
	}{
		{
			"Successful Authenticate",
			&NewUser{
				Email:        "test@uw.edu",
				Password:     "123456",
				PasswordConf: "123456",
				UserName:     "testusername",
				FirstName:    "first",
				LastName:     "last",
			},
			"123456",
			false,
			"",
		},
		{
			"Different Pass Hash",
			&NewUser{
				Email:        "test@uw.edu",
				Password:     "123456",
				PasswordConf: "123456",
				UserName:     "testusername",
				FirstName:    "first",
				LastName:     "last",
			},
			"123457",
			true,
			"pass hash is different from password",
		},
		{
			"Empty Password",
			&NewUser{
				Email:        "test@uw.edu",
				Password:     "123456",
				PasswordConf: "123456",
				UserName:     "testusername",
				FirstName:    "first",
				LastName:     "last",
			},
			"",
			true,
			"passhash is empty",
		},
	}
	for _, c := range cases {
		u, err := c.nu.ToUser()
		if err != nil {
			t.Errorf("case %s: unexpected error %s", c.name, err)
		}
		err = u.Authenticate(c.test)
		switch {
		case !c.expectErr && err != nil:
			t.Errorf("case %s: unexpected error: %s", c.name, err)
		case c.expectErr && err == nil:
			t.Errorf("case %s: expected error: %s", c.name, err)
		}
	}
}

func TestApplyUpdates(t *testing.T) {
	cases := []struct {
		name          string
		u            *User
		updates		 *Updates
		expectErr     bool
		expectedErr   string
		expectedFirst string
		expectedLast  string
	}{
		{
			"Update Both First and Last Name",
			&User {
				FirstName: "Hansol",
				LastName: "Kim",
			},
			&Updates {
				FirstName: "Caleb",
				LastName: "Trapp",
			},
			false,
			"",
			"Caleb",
			"Trapp",
		},
		{
			"Update Last Name",
			&User {
				FirstName: "Hansol",
				LastName: "Kim",
			},
			&Updates {
				FirstName: "Hansol",
				LastName: "Trapp",
			},
			false,
			"",
			"Hansol",
			"Trapp",
		},
		{
			"Empty Update",
			&User {
				FirstName: "Hansol",
				LastName: "Kim",
			},
			&Updates {
				FirstName: "",
				LastName: "",
			},
			true,
			"First Name and Last Name not provided",
			"Hansol",
			"Kim",
		},
	}
	for _, c := range cases {
		err := c.u.ApplyUpdates(c.updates)
		switch {
		case !c.expectErr && err != nil:
			t.Errorf("case %s: unexpected error updating name: %s", c.name, err)
		case c.expectErr && err == nil:
			t.Errorf("case %s: expected error: %s", c.name, c.expectedErr)
		case c.expectedFirst != c.u.FirstName || c.expectedLast != c.u.LastName:
			t.Errorf("case %s: expected error: %s Expected name: %s %s but returned name: %s %s ",
			 		c.name, c.expectedErr, c.expectedLast, c.expectedFirst, c.u.FirstName, c.u.LastName )
		}
	}
}
