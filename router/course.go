package router

import (
	"github.com/gin-gonic/gin"
	"github.com/minhthuy30197/event_sourcing/controller"
)

func setupCourseRoutes(c *controller.Controller, api *gin.RouterGroup) {
	api.POST("/add-teacher", c.AddTeacherToClass)
	api.GET("/play-back/:id", c.Playback)
	api.PUT("/remove-teacher", c.RemoveTeacherFromClass)
	api.GET("/get-teacher/:id", c.GetTeachersOfClass)
	api.POST("/get-history/:id", c.GetHistory)
}
