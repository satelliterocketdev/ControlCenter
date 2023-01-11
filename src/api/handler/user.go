package handler

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"firebase.google.com/go/auth"
	"fmt"
	"github.com/gin-gonic/gin"
	"i69social/src/api/firebase_admin"
	"i69social/src/api/jwt_auth"
	"i69social/src/api/models/user"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type UsersListResponse struct {
	Users           []user.User `json:"users"`
	CurrentPageSize int         `json:"current_page_size"`
	Offset          int         `json:"offset"`
	Size            int         `json:"size"`
}

type CreateUserRequest struct {
	Email string `json:"email"`
}

type UpdateUserRequest struct {
	Main    user.User       `json:"main"`
	Details json.RawMessage `json:"details"`
}

type ChangeUserRequest struct {
	Uid string `json:"uid"`
}

func (ctrl *Controller) exportUsers(c *gin.Context) {
	tokenParam, ok := c.GetQuery("token")
	if !ok || len(tokenParam) == 0 {
		c.Status(http.StatusUnauthorized)
		return
	}

	_, err := jwt_auth.ParseToken(tokenParam)
	if err != nil {
		log.Printf("token is not valid: `%s`", err)
		c.AbortWithStatus(401)
		return
	}

	b := &bytes.Buffer{}
	w := csv.NewWriter(b)

	if err = w.Write([]string{"email"}); err != nil {
		log.Fatalln("error writing record to csv:", err)
	}

	err = firebase_admin.ExportUsersEmails(w)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	w.Flush()

	if err := w.Error(); err != nil {
		log.Fatal(err)
	}

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", "attachment; filename=users.csv")
	c.Data(http.StatusOK, "text/csv", b.Bytes())
}

func (ctrl *Controller) getUser(c *gin.Context) {
	uid := c.Param("id")
	if len(uid) == 0 {
		c.Status(http.StatusBadRequest)
		return
	}

	u, err := firebase_admin.GetUser(uid)
	if auth.IsUserNotFound(err) {
		c.JSON(http.StatusNotFound, ErrorResponse{Message: fmt.Sprintf("User with uid `%s` not found", uid)})
		return
	}
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	reports, err := firebase_admin.GetUsersReports()
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	preparedUser := user.User{
		Email:       u.Email,
		UID:         u.UID,
		DisplayName: u.DisplayName,
		Disabled:    u.Disabled,
		PhoneNumber: u.PhoneNumber,
		PhotoURL:    u.PhotoURL,
		Reports:     reports[uid],
	}

	if u.UserMetadata != nil {
		preparedUser.SignUp = u.UserMetadata.CreationTimestamp
	}

	c.JSON(http.StatusOK, SuccessResponse{Data: preparedUser})
}

func (ctrl *Controller) getUserDetails(c *gin.Context) {
	uid := c.Param("id")
	if len(uid) == 0 {
		c.Status(http.StatusBadRequest)
		return
	}

	details, err := firebase_admin.GetUserDetails(uid)
	if auth.IsUserNotFound(err) {
		c.JSON(http.StatusNotFound, ErrorResponse{Message: fmt.Sprintf("User with uid `%s` not found", uid)})
		return
	}
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	if len(details) > 0 {
		c.JSON(http.StatusOK, SuccessResponse{Data: details})
	} else {
		c.JSON(http.StatusOK, SuccessResponse{Data: struct{}{}})
	}
}

