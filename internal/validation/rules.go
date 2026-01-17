package validation

// Contract

type Container struct {
	Name     string
	Requests map[string]string
	Limits   map[string]string
}
type Resource struct {
	Kind        string
	Name        string
	Namespace   string
	Labels      map[string]string
	Containers  []Container
	Annotations map[string]string
	Spec        Spec
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

// Verify labels 'app' and 'env' are present i kind is Deployment and StatefulSet

func ruleRequiredLabels(r Resource) []RuleViolation {
	var violations []RuleViolation

	requiredLabels := []string{"app", "env"}
	isMandatoryKind := r.Kind == "Deployment" || r.Kind == "StatefulSet" || r.Kind == "DaemonSet"

	if !isMandatoryKind {
		return violations
	}

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

	if r.Kind != "Deployment" && r.Kind != "StatefulSet" {
		return violations
	}

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

func ruleIngressSpecific(r Resource) []RuleViolation {
	var violations []RuleViolation

	if r.Kind != "Ingress" {
		return violations
	}

	// Require ingressClassName for explicit controller
	if r.Spec.IngressClassName == "" {
		violations = append(violations, RuleViolation{
			Resource: r,
			Message:  "Ingress must specify ingressClassName (e.g., 'nginx')",
		})
	}

	// If TLS spec present, require cert-manager annotation
	if len(r.Spec.TLS) > 0 {
		if _, ok := r.Annotations["cert-manager.io/cluster-issuer"]; !ok {
			violations = append(violations, RuleViolation{
				Resource: r,
				Message:  "Ingress with TLS must include cert-manager.io/cluster-issuer annotation for automatic certs",
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
		allViolations = append(allViolations, ruleIngressSpecific(r)...)
	}

	return allViolations
}
