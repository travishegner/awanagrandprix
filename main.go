package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/travishegner/awanagrandprix/dashboard"
)

func main() {
	db, err := dashboard.NewDashboard()
	if err != nil {
		log.WithError(err).Fatal("Failed to create dashboard")
	}

	err = db.Start()
	if err != nil {
		log.WithError(err).Error("Error while starting the dashboard")
	}
}
