package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

type jsonRequest struct {
	Method  string                 `json:"Method"`
	Query   string                 `json:"Query"`
	Payload map[string]interface{} `json:"Payload"`
}

type jsonResponse struct {
	Err    string      `json:"Err"`
	Status int         `json:"Status"`
	Data   interface{} `json:"Data"`
}

func main() {
	// server connection check
	svrConnLimit := 0
	for {
		conn, err := net.Dial("tcp", ":8181")
		if err != nil {
			fmt.Println("no sever connection...")
			fmt.Printf("trying again in 3 seconds...\n\n")
			time.Sleep(time.Second * 3)

			svrConnLimit++
			if svrConnLimit == 3 {
				fmt.Println("server connection failed - quitting")
				os.Exit(1)
			}
			continue
		}

		conn.Close()
		break
	}

	// option loop
	for {
		var request jsonRequest
		fmt.Println(`choose an option (1, 2, 3, q, etc...):
	1: GET query: "1"
	2: GET query: "2"
	3: GET query: "3"
	4: POST payload: {"1":"hello world", "2":"123", "3":true}
	5: DELETE query: "1"
	6: DELETE query: "2"
	7: DELETE query: "3"
	q: QUIT`)
		reader := bufio.NewReader(os.Stdin)

		fmt.Print("-> ")
		reqtype, _ := reader.ReadString('\n')
		reqtype = strings.Replace(reqtype, "\n", "", -1)

		switch reqtype {
		case "1":
			request = jsonRequest{
				Method: "GET",
				Query:  "1",
			}
		case "2":
			request = jsonRequest{
				Method: "GET",
				Query:  "2",
			}
		case "3":
			request = jsonRequest{
				Method: "GET",
				Query:  "3",
			}
		case "4":
			request = jsonRequest{
				Method: "POST",
				Payload: map[string]interface{}{
					"1": "hello world",
					"2": "123",
					"3": true},
			}
		case "5":
			request = jsonRequest{
				Method: "DELETE",
				Query:  "1",
			}
		case "6":
			request = jsonRequest{
				Method: "DELETE",
				Query:  "2",
			}
		case "7":
			request = jsonRequest{
				Method: "DELETE",
				Query:  "3",
			}
		case "q":
			fmt.Printf("quitting\n")
			os.Exit(0)
		default:
			fmt.Printf("\n\ninvalid option\n\n")
			continue
		}

		// server connection
		conn, err := net.Dial("tcp", ":8181")
		if err != nil {
			panic(err)
		}

		json.NewEncoder(conn).Encode(request)

		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			panic(err)
		}

		var response jsonResponse

		json.Unmarshal(buf[:n], &response)

		conn.Close()
		fmt.Println("\n\nresponse:")
		fmt.Printf("%+v\n\n:", response)
	}
}
