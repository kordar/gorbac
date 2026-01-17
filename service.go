package gorbac

import (
	"fmt"
	"time"
)

type RbacService struct {
	mgr AuthManager
}

func NewRbacService(repos AuthRepository, cache bool) *RbacService {
	manager := NewDefaultManager(repos, cache)
	return NewRbacServiceWithManager(manager)
}

func NewRbacServiceWithManager(mgr AuthManager) *RbacService {
	return &RbacService{mgr: mgr}
}

func (s RbacService) GetAuthManager() AuthManager {
	return s.mgr
}

// ---------------------- Roles ---------------------------

func (s RbacService) Roles() []*Role {
	return s.mgr.GetRoles()
}

func (s RbacService) GetRolesByUser(userId interface{}) []*Role {
	return s.mgr.GetRolesByUser(userId)
}

func (s RbacService) AddRole(name string, description string, ruleName string) bool {
	role := NewRole(name, description, ruleName, "", time.Now(), time.Now())
	return s.mgr.Add(role)
}

func (s RbacService) UpdateRole(name string, newName string, description string, ruleName string) bool {
	role := NewRole(newName, description, ruleName, "", time.Now(), time.Now())
	return s.mgr.Update(name, role)
}

func (s RbacService) DeleteRole(name string) bool {
	role := s.mgr.GetRole(name)
	if role == nil {
		return false
	}
	return s.mgr.Remove(role)
}

// ---------------------- Permissions ---------------------------

func (s RbacService) Permissions() []*Permission {
	return s.mgr.GetPermissions()
}

func (s RbacService) GetPermissionsByUser(userId interface{}) []*Permission {
	return s.mgr.GetPermissionsByUser(userId)
}

func (s RbacService) AddPermission(name string, description string, ruleName string) bool {
	permission := NewPermission(name, description, ruleName, "", time.Now(), time.Now())
	return s.mgr.Add(permission)
}

func (s RbacService) UpdatePermission(name string, newName string, description string, ruleName string) bool {
	permission := NewPermission(newName, description, ruleName, "", time.Now(), time.Now())
	return s.mgr.Update(name, permission)
}

func (s RbacService) DeletePermission(name string) bool {
	permission := s.mgr.GetPermission(name)
	if permission == nil {
		return false
	}
	return s.mgr.Remove(permission)
}

func (s RbacService) CleanChildren(parent string) bool {
	item := s.mgr.GetItem(parent)
	if item == nil {
		return false
	}
	return s.mgr.RemoveChildren(item)
}

func (s RbacService) AssignChildren(parent string, children ...string) error {
	if len(children) == 0 {
		return nil
	}

	role := s.mgr.GetRole(parent)
	if role == nil {
		return fmt.Errorf("role %s not found", parent)
	}

	for _, ss := range children {
		item := s.mgr.GetItem(ss)
		if item != nil {
			_ = s.mgr.AddChild(role, item)
		}
	}

	return nil
}

// ---------------------- Rule ---------------------------

func (s RbacService) Rules() []*Rule {
	return s.mgr.GetRules()
}

func (s RbacService) AddRule(name string, executeName string) bool {
	rule := NewRule(name, executeName, time.Now(), time.Now())
	return s.mgr.AddRule(*rule)
}

func (s RbacService) UpdateRule(name string, newName string, executeName string) bool {
	rule := NewRule(newName, executeName, time.Now(), time.Now())
	return s.mgr.UpdateRule(name, *rule)
}

func (s RbacService) DeleteRule(name string) bool {
	rule := s.mgr.GetRule(name)
	if rule == nil {
		return false
	}
	return s.mgr.RemoveRule(*rule)
}

// ------------------- Assign --------------------------

func (s RbacService) Assign(userId interface{}, name string) bool {
	item := s.mgr.GetItem(name)
	if item == nil {
		return false
	}
	if err := s.mgr.Assign(item, userId); err == nil {
		return true
	} else {
		return false
	}
}

func (s RbacService) CleanAssigns(userId interface{}) {
	_ = s.mgr.RemoveAllAssignmentByUser(userId)
}

func (s RbacService) Assigns(userId interface{}, names ...string) {
	s.mgr.Assigns(userId, names...)
	// for _, name := range names {
	// 	item := s.mgr.GetItem(name)
	// 	s.mgr.Assign(item, userId)
	// }
}

func (s RbacService) GetChildren(name string) []Item {
	return s.mgr.GetChildren(name)
}
