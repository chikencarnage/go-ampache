package main

import (
	"log"

	ampache "github.com/chikencarnage/go-ampache"
)

func main() {
	username := ""
	password := ""
	ampacheHost := "http://localhost"

	client := ampache.NewConnection(ampacheHost)
	err := client.PasswordAuth(username, password)
	if err != nil {
		log.Printf("Error connecting to server: %s", err)
		return
	}

	log.Printf("Successfully connected")
}
