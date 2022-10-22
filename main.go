package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nhuphuocnguyen/vcs_sms/daos"
	"github.com/nhuphuocnguyen/vcs_sms/models"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "25122002"
	dbname   = "vcs-sms"
)

var db *sql.DB
var err error

func connectDB() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	print(psqlInfo)

	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		println("Database errors")
		return
	}
}

func NewServerHandler(c *gin.Context) {
	// Lay thong tin server gui len
	var server models.Server
	if err := c.ShouldBindJSON(&server); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
		return
	}

	// luu vao database
	serverDAO := daos.ServerDAO{Db: db}
	id, err := serverDAO.CreateServer(server)

	// neu co loi => 500
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error()})
		return
	}
	// khong loi => tra ve thong tin server da luu
	server.Server_id = id
	c.JSON(http.StatusOK, server)
}

func GetServerHandler(c *gin.Context) {
	// var newserver daos.ServerDAO
	serverDAO := daos.ServerDAO{Db: db}

	// dem so luong server tu database
	count, err := serverDAO.Count()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error()})
		return
	}
	// lay danh sach server tu database
	array, err := serverDAO.Listserver()
	// xu ly loi
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error()})
		return
	}
	// tra ket qua
	c.JSON(http.StatusOK, gin.H{"total": count, "data": array})
}

func main() {
	connectDB()
	router := gin.Default()
	router.POST("/servers", NewServerHandler)
	router.GET("/servers", GetServerHandler)
	router.Run()
}
