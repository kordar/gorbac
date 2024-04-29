package db

import (
	"github.com/kordar/gorbac/base"
)

// ToItem AuthItem转item对象
func ToItem(item AuthItem) base.Item {
	if base.RoleType.Value() == item.Type {
		return base.NewRole(item.Name, item.Description, item.RuleName, item.ExecuteName, item.CreateTime, item.UpdateTime)
	} else {
		return base.NewPermission(item.Name, item.Description, item.RuleName, item.ExecuteName, item.CreateTime, item.UpdateTime)
	}
}

func ToRole(item AuthItem) base.Role {
	return *base.NewRole(item.Name, item.Description, item.RuleName, item.ExecuteName, item.CreateTime, item.UpdateTime)
}

func ToPermission(item AuthItem) base.Permission {
	return *base.NewPermission(item.Name, item.Description, item.RuleName, item.ExecuteName, item.CreateTime, item.UpdateTime)
}

func ToRule(rule AuthRule) base.Rule {
	return *base.NewRule(rule.Name, rule.ExecuteName, rule.CreateTime, rule.UpdateTime)
}

func ToAuthItem(item base.Item) AuthItem {
	return AuthItem{
		Name:        item.GetName(),
		Type:        item.GetType().Value(),
		Description: item.GetDescription(),
		RuleName:    item.GetRuleName(),
		ExecuteName: item.GetExecuteName(),
		CreateTime:  item.GetCreateTime(),
		UpdateTime:  item.GetUpdateTime(),
	}
}

func ToAuthItemChild(parent string, child string) AuthItemChild {
	return AuthItemChild{
		Parent: parent,
		Child:  child,
	}
}

func ToAuthAssignment(assignment base.Assignment) AuthAssignment {
	return AuthAssignment{
		ItemName:   assignment.ItemName,
		UserId:     assignment.UserId,
		CreateTime: assignment.CreateTime,
	}
}

func ToAssignment(authAssignment AuthAssignment) base.Assignment {
	return base.NewAssignment(authAssignment.UserId, authAssignment.ItemName)
}

func ToAuthRule(rule base.Rule) AuthRule {
	return AuthRule{
		Name:       rule.Name,
		CreateTime: rule.CreateTime,
		UpdateTime: rule.UpdateTime,
	}
}
