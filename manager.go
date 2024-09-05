package gorbac

import (
	"context"
	"errors"
	"fmt"
	"github.com/kordar/gologger"
	"time"
)

type DefaultManager struct {
	mapper AuthRepository
	cache  *DefaultCache
	// a list of role names that are assigned to every user automatically without calling [[assign()]].
	// Note that these roles are applied to users, regardless of their state of authentication.
	defaultRoles            map[string]*Role
	_checkAccessAssignments map[interface{}]map[string]*Assignment
}

func NewDefaultManager(mapper AuthRepository, cache bool) *DefaultManager {
	return &DefaultManager{
		mapper:                  mapper,
		cache:                   NewDefaultCache(cache),
		_checkAccessAssignments: make(map[interface{}]map[string]*Assignment),
		defaultRoles:            make(map[string]*Role),
	}
}

func (manager *DefaultManager) GetItem(name string) Item {
	return manager.cache.GetItem(name, func(n string) Item {
		if item, err := manager.mapper.GetItem(name); err == nil {
			return item
		} else {
			return nil
		}
	})
}

func (manager *DefaultManager) getItems(itemType ItemType) []Item {
	items, err := manager.mapper.GetItemsByType(itemType)
	if err != nil {
		logger.Warnf("Error getting items: %s", err.Error())
		return make([]Item, 0)
	}
	return items
}

func (manager *DefaultManager) GetRule(name string) *Rule {
	return manager.cache.GetRule(name, func(n string) *Rule {
		if rule, err := manager.mapper.GetRule(name); err == nil {
			return rule
		} else {
			return nil
		}
	})
}

func (manager *DefaultManager) GetRules() []*Rule {
	return manager.cache.GetRules(func() []*Rule {
		data, err := manager.mapper.GetRules()
		if err != nil {
			return nil
		}
		var rules []*Rule
		for _, rule := range data {
			rules = append(rules, rule)
		}
		return rules
	})
}

func (manager *DefaultManager) addItem(item Item) bool {
	err := manager.mapper.AddItem(item)
	return manager.cache.refreshInvalidateCache(err == nil)
}

func (manager *DefaultManager) AddRule(rule Rule) bool {
	err := manager.mapper.AddRule(rule)
	return manager.cache.refreshInvalidateCache(err == nil)
}

func (manager *DefaultManager) removeItem(item Item) bool {
	_ = manager.mapper.RemoveItem(item.GetName())
	manager.cache.invalidateCache()
	return true
}

func (manager *DefaultManager) RemoveRule(rule Rule) bool {
	_ = manager.mapper.RemoveRule(rule.Name)
	manager.cache.invalidateCache()
	return true
}

func (manager *DefaultManager) updateItem(name string, item Item) bool {
	err := manager.mapper.UpdateItem(name, item)
	return manager.cache.refreshInvalidateCache(err == nil)
}

func (manager *DefaultManager) UpdateRule(name string, rule Rule) bool {
	err := manager.mapper.UpdateRule(name, rule)
	return manager.cache.refreshInvalidateCache(err == nil)
}

// GetRolesByUser 获取用户角色列表
func (manager *DefaultManager) GetRolesByUser(userId interface{}) []*Role {
	data, err := manager.mapper.FindRolesByUser(userId)
	if err != nil {
		return nil
	}
	var roles []*Role
	for _, item := range data {
		role := ToRole(item)
		roles = append(roles, &role)
	}
	return roles
}

// GetChildRoles 获取角色关联的子角色列表
func (manager *DefaultManager) GetChildRoles(roleName string) []*Role {
	role := manager.GetRole(roleName)
	if role == nil {
		logger.Infof("[rbac] Role %s not found manager.", roleName)
		return nil
	}

	result := make(map[string]bool)
	manager.getChildrenRecursive(roleName, manager.getChildrenList(), result)
	roles := make([]*Role, 0)
	roles = append(roles, role)

	for _, r := range manager.GetRoles() {
		if result[r.Name] == true {
			roles = append(roles, r)
		}
	}

	return roles
}

func (manager *DefaultManager) getChildrenRecursive(name string, childrenList map[string][]string, result map[string]bool) {
	if childrenList[name] != nil {
		// 存在子元素
		for _, childName := range childrenList[name] {
			result[childName] = true
			manager.getChildrenRecursive(childName, childrenList, result)
		}
	}
}

func (manager *DefaultManager) getChildrenList() map[string][]string {
	m := make(map[string][]string)
	list, err := manager.mapper.FindChildrenList()
	if err != nil {
		logger.Warnf("[rbac] get children list err=%v", err)
		return m
	}

	for _, child := range list {
		if m[child.Parent] == nil {
			m[child.Parent] = make([]string, 1)
		}
		m[child.Parent] = append(m[child.Parent], child.Child)
	}
	return m
}

