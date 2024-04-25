package base

import "time"

type Permission struct {
	Type        ItemType  `json:"item_type"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	RuleName    string    `json:"rule_name"`
	ExecuteName string    `json:"execute_name"`
	CreateTime  time.Time `json:"create_time"`
	UpdateTime  time.Time `json:"update_time"`
}

func (permission *Permission) GetType() ItemType {
	return permission.Type
}

func (permission *Permission) GetName() string {
	return permission.Name
}

func (permission *Permission) GetDescription() string {
	return permission.Description
}

func (permission *Permission) GetRuleName() string {
	return permission.RuleName
}

func (permission *Permission) GetExecuteName() string {
	return permission.ExecuteName
}

func (permission *Permission) GetCreateTime() time.Time {
	return permission.CreateTime
}

func (permission *Permission) GetUpdateTime() time.Time {
	return permission.UpdateTime
}

func NewPermission(name string, description string, ruleName string, executeName string, time time.Time, time2 time.Time) *Permission {
	return &Permission{
		Type:        PermissionType,
		Name:        name,
		Description: description,
		RuleName:    ruleName,
		ExecuteName: executeName,
		CreateTime:  time,
		UpdateTime:  time2,
	}
}
