package views

import (
	helpers "blog/Helpers"
	"blog/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
)

// GetUsers returns all users
func GetUsers(c *gin.Context) {
	// gives you whether the user is authenticated and the actual user who is logged in
	//user, ok := Session(c)
	//if ok {
	//	fmt.Println(user)
	//}
	var users []models.User
	helpers.Db.Find(&users)
	c.JSON(http.StatusOK, users)
}

// AddUser is for creating an account
func AddUser(c *gin.Context) {
	//var user models.User
	//if err := c.ShouldBindQuery(&user); err == nil {
	//	fmt.Println("Submitted user object is %v", &user)
	//} else {
	//	fmt.Println("Error -%+v", err)
	//}
	validate := validator.New()

	firstName := c.PostForm("first_name")
	email := c.PostForm("email")
	lastName := c.PostForm("last_name")
	password := c.PostForm("password")
	newPassword := helpers.Harsher(password)
	phone := c.PostForm("phone")
	username := c.PostForm("username")
	user := models.User{
		Email:     &email,
		Username:  username,
		Password:  newPassword,
		LastName:  lastName,
		FirstName: firstName,
		Phone:     phone,
	}

	err := validate.Struct(&user)
	if err != nil {
		helpers.SendResponse(c, helpers.Response{Status: http.StatusBadRequest, Error: []string{string(err.Error())}})
	}

	// check if anyone has the username, email, or phone given
	if helpers.Db.Where("username = ?", &user.Username).Find(&user).RowsAffected > 0 {
		helpers.SendResponse(c, helpers.Response{
			Status: http.StatusBadRequest, Error: []string{"Username provided is already taken"},
		})
	} else if helpers.Db.Where("email = ?", &user.Email).Find(&user).RowsAffected > 0 {
		helpers.SendResponse(c, helpers.Response{
			Status: http.StatusBadRequest, Error: []string{"Email provided is already taken"},
		})
	} else if helpers.Db.Where("phone = ?", &user.Phone).Find(&user).RowsAffected > 0 {
		helpers.SendResponse(c, helpers.Response{
			Status: http.StatusBadRequest, Error: []string{"Phone provided is already taken"},
		})
	} else {
		helpers.Db.Create(&user)
		helpers.SendResponse(c, helpers.Response{Status: http.StatusOK, Error: []string{"User account created successfully"}})
	}

}

// Me si the logged-in user profile
func Me(c *gin.Context) {
	user, _, err := helpers.Session(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, &user)
}

// UserByUsernameOrId returns the profile of the user by their username or id,
// both must be passed as username, the differences will be internal
func UserByUsernameOrId(c *gin.Context) {
	var user models.User
	username := c.Param("username")
	res := helpers.Db.First(&user, "username = ? || id = ?", username, username)
	if res.RowsAffected > 0 {
		c.JSON(http.StatusOK, &user)
		return
	} else {
		c.JSON(http.StatusNotFound, gin.H{"error": "No associated account"})
		return
	}

}

func DeactivateMyAccount(c *gin.Context) {
	user, ok, err := helpers.Session(c)
	if ok {
		var newUser *models.User
		Id := user.ID
		helpers.Db.First(&newUser, Id)
		newUser.IsActive = false
		helpers.Db.Save(&newUser)
		c.JSON(http.StatusOK, gin.H{"success": "Account deactivated successfully"})
	}
	c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
}
