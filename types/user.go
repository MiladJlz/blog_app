package types

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"regexp"
)

const (
	bcryptCost      = 12
	minFirstNameLen = 2
	minLastNameLen  = 2
	minPasswordLen  = 7
)

type UpdateUserParams struct {
	FirstName string `json:"firstName" example:"foo"`
	LastName  string `json:"lastName" example:"baz"`
	FcmToken  string `json:"fcmToken"`
	Password  string `json:"password" example:"verysecurepassword"`
}

func (p UpdateUserParams) ToBSON() bson.M {
	m := bson.M{}
	if len(p.FirstName) > 0 {
		m["firstName"] = p.FirstName
	}
	if len(p.LastName) > 0 {
		m["lastName"] = p.LastName
	}
	if len(p.Password) > 0 {
		encpw, _ := bcrypt.GenerateFromPassword([]byte(p.Password), bcryptCost)

		m["password"] = encpw
	}
	return m
}

type CreateUserParams struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	FCMToken  string `json:"fcmToken"`
	Password  string `json:"password"`
}

func (params CreateUserParams) Validate() map[string]string {
	errors := map[string]string{}
	if len(params.FirstName) < minFirstNameLen {
		errors["firstName"] = fmt.Sprintf("firstName length should be at least %d characters", minFirstNameLen)
	}
	if len(params.LastName) < minLastNameLen {
		errors["lastName"] = fmt.Sprintf("lastName length should be at least %d characters", minLastNameLen)
	}
	if len(params.Password) < minPasswordLen {
		errors["password"] = fmt.Sprintf("password length should be at least %d characters", minPasswordLen)
	}
	if !isEmailValid(params.Email) {
		errors["email"] = fmt.Sprintf("email %s is invalid", params.Email)
	}
	return errors
}

func isEmailValid(e string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(e)
}

type User struct {
	ID        primitive.ObjectID   `bson:"_id,omitempty" json:"id,omitempty" example:"66db2c856699531daa9abc16"`
	FirstName string               `bson:"firstName" json:"firstName" example:"foo"`
	LastName  string               `bson:"lastName" json:"lastName" example:"bar"`
	Email     string               `bson:"email" json:"email" example:"foobar@gmail.com"`
	Password  string               `bson:"password" json:"password" example:"verysecurepassword"`
	FCMToken  string               `bson:"fcmToken" json:"fcMToken"`
	Friends   []primitive.ObjectID `bson:"friends" json:"friends" example:"[66db2c856699531daa9abc16,9bdb2c85156699531daa9abc7]"`
}

func NewUserFromParams(params CreateUserParams) (*User, error) {
	encpw, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcryptCost)
	if err != nil {
		return nil, err
	}
	return &User{
		FirstName: params.FirstName,
		LastName:  params.LastName,
		Email:     params.Email,
		Password:  string(encpw),
		FCMToken:  params.FCMToken,
		Friends:   []primitive.ObjectID{},
	}, nil
}

type AddFriendParam struct {
	UserID string `json:"userID"`
}
