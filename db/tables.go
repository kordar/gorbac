package db

var tableNames = map[string]string{
	"rule":       "auth_rule",
	"base":       "auth_item",
	"base-child": "auth_item_child",
	"assignment": "auth_assignment",
}

// SetTableName 配置覆盖默认表名
func SetTableName(key string, value string) {
	if tableNames[key] != "" {
		tableNames[key] = value
	}
}
