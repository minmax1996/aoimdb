package msg_protocol

type MsgPackRootMessage struct {
	Message string
	Error   error

	AuthResponse   *AuthResponse
	SelectResponse *SelectResponse
	GetResponse    *GetResponse
	SetResponse    *SetResponse
}

type AuthResponse struct {
	Message string
}

type SelectResponse struct {
	SelectedDatabase string
}

type GetResponse struct {
	Key   string
	Value interface{}
}

type SetResponse struct {
	Message string
}
