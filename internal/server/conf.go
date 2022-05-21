package server

import (
	"go.uber.org/zap"
)

var logger *zap.SugaredLogger

type DiscordBotConfig struct {
	AuthToken string
	EnvName   string
	Logger    *zap.SugaredLogger
}
