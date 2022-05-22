package server

import (
	"errors"

	"github.com/azuki774/mawinter-discord/internal/client"
)

type ContextStatus string

var (
	// ContextCloseing ->
	// ContextOrderWaiting -> <stamp> -> ContextPriceWaiting -> <price input>
	// -> ContextCategoryWaiting ..
	ContextClosing      = ContextStatus("Closing")
	ContextOrderWaiting = ContextStatus("OrderWaiting")
	ContextPriceWaiting = ContextStatus("PriceWaiting")
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type discordUsers struct {
	Users []*discordUser
}
type discordUser struct {
	ServerInfo       client.ServerInfo
	ID               string
	Name             string // for memo
	Context          ContextStatus
	ChooseCategoryID int64 // now choosing categoryID
	LastOrderID      int64 // not found = -1
}

func (d *discordUsers) addDiscordUser(sinfo client.ServerInfo, id string, name string) *discordUser {
	// ToDo: ServerInfo をセットする
	newUser := discordUser{ServerInfo: sinfo, ID: id, Name: name, Context: ContextClosing, LastOrderID: -1}
	d.Users = append(d.Users, &newUser)
	return &newUser
}

func (d *discordUsers) getUserInfoByID(targetID string) (*discordUser, error) {
	for _, v := range d.Users {
		if v.ID == targetID {
			return v, nil
		}
	}

	return nil, ErrUserNotFound
}

func (d *discordUser) changeCtxStatus(nextCtx ContextStatus) *discordUser {
	logger.Infow("change status", "userID", d.ID, "username", d.Name, "nowstatus", d.Context, "nextstatus", nextCtx)
	d.Context = nextCtx
	return d
}
