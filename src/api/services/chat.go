package services

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	"i69social/src/api/firebase_admin"
	"i69social/src/api/models/profile"
	"log"
)

type Chat struct {
	repo *profile.ProfileRepository
}

func GetChatService(db *gorm.DB, redisDb *redis.Client) *Chat {
	return &Chat{repo: profile.GetProfileRepository(db, redisDb)}
}

func (s *Chat) AddNewMessage(workerId uint, chatId, fromName, fromUserId, recipientToken, messageType, message string) (err error) {
	err = firebase_admin.AddChatMessage(fromUserId, chatId, messageType, message)
	if err != nil {
		return
	}
	s.repo.StoreWorkerLastMessage(profile.UsersLastMessage, fmt.Sprintf("%d", workerId))

	s.repo.SaveSentMessages(workerId, fromUserId, chatId, 1)

	if len(recipientToken) == 0 {
		return
	}

	err = firebase_admin.SendChatMessageToFcm(recipientToken, fromName, messageType, message)
	if err != nil {
		log.Printf("Error occurred while send notification to %s: err %s", recipientToken, err)
		err = nil
		return
	}

	return
}


func (s *Chat) ReadChatMessages(workerId uint, chatId, userId string, messages []string) (err error) {
	err = firebase_admin.ReadChatMessages(chatId, messages)
	if err != nil {
		return
	}

	s.repo.SaveReceivedMessages(workerId, userId, chatId, uint(len(messages)))
	return
}

func (s *Chat) SaveChatNotes(chatId, notes string) (err error) {
	err = firebase_admin.SaveChatNotes(chatId, notes)
	if err != nil {
		return
	}

	return
}

func (s *Chat) StartNewChat(fakeUserId, realUserId string) (err error) {
	err = firebase_admin.StartNewChat(fakeUserId, realUserId)
	if err != nil {
		return
	}

	return
}
