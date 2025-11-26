package main

import (
	"Blog/internal/config"
	"fmt"
	"log"
)

func main() {
	cnfg, err := config.Read()
	if err != nil {
		log.Fatalf("Error reading config: %v", err)
	}
	cnfg.SetUser("Antonios")
	cnfg, err = config.Read()
	if err != nil {
		log.Fatalf("Error reading config: %v", err)
	}
	fmt.Println("Current User:", cnfg.CurrentUserName)
	fmt.Println("DB URL:", cnfg.DBURL)
}

//go run . This runs the program
// go build . This builds the program into an executable file
