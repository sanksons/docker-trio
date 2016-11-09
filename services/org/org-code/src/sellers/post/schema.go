package post

import ()

type Response struct {
	Name  string `json:"name"`
	SeqId int    `json:"seqId"`
	Error string `json:"error"`
}