func (manager *DefaultManager) GetPermissionsByRole(roleName string) []*Permission {
	childrenList := manager.getChildrenList()
	result := make(map[string]bool)
	manager.getChildrenRecursive(roleName, childrenList, result)

	permissions := make([]*Permission, 0)
	names := make([]string, 0)
	for name, exists := range result {
		if exists == true {
			names = append(names, name)
		}
	}

	if len(names) == 0 {
		return permissions
	}

	if list, err := manager.mapper.GetItemList(PermissionType.Value(), names); err == nil {
		for _, item := range list {
			permission := ToPermission(item)
			permissions = append(permissions, &permission)
		}
	}

	return permissions
}

func (manager *DefaultManager) GetPermissionsByUser(userId interface{}) []*Permission {
	directPermissions := manager.getDirectPermissionsByUser(userId)
	inheritedPermissions := manager.getInheritedPermissionsByUser(userId)
	return MergePermissions(directPermissions, inheritedPermissions)
}

// 直接关联的权限列表
func (manager *DefaultManager) getDirectPermissionsByUser(userId interface{}) []*Permission {
	permissions := make([]*Permission, 0)
	if data, err := manager.mapper.FindPermissionsByUser(userId); err == nil {
		for _, item := range data {
			permission := ToPermission(item)
			permissions = append(permissions, &permission)
		}
	}
	return permissions
}

func (manager *DefaultManager) getInheritedPermissionsByUser(userId interface{}) []*Permission {
	permissions := make([]*Permission, 0)
	if authAssignments, err := manager.mapper.FindAssignmentsByUser(userId); err == nil {
		childrenList := manager.getChildrenList()
		result := make(map[string]bool)
		for _, authAssignment := range authAssignments {
			manager.getChildrenRecursive(authAssignment.ItemName, childrenList, result)
		}

		names := make([]string, 0)
		for name, exists := range result {
			if exists == true {
				names = append(names, name)
			}
		}

		if len(names) == 0 {
			return permissions
		}

		if list, err2 := manager.mapper.GetItemList(PermissionType.Value(), names); err2 == nil {
			for _, item := range list {
				permission := ToPermission(item)
				permissions = append(permissions, &permission)
			}
		}

	}

	return permissions
}

func (manager *DefaultManager) CanAddChild(parent Item, child Item) bool {
	return !manager.detectLoop(parent, child)
}

// 递归遍历是否存在子元素是父元素本身，避免出现环
func (manager *DefaultManager) detectLoop(parent Item, child Item) bool {
	if child.GetName() == parent.GetName() {
		return true
	}

	children := manager.GetChildren(child.GetName())
	if children != nil {
		for _, child2 := range children {
			if manager.detectLoop(parent, child2) {
				return true
			}
		}
	}

	return false
}

func (manager *DefaultManager) AddChild(parent Item, child Item) error {

	if parent.GetName() == child.GetName() {
		return errors.New(fmt.Sprintf("Cannot add '%s' as a child of itself.", parent.GetName()))
	}

	if parent.GetType() == PermissionType && child.GetType() == RoleType {
		return errors.New(fmt.Sprintf("Cannot add a role as a child of a permission."))
	}

	if manager.detectLoop(parent, child) {
		return errors.New(fmt.Sprintf("Cannot add '%s' as a child of '%s'. A loop has been detected.", parent.GetName(), child.GetName()))
	}

	itemChild := NewItemChild(parent.GetName(), child.GetName())
	return manager.mapper.AddItemChild(*itemChild)
}

func (manager *DefaultManager) RemoveChild(parent Item, child Item) bool {
	err := manager.mapper.RemoveChild(parent.GetName(), child.GetName())
	return err == nil
}

func (manager *DefaultManager) RemoveChildren(parent Item) bool {
	err := manager.mapper.RemoveChildren(parent.GetName())
	return err == nil
}

func (manager *DefaultManager) HasChild(parent Item, child Item) bool {
	if parent == nil || child == nil {
		return false
	}
	return manager.mapper.HasChild(parent.GetName(), child.GetName())
}

func (manager *DefaultManager) GetChildren(name string) []Item {
	if data, err := manager.mapper.FindChildren(name); err == nil {
		return data
	} else {
		return make([]Item, 0)
	}
}

func (manager *DefaultManager) Assign(item Item, userId interface{}) *Assignment {
	assignment := NewAssignment(userId, item.GetName())
	err := manager.mapper.Assign(*assignment)
	if err == nil {
		delete(manager._checkAccessAssignments, userId)
		return assignment
	}
	return nil
}

func (manager *DefaultManager) Revoke(item Item, userId interface{}) bool {
	delete(manager._checkAccessAssignments, userId)
	err := manager.mapper.RemoveAssignment(userId, item.GetName())
	return err == nil
}

