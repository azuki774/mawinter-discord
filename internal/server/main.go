package server

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var users *discordUsers

func Start(botConfig *DiscordBotConfig) (err error) {
	logger = botConfig.Logger
	clientrepo = botConfig.MawinterClient
	users = &discordUsers{}

	// TODO: add user system
	users.addDiscordUser("288297568369246208", "azuki")

	discord, err := discordgo.New("Bot " + botConfig.AuthToken)
	if err != nil {
		return err
	}
	discord.AddHandler(messageCreate)
	discord.AddHandler(messageReaction)
	discord.Identify.Intents = discordgo.IntentsDirectMessages | discordgo.IntentsDirectMessageReactions

	err = discord.Open() // Connect Discord
	if err != nil {
		logger.Errorw("failed to start discord bot", "error", err)
		return err
	}

	logger.Info("start discord bot")

	stopBot := make(chan os.Signal, 1)
	signal.Notify(stopBot, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	logger.Infow("catch stop signal", "signal", <-stopBot) // blocking

	err = discord.Close()
	if err != nil {
		logger.Errorw("discord bot close error", "error", err)
		return err
	}
	return nil
}
