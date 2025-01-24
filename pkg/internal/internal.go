package internal

type User struct {
	ID       int    `json:"id,string,omitempty"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Note struct {
	ID     int    `json:"id,string,omitempty"`
	Name   string `json:"name"`
	UserId string `json:"user_id,string"`
	Text   string `json:"text"`
}