func (manager *DefaultManager) RevokeAll(userId interface{}) bool {
	delete(manager._checkAccessAssignments, userId)
	err := manager.mapper.RemoveAllAssignmentByUser(userId)
	return err == nil
}

func (manager *DefaultManager) GetAssignment(roleName string, userId interface{}) *Assignment {
	if assignment, err := manager.mapper.GetAssignment(userId, roleName); err == nil {
		return assignment
	} else {
		return nil
	}
}

func (manager *DefaultManager) GetAssignments(userId interface{}) map[string]*Assignment {
	assignments := make(map[string]*Assignment)
	authAssignments, err := manager.mapper.GetAssignments(userId)
	if err != nil {
		return assignments
	}

	for _, assignment := range authAssignments {
		assignments[assignment.ItemName] = assignment
	}
	return assignments
}

func (manager *DefaultManager) GetUserIdsByRole(roleName string) []interface{} {
	users := make([]interface{}, 0)
	authAssignments, err := manager.mapper.GetAssignmentsByItem(roleName)
	if err != nil {
		logger.Warnf("[rbac] GetUserIdsByRole err = %v", err)
		return users
	}

	for _, authAssignment := range authAssignments {
		users = append(users, authAssignment.UserId)
	}

	return users
}

func (manager *DefaultManager) RemoveAll() {
	_ = manager.mapper.RemoveAll()
	manager.cache.invalidateCache()
}

func (manager *DefaultManager) RemoveAllPermissions() {
	manager.removeAllItems(PermissionType)
}

func (manager *DefaultManager) RemoveAllRoles() {
	manager.removeAllItems(RoleType)
}

func (manager *DefaultManager) removeAllItems(itemType ItemType) {
	items := manager.getItems(itemType)
	if len(items) == 0 {
		return
	}

	names := make([]string, 0)
	for _, item := range items {
		names = append(names, item.GetName())
	}

	_ = manager.mapper.RemoveChildByNames(itemType, names)
	_ = manager.mapper.RemoveAssignmentByNames(names)
	_ = manager.mapper.RemoveItemByType(itemType)
	manager.cache.invalidateCache()
}

func (manager *DefaultManager) RemoveAllRules() {
	_ = manager.mapper.RemoveAllRules()
	manager.cache.invalidateCache()
}

func (manager *DefaultManager) RemoveAllAssignments() {
	manager._checkAccessAssignments = make(map[interface{}]map[string]*Assignment)
	_ = manager.mapper.RemoveAllAssignments()
}

func (manager *DefaultManager) CheckAccess(ctx context.Context, userId interface{}, permissionName string) bool {
	var assignments map[string]*Assignment
	if manager._checkAccessAssignments[userId] != nil {
		assignments = manager._checkAccessAssignments[userId]
	} else {
		assignments = manager.GetAssignments(userId)
		manager._checkAccessAssignments[userId] = assignments
	}

	if manager.hasNoAssignments(assignments) {
		return false
	}

	manager.loadFromCache()

	if manager.cache.items != nil {
		return manager.checkAccessFromCache(ctx, userId, permissionName, assignments)
	} else {
		return manager.checkAccessRecursive(ctx, userId, permissionName, assignments)
	}
}

func (manager *DefaultManager) loadFromCache() {
	if manager.cache.items == nil || !manager.cache.cache {
		logger.Info("[rbac] load from cache fail!!")
		return
	}

	manager.cache.invalidateCache()

	rules, err2 := manager.mapper.GetRules()
	if err2 == nil {
		for _, rule := range rules {
			manager.cache.rules[rule.Name] = NewRule(rule.Name, rule.ExecuteName, rule.CreateTime, rule.UpdateTime)
		}
	}

	authItems, err := manager.mapper.FindAllItems()
	if err != nil {
		logger.Warnf("[rbac] LoadFromCache [findAllItems err] = %v", err)
		return
	}

	for _, item := range authItems {
		manager.cache.items[item.GetName()] = item
	}

	authItemChildren, err := manager.mapper.FindChildrenList()
	if err != nil {
		logger.Warnf("[rbac] LoadFromCache [FindChildrenList err] = %v", err)
		return
	}

	for _, authItemChild := range authItemChildren {
		child := authItemChild.Child
		if manager.cache.items[child] != nil {
			if manager.cache.parents[child] == nil {
				manager.cache.parents[child] = make([]string, 0)
			}
			manager.cache.parents[child] = append(manager.cache.parents[child], authItemChild.Parent)
		}
	}
}

