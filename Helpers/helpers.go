package helpers

import (
	"blog/models"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"hash"
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
