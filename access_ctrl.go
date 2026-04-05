package main

import "github.com/pmezard/adblock/adblock"

func ACLCheck(host string) bool {
	if host == "" {
		return false
	}
	matcher := loadACL()
	return matcher.Match(host)
}

func fetchACL() {
	// config.ACL

}

func loadACL() adblock.Matcher {
	matcher := adblock.NewMatcher()
	file, err := adblock.OpenFile("easylist.txt")
	if err != nil {
		panic(err)
	}

	rules, err := adblock.ParseRules(file)
	if err != nil {
		panic(err)
	}
	for idx, rule := range rules {
		err = matcher.AddRule(rule, idx)
		if err != nil {
			panic(err)
		}
	}
	return matcher
}
