package protocols

type SelectRequest struct {
	Database string `json:"database"`
}

type SelectResponse struct {
	SelectedDatabase string `json:"selected_database"`
}

type KeysRequest struct {
}

type KeysResponse struct {
	Keys []string `json:"keys"`
}
