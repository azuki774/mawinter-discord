package server

import (
	"reflect"
	"testing"

	"github.com/azuki774/mawinter-discord/internal/client"
	"go.uber.org/zap"
)

func TestInit(t *testing.T) {
	slogger, _ := zap.NewDevelopment()
	defer slogger.Sync()
	logger = slogger.Sugar()
}

func Test_discordUsers_addDiscordUser(t *testing.T) {
	type args struct {
		sinfo client.ServerInfo
		id    string
		name  string
	}
	tests := []struct {
		name string
		d    *discordUsers
		args args
		want *discordUser
	}{
		{
			name: "new first user",
			d:    &discordUsers{},
			args: args{sinfo: client.ServerInfo{}, id: "1", name: "test1"},
			want: &discordUser{ServerInfo: client.ServerInfo{}, ID: "1", Name: "test1", Context: ContextClosing, LastOrderID: -1},
		},
		{
			name: "new second user",
			d:    &discordUsers{Users: []*discordUser{{ID: "1", Name: "test1", Context: ContextClosing, LastOrderID: 123}}},
			args: args{sinfo: client.ServerInfo{}, id: "2", name: "test2"},
			want: &discordUser{ServerInfo: client.ServerInfo{}, ID: "2", Name: "test2", Context: ContextClosing, LastOrderID: -1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.addDiscordUser(tt.args.sinfo, tt.args.id, tt.args.name); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("discordUsers.addDiscordUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_discordUsers_getUserInfoByID(t *testing.T) {
	users10 := discordUsers{Users: []*discordUser{
		{ID: "1", Name: "test1", Context: ContextClosing, LastOrderID: 123},
		{ID: "2", Name: "test2", Context: ContextClosing, LastOrderID: 456},
		{ID: "3", Name: "test3", Context: ContextClosing, LastOrderID: -1},
	}}

	type args struct {
		targetID string
	}
	tests := []struct {
		name    string
		d       *discordUsers
		args    args
		want    *discordUser
		wantErr bool
	}{
		{
			name:    "Exists user",
			d:       &users10,
			args:    args{targetID: "1"},
			want:    &discordUser{ID: "1", Name: "test1", Context: ContextClosing, LastOrderID: 123},
			wantErr: false,
		},
		{
			name:    "not Exists user",
			d:       &users10,
			args:    args{targetID: "4"},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.d.getUserInfoByID(tt.args.targetID)
			if (err != nil) != tt.wantErr {
				t.Errorf("discordUsers.getUserInfoByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("discordUsers.getUserInfoByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_discordUser_changeCtxStatus(t *testing.T) {
	type args struct {
		nextCtx ContextStatus
	}
	tests := []struct {
		name string
		d    *discordUser
		args args
		want *discordUser
	}{
		{
			name: "Closing -> OrderWaiting",
			d:    &discordUser{ID: "1", Name: "test1", Context: ContextClosing, LastOrderID: -1},
			args: args{nextCtx: ContextOrderWaiting},
			want: &discordUser{ID: "1", Name: "test1", Context: ContextOrderWaiting, LastOrderID: -1},
		},
		{
			name: "OrderWaiting -> PriceWaiting",
			d:    &discordUser{ID: "1", Name: "test1", Context: ContextOrderWaiting, LastOrderID: -1},
			args: args{nextCtx: ContextPriceWaiting},
			want: &discordUser{ID: "1", Name: "test1", Context: ContextPriceWaiting, LastOrderID: -1},
		},
		{
			name: "PriceWaiting -> OrderWaiting",
			d:    &discordUser{ID: "1", Name: "test1", Context: ContextPriceWaiting, LastOrderID: -1},
			args: args{nextCtx: ContextOrderWaiting},
			want: &discordUser{ID: "1", Name: "test1", Context: ContextOrderWaiting, LastOrderID: -1},
		},
		{
			name: "OrderWaiting -> OrderWaiting",
			d:    &discordUser{ID: "1", Name: "test1", Context: ContextOrderWaiting, LastOrderID: -1},
			args: args{nextCtx: ContextOrderWaiting},
			want: &discordUser{ID: "1", Name: "test1", Context: ContextOrderWaiting, LastOrderID: -1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.changeCtxStatus(tt.args.nextCtx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("discordUser.changeCtxStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}
