package server

import (
	"os"

	"go.uber.org/zap"
)

var logger *zap.SugaredLogger
var runEnviroment string

func GetSugaredLogger() (*zap.SugaredLogger, error) {
	logger, err := zap.NewDevelopment()
	defer logger.Sync()
	if err != nil {
		return nil, err
	}
	sugarLogger := logger.Sugar()
	return sugarLogger, nil
}

func GetEnviroment() (discordAuthToken string, err error) {
	if v, ok := os.LookupEnv("AUTH_TOKEN"); ok {
		discordAuthToken = v
	} else {
		logger.Error("Discord authentication token not found")
		return "", err
	}

	if v, ok := os.LookupEnv("RUNENV"); ok {
		runEnviroment = v
	} else {
		runEnviroment = "undefined"
		logger.Warn("Run environment is undefined")
	}

	return discordAuthToken, nil
}
