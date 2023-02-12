package protocols

import (
	"encoding/json"
	"errors"
	"reflect"
	"task1/internal/logger"
	"task1/internal/store"
	"testing"
)

func Test_protocolHelpers_buildJsonResponse(t *testing.T) {
	type args struct {
		err  error
		data interface{}
	}
	tests := []struct {
		name  string
		args  args
		want  int
		want1 jsonResponse
	}{
		{
			name: "ok no error",
			args: args{
				err:  nil,
				data: "hello world",
			},
			want: 200,
			want1: jsonResponse{
				Err:    "",
				Status: 200,
				Data:   "hello world",
			},
		},
		{
			name: "err key not found in store",
			args: args{
				err:  store.ErrStoreKeyNotFound,
				data: "",
			},
			want: 404,
			want1: jsonResponse{
				Err:    "key not found in store",
				Status: 404,
				Data:   nil,
			},
		},
		{
			name: "err key cannot be empty",
			args: args{
				err:  store.ErrKeyEmpty,
				data: "",
			},
			want: 400,
			want1: jsonResponse{
				Err:    "key cannot be empty",
				Status: 400,
				Data:   nil,
			},
		},
		{
			name: "err method forbidden",
			args: args{
				err:  ErrRouteForbidden,
				data: "",
			},
			want: 405,
			want1: jsonResponse{
				Err:    "method forbidden",
				Status: 405,
				Data:   nil,
			},
		},
		{
			name: "err internal server error",
			args: args{
				err:  errors.New("some other error"),
				data: "",
			},
			want: 500,
			want1: jsonResponse{
				Err:    "some other error",
				Status: 500,
				Data:   nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := logger.NewLogger()
			logger.StartNoopLogger()

			got, got1 := BuildJsonResponse(tt.args.err, tt.args.data, logger)
			if got != tt.want {
				t.Errorf("protocolHelpers.buildJsonResponse() got = %v, want %v", got, tt.want)
			}

			// convert back to jsonResponse else checking bytes vs bytes
			var assertResponse jsonResponse

			if err := json.Unmarshal(got1, &assertResponse); err != nil {
				t.Errorf("json.Unmarshal assertResponse err = %v", err)

				return
			}

			if !reflect.DeepEqual(assertResponse, tt.want1) {
				t.Errorf("protocolHelpers.buildJsonResponse() assertResponse = %v, want %v", assertResponse, tt.want1)
			}
		})
	}
}
