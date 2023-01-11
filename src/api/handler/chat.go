package handler

import (
	"github.com/gin-gonic/gin"
	"i69social/src/api/database"
	"i69social/src/api/jwt_auth"
	"i69social/src/api/services"
	"log"
	"net/http"
)

type AddNewChatMessageRequest struct {
	Message        string `json:"message"`
	FromName       string `json:"from_name"`
	FromUserId     string `json:"from_user_id"`
	MessageType string `json:"type"`
	RecipientToken string `json:"recipient_token"`
}

func (*Controller) addChatMessage(c *gin.Context) {

	chatId := c.Param("chatId")
	if len(chatId) == 0 {
		c.Status(http.StatusBadRequest)
		return
	}

	var request AddNewChatMessageRequest
	err := c.Bind(&request)
	if err != nil {
		log.Printf("given add new chat message request is invalid: %s", err)
		c.Status(http.StatusBadRequest)
		return
	}

	cl, ok := c.Get("claims")
	if !ok {
		c.Status(401)
		return
	}

	userInfo, ok := cl.(*jwt_auth.UserClaims)
	if !ok {
		c.Status(401)
		return
	}
	if userInfo.Id == 0 {
		c.Status(401)
		return
	}

	chatService := services.GetChatService(database.GetMySQLDB(), database.GetRedisDB())
	err = chatService.AddNewMessage(
		userInfo.Id,
		chatId,
		request.FromName,
		request.FromUserId,
		request.RecipientToken,
		request.MessageType,
		request.Message,
	)
	if err != nil {
		log.Printf("add new chat message error: `%s`", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "Error occurred when send new message"})
		return
	}

	c.Status(http.StatusOK)
}

type ReadChatMessagesRequest struct {
	Messages []string `json:"messages"`
	UserId   string   `json:"user_id"`
}

func (*Controller) readChatMessage(c *gin.Context) {

	chatId := c.Param("chatId")
	if len(chatId) == 0 {
		c.Status(http.StatusBadRequest)
		return
	}

	var request ReadChatMessagesRequest
	err := c.Bind(&request)
	if err != nil {
		log.Printf("given add new chat message request is invalid: %s", err)
		c.Status(http.StatusBadRequest)
		return
	}

	cl, ok := c.Get("claims")
	if !ok {
		c.Status(401)
		return
	}

	userInfo, ok := cl.(*jwt_auth.UserClaims)
	if !ok {
		c.Status(401)
		return
	}
	if userInfo.Id == 0 {
		c.Status(401)
		return
	}

	chatService := services.GetChatService(database.GetMySQLDB(), database.GetRedisDB())
	err = chatService.ReadChatMessages(userInfo.Id, chatId, request.UserId, request.Messages)
	if err != nil {
		log.Printf("read chat messages error: `%s`", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "Error occurred when read chat messages"})
		return
	}

	c.Status(http.StatusOK)
}

type SaveChatNotesRequest struct {
	Notes string `json:"notes"`
}

func (*Controller) saveNotes(c *gin.Context) {

	chatId := c.Param("chatId")
	if len(chatId) == 0 {
		c.Status(http.StatusBadRequest)
		return
	}

	var request SaveChatNotesRequest
	err := c.Bind(&request)
	if err != nil {
		log.Printf("given save chat notes request is invalid: %s", err)
		c.Status(http.StatusBadRequest)
		return
	}

	chatService := services.GetChatService(database.GetMySQLDB(), database.GetRedisDB())
	err = chatService.SaveChatNotes(chatId, request.Notes)
	if err != nil {
		log.Printf("read chat messages error: `%s`", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "Error occurred when save chat notes"})
		return
	}

	c.Status(http.StatusOK)
}

type StartNewChatRequest struct {
	FakeUserId string `json:"fake_user_id"`
	RealUserId string `json:"real_user_id"`
}
func (*Controller) startNewChat(c *gin.Context) {
	var request StartNewChatRequest
	err := c.Bind(&request)
	if err != nil {
		log.Printf("given start new chat request is invalid: %s", err)
		c.Status(http.StatusBadRequest)
		return
	}

	chatService := services.GetChatService(database.GetMySQLDB(), database.GetRedisDB())
	err = chatService.StartNewChat(request.FakeUserId, request.RealUserId)
	if err != nil {
		log.Printf("start new new shat error: `%s`", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "Error occurred when start new chat"})
		return
	}

	c.Status(http.StatusOK)
}
