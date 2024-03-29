package server

import (
	"github.com/azuki774/mawinter-discord/internal/client"
	"go.uber.org/zap"
)

var logger *zap.Logger
var clientrepo client.ClientRepository

type DiscordBotConfig struct {
	AuthToken      string
	EnvName        string
	MawinterClient client.ClientRepository
	Logger         *zap.Logger
}
