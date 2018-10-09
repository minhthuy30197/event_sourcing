package controller

import (
	"git.hocngay.com/hocngay/event-sourcing/config"
	"github.com/go-pg/pg"
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
