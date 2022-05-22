//go:build integration
// +build integration

package client

import (
	"os"
	"testing"

	"go.uber.org/zap"
)

var testServerInfo ServerInfo
var testServerInvalidInfo ServerInfo

func TestMain(m *testing.M) {
	slogger, _ := zap.NewDevelopment()
	defer slogger.Sync()
	Logger = slogger.Sugar()

	testServerInfo = ServerInfo{Addr: "http://localhost:8080/", User: "test", Pass: "test"}
	testServerInvalidInfo = ServerInfo{Addr: "http://localhost:8080/", User: "test", Pass: "wrong"}
	code := m.Run()

	os.Exit(code)
}

func Test_clientRepo_PostMawinter(t *testing.T) {
	type args struct {
		info       *ServerInfo
		categoryID int64
		price      int64
	}
	tests := []struct {
		name    string
		c       *clientRepo
		args    args
		want    *RecordsDetails
		wantErr bool
	}{
		{
			name:    "normal",
			c:       &clientRepo{},
			args:    args{info: &testServerInfo, categoryID: 200, price: 1234},
			want:    &RecordsDetails{CategoryId: 200, Price: 1234},
			wantErr: false,
		},
		{
			name:    "invalid categoryID",
			c:       &clientRepo{},
			args:    args{info: &testServerInfo, categoryID: -1, price: 1234},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "invalid access",
			c:       &clientRepo{},
			args:    args{info: &testServerInvalidInfo, categoryID: 200, price: 1234},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.c.PostMawinter(tt.args.info, tt.args.categoryID, tt.args.price)
			if (err != nil) != tt.wantErr {
				t.Errorf("clientRepo.PostMawinter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// if !reflect.DeepEqual(got, tt.want) {
			// 	t.Errorf("clientRepo.PostMawinter() = %v, want %v", got, tt.want)
			// }
			if !tt.wantErr {
				if got.CategoryId != tt.want.CategoryId {
					t.Errorf("clientRepo.PostMawinter() = %v, want %v", got.CategoryId, tt.want.CategoryId)
				}
				if got.Price != tt.want.Price {
					t.Errorf("clientRepo.PostMawinter() = %v, want %v", got.Price, tt.want.Price)
				}
			}
		})
	}
}
