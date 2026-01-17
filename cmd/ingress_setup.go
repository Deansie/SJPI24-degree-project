package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var (
	ingressName          string
	ingressHost          string
	ingressService       string
	ingressPort          int
	ingressNamespace     string
	ingressOutput        string
	issuerName           string
	enableTLS            bool
	noTLS                bool
	ingressClassName     string
	enableAffinity       bool
	sessionCookieName    string
	sessionCookieExpires string
	sessionCookieMaxAge  string
	enableForceRedirect  bool
	limitRPS             int
	limitConnections     int
	proxyConnectTimeout  int
	proxyReadTimeout     int
	enableCORS           bool
	corsAllowOrigin      string
	corsAllowMethods     string
	proxyBodySize        string
	ingressEnv           string
)

var ingressSetupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Generate Kubernetes Ingress YAML with automatic TLS via cert-manager",
	Long: appLogo + `Generates an Ingress manifest based on provided service and domain. Enables automatic TLS by default using cert-manager.

Key features include NGINX-specific annotations for security, performance, and routing. Use flags to customize:

- TLS: Enabled by default; use --no-tls to disable (not recommended for production).
- Affinity: Enable session stickiness with --affinity; override cookie settings like --session-cookie-name=MYSESSION.
- Rate Limiting: Set RPS/connections with --limit-rps=10 or --limit-connections=5.
- Timeouts: Adjust with --proxy-connect-timeout=30 or --proxy-read-timeout=120.
- CORS: Enable with --enable-cors; customize origins/methods like --cors-allow-origin=https://example.com --cors-allow-methods=GET,POST,PUT.
- Body Size: Increase for uploads with --proxy-body-size=8m.

Examples:
  # Basic with defaults (TLS on, no extras):
  k8s-deploy ingress setup --name myapp --host app.example.com --service myapp-service

  # With affinity and custom cookie (for session persistence):
  k8s-deploy ingress setup --name myapp --host app.example.com --service myapp-service --affinity --session-cookie-name=MYSESSION --session-cookie-expires=86400

  # Security-focused (rate limiting, no TLS for internal use):
  k8s-deploy ingress setup --name myapp --host app.example.com --service myapp-service --no-tls --limit-rps=20 --limit-connections=10

  # Performance tweaks (timeouts, body size for uploads):
  k8s-deploy ingress setup --name myapp --host app.example.com --service myapp-service --proxy-connect-timeout=30 --proxy-read-timeout=180 --proxy-body-size=16m

  # CORS for web apps (custom origins/methods):
  k8s-deploy ingress setup --name myapp --host app.example.com --service myapp-service --enable-cors --cors-allow-origin=https://frontend.com,https://api.com --cors-allow-methods=GET,POST,DELETE

  # Output to file with custom issuer/class:
  k8s-deploy ingress setup --name myapp --host app.example.com --service myapp-service --issuer=letsencrypt-prod --class=nginx --output=ingress.yaml`,
	RunE: runIngressSetup,
}

