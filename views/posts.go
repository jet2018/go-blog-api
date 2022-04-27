package views

import (
	helpers "blog/Helpers"
	"blog/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Index(c *gin.Context) {
	var articles []models.Article
	helpers.Db.Find(&articles)
	c.JSON(http.StatusOK, articles)
}

func GetPosts(w http.ResponseWriter, r *http.Request) {
}
