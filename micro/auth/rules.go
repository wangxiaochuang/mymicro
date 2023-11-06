package auth

import (
	"fmt"
	"sort"
	"strings"
)

func Verify(rules []*Rule, acc *Account, res *Resource) error {
	validTypes := []string{"*", res.Type}
	validNames := []string{"*", res.Name}
	validEndpoints := []string{"*", res.Endpoint}
	// /a/b/c /a/* /a/b/* /a/b/c/*
	if comps := strings.Split(res.Endpoint, "/"); len(comps) > 1 {
		for i := 1; i < len(comps)+1; i++ {
			wildcard := fmt.Sprintf("%v/*", strings.Join(comps[0:i], "/"))
			validEndpoints = append(validEndpoints, wildcard)
		}
	}
	filteredRules := make([]*Rule, 0)
	for _, rule := range rules {
		if !include(validTypes, rule.Resource.Type) {
			continue
		}
		if !include(validNames, rule.Resource.Name) {
			continue
		}
		if !include(validEndpoints, rule.Resource.Endpoint) {
			continue
		}
		filteredRules = append(filteredRules, rule)
	}

	sort.SliceStable(filteredRules, func(i, j int) bool {
		return filteredRules[i].Priority > filteredRules[j].Priority
	})

	for _, rule := range filteredRules {
		if rule.Scope == ScopePublic && rule.Access == AccessDenied {
			return ErrForbidden
		} else if rule.Scope == ScopePublic && rule.Access == AccessGranted {
			return nil
		}

		if acc == nil {
			continue
		}

		if rule.Scope == ScopeAccount && rule.Access == AccessDenied {
			return ErrForbidden
		} else if rule.Scope == ScopeAccount && rule.Access == AccessGranted {
			return nil
		}

		if include(acc.Scopes, rule.Scope) && rule.Access == AccessDenied {
			return ErrForbidden
		} else if include(acc.Scopes, rule.Scope) && rule.Access == AccessGranted {
			return nil
		}
	}

	return ErrForbidden
}

func include(slice []string, val string) bool {
	for _, s := range slice {
		if strings.EqualFold(s, val) {
			return true
		}
	}
	return false
}