func init() {
	ingressSetupCmd.Flags().StringVar(&ingressName, "name", "", "Name of the Ingress resource (required)")
	ingressSetupCmd.Flags().StringVar(&ingressHost, "host", "", "Host domain for the Ingress rule (e.g., app.example.com) (required)")
	ingressSetupCmd.Flags().StringVar(&ingressService, "service", "", "Backend service name (required)")
	ingressSetupCmd.Flags().IntVar(&ingressPort, "port", 80, "Backend service port")
	ingressSetupCmd.Flags().StringVar(&ingressNamespace, "namespace", "default", "Namespace for the Ingress")
	ingressSetupCmd.Flags().BoolVar(&enableTLS, "tls", true, "Enable automatic TLS with cert-manager (default: true)")
	ingressSetupCmd.Flags().BoolVar(&noTLS, "no-tls", false, "Disable TLS (not recommended for production; use for internal/dev only)")
	ingressSetupCmd.Flags().StringVar(&ingressOutput, "output", "", "Output file path (default: stdout)")
	ingressSetupCmd.Flags().StringVar(&issuerName, "issuer", "letsencrypt", "cert-manager cluster-issuer name (used if --tls is enabled)")
	ingressSetupCmd.Flags().StringVar(&ingressClassName, "class", "nginx", "Ingress class name (e.g., nginx)")
	ingressSetupCmd.Flags().BoolVar(&enableAffinity, "affinity", false, "Enable cookie-based session affinity")
	ingressSetupCmd.Flags().StringVar(&sessionCookieName, "session-cookie-name", "INGRESSCOOKIE", "Session cookie name for affinity")
	ingressSetupCmd.Flags().StringVar(&sessionCookieExpires, "session-cookie-expires", "172800", "Session cookie expires time (seconds)")
	ingressSetupCmd.Flags().StringVar(&sessionCookieMaxAge, "session-cookie-max-age", "172800", "Session cookie max age (seconds)")
	ingressSetupCmd.Flags().BoolVar(&enableForceRedirect, "force-redirect", true, "Enable force SSL redirect")
	ingressSetupCmd.Flags().IntVar(&limitRPS, "limit-rps", 0, "Max requests per second per IP (0 to disable)")
	ingressSetupCmd.Flags().IntVar(&limitConnections, "limit-connections", 0, "Max concurrent connections per IP (0 to disable)")
	ingressSetupCmd.Flags().IntVar(&proxyConnectTimeout, "proxy-connect-timeout", 60, "Backend connect timeout in seconds")
	ingressSetupCmd.Flags().IntVar(&proxyReadTimeout, "proxy-read-timeout", 60, "Backend read timeout in seconds")
	ingressSetupCmd.Flags().BoolVar(&enableCORS, "enable-cors", false, "Enable CORS headers")
	ingressSetupCmd.Flags().StringVar(&corsAllowOrigin, "cors-allow-origin", "*", "Allowed CORS origins (comma-separated)")
	ingressSetupCmd.Flags().StringVar(&corsAllowMethods, "cors-allow-methods", "GET,POST,OPTIONS", "Allowed CORS methods")
	ingressSetupCmd.Flags().StringVar(&proxyBodySize, "proxy-body-size", "1m", "Max proxy body size (e.g., 8m for larger uploads)")
	ingressSetupCmd.Flags().StringVar(&ingressEnv, "env", "", "Environment label for Ingress (e.g., prod, staging, dev; optional)")

	_ = ingressSetupCmd.MarkFlagRequired("name")
	_ = ingressSetupCmd.MarkFlagRequired("host")
	_ = ingressSetupCmd.MarkFlagRequired("service")
}

type Ingress struct {
	APIVersion string      `yaml:"apiVersion"`
	Kind       string      `yaml:"kind"`
	Metadata   Metadata    `yaml:"metadata"`
	Spec       IngressSpec `yaml:"spec"`
}

type Metadata struct {
	Name        string            `yaml:"name"`
	Namespace   string            `yaml:"namespace,omitempty"`
	Labels      map[string]string `yaml:"labels,omitempty"`
	Annotations map[string]string `yaml:"annotations,omitempty"`
}

type IngressSpec struct {
	IngressClassName *string    `yaml:"ingressClassName,omitempty"`
	TLS              []TLSEntry `yaml:"tls,omitempty"`
	Rules            []Rule     `yaml:"rules"`
}

type TLSEntry struct {
	Hosts      []string `yaml:"hosts"`
	SecretName string   `yaml:"secretName"`
}

type Rule struct {
	Host string `yaml:"host"`
	HTTP HTTP   `yaml:"http"`
}

type HTTP struct {
	Paths []Path `yaml:"paths"`
}

type Path struct {
	Path     string  `yaml:"path"`
	PathType string  `yaml:"pathType"`
	Backend  Backend `yaml:"backend"`
}

type Backend struct {
	Service Service `yaml:"service"`
}

