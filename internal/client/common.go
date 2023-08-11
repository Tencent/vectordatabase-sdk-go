package client

type CommmonResponse struct {
	Code int32  `json:"code,omitempty"`
	Msg  string `json:"msg,omitempty"`
}
