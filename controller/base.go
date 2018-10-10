package controller

import (
	"github.com/go-pg/pg"
	"github.com/minhthuy30197/event_sourcing/config"
)

type Controller struct {
	// DB instance
	DB *pg.DB

	// Event DB instance
	EventDB *pg.DB

	// Configuration
	Config config.Config
}

func NewController() *Controller {
	var c Controller
	return &c
}
