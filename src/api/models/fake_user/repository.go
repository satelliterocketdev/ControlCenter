package fake_user

import (
	"fmt"
	"github.com/jinzhu/gorm"
)

type Repository struct {
	db *gorm.DB
}

func GetRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetById(id uint) (user User, err error) {
	if id == 0 {
		err = fmt.Errorf("id cannot be empty")
		return
	}
	err = r.db.Debug().Where(&User{Id: id}).First(&user).Error
	return
}

func (r *Repository) GetByOwnerId(ownerId uint) (users []User, err error) {
	if ownerId == 0 {
		err = fmt.Errorf("ownerId cannot be empty")
		return
	}
	err = r.db.Debug().Where(&User{OwnerId: ownerId}).Find(&users).Error
	return
}


func (r *Repository) Create(user User) (createdUser User, err error) {
	if user.Id != 0 {
		err = fmt.Errorf("user id should be empty")
		return
	}

	err = r.db.Create(&user).Error
	createdUser = user
	return
}

func (r *Repository) Update(externalId string, name string) (err error) {
	if len(externalId) == 0 {
		err = fmt.Errorf("externalId cannot be empty")
		return
	}

	err = r.db.Model(&User{}).Where(&User{ExternalId:externalId}).Update(User{Name:name}).Error
	return
}

