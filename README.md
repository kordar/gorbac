# gorbac

`RBAC`（`Role-Based Access Control`）权限模型，即：基于角色的权限控制。通过角色关联用户，角色关联权限的方式间接赋予用户权限。

## 接口定义

### 权限数据层接口
（基于`mysql`、`sqlite`、`redis`等作为权限数据层），可通过实现该接口自定义数据层存储逻辑。**相关实现类：**

- [`gorbac-gorm`](https://github.com/kordar/gorbac-gorm)基于[`gorm`](https://github.com/go-gorm/gorm)实现权限数据层
- [`gorbac-redis`](https://github.com/kordar/gorbac-redis)基于[`go-redis`](https://github.com/redis/go-redis)实现权限数据层

```go
type AuthRepository interface {}
```

### 权限接口

该接口定义了`rbac`模型的所有功能方法，同时`DefaultManager`实现了该接口

```go
type AuthManager interface {} 
```

- 使用方式

```go
DefaultManager(AuthRepository, Cache)
```


### 权限校验接口

```go
type Access interface {
    CheckAccess(ctx context.Context, userId interface{}, permission string) bool
}
```

## 基于`rule`自定义实现权限控制

`rbac`权限模型默认仅控制到权限节点，一般权限节点关联访问资源`URI`，如果想进行例如具体数据记录的权限控制，需要自定义`rule`实现类进行访问控制。开发使用如下：

```go
// 1、实现Execute接口
type Executor interface {
    Name() string
    Execute(ctx context.Context, userId interface{}, item Item) bool
}

// 2、添加实现类
AddExecutor(executor Executor)
```

注：`RuleName`关联在`Item`属性下。

## `RbacService`使用

对常用权限功能进行包装，满足日常绝大多数使用场景，开箱即用。

```go
func NewRbacService(mgr AuthRepository, cache bool) *RbacService {
    return &RbacService{mgr: NewDefaultManager(mgr, cache)}
}
```









