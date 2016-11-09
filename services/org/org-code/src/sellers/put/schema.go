package put

import ()

type Response struct {
	SeqId  int  `json:"seqId"`
	Result bool `json:"result"`
}

type StylokoRequest struct {
	Api    int         `json:"api"`
	Method string      `json:"method"`
	Params interface{} `json:"params"`
	Id     int         `json:"id"`
}

type Parameters struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SellerData struct {
	Value string `json:"value"`
	Type  string `json:"type"`
}
