package gorbac

import (
	"log"
	"time"
)

type Executor interface {
	Name() string
	Execute(userId int64, item Item, params map[string]interface{}) bool
}

type DemoExecutor struct {
}

func (d *DemoExecutor) Name() string {
	return "demo"
}

func (d *DemoExecutor) Execute(userId int64, item Item, params map[string]interface{}) bool {
	log.Println("============================================")
	log.Println("==============DEMO===============")
	log.Println("============================================")
	return true
}

// ExecuteManager /******************execute manger*****************************/
var ExecuteManager = container{
	content: make(map[string]Executor),
}

type container struct {
	content map[string]Executor
}

func (container *container) AddExecutor(executor Executor) {
	container.content[executor.Name()] = executor
}

func (container *container) GetExecutor(name string) Executor {
	return container.content[name]
}

/*
 executor item
*/
type ItemType int32

func (t ItemType) Value() int32 {
	switch t {
	case PermissionType:
		return 1
	case RoleType:
		return 2
	}
	return 0
}

const (
	PermissionType ItemType = 1 // 角色
	RoleType       ItemType = 2 // 权限
)

type Item interface {
	GetType() ItemType
	GetName() string
	GetDescription() string
	GetRuleName() string
	GetExecuteName() string
	GetCreateTime() time.Time
	GetUpdateTime() time.Time
}

type Permission struct {
	itemType    ItemType
	Name        string
	Description string
	RuleName    string
	ExecuteName string
	CreateTime  time.Time
	UpdateTime  time.Time
}

func (permission *Permission) GetType() ItemType {
	return permission.itemType
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
		itemType:    PermissionType,
		Name:        name,
		Description: description,
		RuleName:    ruleName,
		ExecuteName: executeName,
		CreateTime:  time,
		UpdateTime:  time2,
	}
}

type Role struct {
	itemType    ItemType
	Name        string
	Description string
	RuleName    string
	ExecuteName string
	CreateTime  time.Time
	UpdateTime  time.Time
}

func NewRole(name string, description string, ruleName string, executeName string, time time.Time, time2 time.Time) *Role {
	return &Role{
		itemType:    RoleType,
		Name:        name,
		Description: description,
		RuleName:    ruleName,
		ExecuteName: executeName,
		CreateTime:  time,
		UpdateTime:  time2,
	}
}

func (role *Role) GetType() ItemType {
	return role.itemType
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

type Assignment struct {
	UserId     int64
	ItemName   string
	CreateTime time.Time
}

func NewAssignment(userId int64, itemName string) Assignment {
	return Assignment{UserId: userId, ItemName: itemName, CreateTime: time.Now()}
}

/**
rule
*/
type Rule struct {
	/**
	 * string name of the rule
	 */
	Name string
	/**
	 * 执行函数
	 */
	ExecuteName string
	/**
	 * int UNIX timestamp representing the rule creation time
	 */
	CreateTime time.Time
	/**
	 * int UNIX timestamp representing the rule updating time
	 */
	UpdateTime time.Time
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
