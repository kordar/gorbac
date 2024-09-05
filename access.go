package gorbac

import "context"

// Access 校验用户权限节点
type Access interface {
	CheckAccess(ctx context.Context, userId interface{}, permission string) bool
}
