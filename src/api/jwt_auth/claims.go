package jwt_auth

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"i69social/src/api/models/profile"

	"time"
)

var signingKey []byte

var keyFn = func(*jwt.Token) (interface{}, error) {
	return signingKey, nil
}

var parser = jwt.Parser{
	ValidMethods: []string{"HS256"},
}

func Setup(key string) {
	signingKey = []byte(key)
}

var ErrNotValidToken = errors.New("given access token is not valid")

func ParseToken(accessToken string) (claims *UserClaims, err error) {
	token, err := parser.ParseWithClaims(accessToken, &UserClaims{}, keyFn)
	if err != nil {
		return
	}

	if !token.Valid {
		err = ErrNotValidToken
		return
	}
	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		err = ErrNotValidToken
		return
	}
	return
}

type Privileges map[string]bool

func (p Privileges) Has(privilege string) bool {
	for privilegeName, privilegeValue := range p {
		if privilegeName == privilege && privilegeValue {
			return true
		}
	}

	return false
}

type UserClaims struct {
	jwt.StandardClaims
	Id         uint       `json:"id"`
	Email      string     `json:"email"`
	LastName   string     `json:"last_name"`
	FirstName  string     `json:"first_name"`
	Privileges Privileges `json:"privileges"`
}

func (c UserClaims) Valid() error {
	err := c.StandardClaims.Valid()
	if err != nil {
		return err
	}

	if c.Id == 0 {
		err = fmt.Errorf("user id cannot be empty")
		return err
	}
	return nil
}

func GetUserClaims(p profile.Profile) jwt.Claims {
	return &UserClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 10).Unix(),
		},
		Id:         p.Id,
		Email:      p.Email,
		LastName:   p.LastName,
		FirstName:  p.FirstName,
		Privileges: p.GetPrivileges(),
	}
}

func GetToken(claims jwt.Claims) (tokenString string, err error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err = token.SignedString(signingKey)
	if err != nil {
		return
	}

	return
}
