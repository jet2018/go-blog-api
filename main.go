package main

import (
	helpers "blog/Helpers"
	"blog/urls"
	"fmt"
)

func main() {
	// handles connection to the database
	helpers.Connection()
	// starts and keeps the server running
	fmt.Println("Starting server")
	urls.RequestHandler()
}
