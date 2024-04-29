package db

import "time"

// AuthRule 规则绑定，实现Execute接口完成特殊权限校验功能
type AuthRule struct {
	Name        string    `gorm:"type:varchar(64);primary_key" json:"name"`
	ExecuteName string    `json:"execute_name"`
	CreateTime  time.Time `json:"create_time"`
	UpdateTime  time.Time `json:"update_time"`
}

func (t *AuthRule) TableName() string {
	return tableNames["rule"]
}

// AuthItem 权限节点
type AuthItem struct {
	Name        string    `gorm:"column:name;primary_key" json:"name"`
	Type        int32     `gorm:"column:type" json:"type"`
	Description string    `gorm:"column:description" json:"description"`
	AuthRules   AuthRule  `gorm:"foreignkey:RuleName;association_foreignkey:Name" json:"auth_rules"`
	RuleName    string    `gorm:"column:rule_name" json:"rule_name"`
	ExecuteName string    `gorm:"column:execute_name" json:"execute_name"`
	CreateTime  time.Time `gorm:"column:create_time" json:"create_time"`
	UpdateTime  time.Time `gorm:"column:update_time" json:"update_time"`
}

func (t *AuthItem) TableName() string {
	return tableNames["item"]
}

// AuthItemChild 权限赋值关系
type AuthItemChild struct {
	AuthParent AuthItem `gorm:"foreignkey:Parent;association_foreignkey:Name" json:"auth_parent"`
	Parent     string   `gorm:"column:parent;primary_key" json:"parent"`
	AuthChild  AuthItem `gorm:"foreignkey:Child;association_foreignkey:Name" json:"auth_child"`
	Child      string   `gorm:"column:child;primary_key;index" json:"child"`
}

func (t *AuthItemChild) TableName() string {
	return tableNames["item-child"]
}

// AuthAssignment 用户赋权，userId->关联权限
type AuthAssignment struct {
	AuthItem   AuthItem    `gorm:"foreignkey:ItemName;association_foreignkey:Name" json:"auth_item"`
	ItemName   string      `gorm:"column:item_name;primary_key" json:"item_name"`
	UserId     interface{} `gorm:"column:user_id;type:varchar(32);primary_key;index" json:"user_id"`
	CreateTime time.Time   `gorm:"column:create_time" json:"create_time"`
}

func (t *AuthAssignment) TableName() string {
	return tableNames["assignment"]
}
