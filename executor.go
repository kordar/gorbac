package gorbac

import (
	"context"
	log "github.com/kordar/gologger"
)

type Executor interface {
	Name() string
	Execute(ctx context.Context, userId interface{}, item Item) bool
}

// ExecuteManager /****************** execute manger *****************************/
var ExecuteManager = container{
	content: make(map[string]Executor),
}

type container struct {
	content map[string]Executor
}

func (container *container) AddExecutor(executor Executor) {
	container.content[executor.Name()] = executor
}

func (container *container) GetExecutor(name string) Executor {
	return container.content[name]
}

type DemoExecutor struct {
}

func (d *DemoExecutor) Name() string {
	return "demo"
}

func (d *DemoExecutor) Execute(ctx context.Context, userId interface{}, item Item) bool {
	log.Info("============================================")
	log.Info("==============DEMO===============")
	log.Info("============================================")
	return true
}
