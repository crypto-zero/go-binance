package common

type WebsocketRequest struct {
	ID     uint64      `json:"id"`
	Method string      `json:"method"`
	Params interface{} `json:"params"`
}

type WebsocketReply struct {
	ID      uint64 `json:"id"`
	Code    int64  `json:"code"`
	Message string `json:"msg"`
	Result  interface{}
}

func (wsr *WebsocketReply) OK() error {
	if wsr.Code == 0 {
		return nil
	}
	return &APIError{Code: wsr.Code, Message: wsr.Message}
}
