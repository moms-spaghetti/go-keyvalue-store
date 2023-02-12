package protocols

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	_ "net/http/pprof"
	"strings"
	"task1/internal/logger"
	"task1/internal/metrics"
	"task1/internal/store"
	"testing"
)

const url = "http://localhost:9001"

func TestHTTPHandlers_rootHandler(t *testing.T) {
	type args struct {
		w            *httptest.ResponseRecorder
		r            *http.Request
		addStoreItem bool
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "GET ok",
			args: args{
				r:            httptest.NewRequest(http.MethodGet, url, strings.NewReader(`{"Query":"1"}`)),
				w:            httptest.NewRecorder(),
				addStoreItem: true,
			},
			want: `{"Err":"","Status":200,"Data":"hello world"}`,
		},
		{
			name: "GET fail - store is empty",
			args: args{
				r: httptest.NewRequest(http.MethodGet, url, strings.NewReader(`{"Query":"1"}`)),
				w: httptest.NewRecorder(),
			},
			want: `{"Err":"store is empty","Status":500,"Data":null}`,
		},
		{
			name: "GET fail - key not found in store",
			args: args{
				r:            httptest.NewRequest(http.MethodGet, url, strings.NewReader(`{"Query":"2"}`)),
				w:            httptest.NewRecorder(),
				addStoreItem: true,
			},
			want: `{"Err":"key not found in store","Status":404,"Data":null}`,
		},
		{
			name: "GET fail - key cannot be empty",
			args: args{
				r:            httptest.NewRequest(http.MethodGet, url, strings.NewReader(`{"Query":""}`)),
				w:            httptest.NewRecorder(),
				addStoreItem: true,
			},
			want: `{"Err":"key cannot be empty","Status":400,"Data":null}`,
		},
		{
			name: "GET fail - unexpected EOF",
			args: args{
				r:            httptest.NewRequest(http.MethodGet, url, strings.NewReader(`{"Query":""`)),
				w:            httptest.NewRecorder(),
				addStoreItem: true,
			},
			want: `{"Err":"unexpected EOF","Status":500,"Data":null}`,
		},
		{
			name: "METHOD fail - unsupported method",
			args: args{
				r:            httptest.NewRequest(http.MethodPatch, url, strings.NewReader(`{"Query":"1"}`)),
				w:            httptest.NewRecorder(),
				addStoreItem: true,
			},
			want: `{"Err":"method forbidden","Status":405,"Data":null}`,
		},
		{
			name: "POST ok - single item",
			args: args{
				r: httptest.NewRequest(http.MethodPost, url, strings.NewReader(`{"Payload":{"1":"hello world"}}`)),
				w: httptest.NewRecorder(),
			},
			want: `{"Err":"","Status":200,"Data":null}`,
		},
		{
			name: "POST ok - multiple items",
			args: args{
				r: httptest.NewRequest(http.MethodPost, url, strings.NewReader(`{"Payload":{"1":"hello", "2":"world"}}`)),
				w: httptest.NewRecorder(),
			},
			want: `{"Err":"","Status":200,"Data":null}`,
		},
		{
			name: "POST fail - unexpected EOF",
			args: args{
				r: httptest.NewRequest(http.MethodPost, url, strings.NewReader(`{"Payload":{"1":""}`)),
				w: httptest.NewRecorder(),
			},
			want: `{"Err":"unexpected EOF","Status":500,"Data":null}`,
		},
		{
			name: "POST fail - key cannot be empty",
			args: args{
				r: httptest.NewRequest(http.MethodPost, url, strings.NewReader(`{"Payload":{"":"missing key"}}`)),
				w: httptest.NewRecorder(),
			},
			want: `{"Err":"key cannot be empty","Status":400,"Data":null}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := logger.NewLogger()
			logger.StartNoopLogger()
			storage := store.NewStorage(logger)
			metrics := metrics.NewMetrics(logger)
			metrics.Start()
			rh := NewHTTP(logger, storage, metrics)

			if tt.args.addStoreItem {
				rh.storage.Post(map[string]interface{}{"1": "hello world"})
			}

			rh.rootHandler(tt.args.w, tt.args.r)
			res := tt.args.w.Result()
			defer res.Body.Close()

			data, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Errorf("rootHandler error = %v", err)
			}

			got := string(data)
			if got != tt.want {
				t.Errorf("rootHandler: got = %v, want = %v", got, tt.want)
			}
		})
	}
}
