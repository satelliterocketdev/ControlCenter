package fake_user

import (
	"crypto/md5"
	"fmt"
	"time"
)

type User struct {
	Id uint `json:"id"`
	Name string `json:"name"`
	ExternalId string `json:"external_id"`
	OwnerId uint `json:"owner_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

func (User) TableName() string {
	return "fake_users"
}

const FakeUserIdPrefix = "x0x_"

func GenerateNewId() string {
	dateBytes, _ := time.Now().MarshalBinary()
	hashBytes := md5.Sum(dateBytes)

	return FakeUserIdPrefix + fmt.Sprintf("%x", hashBytes)[3:]
}
