package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"i69social/src/api/database"
	"i69social/src/api/firebase_admin"
	"i69social/src/api/services"
	"log"
	"net/http"
)

func (ctrl *Controller) getPresets(c *gin.Context) {
	presets, err := firebase_admin.GetPresets()
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{Data: presets})
}

type UpdateFakeUserRequest struct {
	OwnerId uint `json:"owner_id"`
	Data json.RawMessage `json:"data"`
}

func (ctrl *Controller) updateFakeUser(c *gin.Context) {
	uid := c.Param("id")
	if len(uid) == 0 {
		c.Status(http.StatusBadRequest)
		return
	}

	var request UpdateFakeUserRequest
	err := c.Bind(&request)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	svc := services.GetFakeUserService(database.GetMySQLDB())
	err = svc.UpdateFakeUser(uid, request.Data)
	if err == firebase_admin.UpdateRequestInvalidParamsErr {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: fmt.Sprintf("User with uid `%s` is not found", uid)})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

type CreateFakeUserRequest struct {
	OwnerId uint `json:"owner_id"`
	Data json.RawMessage `json:"data"`
}

func (ctrl *Controller) createFakeUser(c *gin.Context) {
	var request CreateFakeUserRequest
	err := c.Bind(&request)
	if err != nil {
		log.Printf("given create fake user requst is invalid: %s", err)
		c.Status(http.StatusBadRequest)
		return
	}

	svc := services.GetFakeUserService(database.GetMySQLDB())
	err = svc.CreateFakeUser(request.OwnerId, request.Data)
	if err != nil {
		log.Printf("error on save fake user: %s", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, SuccessResponse{Data: request})
}

func (ctrl *Controller) getFakeUserByExternalId(c *gin.Context) {
	uid := c.Param("id")
	if len(uid) == 0 {
		c.Status(http.StatusBadRequest)
		return
	}

	details, err := firebase_admin.GetFakeUserDetails(uid)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	if details == nil {
		c.Status(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{Data: details})
}

type ChangeOnlineStatusRequest struct {
	Online bool `json:"online"`
}

func (ctrl *Controller) changeFakeUserOnlineStatus(c *gin.Context) {
	uid := c.Param("id")
	if len(uid) == 0 {
		c.Status(http.StatusBadRequest)
		return
	}

	var request ChangeOnlineStatusRequest
	err := c.Bind(&request)
	if err != nil {
		log.Printf("given change fake user online status request is invalid: %s", err)
		c.Status(http.StatusBadRequest)
		return
	}

	err = firebase_admin.ChangeFakeUserOnlineStatus(uid, request.Online)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}
