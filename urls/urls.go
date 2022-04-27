package urls

import (
	"blog/views"
	"fmt"
	"github.com/gin-gonic/gin"
)

func RequestHandler() {
	// define routes alongside their handlers
	router := gin.Default()

	// version one routes
	v1 := router.Group("/v1")
	{
		// endpoints that can be hit without authentication
		v1.GET("/", views.Index)
		v1.POST("/users", views.AddUser)
		v1.GET("/users/:username", views.UserByUsername)
		v1.POST("/login", views.Login)

		// the endpoints below will require you to be logged in
		v1.Use(views.VerifyToken)
		v1.GET("/users", views.GetUsers)
		v1.GET("/users/me", views.Me)
		v1.GET("/users/me/deactivate", views.DeactivateMyAccount)

	}

	err := router.Run(":8080")
	if err != nil {
		fmt.Println("Server down: ", err)
	}
}
