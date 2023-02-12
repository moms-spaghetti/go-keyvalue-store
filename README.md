# netcat udp helper lines
### GET
{"Method":"GET", "Query":"1"}  
{"Method":"GET", "Query":"2"}  
{"Method":"GET", "Query":"3"}  

### POST
{"Method":"POST", "Payload":{"1":"more random text","2":123,"3":false}}  
{"Method":"POST", "Payload":{"3":true}}  

### DELETE
{"Method":"DELETE", "Query":"1"}  

### Errors
{"Method":"GET", "":"missing key"}  
{"Method":"BADMETHOD", "Query":"3"}  
{"Method":"GET", "Query":"none existent query"}  
{"Method":"GET", "Query":"unexpected end of json input...  
{"Method":"POST", "Payload":{"":"missing key"}}  

# URLs
### http
localhost:8080/  

### udp
localhost:9001/

# JSON
### GET
`{
	"Query":"1"
}`

### POST  
`{
    "Payload": {
        "1": "random text",
        "2": 123,
        "3": true,
        "4": {
            "hello": "world"
        }
    }
}`

# makefile
### commands
make run
make build
make uget (reads query)
make upost (adds three items)

# workbench
### find lines
use `wrkbench find` in cmd/kvstore/README.md

### notes
bench 'GET-get single key' will need item adding to store such as `store := map[string]interface{}{"1": "hello world"}` in `store.go` or will return error json store is empty

# tcp client
`cmd/tcpclient/main.go` simple client with console interface to make requests against tcp protocol  
`makerun runc` to start