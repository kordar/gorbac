package base

// ExecuteManager /******************execute manger*****************************/
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
