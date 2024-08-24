package gorbac

import (
	"errors"
	"fmt"
	"github.com/kordar/gologger"
	"time"
)

type DefaultManager struct {
	mapper AuthRepository
	cache  bool
	// all auth items (name => Item)
	items map[string]Item
	// all auth rules (name => Rule)
	rules map[string]*Rule
	// a list of role names that are assigned to every user automatically without calling [[assign()]].
	// Note that these roles are applied to users, regardless of their state of authentication.
	defaultRoles map[string]*Role
	// auth item parent-child relationships (childName => list of parents)
	parents map[string][]string

	_checkAccessAssignments map[interface{}]map[string]*Assignment
}

func NewDefaultManager(mapper AuthRepository, cache bool) *DefaultManager {
	return &DefaultManager{
		mapper:                  mapper,
		cache:                   cache,
		_checkAccessAssignments: make(map[interface{}]map[string]*Assignment),
		defaultRoles:            make(map[string]*Role),
	}
}

func (manager *DefaultManager) invalidateCache() {
	if manager.cache {
		manager.items = make(map[string]Item)
		manager.rules = make(map[string]*Rule)
		manager.parents = make(map[string][]string)
	}
}

func (manager *DefaultManager) refreshInvalidateCache(operator bool) bool {
	if operator {
		manager.invalidateCache()
		return true
	}
	return false
}

func (manager *DefaultManager) getItem(name string) Item {
	if name == "" {
		return nil
	}

	target := manager.items[name]
	if target != nil {
		return target
	}

	item, err := manager.mapper.GetItem(name)
	if err != nil {
		return nil
	}
	return item
}

func (manager *DefaultManager) getItems(t ItemType) []Item {
	data, err := manager.mapper.GetItems(t.Value())
	if err != nil {
		return nil
	}

	var items []Item
	for _, item := range data {
		items = append(items, item)
	}
	return items
}

func (manager *DefaultManager) GetRule(name string) *Rule {
	if manager.rules != nil {
		return manager.rules[name]
	}

	rule, err := manager.mapper.GetRule(name)
	if err != nil {
		return nil
	}
	return rule
}

func (manager *DefaultManager) GetRules() []*Rule {
	length := len(manager.rules)
	if manager.rules != nil && length > 0 {
		rules := make([]*Rule, length)
		for _, rule := range manager.rules {
			rules = append(rules, rule)
		}
		return rules
	}

	data, err := manager.mapper.GetRules()
	if err != nil {
		return nil
	}
	var rules []*Rule
	for _, rule := range data {
		rules = append(rules, rule)
	}
	return rules
}

func (manager *DefaultManager) addItem(item Item) bool {
	err := manager.mapper.AddItem(item)
	return manager.refreshInvalidateCache(err == nil)
}

func (manager *DefaultManager) AddRule(rule Rule) bool {
	err := manager.mapper.AddRule(rule)
	return manager.refreshInvalidateCache(err == nil)
}

func (manager *DefaultManager) removeItem(item Item) bool {
	_ = manager.mapper.RemoveItem(item.GetName())
	manager.invalidateCache()
	return true
}

func (manager *DefaultManager) RemoveRule(rule Rule) bool {
	_ = manager.mapper.RemoveRule(rule.Name)
	manager.invalidateCache()
	return true
}

func (manager *DefaultManager) updateItem(name string, item Item) bool {
	err := manager.mapper.UpdateItem(name, item)
	return manager.refreshInvalidateCache(err == nil)
}

