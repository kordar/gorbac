package base

import "time"

type Item interface {
	GetType() ItemType
	GetName() string
	GetDescription() string
	GetRuleName() string
	GetExecuteName() string
	GetCreateTime() time.Time
	GetUpdateTime() time.Time
}
