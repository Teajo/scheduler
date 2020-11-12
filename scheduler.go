package main

import (
	"jpb/scheduler/api"
	"jpb/scheduler/config"
	"jpb/scheduler/controller"
)

func main() {
	cfg := config.Get()
	ctrl := controller.New()
	api := api.NewHTTPApi(cfg.Port, ctrl)
	api.Listen()
}
