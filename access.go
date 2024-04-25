package gorbac

// Access 校验用户权限节点
type Access interface {
	CheckAccess(userId interface{}, permissionName string, params map[string]interface{}) bool
}
