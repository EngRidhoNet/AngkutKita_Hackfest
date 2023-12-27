// main.go
package main

import (
	"AngkutKita/RegisterLogin/handlers"
	"AngkutKita/RegisterLogin/models"
	"fmt"
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB
var err error

func initialMigration() {
	db.AutoMigrate(&models.User{})
}

func main() {
	dsn := "root:@tcp(localhost:3306)/angkutkita?parseTime=true"
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		return
	}

	initialMigration()

	router := mux.NewRouter()

	// Register handlers with the DB instance
	router.HandleFunc("/register", handlers.RegisterHandler(db)).Methods("POST")
	router.HandleFunc("/login/{username}/{password}", handlers.LoginHandler(db)).Methods("GET")

	fmt.Println("Server is running on port 8081")
	log.Fatal(http.ListenAndServe(":8081", router))
}
