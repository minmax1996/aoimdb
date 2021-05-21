package msg_protocol

import "github.com/minmax1996/aoimdb/internal/aoimdb/datatypes"

type MsgPackRootMessage struct {
	Message string
	Error   error

	AuthResponse        *AuthResponse
	SelectResponse      *SelectResponse
	GetResponse         *GetResponse
	SetResponse         *SetResponse
	KeysResponse        *KeysResponse
	CreateTableResponse *CreateTableResponse
	InsertTableResponse *InsertTableResponse
	SelectTableResponse *SelectTableResponse
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

type KeysResponse struct {
	Keys []string
}

type CreateTableResponse struct {
	Message string
}

type InsertTableResponse struct {
	Message string
}

type SelectTableResponse struct {
	FieldNames []string
	Rows       []datatypes.Row
}
