package protocols

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"task1/internal/logger"
	"task1/internal/metrics"
	"task1/internal/store"
)

const (
	udpnetwork = "udp"
	udpport    = 9001
	zone       = ""
	bufsize    = 1024
)

type UDPServer struct {
	conn      *net.UDPConn
	closeConn chan struct{}
	logger    *logger.Logger
	storage   *store.Storage
	metrics   *metrics.Metrics
}

func NewUDP(
	logger *logger.Logger,
	storage *store.Storage,
	metrics *metrics.Metrics,
) *UDPServer {
	addr := &net.UDPAddr{
		IP:   []byte{0, 0, 0, 0},
		Port: udpport,
		Zone: zone,
	}

	conn, err := net.ListenUDP(udpnetwork, addr)
	if err != nil {
		panic(err)
	}

	return &UDPServer{
		conn:      conn,
		closeConn: make(chan struct{}),
		logger:    logger,
		storage:   storage,
		metrics:   metrics,
	}
}

func (us UDPServer) Start() {
	go func() {
		log.Printf("udp listening on %s", us.conn.LocalAddr().String())
		buf := make([]byte, bufsize)

		for {
			select {
			case <-us.closeConn:
				us.conn.Close()

				return
			default:
				n, retAddr, err := us.conn.ReadFromUDP(buf)
				if err != nil {
					log.Printf("startUDP error: %v", err)
				}

				go us.UDPHandler(buf, n, retAddr)
			}
		}
	}()
}

func (us UDPServer) Stop() {
	close(us.closeConn)
	log.Print("UDP shutdown ok")
}

func (us UDPServer) UDPHandler(buf []byte, n int, retAddr *net.UDPAddr) {
	var (
		req       jsonRequest
		err       error
		storeData interface{}
	)

	if err == nil {
		err = json.Unmarshal(buf[0:n], &req)
	}

	if err == nil {
		us.metrics.LogMetrics(req.Method)
		switch req.Method {
		case http.MethodGet:
			us.logger.Log("UDP GET request")
			storeData, err = us.storage.Get(req.Query)
		case http.MethodPost:
			us.logger.Log("UDP POST request")
			err = us.storage.Post(req.Payload)
		case http.MethodDelete:
			us.logger.Log("UDP DELETE request")
			err = us.storage.Delete(req.Query)
		default:
			err = ErrRouteForbidden
		}
	}

	_, out := BuildJsonResponse(err, storeData, us.logger)
	us.conn.WriteTo(out, retAddr)
}
