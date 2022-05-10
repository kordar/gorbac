package gorbac

import (
	"errors"
	"fmt"
	"log"
	"time"
)

type CheckAccess interface {
	CheckAccess(userId int64, permissionName string, params map[string]interface{}) bool
}

// AuthManager /**
type AuthManager interface {
	// CreatePermission /**
	/**
	 * Creates a new Role object.
	 * Note that the newly created role is not added to the RBAC system yet.
	 * You must fill in the needed data and call [[add()]] to add it to the system.
	 *
	 * @param name String the role name
	 * @return Role the new Role object
	 */
	CreatePermission(name string) *Permission

	// Add /**
	/**
	 * Adds a role, permission or rule to the RBAC system.
	 *
	 * @param object Permission|Rule $object
	 * @return bool whether the role, permission or rule is successfully added to the system
	 */
	Add(item Item) bool
	AddRule(rule Rule) bool

	// Remove /**
	/**
	 * Removes a role, permission or rule from the RBAC system.
	 *
	 * @param object Role|Permission|Rule
	 * @return bool whether the role, permission or rule is successfully removed
	 */
	Remove(item Item) bool
	RemoveRule(rule Rule) bool

	// Update /*
	/**
	 * Updates the specified role, permission or rule in the system.
	 *
	 * @param name   string $ the old name of the role, permission or rule
	 * @param object Role|Permission|Rule $
	 * @return bool whether the update is successful
	 */
	Update(name string, item Item) bool
	UpdateRule(name string, rule Rule) bool

	// GetRole /**
	//	 * Returns the named role.
	//	 *
	//	 * @param name string $ the role name.
	//	 * @return null|Role the role corresponding to the specified name. Null is returned if no such role.
	//	 */
	GetRole(name string) *Role

	/**
	 * Returns all roles in the system.
	 *
	 * @return Role[] all roles in the system. The array is indexed by the role names.
	 */
	GetRoles() []*Role

	/**
	 * Returns the roles that are assigned to the user via [[assign()]].
	 * Note that child roles that are not assigned directly to the user will not be returned.
	 *
	 * @param userId string|int $ the user ID (see [[\yii\web\User::id]])
	 * @return Role[] all roles directly assigned to the user. The array is indexed by the role names.
	 */
	GetRolesByUser(userId int64) []*Role

	/**
	 * Returns child roles of the role specified. Depth isn't limited.
	 *
	 * @param roleName string $ name of the role to file child roles for
	 * @return Role[] Child roles. The array is indexed by the role names.
	 * First element is an instance of the parent Role itself.
	 */
	GetChildRoles(roleName string) []*Role

	/**
	 * Returns the named permission.
	 *
	 * @param name string $ the permission name.
	 * @return null|Permission the permission corresponding to the specified name. Null is returned if no such permission.
	 */
	GetPermission(name string) *Permission

	/**
	 * Returns all permissions in the system.
	 *
	 * @return Permission[] all permissions in the system. The array is indexed by the permission names.
	 */
	GetPermissions() []*Permission

	/**
	 * Returns all permissions that the specified role represents.
	 *
	 * @param roleName string $ the role name
	 * @return Permission[] all permissions that the role represents. The array is indexed by the permission names.
	 */
	GetPermissionsByRole(roleName string) []*Permission

	/**
	 * Returns all permissions that the user has.
	 *
	 * @param userId string|int $ the user ID (see [[\yii\web\User::id]])
	 * @return Permission[] all permissions that the user has. The array is indexed by the permission names.
	 */
	GetPermissionsByUser(userId int64) []*Permission

	/**
	 * Returns the rule of the specified name.
	 *
	 * @param name string $ the rule name
	 * @return null|Rule the rule object, or null if the specified name does not correspond to a rule.
	 */
	GetRule(name string) *Rule

	/**
	 * Returns all rules available in the system.
	 *
	 * @return Rule[] the rules indexed by the rule names
	 */
	GetRules() []*Rule

	/**
	 * Checks the possibility of adding a child to parent.
	 *
	 * @param parent Item $ the parent item
	 * @param child  Item $ the child item to be added to the hierarchy
	 * @return bool possibility of adding
	 * @since 2.0.8
	 */
	CanAddChild(parent Item, child Item) bool

	/**
	 * Adds an item as a child of another item.
	 *
	 * @param parent Item $
	 * @param child  Item $
	 * @return bool whether the child successfully added
	 */
	AddChild(parent Item, child Item) error

	/**
	 * Removes a child from its parent.
	 * Note, the child item is not deleted. Only the parent-child relationship is removed.
	 *
	 * @param parent Item $
	 * @param child  Item $
	 * @return bool whether the removal is successful
	 */
	RemoveChild(parent Item, child Item) bool

	/**
	 * Removed all children form their parent.
	 * Note, the children items are not deleted. Only the parent-child relationships are removed.
	 *
	 * @param parent Item $
	 * @return bool whether the removal is successful
	 */
	RemoveChildren(parent Item) bool

	/**
	 * Returns a value indicating whether the child already exists for the parent.
	 *
	 * @param parent Item $
	 * @param child  Item $
	 * @return bool whether `$child` is already a child of `$parent`
	 */

	HasChild(parent Item, child Item) bool

	/**
	 * Returns the child permissions and/or roles.
	 *
	 * @param name string $ the parent name
	 * @return Item[] the child permissions and/or roles
	 */
	GetChildren(name string) []Item

	/**
	 * Assigns a role to a user.
	 *
	 * @param item   Role|Permission $
	 * @param userId string|int $ the user ID (see [[\yii\web\User::id]])
	 * @return Assignment the role assignment information.
	 */
	Assign(item Item, userId int64) *Assignment

	/**
	 * Revokes a role from a user.
	 *
	 * @param item   Role|Permission $
	 * @param userId string|int $ the user ID (see [[\yii\web\User::id]])
	 * @return bool whether the revoking is successful
	 */
	Revoke(item Item, userId int64) bool

	/**
	 * Revokes all roles from a user.
	 *
	 * @param userId mixed $ the user ID (see [[\yii\web\User::id]])
	 * @return bool whether the revoking is successful
	 */
	RevokeAll(userId int64) bool

	/**
	 * Returns the assignment information regarding a role and a user.
	 *
	 * @param roleName string $ the role name
	 * @param userId   string|int $ the user ID (see [[\yii\web\User::id]])
	 * @return null|Assignment the assignment information. Null is returned if
	 * the role is not assigned to the user.
	 */
	GetAssignment(roleName string, userId int64) *Assignment

	/**
	 * Returns all role assignment information for the specified user.
	 *
	 * @param userId string|int $ the user ID (see [[\yii\web\User::id]])
	 * @return Assignment[] the assignments indexed by role names. An empty array will be
	 * returned if there is no role assigned to the user.
	 */
	GetAssignments(userId int64) map[string]*Assignment

	/**
	 * Returns all user IDs assigned to the role specified.
	 *
	 * @param roleName string $
	 * @return array array of user ID strings
	 * @since 2.0.7
	 */
	GetUserIdsByRole(roleName string) []int64

	/**
	 * Removes all authorization data, including roles, permissions, rules, and assignments.
	 */
	RemoveAll()

	/**
	 * Removes all permissions.
	 * All parent child relations will be adjusted accordingly.
	 */
	RemoveAllPermissions()

	/**
	 * Removes all roles.
	 * All parent child relations will be adjusted accordingly.
	 */
	RemoveAllRoles()

	/**
	 * Removes all rules.
	 * All roles and permissions which have rules will be adjusted accordingly.
	 */
	RemoveAllRules()

	/**
	 * Removes all role assignments.
	 */
	RemoveAllAssignments()

	/**
	 * Remove all role assignments by user
	 */
	RemoveAllAssignmentByUser(userId int64) error
}

