package main

import (
	"io"
	"log"
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
	log.Printf("fetching ACL from: %s\n", config.ACL)
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
	log.Printf("Loaded %d ACL rules\n", len(rules))
}
