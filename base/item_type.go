package base

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

const (
	PermissionType ItemType = 2 // 权限
	RoleType       ItemType = 1 // 角色
)
