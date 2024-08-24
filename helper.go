package gorbac

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

func ToRole(item Item) Role {
	return *NewRole(item.GetName(), item.GetDescription(), item.GetRuleName(), item.GetExecuteName(), item.GetCreateTime(), item.GetUpdateTime())
}

func ToPermission(item Item) Permission {
	return *NewPermission(item.GetName(), item.GetDescription(), item.GetRuleName(), item.GetExecuteName(), item.GetCreateTime(), item.GetUpdateTime())
}
