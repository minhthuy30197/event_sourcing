package main

import (
	"os"

	"git.hocngay.com/hocngay/event-sourcing/config"
	"git.hocngay.com/hocngay/event-sourcing/controller"
	"git.hocngay.com/hocngay/event-sourcing/model"
	"git.hocngay.com/hocngay/event-sourcing/router"
	"github.com/gin-gonic/gin"
)

// @title Service Auth
// @version 1.0
// @description Các API liên quan đến User.
// @host localhost:8080

func main() {
	controller.Register(model.AddTeacherEvent{})

	ginMode := os.Getenv("GIN_MODE")

	r := setupMiddleware(ginMode)
	config := config.SetupConfig()

	// Khởi tạo controller
	c := controller.NewController()
	c.Config = config

	// Kết nối CSDL
	dbConfig := config.Database
	db := model.ConnectDb(dbConfig.User, dbConfig.Password, dbConfig.Database, dbConfig.Address)
	defer db.Close()
	c.DB = db
	
	eventDbConfig := config.EventDatabase
	eventDb := model.ConnectDb(eventDbConfig.User, eventDbConfig.Password, eventDbConfig.Database, eventDbConfig.Address)
	defer eventDb.Close()
	c.EventDB = eventDb

	if ginMode != "release" {
		model.LogQueryToConsole(db)
	}

	err := model.MigrationDb(db, "course")
	if err != nil {
		panic(err)
	}

	err = model.MigrationEventDb(eventDb, config.ServiceName)
	if err != nil {
		panic(err)
	}

	router.SetupRouter(ginMode, config, r, c)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	err = r.Run(":" + port)
	if err != nil {
		panic(err)
	}
}

func setupMiddleware(ginMode string) *gin.Engine {
	r := gin.New()

	if ginMode == "release" {
		gin.DisableConsoleColor()
	} else {
		r.Use(gin.Logger())
	}

	r.Use(gin.Recovery())

	return r
}
