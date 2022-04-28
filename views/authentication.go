package views

import (
	helpers "blog/Helpers"
	"blog/models"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"net/http"
	"strings"
	"time"
)

func FindUserByUsernameOrEmail(username string) (*models.User, error) {
	var user models.User

	if res := helpers.Db.Where("username = ?", username).Or("email = ?", username).Find(&user); res.Error != nil {
		return nil, res.Error
	}
	if !user.IsActive {
		return nil, errors.New("user account is inactive")
	}
	return &user, nil
}

func FindUserByID(id uint) (*models.User, error) {
	var user models.User
	if res := helpers.Db.Find(&user, id); res.Error != nil {
		return nil, res.Error
	}
	return &user, nil
}

type authClaims struct {
	jwt.StandardClaims
	User   models.User `json:"user"`
	UserId uint        `json:"user_id"`
}

func generateToken(user models.User) (string, error) {
	expiresAt := time.Now().Add(24 * time.Hour).Unix()
	fmt.Println(expiresAt)
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, authClaims{
		StandardClaims: jwt.StandardClaims{
			Subject:   user.Username,
			ExpiresAt: expiresAt,
		},
		UserId: user.ID,
		User:   user,
	})
	tokenString, err := token.SignedString(helpers.JwtKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ValidateToken(tokenString string) (uint, string, error) {
	var claims authClaims
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return helpers.JwtKey, nil
	})
	if err != nil {
		return 0, "", err
	}
	if !token.Valid {
		return 0, "", errors.New("invalid token")
	}
	id := claims.UserId
	username := claims.Subject
	var user models.User
	helpers.Db.First(&user, "is_active = ?", true)
	if !user.IsActive {
		return 0, "", errors.New("user account in inactive")
	}
	return id, username, nil
}

// Login view
// @returns JWT token
func Login(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "incorrect parameters",
		})
		return
	}

	user, err := FindUserByUsernameOrEmail(req.Username)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": fmt.Sprintf("user with %s not found", req.Username),
		})
		return
	}
	if user.Password != helpers.Harsher(req.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "incorrect password",
		})
		return
	}
	token, err := generateToken(*user)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Key generation error: " + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

// VerifyToken middleware for verifying tokens
func VerifyToken(c *gin.Context) {
	token, ok := getToken(c)
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Error": "You have no access to perform this action"})
		return
	}
	id, username, err := ValidateToken(token)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Error": "Invalid token provided"})
		return
	}
	c.Set("id", id)
	c.Set("username", username)
	c.Writer.Header().Set("Authorization", "Bearer "+token)
	c.Next()
}

func getToken(c *gin.Context) (string, bool) {
	authValue := c.GetHeader("Authorization")
	arr := strings.Split(authValue, " ")
	if len(arr) != 2 {
		return "", false
	}
	authType := strings.Trim(arr[0], "\n\r\t")
	if strings.ToLower(authType) != strings.ToLower("Bearer") {
		return "", false
	}
	return strings.Trim(arr[1], "\n\t\r"), true
}
