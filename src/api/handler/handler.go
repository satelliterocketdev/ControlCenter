package handler

import (
	"github.com/gin-gonic/gin"
	"i69social/src/api/database"
	"i69social/src/api/jwt_auth"
	"i69social/src/api/models/profile"
	"i69social/src/api/services"
	"log"
	"net/http"
	"strings"
)

type SuccessResponse struct {
	Data interface{} `json:"data"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

func SetupHandler() http.Handler {
	r := gin.Default()
	r.RedirectFixedPath = true
	r.Static("/static", "src/api/static")
	c := NewController()

	r.GET("/", c.static)

	r.POST("/auth/signIn", c.signIn)
	r.POST("/auth/signUp", c.SignUp)
	/*authorizedPublic := r.Use(authorized())
	authorizedPublic.GET("/auth/renew", c.renew)*/

	r.GET("/internal/users.csv", c.exportUsers)
	// Secure API
	internalApi := r.Group("/internal")
	{
		internalApi.Use(authorized())
		ctrlSecureRegular := internalApi.Use(regularAccess())
		ctrlSecureAdmin := internalApi.Group("/").Use(adminAccess())

		ctrlSecureRegular.GET("/users", c.getUsers)
		ctrlSecureRegular.GET("/users/:id", c.getUser)
		ctrlSecureRegular.GET("/users/:id/details", c.getUserDetails)

		ctrlSecureRegular.POST("/users/", c.createUser)

		ctrlSecureRegular.PUT("/users/:id", c.updateUser)
		ctrlSecureRegular.PUT("/users/:id/block", c.blockUser)
		ctrlSecureRegular.PUT("/users/:id/unblock", c.unblockUser)

		ctrlSecureRegular.DELETE("/users/:id", c.deleteUser)
		ctrlSecureRegular.DELETE("/users/:id/reports", c.deleteUserReports)

		ctrlSecureRegular.GET("/presets/fakeUsers", c.getPresets)
		ctrlSecureRegular.POST("/fakeUsers", c.createFakeUser)
		ctrlSecureRegular.PUT("/fakeUsers/:id/details", c.updateFakeUser)
		ctrlSecureRegular.PUT("/fakeUsers/:id/online", c.changeFakeUserOnlineStatus)
		ctrlSecureRegular.GET("/fakeUsers/:id/details", c.getFakeUserByExternalId)

		ctrlSecureRegular.PUT("/chats/:chatId/send", c.addChatMessage)
		ctrlSecureRegular.PUT("/chats/:chatId/read", c.readChatMessage)
		ctrlSecureRegular.PUT("/chats/:chatId/notes", c.saveNotes)
		ctrlSecureRegular.POST("/chats/new", c.startNewChat)

		ctrlSecureRegular.GET("workers/:id/fakeUsers", c.getWorkerFakeUsers)

		ctrlSecureAdmin.POST("/invitation", c.createInvitation)
		r.GET("/internal/invitation/:key", c.getInvitation)

		ctrlSecureAdmin.GET("/workers", c.getWorkers)

		ctrlSecureSuperAdmin := internalApi.Group("/").Use(superAdminAccess())
		ctrlSecureSuperAdmin.DELETE("/workers/:id", c.deleteWorker)
	}


	return r
}

func authorized() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		authHeaderParts := strings.Split(authHeader, " ")
		if len(authHeaderParts) != 2 {
			log.Printf("auth header has wrong format. auth header: `%s`", authHeader)
			c.AbortWithStatus(401)
			return
		}

		if authHeaderParts[0] != "Bearer" {
			log.Printf("auth header has wrong format. auth marker: `%s`", authHeaderParts[0])
			c.AbortWithStatus(401)
			return
		}

		accessToken := authHeaderParts[1]
		userClaims, err := jwt_auth.ParseToken(accessToken)
		if err != nil {
			log.Printf("token is not valid: `%s`", err)
			c.AbortWithStatus(401)
			return
		}
		go func() {
			p := services.GetProfileService(database.GetMySQLDB(), database.GetRedisDB())
			err = p.StoreWorkerLastContact(userClaims.Id)
			if err != nil {
				log.Printf("Error when store worker last contact %s", err)
			}
		}()

		err = services.GetProfileService(database.GetMySQLDB(), database.GetRedisDB()).GetProfileById(userClaims.Id)
		if err != nil {
			log.Printf("error after check in DB: `%s`", err)
			c.AbortWithStatus(401)
			return
		}

		c.Set("claims", userClaims)
		c.Next()

	}
}

func regularAccess() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, ok := c.Get("claims")
		if !ok {
			c.AbortWithStatus(401)
			return
		}

		userClaims, ok := claims.(*jwt_auth.UserClaims)
		if !ok {
			c.AbortWithStatus(401)
			return
		}

		if !userClaims.Privileges.Has(profile.RegularUserPrivilege) {
			c.AbortWithStatus(403)
			return
		}
		log.Println("Regular")
		c.Next()
	}
}

func adminAccess() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, ok := c.Get("claims")
		if !ok {
			c.AbortWithStatus(401)
			return
		}

		userClaims, ok := claims.(*jwt_auth.UserClaims)
		if !ok {
			c.AbortWithStatus(401)
			return
		}

		if !userClaims.Privileges.Has(profile.AdminUserPrivilege) {
			c.AbortWithStatus(403)
			return
		}
		log.Println("Admin")
		c.Next()
	}
}


func superAdminAccess() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, ok := c.Get("claims")
		if !ok {
			c.AbortWithStatus(401)
			return
		}

		userClaims, ok := claims.(*jwt_auth.UserClaims)
		if !ok {
			c.AbortWithStatus(401)
			return
		}

		if !userClaims.Privileges.Has(profile.SuperAdminUserPrivilege) {
			c.AbortWithStatus(403)
			return
		}
		log.Println("Admin")
		c.Next()
	}
}
