```go
package gorbac

import (
    "github.com/kordar/godbutil"
    "github.com/kordar/goutil"
    "log"
)

func init() {
    goutil.ConfigInit("conf/dev.ini") // 初始化配置
    godbutil.GetSqlitePool().InitDataPool("sys")
}

func main() {
    ExecuteManager.AddExecutor(&DemoExecutor{})
    db := godbutil.GetSqlitePool().Handler("sys")
    /*_ = db.AutoMigrate(&models.AuthRule{})
    _ = db.AutoMigrate(&models.AuthItem{})
    _ = db.AutoMigrate(&models.AuthItemChild{})
    _ = db.AutoMigrate(&models.AuthAssignment{})*/
    mapper := NewSqlRbac(db)
    dbManager := NewDbManager(mapper, true)
    //role := executor.NewRole("aa", "", "", "", time.Now())
    //add := dbManager.Add(role)
    //log.Println(add)
    //roles := dbManager.GetRoles()
    //fmt.Println(fmt.Printf("roles = %v", roles))
    //permissions := dbManager.GetPermissions()
    //log.Println(permissions)
    //role := dbManager.GetRole("role1")
    permission := dbManager.GetPermission("permission1")
    //rule := executor.NewRule("rule2", "demo", time.Now())
    //addRule := dbManager.AddRule(rule)
    //log.Println("add rule", addRule)
    // permission.RuleName = rule.Name
    //dbManager.Update("permission1", permission)

	//rule.Name = "demo"
	//dbManager.UpdateRule("rule2", rule)
	/*
		err := dbManager.AddChild(role, permission)*/
	//assign := dbManager.Assign(role, 1001)
	//child := dbManager.HasChild(role, permission)

	access := dbManager.CheckAccess(1001, permission.GetName(), map[string]interface{}{"aa": "cc"})

	log.Println(access)
}
```