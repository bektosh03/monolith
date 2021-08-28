package api

import (
	"net/http"

	"github.com/bektosh03/monolith/api/handlers"
	"github.com/gin-gonic/gin"
)

const token = "abc123"

func New(h *handlers.Handler) *gin.Engine {
	r := gin.Default()
	gin.SetMode(gin.DebugMode)

	r.Use(func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth != token {
			c.AbortWithStatusJSON(http.StatusForbidden, struct{
				Error string `json:"error"`
			}{Error: "NOT ENOUGH RIGHTS"})
		}
	})

	// Root route
	r.GET("/", h.Hello)
	
	// User create
	r.POST("/user/", h.CreateUser)

	// List Users
	r.GET("/users/", h.GetUsers)
	r.GET("/user/:email/", h.GetUserByEmail)
	r.DELETE("/user/delete/:email/", h.Delete)

	return r
}
