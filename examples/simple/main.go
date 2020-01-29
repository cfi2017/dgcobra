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

	// create and setup new handler
	handler = dgcobra.NewHandler(session)
	// add simple global prefix
	handler.AddPrefix("!")
	// set command factory
	handler.RootFactory = cmd.RootCmdFactory
	// register a prefix lookup function
	handler.PrefixFunc = func(session *discordgo.Session, event *discordgo.MessageCreate) []string {
		if event.GuildID != "" && cmd.Prefixes[event.GuildID] != "" {
			return []string{cmd.Prefixes[event.GuildID], fmt.Sprintf("<@!%s> ", session.State.User.ID)}
		}
		return nil
	}
	// register new handler with discordgo
	handler.Start()

	// add ready handler to add bot mention when ID is available
	session.AddHandlerOnce(onReady)

	// open session
	err = session.Open()
	if err != nil {
		panic(err)
	}

	// cleanup
	log.Println("Bot is running. Press CTRL-C to exit.")
	waitForSig() // wait for termination signal
	err = session.Close()
	if err != nil {
		panic(err)
	}
}

func onReady(_ *discordgo.Session, ready *discordgo.Ready) {
	// register bot mention as new global prefix
	handler.AddPrefix(fmt.Sprintf("<@!%s> ", ready.User.ID))
}

func waitForSig() {
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}
