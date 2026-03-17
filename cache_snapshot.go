package gorbac

import "time"

type cacheSnapshot struct {
	Version int                 `json:"version"`
	Items   []cacheSnapshotItem `json:"items"`
	Rules   []cacheSnapshotRule `json:"rules"`
	Parents map[string][]string `json:"parents"`
}

type cacheSnapshotItem struct {
	Name        string    `json:"name"`
	Type        int32     `json:"type"`
	Description string    `json:"description"`
	RuleName    string    `json:"rule_name"`
	ExecuteName string    `json:"execute_name"`
	CreateTime  time.Time `json:"create_time"`
	UpdateTime  time.Time `json:"update_time"`
}

type cacheSnapshotRule struct {
	Name        string    `json:"name"`
	ExecuteName string    `json:"execute_name"`
	CreateTime  time.Time `json:"create_time"`
	UpdateTime  time.Time `json:"update_time"`
}

