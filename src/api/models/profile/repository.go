package profile

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	"time"
)

type ProfileRepository struct {
	db *gorm.DB
	redisDb *redis.Client
}

const (
	secret  = "secret"
	InvitationExpiration = 24 * time.Hour
)

func GetProfileRepository(db *gorm.DB, redisDb *redis.Client) *ProfileRepository {
	return &ProfileRepository{db: db, redisDb: redisDb}
}

func (r *ProfileRepository) GetAll() (users []Profile, err error) {
	err = r.db.Debug().Omit("password").Find(&users).Error
	return
}

func (r *ProfileRepository) GetAllNotIn(exceptIds []uint) (users []Profile, err error) {
	if len(exceptIds) == 0 {
		err = fmt.Errorf("user ids cannot be empty")
		return
	}
	err = r.db.Debug().Preload("FakeUsers").Omit("password").Where("id not in (?)", exceptIds).Find(&users).Error
	return
}

func (r *ProfileRepository) GetByIds(ids []uint) (users []Profile, err error) {
	if len(ids) == 0 {
		err = fmt.Errorf("user ids cannot be empty")
		return
	}
	err = r.db.Debug().Where("id in (?)", ids).Find(&users).Error
	return
}

func (r *ProfileRepository) GetById(id uint) (user Profile, err error) {
	if id == 0 {
		err = fmt.Errorf("user id cannot be empty")
		return
	}
	err = r.db.Debug().Where(&Profile{Id: id}).First(&user).Error
	return
}

func (r *ProfileRepository) DeleteById(id uint) (err error) {
	if id == 0 {
		err = fmt.Errorf("user id cannot be empty")
		return
	}
	err = r.db.Debug().Where(&Profile{Id: id}).Delete(&Profile{}).Error
	return
}

func (r *ProfileRepository) GetByEmail(email string) (user Profile, err error) {
	if len(email) == 0 {
		err = fmt.Errorf("user email cannot be empty")
		return
	}

	err = r.db.Where(&Profile{Email: email}).First(&user).Error
	return
}


func (r *ProfileRepository) Create(user Profile) (createdUser Profile, err error) {
	if user.Id != 0 {
		err = fmt.Errorf("user id should be empty")
		return
	}

	err = r.db.Exec("INSERT INTO users (email, password, first_name, last_name, roles, is_active) "+
		"VALUES(?, AES_ENCRYPT(?, ?), ?, ?, ?, ?);",
		user.Email, user.Password, secret, user.FirstName, user.LastName, string(user.Roles.RawMessage), true).Error
	if err != nil {
		return
	}

	createdUser, err = r.GetByUserParams(user)
	return
}

func (r *ProfileRepository) GetByUserParams(userParams Profile) (user Profile, err error) {
	if len(userParams.Email) == 0 {
		err = fmt.Errorf("user username cannot be empty")
		return
	}

	err = r.db.Where(&Profile{Email: userParams.Email}).
		Where("AES_DECRYPT(password, ?) = ?", secret, userParams.Password).First(&user).Error
	return
}

func (r *ProfileRepository) StoreInvitationParams(key string, invitationParams GenerateInvitationRequest) (err error) {
	serializedValue, err := json.Marshal(invitationParams)
	if err != nil {
		return
	}

	_, err = r.redisDb.Set(key, string(serializedValue), InvitationExpiration).Result()
	if err != nil {
		return
	}

	return
}

var KeyNotExists = errors.New("key not exists in storage")

func (r *ProfileRepository) GetInvitationParams(key string) (invitationParams GenerateInvitationRequest, err error) {
	value, err := r.redisDb.Get(key).Result()
	if err == redis.Nil {
		err = KeyNotExists
		return
	}
	if err != nil {
		return
	}

	err = json.Unmarshal([]byte(value), &invitationParams)

	return
}

func (r *ProfileRepository) RemoveInvitationKey(key string) (err error) {
	_, err = r.redisDb.Del(key).Result()

	if err != nil {
		return
	}

	return
}

const (
	UsersLastContact = "users_last_contact"
	UsersLastMessage = "users_last_message"
)

func (r *ProfileRepository) StoreWorkerLastContact(key string, userId string) (err error) {
	_, err = r.redisDb.HSet(key, userId, time.Now().UnixNano()).Result()
	if err != nil {
		return
	}

	return
}

func (r *ProfileRepository) GetWorkersLastItem(key string) (m map[string]string, err error) {
	m, err = r.redisDb.HGetAll(key).Result()
	if err != nil {
		return
	}

	return
}

func (r *ProfileRepository) StoreWorkerLastMessage(key string, userId string) (err error) {
	_, err = r.redisDb.HSet(key, userId, time.Now().UnixNano()).Result()
	if err != nil {
		return
	}

	return
}

func (r *ProfileRepository) SaveSentMessages(workerId uint, fakeUserId, chatId string, count uint) (user Profile, err error) {
	nowDate := time.Now().UTC().Format("2006-01-02")
	err = r.db.Exec(`INSERT INTO worker_stats(worker_id, fake_user_id, chat_id, sent_messages, date)
VALUES
    (?, ?, ?, ?, ?)
ON DUPLICATE KEY UPDATE
    sent_messages = sent_messages + ?;`, workerId, fakeUserId, chatId, count, nowDate, count).Error
	return
}

func (r *ProfileRepository) SaveReceivedMessages(workerId uint, fakeUserId, chatId string, count uint) (user Profile, err error) {
	nowDate := time.Now().UTC().Format("2006-01-02")
	err = r.db.Exec(`INSERT INTO worker_stats(worker_id, fake_user_id, chat_id, received_messages, date)
VALUES
    (?, ?, ?, ?, ?)
ON DUPLICATE KEY UPDATE
    received_messages = received_messages + ?;`, workerId, fakeUserId, chatId, count, nowDate, count).Error
	return
}

type WorkerStat struct {
	WorkerId uint `json:"worker_id"`
	Date time.Time `json:"date"`
	SentMessages uint `json:"sent_messages"`
	ReceivedMessages uint `json:"received_messages"`
}

func (s WorkerStat) ToExported() WorkerStatExported {
	return WorkerStatExported{
		WorkerId: s.WorkerId,
		Date: s.Date.Format("2006-01-02"),
		SentMessages: s.SentMessages,
		ReceivedMessages: s.ReceivedMessages,
	}
}

type WorkerStatExported struct {
	WorkerId uint `json:"worker_id"`
	Date string `json:"date"`
	SentMessages uint `json:"sent_messages"`
	ReceivedMessages uint `json:"received_messages"`
}

func (r *ProfileRepository) GetWorkersStats() (workers []WorkerStat, err error) {
	err = r.db.Raw("SELECT worker_id, date, sum(sent_messages) sent_messages, sum(received_messages) received_messages " +
		"FROM worker_stats where date > ? " +
		"group by worker_id, date;", time.Now().UTC().Add(-7 * time.Hour * 24).Format("2006-01-02")).Scan(&workers).Error
	return
}

type WorkerActiveStat struct {
	WorkerId uint `json:"worker_id"`
	ActiveChats uint `json:"active_chats"`
}

func (r *ProfileRepository) GetActiveChats() (workers []WorkerActiveStat, err error) {
	err = r.db.Raw("SELECT worker_id, count(distinct(chat_id)) active_chats " +
		"FROM worker_stats where date = ? and sent_messages > 0 group by worker_id;", time.Now().UTC().Format("2006-01-02")).Scan(&workers).Error
	return
}
