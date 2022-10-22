package models

type Server struct {
	Server_id    int    `json:"server_id"`
	Server_name  string `json:"server_name"`
	Status       bool   `json:"status"`
	Created_time int    `json:"created_time"`
	Last_updated int    `json:"last_updated"`
	Ipv4         string `json:"ipv4"`
}
