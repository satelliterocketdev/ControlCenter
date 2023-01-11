package services

import (
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	"i69social/src/api/firebase_admin"
	"i69social/src/api/models/fake_user"
)

type FakeUser struct {
	repo *fake_user.Repository
}

func GetFakeUserService(db *gorm.DB) *FakeUser {
	return &FakeUser{repo: fake_user.GetRepository(db)}
}

func (s *FakeUser) CreateFakeUser(ownerId uint, data json.RawMessage) (err error) {
	details, err := firebase_admin.CreateFakeUser(data)
	if err != nil {
		return
	}

	id, name, err := s.getFakeUserDetailsToSave(details)
	if err != nil {
		return
	}

	newFakeUser := fake_user.User{
		Name:       name,
		ExternalId: id,
		OwnerId:    ownerId,
	}

	_, err = s.repo.Create(newFakeUser)
	return
}

const (
	NameKey = "name"
	IDKey   = "id"
)

func (s *FakeUser) getFakeUserDetailsToSave(details map[string]interface{}) (id, name string, err error) {
	idVal, ok := details[IDKey]
	if !ok {
		err = fmt.Errorf("key %s is not present in user details", IDKey)
		return
	}

	id, ok = idVal.(string)
	if !ok {
		err = fmt.Errorf("value of key %s is not valid string: given %T", IDKey, idVal)
		return
	}

	nameVal, ok := details[NameKey]
	if !ok {
		err = fmt.Errorf("key %s is not present in user details", NameKey)
		return
	}

	name, ok = nameVal.(string)
	if !ok {
		err = fmt.Errorf("value of key %s is not valid string: given %T", NameKey, idVal)
		return
	}

	return
}

func (s *FakeUser) UpdateFakeUser(uid string, data json.RawMessage) (err error) {
	details, err := firebase_admin.UpdateFakeUser(uid, data)
	if err != nil {
		return
	}

	_, name, err := s.getFakeUserDetailsToSave(details)
	if err != nil {
		return
	}


	err = s.repo.Update(uid, name)
	return
}
