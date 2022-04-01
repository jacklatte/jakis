package main

import (
	"github.com/jackelatte/jakis"
)

func main() {
	addr := "0.0.0.0"
	port := 8080
	srv := jakis.NewServer(addr, port)
	srv.Run()
}
