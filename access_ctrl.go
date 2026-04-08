package main

import (
	"io"
	"log"
	"net/http"
	"time"

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

func fetchACL() io.ReadCloser {
	log.Printf("fetching ACL from: %s\n", config.ACL)
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(config.ACL)
	if err != nil {
		panic(err)
	}
	return resp.Body

}

func loadACL() {
	matcher = adblock.NewMatcher()
	// file, err := os.Open("ACL.txt")
	// if err != nil {
	// 	panic(err)
	// }

	// rules, err := adblock.ParseRules(file)

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
