package base

import (
	log "github.com/kordar/gologger"
)

type Executor interface {
	Name() string
	Execute(userId interface{}, item Item, params map[string]interface{}) bool
}

type DemoExecutor struct {
}

func (d *DemoExecutor) Name() string {
	return "demo"
}

func (d *DemoExecutor) Execute(userId interface{}, item Item, params map[string]interface{}) bool {
	log.Info("============================================")
	log.Info("==============DEMO===============")
	log.Info("============================================")
	return true
}