func (manager *DefaultManager) UpdateRule(name string, rule Rule) bool {
	err := manager.mapper.UpdateRule(name, rule)
	return manager.refreshInvalidateCache(err == nil)
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
	roles := make([]*Role, 1)
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
	if authAssignments, err := manager.mapper.FindAssignmentByUser(userId); err == nil {
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
	data, err := manager.mapper.FindChildren(name)
	if err != nil {
		return nil
	}

	var items []Item
	for _, item := range data {
		items = append(items, item)
	}
	return items
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
	authAssignments, err := manager.mapper.GetAssignmentByItems(roleName)
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
	manager.invalidateCache()
}

func (manager *DefaultManager) RemoveAllPermissions() {
	manager.removeAllItems(PermissionType)
}

func (manager *DefaultManager) RemoveAllRoles() {
	manager.removeAllItems(RoleType)
}

func (manager *DefaultManager) removeAllItems(t ItemType) {
	items, err := manager.mapper.GetItems(t.Value())
	if err != nil {
		logger.Warnf("[rbac] RemoveAllItems err = %v", err)
		return
	}

	if len(items) == 0 {
		return
	}

	names := make([]string, 0)
	for _, item := range items {
		names = append(names, item.GetName())
	}

	key := "parent"
	if t == PermissionType {
		key = "child"
	}

	_ = manager.mapper.RemoveChildByNames(key, names)
	_ = manager.mapper.RemoveAssignmentByName(names)
	_ = manager.mapper.RemoveItemByType(t.Value())
	manager.invalidateCache()
}

func (manager *DefaultManager) RemoveAllRules() {
	_ = manager.mapper.RemoveAllRules()
	manager.invalidateCache()
}

func (manager *DefaultManager) RemoveAllAssignments() {
	manager._checkAccessAssignments = make(map[interface{}]map[string]*Assignment)
	_ = manager.mapper.RemoveAllAssignments()
}

func (manager *DefaultManager) CheckAccess(userId interface{}, permissionName string, params map[string]interface{}) bool {
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

	if manager.items != nil {
		return manager.checkAccessFromCache(userId, permissionName, params, assignments)
	} else {
		return manager.checkAccessRecursive(userId, permissionName, params, assignments)
	}
}

func (manager *DefaultManager) loadFromCache() {
	if manager.items == nil || !manager.cache {
		logger.Info("[rbac] load from cache fail!!")
		return
	}

	manager.invalidateCache()

	rules, err2 := manager.mapper.GetRules()
	if err2 == nil {
		for _, rule := range rules {
			manager.rules[rule.Name] = NewRule(rule.Name, rule.ExecuteName, rule.CreateTime, rule.UpdateTime)
		}
	}

	authItems, err := manager.mapper.FindAllItems()
	if err != nil {
		logger.Warnf("[rbac] LoadFromCache [findAllItems err] = %v", err)
		return
	}

	for _, item := range authItems {
		manager.items[item.GetName()] = item
	}

	authItemChildren, err := manager.mapper.FindChildrenList()
	if err != nil {
		logger.Warnf("[rbac] LoadFromCache [FindChildrenList err] = %v", err)
		return
	}

	for _, authItemChild := range authItemChildren {
		child := authItemChild.Child
		if manager.items[child] != nil {
			if manager.parents[child] == nil {
				manager.parents[child] = make([]string, 1)
			}
			manager.parents[child] = append(manager.parents[child], authItemChild.Parent)
		}
	}
}

func (manager *DefaultManager) checkAccessFromCache(userId interface{}, itemName string, params map[string]interface{}, assignments map[string]*Assignment) bool {
	if manager.items[itemName] == nil {
		return false
	}

	item := manager.items[itemName]
	// logger.debug(item instanceof Role ? "Checking role: " + itemName : "Checking permission: " + itemName);
	if manager.executeRule(userId, item, params) == false {
		return false
	}

	if assignments[itemName] != nil || manager.defaultRoles[itemName] != nil {
		return true
	}

	parents := manager.parents[itemName]
	if parents != nil {
		for _, parent := range parents {
			if manager.checkAccessFromCache(userId, parent, params, assignments) {
				return true
			}
		}
	}

	return false
}

func (manager *DefaultManager) checkAccessRecursive(userId interface{}, itemName string, params map[string]interface{}, assignments map[string]*Assignment) bool {
	if manager.items[itemName] == nil {
		return false
	}

	item := manager.items[itemName]
	// logger.debug(item instanceof Role ? "Checking role: " + itemName : "Checking permission: " + itemName);
	if !manager.executeRule(userId, item, params) {
		return false
	}

	if assignments[itemName] != nil || manager.defaultRoles[itemName] != nil {
		return true
	}

	if authChildren, err := manager.mapper.FindChildrenFormChild(itemName); err == nil {
		for _, authChild := range authChildren {
			if manager.checkAccessRecursive(userId, authChild.Parent, params, assignments) {
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
	manager.invalidateCache()
	return err
}

func (manager *DefaultManager) Update(name string, item Item) bool {
	// TODO if the rule of the object is not alive in the system, then to create it to the system
	manager.checkRuleExits(item.GetRuleName())
	return manager.updateItem(name, item)
}

func (manager *DefaultManager) GetRole(name string) *Role {
	item := manager.getItem(name)
	if item == nil {
		return nil
	}
	return NewRole(item.GetName(), item.GetDescription(), item.GetRuleName(), item.GetExecuteName(), item.GetCreateTime(), item.GetUpdateTime())
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
	item := manager.getItem(name)
	if item == nil {
		return nil
	}
	return NewPermission(item.GetName(), item.GetDescription(), item.GetRuleName(), item.GetExecuteName(), item.GetCreateTime(), item.GetUpdateTime())
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

func (manager *DefaultManager) executeRule(userId interface{}, item Item, params map[string]interface{}) bool {
	if item.GetRuleName() == "" {
		return true
	}

	rule := manager.GetRule(item.GetRuleName())
	if rule == nil {
		logger.Warn("[rbac] Rule not found: " + item.GetRuleName())
		return false
	}

	return rule.GetExecutor().Execute(userId, item, params)
}

func (manager *DefaultManager) hasNoAssignments(assignments map[string]*Assignment) bool {
	return len(assignments) == 0 && len(manager.defaultRoles) == 0
}