type Service struct {
	Name string `yaml:"name"`
	Port Port   `yaml:"port"`
}

type Port struct {
	Number int `yaml:"number"`
}

func runIngressSetup(cmd *cobra.Command, args []string) error {
	if ingressName == "" || ingressHost == "" || ingressService == "" {
		return fmt.Errorf("required flags: --name, --host, --service")
	}

	useTLS := enableTLS && !noTLS
	if !useTLS {
		fmt.Println("Warning: TLS is disabled. Consider enabling for security.")
	}

	annotations := make(map[string]string)
	annotations["nginx.ingress.kubernetes.io/rewrite-target"] = "/"
	annotations["nginx.ingress.kubernetes.io/proxy-body-size"] = proxyBodySize
	annotations["nginx.ingress.kubernetes.io/proxy-connect-timeout"] = fmt.Sprintf("%d", proxyConnectTimeout)
	annotations["nginx.ingress.kubernetes.io/proxy-read-timeout"] = fmt.Sprintf("%d", proxyReadTimeout)
	if useTLS {
		annotations["cert-manager.io/cluster-issuer"] = issuerName
		annotations["nginx.ingress.kubernetes.io/ssl-redirect"] = "true"
		if enableForceRedirect {
			annotations["nginx.ingress.kubernetes.io/force-ssl-redirect"] = "true"
		}
	}
	if enableAffinity {
		annotations["nginx.ingress.kubernetes.io/affinity"] = "cookie"
		annotations["nginx.ingress.kubernetes.io/affinity-mode"] = "persistent"
		annotations["nginx.ingress.kubernetes.io/session-cookie-name"] = sessionCookieName
		annotations["nginx.ingress.kubernetes.io/session-cookie-expires"] = sessionCookieExpires
		annotations["nginx.ingress.kubernetes.io/session-cookie-max-age"] = sessionCookieMaxAge
	}
	if limitRPS > 0 {
		annotations["nginx.ingress.kubernetes.io/limit-rps"] = fmt.Sprintf("%d", limitRPS)
	}
	if limitConnections > 0 {
		annotations["nginx.ingress.kubernetes.io/limit-connections"] = fmt.Sprintf("%d", limitConnections)
	}
	if enableCORS {
		annotations["nginx.ingress.kubernetes.io/enable-cors"] = "true"
		annotations["nginx.ingress.kubernetes.io/cors-allow-origin"] = corsAllowOrigin
		annotations["nginx.ingress.kubernetes.io/cors-allow-methods"] = corsAllowMethods
	}

	labels := make(map[string]string)
	if ingressEnv != "" {
		labels["app"] = ingressName
		labels["env"] = ingressEnv
	}

	ing := Ingress{
		APIVersion: "networking.k8s.io/v1",
		Kind:       "Ingress",
		Metadata: Metadata{
			Name:        ingressName,
			Namespace:   ingressNamespace,
			Labels:      labels,
			Annotations: annotations,
		},
		Spec: IngressSpec{
			IngressClassName: &ingressClassName,
			Rules: []Rule{
				{
					Host: ingressHost,
					HTTP: HTTP{
						Paths: []Path{
							{
								Path:     "/",
								PathType: "Prefix",
								Backend: Backend{
									Service: Service{
										Name: ingressService,
										Port: Port{
											Number: ingressPort,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	if useTLS {
		ing.Spec.TLS = []TLSEntry{
			{
				Hosts:      []string{ingressHost},
				SecretName: ingressName + "-tls",
			},
		}
	}

	yamlData, err := yaml.Marshal(ing)
	if err != nil {
		return fmt.Errorf("failed to generate YAML: %w", err)
	}

	if ingressOutput != "" {
		if err := os.WriteFile(ingressOutput, yamlData, 0644); err != nil {
			return fmt.Errorf("failed to write to file: %w", err)
		}
		fmt.Printf("Ingress YAML written to %s\n", ingressOutput)
	} else {
		fmt.Println(string(yamlData))
	}

	return nil
}
