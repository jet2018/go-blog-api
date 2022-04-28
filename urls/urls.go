package urls

import (
	"blog/views"
	"fmt"
	"github.com/gin-gonic/gin"
)

func RequestHandler() {
	// define routes alongside their handlers
	router := gin.Default()
	router.Static("/static", "./static")
	// version one routes
	v1 := router.Group("/v1")
	{
		// endpoints that can be hit without authentication
		v1.GET("/", views.Index)
		v1.POST("/users", views.AddUser)
		v1.GET("/users/:username", views.UserByUsernameOrId)
		v1.POST("/login", views.Login)
		v1.GET("/articles", views.ListArticles)
		v1.GET("/articles/:id", views.SingleArticle)

		// the endpoints below will require you to be logged in
		v1.Use(views.VerifyToken)
		v1.GET("/users", views.GetUsers)
		v1.GET("/users/me", views.Me)
		v1.GET("/users/me/deactivate", views.DeactivateMyAccount)
		v1.POST("/articles", views.CreateArticle)
		v1.DELETE("/articles/:id", views.DeleteAnArticle)

	}

	err := router.Run(":8080")
	if err != nil {
		fmt.Println("Server down: ", err)
	}
}
