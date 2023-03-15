package main

import (
	"os"

	"github.com/azuki774/mawinter-discord/internal/client"
	"github.com/azuki774/mawinter-discord/internal/server"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var logger *zap.Logger

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		return runServer()
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func runServer() (err error) {
	logger, err = GetLogger()
	if err != nil {
		return err
	}
	defer logger.Sync()
	logger.Info("set logger")

	botConfig, err := GetEnviroment()
	if err != nil {
		logger.Error("get environment error")
		return err
	}

	logger.Info("Run enviromenment", zap.String("env", botConfig.EnvName))
	botConfig.Logger = logger

	if err = server.Start(botConfig); err != nil {
		return err
	}

	return nil
}

func GetLogger() (*zap.Logger, error) {
	lg, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}
	return lg, nil
}

func GetEnviroment() (botConfig *server.DiscordBotConfig, err error) {
	botConfig = &server.DiscordBotConfig{}
	if v, ok := os.LookupEnv("AUTH_TOKEN"); ok {
		botConfig.AuthToken = v
		logger.Debug("use token", zap.String("token", v))
	} else {
		logger.Error("discord authentication token not found", zap.Error(err))
		return nil, err
	}

	if v, ok := os.LookupEnv("RUNENV"); ok {
		botConfig.EnvName = v
	} else {
		botConfig.EnvName = "undefined"
		logger.Warn("run environment is undefined")
	}

	if os.Getenv("USE_MOCK") == "0" {
		logger.Info("use mawinter server")
		botConfig.MawinterClient = client.NewClientRepo()
	} else {
		logger.Info("use mawinter stub")
		botConfig.MawinterClient = client.NewMockClientRepo()
	}

	return botConfig, nil
}
