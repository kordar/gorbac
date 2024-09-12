package gorbac

type DefaultCache struct {
	enable bool
	// all auth items (name => Item)
	items map[string]Item
	// all auth rules (name => Rule)
	rules map[string]*Rule
	// auth item parent-child relationships (childName => list of parents)
	parents map[string][]string
}

func NewDefaultCache(cache bool) *DefaultCache {
	return &DefaultCache{enable: cache}
}

func (manager *DefaultCache) invalidateCache() {
	if manager.enable {
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
	if manager.rules == nil {
		return nil
	}
	rule := f(name)
	if rule == nil {
		return nil
	}
	manager.rules[name] = rule
	return manager.rules[name]
}

func (manager *DefaultCache) GetRules(f func() []*Rule) []*Rule {
	if manager.rules == nil {
		return make([]*Rule, 0)
	}

	if len(manager.rules) == 0 {
		rules := f()
		for _, rule := range rules {
			manager.rules[rule.Name] = rule
		}
		return rules
	}

	rules := make([]*Rule, 0, len(manager.rules))
	for _, rule := range manager.rules {
		rules = append(rules, rule)
	}
	return rules
}
