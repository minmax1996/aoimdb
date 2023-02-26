package datatypes

type ClientTypes string

const (
	ClientTypesString  ClientTypes = "str"
	ClientTypesInt32   ClientTypes = "int32"
	ClientTypesFloat32 ClientTypes = "float32"
)

func (c ClientTypes) String() string {
	return string(c)
}
