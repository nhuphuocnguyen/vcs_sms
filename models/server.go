package models

type Server struct {
	Server_id    string `json:"server_id"`
	Server_name  string `json:"server_name"`
	Status       string `json:"status"`
	Created_time int    `json:"created_time"`
	Last_updated int    `json:"last_updated"`
	Ipv4         string `json:"ipv4"`
}

type ImportExcel struct {
	Server_id   string `json:"server_id"`
	Server_name string `json:"server_name"`
	Status      string `json:"status"`
	Ipv4        string `json:"ipv4"`
}
