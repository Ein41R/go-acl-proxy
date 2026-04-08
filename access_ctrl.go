package main

import (
	"os"

	"github.com/pmezard/adblock/adblock"
)

var matcher *adblock.RuleMatcher

func ACLCheck(host string) bool {
	if host == "" {
		return false
	}
	req := &adblock.Request{URL: host}
	match, _, err := matcher.Match(req)
	if err != nil {
		panic(err)
	}
	return match
}

func fetchACL() {
	// config.ACL

}

func loadACL() {
	matcher = adblock.NewMatcher()
	file, err := os.Open("https://easylist.to/easylist/easylist.txt")
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
}
