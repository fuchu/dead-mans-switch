package main

import (
	"dms/controller"
	"dms/router"
	"flag"

	"dms/pkg/logger"
)

var alertManager *controller.AlertManager
var configPath = flag.String("config", "/etc/deadmansswitch/config.yaml", "Path to config.yaml.")

func init() {
	flag.Parse()
	var err error
	alertManager, err = controller.NewAlertManager(*configPath)
	if err != nil {
		logger.Log.Fatalln("Can not find config file!")
	}

	alertManager.StartTimers()
}

func main() {
	r := router.Router(alertManager)
	r.Run()
}
