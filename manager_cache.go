package gorbac

type DefaultCache struct {
	cache bool
	// all auth items (name => Item)
	items map[string]Item
	// all auth rules (name => Rule)
	rules map[string]*Rule
	// auth item parent-child relationships (childName => list of parents)
	parents map[string][]string
}

func NewDefaultCache(cache bool) *DefaultCache {
	return &DefaultCache{cache: cache}
}

func (manager *DefaultCache) invalidateCache() {
	if manager.cache {
		manager.items = make(map[string]Item)
		manager.rules = make(map[string]*Rule)
		manager.parents = make(map[string][]string)
	}
}

func (manager *DefaultCache) refreshInvalidateCache(operator bool) bool {
	if operator {
		manager.invalidateCache()
		return true
	}
	return false
}

func (manager *DefaultCache) GetItem(name string, f func(n string) Item) Item {
	if name == "" {
		return nil
	}
	if manager.items != nil {
		return manager.items[name]
	}
	return f(name)
}

func (manager *DefaultCache) GetRule(name string, f func(n string) *Rule) *Rule {
	if manager.rules != nil {
		return manager.rules[name]
	}
	return f(name)
}

func (manager *DefaultCache) GetRules(f func() []*Rule) []*Rule {
	if manager.rules != nil {
		rules := make([]*Rule, 0, len(manager.rules))
		for _, rule := range manager.rules {
			rules = append(rules, rule)
		}
		return rules
	}
	return f()
}
