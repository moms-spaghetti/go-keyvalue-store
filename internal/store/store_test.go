package store

import (
	"reflect"
	"task1/internal/logger"
	"testing"
)

func TestService_Get(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name         string
		args         args
		addStoreItem bool
		want         interface{}
		wantErr      bool
	}{
		{
			name: "GET - ok",
			args: args{
				key: "1",
			},
			addStoreItem: true,
			want:         "hello world",
		},
		{
			name: "GET fail - key empty",
			args: args{
				key: "",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "GET fail - store empty",
			args: args{
				key: "1",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "GET fail - key not found",
			args: args{
				key: "2",
			},
			addStoreItem: true,
			want:         nil,
			wantErr:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := logger.NewLogger()
			logger.Start()
			kv := NewStorage(logger)

			if tt.addStoreItem {
				kv.Post(map[string]interface{}{"1": "hello world"})
			}

			got, err := kv.Get(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Service.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_Post(t *testing.T) {
	type args struct {
		data StoreData
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "POST - ok",
			args: args{
				data: map[string]interface{}{"1": "hello world"},
			},
		},
		{
			name: "POST fail - key empty",
			args: args{
				data: map[string]interface{}{"": "hello world"},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := logger.NewLogger()
			logger.Start()
			kv := NewStorage(logger)

			if err := kv.Post(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("Service.Post() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
