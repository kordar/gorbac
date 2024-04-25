package base

func MergePermissions(list1 []*base.Permission, list2 []*base.Permission) []*base.Permission {
	m := make(map[string]*base.Permission)
	for _, permission := range list1 {
		m[permission.Name] = permission
	}
	for _, permission := range list2 {
		m[permission.Name] = permission
	}
	permissions := make([]*base.Permission, len(m))
	for _, permission := range m {
		permissions = append(permissions, permission)
	}
	return permissions
}
