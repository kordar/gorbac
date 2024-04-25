package base

import "time"

type Assignment struct {
	UserId     interface{} `json:"user_id"`
	ItemName   string      `json:"item_name"`
	CreateTime time.Time   `json:"create_time"`
}

func NewAssignment(userId interface{}, itemName string) Assignment {
	return Assignment{UserId: userId, ItemName: itemName, CreateTime: time.Now()}
}
