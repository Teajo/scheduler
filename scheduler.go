package main

import (
	"fmt"
	"jpb/scheduler/api"
	"jpb/scheduler/config"
	"jpb/scheduler/controller"
	"jpb/scheduler/db"
	"jpb/scheduler/events"
	"jpb/scheduler/publisher"
	"jpb/scheduler/taskqueue"
	"jpb/scheduler/ticker"
	"os"

	"github.com/facebookgo/inject"
)

func main() {
	var g inject.Graph

	var ctrl controller.Ctrl
	var api api.HTTPApi
	var ticker ticker.Ticker
	var pm publisher.PubManager
	var tq taskqueue.TaskQueue
	cfg := config.Get()
	db := db.Getdb(cfg.DbDriver)
	bus := events.New()

	err := g.Provide(
		&inject.Object{Value: db},
		&inject.Object{Value: cfg},
		&inject.Object{Value: bus},
		&inject.Object{Value: &ctrl},
		&inject.Object{Value: &api},
		&inject.Object{Value: &ticker},
		&inject.Object{Value: &pm},
		&inject.Object{Value: &tq},
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if err := g.Populate(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	tq.Start()
	pm.Start()
	ticker.Start(cfg.TimeChunk)
	api.Listen()
}