/**
base interface
*/
type BaseManagerInterface interface {
	/**
	 * Returns the named auth item.
	 *
	 * @param name String the auth item name.
	 * @return Item the auth item corresponding to the specified name. Null is returned if no such item.
	 */
	getItem(name string) Item
	/**
	 * Returns the items of the specified type.
	 *
	 * @param type int the auth item type (either [[Item::TYPE_ROLE]] or [[Item::TYPE_PERMISSION]]
	 * @return Item[] the auth items of the specified type.
	 */
	getItems(t ItemType) []Item
	/**
	 * Adds an auth item to the RBAC system.
	 *
	 * @param item the item to add
	 * @return bool whether the auth item is successfully added to the system.
	 */
	addItem(item Item) bool
	/**
	 * Remove an auth item from the RBAC system.
	 *
	 * @param item the item to remove
	 * @return bool whether the role or permission is successfully removed.
	 */
	removeItem(item Item) bool
	/**
	 * Updates an auth item in the RBAC system
	 *
	 * @param name String name the name of the item being updated
	 * @param item the updated item
	 * @return bool whether the auth item is successfully updated.
	 */
	updateItem(name string, item Item) bool
}

type DbManager struct {
	/**
	 */
	mapper AuthRepository
	/**

	 */
	cache bool
	/**
	 * Item[] all auth items (name => Item)
	 */
	items map[string]Item
	/**
	 * Rule[] all auth rules (name => Rule)
	 */
	rules map[string]*Rule
	/**
	 * @var array a list of role names that are assigned to every user automatically without calling [[assign()]].
	 * Note that these roles are applied to users, regardless of their state of authentication.
	 */
	defaultRoles map[string]*Role
	/**
	 * array auth item parent-child relationships (childName => list of parents)
	 */
	parents map[string][]string

	_checkAccessAssignments map[int64]map[string]*Assignment
}

