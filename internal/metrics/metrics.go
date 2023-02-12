package metrics

import (
	"log"
	"task1/internal/logger"
	"time"

	"fmt"
)

const (
	statGet    = "GET"
	statPost   = "POST"
	statDelete = "DELETE"
	printDelay = 10 * time.Second
)

type Metrics struct {
	metrics chan string
	done    chan struct{}
	stats   *stats
	logger  *logger.Logger
}

type stats struct {
	get     int
	post    int
	delete  int
	unknown int
}

func NewMetrics(logger *logger.Logger) *Metrics {
	metrics := &Metrics{
		metrics: make(chan string),
		done:    make(chan struct{}),
		stats: &stats{
			get:     0,
			post:    0,
			delete:  0,
			unknown: 0,
		},
		logger: logger,
	}

	return metrics
}

func (m *Metrics) Start() {
	log.Print("metrics started")
	go func() {
		for {
			select {
			case stat := <-m.metrics:
				switch stat {
				case statGet:
					m.stats.get++
				case statPost:
					m.stats.post++
				case statDelete:
					m.stats.delete++
				default:
					m.stats.unknown++
				}

				m.PrintMetrics()
			case <-m.done:
				return
			}
		}
	}()
}

func (m *Metrics) StartNoopMetrics() {
	go func() {
		for {
			select {
			case <-m.metrics:
			case <-m.done:
				return
			}
		}
	}()
}

func (m *Metrics) LogMetrics(method string) {
	m.metrics <- method
}

func (m *Metrics) PrintMetrics() {
	out := fmt.Sprintf("\nMETRICS - GET: %d, POST: %d, DELETE: %d, UNKNOWN: %d ",
		m.stats.get,
		m.stats.post,
		m.stats.delete,
		m.stats.unknown,
	)
	m.logger.Log(out)
}

func (m *Metrics) Stop() {
	close(m.metrics)
	close(m.done)
	log.Print("metrics shutdown ok")
}
