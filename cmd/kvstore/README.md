# workbench
### find lines
wrkbench -n "POST-adds 3 items" -e "http://localhost:8080/" -v false -H "Content-Type: application/json" -m "POST" -d '{"Payload":{"1":"more random text","2":123,"3":false}}'
wrkbench -n "GET-get single key" -e "http://localhost:8080/" -v false -H "Content-Type: application/json" -m "GET" -d '{"Query":"1"}'