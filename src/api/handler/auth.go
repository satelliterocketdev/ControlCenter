package handler

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"i69social/src/api/database"
	"i69social/src/api/jwt_auth"
	"i69social/src/api/models/profile"
	"i69social/src/api/models/tools"
	"i69social/src/api/services"
	"log"
	"net/http"
)

type AuthResponse struct {
	Token string `json:"token"`
}

func (*Controller) SignUp(c *gin.Context) {

	var params profile.SignUpRequest
	err := c.ShouldBindJSON(&params)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: fmt.Sprintf("Given request to create user is invalid. Orig err: `%s`", err)})
		return
	}
	params.TrimSpaces()

	err = tools.Validator.Struct(params)
	if err != nil {
		c.Writer.WriteHeader(http.StatusUnauthorized)
		log.Printf("SignIn Credentials are invalid. Orig err: `%s`", err)
		return
	}

	userService := services.GetProfileService(database.GetMySQLDB(), database.GetRedisDB())
	_, err = userService.Create(params)
	if err != nil {
		switch err.(type) {
		case *tools.ValidationErr:
			log.Printf("create new user validation error: `%s`", err)
			c.JSON(http.StatusBadRequest, ErrorResponse{Message: fmt.Sprintf("Given request is invalid. Orig err: `%s`", err)})
			return
		}

		log.Printf("create new user error: `%s`", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "Error occurred when create new user"})
		return
	}



	c.Status(http.StatusOK)
}

func (ctrl *Controller) signIn(c *gin.Context) {

	var credentials profile.SignInRequest
	err := c.Bind(&credentials)
	if err != nil {
		log.Println(err)
		//c.Writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	err = tools.Validator.Struct(credentials)
	if err != nil {
		c.Writer.WriteHeader(http.StatusUnauthorized)
		log.Printf("SignIn Credentials are invalid. Orig err: `%s`", err)
		return
	}

	profileSrc := services.GetProfileService(database.GetMySQLDB(),  database.GetRedisDB())
	p, err := profileSrc.FindUser(credentials)
	if err != nil {
		var validationErr *tools.ValidationErr
		if errors.As(err, &validationErr) {
			log.Printf("validate error %s", err)
			c.JSON(http.StatusUnauthorized, ErrorResponse{Message: fmt.Sprintf("Given credentials is invalid.")})
			return
		}


		log.Printf("internal error %s", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "Error occurred when get user"})
		return
	}

	token, err := jwt_auth.GetToken(jwt_auth.GetUserClaims(p))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "Error occurred when generate access token"})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{Data: AuthResponse{Token: token}})
}


func (ctrl *Controller) renew(c *gin.Context) {

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

	token, err := jwt_auth.GetToken(userInfo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "Error occurred when generate access token"})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{Data: AuthResponse{Token: token}})
}
