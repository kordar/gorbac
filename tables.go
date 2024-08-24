package gorbac

var tableNames = map[string]string{
	"rule":       "auth_rule",
	"item":       "auth_item",
	"item-child": "auth_item_child",
	"assignment": "auth_assignment",
}

// SetTableName 配置覆盖默认表名
func SetTableName(key string, value string) {
	if tableNames[key] != "" {
		tableNames[key] = value
	}
}

func GetTableName(key string) string {
	return tableNames[key]
}
