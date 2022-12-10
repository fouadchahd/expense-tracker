package main

import (
	"github.com/fouadelhamri/expense-tracker/controller"
	"github.com/fouadelhamri/expense-tracker/model"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strings"
)

var db *gorm.DB

const DateFormat = "2006-01-02"

func getDns() string {
	var envs map[string]string
	envs, err := godotenv.Read(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	username, ok1 := envs["SERVER_USERNAME"]
	pass, ok2 := envs["SERVER_PASS"]
	dbName, ok3 := envs["DB_NAME"]

	if !(ok1 && ok2 && ok3) {
		log.Printf(" %v %v %v", username, pass, dbName)
		log.Fatal("Missing Env variable to start")
	}
	var dns strings.Builder
	dns.WriteString(username)
	dns.WriteString(":" + pass)
	dns.WriteString("@tcp(localhost:8889)/")
	dns.WriteString(dbName)
	dns.WriteString("?charset=utf8mb4&parseTime=True&loc=Local")
	return dns.String()
}
func ConnectToDb(connectionString string) *gorm.DB {
	db, err := gorm.Open(mysql.Open(connectionString), &gorm.Config{})
	if err != nil {
		panic("can't connect to database")
	}
	return db
}

func init() {
	var Dns = getDns()
	db = ConnectToDb(Dns)
}

func main() {
	r := mux.NewRouter().StrictSlash(true)
	//					Migrate					//
	db.AutoMigrate(&model.User{})
	db.AutoMigrate(&model.Group{})
	db.AutoMigrate(&model.MetaItem{})
	db.AutoMigrate(&model.Category{})
	db.AutoMigrate(&model.Transaction{})
	//					Routes					//

	//				  User Routes				//
	userController := controller.NewUserController(db)
	r.HandleFunc("/register", userController.RegisterController).Methods(http.MethodPost).Name("registering")
	r.HandleFunc("/users/{id}/{password}", userController.CheckAuthorization).Methods(http.MethodGet).Name("check_authorization")
	r.HandleFunc("/users/{id}", userController.GetSingleUser).Methods(http.MethodGet).Name("get_user_by_ID")
	r.HandleFunc("/login", userController.LoginController).Methods(http.MethodPost)
	//				  Category Routes				//
	categoryController := controller.NewCategoryController(db)
	r.HandleFunc("/categories", categoryController.CreateCategory).Methods(http.MethodPost)
	r.HandleFunc("/categories", categoryController.GetCategories).Methods(http.MethodGet)
	//				  Transaction Routes				//
	transactionController := controller.NewTransactionController(db)
	r.HandleFunc("/transactions", transactionController.CreateTransaction).Methods(http.MethodPost)
	r.HandleFunc("/transactions", transactionController.GetTransactionsByUser).Methods(http.MethodGet)
	r.HandleFunc("/transactions", transactionController.GetTransactionsByUser).Queries("year", "{year:[0-9]{4}}").Methods(http.MethodGet)
	r.HandleFunc("/transactions", transactionController.GetTransactionsByUser).Queries("year", "{year:[0-9]{4}}", "month", "{month:[0-9]{1,2}}").Methods(http.MethodGet)
	r.HandleFunc("/transactions", transactionController.GetTransactionsByUser).Queries("year", "{year:[0-9]{4}}", "month", "{month:[0-9]{1,2}}", "day", "{day:[0-9]{1,2}}").Methods(http.MethodGet)
	r.HandleFunc("/groups/{id:[0-9a-zA-Z]+}/transactions", transactionController.GetTransactionGroup).Methods(http.MethodGet)
	//				  Group Routes				//
	groupController := controller.NewGroupController(db)
	r.HandleFunc("/groups/{id:[0-9a-zA-Z]+}/join", groupController.JoinGroup).Methods(http.MethodPost)
	r.HandleFunc("/groups/{id:[0-9a-zA-Z]+}/left", groupController.LeftGroup).Methods(http.MethodPost)
	r.HandleFunc("/groups", groupController.CreateGroupByUser).Methods(http.MethodPost)

	//				  Init Server				//
	log.Println("Listening On :8000 ...")
	log.Fatal(http.ListenAndServe(":8000", r))
}
