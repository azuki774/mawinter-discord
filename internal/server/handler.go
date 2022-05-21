package server

import (
	"errors"

	"github.com/bwmarrin/discordgo"
)

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	loggerID := logger.With("name", "MessageCreate", "userID", m.Author.ID)
	loggerID.Debugw("handler called")
	if m.Author.ID == s.State.User.ID {
		return
	}

	user, err := users.getUserInfoByID(m.Author.ID)
	if errors.Is(err, ErrUserNotFound) {
		loggerID.Warn("unknown user")
	} else if err != nil {
		loggerID.Errorw("unknown error", "error", err)
	}
	loggerIDName := logger.With("username", user.Name)
	loggerIDName.Infow("post Message")

	switch user.Context {
	case ContextClosing:
		loggerIDName.Debugw("check now status", "status", ContextClosing)
		closingHandle(s, m, user)
	case ContextOrderWaiting:
		loggerIDName.Debugw("check now status", "status", ContextOrderWaiting)
		orderWaitingHandle(s, m, user)
	case ContextPriceWaiting:
		loggerIDName.Debugw("check now status", "status", ContextPriceWaiting)
		priceWaitingHandle(s, m, user)
	}

}

func messageReaction(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	logger.Debugw("handler called", "name", "MessageReactionAdd", "userID", r.UserID)
	if r.UserID == s.State.User.ID {
		return
	}
}

func closingHandle(s *discordgo.Session, m *discordgo.MessageCreate, user *discordUser) {
	mes, err := s.ChannelMessageSend(m.ChannelID, "start")
	if err != nil {
		logger.Errorw("failed to send message", "error", err)
		return
	}
	chanID := mes.ChannelID
	mesID := mes.ID
	emoji := "❤️‍🔥"
	err = s.MessageReactionAdd(chanID, mesID, emoji)
	if err != nil {
		logger.Errorw("failed to add reaction", "error", err)
		return
	}

	user.changeCtxStatus(ContextOrderWaiting)
}

func orderWaitingHandle(s *discordgo.Session, m *discordgo.MessageCreate, user *discordUser) {

}

func priceWaitingHandle(s *discordgo.Session, m *discordgo.MessageCreate, user *discordUser) {

}
