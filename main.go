package main

import (
	asciiartweb "asciiartweb/backend/server"
	"log"
)

func main() {
	err := asciiartweb.Server()
	log.Fatal(err)
}
