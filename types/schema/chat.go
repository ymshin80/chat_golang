package schema

import "time"

type Chat struct{
	ID      int64  `json:"id"`
	Room    string `json:"room"`
	Name    string `json:"name"`
	Message string `json:"message"`
	SendDtm    time.Time `json:"sendDtm"`
}

