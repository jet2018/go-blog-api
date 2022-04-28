package views

import (
	helpers "blog/Helpers"
	"blog/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
)

func ListArticles(c *gin.Context) {
	var articles []models.Article
	helpers.Db.Find(&articles)
	c.JSON(http.StatusOK, articles)
}

func CreateArticle(c *gin.Context) {
	// check if user is logged in
	checkAuth := helpers.IsAuthenticated(c)
	if checkAuth == nil {
		return
	}
	Id := checkAuth.ID
	coverImage, _ := c.FormFile("cover_photo")
	body := c.PostForm("body")
	title := c.PostForm("title")
	filePath := ""

	// check if someone has an image and save it, otherwise, bypass this.
	if coverImage != nil {
		err, msg, path := helpers.Upload(c, coverImage)
		if err != nil {
			helpers.SendResponse(c, helpers.Response{
				Status: http.StatusInternalServerError,
				Error:  []string{msg + "::" + err.Error()},
			})
		}
		filePath = path
	}

	validate := validator.New()
	article := models.Article{
		UserId: int(Id),
		Cover:  filePath,
		Body:   body,
		Title:  title,
	}
	err := validate.Struct(&article)
	if err != nil {
		helpers.SendResponse(c, helpers.Response{Status: http.StatusBadRequest, Error: []string{string(err.Error())}})
		return
	}
	helpers.Db.Create(&article)
	c.JSON(http.StatusOK, gin.H{
		"message": "Your article has been successfully created.",
		"article": &article,
	})
}

func DeleteAnArticle(c *gin.Context) {
	Id := c.Param("id")
	user := helpers.IsAuthenticated(c)

	if user == nil {
		return
	}

	var article models.Article

	res := helpers.Db.Find(&article, Id)
	if res != nil {
		helpers.Db.Delete(&article)
		helpers.SendResponse(c, helpers.Response{
			Status: http.StatusNoContent,
			Error:  []string{"article deleted"},
		})
		return
	} else {
		helpers.SendResponse(c, helpers.Response{
			Status: http.StatusNotFound,
			Error:  []string{"article not found"},
		})
		return
	}

}

func SingleArticle(c *gin.Context) {
	id := c.Param("id")
	var article models.Article
	var count int64
	res := helpers.Db.First(&article, "id = ?", id)
	res.Count(&count)
	if count > 0 {
		c.JSON(http.StatusOK, &article)
	} else {
		helpers.SendResponse(c, helpers.Response{
			Status: http.StatusNotFound,
			Error:  []string{"No associated article"},
		})
	}
}
