package router

import (
	"git.hocngay.com/hocngay/event-sourcing/controller"
	"github.com/gin-gonic/gin"
)

func setupCourseRoutes(c *controller.Controller, api *gin.RouterGroup) {
	api.POST("/add-teacher", c.AddTeacherToClass)
	api.GET("/play-back/:id", c.Playback)
	api.PUT("/remove-teacher", c.RemoveTeacherFromClass)
	api.GET("/get-teacher/:id", c.GetTeachersOfClass)
	api.GET("/get-history/:id", c.GetHistory)
}