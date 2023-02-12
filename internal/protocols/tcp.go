package protocols

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"task1/internal/logger"
	"task1/internal/metrics"
	"task1/internal/store"
)

const (
	tcpaddr    = ":8181"
	tcpnetwork = "tcp"
)

type TCPServer struct {
	listener net.Listener
	conns    map[string]net.Conn
	done     chan struct{}
	logger   *logger.Logger
	storage  *store.Storage
	metrics  *metrics.Metrics
}

func NewTCP(
	logger *logger.Logger,
	storage *store.Storage,
	metrics *metrics.Metrics,
) *TCPServer {

	lis, err := net.Listen(tcpnetwork, tcpaddr)
	if err != nil {
		panic(err)
	}

	return &TCPServer{
		listener: lis,
		conns:    make(map[string]net.Conn),
		done:     make(chan struct{}),
		logger:   logger,
		storage:  storage,
		metrics:  metrics,
	}
}

func (ts TCPServer) Start() {
	go func() {
		for {
			log.Printf("tcp listning on %s", tcpaddr)
			conn, err := ts.listener.Accept()
			if err != nil {
				select {
				case <-ts.done:
					if len(ts.conns) > 0 {
						ts.logger.Log("closing active conns")

						for n, c := range ts.conns {
							ts.logger.Log(fmt.Sprintf("active conn %s closed", n))
							c.Close()
						}
					}
					return
				default:
					log.Printf("listener error: %v", err)

					return
				}
			}

			connID := createConnID()
			ts.addConn(conn, connID)

			go ts.tcpHandler(conn, connID)
		}
	}()
}

func (ts TCPServer) Stop() {
	close(ts.done)
	if err := ts.listener.Close(); err != nil {
		log.Printf("listener close err: %v", err)
	}
	log.Print("TCP shutdown ok")
}

func (ts TCPServer) addConn(conn net.Conn, connID string) {
	ts.logger.Log(fmt.Sprintf("conn %s added", connID))
	ts.conns[connID] = conn
}

func (ts TCPServer) removeConn(connID string) {
	ts.logger.Log(fmt.Sprintf("conn %s removed", connID))
	delete(ts.conns, connID)
}

func createConnID() string {
	// horrible dirty ID creation
	id := make([]byte, 4)
	rand.Read(id)

	return fmt.Sprintf("%X", id[0:4])
}

func (ts TCPServer) tcpHandler(conn net.Conn, connID string) {
	var (
		err       error
		n         int
		req       jsonRequest
		storeData interface{}
	)

	buf := make([]byte, 1024)
	if err == nil {
		n, err = conn.Read(buf)
	}

	if err == nil {
		err = json.Unmarshal(buf[:n], &req)
	}

	if err == nil {
		ts.metrics.LogMetrics(req.Method)
		switch req.Method {
		case http.MethodGet:
			ts.logger.Log("TCP GET request")
			storeData, err = ts.storage.Get(req.Query)
		case http.MethodPost:
			ts.logger.Log("TCP POST request")
			err = ts.storage.Post(req.Payload)
		case http.MethodDelete:
			ts.logger.Log("TCP DELETE request")
			err = ts.storage.Delete(req.Query)
		default:
			err = ErrRouteForbidden
		}
	}

	_, response := BuildJsonResponse(err, storeData, ts.logger)
	conn.Write(response)
	conn.Close()
	ts.removeConn(connID)
}
