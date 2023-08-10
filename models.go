package gorbac

import "time"

var tableNames = map[string]string{
	"rule":       "auth_rule",
	"item":       "auth_item",
	"item-child": "auth_item_child",
	"assignment": "auth_assignment",
}

func SetTableName(key string, value string) {
	if tableNames[key] != "" {
		tableNames[key] = value
	}
}

type AuthRule struct {
	Name        string `gorm:"type:varchar(64);primary_key"`
	ExecuteName string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (t *AuthRule) TableName() string {
	return tableNames["rule"]
}

type AuthItem struct {
	Name        string   `gorm:"type:varchar(64);primary_key"`
	Type        int32    `gorm:"index"`
	Description string   `gorm:"text"`
	AuthRules   AuthRule `gorm:"foreignkey:RuleName;association_foreignkey:Name"`
	RuleName    string   `gorm:"type:varchar(64);index"`
	ExecuteName string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (t *AuthItem) TableName() string {
	return tableNames["item"]
}

type AuthItemChild struct {
	AuthParent AuthItem `gorm:"foreignkey:Parent;association_foreignkey:Name"`
	Parent     string   `gorm:"type:varchar(64);primary_key"`
	AuthChild  AuthItem `gorm:"foreignkey:Child;association_foreignkey:Name"`
	Child      string   `gorm:"type:varchar(64);primary_key;index"`
}

func (t *AuthItemChild) TableName() string {
	return tableNames["item-child"]
}

type AuthAssignment struct {
	AuthItem  AuthItem `gorm:"foreignkey:ItemName;association_foreignkey:Name"`
	ItemName  string   `gorm:"type:varchar(64);primary_key"`
	UserId    int64    `gorm:"primary_key;index"`
	CreatedAt time.Time
}

func (t *AuthAssignment) TableName() string {
	return tableNames["assignment"]
}