func (ctrl *Controller) getUsers(c *gin.Context) {
	filter, ok := c.GetQuery("filter")
	if ok && len(filter) == 0 {
		c.Status(http.StatusBadRequest)
		return
	}

	if ok {
		searchByFilter(c, filter)
		return
	}

	orderedFlag, _ := c.GetQuery("ordered")
	var ordered bool
	var direction bool
	if orderedFlag == "1" {
		ordered = true
		direction = true
	}
	if orderedFlag == "-1" {
		ordered = true
		direction = false
	}
	var err error
	limit := 25
	offset := 0
	limitFlag, ok := c.GetQuery("limit")
	if ok {
		limit, err = strconv.Atoi(limitFlag)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		if limit <= 0 {
			c.Status(http.StatusBadRequest)
			return
		}
	}

	offsetFlag, ok := c.GetQuery("offset")
	if ok {
		offset, err = strconv.Atoi(offsetFlag)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		if offset < 0 {
			c.Status(http.StatusBadRequest)
			return
		}
	}

	var fieldFlag string
	if ordered {
		fieldFlag, ok = c.GetQuery("field")
		if !ok {
			c.Status(http.StatusBadRequest)
			return
		}
		fieldFlag = strings.TrimSpace(fieldFlag)
		if len(fieldFlag) == 0 {
			c.Status(http.StatusBadRequest)
			return
		}
	}

	users, size, err := firebase_admin.GetUsers(limit, offset, ordered, direction, fieldFlag)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: fmt.Sprintf("Error occurred when get users.")})
		log.Printf("Error occurred when get users. Orig err: `%s`", err)
		return
	}

	response := UsersListResponse{Users: users, CurrentPageSize: len(users), Size: size, Offset: offset}
	c.JSON(http.StatusOK, SuccessResponse{Data: response})
	return
}

var searchByFilter = func(c *gin.Context, filter string) {
	users, err := firebase_admin.GetUsersByFilter(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: fmt.Sprintf("Error occurred when get users by filter `%s`. Orig err: `%s`", filter, err.Error())})
		return
	}

	response := UsersListResponse{Users: users}
	c.JSON(http.StatusOK, SuccessResponse{Data: response})
}


func (ctrl *Controller) createUser(c *gin.Context) {
	var request CreateUserRequest
	err := c.Bind(&request)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	u, err := firebase_admin.CreateUser(request.Email)
	if auth.IsEmailAlreadyExists(err) {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: fmt.Sprintf("User with email `%s` is already exists", request.Email)})
		return
	}
	if err != nil {
		log.Printf("error when create user: %s", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	newUser := user.User{
		Email:       u.Email,
		UID:         u.UID,
		DisplayName: u.DisplayName,
		Disabled:    u.Disabled,
		PhoneNumber: u.PhoneNumber,
		PhotoURL:    u.PhotoURL,
	}

	if u.UserMetadata != nil {
		newUser.SignUp = u.UserMetadata.CreationTimestamp
	}

	c.JSON(http.StatusOK, SuccessResponse{Data: newUser})
}

func (ctrl *Controller) updateUser(c *gin.Context) {
	uid := c.Param("id")
	if len(uid) == 0 {
		c.Status(http.StatusBadRequest)
		return
	}

	var request UpdateUserRequest
	err := c.Bind(&request)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	err = firebase_admin.UpdateUser(uid, request.Main, request.Details)
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

func (ctrl *Controller) blockUser(c *gin.Context) {
	uid := c.Param("id")
	if uid == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "User uid cannot be empty"})
		return
	}

	err := firebase_admin.BlockUser(uid)
	if auth.IsUserNotFound(err) {
		c.JSON(http.StatusNotFound, ErrorResponse{Message: fmt.Sprintf("User with uid `%s` not found", uid)})
		return
	}
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}

func (ctrl *Controller) unblockUser(c *gin.Context) {
	uid := c.Param("id")
	if uid == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "User uid cannot be empty"})
		return
	}

	err := firebase_admin.UnBlockUser(uid)
	if auth.IsUserNotFound(err) {
		c.JSON(http.StatusNotFound, ErrorResponse{Message: fmt.Sprintf("User with uid `%s` not found", uid)})
		return
	}
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}


func (ctrl *Controller) deleteUser(c *gin.Context) {
	uid := c.Param("id")
	if uid == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "User uid cannot be empty"})
		return
	}

	err := firebase_admin.RemoveUser(uid)
	if auth.IsUserNotFound(err) {
		c.JSON(http.StatusNoContent, ErrorResponse{Message: fmt.Sprintf("User with uid `%s` not found", uid)})
		return
	}
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	firebase_admin.RemoveUserDetails(uid)

	c.Status(http.StatusNoContent)
}

func (ctrl *Controller) deleteUserReports(c *gin.Context) {
	uid := c.Param("id")
	if uid == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "User uid cannot be empty"})
		return
	}

	err := firebase_admin.RemoveUserReports(uid)
	if auth.IsUserNotFound(err) {
		c.JSON(http.StatusNoContent, ErrorResponse{Message: fmt.Sprintf("User with uid `%s` not found", uid)})
		return
	}
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusNoContent)
}

