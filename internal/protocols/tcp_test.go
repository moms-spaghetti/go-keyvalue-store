package protocols

import (
	"encoding/json"
	"net"
	"reflect"
	"task1/internal/logger"
	"task1/internal/metrics"
	"task1/internal/store"
	"testing"
)

var (
	tcpServer     *TCPServer
	tcpstorage    *store.Storage
	tcplogService *logger.Logger
)

func init() {
	tcplogService = logger.NewLogger()
	tcplogService.StartNoopLogger()
	tcpstorage = store.NewStorage(tcplogService)
	tcpmetrics := metrics.NewMetrics(tcplogService)
	tcpmetrics.StartNoopMetrics()
	tcpServer = NewTCP(tcplogService, tcpstorage, tcpmetrics)

	tcpServer.Start()
}

func TestTCPServer_TCPStart(t *testing.T) {
	type args struct {
		network string
		addr    string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Dial TCP Server - ok",
			args: args{
				network: "tcp",
				addr:    ":8181",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			conn, err := net.Dial(tt.args.network, tt.args.addr)
			if err != nil {
				t.Errorf("failed to dial tcp server error: %v", err)
			}

			defer func() {
				conn.Close()
			}()
		})
	}
}

func TestTCPServer_tcpHandler(t *testing.T) {
	type args struct {
		addStoreItem bool
		data         map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    jsonResponse
		wantErr bool
	}{
		{
			name: "GET ok",
			args: args{
				addStoreItem: true,
				data: map[string]interface{}{
					"Query":  "1",
					"Method": "GET",
				},
			},
			want: jsonResponse{
				Err:    "",
				Status: 200,
				Data:   "hello world",
			},
		},
		{
			name: "GET fail - store is empty",
			args: args{
				data: map[string]interface{}{
					"Query":  "1",
					"Method": "GET",
				},
			},
			want: jsonResponse{
				Err:    "store is empty",
				Status: 500,
				Data:   nil,
			},
		},
		{
			name: "GET fail - key not found in store",
			args: args{
				addStoreItem: true,
				data: map[string]interface{}{
					"Query":  "2",
					"Method": "GET",
				},
			},
			want: jsonResponse{
				Err:    "key not found in store",
				Status: 404,
				Data:   nil,
			},
		},
		{
			name: "GET fail - key cannot be empty",
			args: args{
				addStoreItem: true,
				data: map[string]interface{}{
					"Query":  "",
					"Method": "GET",
				},
			},
			want: jsonResponse{
				Err:    "key cannot be empty",
				Status: 400,
				Data:   nil,
			},
		},
		{
			name: "GET fail - unsupported method",
			args: args{
				addStoreItem: true,
				data: map[string]interface{}{
					"Query":  "1",
					"Method": "PATCH",
				},
			},
			want: jsonResponse{
				Err:    "method forbidden",
				Status: 405,
				Data:   nil,
			},
		},
		{
			name: "POST ok - single item",
			args: args{
				data: map[string]interface{}{
					"Method":  "POST",
					"Payload": map[string]interface{}{"1": "hello"},
				},
			},
			want: jsonResponse{
				Err:    "",
				Status: 200,
				Data:   nil,
			},
		},
		{
			name: "POST ok - multiple items",
			args: args{
				data: map[string]interface{}{
					"Method":  "POST",
					"Payload": map[string]interface{}{"1": "hello", "2": "world"},
				},
			},
			want: jsonResponse{
				Err:    "",
				Status: 200,
				Data:   nil,
			},
		},
		{
			name: "POST fail - key cannot be empty",
			args: args{
				data: map[string]interface{}{
					"Method":  "POST",
					"Payload": map[string]interface{}{"": "missing key"},
				},
			},
			want: jsonResponse{
				Err:    "key cannot be empty",
				Status: 400,
				Data:   nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			*tcpstorage = *store.NewStorage(tcplogService)

			if tt.args.addStoreItem {
				tcpstorage.Post(map[string]interface{}{"1": "hello world"})
			}

			conn, err := net.Dial("tcp", ":8181")
			if err != nil {
				panic(err)
			}

			defer conn.Close()

			req, err := json.Marshal(tt.args.data)
			if err != nil {
				t.Errorf("marshal data error: %v", err)
			}

			if _, err := conn.Write(req); err != nil {
				t.Errorf("write to tcp error: %v", err)
			}

			buf := make([]byte, 1024)
			n, err := conn.Read(buf)
			if err != nil {
				t.Error("response did match expected output")
			}

			var got jsonResponse
			json.Unmarshal(buf[0:n], &got)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TCPHandler got = %v, want %v", got, tt.want)
			}
		})
	}
}
