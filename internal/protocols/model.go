package protocols

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
