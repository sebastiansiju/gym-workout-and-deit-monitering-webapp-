package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/Cawlumm/sebu-backend/config"
	"github.com/Cawlumm/sebu-backend/controllers"
	"github.com/Cawlumm/sebu-backend/db"
	"github.com/Cawlumm/sebu-backend/routes"
	"github.com/Cawlumm/sebu-backend/seed"
	"github.com/Cawlumm/sebu-backend/stores"
	"github.com/gin-gonic/gin"
)

func main() {
	showVersion := flag.Bool("version", false, "print the build version and exit")
	flag.Parse()
	if *showVersion {
		fmt.Printf("sebu %s\n", config.Version())
		os.Exit(0)
	}

	config.Load()
	db.Connect()
	seed.DemoUser(db.DB)
	seed.Exercises(db.DB)
	go seed.DemoData(db.DB)

	if config.C.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()
	s := stores.New(db.DB)
	h := controllers.NewHandler(s)
	routes.Setup(r, h)

	addr := ":" + config.C.Port
	log.Printf("sebu API listening on %s (env=%s)", addr, config.C.Env)
	if err := r.Run(addr); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
