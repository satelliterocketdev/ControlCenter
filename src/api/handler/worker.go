package handler

import (
	"github.com/gin-gonic/gin"
	"i69social/src/api/database"
	"i69social/src/api/jwt_auth"
	"i69social/src/api/models/fake_user"
	"i69social/src/api/models/profile"
	"i69social/src/api/services"
	"log"
	"net/http"
	"strconv"
)

type WorkersResponse struct {
	Workers     []profile.Profile `json:"workers"`
	LastContact map[string]string `json:"last_contact"`
}

func (*Controller) getWorkers(c *gin.Context) {

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

	userService := services.GetProfileService(database.GetMySQLDB(), database.GetRedisDB())
	workers, err := userService.GetWorkers()
	if err != nil {
		log.Printf("get workers error: `%s`", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "Error occurred when get workers list"})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{Data: workers})
}

func (*Controller) deleteWorker(c *gin.Context) {

	workerIdValue := c.Param("id")
	if len(workerIdValue) == 0 {
		c.Status(http.StatusBadRequest)
		return
	}

	workerId, err := strconv.Atoi(workerIdValue)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	userService := services.GetProfileService(database.GetMySQLDB(), database.GetRedisDB())
	err = userService.DeleteWorker(uint(workerId))
	if err != nil {
		log.Printf("get workers error: `%s`", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "Error occurred when get workers list"})
		return
	}

	c.Status(http.StatusNoContent)
}

type FakeUsersResponse struct {
	FakeUsers []fake_user.User `json:"fake_users"`
}

func (*Controller) getWorkerFakeUsers(c *gin.Context) {

	workerIdValue := c.Param("id")
	if len(workerIdValue) == 0 {
		c.Status(http.StatusBadRequest)
		return
	}

	workerId, err := strconv.Atoi(workerIdValue)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	userService := services.GetProfileService(database.GetMySQLDB(), database.GetRedisDB())
	users, err := userService.GetWorkerFakeUsers(uint(workerId))
	if err != nil {
		log.Printf("get fake users error: `%s`", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "Error occurred when get fake users list"})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{Data: FakeUsersResponse{FakeUsers: users}})
}