func (manager *DefaultManager) checkAccessFromCache(ctx context.Context, userId interface{}, itemName string, assignments map[string]*Assignment) bool {
	if manager.cache.items[itemName] == nil {
		return false
	}

	item := manager.cache.items[itemName]
	// logger.debug(item instanceof Role ? "Checking role: " + itemName : "Checking permission: " + itemName);
	if manager.executeRule(ctx, userId, item) == false {
		return false
	}

	if assignments[itemName] != nil || manager.defaultRoles[itemName] != nil {
		return true
	}

	parents := manager.cache.parents[itemName]
	if parents != nil {
		for _, parent := range parents {
			if manager.checkAccessFromCache(ctx, userId, parent, assignments) {
				return true
			}
		}
	}

	return false
}

func (manager *DefaultManager) checkAccessRecursive(ctx context.Context, userId interface{}, itemName string, assignments map[string]*Assignment) bool {
	if manager.cache.items[itemName] == nil {
		return false
	}

	item := manager.cache.items[itemName]
	// logger.debug(item instanceof Role ? "Checking role: " + itemName : "Checking permission: " + itemName);
	if !manager.executeRule(ctx, userId, item) {
		return false
	}

	if assignments[itemName] != nil || manager.defaultRoles[itemName] != nil {
		return true
	}

	if authChildren, err := manager.mapper.FindChildrenFormChild(itemName); err == nil {
		for _, authChild := range authChildren {
			if manager.checkAccessRecursive(ctx, userId, authChild.Parent, assignments) {
				return true
			}
		}
	}

	return false
}

func (manager *DefaultManager) CreateRole(name string) *Role {
	return NewRole(name, "", "", "", time.Now(), time.Now())
}

func (manager *DefaultManager) CreatePermission(name string) *Permission {
	return NewPermission(name, "", "", "", time.Now(), time.Now())
}

func (manager *DefaultManager) Add(item Item) bool {
	// TODO if the rule of the object is not alive in the system, then to create it to the system
	manager.checkRuleExits(item.GetRuleName())
	return manager.addItem(item)
}

func (manager *DefaultManager) Remove(item Item) bool {
	return manager.removeItem(item)
}

func (manager *DefaultManager) RemoveAllAssignmentByUser(userId interface{}) error {
	err := manager.mapper.RemoveAllAssignmentByUser(userId)
	manager.cache.invalidateCache()
	return err
}

func (manager *DefaultManager) Update(name string, item Item) bool {
	// TODO if the rule of the object is not alive in the system, then to create it to the system
	manager.checkRuleExits(item.GetRuleName())
	return manager.updateItem(name, item)
}

func (manager *DefaultManager) GetRole(name string) *Role {
	item := manager.GetItem(name)
	if item == nil {
		return nil
	}
	if item.GetType() != RoleType {
		return nil
	} else {
		role := ToRole(item)
		return &role
	}
}

func (manager *DefaultManager) GetRoles() []*Role {
	var roles []*Role
	items := manager.getItems(RoleType)
	for _, item := range items {
		role := NewRole(item.GetName(), item.GetDescription(), item.GetRuleName(), item.GetExecuteName(), item.GetCreateTime(), item.GetUpdateTime())
		roles = append(roles, role)
	}
	return roles
}

func (manager *DefaultManager) GetPermission(name string) *Permission {
	item := manager.GetItem(name)
	if item == nil {
		return nil
	}
	if item.GetType() != PermissionType {
		return nil
	} else {
		permission := ToPermission(item)
		return &permission
	}
}

func (manager *DefaultManager) GetPermissions() []*Permission {
	var permissions []*Permission
	items := manager.getItems(PermissionType)
	for _, item := range items {
		permission := NewPermission(item.GetName(), item.GetDescription(), item.GetRuleName(), item.GetExecuteName(), item.GetCreateTime(), item.GetUpdateTime())
		permissions = append(permissions, permission)
	}
	return permissions
}

func (manager *DefaultManager) checkRuleExits(name string) {
	if name != "" && manager.GetRule(name) == nil {
		rule := NewRule(name, "", time.Now(), time.Now())
		manager.AddRule(*rule)
	}
}

func (manager *DefaultManager) SetDefaultRoles(roles ...*Role) {
	for _, role := range roles {
		manager.defaultRoles[role.GetName()] = role
	}
}

func (manager *DefaultManager) getDefaultRoles() map[string]*Role {
	return manager.defaultRoles
}

func (manager *DefaultManager) executeRule(ctx context.Context, userId interface{}, item Item) bool {
	if item.GetRuleName() == "" {
		return true
	}

	rule := manager.GetRule(item.GetRuleName())
	if rule == nil {
		logger.Warn("[rbac] Rule not found: " + item.GetRuleName())
		return false
	}

	return rule.GetExecutor().Execute(ctx, userId, item)
}

func (manager *DefaultManager) hasNoAssignments(assignments map[string]*Assignment) bool {
	return len(assignments) == 0 && len(manager.defaultRoles) == 0
}
