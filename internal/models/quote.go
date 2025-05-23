package models

type Quote struct {
	Id     uint64 `json:"id"`
	Author string `json:"author"`
	Quote  string `json:"quote"`
}
