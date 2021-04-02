package dto

type CardDTO struct {
	Id int64 `json:"id"`
	Type string `json:"type"`
	Number string `json:"number"`
}

type Result struct{
	Result string `json:"result,omitempty"`
	ErrorDescription string `json:"errorDesc,omitempty"`
	Cards []*CardDTO `json:"cards,omitempty"`
}