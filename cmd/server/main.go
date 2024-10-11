package main

import (
	"log"

	"github.com/vikraj01/terraplay/internals/server"
)


func main(){
	done := make(chan bool)

    go func() {
		server.StartServer()
        done <- true 
    }()

    <-done
    log.Println("CLI has exited. Performing cleanup.")
}