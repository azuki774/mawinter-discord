package server

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/azuki774/mawinter-discord/internal/client"
	"github.com/bwmarrin/discordgo"
)

var users *discordUsers
var envName string

func Start(botConfig *DiscordBotConfig) (err error) {
	logger = botConfig.Logger
	envName = botConfig.EnvName
	client.Logger = logger
	clientrepo = botConfig.MawinterClient
	users = &discordUsers{}

	// TODO: add multi user system
	recordUserInfo()

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

func recordUserInfo() {
	users.addDiscordUser(
		client.ServerInfo{
			Addr: os.Getenv("USER_MAWINTER_PATH"),
			User: os.Getenv("USER_MAWINTER_USER"),
			Pass: os.Getenv("USER_MAWINTER_PASS")},
		os.Getenv("USER_DISCORD_ID"),
		os.Getenv("USER_DISCORD_NAME"))
	logger.Infow("userinfo loaded", "addr", os.Getenv("USER_MAWINTER_PATH"), "user", os.Getenv("USER_MAWINTER_USER"), "path", os.Getenv("USER_MAWINTER_PASS"),
		"discord_id", os.Getenv("USER_DISCORD_ID"), "name", os.Getenv("USER_DISCORD_NAME"))
}
