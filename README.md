# GoRBAC 权限管理库

一个简单灵活的 **基于角色的访问控制 (RBAC)** 系统，适用于 Go 项目。  
支持角色、权限、规则、用户分配、父子继承关系以及动态权限校验，可选缓存加速。

---

## 功能特点

- 定义 **角色（Role）** 和 **权限（Permission）**
- 给用户分配 **角色/权限**
- 支持 **角色与权限的父子继承关系**
- 支持 **规则（Rule）/执行器（Executor）** 实现动态权限校验
- 默认角色支持所有用户自动拥有
- 权限继承机制
- 可选缓存，提高权限校验性能

---

## 安装

```bash
go get github.com/kodar/gorbac
```

------

## 核心概念

### ItemType

```go
const (
    PermissionType ItemType = 2 // 权限
    RoleType       ItemType = 1 // 角色
)
```

### Item 接口

`Role` 和 `Permission` 都实现了 `Item` 接口：

```go
type Item interface {
    GetType() ItemType
    GetName() string
    GetDescription() string
    GetRuleName() string
    GetExecuteName() string
    GetCreateTime() time.Time
    GetUpdateTime() time.Time
}
```

### Assignment

用户与角色/权限的绑定：

```go
type Assignment struct {
    UserId     interface{}
    ItemName   string
    CreateTime time.Time
}
```

### Rule 与 Executor

- Rule 可绑定 Executor，在运行时动态判断权限
- Executor 实现接口：

```go
type Executor interface {
    Name() string
    Execute(ctx context.Context, userId interface{}, item Item) bool
}
```

------

## 快速使用示例

```go
package main

import (
    "context"
    "fmt"

    "github.com/yourusername/gorbac"
)

func main() {
    // 创建仓库（需要实现 AuthRepository 接口）
    repo := NewMemoryAuthRepository() // 你的实现

    // 初始化 RBAC 服务（开启缓存）
    service := gorbac.NewRbacService(repo, true)

    // 创建角色和权限
    service.AddRole("admin", "管理员角色", "")
    service.AddPermission("view_dashboard", "查看仪表盘权限", "")

    // 将权限分配给角色
    err := service.AssignChildren("admin", "view_dashboard")
    if err != nil {
        panic(err)
    }

    // 给用户分配角色
    userId := 123
    service.GetAuthManager().Assign(service.GetAuthManager().GetRole("admin"), userId)

    // 权限检查
    ctx := context.Background()
    if service.GetAuthManager().CheckAccess(ctx, userId, "view_dashboard") {
        fmt.Println("用户有权限查看仪表盘")
    } else {
        fmt.Println("权限不足")
    }
}
```

------

## API 概览

### RbacService

- `Roles() []*Role` – 获取所有角色
- `AddRole(name, desc, ruleName string) bool` – 添加角色
- `UpdateRole(name, newName, desc, ruleName string) bool` – 更新角色
- `DeleteRole(name string) bool` – 删除角色
- `Permissions() []*Permission` – 获取所有权限
- `AddPermission(name, desc, ruleName string) bool` – 添加权限
- `UpdatePermission(name, newName, desc, ruleName string) bool` – 更新权限
- `DeletePermission(name string) bool` – 删除权限
- `AssignChildren(parent string, children ...string) error` – 分配子项
- `CleanChildren(parent string) bool` – 清空子项
- `GetRolesByUser(userId interface{}) []*Role` – 获取用户角色
- `GetPermissionsByUser(userId interface{}) []*Permission` – 获取用户权限

### AuthManager

- `Assign(item Item, userId interface{}) *Assignment` – 给用户分配角色/权限
- `Revoke(item Item, userId interface{}) bool` – 回收角色/权限
- `CheckAccess(ctx context.Context, userId interface{}, permission string) bool` – 权限校验
- `SetDefaultRoles(roles ...*Role)` – 设置默认角色
- `GetDefaultRoles() []*Role` – 获取默认角色

------

## 高级特性

### 规则与执行器（Rule & Executor）

```go
type DemoExecutor struct {}

func (d *DemoExecutor) Name() string { return "demo" }
func (d *DemoExecutor) Execute(ctx context.Context, userId interface{}, item Item) bool {
    return true
}

// 注册执行器
gorbac.ExecuteManager.AddExecutor(&DemoExecutor{})
```

- 将 `RuleName` 绑定到角色或权限上
- 在 `CheckAccess` 时动态执行规则

------

## License

MIT License
