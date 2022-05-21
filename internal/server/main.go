package server

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

func Start() (err error) {
	logger, err = GetSugaredLogger()
	if err != nil {
		return err
	}

	logger.Info("Set logger")

	token, err := GetEnviroment()
	if err != nil {
		logger.Error("Get environment error")
		return err
	}

	logger.Infow("Run enviromenment", "env", runEnviroment)

	if err = runServer(token); err != nil {
		return err
	}

	return nil
}

func runServer(token string) (err error) {
	discord, err := discordgo.New("Bot " + token)
	if err != nil {
		return err
	}
	discord.AddHandler(messageCreate)
	discord.Identify.Intents = discordgo.IntentsGuildMessages

	err = discord.Open()
	if err != nil {
		logger.Errorw("failed to start discord bot", "error", err)
		return err
	}

	logger.Info("start discord bot")

	stopBot := make(chan os.Signal, 1)
	signal.Notify(stopBot, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	logger.Infow("catch stop signal", "signal", <-stopBot)

	err = discord.Close()
	if err != nil {
		logger.Errorw("discord bot close error", "error", err)
		return err
	}
	return nil
}
