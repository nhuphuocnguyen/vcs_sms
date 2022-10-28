package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
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
	now := time.Now().Unix()
	server.Created_time = int(now)
	server.Last_updated = int(now)

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
	from, err := strconv.Atoi(c.Query("from"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error()})
		return
	}
	size, err := strconv.Atoi(c.Query("size"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error()})
		return
	}
	sort := c.Query("sort")
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
	array, err := serverDAO.Listserver(sort, from, size)
	// xu ly loi
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error()})
		return
	}
	// tra ket qua
	c.JSON(http.StatusOK, gin.H{"total": count, "data": array})
}

func UpdateServerHandler(c *gin.Context) {
	id, errs := strconv.Atoi(c.Param("server_id"))
	if errs != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
		return
	}
	var server models.Server
	now := time.Now().Unix()
	if err := c.ShouldBindJSON(&server); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
		return
	}
	serverDAO := daos.ServerDAO{Db: db}
	server.Last_updated = int(now)
	id, err := serverDAO.UpdateServer(server, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error()})
		return
	}
	result, err := serverDAO.Get(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, *result)

}

func DeleteServerHandler(c *gin.Context) {
	id, errs := strconv.Atoi(c.Param("server_id"))
	if errs != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
		return
	}
	var server models.Server
	serverDAO := daos.ServerDAO{Db: db}
	id, err := serverDAO.DeleteServer(server, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error()})
		return
	}
	server.Server_id = id
	c.JSON(http.StatusOK, server)

}

func ExportExcelHandle(c *gin.Context) {
	from, err := strconv.Atoi(c.Query("from"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error()})
		return
	}
	size, err := strconv.Atoi(c.Query("size"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error()})
		return
	}
	sort := c.Query("sort")
	f := excelize.NewFile()
	index := f.NewSheet("Sheet1")
	serverDAO := daos.ServerDAO{Db: db}

	var servers []models.Server

	servers, err = serverDAO.Listserver(sort, from, size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error()})
		return
	}
	for i, ctx := range servers {
		f.SetCellValue("Sheet1", "A"+strconv.Itoa(i+1), ctx.Server_id)
		f.SetCellValue("Sheet1", "B"+strconv.Itoa(i+1), ctx.Server_name)
		f.SetCellValue("Sheet1", "C"+strconv.Itoa(i+1), ctx.Status)
		f.SetCellValue("Sheet1", "E"+strconv.Itoa(i+1), ctx.Created_time)
		f.SetCellValue("Sheet1", "F"+strconv.Itoa(i+1), ctx.Last_updated)
		f.SetCellValue("Sheet1", "G"+strconv.Itoa(i+1), ctx.Ipv4)
	}
	f.SetActiveSheet(index)
	if err := f.SaveAs("New Server.xlsx"); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "file has been created successfully"})

}
func main() {
	connectDB()
	router := gin.Default()
	router.POST("/servers", NewServerHandler)
	router.GET("/servers", GetServerHandler)
	router.PUT("/servers/:server_id", UpdateServerHandler)
	router.DELETE("servers/:server_id", DeleteServerHandler)
	router.Run()
}
