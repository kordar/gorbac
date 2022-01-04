package gorbac

import "time"

type AuthRule struct {
	Name        string `gorm:"type:varchar(64);primary_key"`
	ExecuteName string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (t *AuthRule) TableName() string {
	return "auth_rule"
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
	return "auth_item"
}

type AuthItemChild struct {
	AuthParent AuthItem `gorm:"foreignkey:Parent;association_foreignkey:Name"`
	Parent     string   `gorm:"type:varchar(64);primary_key"`
	AuthChild  AuthItem `gorm:"foreignkey:Child;association_foreignkey:Name"`
	Child      string   `gorm:"type:varchar(64);primary_key;index"`
}

func (t *AuthItemChild) TableName() string {
	return "auth_item_child"
}

type AuthAssignment struct {
	AuthItem  AuthItem `gorm:"foreignkey:ItemName;association_foreignkey:Name"`
	ItemName  string   `gorm:"type:varchar(64);primary_key"`
	UserId    int64   `gorm:"primary_key;index"`
	CreatedAt time.Time
}

func (t *AuthAssignment) TableName() string {
	return "auth_assignment"
}