package main

import (
	"io"
	"net/http"
	"time"

	"github.com/pmezard/adblock/adblock"
)

// WARNING: global variable, consider refactoring to avoid it
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

func fetchACL() io.ReadCloser {
	l.Infof("fetching ACL from: %s", config.ACL)
	client := &http.Client{
		Timeout: config.Timeout * time.Second,
	}

	resp, err := client.Get(config.ACL)
	if err != nil {
		panic(err)
	}
	return resp.Body

}

// NOTE: runs in goroutine, therefore panic() justfified?
func loadACL() {
	matcher = adblock.NewMatcher()

	rules, err := adblock.ParseRules(fetchACL())
	if err != nil {
		panic(err)
	}
	for idx, rule := range rules {
		err = matcher.AddRule(rule, idx)
		if err != nil {
			panic(err)
		}
	}
	l.Infof("Loaded %d ACL rules", len(rules))
}
