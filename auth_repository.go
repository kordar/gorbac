package gorbac

import "gorm.io/gorm"

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
	FindRolesByUser(userId int64) ([]*AuthItem, error)
	FindChildrenList() ([]*AuthItemChild, error)
	FindChildrenFormChild(child string) ([]*AuthItemChild, error)
	GetItemList(t int32, names []string) ([]*AuthItem, error)
	FindPermissionsByUser(userId int64) ([]*AuthItem, error)
	FindAssignmentByUser(userId int64) ([]*AuthAssignment, error)
	AddItemChild(itemChild AuthItemChild) error
	RemoveChild(parent string, child string) error
	RemoveChildren(parent string) error
	HasChild(parent string, child string) bool
	FindChildren(name string) ([]*AuthItem, error)
	Assign(assignment AuthAssignment) error
	RemoveAssignment(userId int64, name string) error
	RemoveAllAssignmentByUser(userId int64) error
	RemoveAllAssignments() error
	GetAssignment(userId int64, name string) (*AuthAssignment, error)
	GetAssignmentByItems(name string) ([]*AuthAssignment, error)
	GetAssignments(userId int64) ([]*AuthAssignment, error)
	GetAllAssignment() ([]*AuthAssignment, error)
	RemoveAll() error
	RemoveChildByNames(key string, names []string) error
	RemoveAssignmentByName(names []string) error
	RemoveItemByType(t int32) error
	RemoveAllRules() error
}

/**
 * @Description:
 * @receiver rbac
 * @param authItem
 * @return error
 */
type SqlRbac struct {
	db *gorm.DB
}

func NewSqlRbac(db *gorm.DB) *SqlRbac {
	return &SqlRbac{db: db}
}

func (rbac *SqlRbac) AddItem(authItem AuthItem) error {
	return rbac.db.Create(&authItem).Error
}

func (rbac *SqlRbac) GetItem(name string) (*AuthItem, error) {
	item := AuthItem{}
	err := rbac.db.Where("name = ?", name).First(&item).Error
	return &item, err
}

func (rbac *SqlRbac) GetItems(t int32) ([]*AuthItem, error) {
	var items []*AuthItem
	err := rbac.db.Where("type = ?", t).Find(&items).Error
	return items, err
}

func (rbac *SqlRbac) FindAllItems() ([]*AuthItem, error) {
	var items []*AuthItem
	err := rbac.db.Find(&items).Error
	return items, err
}

func (rbac *SqlRbac) AddRule(rule AuthRule) error {
	return rbac.db.Create(&rule).Error
}

func (rbac *SqlRbac) GetRule(name string) (*AuthRule, error) {
	var rule AuthRule
	err := rbac.db.Where("name = ?", name).First(&rule).Error
	return &rule, err
}

func (rbac *SqlRbac) GetRules() ([]*AuthRule, error) {
	var rules []*AuthRule
	err := rbac.db.Find(&rules).Error
	return rules, err
}

func (rbac *SqlRbac) RemoveItem(name string) error {
	return rbac.db.Transaction(func(tx *gorm.DB) error {
		itemChild := AuthItemChild{}
		tx.Where("parent = ? or child = ?", name, name).Delete(&itemChild)
		assignment := AuthAssignment{}
		tx.Where("item_name = ?", name).Delete(&assignment)
		item := AuthItem{}
		tx.Where("name = ?", name).Delete(&item)
		return nil
	})
}

func (rbac *SqlRbac) RemoveRule(ruleName string) error {
	return rbac.db.Transaction(func(tx *gorm.DB) error {
		item := AuthItem{}
		tx.Model(&item).Where("rule_name = ?", ruleName).Update("rule_name", nil)
		rule := AuthRule{}
		tx.Where("name = ?", ruleName).Delete(&rule)
		return nil
	})
}

func (rbac *SqlRbac) UpdateItem(itemName string, updateItem AuthItem) error {
	return rbac.db.Transaction(func(tx *gorm.DB) error {
		if itemName != updateItem.Name {
			child := AuthItemChild{}
			assignment := AuthAssignment{}
			tx.Model(&child).Where("parent = ?", itemName).Update("parent", updateItem.Name)
			tx.Model(&child).Where("child = ?", itemName).Update("child", updateItem.Name)
			tx.Model(&assignment).Where("item_name = ?", itemName).Update("item_name", updateItem.Name)
		}
		authItem := AuthItem{}
		return tx.Model(&authItem).Where("name = ?", itemName).Omit("create_at").Updates(&updateItem).Error
	})
}

func (rbac *SqlRbac) UpdateRule(ruleName string, updateRule AuthRule) error {
	return rbac.db.Transaction(func(tx *gorm.DB) error {
		if ruleName != updateRule.Name {
			item := AuthItem{}
			tx.Model(&item).Where("rule_name = ?", ruleName).Update("rule_name", updateRule.Name)
		}
		rule := AuthRule{}
		return tx.Model(&rule).Where("name = ?", ruleName).Omit("create_at").Updates(&updateRule).Error
	})
}

func (rbac *SqlRbac) FindRolesByUser(userId int64) ([]*AuthItem, error) {
	assignment := AuthAssignment{}
	var items []*AuthItem
	err := rbac.db.Model(&assignment).
		Joins("inner join auth_item on auth_assignment.item_name = auth_item.name").
		Where("auth_assignment.user_id = ? and auth_item.type = 1", userId).
		Find(&items).Error
	return items, err
}

func (rbac *SqlRbac) FindChildrenList() ([]*AuthItemChild, error) {
	var children []*AuthItemChild
	err := rbac.db.Find(&children).Error
	return children, err
}

