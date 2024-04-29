package db

type AuthRepository interface {
	AddItem(authItem AuthItem) error
	GetItem(name string) (*AuthItem, error)
	GetItems(t int32) ([]*AuthItem, error)
	FindAllItems() ([]*AuthItem, error)
	AddRule(rule AuthRule) error
	GetRule(name string) (*AuthRule, error)
	GetRules() ([]*AuthRule, error)
	RemoveItem(name string) error
	RemoveRule(ruleName string) error
	UpdateItem(itemName string, item AuthItem) error
	UpdateRule(ruleName string, rule AuthRule) error
	FindRolesByUser(userId interface{}) ([]*AuthItem, error)
	FindChildrenList() ([]*AuthItemChild, error)
	FindChildrenFormChild(child string) ([]*AuthItemChild, error)
	GetItemList(t int32, names []string) ([]*AuthItem, error)
	FindPermissionsByUser(userId interface{}) ([]*AuthItem, error)
	FindAssignmentByUser(userId interface{}) ([]*AuthAssignment, error)
	AddItemChild(itemChild AuthItemChild) error
	RemoveChild(parent string, child string) error
	RemoveChildren(parent string) error
	HasChild(parent string, child string) bool
	FindChildren(name string) ([]*AuthItem, error)
	Assign(assignment AuthAssignment) error
	RemoveAssignment(userId interface{}, name string) error
	RemoveAllAssignmentByUser(userId interface{}) error
	RemoveAllAssignments() error
	GetAssignment(userId interface{}, name string) (*AuthAssignment, error)
	GetAssignmentByItems(name string) ([]*AuthAssignment, error)
	GetAssignments(userId interface{}) ([]*AuthAssignment, error)
	GetAllAssignment() ([]*AuthAssignment, error)
	RemoveAll() error
	RemoveChildByNames(key string, names []string) error
	RemoveAssignmentByName(names []string) error
	RemoveItemByType(t int32) error
	RemoveAllRules() error
}
