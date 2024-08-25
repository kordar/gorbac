package gorbac

import "time"

const (
	NoneType       ItemType = 0
	PermissionType ItemType = 2 // 权限
	RoleType       ItemType = 1 // 角色
)

// ItemType /*
type ItemType int32

func (t ItemType) Value() int32 {
	switch t {
	case PermissionType:
		return 2
	case RoleType:
		return 1
	}
	return 0
}

type Item interface {
	GetType() ItemType
	GetName() string
	GetDescription() string
	GetRuleName() string
	GetExecuteName() string
	GetCreateTime() time.Time
	GetUpdateTime() time.Time
}

type Rule struct {
	Name        string    `json:"name"`
	ExecuteName string    `json:"execute_name"`
	CreateTime  time.Time `json:"create_time"`
	UpdateTime  time.Time `json:"update_time"`
}

func NewRule(name string, executeName string, createTime time.Time, updateTime time.Time) *Rule {
	return &Rule{Name: name, ExecuteName: executeName, CreateTime: createTime, UpdateTime: updateTime}
}

func (rule *Rule) SetName(name string) {
	rule.Name = name
}

func (rule *Rule) SetExecutor(executor Executor) {
	rule.ExecuteName = executor.Name()
}

func (rule *Rule) GetExecutor() Executor {
	return ExecuteManager.GetExecutor(rule.ExecuteName)
}

// ------------- role

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

// --------------- permission

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

type Assignment struct {
	UserId     interface{} `json:"user_id"`
	ItemName   string      `json:"item_name"`
	CreateTime time.Time   `json:"create_time"`
}

func NewAssignment(userId interface{}, itemName string) *Assignment {
	return &Assignment{UserId: userId, ItemName: itemName, CreateTime: time.Now()}
}

type ItemChild struct {
	Parent string `json:"parent"`
	Child  string `json:"child"`
}

func NewItemChild(parent string, child string) *ItemChild {
	return &ItemChild{Parent: parent, Child: child}
}
