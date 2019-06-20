package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

const (
	usage = "-prod to launch in production mode"
)

func main() {
	// Check if the test flag is perseent
	var dbHost string
	prodEnv := flag.Bool("prod", false, usage)
	flag.Parse()
	if *prodEnv == false {
		log.Println("[INFO] In Test env")
		dbHost = os.Getenv("MONGODB_TEST")
	} else {
		log.Println("[INFO] In Prod env")
		dbHost = os.Getenv("MONGODB_PROD")
	}
	http.Handle("/", handlers.Adapt(
		r,
		handlers.WithMailSender(sender),
		handlers.WithRepos(dbSession),
		handlers.Logging(),
		handlers.WithAPIVersion(),
	))
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal("[FATAL] ListenAndServe: ", err)
	}
}
