package protocols

type SetRequest struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

type SetResponse struct {
	Message string `json:"message"`
}
