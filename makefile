.PHONY: run build uget upost runc buildc
run:
		go run cmd/kvstore/main.go

build:
		go build cmd/kvstore/main.go

uget:
		@read -p "query: " q; \
		(echo '{"Method":"GET", "Query":"'$$q'"}'; sleep 0.75) | nc -u 0.0.0.0 9001

upost:
		(echo '{"Method":"POST", "Payload":{"1":"more random text","2":123,"3":false}}'; sleep 0.75) | nc -u 0.0.0.0 9001

runc:
		go run cmd/tcpclient/main.go

buildc:
		go build cmd/tcpclient/main.go