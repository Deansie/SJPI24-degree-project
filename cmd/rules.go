package cmd

// Contract

type Resource struct {
	Kind      string
	Name      string
	Namespace string
}

type RuleViolation struct {
	Resource Resource
	Message  string
}

// Verify namespace is set and explicitly not to 'default'

func ruleNamespaceRequired(r Resource) []RuleViolation {
	var violations []RuleViolation

	if r.Namespace == "" || r.Namespace == "default" {
		violations = append(violations, RuleViolation{
			Resource: r,
			Message:  "resource must explicitly define a non-default namespace ('default' is not allowed)",
		})
	}

	return violations
}

// Rule runner

func ValidateRules(resources []Resource) []RuleViolation {
	var allViolations []RuleViolation

	for _, r := range resources {
		allViolations = append(allViolations, ruleNamespaceRequired(r)...)
	}

	return allViolations
}
