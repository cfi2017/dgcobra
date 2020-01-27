package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/cfi2017/dgcobra"
	"github.com/cfi2017/dgcobra/examples/simple/cmd"
)

var token = flag.String("token", "", "bot token")
var session *discordgo.Session
var handler *dgcobra.Handler

func main() {
	flag.Parse()
	if *token == "" {
		log.Fatal("missing token")
	}

	var err error
	session, err = discordgo.New("Bot " + *token)
	if err != nil {
		panic(err)
	}

	handler = dgcobra.NewHandler(session)

	handler.AddPrefix("!")
	handler.RootFactory = cmd.RootCmdFactory
	handler.Start()

	session.AddHandlerOnce(onReady)

	err = session.Open()
	if err != nil {
		panic(err)
	}

	log.Println("Bot is running. Press CTRL-C to exit.")
	waitForSig() // wait for termination signal

	err = session.Close()
	if err != nil {
		panic(err)
	}
}

func onReady(_ *discordgo.Session, ready *discordgo.Ready) {
	handler.AddPrefix(fmt.Sprintf("<@!%s> ", ready.User.ID))
}

func waitForSig() {
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}
