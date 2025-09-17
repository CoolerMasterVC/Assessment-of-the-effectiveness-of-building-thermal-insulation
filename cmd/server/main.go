package main

import (
	"Lab1/internal/api"
	"log"
)

func main() {
	log.Println("Application Start")
	api.StartServer()
	log.Println("Application terminated!")

}
