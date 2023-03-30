package main

import (
	"log"
	"weather"
)

func init() {
	log.SetFlags(0)
}

func main() {
	weather.RunCLI()
}
