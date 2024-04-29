package gorbac

import (
	"errors"
	"fmt"
	log "github.com/kordar/gologger"
	"github.com/kordar/gorbac/base"
	"github.com/kordar/gorbac/db"
	"time"
)

type DbManager struct {
	/**
	 */
	mapper db.AuthRepository
	/**

	 */
	cache bool
	/**
	 * Item[] all auth items (name => Item)
	 */
	items map[string]base.Item
	/**
	 * Rule[] all auth rules (name => Rule)
	 */
	rules map[string]*base.Rule
	/**
	 * @var array a list of role names that are assigned to every user automatically without calling [[assign()]].
	 * Note that these roles are applied to users, regardless of their state of authentication.
	 */
	defaultRoles map[string]*base.Role
	/**
	 * array auth item parent-child relationships (childName => list of parents)
	 */
	parents map[string][]string

	_checkAccessAssignments map[interface{}]map[string]*base.Assignment
}

func NewDbManager(mapper db.AuthRepository, cache bool) *DbManager {
	return &DbManager{
		mapper:                  mapper,
		cache:                   cache,
		_checkAccessAssignments: make(map[interface{}]map[string]*base.Assignment),
		defaultRoles:            make(map[string]*base.Role),
	}
}

func (manager *DbManager) invalidateCache() {
	if manager.cache {
		manager.items = make(map[string]base.Item)
		manager.rules = make(map[string]*base.Rule)
		manager.parents = make(map[string][]string)
	}
}

func (manager *DbManager) refreshInvalidateCache(operator bool) bool {
	if operator {
		manager.invalidateCache()
		return true
	}
	return false
}

func (manager *DbManager) getItem(name string) base.Item {
	if name == "" {
		return nil
	}

	target := manager.items[name]
	if target != nil {
		return target
	}

	authItem, err := manager.mapper.GetItem(name)
	if err != nil {
		return nil
	}
	item := db.ToItem(*authItem)
	return item
}

func (manager *DbManager) getItems(t base.ItemType) []base.Item {
	data, err := manager.mapper.GetItems(t.Value())
	if err != nil {
		return nil
	}

	var items []base.Item
	for _, authItem := range data {
		item := db.ToItem(*authItem)
		items = append(items, item)
	}
	return items
}

func (manager *DbManager) GetRule(name string) *base.Rule {
	if manager.rules != nil {
		return manager.rules[name]
	}

	authRule, err := manager.mapper.GetRule(name)
	if err != nil {
		return nil
	}
	rule := db.ToRule(*authRule)
	return &rule
}

func (manager *DbManager) GetRules() []*base.Rule {
	length := len(manager.rules)
	if manager.rules != nil && length > 0 {
		rules := make([]*base.Rule, length)
		for _, rule := range manager.rules {
			rules = append(rules, rule)
		}
		return rules
	}

	data, err := manager.mapper.GetRules()
	if err != nil {
		return nil
	}
	var rules []*base.Rule
	for _, authRule := range data {
		rule := db.ToRule(*authRule)
		rules = append(rules, &rule)
	}
	return rules
}

func (manager *DbManager) addItem(item base.Item) bool {
	authItem := db.ToAuthItem(item)
	err := manager.mapper.AddItem(authItem)
	return manager.refreshInvalidateCache(err == nil)
}

func (manager *DbManager) AddRule(rule base.Rule) bool {
	authRule := db.ToAuthRule(rule)
	err := manager.mapper.AddRule(authRule)
	return manager.refreshInvalidateCache(err == nil)
}

func (manager *DbManager) removeItem(item base.Item) bool {
	_ = manager.mapper.RemoveItem(item.GetName())
	manager.invalidateCache()
	return true
}

func (manager *DbManager) RemoveRule(rule base.Rule) bool {
	_ = manager.mapper.RemoveRule(rule.Name)
	manager.invalidateCache()
	return true
}

func (manager *DbManager) updateItem(name string, item base.Item) bool {
	authItem := db.ToAuthItem(item)
	err := manager.mapper.UpdateItem(name, authItem)
	return manager.refreshInvalidateCache(err == nil)
}

func (manager *DbManager) UpdateRule(name string, rule base.Rule) bool {
	authRule := db.ToAuthRule(rule)
	err := manager.mapper.UpdateRule(name, authRule)
	return manager.refreshInvalidateCache(err == nil)
}

// GetRolesByUser 获取用户角色列表
func (manager *DbManager) GetRolesByUser(userId interface{}) []*base.Role {
	data, err := manager.mapper.FindRolesByUser(userId)
	if err != nil {
		return nil
	}
	var roles []*base.Role
	for _, authItem := range data {
		role := db.ToRole(*authItem)
		roles = append(roles, &role)
	}
	return roles
}

