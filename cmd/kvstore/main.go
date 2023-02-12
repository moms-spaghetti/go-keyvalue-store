package main

import (
	"os"
	"os/signal"
	"syscall"
	"task1/internal/logger"
	"task1/internal/metrics"
	"task1/internal/protocols"
	"task1/internal/store"
)

func main() {
	logger := logger.NewLogger()
	metrics := metrics.NewMetrics(logger)
	storage := store.NewStorage(logger)

	udp := *protocols.NewUDP(logger, storage, metrics)
	http := *protocols.NewHTTP(logger, storage, metrics)
	tcp := *protocols.NewTCP(logger, storage, metrics)

	starts := []func(){
		logger.Start,
		udp.Start,
		http.Start,
		tcp.Start,
		metrics.Start,
	}

	stops := []func(){
		udp.Stop,
		http.Stop,
		tcp.Stop,
		metrics.Stop,
		logger.Stop,
	}

	wait := make(chan os.Signal, 1)
	signal.Notify(wait, syscall.SIGINT)

	run(starts)

	<-wait
	defer close(wait)

	run(stops)
}

func run(fn []func()) {
	for _, f := range fn {
		f()
	}
}
