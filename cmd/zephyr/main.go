package main

import (
    "log"

    "github.com/vikraj01/terraplay/internals/cli"
)

func main() {
    done := make(chan bool)

    go func() {
        cli.StartCli()
        done <- true 
    }()

    <-done
    log.Println("CLI has exited. Performing cleanup.")
}
