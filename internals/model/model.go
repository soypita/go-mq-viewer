package model

// CreateConnectionRequest for establish new connection
type CreateConnectionRequest struct {
	Username    string `json:"user"`
	Password    string `json:"password"`
	ManagerName string `json:"manager"`
	Host        string `json:"host"`
	Port        string `json:"port"`
	ChannelName string `json:"channel"`
	Queue       string `json:"queue"`
}

// Message response for single message from queueu
type Message struct {
	MessageID   string `json:"messageId"`
	Timestamp   int64  `json:"timestamp"`
	MessageBody string `json:"msg"`
}
