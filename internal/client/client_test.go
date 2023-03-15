package client

import (
	"reflect"
	"testing"

	"go.uber.org/zap"
)

func TestInit(t *testing.T) {
	Logger, _ = zap.NewProduction()
}

func TestNewClientRepo(t *testing.T) {
	tests := []struct {
		name string
		want *clientRepo
	}{
		{
			name: "create",
			want: &clientRepo{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewClientRepo(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewClientRepo() = %v, want %v", got, tt.want)
			}
		})
	}
}
