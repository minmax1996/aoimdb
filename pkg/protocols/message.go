package protocols

type Request struct {
	MessageType MessageType
	Message     interface{}
}

type Response struct {
	Error       error
	Message     string
	MessageType MessageType

	AuthResponse        *AuthResponse
	SelectResponse      *SelectResponse
	GetResponse         *GetResponse
	SetResponse         *SetResponse
	KeysResponse        *KeysResponse
	CreateTableResponse *CreateTableResponse
	InsertTableResponse *InsertTableResponse
	SelectTableResponse *SelectTableResponse
}

type MessageType string

const (
	MessageTypeAuth        MessageType = "auth"
	MessageTypeSelect      MessageType = "select"
	MessageTypeGet         MessageType = "get"
	MessageTypeSet         MessageType = "set"
	MessageTypeKeys        MessageType = "keys"
	MessageTypeCreateTable MessageType = "create_table"
	MessageTypeInsertTable MessageType = "insert_table"
	MessageTypeSelectTable MessageType = "select_table"
)
