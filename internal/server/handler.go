package server

import (
	"errors"
	"fmt"
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

	categoryIDFood          = 210
	categoryIDLife          = 240
	categoryIDEntertainment = 250
	categoryIDFriends       = 251
	categoryIDCompute       = 230
	categoryIDTrans         = 270
	categoryIDStudy         = 260
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

	var price int
	switch user.Context {
	case ContextClosing:
		loggerIDName.Debugw("check now status", "status", ContextClosing)
		user.changeCtxStatus(ContextOrderWaiting)
		categoryChoicePost(s, m, user)
	case ContextOrderWaiting:
		loggerIDName.Debugw("check now status", "status", ContextOrderWaiting)
	case ContextPriceWaiting:
		loggerIDName.Debugw("check now status", "status", ContextPriceWaiting)
		if m.Content == "r" {
			user.changeCtxStatus(ContextOrderWaiting)
			categoryChoicePost(s, m, user)
		}
		// parse price
		price, err = strconv.Atoi(m.Content)
		if err != nil {
			loggerIDName.Warnw("invalid price value", "error", err)
			_, err = s.ChannelMessageSend(m.ChannelID, "invalid value"+"\n"+"what value? (type 'r' to back to category select.)")
			if err != nil {
				logger.Errorw("failed to send message", "error", err)
				return
			}
			return
		}

		// POST to mawinter-server
		res, err := clientrepo.PostMawinter(&user.ServerInfo, user.ChooseCategoryID, int64(price))
		if err != nil {
			// TODO: invalid data or connection lost かで分ける
			logger.Errorw("failed to send order to mawinter-server", "error", err)
			return
		}

		logger.Infow("receive response from mawinter", "response", *res, "addr", user.ServerInfo.Addr)
		resText := fmt.Sprintf("ID: %v, categoryID: %v, Price: %v", res.Id, res.CategoryId, res.Price)
		_, err = s.ChannelMessageSend(m.ChannelID, resText)
		if err != nil {
			logger.Errorw("failed to send message", "error", err)
			return
		}

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
		user.ChooseCategoryID = categoryIDFood
	case categoryLife:
		loggerIDName.Infow("category choose", "category", "life")
		user.ChooseCategoryID = categoryIDLife
	case categoryEntertainment:
		loggerIDName.Infow("category choose", "category", "entertainment")
		user.ChooseCategoryID = categoryIDEntertainment
	case categoryFriends:
		loggerIDName.Infow("category choose", "category", "friends")
		user.ChooseCategoryID = categoryIDFriends
	case categoryCompute:
		loggerIDName.Infow("category choose", "category", "compute")
		user.ChooseCategoryID = categoryIDCompute
	case categoryTrans:
		loggerIDName.Infow("category choose", "category", "trans")
		user.ChooseCategoryID = categoryIDTrans
	case categoryStudy:
		loggerIDName.Infow("category choose", "category", "study")
		user.ChooseCategoryID = categoryIDStudy
	default:
		loggerIDName.Infow("unknown choose", "category", r.Emoji.Name)
		user.ChooseCategoryID = -1
		return
	}

	user.changeCtxStatus(ContextPriceWaiting)
	_, err = s.ChannelMessageSend(r.ChannelID, "you choose "+r.Emoji.Name+"\n"+"what value? (type 'r' to back to category select.)")
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
