package gorbac

func ConversionToItem(item AuthItem) Item {
	switch item.Type {
	case RoleType.Value():
		return NewRole(item.Name, item.Description, item.RuleName, item.ExecuteName, item.CreatedAt, item.UpdatedAt)
	}

	return NewPermission(item.Name, item.Description, item.RuleName, item.ExecuteName, item.CreatedAt, item.UpdatedAt)
}

func ConversionToRole(item AuthItem) Role {
	return *NewRole(item.Name, item.Description, item.RuleName, item.ExecuteName, item.CreatedAt, item.UpdatedAt)
}

func ConversionToPermission(item AuthItem) Permission {
	return *NewPermission(item.Name, item.Description, item.RuleName, item.ExecuteName, item.CreatedAt, item.UpdatedAt)
}

func ConversionToRule(rule AuthRule) Rule {
	return *NewRule(rule.Name, rule.ExecuteName, rule.CreatedAt, rule.UpdatedAt)
}

func ConversionToAuthItem(item Item) AuthItem {
	return AuthItem{
		Name:        item.GetName(),
		Type:        item.GetType().Value(),
		Description: item.GetDescription(),
		RuleName:    item.GetRuleName(),
		ExecuteName: item.GetExecuteName(),
		CreatedAt:   item.GetCreateTime(),
		UpdatedAt:   item.GetUpdateTime(),
	}
}

func ConversionToAuthItemChild(parent string, child string) AuthItemChild {
	return AuthItemChild{
		Parent: parent,
		Child:  child,
	}
}

func ConversionToAuthAssignment(assignment Assignment) AuthAssignment {
	return AuthAssignment{
		ItemName:  assignment.ItemName,
		UserId:    assignment.UserId,
		CreatedAt: assignment.CreateTime,
	}
}

func ConversionToAssignment(authAssignment AuthAssignment) Assignment {
	return NewAssignment(authAssignment.UserId, authAssignment.ItemName)
}

func ConversionToAuthRule(rule Rule) AuthRule {
	return AuthRule{
		Name:      rule.Name,
		CreatedAt: rule.CreateTime,
		UpdatedAt: rule.UpdateTime,
	}
}

func MergePermissions(list1 []*Permission, list2 []*Permission) []*Permission {
	m := make(map[string]*Permission)
	for _, permission := range list1 {
		m[permission.Name] = permission
	}
	for _, permission := range list2 {
		m[permission.Name] = permission
	}
	permissions := make([]*Permission, len(m))
	for _, permission := range m {
		permissions = append(permissions, permission)
	}
	return permissions
}
