package base

import "time"

type Rule struct {

	/**
	 * string name of the rule
	 */
	Name string `json:"name"`

	/**
	 * 执行函数
	 */
	ExecuteName string `json:"execute_name"`

	/**
	 * int UNIX timestamp representing the rule creation time
	 */
	CreateTime time.Time `json:"create_time"`

	/**
	 * int UNIX timestamp representing the rule updating time
	 */
	UpdateTime time.Time `json:"update_time"`
}

func NewRule(name string, executeName string, createTime time.Time, updateTime time.Time) *Rule {
	return &Rule{Name: name, ExecuteName: executeName, CreateTime: createTime, UpdateTime: updateTime}
}

func (rule *Rule) SetName(name string) {
	rule.Name = name
}

func (rule *Rule) SetExecutor(executor Executor) {
	rule.ExecuteName = executor.Name()
}

func (rule *Rule) GetExecutor() Executor {
	return ExecuteManager.GetExecutor(rule.ExecuteName)
}
