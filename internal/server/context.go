package server

import "errors"

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
	ID          string
	Name        string // for memo
	Context     ContextStatus
	LastOrderID int64 // not found = -1
}

func (d *discordUsers) addDiscordUser(id string, name string) *discordUser {
	newUser := discordUser{ID: id, Name: name, Context: ContextClosing, LastOrderID: -1}
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

func (d *discordUser) changeCtxStatus(nextCtx ContextStatus) {
	logger.Infow("change status", "userID", d.ID, "username", d.Name, "nowstatus", d.Context, "nextstatus", nextCtx)
	d.Context = nextCtx
}
