package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/dmfed/conf"
	"gopkg.in/tucnak/telebot.v2"
)

func main() {
	var (
		flagConfigFile = flag.String("c", "/usr/local/etc/pesik.conf", "configuration file to use")
	)
	flag.Parse()
	cfg, err := conf.ParseFile(*flagConfigFile)
	if err != nil {
		log.Println("could not read config file:", *flagConfigFile)
		return
	}
	token := cfg.Get("token").String()
	settings := telebot.Settings{Token: token,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second}}

	bot, err := telebot.NewBot(settings)
	if err != nil {
		log.Println(err)
		return
	}
	id, _ := cfg.Get("ownerid").Int()
	idstr := cfg.Get("ownerid").String()
	username := cfg.Get("ownerusername").String()
	pesik := pesikobot{bot, botOwner{id, idstr, username}}
	pesik.Handle(telebot.OnText, pesik.all)
	// Let's handle system signals
	interrupts := make(chan os.Signal, 1)
	signal.Notify(interrupts, os.Interrupt, os.Kill)
	go func() {
		sig := <-interrupts
		log.Printf("exiting on signal: %v", sig)
		pesik.Stop()
	}()
	// Actually starting the bot
	pesik.Start()
}
