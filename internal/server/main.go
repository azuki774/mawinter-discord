package server

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/azuki774/mawinter-discord/internal/client"
	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
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
		logger.Error("failed to start discord bot", zap.Error(err))
		return err
	}

	logger.Info("start discord bot")

	stopBot := make(chan os.Signal, 1)
	signal.Notify(stopBot, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)

	err = discord.Close()
	if err != nil {
		logger.Error("discord bot close error", zap.Error(err))
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
	logger.Info("userinfo loaded", zap.String("addr", os.Getenv("USER_MAWINTER_PATH")), zap.String("user", os.Getenv("USER_MAWINTER_USER")), zap.String("path", os.Getenv("USER_MAWINTER_PASS")),
		zap.String("discord_id", os.Getenv("USER_DISCORD_ID")), zap.String("name", os.Getenv("USER_DISCORD_NAME")))
}
