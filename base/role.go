package base

import "time"

type Role struct {
	Type        ItemType  `json:"item_type"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	RuleName    string    `json:"rule_name"`
	ExecuteName string    `json:"execute_name"`
	CreateTime  time.Time `json:"create_time"`
	UpdateTime  time.Time `json:"update_time"`
}

func NewRole(name string, description string, ruleName string, executeName string, time time.Time, time2 time.Time) *Role {
	return &Role{
		Type:        RoleType,
		Name:        name,
		Description: description,
		RuleName:    ruleName,
		ExecuteName: executeName,
		CreateTime:  time,
		UpdateTime:  time2,
	}
}

func (role *Role) GetType() ItemType {
	return role.Type
}

func (role *Role) GetName() string {
	return role.Name
}

func (role *Role) GetDescription() string {
	return role.Description
}

func (role *Role) GetRuleName() string {
	return role.RuleName
}

func (role *Role) GetExecuteName() string {
	return role.ExecuteName
}

func (role *Role) GetCreateTime() time.Time {
	return role.CreateTime
}

func (role *Role) GetUpdateTime() time.Time {
	return role.UpdateTime
}
