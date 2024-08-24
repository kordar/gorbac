package gorbac

import (
	"errors"
	"fmt"
	"time"
)

type RbacService struct {
	mgr *DefaultManager
}

func NewRbacService(mgr AuthRepository, cache bool) *RbacService {
	return &RbacService{mgr: NewDefaultManager(mgr, cache)}
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

func (s RbacService) AssignPermission(parent string, child string) error {
	role := s.mgr.GetRole(parent)
	if role == nil {
		return errors.New(fmt.Sprintf("role %s not found", parent))
	}
	permission := s.mgr.GetPermission(child)
	if permission == nil {
		return errors.New(fmt.Sprintf("permission %s not found", child))
	}
	return s.mgr.AddChild(role, permission)
}

func (s RbacService) AssignRole(parent string, child string) error {
	role := s.mgr.GetRole(parent)
	if role == nil {
		return errors.New(fmt.Sprintf("role %s not found", parent))
	}
	role2 := s.mgr.GetRole(child)
	if role2 == nil {
		return errors.New(fmt.Sprintf("role %s not found", child))
	}
	return s.mgr.AddChild(role, role2)
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

func (s RbacService) AssignRoleToUser(name string, userId interface{}) bool {
	role := s.mgr.GetRole(name)
	if role == nil {
		return false
	}
	if err := s.mgr.Assign(role, userId); err == nil {
		return true
	} else {
		return false
	}
}

func (s RbacService) AssignItemToUser(roles []string, permissions []string, userId interface{}) {
	_ = s.mgr.RemoveAllAssignmentByUser(userId)
	for _, pname := range permissions {
		permission := s.mgr.GetPermission(pname)
		if permission == nil {
			continue
		}
		s.mgr.Assign(permission, userId)
	}
	for _, rname := range roles {
		role := s.mgr.GetRole(rname)
		if role == nil {
			continue
		}
		s.mgr.Assign(role, userId)
	}
}

func (s RbacService) GetChildren(name string) []Item {
	return s.mgr.GetChildren(name)
}

func (s RbacService) AddPermissionChildren(parentName string, children []string) {
	role := s.mgr.GetRole(parentName)
	if role == nil {
		return
	}
	s.mgr.RemoveChildren(role)
	for _, name := range children {
		permission := s.mgr.GetPermission(name)
		_ = s.mgr.AddChild(role, permission)
	}
}
