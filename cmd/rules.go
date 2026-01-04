package cmd

// Contract

type Container struct {
	Name     string
	Requests map[string]string
	Limits   map[string]string
}
type Resource struct {
	Kind       string
	Name       string
	Namespace  string
	Labels     map[string]string
	Containers []Container
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

// Verify resource requests and limits for container are present

func ruleContainerResources(r Resource) []RuleViolation {
	var violations []RuleViolation

	for _, c := range r.Containers {
		if c.Requests == nil || c.Limits == nil {
			violations = append(violations, RuleViolation{
				Resource: r,
				Message:  "container \"" + c.Name + "\" must define cpu and memory requests and limits",
			})
			continue
		}

		if _, ok := c.Requests["cpu"]; !ok {
			violations = append(violations, RuleViolation{
				Resource: r,
				Message:  "container \"" + c.Name + "\" is missing cpu request",
			})
		}

		if _, ok := c.Requests["memory"]; !ok {
			violations = append(violations, RuleViolation{
				Resource: r,
				Message:  "container \"" + c.Name + "\" is missing memory request",
			})
		}

		if _, ok := c.Limits["cpu"]; !ok {
			violations = append(violations, RuleViolation{
				Resource: r,
				Message:  "container \"" + c.Name + "\" is missing cpu limit",
			})
		}

		if _, ok := c.Limits["memory"]; !ok {
			violations = append(violations, RuleViolation{
				Resource: r,
				Message:  "container \"" + c.Name + "\" is missing memory limit",
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
		allViolations = append(allViolations, ruleContainerResources(r)...)
	}

	return allViolations
}
