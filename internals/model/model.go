package model

type CreateConnectionRequest struct {
	Username string `json:"user"`
	Password string `json:"password"`
	ManagerName string `json:"manager"`
	Host string `json:"host"`
	Port string `json:"port"`
	ChannelName string `json:"channel"`
}
