package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/demonwolfdev/community-bot/pkg/discord"
	twitchgo "github.com/gempir/go-twitch-irc/v2"
	"github.com/joho/godotenv"
)

var discordClient, _ = discordgo.New()

func init() {
	godotenv.Load()
}

func main() {
	discordClient.Token = os.Getenv("DISCORD_TOKEN")

	// Discord CLient
	discordClient.AddHandler(discord.Ready)
	discordClient.AddHandler(discord.MessageCreate)
	discordClient.Identify.Intents = discordgo.IntentsGuildMessages

	err := discordClient.Open()
	if err != nil {
		log.Printf("error opening connection to Discord, %s\n", err)
		os.Exit(1)
	}
	twitchClient := twitchgo.NewClient(os.Getenv("TWITCH_USERNAME"), os.Getenv("TWITCH_OAUTH"))

	twitchClient.Join(os.Getenv("TWITCH_CHANNEL"))

	twitchClient.OnPrivateMessage(func(message twitchgo.PrivateMessage) {
	})

	err = twitchClient.Connect()
	if err != nil {
		panic(err)
	}

	log.Printf(`Now running. Press CTRL-C to exit.`)
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Clean up
	discordClient.Close()
}