func NewDbManager(mapper AuthRepository, cache bool) *DbManager {
	return &DbManager{mapper: mapper, cache: cache, _checkAccessAssignments: make(map[int64]map[string]*Assignment), defaultRoles: make(map[string]*Role)}
}

func (manager *DbManager) invalidateCache() {
	if manager.cache {
		manager.items = make(map[string]Item)
		manager.rules = make(map[string]*Rule)
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

func (manager *DbManager) getItem(name string) Item {
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
	item := ConversionToItem(*authItem)
	return item
}

func (manager *DbManager) getItems(t ItemType) []Item {
	data, err := manager.mapper.GetItems(t.Value())
	if err != nil {
		return nil
	}

	var items []Item
	for _, authItem := range data {
		item := ConversionToItem(*authItem)
		items = append(items, item)
	}
	return items
}

func (manager *DbManager) GetRule(name string) *Rule {
	if manager.rules != nil {
		return manager.rules[name]
	}

	authRule, err := manager.mapper.GetRule(name)
	if err != nil {
		return nil
	}
	rule := ConversionToRule(*authRule)
	return &rule
}

func (manager *DbManager) GetRules() []*Rule {
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
	for _, authRule := range data {
		rule := ConversionToRule(*authRule)
		rules = append(rules, &rule)
	}
	return rules
}

func (manager *DbManager) addItem(item Item) bool {
	authItem := ConversionToAuthItem(item)
	err := manager.mapper.AddItem(authItem)
	return manager.refreshInvalidateCache(err == nil)
}

func (manager *DbManager) AddRule(rule Rule) bool {
	authRule := ConversionToAuthRule(rule)
	err := manager.mapper.AddRule(authRule)
	return manager.refreshInvalidateCache(err == nil)
}

func (manager *DbManager) removeItem(item Item) bool {
	_ = manager.mapper.RemoveItem(item.GetName())
	manager.invalidateCache()
	return true
}

func (manager *DbManager) RemoveRule(rule Rule) bool {
	_ = manager.mapper.RemoveRule(rule.Name)
	manager.invalidateCache()
	return true
}

func (manager *DbManager) updateItem(name string, item Item) bool {
	authItem := ConversionToAuthItem(item)
	err := manager.mapper.UpdateItem(name, authItem)
	return manager.refreshInvalidateCache(err == nil)
}

func (manager *DbManager) UpdateRule(name string, rule Rule) bool {
	authRule := ConversionToAuthRule(rule)
	err := manager.mapper.UpdateRule(name, authRule)
	return manager.refreshInvalidateCache(err == nil)
}

/**
获取用户角色列表
*/
func (manager *DbManager) GetRolesByUser(userId int64) []*Role {
	if userId < 0 {
		return nil
	}
	data, err := manager.mapper.FindRolesByUser(userId)
	if err != nil {
		return nil
	}
	var roles []*Role
	for _, authItem := range data {
		role := ConversionToRole(*authItem)
		roles = append(roles, &role)
	}
	return roles
}

/**
获取角色关联的子角色列表
*/
func (manager *DbManager) GetChildRoles(roleName string) []*Role {
	role := manager.GetRole(roleName)
	if role == nil {
		log.Println(fmt.Sprintf("Role %s not founmanager.", roleName))
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
		log.Println(fmt.Sprintf("getChildrenList err=%v", err))
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

func (manager *DbManager) GetPermissionsByRole(roleName string) []*Permission {
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
			permission := NewPermission(item.Name, item.Description, item.RuleName, item.ExecuteName, item.CreatedAt, item.UpdatedAt)
			permissions = append(permissions, permission)
		}
	}

	return permissions
}

func (manager *DbManager) GetPermissionsByUser(userId int64) []*Permission {
	directPermissions := manager.getDirectPermissionsByUser(userId)
	inheritedPermissions := manager.getInheritedPermissionsByUser(userId)
	return MergePermissions(directPermissions, inheritedPermissions)
}

/**
直接关联的权限列表
*/
func (manager *DbManager) getDirectPermissionsByUser(userId int64) []*Permission {
	permissions := make([]*Permission, 0)
	if userId <= 0 {
		return permissions
	}
	if data, err := manager.mapper.FindPermissionsByUser(userId); err == nil {
		for _, authItem := range data {
			permission := NewPermission(authItem.Name, authItem.Description, authItem.RuleName, authItem.ExecuteName, authItem.CreatedAt, authItem.UpdatedAt)
			permissions = append(permissions, permission)
		}
	}
	return permissions
}

func (manager *DbManager) getInheritedPermissionsByUser(userId int64) []*Permission {
	permissions := make([]*Permission, 0)
	if userId <= 0 {
		return permissions
	}

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

		if list, err := manager.mapper.GetItemList(PermissionType.Value(), names); err == nil {
			for _, item := range list {
				permission := NewPermission(item.Name, item.Description, item.RuleName, item.ExecuteName, item.CreatedAt, item.UpdatedAt)
				permissions = append(permissions, permission)
			}
		}

	}

	return permissions
}

func (manager *DbManager) CanAddChild(parent Item, child Item) bool {
	return !manager.detectLoop(parent, child)
}

/**
递归遍历是否存在子元素是父元素本身，避免出现环
*/
func (manager *DbManager) detectLoop(parent Item, child Item) bool {
	if child.GetName() == parent.GetName() {
		return true
	}

	children := manager.GetChildren(child.GetName())
	if children != nil {
		for _, child := range children {
			if manager.detectLoop(parent, child) {
				return true
			}
		}
	}

	return false
}

func (manager *DbManager) AddChild(parent Item, child Item) error {

	if parent.GetName() == child.GetName() {
		return errors.New(fmt.Sprintf("Cannot add '%s' as a child of itself.", parent.GetName()))
	}

	if parent.GetType() == PermissionType && child.GetType() == RoleType {
		return errors.New(fmt.Sprintf("Cannot add a role as a child of a permission."))
	}

	if manager.detectLoop(parent, child) {
		return errors.New(fmt.Sprintf("Cannot add '%s' as a child of '%s'. A loop has been detected.", parent.GetName(), child.GetName()))
	}

	itemChild := ConversionToAuthItemChild(parent.GetName(), child.GetName())
	return manager.mapper.AddItemChild(itemChild)
}

func (manager *DbManager) RemoveChild(parent Item, child Item) bool {
	err := manager.mapper.RemoveChild(parent.GetName(), child.GetName())
	return err == nil
}

func (manager *DbManager) RemoveChildren(parent Item) bool {
	err := manager.mapper.RemoveChildren(parent.GetName())
	return err == nil
}

func (manager *DbManager) HasChild(parent Item, child Item) bool {
	if parent == nil || child == nil {
		return false
	}
	return manager.mapper.HasChild(parent.GetName(), child.GetName())
}

func (manager *DbManager) GetChildren(name string) []Item {
	data, err := manager.mapper.FindChildren(name)
	if err != nil {
		return nil
	}

	var items []Item
	for _, authItem := range data {
		item := ConversionToItem(*authItem)
		items = append(items, item)
	}
	return items
}

func (manager *DbManager) Assign(item Item, userId int64) *Assignment {
	assignment := NewAssignment(userId, item.GetName())
	authAssignment := ConversionToAuthAssignment(assignment)
	err := manager.mapper.Assign(authAssignment)
	if err == nil {
		delete(manager._checkAccessAssignments, userId)
		return &assignment
	}
	return nil
}

func (manager *DbManager) Revoke(item Item, userId int64) bool {
	if userId <= 0 {
		return false
	}
	delete(manager._checkAccessAssignments, userId)
	err := manager.mapper.RemoveAssignment(userId, item.GetName())
	return err == nil
}

func (manager *DbManager) RevokeAll(userId int64) bool {
	if userId <= 0 {
		return false
	}
	delete(manager._checkAccessAssignments, userId)
	err := manager.mapper.RemoveAllAssignmentByUser(userId)
	return err == nil
}

func (manager *DbManager) GetAssignment(roleName string, userId int64) *Assignment {
	if userId <= 0 {
		return nil
	}

	authAssignment, err := manager.mapper.GetAssignment(userId, roleName)
	if err != nil {
		return nil
	}
	assignment := ConversionToAssignment(*authAssignment)
	return &assignment
}

func (manager *DbManager) GetAssignments(userId int64) map[string]*Assignment {
	assignments := make(map[string]*Assignment)
	if userId <= 0 {
		return assignments
	}

	authAssignments, err := manager.mapper.GetAssignments(userId)
	if err != nil {
		return assignments
	}

	for _, authAssignment := range authAssignments {
		assignment := ConversionToAssignment(*authAssignment)
		assignments[assignment.ItemName] = &assignment
	}
	return assignments
}

func (manager *DbManager) GetUserIdsByRole(roleName string) []int64 {
	users := make([]int64, 0)
	authAssignments, err := manager.mapper.GetAssignmentByItems(roleName)
	if err != nil {
		log.Println(fmt.Sprintf("GetUserIdsByRole err = %v", err))
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
	manager.removeAllItems(PermissionType)
}

func (manager *DbManager) RemoveAllRoles() {
	manager.removeAllItems(RoleType)
}

func (manager *DbManager) removeAllItems(t ItemType) {
	items, err := manager.mapper.GetItems(t.Value())
	if err != nil {
		log.Println(fmt.Sprintf("removeAllItems err = %v", err))
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
	if t == PermissionType {
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
	manager._checkAccessAssignments = make(map[int64]map[string]*Assignment)
	_ = manager.mapper.RemoveAllAssignments()
}

func (manager *DbManager) CheckAccess(userId int64, permissionName string, params map[string]interface{}) bool {
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

func (manager *DbManager) loadFromCache() {
	if manager.items == nil || !manager.cache {
		log.Println("load from cache fail!!")
		return
	}

	manager.invalidateCache()

	rules, err2 := manager.mapper.GetRules()
	if err2 == nil {
		for _, rule := range rules {
			manager.rules[rule.Name] = NewRule(rule.Name, rule.ExecuteName, rule.CreatedAt, rule.UpdatedAt)
		}
	}

	authItems, err := manager.mapper.FindAllItems()
	if err != nil {
		log.Println(fmt.Sprintf("loadFromCache [findAllItems err] = %v", err))
		return
	}

	for _, authItem := range authItems {
		item := ConversionToItem(*authItem)
		manager.items[item.GetName()] = item
	}

	authItemChildren, err := manager.mapper.FindChildrenList()
	if err != nil {
		log.Println(fmt.Sprintf("loadFromCache [FindChildrenList err] = %v", err))
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

func (manager *DbManager) checkAccessFromCache(userId int64, itemName string, params map[string]interface{}, assignments map[string]*Assignment) bool {
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

func (manager *DbManager) checkAccessRecursive(userId int64, itemName string, params map[string]interface{}, assignments map[string]*Assignment) bool {
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

func (manager *DbManager) CreateRole(name string) *Role {
	return NewRole(name, "", "", "", time.Now(), time.Now())
}

func (manager *DbManager) CreatePermission(name string) *Permission {
	return NewPermission(name, "", "", "", time.Now(), time.Now())
}

func (manager *DbManager) Add(item Item) bool {
	// TODO if the rule of the object is not alive in the system, then to create it to the system
	manager.checkRuleExits(item.GetRuleName())
	return manager.addItem(item)
}

func (manager *DbManager) Remove(item Item) bool {
	return manager.removeItem(item)
}

func (manager *DbManager) RemoveAllAssignmentByUser(userId int64) error {
	err := manager.mapper.RemoveAllAssignmentByUser(userId)
	manager.invalidateCache()
	return err
}

func (manager *DbManager) Update(name string, item Item) bool {
	// TODO if the rule of the object is not alive in the system, then to create it to the system
	manager.checkRuleExits(item.GetRuleName())
	return manager.updateItem(name, item)
}

func (manager *DbManager) GetRole(name string) *Role {
	item := manager.getItem(name)
	return NewRole(item.GetName(), item.GetDescription(), item.GetRuleName(), item.GetExecuteName(), item.GetCreateTime(), item.GetUpdateTime())
}

func (manager *DbManager) GetRoles() []*Role {
	var roles []*Role
	items := manager.getItems(RoleType)
	for _, item := range items {
		role := NewRole(item.GetName(), item.GetDescription(), item.GetRuleName(), item.GetExecuteName(), item.GetCreateTime(), item.GetUpdateTime())
		roles = append(roles, role)
	}
	return roles
}

func (manager *DbManager) GetPermission(name string) *Permission {
	item := manager.getItem(name)
	if item == nil {
		return nil
	}
	return NewPermission(item.GetName(), item.GetDescription(), item.GetRuleName(), item.GetExecuteName(), item.GetCreateTime(), item.GetUpdateTime())
}

func (manager *DbManager) GetPermissions() []*Permission {
	var permissions []*Permission
	items := manager.getItems(PermissionType)
	for _, item := range items {
		permission := NewPermission(item.GetName(), item.GetDescription(), item.GetRuleName(), item.GetExecuteName(), item.GetCreateTime(), item.GetUpdateTime())
		permissions = append(permissions, permission)
	}
	return permissions
}

func (manager *DbManager) checkRuleExits(name string) {
	if name != "" && manager.GetRule(name) == nil {
		rule := NewRule(name, "", time.Now(), time.Now())
		manager.AddRule(*rule)
	}
}

func (manager *DbManager) SetDefaultRoles(roles ...*Role) {
	for _, role := range roles {
		manager.defaultRoles[role.GetName()] = role
	}
}

func (manager *DbManager) getDefaultRoles() map[string]*Role {
	return manager.defaultRoles
}

func (manager *DbManager) executeRule(userId int64, item Item, params map[string]interface{}) bool {
	if item.GetRuleName() == "" {
		return true
	}

	rule := manager.GetRule(item.GetRuleName())
	if rule == nil {
		log.Println("Rule not found: " + item.GetRuleName())
		return false
	}

	return rule.GetExecutor().Execute(userId, item, params)
}

func (manager *DbManager) hasNoAssignments(assignments map[string]*Assignment) bool {
	return len(assignments) == 0 && len(manager.defaultRoles) == 0
}
