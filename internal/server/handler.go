package server

import (
	"errors"
	"strconv"

	"github.com/bwmarrin/discordgo"
)

const (
	categoryFood          = "🍏"
	categoryLife          = "🧺"
	categoryEntertainment = "🎲"
	categoryFriends       = "😀"
	categoryCompute       = "🖥️"
	categoryTrans         = "🚉"
	categoryStudy         = "📗"
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
	loggerIDName.Infow("post Message", "message", m.Content)

	switch user.Context {
	case ContextClosing:
		loggerIDName.Debugw("check now status", "status", ContextClosing)
		categoryChoicePost(s, m, user)
		user.changeCtxStatus(ContextOrderWaiting)
	case ContextOrderWaiting:
		loggerIDName.Debugw("check now status", "status", ContextOrderWaiting)
	case ContextPriceWaiting:
		loggerIDName.Debugw("check now status", "status", ContextPriceWaiting)
		if m.Content == "r" {
			user.changeCtxStatus(ContextOrderWaiting)
			categoryChoicePost(s, m, user)
		}
		// parse price
		_, err = strconv.Atoi(m.Content)
		if err != nil {
			loggerIDName.Warnw("invalid price value", "error", err)
			return
		}

		// POST to mawinter-server

		// TODO: real Response

		user.changeCtxStatus(ContextOrderWaiting)
		categoryChoicePost(s, m, user)
	}

}

func messageReaction(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	logger.Debugw("handler called", "name", "MessageReactionAdd", "userID", r.UserID)
	loggerID := logger.With("name", "MessageReactionAdd", "userID", r.UserID, "emoji", r.Emoji)
	if r.UserID == s.State.User.ID {
		return
	}

	user, err := users.getUserInfoByID(r.UserID)
	if errors.Is(err, ErrUserNotFound) {
		loggerID.Warn("unknown user")
	} else if err != nil {
		loggerID.Errorw("unknown error", "error", err)
	}

	loggerIDName := logger.With("username", user.Name)
	switch r.Emoji.Name {
	case categoryFood:
		loggerIDName.Infow("category choose", "category", "food")
	case categoryLife:
		loggerIDName.Infow("category choose", "category", "life")
	case categoryEntertainment:
		loggerIDName.Infow("category choose", "category", "entertainment")
	case categoryFriends:
		loggerIDName.Infow("category choose", "category", "friends")
	case categoryCompute:
		loggerIDName.Infow("category choose", "category", "compute")
	case categoryTrans:
		loggerIDName.Infow("category choose", "category", "trans")
	case categoryStudy:
		loggerIDName.Infow("category choose", "category", "study")
	default:
		loggerIDName.Infow("unknown choose", "category", r.Emoji.Name)
		return
	}

	user.changeCtxStatus(ContextPriceWaiting)
	_, err = s.ChannelMessageSend(r.ChannelID, "you choose "+r.Emoji.Name+"\n"+"what value?")
	if err != nil {
		logger.Errorw("failed to send message", "error", err)
		return
	}
}

func categoryChoicePost(s *discordgo.Session, m *discordgo.MessageCreate, user *discordUser) {
	mes, err := s.ChannelMessageSend(m.ChannelID, "choose category")
	if err != nil {
		logger.Errorw("failed to send message", "error", err)
		return
	}
	chanID := mes.ChannelID
	mesID := mes.ID
	s.MessageReactionAdd(chanID, mesID, categoryFood)
	s.MessageReactionAdd(chanID, mesID, categoryLife)
	s.MessageReactionAdd(chanID, mesID, categoryEntertainment)
	s.MessageReactionAdd(chanID, mesID, categoryFriends)
	s.MessageReactionAdd(chanID, mesID, categoryCompute)
	s.MessageReactionAdd(chanID, mesID, categoryTrans)
	s.MessageReactionAdd(chanID, mesID, categoryStudy)
}
