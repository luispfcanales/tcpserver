package model

type MessageTCP struct {
	From    string
	Payload []byte
}

type RequestMessage struct {
	EventName string `json:"event_name,omitempty"`
	Message   string `json:"message,omitempty"`
}

type ResponseMessage struct {
	Status  int         `json:"status,omitempty"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}
