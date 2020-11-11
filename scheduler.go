package main

import (
	"jpb/scheduler/api"
	"jpb/scheduler/config"
)

func main() {
	cfg := config.Get()
	api := api.NewHTTPApi(cfg.Port)
	api.Listen()
}
