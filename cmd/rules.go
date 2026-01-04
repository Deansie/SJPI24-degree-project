package cmd

// Contract

type Resource struct {
	Kind      string
	Name      string
	Namespace string
	Labels    map[string]string
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

// Verify labels 'app' and 'env' are present

func ruleRequiredLabels(r Resource) []RuleViolation {
	var violations []RuleViolation

	requiredLabels := []string{"app", "env"}

	for _, label := range requiredLabels {
		if r.Labels == nil {
			violations = append(violations, RuleViolation{
				Resource: r,
				Message:  "missing required label: " + label,
			})
			continue
		}

		if _, ok := r.Labels[label]; !ok {
			violations = append(violations, RuleViolation{
				Resource: r,
				Message:  "missing required label: " + label,
			})
		}
	}

	return violations
}

// Rule runner

func ValidateRules(resources []Resource) []RuleViolation {
	var allViolations []RuleViolation

	for _, r := range resources {
		allViolations = append(allViolations, ruleNamespaceRequired(r)...)
		allViolations = append(allViolations, ruleRequiredLabels(r)...)
	}

	return allViolations
}
