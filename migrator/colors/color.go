package main

import (
	"fmt"
	"github.com/softdev9/trendee-api-master/Godeps/_workspace/src/gopkg.in/mgo.v2"
	"github.com/softdev9/trendee-api-master/repos"
	"os"
)

func main() {
	fmt.Println("Are we good for inserting colors")

	colorToInsert := []repos.ColorRecord{
		repos.ColorRecord{HexCode: "#ff0000", Name: "red", TranslationFr: "rouge"},
		repos.ColorRecord{HexCode: "#00ff00", Name: "green", TranslationFr: "vert"},
		repos.ColorRecord{HexCode: "#0000ff", Name: "blue", TranslationFr: "bleu"},
		repos.ColorRecord{HexCode: "#ffff00", Name: "yellow", TranslationFr: "jaune"},
		repos.ColorRecord{HexCode: "#ff00ff", Name: "purple", TranslationFr: "violet"},
		repos.ColorRecord{HexCode: "#00ffff", Name: "magenta", TranslationFr: "magenta"},
		repos.ColorRecord{HexCode: "#000000", Name: "black", TranslationFr: "noir"},
		repos.ColorRecord{HexCode: "#ffffff", Name: "white", TranslationFr: "blanc"},
	}

	// Connect to the database
	dbHost := os.Getenv("MONGODB_PROD")
	fmt.Println("Connecting to ", dbHost)
	dbConn, err := mgo.Dial(dbHost)
	if err != nil {
		fmt.Println("Not able to connect to the db", err.Error())
		return
	}
	defer dbConn.Close()
	// Create the category repo
	repo := repos.NewColorRepo(dbConn)
	for _, c := range colorToInsert {
		repo.SaveColor(&c)
	}
	fmt.Println("Color inserted")
}
