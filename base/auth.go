package base

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

	// GetRole
	//	 * Returns the named role.
	//	 *
	//	 * @param name string $ the role name.
	//	 * @return null|Role the role corresponding to the specified name. Null is returned if no such role.
	//	 */
	GetRole(name string) *Role

	// GetRoles
	/**
	 * Returns all roles in the system.
	 *
	 * @return Role[] all roles in the system. The array is indexed by the role names.
	 */
	GetRoles() []*Role

	// GetRolesByUser
	/**
	 * Returns the roles that are assigned to the user via [[assign()]].
	 * Note that child roles that are not assigned directly to the user will not be returned.
	 *
	 * @param userId string|int $ the user ID (see [[\yii\web\User::id]])
	 * @return Role[] all roles directly assigned to the user. The array is indexed by the role names.
	 */
	GetRolesByUser(userId interface{}) []*Role

	// GetChildRoles
	/**
	 * Returns child roles of the role specified. Depth isn't limited.
	 *
	 * @param roleName string $ name of the role to file child roles for
	 * @return Role[] Child roles. The array is indexed by the role names.
	 * First base is an instance of the parent Role itself.
	 */
	GetChildRoles(roleName string) []*Role

	// GetPermission
	/**
	 * Returns the named permission.
	 *
	 * @param name string $ the permission name.
	 * @return null|Permission the permission corresponding to the specified name. Null is returned if no such permission.
	 */
	GetPermission(name string) *Permission

	// GetPermissions
	/**
	 * Returns all permissions in the system.
	 *
	 * @return Permission[] all permissions in the system. The array is indexed by the permission names.
	 */
	GetPermissions() []*Permission

	// GetPermissionsByRole
	/**
	 * Returns all permissions that the specified role represents.
	 *
	 * @param roleName string $ the role name
	 * @return Permission[] all permissions that the role represents. The array is indexed by the permission names.
	 */
	GetPermissionsByRole(roleName string) []*Permission

	// GetPermissionsByUser
	/**
	 * Returns all permissions that the user has.
	 *
	 * @param userId string|int $ the user ID (see [[\yii\web\User::id]])
	 * @return Permission[] all permissions that the user has. The array is indexed by the permission names.
	 */
	GetPermissionsByUser(userId interface{}) []*Permission

	// GetRule
	/**
	 * Returns the rule of the specified name.
	 *
	 * @param name string $ the rule name
	 * @return null|Rule the rule object, or null if the specified name does not correspond to a rule.
	 */
	GetRule(name string) *Rule

	// GetRules
	/**
	 * Returns all rules available in the system.
	 *
	 * @return Rule[] the rules indexed by the rule names
	 */
	GetRules() []*Rule

	// CanAddChild
	/**
	 * Checks the possibility of adding a child to parent.
	 *
	 * @param parent Item $ the parent item
	 * @param child  Item $ the child item to be added to the hierarchy
	 * @return bool possibility of adding
	 * @since 2.0.8
	 */
	CanAddChild(parent Item, child Item) bool

	// AddChild
	/**
	 * Adds an item as a child of another item.
	 *
	 * @param parent Item $
	 * @param child  Item $
	 * @return bool whether the child successfully added
	 */
	AddChild(parent Item, child Item) error

	// RemoveChild
	/**
	 * Removes a child from its parent.
	 * Note, the child item is not deleted. Only the parent-child relationship is removed.
	 *
	 * @param parent Item $
	 * @param child  Item $
	 * @return bool whether the removal is successful
	 */
	RemoveChild(parent Item, child Item) bool

	// RemoveChildren
	/**
	 * Removed all children form their parent.
	 * Note, the children items are not deleted. Only the parent-child relationships are removed.
	 *
	 * @param parent Item $
	 * @return bool whether the removal is successful
	 */
	RemoveChildren(parent Item) bool

	// HasChild
	/**
	 * Returns a value indicating whether the child already exists for the parent.
	 *
	 * @param parent Item $
	 * @param child  Item $
	 * @return bool whether `$child` is already a child of `$parent`
	 */
	HasChild(parent Item, child Item) bool

	// GetChildren
	/**
	 * Returns the child permissions and/or roles.
	 *
	 * @param name string $ the parent name
	 * @return Item[] the child permissions and/or roles
	 */
	GetChildren(name string) []Item

	// Assign
	/**
	 * Assigns a role to a user.
	 *
	 * @param item   Role|Permission $
	 * @param userId string|int $ the user ID (see [[\yii\web\User::id]])
	 * @return Assignment the role assignment information.
	 */
	Assign(item Item, userId interface{}) *Assignment

	// Revoke
	/**
	 * Revokes a role from a user.
	 *
	 * @param item   Role|Permission $
	 * @param userId string|int $ the user ID (see [[\yii\web\User::id]])
	 * @return bool whether the revoking is successful
	 */
	Revoke(item Item, userId interface{}) bool

	// RevokeAll
	/**
	 * Revokes all roles from a user.
	 *
	 * @param userId mixed $ the user ID (see [[\yii\web\User::id]])
	 * @return bool whether the revoking is successful
	 */
	RevokeAll(userId interface{}) bool

	// GetAssignment
	/**
	 * Returns the assignment information regarding a role and a user.
	 *
	 * @param roleName string $ the role name
	 * @param userId   string|int $ the user ID (see [[\yii\web\User::id]])
	 * @return null|Assignment the assignment information. Null is returned if
	 * the role is not assigned to the user.
	 */
	GetAssignment(roleName string, userId interface{}) *Assignment

	// GetAssignments
	/**
	 * Returns all role assignment information for the specified user.
	 *
	 * @param userId string|int $ the user ID (see [[\yii\web\User::id]])
	 * @return Assignment[] the assignments indexed by role names. An empty array will be
	 * returned if there is no role assigned to the user.
	 */
	GetAssignments(userId interface{}) map[string]*Assignment

	// GetUserIdsByRole
	/**
	 * Returns all user IDs assigned to the role specified.
	 *
	 * @param roleName string $
	 * @return array array of user ID strings
	 * @since 2.0.7
	 */
	GetUserIdsByRole(roleName string) []interface{}

	// RemoveAll
	/**
	 * Removes all authorization data, including roles, permissions, rules, and assignments.
	 */
	RemoveAll()

	// RemoveAllPermissions
	/**
	 * Removes all permissions.
	 * All parent child relations will be adjusted accordingly.
	 */
	RemoveAllPermissions()

	// RemoveAllRoles
	/**
	 * Removes all roles.
	 * All parent child relations will be adjusted accordingly.
	 */
	RemoveAllRoles()

	// RemoveAllRules
	/**
	 * Removes all rules.
	 * All roles and permissions which have rules will be adjusted accordingly.
	 */
	RemoveAllRules()

	// RemoveAllAssignments
	/**
	 * Removes all role assignments.
	 */
	RemoveAllAssignments()

	// RemoveAllAssignmentByUser
	/**
	 * Remove all role assignments by user
	 */
	RemoveAllAssignmentByUser(userId interface{}) error
}

type ManagerInterface interface {
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
