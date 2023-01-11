package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"i69social/src/api/database"
	"i69social/src/api/models/profile"
	"i69social/src/api/models/tools"
	"i69social/src/api/services"
	"log"
	"net/http"
)


type GenerateInvitationResponse struct {
	Link string `json:"link"`
}

func (*Controller) createInvitation(c *gin.Context) {

	var params profile.GenerateInvitationRequest
	err := c.ShouldBindJSON(&params)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: fmt.Sprintf("Given request to generate invitation is invalid. Orig err: `%s`", err)})
		return
	}


	params.TrimSpaces()

	userService := services.GetProfileService(database.GetMySQLDB(), database.GetRedisDB())
	link, err := userService.GenerateInvitation(params)
	if err != nil {
		switch err.(type) {
		case *tools.ValidationErr:
			log.Printf("generate ivitation validation error: `%s`", err)
			c.JSON(http.StatusBadRequest, ErrorResponse{Message: fmt.Sprintf("Given request is invalid. Orig err: `%s`", err)})
			return
		}

		log.Printf("generate ivitation error: `%s`", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "Error occurred when generate invitation link"})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{Data: GenerateInvitationResponse{Link: link}})
}

type GetInvitationResponse struct {
	Key string `json:"key"`
	Email string `json:"email"`
}

func (*Controller) getInvitation(c *gin.Context) {

	key := c.Param("key")
	if len(key) != 10 {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: fmt.Sprintf("Given invitation is invalid.")})
		return
	}



	userService := services.GetProfileService(database.GetMySQLDB(), database.GetRedisDB())
	params, err := userService.GetInvitationParams(key)
	if err != nil {
		switch err.(type) {
		case *tools.ValidationErr:
			log.Printf("generate invitation validation error: `%s`", err)
			c.JSON(http.StatusBadRequest, ErrorResponse{Message: fmt.Sprintf("Given request is invalid. Orig err: `%s`", err)})
			return
		}

		log.Printf("generate invitation error: `%s`", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "Error occurred when generate invitation link"})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{Data: GetInvitationResponse{Key: key, Email: params.Email}})
}