// GetChildRoles 获取角色关联的子角色列表
func (manager *DbManager) GetChildRoles(roleName string) []*base.Role {
	role := manager.GetRole(roleName)
	if role == nil {
		log.Infof("[rbac] Role %s not found manager.", roleName)
		return nil
	}

	result := make(map[string]bool)
	manager.getChildrenRecursive(roleName, manager.getChildrenList(), result)
	roles := make([]*base.Role, 1)
	roles = append(roles, role)

	for _, r := range manager.GetRoles() {
		if result[r.Name] == true {
			roles = append(roles, r)
		}
	}

	return roles
}

func (manager *DbManager) getChildrenRecursive(name string, childrenList map[string][]string, result map[string]bool) {
	if childrenList[name] != nil {
		// 存在子元素
		for _, childName := range childrenList[name] {
			result[childName] = true
			manager.getChildrenRecursive(childName, childrenList, result)
		}
	}
}

func (manager *DbManager) getChildrenList() map[string][]string {
	m := make(map[string][]string)
	list, err := manager.mapper.FindChildrenList()
	if err != nil {
		log.Warnf("[rbac] get children list err=%v", err)
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

func (manager *DbManager) GetPermissionsByRole(roleName string) []*base.Permission {
	childrenList := manager.getChildrenList()
	result := make(map[string]bool)
	manager.getChildrenRecursive(roleName, childrenList, result)

	permissions := make([]*base.Permission, 0)
	names := make([]string, 0)
	for name, exists := range result {
		if exists == true {
			names = append(names, name)
		}
	}

	if len(names) == 0 {
		return permissions
	}

	if list, err := manager.mapper.GetItemList(base.PermissionType.Value(), names); err == nil {
		for _, item := range list {
			permission := base.NewPermission(item.Name, item.Description, item.RuleName, item.ExecuteName, item.CreateTime, item.UpdateTime)
			permissions = append(permissions, permission)
		}
	}

	return permissions
}

func (manager *DbManager) GetPermissionsByUser(userId interface{}) []*base.Permission {
	directPermissions := manager.getDirectPermissionsByUser(userId)
	inheritedPermissions := manager.getInheritedPermissionsByUser(userId)
	return base.MergePermissions(directPermissions, inheritedPermissions)
}

// 直接关联的权限列表
func (manager *DbManager) getDirectPermissionsByUser(userId interface{}) []*base.Permission {
	permissions := make([]*base.Permission, 0)
	if data, err := manager.mapper.FindPermissionsByUser(userId); err == nil {
		for _, authItem := range data {
			permission := base.NewPermission(authItem.Name, authItem.Description, authItem.RuleName, authItem.ExecuteName, authItem.CreateTime, authItem.UpdateTime)
			permissions = append(permissions, permission)
		}
	}
	return permissions
}

func (manager *DbManager) getInheritedPermissionsByUser(userId interface{}) []*base.Permission {
	permissions := make([]*base.Permission, 0)
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

		if list, err := manager.mapper.GetItemList(base.PermissionType.Value(), names); err == nil {
			for _, item := range list {
				permission := base.NewPermission(item.Name, item.Description, item.RuleName, item.ExecuteName, item.CreateTime, item.UpdateTime)
				permissions = append(permissions, permission)
			}
		}

	}

	return permissions
}

func (manager *DbManager) CanAddChild(parent base.Item, child base.Item) bool {
	return !manager.detectLoop(parent, child)
}

// 递归遍历是否存在子元素是父元素本身，避免出现环
func (manager *DbManager) detectLoop(parent base.Item, child base.Item) bool {
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

func (manager *DbManager) AddChild(parent base.Item, child base.Item) error {

	if parent.GetName() == child.GetName() {
		return errors.New(fmt.Sprintf("Cannot add '%s' as a child of itself.", parent.GetName()))
	}

	if parent.GetType() == base.PermissionType && child.GetType() == base.RoleType {
		return errors.New(fmt.Sprintf("Cannot add a role as a child of a permission."))
	}

	if manager.detectLoop(parent, child) {
		return errors.New(fmt.Sprintf("Cannot add '%s' as a child of '%s'. A loop has been detected.", parent.GetName(), child.GetName()))
	}

	itemChild := db.ToAuthItemChild(parent.GetName(), child.GetName())
	return manager.mapper.AddItemChild(itemChild)
}

func (manager *DbManager) RemoveChild(parent base.Item, child base.Item) bool {
	err := manager.mapper.RemoveChild(parent.GetName(), child.GetName())
	return err == nil
}

func (manager *DbManager) RemoveChildren(parent base.Item) bool {
	err := manager.mapper.RemoveChildren(parent.GetName())
	return err == nil
}

func (manager *DbManager) HasChild(parent base.Item, child base.Item) bool {
	if parent == nil || child == nil {
		return false
	}
	return manager.mapper.HasChild(parent.GetName(), child.GetName())
}

func (manager *DbManager) GetChildren(name string) []base.Item {
	data, err := manager.mapper.FindChildren(name)
	if err != nil {
		return nil
	}

	var items []base.Item
	for _, authItem := range data {
		item := db.ToItem(*authItem)
		items = append(items, item)
	}
	return items
}

func (manager *DbManager) Assign(item base.Item, userId interface{}) *base.Assignment {
	assignment := base.NewAssignment(userId, item.GetName())
	authAssignment := db.ToAuthAssignment(assignment)
	err := manager.mapper.Assign(authAssignment)
	if err == nil {
		delete(manager._checkAccessAssignments, userId)
		return &assignment
	}
	return nil
}

func (manager *DbManager) Revoke(item base.Item, userId interface{}) bool {
	delete(manager._checkAccessAssignments, userId)
	err := manager.mapper.RemoveAssignment(userId, item.GetName())
	return err == nil
}

func (manager *DbManager) RevokeAll(userId interface{}) bool {
	delete(manager._checkAccessAssignments, userId)
	err := manager.mapper.RemoveAllAssignmentByUser(userId)
	return err == nil
}

func (manager *DbManager) GetAssignment(roleName string, userId interface{}) *base.Assignment {
	authAssignment, err := manager.mapper.GetAssignment(userId, roleName)
	if err != nil {
		return nil
	}
	assignment := db.ToAssignment(*authAssignment)
	return &assignment
}

func (manager *DbManager) GetAssignments(userId interface{}) map[string]*base.Assignment {
	assignments := make(map[string]*base.Assignment)
	authAssignments, err := manager.mapper.GetAssignments(userId)
	if err != nil {
		return assignments
	}

	for _, authAssignment := range authAssignments {
		assignment := db.ToAssignment(*authAssignment)
		assignments[assignment.ItemName] = &assignment
	}
	return assignments
}

func (manager *DbManager) GetUserIdsByRole(roleName string) []interface{} {
	users := make([]interface{}, 0)
	authAssignments, err := manager.mapper.GetAssignmentByItems(roleName)
	if err != nil {
		log.Warnf("[rbac] GetUserIdsByRole err = %v", err)
		return users
	}

	for _, authAssignment := range authAssignments {
		users = append(users, authAssignment.UserId)
	}

	return users
}

func (manager *DbManager) RemoveAll() {
	_ = manager.mapper.RemoveAll()
	manager.invalidateCache()
}

func (manager *DbManager) RemoveAllPermissions() {
	manager.removeAllItems(base.PermissionType)
}

func (manager *DbManager) RemoveAllRoles() {
	manager.removeAllItems(base.RoleType)
}

func (manager *DbManager) removeAllItems(t base.ItemType) {
	items, err := manager.mapper.GetItems(t.Value())
	if err != nil {
		log.Warnf("[rbac] RemoveAllItems err = %v", err)
		return
	}

	if len(items) == 0 {
		return
	}

	names := make([]string, 0)
	for _, item := range items {
		names = append(names, item.Name)
	}

	key := "parent"
	if t == base.PermissionType {
		key = "child"
	}

	_ = manager.mapper.RemoveChildByNames(key, names)
	_ = manager.mapper.RemoveAssignmentByName(names)
	_ = manager.mapper.RemoveItemByType(t.Value())
	manager.invalidateCache()
}

func (manager *DbManager) RemoveAllRules() {
	_ = manager.mapper.RemoveAllRules()
	manager.invalidateCache()
}

func (manager *DbManager) RemoveAllAssignments() {
	manager._checkAccessAssignments = make(map[interface{}]map[string]*base.Assignment)
	_ = manager.mapper.RemoveAllAssignments()
}

func (manager *DbManager) CheckAccess(userId interface{}, permissionName string, params map[string]interface{}) bool {
	var assignments map[string]*base.Assignment
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

func (manager *DbManager) loadFromCache() {
	if manager.items == nil || !manager.cache {
		log.Info("[rbac] load from cache fail!!")
		return
	}

	manager.invalidateCache()

	rules, err2 := manager.mapper.GetRules()
	if err2 == nil {
		for _, rule := range rules {
			manager.rules[rule.Name] = base.NewRule(rule.Name, rule.ExecuteName, rule.CreateTime, rule.UpdateTime)
		}
	}

	authItems, err := manager.mapper.FindAllItems()
	if err != nil {
		log.Warnf("[rbac] LoadFromCache [findAllItems err] = %v", err)
		return
	}

	for _, authItem := range authItems {
		item := db.ToItem(*authItem)
		manager.items[item.GetName()] = item
	}

	authItemChildren, err := manager.mapper.FindChildrenList()
	if err != nil {
		log.Warnf("[rbac] LoadFromCache [FindChildrenList err] = %v", err)
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

func (manager *DbManager) checkAccessFromCache(userId interface{}, itemName string, params map[string]interface{}, assignments map[string]*base.Assignment) bool {
	if manager.items[itemName] == nil {
		return false
	}

	item := manager.items[itemName]
	// log.debug(item instanceof Role ? "Checking role: " + itemName : "Checking permission: " + itemName);
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

func (manager *DbManager) checkAccessRecursive(userId interface{}, itemName string, params map[string]interface{}, assignments map[string]*base.Assignment) bool {
	if manager.items[itemName] == nil {
		return false
	}

	item := manager.items[itemName]
	// log.debug(item instanceof Role ? "Checking role: " + itemName : "Checking permission: " + itemName);
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

func (manager *DbManager) CreateRole(name string) *base.Role {
	return base.NewRole(name, "", "", "", time.Now(), time.Now())
}

func (manager *DbManager) CreatePermission(name string) *base.Permission {
	return base.NewPermission(name, "", "", "", time.Now(), time.Now())
}

func (manager *DbManager) Add(item base.Item) bool {
	// TODO if the rule of the object is not alive in the system, then to create it to the system
	manager.checkRuleExits(item.GetRuleName())
	return manager.addItem(item)
}

func (manager *DbManager) Remove(item base.Item) bool {
	return manager.removeItem(item)
}

func (manager *DbManager) RemoveAllAssignmentByUser(userId interface{}) error {
	err := manager.mapper.RemoveAllAssignmentByUser(userId)
	manager.invalidateCache()
	return err
}

func (manager *DbManager) Update(name string, item base.Item) bool {
	// TODO if the rule of the object is not alive in the system, then to create it to the system
	manager.checkRuleExits(item.GetRuleName())
	return manager.updateItem(name, item)
}

func (manager *DbManager) GetRole(name string) *base.Role {
	item := manager.getItem(name)
	if item == nil {
		return nil
	}
	return base.NewRole(item.GetName(), item.GetDescription(), item.GetRuleName(), item.GetExecuteName(), item.GetCreateTime(), item.GetUpdateTime())
}

func (manager *DbManager) GetRoles() []*base.Role {
	var roles []*base.Role
	items := manager.getItems(base.RoleType)
	for _, item := range items {
		role := base.NewRole(item.GetName(), item.GetDescription(), item.GetRuleName(), item.GetExecuteName(), item.GetCreateTime(), item.GetUpdateTime())
		roles = append(roles, role)
	}
	return roles
}

func (manager *DbManager) GetPermission(name string) *base.Permission {
	item := manager.getItem(name)
	if item == nil {
		return nil
	}
	return base.NewPermission(item.GetName(), item.GetDescription(), item.GetRuleName(), item.GetExecuteName(), item.GetCreateTime(), item.GetUpdateTime())
}

func (manager *DbManager) GetPermissions() []*base.Permission {
	var permissions []*base.Permission
	items := manager.getItems(base.PermissionType)
	for _, item := range items {
		permission := base.NewPermission(item.GetName(), item.GetDescription(), item.GetRuleName(), item.GetExecuteName(), item.GetCreateTime(), item.GetUpdateTime())
		permissions = append(permissions, permission)
	}
	return permissions
}

func (manager *DbManager) checkRuleExits(name string) {
	if name != "" && manager.GetRule(name) == nil {
		rule := base.NewRule(name, "", time.Now(), time.Now())
		manager.AddRule(*rule)
	}
}

func (manager *DbManager) SetDefaultRoles(roles ...*base.Role) {
	for _, role := range roles {
		manager.defaultRoles[role.GetName()] = role
	}
}

func (manager *DbManager) getDefaultRoles() map[string]*base.Role {
	return manager.defaultRoles
}

func (manager *DbManager) executeRule(userId interface{}, item base.Item, params map[string]interface{}) bool {
	if item.GetRuleName() == "" {
		return true
	}

	rule := manager.GetRule(item.GetRuleName())
	if rule == nil {
		log.Warn("[rbac] Rule not found: " + item.GetRuleName())
		return false
	}

	return rule.GetExecutor().Execute(userId, item, params)
}

func (manager *DbManager) hasNoAssignments(assignments map[string]*base.Assignment) bool {
	return len(assignments) == 0 && len(manager.defaultRoles) == 0
}
