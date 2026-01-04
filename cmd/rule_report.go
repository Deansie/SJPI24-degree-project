package cmd

import "fmt"

func reportRuleViolations(violations []RuleViolation) {
	for _, v := range violations {
		fmt.Printf("✗ %s %s\n", v.Resource.Kind, v.Resource.Name)
		fmt.Printf("  - %s\n", v.Message)
	}
}
