package router

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/minhthuy30197/event_sourcing/config"
	"github.com/minhthuy30197/event_sourcing/controller"
)

type JsonDataRoute struct {
	Tags        []string `json:"tags"`
	Summary     string   `json:"summary"`
	Description string   `json:"description"`
}

type JsonData struct {
	Paths map[string]map[string]JsonDataRoute
}

type Route struct {
	Service     string
	Endpoint    string
	HttpMethod  string
	Description string
	IsPublic    bool
	IsAdmin     bool
}

type Message struct {
	Id     string
	Routes []Route
	Time   time.Time
}

func SetupRouter(ginMode string, config config.Config, r *gin.Engine, c *controller.Controller) {
	// Mọi routes đều bắt đầu bởi prefix ServiceName
	api := r.Group("/" + config.ServiceName)
	{
		setupCourseRoutes(c, api)
	}
}
