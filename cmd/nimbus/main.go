package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/vikraj01/terraplay/internals/bot"
	"github.com/vikraj01/terraplay/internals/webhook"
)

func main() {
	go func(){
		bot.StartBot()
	}()

	go func(){
		http.HandleFunc("/webhook", webhook.HandleWebhook)
		log.Println("Webhook server listening on port: 8080")
		err := http.ListenAndServe(":8080", nil)
		if err != nil {
			log.Fatalf("Webhook server failed: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	log.Println("Shutting down both Discord bot and webhook listener...")
}
