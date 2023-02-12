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
	udpServer     *UDPServer
	udpstorage    *store.Storage
	udplogService *logger.Logger
)

func init() {
	udplogService = logger.NewLogger()
	udplogService.StartNoopLogger()
	udpstorage = store.NewStorage(udplogService)
	udpmetrics := metrics.NewMetrics(udplogService)
	udpmetrics.StartNoopMetrics()
	udpServer = NewUDP(udplogService, udpstorage, udpmetrics)

	udpServer.Start()
}

func TestUDPHandlers_UDPStart(t *testing.T) {
	type args struct {
		network string
		rAddr   *net.UDPAddr
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Dial UDP Server - ok",
			args: args{
				network: "udp",
				rAddr: &net.UDPAddr{
					IP:   []byte{0, 0, 0, 0},
					Port: 9001,
					Zone: "",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			conn, err := net.DialUDP(tt.args.network, nil, tt.args.rAddr)
			if err != nil {
				t.Errorf("failed to dial udp server error: %v", err)
			}

			defer func() {
				conn.Close()
			}()
		})
	}
}

func TestUDPHandlers_UDPHandler(t *testing.T) {
	type args struct {
		network      string
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
				network:      "udp",
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
				network: "udp",
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
				network:      "udp",
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
				network:      "udp",
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
				network:      "udp",
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
				network: "udp",
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
				network: "udp",
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
				network: "udp",
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

	lAddr := &net.UDPAddr{
		IP:   []byte{0, 0, 0, 0},
		Port: 9002,
		Zone: "",
	}
	rAddr := &net.UDPAddr{
		IP:   []byte{0, 0, 0, 0},
		Port: 9001,
		Zone: "",
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			*udpstorage = *store.NewStorage(udplogService)

			if tt.args.addStoreItem {
				udpstorage.Post(map[string]interface{}{"1": "hello world"})
			}

			conn, err := net.DialUDP(tt.args.network, lAddr, rAddr)
			if err != nil {
				t.Errorf("dial udp error: %v", err)
			}

			defer conn.Close()

			req, err := json.Marshal(tt.args.data)
			if err != nil {
				t.Errorf("marshal data error: %v", err)
			}

			if _, err := conn.Write(req); err != nil {
				t.Errorf("write to udp error: %v", err)
			}

			buf := make([]byte, 1024)
			n, err := conn.Read(buf)
			if err != nil {
				t.Error("response did match expected output")
			}

			var got jsonResponse
			json.Unmarshal(buf[0:n], &got)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UDPHandler got = %v, want %v", got, tt.want)
			}
		})
	}
}
