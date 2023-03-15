package server

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
)

const (
	categoryFood          = "🍏"
	categoryLife          = "🧺"
	categoryEntertainment = "🎲"
	categoryFriends       = "😀"
	categoryCompute       = "🖥️"
	categoryTrans         = "🚉"
	categoryStudy         = "📗"

	operateCancel = "❌"

	categoryIDFood          = 210
	categoryIDLife          = 240
	categoryIDEntertainment = 250
	categoryIDFriends       = 251
	categoryIDCompute       = 230
	categoryIDTrans         = 270
	categoryIDStudy         = 260
)

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	loggerID := logger.With(zap.String("name", "MessageCreate"), zap.String("userID", m.Author.ID))
	loggerID.Debug("handler called")
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == "env" { // show environment name
		_, err := s.ChannelMessageSend(m.ChannelID, envName)
		if err != nil {
			logger.Error("failed to send message", zap.Error(err))
		}
		return
	}

	user, err := users.getUserInfoByID(m.Author.ID)
	if errors.Is(err, ErrUserNotFound) {
		loggerID.Warn("unknown user")
	} else if err != nil {
		loggerID.Error("unknown error", zap.Error(err))
	}
	loggerIDName := logger.With(zap.String("username", user.Name))
	loggerIDName.Info("post Message", zap.String("message", m.Content))

	var price int
	switch user.Context {
	case ContextClosing:
		user.changeCtxStatus(ContextOrderWaiting)
		categoryChoicePost(s, m, user)
	case ContextOrderWaiting:

	case ContextPriceWaiting:

		if m.Content == "r" {
			user.changeCtxStatus(ContextOrderWaiting)
			categoryChoicePost(s, m, user)
			return
		}
		// parse price
		price, err = strconv.Atoi(m.Content)
		if err != nil {
			loggerIDName.Warn("invalid price value", zap.Error(err))
			_, err = s.ChannelMessageSend(m.ChannelID, "invalid value"+"\n"+"what value? (type 'r' to back to category select.)")
			if err != nil {
				logger.Error("failed to send message", zap.Error(err))
				return
			}
			return
		}

		// POST to mawinter-server
		res, err := clientrepo.PostMawinter(&user.ServerInfo, user.ChooseCategoryID, int64(price))
		if err != nil {
			// TODO: invalid data or connection lost かで分ける
			_, nerr := s.ChannelMessageSend(m.ChannelID, "internal error")
			if nerr != nil {
				logger.Error("failed to send message", zap.Error(nerr))
				return
			}
			logger.Error("failed to send order to mawinter-server", zap.Error(err))
			user.changeCtxStatus(ContextOrderWaiting)
			categoryChoicePost(s, m, user)
			return
		}

		logger.Info("receive response from mawinter", zap.String("addr", user.ServerInfo.Addr))
		resText := fmt.Sprintf("ID: %v, Name: %v, Price: %v", res.Id, res.Name, res.Price)
		user.LastOrderID = res.Id
		_, err = s.ChannelMessageSend(m.ChannelID, resText)
		if err != nil {
			logger.Error("failed to send message", zap.Error(err))
			return
		}

		user.changeCtxStatus(ContextOrderWaiting)
		categoryChoicePost(s, m, user)
	}

}

func messageReaction(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	logger.Debug("handler called", zap.String("name", "MessageReactionAdd"), zap.String("userID", r.UserID))
	loggerID := logger.With(zap.String("name", "MessageCreate"), zap.String("userID", r.UserID))
	if r.UserID == s.State.User.ID {
		return
	}

	user, err := users.getUserInfoByID(r.UserID)
	if errors.Is(err, ErrUserNotFound) {
		loggerID.Warn("unknown user")
	} else if err != nil {
		loggerID.Error("unknown error", zap.Error(err))
	}

	loggerIDName := logger.With(zap.String("username", user.Name))

	switch r.Emoji.Name {
	case categoryFood:
		loggerIDName.Info("category choose", zap.String("category", "food"))
		user.ChooseCategoryID = categoryIDFood
	case categoryLife:
		loggerIDName.Info("category choose", zap.String("category", "life"))
		user.ChooseCategoryID = categoryIDLife
	case categoryEntertainment:
		loggerIDName.Info("category choose", zap.String("category", "entertainment"))
		user.ChooseCategoryID = categoryIDEntertainment
	case categoryFriends:
		loggerIDName.Info("category choose", zap.String("category", "friends"))
		user.ChooseCategoryID = categoryIDFriends
	case categoryCompute:
		loggerIDName.Info("category choose", zap.String("category", "compute"))
		user.ChooseCategoryID = categoryIDCompute
	case categoryTrans:
		loggerIDName.Info("category choose", zap.String("category", "trans"))
		user.ChooseCategoryID = categoryIDTrans
	case categoryStudy:
		loggerIDName.Info("category choose", zap.String("category", "study"))
		user.ChooseCategoryID = categoryIDStudy
	case operateCancel:
		// 直前の登録をキャンセルする
		loggerIDName.Info("category choose", zap.String("operate", "cancel"))

		// lastID = -1 --> not recorded
		if user.LastOrderID == -1 {
			_, err = s.ChannelMessageSend(r.ChannelID, "not recorded last order")
			if err != nil {
				logger.Error("failed to send message", zap.Error(err))
			}
			return
		}

		err = clientrepo.DeleteMawinter(&user.ServerInfo, user.LastOrderID)
		if err != nil {
			logger.Error("failed to delete the last record", zap.Error(err))
			s.ChannelMessageSend(r.ChannelID, "failed to delete the last record")
		}

		s.ChannelMessageSend(r.ChannelID, "deleted the last record")
		if err != nil {
			logger.Error("failed to send message", zap.Error(err))
		}

		user.LastOrderID = -1
		user.changeCtxStatus(ContextOrderWaiting)
		return
	default:
		loggerIDName.Info("unknown choose", zap.String("category", r.Emoji.Name))
		user.ChooseCategoryID = -1
		return
	}

	user.changeCtxStatus(ContextPriceWaiting)
	_, err = s.ChannelMessageSend(r.ChannelID, "you choose "+r.Emoji.Name+"\n"+"what value? (type 'r' to back to category select.)")
	if err != nil {
		logger.Error("failed to send message", zap.Error(err))
		return
	}
}

func categoryChoicePost(s *discordgo.Session, m *discordgo.MessageCreate, user *discordUser) {
	mes, err := s.ChannelMessageSend(m.ChannelID, "choose category")
	if err != nil {
		logger.Error("failed to send message", zap.Error(err))
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
	s.MessageReactionAdd(chanID, mesID, operateCancel)
}
