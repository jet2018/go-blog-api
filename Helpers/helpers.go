package helpers

import (
	"blog/models"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"hash"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
)

// database connection, update from here to your new database
var dsn = "root:peacebewithyouall2020@tcp(127.0.0.1:3306)/blog?charset=utf8mb4&parseTime=True&loc=Local"
var Db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

// Connection provides the global connection to the database
func Connection() *gorm.DB {
	if err != nil {
		fmt.Printf("Failed to establish connection, %v\n", &err)
	}
	err = Db.AutoMigrate(&models.User{}, &models.Category{}, &models.Article{}, &models.Comment{})
	if err != nil {
		fmt.Println(err)
	}
	return Db
}

type Response struct {
	Status  int
	Message []string
	Error   []string
}

// SendResponse is the clean way of handling errors quickly
func SendResponse(c *gin.Context, response Response) {
	if len(response.Message) > 0 {
		c.JSON(response.Status, map[string]interface{}{"message": strings.Join(response.Message, "; ")})
	} else if len(response.Error) > 0 {
		c.JSON(response.Status, map[string]interface{}{"error": strings.Join(response.Error, "; ")})
	}
}

// Harsher Hashes the given raw text and returns its hashed string format
func Harsher(raw string) string {
	algorithm := sha256.New()
	return stringHarsher(algorithm, raw)
}

// Handles the hashing.
func stringHarsher(algorithm hash.Hash, text string) string {
	algorithm.Write([]byte(text))
	return hex.EncodeToString(algorithm.Sum(nil))
}

// JwtKey will be used to set up the jwt token
var JwtKey = []byte("mvpDQe5RgbnHs7+s8Bn7gOAqu1rt2cNbjkn2COFRc2NCIV1YYRG+AdyAjgoVd4oC0HUL5+TDSYxrn0RJ03i0e5bB0xk+m0AQQiyHLASmOAMdLUA9Vlr/2UnBkx8xc5VFCt7yIZDJCB2qWY3pePVYex6SWGgc/JczAl0be/Pg+vCQyH4ehFNiMXeTv4CrDlYdwqEfnxeEyHaXOd0tPdie1o2FHDsb61kUApVYFA")

var UploadTo = "static/"

// Upload uploads images to static folder
func Upload(c *gin.Context, image *multipart.FileHeader) (error, string, string) {

	// check if UploadTO folder exists, if does not, then create it.
	if _, err := os.Stat(UploadTo); os.IsNotExist(err) {
		if err := os.Mkdir(UploadTo, os.ModePerm); err != nil {
			log.Fatal(err)
		}
	}

	// the new path to the file saved, will be returned so that it can be further
	// acted on, like saving it to the database ot returning it directly to the
	// front-end
	path := UploadTo + image.Filename
	if err := c.SaveUploadedFile(image, path); err != nil {
		return err, "File couldn't be saved successfully", ""
	}
	return nil, "File saved successfully", path
}

func Session(c *gin.Context) (*models.User, bool, error) {
	var user models.User
	id, ok := c.Get("id")
	if !ok {
		return nil, false, errors.New("not authorised to perform this action")
	}
	Db.First(&user, id)
	// check if user's account is active
	if !user.IsActive {
		return nil, false, errors.New("not authorised to perform this action, account is still deactivated")
	}
	return &user, true, nil
}

func IsAuthenticated(c *gin.Context) *models.User {
	var user, ok, err = Session(c)
	if !ok || err != nil || user == nil {
		SendResponse(c, Response{
			Status: http.StatusUnauthorized,
			Error:  []string{err.Error()},
		})
		return nil
	}
	return user
}
