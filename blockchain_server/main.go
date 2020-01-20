package main

import (
	"flag"
	"log"
)

func init() {
	log.SetPrefix("Blockchain: ")
}

func main() {
	port := flag.Uint("port", 5000, "TCP port  Number for Blockchain server.")
	flag.Parse()
	app := NewBlockchainserver(uint16(*port))
	app.Run()
}
