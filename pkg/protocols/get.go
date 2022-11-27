package protocols

type GetRequest struct {
	Key string `json:"key"`
}

type GetResponse struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}