func (rbac *SqlRbac) FindChildrenFormChild(child string) ([]*AuthItemChild, error) {
	var children []*AuthItemChild
	err := rbac.db.Where("child = ?", child).Find(&children).Error
	return children, err
}

func (rbac *SqlRbac) GetItemList(t int32, names []string) ([]*AuthItem, error) {
	var items []*AuthItem
	if len(names) > 0 {
		err := rbac.db.Where("type = ? and name in ?", t, names).Find(&items).Error
		return items, err
	} else {
		err := rbac.db.Where("type = ?", t).Find(&items).Error
		return items, err
	}
}

func (rbac *SqlRbac) FindPermissionsByUser(userId int64) ([]*AuthItem, error) {
	assignment := AuthAssignment{}
	var items []*AuthItem
	err := rbac.db.Model(&assignment).
		Joins("inner join auth_item on auth_assignment.item_name = auth_item.name").
		Where("auth_assignment.user_id = ? and auth_item.type = 2", userId).
		Find(&items).Error
	return items, err
}

func (rbac *SqlRbac) FindAssignmentByUser(userId int64) ([]*AuthAssignment, error) {
	var assignments []*AuthAssignment
	err := rbac.db.Where("user_id = ?", userId).Find(&assignments).Error
	return assignments, err
}

func (rbac *SqlRbac) AddItemChild(itemChild AuthItemChild) error {
	return rbac.db.Create(&itemChild).Error
}

func (rbac *SqlRbac) RemoveChild(parent string, child string) error {
	var itemChild AuthItemChild
	return rbac.db.Where("parent = ? and child = ?", parent, child).Delete(&itemChild).Error
}

func (rbac *SqlRbac) RemoveChildren(parent string) error {
	var itemChild AuthItemChild
	return rbac.db.Where("parent = ?", parent).Delete(&itemChild).Error
}

func (rbac *SqlRbac) HasChild(parent string, child string) bool {
	var itemChild AuthItemChild
	first := rbac.db.Model(&itemChild).Where("parent = ? and child = ?", parent, child).First(&itemChild)
	return first.Error == nil
}

func (rbac *SqlRbac) FindChildren(name string) ([]*AuthItem, error) {
	var items []*AuthItem
	item := AuthItem{}
	err := rbac.db.Model(&item).
		Joins("inner join auth_item_child on auth_item.name = auth_item_child.child").
		Where("auth_item_child.parent = ?", name).Error
	return items, err
}

func (rbac *SqlRbac) Assign(assignment AuthAssignment) error {
	return rbac.db.Create(&assignment).Error
}

func (rbac *SqlRbac) RemoveAssignment(userId int64, name string) error {
	var assignment AuthAssignment
	return rbac.db.Where("user_id = ? and item_name = ?", userId, name).Delete(&assignment).Error
}

func (rbac *SqlRbac) RemoveAllAssignmentByUser(userId int64) error {
	var assignment AuthAssignment
	return rbac.db.Where("user_id = ?", userId).Delete(&assignment).Error
}

func (rbac *SqlRbac) RemoveAllAssignments() error {
	var assignment AuthAssignment
	return rbac.db.Delete(&assignment).Error
}

func (rbac *SqlRbac) GetAssignment(userId int64, name string) (*AuthAssignment, error) {
	var assignments *AuthAssignment
	err := rbac.db.Where("user_id = ? and item_name = ?", userId, name).First(assignments).Error
	return assignments, err
}

func (rbac *SqlRbac) GetAssignmentByItems(name string) ([]*AuthAssignment, error) {
	var assignments []*AuthAssignment
	err := rbac.db.Where("item_name = ?", name).Find(&assignments).Error
	return assignments, err
}

func (rbac *SqlRbac) GetAssignments(userId int64) ([]*AuthAssignment, error) {
	var assignments []*AuthAssignment
	err := rbac.db.Where("user_id = ?", userId).Find(&assignments).Error
	return assignments, err
}

func (rbac *SqlRbac) GetAllAssignment() ([]*AuthAssignment, error) {
	var assignments []*AuthAssignment
	err := rbac.db.Find(&assignments).Error
	return assignments, err
}

func (rbac *SqlRbac) RemoveAll() error {
	return rbac.db.Transaction(func(tx *gorm.DB) error {
		var assignment AuthAssignment
		tx.Delete(&assignment)
		var item AuthItem
		tx.Delete(&item)
		var rule AuthRule
		tx.Delete(&rule)
		return nil
	})
}

func (rbac *SqlRbac) RemoveChildByNames(key string, names []string) error {
	if names != nil && len(names) > 0 {
		var itemChild AuthItemChild
		return rbac.db.Where(key+" in (?)", names).Delete(&itemChild).Error
	}
	return nil
}

func (rbac *SqlRbac) RemoveAssignmentByName(names []string) error {
	if names != nil && len(names) > 0 {
		var assignments AuthAssignment
		return rbac.db.Where("item_name in (?)", names).Delete(&assignments).Error
	}
	return nil
}

func (rbac *SqlRbac) RemoveItemByType(t int32) error {
	var item AuthItem
	return rbac.db.Where("type = ?", t).Delete(&item).Error
}

func (rbac *SqlRbac) RemoveAllRules() error {
	return rbac.db.Transaction(func(tx *gorm.DB) error {
		var item AuthItem
		tx.Model(&item).Update("rule_name", nil)
		var rule AuthRule
		tx.Delete(&rule)
		return nil
	})
}
