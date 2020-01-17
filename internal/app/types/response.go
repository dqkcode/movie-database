package types

type (
	Response struct {
		Code  string
		Data  interface{}
		Error string
	}
)

const (
	CodeSuccess = "0000"
	CodeFail    = "1000"
)
