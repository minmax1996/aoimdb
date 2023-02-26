package protocols

type CreateTableRequest struct {
}

type InsertTableRequst struct {
}

type SelectTableRequest struct {
	Table string
}

type CreateTableResponse struct {
	Message string
}

type InsertTableResponse struct {
	Message string
}

type SelectTableResponse struct {
	FieldNames []string
	Rows       [][]string
}
