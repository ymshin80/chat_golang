package types

import "strings"


type header struct {
	Result int    `json:"result"`
	Data   string `json:"data"`
}

func newHeader(result int, data ...string) *header {
	return &header{
		Result: result,
		Data:   strings.Join(data, ","),
	}
}


type response struct {
	*header
	Result interface{} `json:"result"`
}

func NewRes(result int, res interface{}, data ...string) *response {
	return &response{
		header: newHeader(result, data...),
		Result: res,
	}
}