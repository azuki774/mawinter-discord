package main

import (
	"os"

	"github.com/azuki774/mawinter-discord/internal/client"
	"github.com/azuki774/mawinter-discord/internal/server"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var logger *zap.SugaredLogger

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
	logger, err := GetSugaredLogger()
	if err != nil {
		return err
	}

	logger.Info("set logger")

	botConfig, err := GetEnviroment()
	if err != nil {
		logger.Error("get environment error")
		return err
	}

	logger.Infow("Run enviromenment", "env", botConfig.EnvName)
	botConfig.Logger = logger

	if err = server.Start(botConfig); err != nil {
		return err
	}

	return nil
}

func GetSugaredLogger() (*zap.SugaredLogger, error) {
	logger, err := zap.NewDevelopment()
	defer logger.Sync()
	if err != nil {
		return nil, err
	}
	sugarLogger := logger.Sugar()
	return sugarLogger, nil
}

func GetEnviroment() (botConfig *server.DiscordBotConfig, err error) {
	botConfig = &server.DiscordBotConfig{}
	if v, ok := os.LookupEnv("AUTH_TOKEN"); ok {
		botConfig.AuthToken = v
	} else {
		logger.Error("discord authentication token not found")
		return nil, err
	}

	if v, ok := os.LookupEnv("RUNENV"); ok {
		botConfig.EnvName = v
	} else {
		botConfig.EnvName = "undefined"
		logger.Warn("run environment is undefined")
	}

	// TODO : USE_MOCK Optionの処理

	botConfig.MawinterClient = client.NewClientRepo()
	return botConfig, nil
}
