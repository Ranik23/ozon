package item

import (
	"encoding/json"
	"log/slog"
	"time"
)


type Item struct {
	Id 				int			`json:"id"`
	Receiver_Id 	int			`json:"receiver_id"`
	Expiration_Date time.Time	`json:"expiration"`
	Received 		bool		`json:"received"`
	Status 			bool		`json:"status"`
}


func (i Item) String() string {
	jsonData, err := json.MarshalIndent(i, "", "	")
	if err != nil {
		slog.Error("error marshalling", slog.String("err", err.Error()))
		return ""
	}
	return string(jsonData)
}


